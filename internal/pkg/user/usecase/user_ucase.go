package usecase

import (
	"context"
	"crypto/md5"
	"dripapp/internal/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
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

func sendResp(resp models.JSON, w http.ResponseWriter) {
	byteResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(byteResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *userUsecase) getUserByCookie(c context.Context, sessionCookie string) (models.User, error) {
	userID, err := h.Session.GetUserIDByCookie(sessionCookie)
	if err != nil {
		return models.User{},
			err
	}

	user, err := h.UserRepo.GetUserByID(c, userID)
	if err != nil {
		return models.User{}, errors.New("error db: getUserByID")
	}

	return *user, nil
}

func (h *userUsecase) CurrentUser(c context.Context, r *http.Request) (models.User, int) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	session, err := r.Cookie("sessionId")
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, err)
		return models.User{}, StatusNotFound
	}

	currentUser, err := h.getUserByCookie(ctx, session.Value)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, err)
		return models.User{}, StatusNotFound
	}

	log.Printf("CODE %d", StatusOK)
	return currentUser, StatusOK
}

func (h *userUsecase) EditProfile(c context.Context, newUserData models.User, r *http.Request) (models.User, int) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	session, err := r.Cookie("sessionId")
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, err)
		return models.User{}, StatusNotFound
	}

	fmt.Println("get user", session.Value)
	currentUser, err := h.getUserByCookie(ctx, session.Value)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, err)
		return models.User{}, StatusNotFound
	}

	err = currentUser.FillProfile(&newUserData)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusBadRequest, err)
		return models.User{}, StatusBadRequest
	}

	err = h.UserRepo.UpdateUser(c, &currentUser)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusInternalServerError, err)
		return models.User{}, StatusInternalServerError
	}

	log.Printf("CODE %d", StatusOK)

	return currentUser, StatusOK
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

		if !h.Session.IsSessionByUserID(identifiableUser.ID) {
			http.SetCookie(w, &cookie)
			err = h.Session.NewSessionCookie(cookie.Value, identifiableUser.ID)
			if err != nil {
				log.Printf("CODE %d ERROR %s", StatusInternalServerError, err)
				return models.User{}, StatusInternalServerError
			}
		}

		return *identifiableUser, StatusOK
	} else {
		log.Printf("CODE %d ERROR %s", StatusNotFound, errors.New("not correct password"))
		return models.User{}, StatusNotFound
	}
}

func (h *userUsecase) Logout(c context.Context, w http.ResponseWriter, r *http.Request) int {
	session, err := r.Cookie("sessionId")
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, err)
		return StatusNotFound
	}

	err = h.Session.DeleteSessionCookie(session.Value)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusInternalServerError, err)
		return StatusInternalServerError
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

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

	user, err := h.UserRepo.CreateUser(c, &logUserData)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusInternalServerError, err)
		return StatusInternalServerError
	}

	cookie := createSessionCookie(logUserData)

	if !h.Session.IsSessionByUserID(identifiableUser.ID) {
		http.SetCookie(w, &cookie)
		err = h.Session.NewSessionCookie(cookie.Value, user.ID)
		if err != nil {
			log.Printf("CODE %d ERROR %s", StatusInternalServerError, err)
			return StatusInternalServerError
		}
	}

	log.Printf("CODE %d", StatusOK)

	return StatusOK
}

func (h *userUsecase) NextUser(c context.Context, swipedUserData models.SwipedUser, r *http.Request) (models.User, int) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	// get current user by cookie
	session, err := r.Cookie("sessionId")
	if err == http.ErrNoCookie {
		log.Printf("CODE %d ERROR %s", StatusNotFound, err)
		return models.User{}, StatusNotFound
	}
	currentUser, err := h.getUserByCookie(ctx, session.Value)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, err)
		return models.User{}, StatusNotFound
	}

	// add in swaped users map for current user
	err = h.UserRepo.AddSwipedUsers(ctx, currentUser.ID, swipedUserData.Id)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, err)
		return models.User{}, StatusNotFound
	}
	// find next user for swipe
	nextUser, err := h.UserRepo.GetNextUserForSwipe(ctx, currentUser.ID)
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, err)
		return models.User{}, StatusNotFound
	}

	log.Printf("CODE %d", StatusOK)

	return *nextUser, StatusOK
}

func (h *userUsecase) GetAllTags(c context.Context, r *http.Request) (models.Tags, int) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	allTags := h.UserRepo.GetTags(ctx)
	var respTag models.Tag
	var currentAllTags = make(map[uint64]models.Tag)
	var respAllTags models.Tags
	counter := 0

	for key, value := range allTags {
		respTag.Id = key
		respTag.TagText = value
		currentAllTags[uint64(counter)] = respTag
		counter++
	}

	respAllTags.AllTags = currentAllTags
	respAllTags.Count = uint64(counter)

	return respAllTags, StatusOK
}
