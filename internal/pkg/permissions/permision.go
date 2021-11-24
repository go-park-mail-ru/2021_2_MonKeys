package permissions

import (
	"context"
	"dripapp/configs"
	_userModels "dripapp/internal/dripapp/models"
	_authClient "dripapp/internal/microservices/auth/delivery/grpc/client"
	_sessionModels "dripapp/internal/microservices/auth/models"
	"dripapp/internal/pkg/logger"
	"dripapp/internal/pkg/responses"

	"net/http"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

type Permission struct {
	AuthClient _authClient.SessionClient
}

func (perm *Permission) CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userSession _sessionModels.Session
		session, err := r.Cookie("sessionId")
		if err != nil {
			responses.SendError(w, _userModels.HTTPError{
				Code:    http.StatusForbidden,
				Message: _userModels.ErrAuth,
			}, logger.DripLogger.WarnLogging)
			return
		} else {
			userSession, err = perm.AuthClient.GetFromSession(r.Context(), session.Value)
			if err != nil {
				responses.SendError(w, _userModels.HTTPError{
					Code:    http.StatusForbidden,
					Message: _userModels.ErrAuth,
				}, logger.DripLogger.WarnLogging)
				return
			}
		}
		r = r.WithContext(context.WithValue(r.Context(), configs.ContextUserID, userSession))
		next.ServeHTTP(w, r)
	})
}

func (perm *Permission) GetCurrentUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		logger.DripLogger.DebugLogging("get current")
		ctxSession := r.Context().Value(configs.ContextUserID)
		if ctxSession == nil {
			responses.SendError(w, _userModels.HTTPError{
				Code:    http.StatusForbidden,
				Message: _userModels.ErrExtractContext,
			}, logger.DripLogger.ErrorLogging)
			return
		}
		currentSession, ok := ctxSession.(_sessionModels.Session)
		if !ok {
			responses.SendError(w, _userModels.HTTPError{
				Code:    http.StatusForbidden,
				Message: _userModels.ErrExtractContext,
			}, logger.DripLogger.ErrorLogging)
			return
		}

		currentUser, err := perm.AuthClient.GetById(r.Context(), currentSession)
		if err != nil {
			responses.SendError(w, _userModels.HTTPError{
				Code:    http.StatusNotFound,
				Message: err,
			}, logger.DripLogger.ErrorLogging)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), configs.ContextUser, currentUser))
		next.ServeHTTP(w, r)
	})
}

func generateCsrfLogic(w http.ResponseWriter) {
	csrf, err := uuid.NewV4()
	if err != nil {
		responses.SendError(w, _userModels.HTTPError{
			Code:    http.StatusForbidden,
			Message: _userModels.ErrNoPermission,
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
			// csrf := r.Header.Get("x-csrf-Token")
			// csrfCookie, err := r.Cookie("csrf")

			// if err != nil || csrf == "" || csrfCookie.Value == "" || csrfCookie.Value != csrf {
			// 	responses.SendError(w, models.HTTPError{
			// 		Code:    models.StatusCsrfProtection,
			// 		Message: models.ErrCSRF,
			// 	}, logger.DripLogger.ErrorLogging)
			// 	return
			// }
			// generateCsrfLogic(w)
			next.ServeHTTP(w, r)
		})

}
