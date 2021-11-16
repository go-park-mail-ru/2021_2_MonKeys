package permissions

import (
	"context"
	"dripapp/configs"
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/logger"
	"dripapp/internal/pkg/responses"
	"fmt"

	"net/http"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

type UserMiddlware struct {
	UserRepo models.UserRepository
}

func CheckAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			fmt.Println("check middlware")
			session, ok := r.Context().Value(configs.ContextUserID).(models.Session)
			if !ok {
				responses.SendErrorResponse(w, models.HTTPError{
					Code:    http.StatusForbidden,
					Message: models.ErrExtractContext,
				}, logger.DripLogger.ErrorLogging)
				return
			}
			if session.UserID == 0 {
				responses.SendErrorResponse(w, models.HTTPError{
					Code:    http.StatusForbidden,
					Message: models.ErrAuth,
				}, logger.DripLogger.WarnLogging)
				return
			}

			next.ServeHTTP(w, r)
		})
}

func (us *UserMiddlware) GetCurrentUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			fmt.Println("get current")
			ctxSession := r.Context().Value(configs.ContextUserID)
			if ctxSession == nil {
				responses.SendErrorResponse(w, models.HTTPError{
					Code:    http.StatusForbidden,
					Message: models.ErrExtractContext,
				}, logger.DripLogger.ErrorLogging)
				return
			}
			currentSession, ok := ctxSession.(models.Session)
			if !ok {
				responses.SendErrorResponse(w, models.HTTPError{
					Code:    http.StatusForbidden,
					Message: models.ErrExtractContext,
				}, logger.DripLogger.ErrorLogging)
				return
			}

			currentUser, err := us.UserRepo.GetUserByID(r.Context(), currentSession.UserID)
			if err != nil {
				responses.SendErrorResponse(w, models.HTTPError{
					Code:    http.StatusNotFound,
					Message: err.Error(),
				}, logger.DripLogger.ErrorLogging)
				return
			}

			if len(currentUser.Date) != 0 {
				currentUser.Age, err = models.GetAgeFromDate(currentUser.Date)
				if err != nil {
					responses.SendErrorResponse(w, models.HTTPError{
						Code:    http.StatusNotFound,
						Message: err.Error(),
					}, logger.DripLogger.ErrorLogging)
					return
				}
			}

			r = r.WithContext(context.WithValue(r.Context(), configs.ContextUser, currentUser))
			next.ServeHTTP(w, r)
		})
}

func generateCsrfLogic(w http.ResponseWriter) {
	csrf, err := uuid.NewV4()
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusForbidden,
			Message: "no permission",
		}, logger.DripLogger.ErrorLogging)
		return
	}
	timeDelta := time.Now().Add(time.Hour * 24 * 30)
	csrfCookie := &http.Cookie{Name: "csrf", Value: csrf.String(), Path: "/", HttpOnly: true, Expires: timeDelta}

	http.SetCookie(w, csrfCookie)
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
			csrf := r.Header.Get("x-csrf-Token")
			csrfCookie, err := r.Cookie("csrf")

			if err != nil || csrf == "" || csrfCookie.Value == "" || csrfCookie.Value != csrf {
				responses.SendErrorResponse(w, models.HTTPError{
					Code:    http.StatusInternalServerError,
					Message: "csrf-protection",
				}, logger.DripLogger.ErrorLogging)
				return
			}
			generateCsrfLogic(w)
			next.ServeHTTP(w, r)
		})

}
