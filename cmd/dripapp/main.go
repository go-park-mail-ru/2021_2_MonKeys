package main

import (
	"dripapp/configs"
	"dripapp/internal/dripapp/middleware"
	"dripapp/internal/pkg/session"
	_userDelivery "dripapp/internal/pkg/user/delivery"
	_userRepo "dripapp/internal/pkg/user/repository"
	_userUsecase "dripapp/internal/pkg/user/usecase"
	"log"
	"net/http"
	"os"

	_ "dripapp/docs"

	"github.com/gorilla/mux"
)

const StatusEmailAlreadyExists = 1001

// @title Drip API
// @version 1.0
// @description API for Drip.
// @termsOfService http://swagger.io/terms/

// @host api.ijia.me
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Set-Cookie
func main() {
	// logfile
	logFile, err := os.OpenFile("../../logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(logFile)

	configs.SetConfig()

	// handlers
	userRepo := _userRepo.NewMockDB()
	userRepo.MockDB()
	// sm := session.NewSessionDB()
	sm, err := session.NewTarantoolConnection(configs.Tarantool)
	if err != nil {
		log.Fatal(err)
	}

	timeoutContext := configs.Timeouts.ContextTimeout
	userUCase := _userUsecase.NewUserUsecase(userRepo, sm, timeoutContext)

	// router
	router := mux.NewRouter()

	// middleware
	router.Use(middleware.Logger(logFile))
	router.Use(middleware.CORS)
	router.Use(middleware.PanicRecovery)

	_userDelivery.SetRouting(router, userUCase)

	srv := &http.Server{
		Handler:      router,
		Addr:         configs.Server.Port,
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	log.Printf("STD starting server at %s\n", srv.Addr)

	log.Fatal(srv.ListenAndServe())
	// log.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
}