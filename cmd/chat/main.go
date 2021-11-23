package main

import (
	"dripapp/configs"
	"dripapp/internal/dripapp/middleware"
	_sessionRepo "dripapp/internal/dripapp/session/repository"
	_userRepo "dripapp/internal/dripapp/user/repository"
	_chatDelivery "dripapp/internal/microservices/chat/delivery"
	"dripapp/internal/microservices/chat/models"
	_chatRepo "dripapp/internal/microservices/chat/repository"
	_chatUsecase "dripapp/internal/microservices/chat/usecase"
	"dripapp/internal/pkg/logger"
	"log"
	"net/http"
	"os"

	_ "dripapp/docs"

	"github.com/gorilla/mux"
)

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

	sm, err := _sessionRepo.NewTarantoolConnection(configs.Tarantool)
	if err != nil {
		log.Fatal(err)
	}

	timeoutContext := configs.Timeouts.ContextTimeout

	// chat
	// repository
	userRepo, err := _userRepo.NewPostgresUserRepository(configs.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	hub := models.NewHub()
	go hub.Run()
	chatRepo, err := _chatRepo.NewPostgresChatRepository(configs.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	// usecase
	chatUseCase := _chatUsecase.NewChatUseCase(chatRepo, hub, timeoutContext)

	// delivery
	_chatDelivery.SetChatRouting(logger.DripLogger, router, chatUseCase, userRepo)

	// middleware
	middleware.NewMiddleware(router, sm, logFile, logger.DripLogger)

	srv := &http.Server{
		Handler:      router,
		Addr:         configs.ChatServer.Port,
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	log.Printf("STD starting server at %s\n", srv.Addr)

	// for local
	log.Fatal(srv.ListenAndServe())
	// for deploy
	// log.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
}
