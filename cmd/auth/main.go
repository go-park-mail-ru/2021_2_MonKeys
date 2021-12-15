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
	"fmt"
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
	userUCase := _userUCase.NewUserUsecase(userRepo, fileManager, timeoutContext, nil)

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


	mode:= os.Getenv("DRIPAPP")

	if mode=="LOCAL" {
		log.Fatal(srv.ListenAndServe())
	} else if mode=="DEPLOY" {
		log.Fatal(srv.ListenAndServeTLS("star.monkeys.team.crt", "star.monkeys.team.key"))
	} else {
		log.Printf("NO MODE SPECIFIED.SET ENV VAR DRIPAPP TO \"LOCAL\" or \"DEPLOY\"")
	}

	fmt.Printf("STD starting server(%s) at %s\n",mode, srv.Addr)

	// for local

	// for deploy
	// log.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
}
