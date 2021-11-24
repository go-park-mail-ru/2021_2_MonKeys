package main

import (
	"context"
	"dripapp/configs"
	"dripapp/internal/dripapp/file"
	_fileDelivery "dripapp/internal/dripapp/file/delivery"
	"dripapp/internal/dripapp/middleware"
	"dripapp/internal/dripapp/models"
	_sessionDelivery "dripapp/internal/dripapp/session/delivery"
	_sessionRepo "dripapp/internal/dripapp/session/repository"
	_sessionUcase "dripapp/internal/dripapp/session/usecase"
	_userDelivery "dripapp/internal/dripapp/user/delivery"
	_userRepo "dripapp/internal/dripapp/user/repository"
	_userUsecase "dripapp/internal/dripapp/user/usecase"
	"dripapp/internal/pkg/hasher"
	"dripapp/internal/pkg/logger"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "dripapp/docs"

	"github.com/gorilla/mux"
)

func createUser(r models.UserRepository, f models.FileRepository, name string) uint64 {
	loginData := models.LoginUser{
		Email:    name + "@mail.ru",
		Password: hasher.HashAndSalt(nil, "qweQWE12"),
	}
	user, err := r.CreateUser(context.Background(), loginData)
	fmt.Println("CreateUser: ", err)

	err = f.CreateFoldersForNewUser(user)
	fmt.Println("CreateFoldersForNewUser: ", err)

	userData := models.User{
		ID:          user.ID,
		Email:       user.Email,
		Password:    user.Password,
		Name:        name,
		Gender:      "male",
		FromAge:     18,
		ToAge:       100,
		Date:        "2004-01-02",
		Description: "Всем привет меня зовут" + name,
		Imgs:        []string{"wsx.webp"},
	}
	fmt.Println("FillProfile: ", err)

	_, _ = r.UpdateUser(context.Background(), userData)

	return user.ID
}

func startRepo(r models.UserRepository, f models.FileRepository) {
	time.Sleep(3 * time.Second)

	userID1 := createUser(r, f, "Ilyagu")
	userID2 := createUser(r, f, "Vova")
	_ = createUser(r, f, "Misha")

	ctx := context.Background()
	r.AddReaction(ctx, userID1, userID2, 1)
	//r.AddReaction(ctx, userID2, userID1, 1)

	//fmt.Println("SendMessage: ", err)
}

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
	middleware.NewMiddleware(router, sm, logFile, logger.DripLogger)

	srv := &http.Server{
		Handler:      router,
		Addr:         configs.Server.Port,
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	log.Printf("STD starting server at %s\n", srv.Addr)

	go startRepo(userRepo, fileManager)
	// for local
	log.Fatal(srv.ListenAndServe())
	// for deploy
	// log.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
}
