package http

import (
	"crypto/md5"
	_userModels "dripapp/internal/dripapp/models"
	_sessionModels "dripapp/internal/microservices/auth/models"
	"dripapp/internal/pkg/logger"
	"dripapp/internal/pkg/responses"
	"fmt"
	"net/http"
	"time"
)

type SessionHandler struct {
	Logger       logger.Logger
	UserUCase    _userModels.UserUsecase
	SessionUcase _sessionModels.SessionUsecase
}

func createSessionCookie(user _userModels.LoginUser) http.Cookie {
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

func (h *SessionHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var logUserData _userModels.LoginUser
	err := responses.ReadJSON(r, &logUserData)
	if err != nil {
		responses.SendError(w, _userModels.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	user, err := h.UserUCase.Login(r.Context(), logUserData)
	if err != nil {
		responses.SendError(w, _userModels.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.WarnLogging)
		return
	}

	sessionCookie := createSessionCookie(logUserData)

	sess := _sessionModels.Session{
		Cookie: sessionCookie.Value,
		UserID: user.ID,
	}
	err = h.SessionUcase.AddSession(r.Context(), sess)
	if err != nil {
		responses.SendError(w, _userModels.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err,
		}, h.Logger.WarnLogging)
		return
	}

	http.SetCookie(w, &sessionCookie)

	responses.SendData(w, user)
}

func (h *SessionHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := h.SessionUcase.DeleteSession(r.Context())
	if err != nil {
		responses.SendError(w, _userModels.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.WarnLogging)
		return
	}

	authCookie := &http.Cookie{
		Name:     "sessionId",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(w, authCookie)

	csrfCookie := &http.Cookie{
		Name:     "csrf",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(w, csrfCookie)

	responses.SendOK(w)
}

func (h *SessionHandler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.UserUCase.CurrentUser(r.Context())
	if err != nil {
		responses.SendError(w, _userModels.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendData(w, user)
}

func (h *SessionHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var logUserData _userModels.LoginUser
	err := responses.ReadJSON(r, &logUserData)
	if err != nil {
		responses.SendError(w, _userModels.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	user, err := h.UserUCase.Signup(r.Context(), logUserData)
	if err != nil {
		code := http.StatusNotFound
		if err == _userModels.ErrEmailAlreadyExists {
			code = _userModels.StatusEmailAlreadyExists
		}
		responses.SendError(w, _userModels.HTTPError{
			Code:    code,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}
	cookie := _sessionModels.CreateSessionCookie(logUserData)

	sess := _sessionModels.Session{
		Cookie: cookie.Value,
		UserID: user.ID,
	}

	err = h.SessionUcase.AddSession(r.Context(), sess)
	if err != nil {
		responses.SendError(w, _userModels.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err,
		}, h.Logger.WarnLogging)
		return
	}

	http.SetCookie(w, &cookie)

	responses.SendData(w, user)
}
