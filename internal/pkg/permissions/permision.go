package permissions

import (
	"dripapp/configs"
	"dripapp/internal/pkg/models"
	"dripapp/internal/pkg/responses"

	"net/http"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

func CheckAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			session, ok := r.Context().Value(configs.ForContext).(models.Session)
			if !ok {
				responses.SendErrorResponse(w, &models.HTTPError{
					Code:    http.StatusForbidden,
					Message: "context extract error",
				})
				return
			}
			if session.UserID == 0 {
				responses.SendErrorResponse(w, &models.HTTPError{
					Code:    http.StatusForbidden,
					Message: "auth errror",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
}

func generateCsrfLogic(w http.ResponseWriter) {
	csrf, err := uuid.NewV4()
	if err != nil {
		responses.SendErrorResponse(w, &models.HTTPError{
			Code:    http.StatusForbidden,
			Message: "no permission",
		})
		return
	}
	timeDelta := time.Now().Add(time.Hour * 24 * 30)
	cookie1 := &http.Cookie{Name: "csrf", Value: csrf.String(), Path: "/", HttpOnly: true, Expires: timeDelta}

	http.SetCookie(w, cookie1)
	w.Header().Set("csrf", csrf.String())

}

func SetCSRF(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			generateCsrfLogic(w)
			next.ServeHTTP(w, r)
		})
}

func CheckCSRF(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			csrf := r.Header.Get("X-Csrf-Token")
			csrfCookie, err := r.Cookie("csrf")

			if err != nil || csrf == "" || csrfCookie.Value == "" || csrfCookie.Value != csrf {
				responses.SendErrorResponse(w, &models.HTTPError{
					Code:    http.StatusInternalServerError,
					Message: "csrf-protection",
				})
				return
			}
			generateCsrfLogic(w)
			next.ServeHTTP(w, r)
		})

}
