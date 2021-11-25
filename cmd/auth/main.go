package main

import (
	"dripapp/configs"
	"dripapp/internal/dripapp/file"
	"dripapp/internal/dripapp/middleware"
	_userRepo "dripapp/internal/dripapp/user/repository"
	_userUCase "dripapp/internal/dripapp/user/usecase"
	_authClient "dripapp/internal/microservices/auth/delivery/grpc/client"
	grpcServer "dripapp/internal/microservices/auth/delivery/grpc/grpc_server"
	_sessionDelivery "dripapp/internal/microservices/auth/delivery/http"
	_sessionRepo "dripapp/internal/microservices/auth/repository"
	_sessionUCase "dripapp/internal/microservices/auth/usecase"
	"dripapp/internal/pkg/logger"
	"log"
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

	timeoutContext := configs.Timeouts.ContextTimeout

	// repository
	sm, err := _sessionRepo.NewTarantoolConnection(configs.Tarantool)
	if err != nil {
		log.Fatal(err)
	}
	userRepo, err := _userRepo.NewPostgresUserRepository(configs.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	fileManager, err := file.NewFileManager(configs.FileStorage)
	if err != nil {
		log.Fatal(err)
	}

	// usecase
	sessionUCase := _sessionUCase.NewSessionUsecase(sm, timeoutContext)
	userUCase := _userUCase.NewUserUsecase(userRepo, fileManager, timeoutContext)

	// new auth server
	go grpcServer.StartAuthGrpcServer(sm, userRepo, configs.AuthServer.GrpcUrl)

	// auth client
	grpcConn, _ := grpc.Dial(configs.AuthServer.GrpcUrl, grpc.WithInsecure())
	grpcAuthClient := _authClient.NewAuthClient(grpcConn)

	// delivery
	_sessionDelivery.SetSessionRouting(logger.DripLogger, router, userUCase, sessionUCase, *grpcAuthClient)

	// middleware
	middleware.NewMiddleware(router, sm, logFile, logger.DripLogger)

	srv := &http.Server{
		Handler:      router,
		Addr:         configs.AuthServer.HttpPort,
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	log.Printf("STD starting server at %s\n", srv.Addr)

	// for local
	log.Fatal(srv.ListenAndServe())
	// for deploy
	// log.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
}
