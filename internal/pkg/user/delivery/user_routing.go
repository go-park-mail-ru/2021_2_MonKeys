package delivery

import (
	"dripapp/internal/pkg/models"
	"dripapp/internal/pkg/permissions"
	_sessionDelivery "dripapp/internal/pkg/session/delivery"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetRouting(router *mux.Router, us models.UserUsecase, su models.SessionUsecase) {
	userHandler := &UserHandler{
		UserUCase:    us,
		SessionUcase: su,
	}
	sessionHandler := &_sessionDelivery.SessionHandler{
		UserUCase:    us,
		SessionUcase: su,
	}

	router.HandleFunc("/api/v1/session",
		permissions.SetCSRF(sessionHandler.LoginHandler)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/session",
		permissions.CheckCSRF(permissions.CheckAuthenticated(sessionHandler.LogoutHandler))).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/api/v1/profile",
		permissions.CheckCSRF(permissions.CheckAuthenticated(userHandler.CurrentUser))).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/profile",
		permissions.CheckCSRF(permissions.CheckAuthenticated(userHandler.EditProfileHandler))).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/v1/profile",
		permissions.SetCSRF(userHandler.SignupHandler)).Methods("POST", "OPTIONS")

	router.HandleFunc("/api/v1/profile/photo",
		permissions.CheckAuthenticated(userHandler.UploadPhoto)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/profile/photo",
		permissions.CheckAuthenticated(userHandler.DeletePhoto)).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/api/v1/user/cards",
		permissions.CheckCSRF(permissions.CheckAuthenticated(userHandler.NextUserHandler))).Methods("GET", "OPTIONS")

	router.HandleFunc("/api/v1/matches",
		permissions.CheckCSRF(permissions.CheckAuthenticated(userHandler.MatchesHandler))).Methods("GET", "OPTIONS")

	router.HandleFunc("/api/v1/tags",
		permissions.CheckCSRF(permissions.CheckAuthenticated(userHandler.GetAllTags))).Methods("GET", "OPTIONS")

	router.PathPrefix("/api/documentation/").Handler(httpSwagger.WrapHandler)

	/*
		router.HandleFunc("/api/v1/session", sessionHandler.LoginHandler).Methods("POST", "OPTIONS")
		router.HandleFunc("/api/v1/session", sessionHandler.LogoutHandler).Methods("DELETE", "OPTIONS")

		router.HandleFunc("/api/v1/profile", userHandler.CurrentUser).Methods("GET", "OPTIONS")
		router.HandleFunc("/api/v1/profile", userHandler.EditProfileHandler).Methods("PUT", "OPTIONS")
		router.HandleFunc("/api/v1/profile", userHandler.SignupHandler).Methods("POST", "OPTIONS")

		router.HandleFunc("/api/v1/profile/photo", userHandler.UploadPhoto).Methods("POST", "OPTIONS")
		router.HandleFunc("/api/v1/profile/photo", userHandler.DeletePhoto).Methods("DELETE", "OPTIONS")

		router.HandleFunc("/api/v1/user/cards", userHandler.NextUserHandler).Methods("GET", "OPTIONS")

		router.HandleFunc("/api/v1/matches", userHandler.MatchesHandler).Methods("GET", "OPTIONS")

		router.HandleFunc("/api/v1/tags", userHandler.GetAllTags).Methods("GET", "OPTIONS")

		router.PathPrefix("/api/documentation/").Handler(httpSwagger.WrapHandler)]
	*/
}
