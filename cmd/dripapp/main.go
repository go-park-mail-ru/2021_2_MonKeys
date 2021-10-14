package main

import (
	"log"
	"net/http"
	"server/Handlers"
	"server/MockDB"
	"server/Models"
	"server/internal/dripapp/middleware"

	_ "server/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/tarantool/go-tarantool"
)

const StatusEmailAlreadyExists = 1001
const (
	certFile = "api.ijia.me.crt"
	keyFile  = "api.ijia.me.key"
)

var (
	db = MockDB.NewMockDB()
)

func init() {
	db.CreateUserAndProfile(Models.User{
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
	db.CreateUserAndProfile(Models.User{
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
	db.CreateUserAndProfile(Models.User{
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
	db.CreateUserAndProfile(Models.User{
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
	db.CreateTag("anime")
	db.CreateTag("netflix")
	db.CreateTag("games")
	db.CreateTag("walk")
	db.CreateTag("JS")
	db.CreateTag("baumanka")
	db.CreateTag("music")
	db.CreateTag("sport")
}

func Router(env *Handlers.Env) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/profile", env.CurrentUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/session", env.LoginHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/profile", env.EditProfileHandler).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/v1/profile", env.SignupHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/session", env.LogoutHandler).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/v1/feed", env.NextUserHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/tags", env.GetAllTags).Methods("GET", "OPTIONS")

	// middleware
	router.Use(middleware.Logger)
	router.Use(middleware.CORS)
	router.Use(middleware.PanicRecovery)

	router.PathPrefix("/api/documentation/").Handler(httpSwagger.WrapHandler)

	return router
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

	conn, err := tarantool.Connect("127.0.0.1:3301", tarantool.Opts{
		User: "admin",
		Pass: "pass",
	})
	if err != nil {
		log.Fatalf("Connection refused")
	}
	defer conn.Close()

	env := &Handlers.Env{
		DB:        db, // NewMockDB()
		SessionDB: MockDB.NewSessionDB(),
	}

	router := Router(env)

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
