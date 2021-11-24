package http

import (
	_userModels "dripapp/internal/dripapp/models"
	_authClient "dripapp/internal/microservices/auth/delivery/grpc/client"
	_sessionModels "dripapp/internal/microservices/auth/models"
	"dripapp/internal/pkg/logger"
	_p "dripapp/internal/pkg/permissions"

	"github.com/gorilla/mux"
)

func SetSessionRouting(loggger logger.Logger, router *mux.Router, us _userModels.UserUsecase, su _sessionModels.SessionUsecase, sc _authClient.SessionClient) {
	sessionHandler := &SessionHandler{
		Logger:       loggger,
		UserUCase:    us,
		SessionUcase: su,
	}

	perm := _p.Permission{
		AuthClient: sc,
	}

	router.HandleFunc("/api/v1/session", _p.SetCSRF(sessionHandler.LoginHandler)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/session", perm.CheckAuth(sessionHandler.LogoutHandler)).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/api/v1/profile", _p.SetCSRF(perm.CheckAuth(perm.GetCurrentUser(sessionHandler.CurrentUser)))).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/profile", _p.SetCSRF(sessionHandler.SignupHandler)).Methods("POST", "OPTIONS")
}
