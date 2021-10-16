package usecase

import (
	"crypto/md5"
	"dripapp/internal/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type userUsecase struct {
	UserRepo models.UserRepository
	Session  models.SessionRepository
	// contextTimeout time.Duration
}

const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusInternalServerError = 500
	StatusEmailAlreadyExists  = 1001
)

func NewUserUsecase(ur models.UserRepository, sess models.SessionRepository) models.UserUsecase {
	return &userUsecase{
		UserRepo: ur,
		Session:  sess,
		// contextTimeout: timeout,
	}
}

func createSessionCookie(user *models.LoginUser) http.Cookie {
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

func (h *userUsecase) getUserByCookie(sessionCookie string) (*models.User, error) {
	userID, err := h.Session.GetUserIDByCookie(sessionCookie)
	if err != nil {
		return &models.User{}, errors.New("error sessionDB: GetUserIDByCookie")
	}

	user, err := h.UserRepo.GetUserByID(userID)
	if err != nil {
		return &models.User{}, errors.New("error db: getUserByID")
	}

	return user, nil
}

func (h *userUsecase) CurrentUser(w http.ResponseWriter, r *http.Request) {
	var resp models.JSON
	session, err := r.Cookie("sessionId")
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)

		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	currentUser, err := h.getUserByCookie(session.Value)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	resp.Status = StatusOK
	resp.Body = currentUser

	sendResp(resp, w)

	log.Printf("CODE %d", resp.Status)
}

func (h *userUsecase) EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	var resp models.JSON
	session, err := r.Cookie("sessionId")
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	currentUser, err := h.getUserByCookie(session.Value)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	var newUserData models.User
	err = json.Unmarshal(byteReq, &newUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	err = currentUser.FillProfile(&newUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	err = h.UserRepo.UpdateUser(currentUser)
	if err != nil {
		resp.Status = StatusInternalServerError
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	resp.Status = StatusOK
	resp.Body = currentUser

	sendResp(resp, w)

	log.Printf("CODE %d", resp.Status)
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
func (h *userUsecase) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var resp models.JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	var logUserData *models.LoginUser
	err = json.Unmarshal(byteReq, &logUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	identifiableUser, err := h.UserRepo.GetUser(logUserData.Email)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	status := StatusOK
	if identifiableUser.IsCorrectPassword(logUserData.Password) {
		cookie := createSessionCookie(logUserData)
		err = h.Session.NewSessionCookie(cookie.Value, identifiableUser.ID)
		if err != nil {
			resp.Status = StatusInternalServerError
			sendResp(resp, w)
			log.Printf("CODE %d ERROR %s", resp.Status, err)
			return
		}

		http.SetCookie(w, &cookie)

		resp.Body = identifiableUser
	} else {
		status = StatusNotFound
	}

	resp.Status = status
	sendResp(resp, w)

	log.Printf("CODE %d", resp.Status)
}

func (h *userUsecase) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("sessionId")
	if err != nil {
		sendResp(models.JSON{Status: StatusNotFound}, w)
		log.Printf("CODE %d ERROR %s", StatusNotFound, err)
		return
	}

	err = h.Session.DeleteSessionCookie(session.Value)
	if err != nil {
		sendResp(models.JSON{Status: StatusInternalServerError}, w)
		log.Printf("CODE %d ERROR %s", StatusInternalServerError, err)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	log.Printf("CODE %d", StatusOK)
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
func (h *userUsecase) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var resp models.JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	var logUserData *models.LoginUser
	err = json.Unmarshal(byteReq, &logUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	identifiableUser, _ := h.UserRepo.GetUser(logUserData.Email)
	if !identifiableUser.IsEmpty() {
		resp.Status = StatusEmailAlreadyExists
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	user, err := h.UserRepo.CreateUser(logUserData)
	if err != nil {
		resp.Status = StatusInternalServerError
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	cookie := createSessionCookie(logUserData)
	err = h.Session.NewSessionCookie(cookie.Value, user.ID)
	if err != nil {
		resp.Status = StatusInternalServerError
		sendResp(resp, w)
		return
	}

	http.SetCookie(w, &cookie)

	resp.Status = StatusOK
	sendResp(resp, w)

	log.Printf("CODE %d", resp.Status)
}

func (h *userUsecase) NextUserHandler(w http.ResponseWriter, r *http.Request) {
	var resp models.JSON

	// get current user by cookie
	session, err := r.Cookie("sessionId")
	if err == http.ErrNoCookie {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}
	currentUser, err := h.getUserByCookie(session.Value)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	// get swiped usedata for registrationr id from json
	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}
	var swipedUserData models.SwipedUser
	err = json.Unmarshal(byteReq, &swipedUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	// add in swaped users map for current user
	err = h.UserRepo.AddSwipedUsers(currentUser.ID, swipedUserData.Id)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}
	// find next user for swipe
	nextUser, err := h.UserRepo.GetNextUserForSwipe(currentUser.ID)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	resp.Status = StatusOK
	resp.Body = nextUser

	sendResp(resp, w)

	log.Printf("CODE %d", resp.Status)
}

func (h *userUsecase) GetAllTags(w http.ResponseWriter, r *http.Request) {
	allTags := h.UserRepo.GetTags()
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

	var resp models.JSON

	resp.Status = http.StatusOK
	resp.Body = respAllTags

	byteResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(byteResp)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
