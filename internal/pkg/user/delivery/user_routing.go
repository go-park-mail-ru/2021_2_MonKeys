package delivery

import (
	"dripapp/internal/pkg/models"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetRouting(router *mux.Router, us models.UserUsecase) {
	userHandler := &UserHandler{
		UserUCase: us,
	}

	router.HandleFunc("/api/v1/session", userHandler.LoginHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/session", userHandler.LogoutHandler).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/api/v1/profile", userHandler.CurrentUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/profile", userHandler.EditProfileHandler).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/v1/profile", userHandler.SignupHandler).Methods("POST", "OPTIONS")

	router.HandleFunc("/api/v1/feed", userHandler.NextUserHandler).Methods("GET", "OPTIONS")

	router.HandleFunc("/api/v1/tags", userHandler.GetAllTags).Methods("GET", "OPTIONS")

	router.PathPrefix("/api/documentation/").Handler(httpSwagger.WrapHandler)
}
