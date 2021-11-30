package main

import (
	"dripapp/configs"
	"dripapp/internal/dripapp/file"
	_fileDelivery "dripapp/internal/dripapp/file/delivery"
	"dripapp/internal/dripapp/middleware"
	_userDelivery "dripapp/internal/dripapp/user/delivery"
	_userRepo "dripapp/internal/dripapp/user/repository"
	_userUsecase "dripapp/internal/dripapp/user/usecase"
	_authClient "dripapp/internal/microservices/auth/delivery/grpc/client"
	_sessionRepo "dripapp/internal/microservices/auth/repository"
	_sessionUcase "dripapp/internal/microservices/auth/usecase"
	_chatDelivery "dripapp/internal/microservices/chat/delivery"
	"dripapp/internal/microservices/chat/models"
	_chatRepo "dripapp/internal/microservices/chat/repository"
	_chatUsecase "dripapp/internal/microservices/chat/usecase"
	"log"

	"dripapp/internal/pkg/logger"
	"net/http"
	"os"

	_ "dripapp/docs"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
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
	hub := models.NewHub()
	go hub.Run()
	chatRepo, err := _chatRepo.NewPostgresChatRepository(configs.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	// usecase
	chatUseCase := _chatUsecase.NewChatUseCase(chatRepo, hub, timeoutContext)

	sessionUcase := _sessionUcase.NewSessionUsecase(sm, timeoutContext)
	userUCase := _userUsecase.NewUserUsecase(
		userRepo,
		fileManager,
		timeoutContext,
	)

	// auth client
	grpcConn, _ := grpc.Dial(configs.AuthServer.GrpcUrl, grpc.WithInsecure())
	grpcAuthClient := _authClient.NewAuthClient(grpcConn)

	// delivery
	_userDelivery.SetUserRouting(logger.DripLogger, router, userUCase, sessionUcase, *grpcAuthClient)
	_fileDelivery.SetFileRouting(router, *fileManager)
	_chatDelivery.SetChatRouting(logger.DripLogger, router, chatUseCase, *grpcAuthClient)

	// middleware
	middleware.NewMiddleware(router, sm, logFile, logger.DripLogger)

	srv := &http.Server{
		Handler:      router,
		Addr:         configs.Server.HttpPort,
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	log.Printf("STD starting server at %s\n", srv.Addr)

	// for local
	// log.Fatal(srv.ListenAndServe())
	// for deploy
	log.Fatal(srv.ListenAndServeTLS(configs.Server.CertFile, configs.Server.KeyFile))
}
