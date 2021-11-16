package main

import (
	"dripapp/configs"
	"dripapp/internal/dripapp/file"
	_fileDelivery "dripapp/internal/dripapp/file/delivery"
	"dripapp/internal/dripapp/middleware"
	_sessionDelivery "dripapp/internal/dripapp/session/delivery"
	_sessionRepo "dripapp/internal/dripapp/session/repository"
	_sessionUcase "dripapp/internal/dripapp/session/usecase"
	_userDelivery "dripapp/internal/dripapp/user/delivery"
	_userRepo "dripapp/internal/dripapp/user/repository"
	_userUsecase "dripapp/internal/dripapp/user/usecase"
	"dripapp/internal/pkg/logger"
	"log"
	"net/http"
	"os"

	_ "dripapp/docs"

	"github.com/gorilla/mux"
)

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
	logFile, err := os.OpenFile("logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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

	// router
	router := mux.NewRouter()

	// repository
	userRepo, err := _userRepo.NewPostgresUserRepository(configs.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	sm, err := _sessionRepo.NewTarantoolConnection(configs.Tarantool)
	if err != nil {
		log.Fatal(err)
	}

	fileManager, err := file.NewFileManager(configs.FileStorage)
	if err != nil {
		log.Fatal(err)
	}

	timeoutContext := configs.Timeouts.ContextTimeout

	// usecase
	sessionUcase := _sessionUcase.NewSessionUsecase(sm, timeoutContext)
	userUCase := _userUsecase.NewUserUsecase(
		userRepo,
		fileManager,
		timeoutContext,
	)

	// delivery
	_userDelivery.SetUserRouting(logger.DripLogger, router, userUCase, sessionUcase, userRepo)
	_sessionDelivery.SetSessionRouting(logger.DripLogger, router, userUCase, sessionUcase)
	_fileDelivery.SetFileRouting(router, *fileManager)

	// middleware
	middleware.NewMiddleware(router, sm, logFile)

	srv := &http.Server{
		Handler:      router,
		Addr:         configs.Server.Port,
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	log.Printf("STD starting server at %s\n", srv.Addr)

	// for local
	log.Fatal(srv.ListenAndServe())
	// for deploy
	// log.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
}
