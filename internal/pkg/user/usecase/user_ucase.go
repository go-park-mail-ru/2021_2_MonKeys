package usecase

import (
	"context"
	"crypto/md5"
	"dripapp/internal/pkg/models"
	"dripapp/internal/pkg/responses"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type userUsecase struct {
	UserRepo       models.UserRepository
	Session        models.SessionRepository
	contextTimeout time.Duration
}

const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusInternalServerError = 500
	StatusEmailAlreadyExists  = 1001
)

const maxPhotoSize = 20 * 1024 * 1025 // - это из доставки. Пока пусть будет здесь для AddPhoto()

func NewUserUsecase(ur models.UserRepository, sess models.SessionRepository, timeout time.Duration) models.UserUsecase {
	return &userUsecase{
		UserRepo:       ur,
		Session:        sess,
		contextTimeout: timeout,
	}
}

func createSessionCookie(user models.LoginUser) http.Cookie {
	expiration := time.Now().Add(10 * time.Hour)

	data := user.Password + time.Now().String()
	md5CookieValue := fmt.Sprintf("%x", md5.Sum([]byte(data)))

	cookie := http.Cookie{
		Name:     "sessionId",
		Value:    md5CookieValue,
		Expires:  expiration,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	return cookie
}

func (h *userUsecase) CurrentUser(c context.Context) (models.User, int) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	ctxSession := ctx.Value("userID")
	if ctxSession == nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, errors.New("context nil error"))
		return models.User{}, StatusNotFound
	}
	currentSession, ok := ctxSession.(models.Session)
	if !ok {
		log.Printf("CODE %d ERROR %s", StatusNotFound, errors.New("convert to model session error"))
		return models.User{}, StatusNotFound
	}

	currentUser, err := h.UserRepo.GetUserByID(c, currentSession.UserID)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, err)
		return models.User{}, StatusNotFound
	}

	log.Printf("CODE %d", StatusOK)
	return currentUser, StatusOK
}

func (h *userUsecase) EditProfile(c context.Context, newUserData models.User) (models.User, int) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	ctxSession := ctx.Value("userID")
	if ctxSession == nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, errors.New("context nil error"))
		return models.User{}, StatusNotFound
	}
	currentSession, ok := ctxSession.(models.Session)
	if !ok {
		log.Printf("CODE %d ERROR %s", StatusNotFound, errors.New("convert to model session error"))
		return models.User{}, StatusNotFound
	}

	currentUser, err := h.UserRepo.GetUserByID(c, currentSession.UserID)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, err)
		return models.User{}, StatusNotFound
	}

	err = currentUser.FillProfile(&newUserData)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusBadRequest, err)
		return models.User{}, StatusBadRequest
	}

	_, err = h.UserRepo.UpdateUser(c, currentUser)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusInternalServerError, err)
		return models.User{}, StatusInternalServerError
	}

	log.Printf("CODE %d", StatusOK)

	return currentUser, StatusOK
}

func (h *userUsecase) AddPhoto(c context.Context, w http.ResponseWriter, r *http.Request) {
	var resp models.JSON
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUserId := ctx.Value("userID")
	if currentUserId == nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, errors.New("context nil error"))
		return
	}
	userId := currentUserId.(uint64)

	currentUser, err := h.UserRepo.GetUserByID(c, userId)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, err)
		return
	}

	err = r.ParseMultipartForm(maxPhotoSize)
	if err != nil {
		resp.Status = StatusBadRequest
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
		resp.Status = StatusInternalServerError
		responses.SendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	resp.Status = StatusOK
	resp.Body = models.Photo{Title: currentUser.GetLastPhoto()}
	responses.SendResp(resp, w)
}

func (h *userUsecase) DeletePhoto(c context.Context, w http.ResponseWriter, r *http.Request) {
	var resp models.JSON
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	currentUserId := ctx.Value("userID")
	if currentUserId == nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, errors.New("context nil error"))
		return
	}
	userId := currentUserId.(uint64)

	currentUser, err := h.UserRepo.GetUserByID(c, userId)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, err)
		return
	}

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		responses.SendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	var photo *models.Photo
	err = json.Unmarshal(byteReq, &photo)
	if err != nil {
		resp.Status = StatusBadRequest
		responses.SendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	if currentUser.IsHavePhoto(photo.Title) {
		resp.Status = StatusBadRequest
		responses.SendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	err = h.UserRepo.DeletePhoto(c, currentUser, photo.Title)
	if err != nil {
		resp.Status = StatusInternalServerError
		responses.SendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	resp.Status = StatusOK
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
func (h *userUsecase) Login(c context.Context, logUserData models.LoginUser, w http.ResponseWriter) (models.User, int) {

	identifiableUser, err := h.UserRepo.GetUser(c, logUserData.Email)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, err)
		return models.User{}, StatusNotFound
	}

	if identifiableUser.IsCorrectPassword(logUserData.Password) {
		cookie := createSessionCookie(logUserData)

		if !h.Session.IsSessionByCookie(cookie.Value) {
			err = h.Session.NewSessionCookie(cookie.Value, identifiableUser.ID)
			if err != nil {
				log.Printf("CODE %d ERROR %s", StatusInternalServerError, err)
				return models.User{}, StatusInternalServerError
			}
		}
		http.SetCookie(w, &cookie)

		log.Printf("CODE %d", StatusOK)
		return identifiableUser, StatusOK
	} else {
		log.Printf("CODE %d ERROR %s", StatusNotFound, errors.New("not correct password"))
		return models.User{}, StatusNotFound
	}
}

func (h *userUsecase) Logout(c context.Context) int {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	ctxSession := ctx.Value("userID")
	if ctxSession == nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, errors.New("context nil error"))
		return StatusNotFound
	}
	currentSession, ok := ctxSession.(models.Session)
	if !ok {
		log.Printf("CODE %d ERROR %s", StatusNotFound, errors.New("convert to model session error"))
		return StatusNotFound
	}

	err := h.Session.DeleteSessionCookie(currentSession.Cookie)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusInternalServerError, err)
		return StatusInternalServerError
	}

	log.Printf("CODE %d", StatusOK)
	return StatusOK
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
func (h *userUsecase) Signup(c context.Context, logUserData models.LoginUser, w http.ResponseWriter) int {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	identifiableUser, _ := h.UserRepo.GetUser(ctx, logUserData.Email)
	if !identifiableUser.IsEmpty() {
		log.Printf("CODE %d ERROR %s", StatusEmailAlreadyExists, models.ErrEmailAlreadyExists)
		return StatusEmailAlreadyExists
	}

	logUserData.Password = models.HashPassword(logUserData.Password)
	user, err := h.UserRepo.CreateUser(c, logUserData)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusInternalServerError, err)
		return StatusInternalServerError
	}

	cookie := createSessionCookie(logUserData)

	err = h.Session.NewSessionCookie(cookie.Value, user.ID)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusInternalServerError, err)
		return StatusInternalServerError
	}

	http.SetCookie(w, &cookie)

	log.Printf("CODE %d", StatusOK)

	return StatusOK
}

func (h *userUsecase) NextUser(c context.Context) ([]models.User, int) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	ctxSession := ctx.Value("userID")
	if ctxSession == nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, errors.New("context nil error"))
		return nil, StatusNotFound
	}
	currentSession, ok := ctxSession.(models.Session)
	if !ok {
		log.Printf("CODE %d ERROR %s", StatusNotFound, errors.New("convert to model session error"))
		return nil, StatusNotFound
	}

	currentUser, err := h.UserRepo.GetUserByID(c, currentSession.UserID)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, err)
		return []models.User{}, StatusNotFound
	}

	// add in swaped users map for current user
	// err = h.UserRepo.AddSwipedUsers(ctx, currentUser.ID, swipedUserData.Id, "like")
	// if err != nil {
	// 	log.Printf("CODE %d ERROR %s", StatusNotFound, err)
	// 	return models.User{}, StatusNotFound
	// }
	// find next user for swipe
	nextUsers, err := h.UserRepo.GetNextUserForSwipe(ctx, currentUser.ID)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, err)
		return []models.User{}, StatusNotFound
	}

	log.Printf("CODE %d", StatusOK)

	return nextUsers, StatusOK
}

func (h *userUsecase) GetAllTags(c context.Context) (models.Tags, int) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	ctxSession := ctx.Value("userID")
	if ctxSession == nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, errors.New("context nil error"))
		return models.Tags{}, StatusNotFound
	}
	currentSession, ok := ctxSession.(models.Session)
	if !ok {
		log.Printf("CODE %d ERROR %s", StatusNotFound, errors.New("convert to model session error"))
		return models.Tags{}, StatusNotFound
	}

	_, err := h.UserRepo.GetUserByID(c, currentSession.UserID)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, err)
		return models.Tags{}, StatusNotFound
	}

	allTags, err := h.UserRepo.GetTags(ctx)
	if err != nil {
		return models.Tags{}, StatusNotFound
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

	return respAllTags, StatusOK
}
