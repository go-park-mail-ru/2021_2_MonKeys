package Handlers

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"server/Models"
	"time"
)

const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusInternalServerError = 500
	StatusEmailAlreadyExists  = 1001
)

func sendResp(resp Models.JSON, w http.ResponseWriter) {
	byteResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteResp)
}

func createSessionCookie(user Models.LoginUser) http.Cookie {
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

type Env struct {
	DB interface {
		GetUser(email string) (Models.User, error)
		GetUserByID(userID uint64) (Models.User, error)
		CreateUser(logUserData Models.LoginUser) (Models.User, error)
		AddSwipedUsers(currentUserId, swipedUserId uint64) error
		GetNextUserForSwipe(currentUserId uint64) (Models.User, error)
		UpdateUser(newUserData Models.User) error
	}
	SessionDB interface {
		GetUserIDByCookie(sessionCookie string) (uint64, error)
		NewSessionCookie(sessionCookie string, userId uint64) error
		DeleteSessionCookie(sessionCookie string) error
	}
}

func (env *Env) CurrentUser(w http.ResponseWriter, r *http.Request) {
	var resp Models.JSON
	session, err := r.Cookie("sessionId")
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		return
	}

	currentUser, err := env.getUserByCookie(session.Value)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		return
	}

	resp.Status = StatusOK
	resp.Body = currentUser

	sendResp(resp, w)
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
func (env *Env) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var resp Models.JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		return
	}

	var logUserData Models.LoginUser
	err = json.Unmarshal(byteReq, &logUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		return
	}

	identifiableUser, err := env.DB.GetUser(logUserData.Email)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		return
	}

	status := StatusOK
	if identifiableUser.IsCorrectPassword(logUserData.Password) {
		cookie := createSessionCookie(logUserData)
		err = env.SessionDB.NewSessionCookie(cookie.Value, identifiableUser.ID)
		if err != nil {
			resp.Status = StatusInternalServerError
			sendResp(resp, w)
			return
		}

		http.SetCookie(w, &cookie)

		resp.Body = identifiableUser
	} else {
		status = StatusNotFound
	}

	resp.Status = status
	sendResp(resp, w)
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
func (env *Env) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var resp Models.JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		return
	}

	var logUserData Models.LoginUser
	err = json.Unmarshal(byteReq, &logUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		return
	}

	identifiableUser, _ := env.DB.GetUser(logUserData.Email)
	if !identifiableUser.IsEmpty() {
		resp.Status = StatusEmailAlreadyExists
		sendResp(resp, w)
		return
	}

	user, err := env.DB.CreateUser(logUserData)
	if err != nil {
		resp.Status = StatusInternalServerError
		sendResp(resp, w)
		return
	}

	cookie := createSessionCookie(logUserData)
	err = env.SessionDB.NewSessionCookie(cookie.Value, user.ID)
	if err != nil {
		resp.Status = StatusInternalServerError
		sendResp(resp, w)
		return
	}

	http.SetCookie(w, &cookie)

	resp.Status = StatusOK
	sendResp(resp, w)
}

func (env *Env) EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	var resp Models.JSON
	session, err := r.Cookie("sessionId")
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		return
	}

	currentUser, err := env.getUserByCookie(session.Value)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		return
	}

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		return
	}

	var newUserData Models.User
	err = json.Unmarshal(byteReq, &newUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		return
	}

	err = currentUser.FillProfile(newUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		return
	}

	err = env.DB.UpdateUser(currentUser)
	if err != nil {
		resp.Status = StatusInternalServerError
		sendResp(resp, w)
		return
	}

	resp.Status = StatusOK
	resp.Body = currentUser

	sendResp(resp, w)
}

func (env *Env) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("sessionId")
	if err != nil {
		sendResp(Models.JSON{Status: StatusNotFound}, w)
		return
	}

	err = env.SessionDB.DeleteSessionCookie(session.Value)
	if err != nil {
		sendResp(Models.JSON{Status: StatusInternalServerError}, w)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
}

func (env *Env) NextUserHandler(w http.ResponseWriter, r *http.Request) {
	var resp Models.JSON

	// get current user by cookie
	session, err := r.Cookie("sessionId")
	if err == http.ErrNoCookie {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		return
	}
	currentUser, err := env.getUserByCookie(session.Value)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		return
	}

	// get swiped usedata for registrationr id from json
	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		return
	}
	var swipedUserData Models.SwipedUser
	err = json.Unmarshal(byteReq, &swipedUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		return
	}

	// add in swaped users map for current user
	err = env.DB.AddSwipedUsers(currentUser.ID, swipedUserData.Id)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		return
	}
	// find next user for swipe
	nextUser, err := env.DB.GetNextUserForSwipe(currentUser.ID)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		return
	}

	resp.Status = StatusOK
	resp.Body = nextUser

	sendResp(resp, w)
}

func (env Env) getUserByCookie(sessionCookie string) (Models.User, error) {
	userID, err := env.SessionDB.GetUserIDByCookie(sessionCookie)
	if err != nil {
		return Models.User{}, errors.New("error sessionDB: GetUserIDByCookie")
	}

	user, err := env.DB.GetUserByID(userID)
	if err != nil {
		return Models.User{}, errors.New("error db: getUserByID")
	}

	return user, nil
}
