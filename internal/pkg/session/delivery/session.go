package delivery

import (
	"crypto/md5"
	"dripapp/internal/pkg/models"
	"dripapp/internal/pkg/responses"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusInternalServerError = 500
	StatusEmailAlreadyExists  = 1001
)

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
	var resp models.JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		responses.SendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	var logUserData *models.LoginUser
	err = json.Unmarshal(byteReq, &logUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		responses.SendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}
	user, status := h.UserUCase.Login(r.Context(), *logUserData)
	resp.Status = status
	if status == StatusOK {
		cookie := createSessionCookie(*logUserData)

		sess := models.Session{
			Cookie: cookie.Value,
			UserID: user.ID,
		}
		err = h.SessionUcase.AddSession(r.Context(), sess)
		if err != nil {
			resp.Status = StatusInternalServerError
			log.Printf("CODE %d ERROR %s", resp.Status, err)
			responses.SendResp(resp, w)
			return
		}
		resp.Body = user
		http.SetCookie(w, &cookie)
	}

	responses.SendResp(resp, w)
}

func (h *SessionHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := h.SessionUcase.DeleteSession(r.Context())
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, err)
		responses.SendResp(models.JSON{Status: StatusNotFound}, w)
	}
	session, err := r.Cookie("sessionId")
	if err != nil {
		log.Printf("CODE %d ERROR %s", StatusNotFound, err)
		responses.SendResp(models.JSON{Status: StatusNotFound}, w)
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	responses.SendResp(models.JSON{Status: StatusOK}, w)
}
