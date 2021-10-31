package delivery

import (
	"crypto/md5"
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/responses"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type SessionHandler struct {
	// Logger    *zap.SugaredLogger
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

// @Summary LogIn
// @Description log in
// @Tags login
// @Accept json
// @Produce json
// @Param input body LoginUser true "data for login"
// @Success 200 {object} JSON
// @Failure 400,404,500
// @Router /login [post]
func (h *SessionHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var resp responses.JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	var logUserData *models.LoginUser
	err = json.Unmarshal(byteReq, &logUserData)
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	user, status := h.UserUCase.Login(r.Context(), *logUserData)
	resp.Status = status.Code
	if status.Code == http.StatusOK {
		cookie := createSessionCookie(*logUserData)

		sess := models.Session{
			Cookie: cookie.Value,
			UserID: user.ID,
		}
		err = h.SessionUcase.AddSession(r.Context(), sess)
		if err != nil {
			responses.SendErrorResponse(w, models.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}
		resp.Body = user
		http.SetCookie(w, &cookie)
	}

	responses.SendOKResp(resp, w)
}

func (h *SessionHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := h.SessionUcase.DeleteSession(r.Context())
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}
	session, err := r.Cookie("sessionId")
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	session.Secure = true
	session.HttpOnly = true
	session.SameSite = http.SameSiteNoneMode
	http.SetCookie(w, session)

	responses.SendOKResp(responses.JSON{Status: http.StatusOK}, w)
}