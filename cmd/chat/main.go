package main

import (
	"dripapp/configs"
	"dripapp/internal/dripapp/middleware"
	_authClient "dripapp/internal/microservices/auth/delivery/grpc/client"
	_sessionRepo "dripapp/internal/microservices/auth/repository"
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

	sm, err := _sessionRepo.NewTarantoolConnection(configs.Tarantool)
	if err != nil {
		log.Fatal(err)
	}

	timeoutContext := configs.Timeouts.ContextTimeout

	// chat
	// repository
	hub := models.NewHub()
	go hub.Run()
	chatRepo, err := _chatRepo.NewPostgresChatRepository(configs.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	// usecase
	chatUseCase := _chatUsecase.NewChatUseCase(chatRepo, hub, timeoutContext)

	// auth client
	grpcConn, _ := grpc.Dial(configs.AuthServer.GrpcUrl, grpc.WithInsecure())
	grpcAuthClient := _authClient.NewAuthClient(grpcConn)

	// delivery
	_chatDelivery.SetChatRouting(logger.DripLogger, router, chatUseCase, *grpcAuthClient)

	// middleware
	middleware.NewMiddleware(router, sm, logFile, logger.DripLogger)

	srv := &http.Server{
		Handler:      router,
		Addr:         configs.ChatServer.HttpPort,
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}


	mode:= os.Getenv("DRIPAPP")
	log.Printf(mode);
	if mode=="LOCAL" {
		log.Fatal(srv.ListenAndServe())
	} else if mode=="DEPLOY" {
		log.Fatal(srv.ListenAndServeTLS("star.monkeys.team.crt", "star.monkeys.team.key"))
	} else {
		log.Printf("NO MODE SPECIFIED.SET ENV VAR DRIPAPP TO \"LOCAL\" or \"DEPLOY\"")
	}
	
	log.Printf("STD starting server(%s) at %s\n",mode, srv.Addr)
}
