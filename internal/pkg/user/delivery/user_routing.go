package delivery

import (
	"dripapp/internal/pkg/models"
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

	router.PathPrefix("/api/documentation/").Handler(httpSwagger.WrapHandler)
}
