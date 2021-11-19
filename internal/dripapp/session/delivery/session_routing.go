package delivery

import (
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/logger"
	"dripapp/internal/pkg/permissions"

	"github.com/gorilla/mux"
)

func SetSessionRouting(loggger logger.Logger, router *mux.Router, us models.UserUsecase, su models.SessionUsecase) {
	sessionHandler := &SessionHandler{
		Logger:       loggger,
		UserUCase:    us,
		SessionUcase: su,
	}

	router.HandleFunc("/api/v1/session", permissions.SetCSRF(sessionHandler.LoginHandler)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/session", permissions.CheckAuthenticated(sessionHandler.LogoutHandler)).Methods("DELETE", "OPTIONS")
}
