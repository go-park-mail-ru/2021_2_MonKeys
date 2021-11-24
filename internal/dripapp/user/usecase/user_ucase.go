package usecase

import (
	"context"
	"dripapp/configs"
	"dripapp/internal/dripapp/models"
	_sessionModels "dripapp/internal/microservices/auth/models"
	"dripapp/internal/pkg/hasher"
	"fmt"
	"io"
	"strconv"
	"time"
)

type userUsecase struct {
	UserRepo       models.UserRepository
	Session        _sessionModels.SessionRepository
	File           models.FileRepository
	contextTimeout time.Duration
}

func NewUserUsecase(
	ur models.UserRepository,
	fileManager models.FileRepository,
	timeout time.Duration) models.UserUsecase {
	return &userUsecase{
		UserRepo:       ur,
		File:           fileManager,
		contextTimeout: timeout,
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
	if err != nil {
		return models.User{}, err
	}

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

	return currMath, nil
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
		banId, err := h.UserRepo.GetReportsWithMaxCountCount(ctx, report.ToId)
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
			reportStatus = "FAKE"
		case models.AggressionReport:
			reportStatus = "AGGRESSION"
		case models.SkamReport:
			reportStatus = "SKAM"
		case models.UnderageReport:
			reportStatus = "UNDERAGE"
		}

		err = h.UserRepo.UpdateReportStatus(ctx, report.ToId, reportStatus)
		if err != nil {
			return err
		}
	}

	return nil
}
