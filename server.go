package main

import (
	"log"
	"net/http"
	"server/Handlers"
	"server/MockDB"
	"server/Models"

	_ "server/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	certFile = "api.ijia.me.crt"
	keyFile  = "api.ijia.me.key"
)

var (
	db = MockDB.NewMockDB()
)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://ijia.me")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept,"+
				"Content-Type,"+
				"Content-Length,"+
				"Accept-Encoding,"+
				"X-CSRF-Token,"+
				"Authorization,"+
				"Allow-Credentials,"+
				"Set-Cookie,"+
				"Access-Control-Allow-Credentials,"+
				"Access-Control-Allow-Origin")
		next.ServeHTTP(w, r)
	})
}

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
}

func Router(env *Handlers.Env) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/currentuser", env.CurrentUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/login", env.LoginHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/editprofile", env.EditProfileHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/signup", env.SignupHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/logout", env.LogoutHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/nextswipeuser", env.NextUserHandler).Methods("POST", "OPTIONS")
	router.Use(CORSMiddleware)

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
	env := &Handlers.Env{
		DB:        db, // NewMockDB()
		SessionDB: MockDB.NewSessionDB(),
	}

	router := Router(env)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":443",
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	log.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
}
