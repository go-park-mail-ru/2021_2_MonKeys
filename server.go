package main

import (
	"log"
	"net/http"
	"strings"
	"server/Handlers"
	"server/MockDB"
	"server/Models"
	"time"

	_ "server/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

const StatusEmailAlreadyExists = 1001
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
		var sb strings.Builder
		sb.WriteString("Accept,")
		sb.WriteString("Content-Type,")
		sb.WriteString("Content-Length,")
		sb.WriteString("Accept-Encoding,")
		sb.WriteString("X-CSRF-Token,")
		sb.WriteString("Authorization,")
		sb.WriteString("Allow-Credentials,")
		sb.WriteString("Set-Cookie,")
		sb.WriteString("Access-Control-Allow-Credentials,")
		sb.WriteString("Access-Control-Allow-Origin")
		w.Header().Set("Access-Control-Allow-Headers", sb.String())
		next.ServeHTTP(w, r)

		log.Printf("LOG [%s] %s, %s %s",
			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
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
  
  router.HandleFunc("/api/v1/profile", env.currentUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/auth", env.loginHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/profile", env.editProfileHandler).Methods("PATCH", "OPTIONS")
	router.HandleFunc("/api/v1/profile", env.signupHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/auth", env.logoutHandler).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/v1/feed", env.nextUserHandler).Methods("GET", "OPTIONS")
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

	log.Printf("STD starting server at %s\n", srv.Addr)

	log.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
}
