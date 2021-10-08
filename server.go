package main

import (
	"fmt"
	"log"
	"net/http"
	"server/Handlers"
	"server/MockDB"
	"server/Models"
	"strings"
	"time"

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

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://ijia.me")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, PUT, OPTIONS")
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
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("LOG [%s] %s, %s %s",
			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))

		next.ServeHTTP(w, r)
	})
}

func PaincRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("Recovered from panic with err: %s on %s", err, r.RequestURI)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
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
	router.Use(Logger)
	router.Use(CORS)
	router.Use(PaincRecovery)

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
