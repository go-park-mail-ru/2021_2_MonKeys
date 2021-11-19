package delivery

import (
	"crypto/md5"
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/logger"
	"dripapp/internal/pkg/responses"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type SessionHandler struct {
	Logger       logger.Logger
	UserUCase    models.UserUsecase
	SessionUcase models.SessionUsecase
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

func (h *SessionHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var resp responses.JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	var logUserData models.LoginUser
	err = json.Unmarshal(byteReq, &logUserData)
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	user, status := h.UserUCase.Login(r.Context(), logUserData)
	if status != nil {
		resp.Status = http.StatusNotFound
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    resp.Status,
			Message: status,
		}, h.Logger.WarnLogging)
		return
	}

	sessionCookie := createSessionCookie(logUserData)

	sess := models.Session{
		Cookie: sessionCookie.Value,
		UserID: user.ID,
	}
	err = h.SessionUcase.AddSession(r.Context(), sess)
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err,
		}, h.Logger.WarnLogging)
		return
	}

	http.SetCookie(w, &sessionCookie)

	resp.Status = http.StatusOK
	resp.Body = user

	responses.SendOKResp(resp, w)
}

func (h *SessionHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := h.SessionUcase.DeleteSession(r.Context())
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
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

	responses.SendOKResp(responses.JSON{Status: http.StatusOK}, w)
}
