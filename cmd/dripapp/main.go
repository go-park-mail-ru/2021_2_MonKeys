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
	_sessionDelivery "dripapp/internal/microservices/auth/delivery/http"
	_sessionRepo "dripapp/internal/microservices/auth/repository"
	_sessionUcase "dripapp/internal/microservices/auth/usecase"
	"dripapp/internal/pkg/logger"
	"log"
	"net/http"
	"os"

	_ "dripapp/docs"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

// func createUser(r models.UserRepository, f models.FileRepository, name string) uint64 {
// 	loginData := models.LoginUser{
// 		Email:    name + "@mail.ru",
// 		Password: hasher.HashAndSalt(nil, "qweQWE12"),
// 	}
// 	user, err := r.CreateUser(context.Background(), loginData)
// 	fmt.Println("CreateUser: ", err)

// 	err = f.CreateFoldersForNewUser(user)
// 	fmt.Println("CreateFoldersForNewUser: ", err)

// 	userData := models.User{
// 		ID:          user.ID,
// 		Email:       user.Email,
// 		Password:    user.Password,
// 		Name:        name,
// 		Gender:      "male",
// 		FromAge:     18,
// 		ToAge:       100,
// 		Date:        "2004-01-02",
// 		Description: "Всем привет меня зовут" + name,
// 		Imgs:        []string{"wsx.webp"},
// 	}
// 	fmt.Println("FillProfile: ", err)

// 	_, _ = r.UpdateUser(context.Background(), userData)

// 	return user.ID
// }

// func startRepo(r models.UserRepository, cr models.ChatRepository, f models.FileRepository) {
// 	time.Sleep(3 * time.Second)

// 	userID1 := createUser(r, f, "Ilyagu")
// 	userID2 := createUser(r, f, "Vova")
// 	userID3 := createUser(r, f, "Misha")

// 	// Message
// 	_, err := cr.SaveMessage(userID1, userID2, "")
// 	_, err = cr.SaveMessage(userID2, userID1, "")
// 	_, err = cr.SaveMessage(userID1, userID2, "AAAAAAAAA!")

// 	_, err = cr.SaveMessage(userID1, userID3, "")
// 	_, err = cr.SaveMessage(userID3, userID1, "")
// 	_, err = cr.SaveMessage(userID1, userID3, "первое сообщение!")
// 	fmt.Println("SendMessage: ", err)
// }

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

	// auth client
	grpcConn, _ := grpc.Dial(configs.AuthServer.GrpcUrl, grpc.WithInsecure())
	grpcAuthClient := _authClient.NewStaffClient(grpcConn)

	// delivery
	_userDelivery.SetUserRouting(logger.DripLogger, router, userUCase, sessionUcase, *grpcAuthClient)
	_sessionDelivery.SetSessionRouting(logger.DripLogger, router, userUCase, sessionUcase, *grpcAuthClient)
	_fileDelivery.SetFileRouting(router, *fileManager)

	// middleware
	middleware.NewMiddleware(router, sm, logFile, logger.DripLogger)

	srv := &http.Server{
		Handler:      router,
		Addr:         configs.Server.HttpPort,
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	log.Printf("STD starting server at %s\n", srv.Addr)

	// go startRepo(userRepo, chatRepo, fileManager)
	// for local
	log.Fatal(srv.ListenAndServe())
	// for deploy
	// log.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
}
