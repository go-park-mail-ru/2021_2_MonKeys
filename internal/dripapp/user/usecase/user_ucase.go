package usecase

import (
	"bytes"
	"context"
	"dripapp/configs"
	"dripapp/internal/dripapp/models"
	_sessionModels "dripapp/internal/microservices/auth/models"
	"dripapp/internal/pkg/hasher"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const (
	FAKE       = "ФЭЙК"
	AGGRESSION = "АГРЕССИЯ"
	SCAM       = "СКАМ"
	UNDERAGE   = "НЕСОВЕРШЕННОЛЕТНИЙ"
)

type userUsecase struct {
	UserRepo       models.UserRepository
	Session        _sessionModels.SessionRepository
	File           models.FileRepository
	contextTimeout time.Duration
	hub            *models.Hub
}

func NewUserUsecase(
	ur models.UserRepository,
	fileManager models.FileRepository,
	timeout time.Duration,
	hub *models.Hub) models.UserUsecase {
	return &userUsecase{
		UserRepo:       ur,
		File:           fileManager,
		contextTimeout: timeout,
		hub:            hub,
	}
}

func (h *userUsecase) CurrentUser(c context.Context) (models.User, error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUser, ok := ctx.Value(configs.ContextUser).(models.User)
	fmt.Println(currentUser, ok)
	if !ok {
		return models.User{}, models.ErrContextNilError
	}

	return currentUser, nil
}

func (h *userUsecase) EditProfile(c context.Context, newUserData models.User) (updatedUser models.User, err error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUser, ok := ctx.Value(configs.ContextUser).(models.User)
	if !ok {
		return models.User{}, models.ErrContextNilError
	}

	newUserData.ID = currentUser.ID
	newUserData.Email = currentUser.Email

	updatedUser, err = h.UserRepo.UpdateUser(c, newUserData)
	if err != nil {
		return models.User{}, err
	}

	return updatedUser, nil
}

func (h *userUsecase) AddPhoto(c context.Context, photo io.Reader, fileName string) (models.Photo, error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUser, ok := ctx.Value(configs.ContextUser).(models.User)
	if !ok {
		return models.Photo{}, models.ErrContextNilError
	}

	photoPath, err := h.File.SaveUserPhoto(currentUser, photo, fileName)
	if err != nil {
		return models.Photo{}, err
	}

	currentUser.AddNewPhoto(photoPath)

	err = h.UserRepo.UpdateImgs(c, currentUser.ID, currentUser.Imgs)
	if err != nil {
		return models.Photo{}, err
	}

	return models.Photo{Path: photoPath}, nil
}

func (h *userUsecase) DeletePhoto(c context.Context, photo models.Photo) error {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUser, ok := ctx.Value(configs.ContextUser).(models.User)
	if !ok {
		return models.ErrContextNilError
	}

	err := currentUser.DeletePhoto(photo)
	if err != nil {
		return err
	}

	err = h.UserRepo.UpdateImgs(c, currentUser.ID, currentUser.Imgs)
	if err != nil {
		return err
	}

	err = h.File.Delete(photo.Path)
	if err != nil {
		return err
	}

	return nil
}

func (h *userUsecase) Login(c context.Context, logUserData models.LoginUser) (models.User, error) {
	identifiableUser, err := h.UserRepo.GetUser(c, logUserData.Email)
	if err != nil {
		return models.User{}, err
	}

	if !hasher.CheckWithHash(identifiableUser.Password, logUserData.Password) {
		return models.User{}, models.ErrMismatch
	}

	return identifiableUser, nil
}

func (h *userUsecase) Signup(c context.Context, logUserData models.LoginUser) (models.User, error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	identifiableUser, _ := h.UserRepo.GetUser(ctx, logUserData.Email)
	if len(identifiableUser.Email) != 0 {
		return models.User{}, models.ErrEmailAlreadyExists
	}

	logUserData.Password = hasher.HashAndSalt(nil, logUserData.Password)

	user, err := h.UserRepo.CreateUser(c, logUserData)
	if err != nil {
		return models.User{}, err
	}

	err = h.File.CreateFoldersForNewUser(user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (h *userUsecase) NextUser(c context.Context) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUser, ok := ctx.Value(configs.ContextUser).(models.User)
	if !ok {
		return nil, models.ErrContextNilError
	}

	nextUsers, err := h.UserRepo.GetNextUserForSwipe(ctx, currentUser)
	if err != nil {
		return nil, err
	}

	return nextUsers, nil
}

func (h *userUsecase) GetAllTags(c context.Context) (models.Tags, error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	allTags, err := h.UserRepo.GetTags(ctx)
	if err != nil {
		return models.Tags{}, err
	}
	var respTag models.Tag
	var currentAllTags = make(map[uint64]models.Tag)
	var respAllTags models.Tags
	counter := 0

	for _, value := range allTags {
		respTag.TagName = value
		currentAllTags[uint64(counter)] = respTag
		counter++
	}

	respAllTags.AllTags = currentAllTags
	respAllTags.Count = uint64(counter)

	return respAllTags, nil
}

func (h *userUsecase) UsersMatches(c context.Context) (models.Matches, error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUser, ok := ctx.Value(configs.ContextUser).(models.User)
	if !ok {
		return models.Matches{}, models.ErrContextNilError
	}

	// find matches
	mathes, err := h.UserRepo.GetUsersMatches(ctx, currentUser.ID)
	if err != nil {
		return models.Matches{}, err
	}

	// count
	counter := 0
	var allMathesMap = make(map[uint64]models.User)
	for _, value := range mathes {
		allMathesMap[uint64(counter)] = value
		counter++
	}

	var allMatches models.Matches
	allMatches.AllUsers = allMathesMap
	allMatches.Count = strconv.Itoa(counter)

	return allMatches, nil
}

func (h *userUsecase) UsersMatchesWithSearching(c context.Context, searchData models.Search) (models.Matches, error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUser, ok := ctx.Value(configs.ContextUser).(models.User)
	if !ok {
		return models.Matches{}, models.ErrContextNilError
	}

	// find matches
	mathes, err := h.UserRepo.GetUsersMatchesWithSearching(ctx, currentUser.ID, searchData.SearchingTmpl)
	if err != nil {
		return models.Matches{}, err
	}

	// count
	counter := 0
	var allMathesMap = make(map[uint64]models.User)
	for _, value := range mathes {
		allMathesMap[uint64(counter)] = value
		counter++
	}

	var allMatches models.Matches
	allMatches.AllUsers = allMathesMap
	allMatches.Count = strconv.Itoa(counter)

	return allMatches, nil
}

func (h *userUsecase) Reaction(c context.Context, reactionData models.UserReaction) (models.Match, error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUser, ok := ctx.Value(configs.ContextUser).(models.User)
	if !ok {
		return models.Match{}, models.ErrContextNilError
	}

	// added reaction in db
	err := h.UserRepo.AddReaction(ctx, currentUser.ID, reactionData.Id, reactionData.Reaction)
	if err != nil {
		return models.Match{}, err
	}

	// no match if dislike
	var currMath models.Match
	currMath.Match = false
	if reactionData.Reaction != 1 {
		return currMath, nil
	}

	// get users who liked current user
	var likes []uint64
	likes, err = h.UserRepo.GetLikes(ctx, currentUser.ID)
	if err != nil {
		return models.Match{}, err
	}

	for _, value := range likes {
		if value == reactionData.Id {
			currMath.Match = true
			err = h.UserRepo.DeleteReaction(ctx, currentUser.ID, reactionData.Id)
			if err != nil {
				return models.Match{}, err
			}
			err = h.UserRepo.AddMatch(ctx, currentUser.ID, reactionData.Id)
			if err != nil {
				return models.Match{}, err
			}
		}
	}

	// notifications
	if currMath.Match {
		h.hub.NotifyAboutMatchWith(reactionData.Id, currentUser)
	}

	return currMath, nil
}

func (h *userUsecase) ClientHandler(c context.Context, notifications models.Notifications) error {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUser, ok := ctx.Value(configs.ContextUser).(models.User)
	if !ok {
		return models.ErrContextNilError
	}

	models.NewClient(currentUser, h.hub, notifications)

	return nil
}

func (h *userUsecase) UserLikes(c context.Context) (models.Likes, error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUser, ok := ctx.Value(configs.ContextUser).(models.User)
	if !ok {
		return models.Likes{}, models.ErrContextNilError
	}

	// find likes
	likes, err := h.UserRepo.GetUsersLikes(ctx, currentUser.ID)
	if err != nil {
		return models.Likes{}, err
	}

	// count
	counter := 0
	var allMathesMap = make(map[uint64]models.User)
	for _, value := range likes {
		allMathesMap[uint64(counter)] = value
		counter++
	}

	var allLikes models.Likes
	allLikes.AllUsers = allMathesMap
	allLikes.Count = strconv.Itoa(counter)

	return allLikes, nil
}

func (h *userUsecase) GetAllReports(c context.Context) (models.Reports, error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	allReports, err := h.UserRepo.GetReports(ctx)
	if err != nil {
		return models.Reports{}, err
	}
	var respReport models.Report
	var currentAllReports = make(map[uint64]models.Report)
	var respAllReports models.Reports
	counter := 0

	for _, value := range allReports {
		respReport.ReportDesc = value
		currentAllReports[uint64(counter)] = respReport
		counter++
	}

	respAllReports.AllReports = currentAllReports
	respAllReports.Count = uint64(counter)

	return respAllReports, nil
}

func (h *userUsecase) AddReport(c context.Context, report models.NewReport) error {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUser, ok := ctx.Value(configs.ContextUser).(models.User)
	if !ok {
		return models.ErrContextNilError
	}

	// add new report
	err := h.UserRepo.AddReport(ctx, report)
	if err != nil {
		return err
	}

	// delete likes with this user
	err = h.UserRepo.DeleteReaction(ctx, currentUser.ID, report.ToId)
	if err != nil {
		return err
	}
	// delete matches with this user
	err = h.UserRepo.DeleteMatches(ctx, currentUser.ID, report.ToId)
	if err != nil {
		return err
	}

	// added dislike(2) reaction in db
	err = h.UserRepo.AddReaction(ctx, currentUser.ID, report.ToId, models.DislikeReaction)
	if err != nil {
		return err
	}

	// report's count ToId user
	curCount, err := h.UserRepo.GetReportsCount(ctx, report.ToId)
	if err != nil {
		return err
	}

	// if report's count > limit -> ban
	if curCount > models.ReportLimit {
		banId, err := h.UserRepo.GetReportsWithMaxCount(ctx, report.ToId)
		if err != nil {
			return err
		}
		banDesc, err := h.UserRepo.GetReportDesc(ctx, banId)
		if err != nil {
			return err
		}

		var reportStatus string
		switch banDesc {
		case models.FakeReport:
			reportStatus = FAKE
		case models.AggressionReport:
			reportStatus = AGGRESSION
		case models.SkamReport:
			reportStatus = SCAM
		case models.UnderageReport:
			reportStatus = UNDERAGE
		}

		err = h.UserRepo.UpdateReportStatus(ctx, report.ToId, reportStatus)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *userUsecase) CreatePayment(c context.Context, newPayment models.Payment) (models.RedirectUrl, error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUser, ok := ctx.Value(configs.ContextUser).(models.User)
	if !ok {
		return models.RedirectUrl{}, models.ErrContextNilError
	}

	var amount = make(map[string]string)
	amount["value"] = newPayment.Amount
	amount["currency"] = configs.Payment.Currency
	var confirmation = make(map[string]string)
	confirmation["type"] = "redirect"
	confirmation["return_url"] = configs.Payment.ReturnUrl
	var paymentInfo models.PaymentInfo
	paymentInfo.Amount = amount
	paymentInfo.Capture = true
	paymentInfo.Confirmation = confirmation

	paymentInfoJSON, err := json.Marshal(paymentInfo)
	if err != nil {
		return models.RedirectUrl{}, err
	}

	paymentRequest, err := http.NewRequest("POST", configs.Payment.YooKassaUrl, bytes.NewBuffer(paymentInfoJSON))
	if err != nil {
		return models.RedirectUrl{}, err
	}
	paymentRequest.Header.Set("Authorization", "Basic "+configs.Payment.AuthToken)
	paymentRequest.Header.Set("Idempotence-Key", uuid.NewString())
	paymentRequest.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient

	resp, err := client.Do(paymentRequest)
	if err != nil {
		return models.RedirectUrl{}, err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.RedirectUrl{}, err
	}
	var yooKassaResponse models.YooKassaResponse
	err = json.Unmarshal(buf, &yooKassaResponse)
	if err != nil {
		return models.RedirectUrl{}, err
	}

	var curRedirect models.RedirectUrl
	curRedirect.URL = yooKassaResponse.Confirmation.ConfirmationUrl

	err = h.UserRepo.CreatePayment(ctx, yooKassaResponse.Id, yooKassaResponse.Status, yooKassaResponse.Amount.Value+yooKassaResponse.Amount.Currency, currentUser.ID)
	if err != nil {
		return models.RedirectUrl{}, err
	}

	periodStart, err := time.Parse(models.DateLayout, yooKassaResponse.CreatedAt)
	if err != nil {
		return models.RedirectUrl{}, err
	}
	periodEnd := periodStart.AddDate(0, int(newPayment.Period), 0)

	err = h.UserRepo.CreateSubscription(ctx, periodStart, periodEnd, currentUser.ID, yooKassaResponse.Id)
	if err != nil {
		return models.RedirectUrl{}, err
	}

	return curRedirect, nil
}

func (h *userUsecase) UpdatePayment(c context.Context, paymentNotificationData models.PaymentNotification) error {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	err := h.UserRepo.UpdatePayment(ctx, paymentNotificationData.Object.Id, paymentNotificationData.Object.Status)
	if err != nil {
		return err
	}

	if paymentNotificationData.Object.Status == models.PaymentStatusSuccessString {
		err = h.UserRepo.UpdateSubscription(ctx, paymentNotificationData.Object.Id, models.PaymentStatusSuccess)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *userUsecase) CheckSubscription(c context.Context) (models.Subscription, error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUser, ok := ctx.Value(configs.ContextUser).(models.User)
	if !ok {
		return models.Subscription{}, models.ErrContextNilError
	}

	isActive, err := h.UserRepo.CheckSubscription(ctx, currentUser.ID)
	if err != nil {
		return models.Subscription{}, err
	}

	return models.Subscription{SubscriptionActive: isActive}, nil
}
