package delivery

import (
	"dripapp/internal/dripapp/models"
	_authClient "dripapp/internal/microservices/auth/delivery/grpc/client"
	_sessionModels "dripapp/internal/microservices/auth/models"
	"dripapp/internal/pkg/logger"
	_p "dripapp/internal/pkg/permissions"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetUserRouting(loggger logger.Logger, router *mux.Router, us models.UserUsecase, su _sessionModels.SessionUsecase, sc _authClient.SessionClient) {
	userHandler := &UserHandler{
		Logger:       loggger,
		UserUCase:    us,
		SessionUcase: su,
	}

	perm := _p.Permission{
		AuthClient: sc,
	}

	router.HandleFunc("/api/v1/profile", _p.CheckCSRF(perm.CheckAuth(perm.GetCurrentUser(userHandler.EditProfileHandler)))).Methods("PUT", "OPTIONS")

	router.HandleFunc("/api/v1/profile/photo", _p.CheckCSRF(perm.CheckAuth(perm.GetCurrentUser(userHandler.UploadPhoto)))).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/profile/photo", _p.CheckCSRF(perm.CheckAuth(perm.GetCurrentUser(userHandler.DeletePhoto)))).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/api/v1/user/cards", _p.SetCSRF(perm.CheckAuth(perm.GetCurrentUser(userHandler.NextUserHandler)))).Methods("GET", "OPTIONS")

	router.HandleFunc("/api/v1/user/likes", _p.SetCSRF(perm.CheckAuth(perm.GetCurrentUser(userHandler.LikesHandler)))).Methods("GET", "OPTIONS")

	router.HandleFunc("/api/v1/matches", _p.SetCSRF(perm.CheckAuth(perm.GetCurrentUser(userHandler.MatchesHandler)))).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/matches", _p.SetCSRF((perm.CheckAuth(perm.GetCurrentUser(userHandler.SearchMatchesHandler))))).Methods("POST", "OPTIONS")

	router.HandleFunc("/api/v1/likes", _p.CheckCSRF(perm.CheckAuth(perm.GetCurrentUser(userHandler.ReactionHandler)))).Methods("POST", "OPTIONS")

	router.HandleFunc("/api/v1/tags", _p.SetCSRF(perm.CheckAuth(userHandler.GetAllTags))).Methods("GET", "OPTIONS")

	router.HandleFunc("/api/v1/reports", _p.SetCSRF(perm.CheckAuth(userHandler.GetAllReports))).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/reports", _p.SetCSRF(perm.CheckAuth(perm.GetCurrentUser(userHandler.AddReport)))).Methods("POST", "OPTIONS")

	router.HandleFunc("/api/v1/payment/{id:[0-9]+}", _p.SetCSRF(perm.CheckAuth(perm.GetCurrentUser(userHandler.UpdatePayment)))).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/v1/payment", _p.SetCSRF(perm.CheckAuth(perm.GetCurrentUser(userHandler.CreatePayment)))).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/payment", _p.SetCSRF(perm.CheckAuth(perm.GetCurrentUser(userHandler.CheckPayment)))).Methods("GET", "OPTIONS")

	router.PathPrefix("/api/documentation/").Handler(httpSwagger.WrapHandler)
}
