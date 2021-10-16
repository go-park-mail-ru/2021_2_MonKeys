package dripapp

import (
	"dripapp/internal/dripapp/middleware"
	"dripapp/internal/pkg/models"
	"dripapp/internal/pkg/session"
	_userDelivery "dripapp/internal/pkg/user/delivery"
	_userRepo "dripapp/internal/pkg/user/repository"
	_userUsecase "dripapp/internal/pkg/user/usecase"
	"log"
	"net/http"
	"os"

	_ "dripapp/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

const StatusEmailAlreadyExists = 1001
const (
	certFile = "api.ijia.me.crt"
	keyFile  = "api.ijia.me.key"
)

var (
	userRepo = _userRepo.NewMockDB()
)

func init() {
	userRepo.CreateUserAndProfile(&models.User{
		ID:          1,
		Name:        "Mikhail",
		Email:       "mumeu222@mail.ru",
		Password:    "af57966e1958f52e41550e822dd8e8a4", //VBif222!
		Date:        "2012-12-12",
		Age:         20,
		Description: "Hahahahaha",
		ImgSrc:      "/img/Yachty-tout.jpg",
		Tags:        []string{"soccer", "anime"},
	})
	userRepo.CreateUserAndProfile(&models.User{
		ID:          2,
		Name:        "Mikhail2",
		Email:       "mumeu222@mail.ru",
		Password:    "af57966e1958f52e41550e822dd8e8a4", //VBif222!
		Date:        "2012-12-12",
		Age:         20,
		Description: "Hahahahaha",
		ImgSrc:      "/img/Yachty-tout.jpg",
		Tags:        []string{"soccer", "anime"},
	})
	userRepo.CreateUserAndProfile(&models.User{
		ID:          3,
		Name:        "Mikhail3",
		Email:       "mumeu222@mail.ru",
		Password:    "af57966e1958f52e41550e822dd8e8a4", //VBif222!
		Date:        "2012-12-12",
		Age:         20,
		Description: "Hahahahaha",
		ImgSrc:      "/img/Yachty-tout.jpg",
		Tags:        []string{"soccer", "anime"},
	})
	userRepo.CreateUserAndProfile(&models.User{
		ID:          4,
		Name:        "Mikhail4",
		Email:       "mumeu222@mail.ru",
		Password:    "af57966e1958f52e41550e822dd8e8a4", //VBif222!
		Date:        "2012-12-12",
		Age:         20,
		Description: "Hahahahaha",
		ImgSrc:      "/img/Yachty-tout.jpg",
		Tags:        []string{"soccer", "anime"},
	})
	userRepo.CreateTag("anime")
	userRepo.CreateTag("netflix")
	userRepo.CreateTag("games")
	userRepo.CreateTag("walk")
	userRepo.CreateTag("JS")
	userRepo.CreateTag("baumanka")
	userRepo.CreateTag("music")
	userRepo.CreateTag("sport")

	// viper.SetConfigFile(`../../config.json`)
	// err := viper.ReadInConfig()
	// if err != nil {
	// 	panic(err)
	// }

	// if viper.GetBool(`debug`) {
	// 	log.Println("Service RUN on DEBUG mode")
	// }
}

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
	// conn, err := tarantool.Connect("127.0.0.1:3301", tarantool.Opts{
	// 	User: "admin",
	// 	Pass: "pass",
	// })
	// if err != nil {
	// 	log.Fatalf("Connection refused")
	// }
	// defer func(conn *tarantool.Connection) {
	// 	err := conn.Close()
	// 	if err != nil {

	// 	}
	// }(conn)

	// logfile
	logFile, err := os.OpenFile("../../logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {

		}
	}(logFile)

	// handlers
	// userRepo := _userRepo.NewMockDB()
	// sm, err := session.NewTarantoolConnection()
	sess := session.NewSessionDB()
	if err != nil {
		log.Fatal(err)
	}

	userUCase := _userUsecase.NewUserUsecase(userRepo, sess)

	userHandler := &_userDelivery.UserHandler{
		UserUCase: userUCase,
	}

	// router
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/session", userHandler.LoginHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/session", userHandler.LogoutHandler).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/api/v1/profile", userHandler.CurrentUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/profile", userHandler.EditProfileHandler).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/v1/profile", userHandler.SignupHandler).Methods("POST", "OPTIONS")

	router.HandleFunc("/api/v1/feed", userHandler.NextUserHandler).Methods("GET", "OPTIONS")

	router.HandleFunc("/api/v1/tags", userHandler.GetAllTags).Methods("GET", "OPTIONS")

	router.PathPrefix("/api/documentation/").Handler(httpSwagger.WrapHandler)

	// middleware
	router.Use(middleware.Logger(logFile))
	router.Use(middleware.CORS)
	router.Use(middleware.PanicRecovery)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8000",
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	log.Printf("STD starting server at %s\n", srv.Addr)

	log.Fatal(srv.ListenAndServe())
	// log.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
}
