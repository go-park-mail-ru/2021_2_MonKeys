package usecase

import (
	"context"
	"dripapp/configs"
	"dripapp/internal/pkg/models"
	"dripapp/internal/pkg/responses"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type userUsecase struct {
	UserRepo       models.UserRepository
	Session        models.SessionRepository
	contextTimeout time.Duration
}

const maxPhotoSize = 20 * 1024 * 1025 // - это из доставки. Пока пусть будет здесь для AddPhoto()

func NewUserUsecase(ur models.UserRepository, sess models.SessionRepository, timeout time.Duration) models.UserUsecase {
	return &userUsecase{
		UserRepo:       ur,
		Session:        sess,
		contextTimeout: timeout,
	}
}

func (h *userUsecase) CurrentUser(c context.Context) (models.User, int) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	ctxSession := ctx.Value(configs.ForContext)
	if ctxSession == nil {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, errors.New("context nil error"))
		return models.User{}, http.StatusNotFound
	}
	currentSession, ok := ctxSession.(models.Session)
	if !ok {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, errors.New("convert to model session error"))
		return models.User{}, http.StatusNotFound
	}

	currentUser, err := h.UserRepo.GetUserByID(c, currentSession.UserID)
	if err != nil {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, err)
		return models.User{}, http.StatusNotFound
	}

	log.Printf("CODE %d", http.StatusOK)
	return currentUser, http.StatusOK
}

func (h *userUsecase) EditProfile(c context.Context, newUserData models.User) (models.User, int) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	ctxSession := ctx.Value(configs.ForContext)
	if ctxSession == nil {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, errors.New("context nil error"))
		return models.User{}, http.StatusNotFound
	}
	currentSession, ok := ctxSession.(models.Session)
	if !ok {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, errors.New("convert to model session error"))
		return models.User{}, http.StatusNotFound
	}

	currentUser, err := h.UserRepo.GetUserByID(c, currentSession.UserID)
	if err != nil {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, err)
		return models.User{}, http.StatusNotFound
	}

	err = currentUser.FillProfile(&newUserData)
	if err != nil {
		log.Printf("CODE %d ERROR %s", http.StatusBadRequest, err)
		return models.User{}, http.StatusBadRequest
	}

	_, err = h.UserRepo.UpdateUser(c, currentUser)
	if err != nil {
		log.Printf("CODE %d ERROR %s", http.StatusInternalServerError, err)
		return models.User{}, http.StatusInternalServerError
	}

	log.Printf("CODE %d", http.StatusOK)

	return currentUser, http.StatusOK
}

func (h *userUsecase) AddPhoto(c context.Context, w http.ResponseWriter, r *http.Request) {
	var resp models.JSON
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUserId := ctx.Value(configs.ForContext)
	if currentUserId == nil {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, errors.New("context nil error"))
		return
	}
	userId := currentUserId.(uint64)

	currentUser, err := h.UserRepo.GetUserByID(c, userId)
	if err != nil {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, err)
		return
	}

	err = r.ParseMultipartForm(maxPhotoSize)
	if err != nil {
		resp.Status = http.StatusBadRequest
		responses.SendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	uploadedPhoto, _, err := r.FormFile("photo")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	defer uploadedPhoto.Close()

	currentUser.SaveNewPhoto()

	err = h.UserRepo.AddPhoto(c, currentUser, uploadedPhoto)
	if err != nil {
		resp.Status = http.StatusInternalServerError
		responses.SendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	resp.Status = http.StatusOK
	resp.Body = models.Photo{Title: currentUser.GetLastPhoto()}
	responses.SendResp(resp, w)
}

func (h *userUsecase) DeletePhoto(c context.Context, w http.ResponseWriter, r *http.Request) {
	var resp models.JSON
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUserId := ctx.Value(configs.ForContext)
	if currentUserId == nil {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, errors.New("context nil error"))
		return
	}
	userId := currentUserId.(uint64)

	currentUser, err := h.UserRepo.GetUserByID(c, userId)
	if err != nil {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, err)
		return
	}

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = http.StatusBadRequest
		responses.SendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	var photo *models.Photo
	err = json.Unmarshal(byteReq, &photo)
	if err != nil {
		resp.Status = http.StatusBadRequest
		responses.SendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	if currentUser.IsHavePhoto(photo.Title) {
		resp.Status = http.StatusBadRequest
		responses.SendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	err = h.UserRepo.DeletePhoto(c, currentUser, photo.Title)
	if err != nil {
		resp.Status = http.StatusInternalServerError
		responses.SendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	resp.Status = http.StatusOK
	responses.SendResp(resp, w)
}

// @Summary LogIn
// @Description log in
// @Tags login
// @Accept json
// @Produce json
// @Param input body LoginUser true "data for login"
// @Success 200 {object} JSON
// @Failure 400,404,500
// @Router /login [post]
func (h *userUsecase) Login(c context.Context, logUserData models.LoginUser) (models.User, int) {

	identifiableUser, err := h.UserRepo.GetUser(c, logUserData.Email)
	if err != nil {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, err)
		return models.User{}, http.StatusNotFound
	}

	if identifiableUser.IsCorrectPassword(logUserData.Password) {
		log.Printf("CODE %d", http.StatusOK)
		return identifiableUser, http.StatusOK
	} else {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, errors.New("not correct password"))
		return models.User{}, http.StatusNotFound
	}
}

// @Summary SignUp
// @Description registration user
// @Tags registration
// @Accept json
// @Produce json
// @Param input body LoginUser true "data for registration"
// @Success 200 {object} JSON
// @Failure 400,404,500
// @Router /signup [post]
func (h *userUsecase) Signup(c context.Context, logUserData models.LoginUser) (models.User, int) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	identifiableUser, _ := h.UserRepo.GetUser(ctx, logUserData.Email)
	if !identifiableUser.IsEmpty() {
		log.Printf("CODE %d ERROR %s", models.StatusEmailAlreadyExists, models.ErrEmailAlreadyExists)
		return models.User{}, models.StatusEmailAlreadyExists
	}

	logUserData.Password = models.HashPassword(logUserData.Password)
	user, err := h.UserRepo.CreateUser(c, logUserData)
	if err != nil {
		log.Printf("CODE %d ERROR %s", http.StatusInternalServerError, err)
		return models.User{}, http.StatusInternalServerError
	}

	log.Printf("CODE %d", http.StatusOK)

	return user, http.StatusOK
}

func (h *userUsecase) NextUser(c context.Context) ([]models.User, int) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	ctxSession := ctx.Value(configs.ForContext)
	if ctxSession == nil {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, errors.New("context nil error"))
		return nil, http.StatusNotFound
	}
	currentSession, ok := ctxSession.(models.Session)
	if !ok {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, errors.New("convert to model session error"))
		return nil, http.StatusNotFound
	}

	currentUser, err := h.UserRepo.GetUserByID(c, currentSession.UserID)
	if err != nil {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, err)
		return []models.User{}, http.StatusNotFound
	}

	// add in swaped users map for current user
	// err = h.UserRepo.AddSwipedUsers(ctx, currentUser.ID, swipedUserData.Id, "like")
	// if err != nil {
	// 	log.Printf("CODE %d ERROR %s", http.StatusNotFound, err)
	// 	return models.User{}, http.StatusNotFound
	// }
	nextUsers, err := h.UserRepo.GetNextUserForSwipe(ctx, currentUser.ID)
	if err != nil {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, err)
		return []models.User{}, http.StatusNotFound
	}

	log.Printf("CODE %d", http.StatusOK)

	return nextUsers, http.StatusOK
}

func (h *userUsecase) GetAllTags(c context.Context) (models.Tags, int) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	ctxSession := ctx.Value(configs.ForContext)
	if ctxSession == nil {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, errors.New("context nil error"))
		return models.Tags{}, http.StatusNotFound
	}
	currentSession, ok := ctxSession.(models.Session)
	if !ok {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, errors.New("convert to model session error"))
		return models.Tags{}, http.StatusNotFound
	}

	_, err := h.UserRepo.GetUserByID(c, currentSession.UserID)
	if err != nil {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, err)
		return models.Tags{}, http.StatusNotFound
	}

	allTags, err := h.UserRepo.GetTags(ctx)
	if err != nil {
		return models.Tags{}, http.StatusNotFound
	}
	var respTag models.Tag
	var currentAllTags = make(map[uint64]models.Tag)
	var respAllTags models.Tags
	counter := 0

	for _, value := range allTags {
		respTag.Tag_Name = value
		currentAllTags[uint64(counter)] = respTag
		counter++
	}

	respAllTags.AllTags = currentAllTags
	respAllTags.Count = uint64(counter)

	return respAllTags, http.StatusOK
}

func (h *userUsecase) UsersMatches(c context.Context, r *http.Request) (models.Matches, int) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	ctxSession := ctx.Value(configs.ForContext)
	if ctxSession == nil {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, errors.New("context nil error"))
		return models.Matches{}, http.StatusNotFound
	}
	currentSession, ok := ctxSession.(models.Session)
	if !ok {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, errors.New("convert to model session error"))
		return models.Matches{}, http.StatusNotFound
	}

	// find matches
	mathes, err := h.UserRepo.GetUsersMatches(ctx, currentSession.UserID)
	if err != nil {
		log.Printf("CODE %d ERROR %s", http.StatusNotFound, err)
		return models.Matches{}, http.StatusNotFound
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

	log.Printf("CODE %d", http.StatusOK)

	return allMatches, http.StatusOK
}
