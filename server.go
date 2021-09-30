package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "server/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusInternalServerError = 500
	StatusEmailAlreadyExists  = 1001
)

func sendResp(resp JSON, w http.ResponseWriter) {
	byteResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteResp)
}

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

func createSessionCookie(user LoginUser) http.Cookie {
	expiration := time.Now().Add(10 * time.Hour)

	data := user.Password + time.Now().String()
	md5CookieValue := fmt.Sprintf("%x", md5.Sum([]byte(data)))

	cookie := http.Cookie{
		Name:     "sessionId",
		Value:    md5CookieValue,
		Expires:  expiration,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	return cookie
}

func (env *Env) currentUser(w http.ResponseWriter, r *http.Request) {
	var resp JSON
	session, err := r.Cookie("sessionId")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	currentUser, err := env.getUserByCookie(session.Value)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp.Status = StatusOK
	resp.Body = currentUser

	sendResp(resp, w)
}

// @Summary LogIn
// @Description log in
// @Tags login
// @Accept json
// @Produce json
// @Param input body LoginUser true "data for login"
// @Success 200 {object} JSON
// @Failure 400,404,500
// @Router /login [post]
func (env *Env) loginHandler(w http.ResponseWriter, r *http.Request) {
	var resp JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var logUserData LoginUser
	err = json.Unmarshal(byteReq, &logUserData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	identifiableUser, err := env.db.getUser(logUserData.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	status := StatusOK
	if identifiableUser.isCorrectPassword(logUserData.Password) {
		cookie := createSessionCookie(logUserData)
		err = env.sessionDB.newSessionCookie(cookie.Value, identifiableUser.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &cookie)

		resp.Body = identifiableUser
	} else {
		status = StatusNotFound
	}

	resp.Status = status
	sendResp(resp, w)
}

// @Summary SignUp
// @Description registration user
// @Tags registration
// @Accept json
// @Produce json
// @Param input body LoginUser true "data for registration"
// @Success 200 {object} JSON
// @Failure 400,404,500
// @Router /signup [post]
func (env *Env) signupHandler(w http.ResponseWriter, r *http.Request) {
	var resp JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var logUserData LoginUser
	err = json.Unmarshal(byteReq, &logUserData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	identifiableUser, _ := env.db.getUser(logUserData.Email)
	if !identifiableUser.isEmpty() {
		resp.Status = StatusEmailAlreadyExists
		sendResp(resp, w)
		return
	}

	user, err := env.db.createUser(logUserData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := createSessionCookie(logUserData)
	err = env.sessionDB.newSessionCookie(cookie.Value, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &cookie)

	resp.Status = StatusOK
	sendResp(resp, w)
}

func (env *Env) editProfileHandler(w http.ResponseWriter, r *http.Request) {
	var resp JSON
	session, err := r.Cookie("sessionId")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	currentUser, err := env.getUserByCookie(session.Value)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var newUserData User
	err = json.Unmarshal(byteReq, &newUserData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = currentUser.fillProfile(newUserData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = env.db.updateUser(currentUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.Status = StatusOK
	resp.Body = currentUser

	sendResp(resp, w)
}

func (env *Env) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("sessionId")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = env.sessionDB.deleteSessionCookie(session.Value)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
}

func (env *Env) nextUserHandler(w http.ResponseWriter, r *http.Request) {
	var resp JSON

	// get current user by cookie
	session, err := r.Cookie("sessionId")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	currentUser, err := env.getUserByCookie(session.Value)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// get swiped usedata for registrationr id from json
	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var swipedUserData SwipedUser
	err = json.Unmarshal(byteReq, &swipedUserData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// add in swaped users map for current user
	err = env.db.addSwipedUsers(currentUser.ID, swipedUserData.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// find next user for swipe
	nextUser, err := env.db.getNextUserForSwipe(currentUser.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp.Status = StatusOK
	resp.Body = nextUser

	sendResp(resp, w)
}

type Env struct {
	db interface {
		getUser(email string) (User, error)
		getUserByID(userID uint64) (User, error)
		createUser(logUserData LoginUser) (User, error)
		addSwipedUsers(currentUserId, swipedUserId uint64) error
		getNextUserForSwipe(currentUserId uint64) (User, error)
		updateUser(newUserData User) error
	}
	sessionDB interface {
		getUserIDByCookie(sessionCookie string) (uint64, error)
		newSessionCookie(sessionCookie string, userId uint64) error
		deleteSessionCookie(sessionCookie string) error
	}
}

func (env Env) getUserByCookie(sessionCookie string) (User, error) {
	userID, err := env.sessionDB.getUserIDByCookie(sessionCookie)
	if err != nil {
		return User{}, errors.New("error sessionDB: getUserIDByCookie")
	}

	user, err := env.db.getUserByID(userID)
	if err != nil {
		return User{}, errors.New("error db: getUserByID")
	}

	return user, nil
}

var (
	db = NewMockDB()
)

func init() {
	db.users[1] = User{
		ID:          1,
		Name:        "Mikhail",
		Email:       "mumeu222@mail.ru",
		Password:    "af57966e1958f52e41550e822dd8e8a4", //VBif222!
		Date:        "2012-12-12",
		Age:         20,
		Description: "Hahahahaha",
		ImgSrc:      "/img/Yachty-tout.jpg",
		Tags:        []string{"soccer", "anime"},
	}
	db.users[2] = User{
		ID:          2,
		Name:        "Mikhail2",
		Email:       "mumeu222@mail.ru",
		Password:    "af57966e1958f52e41550e822dd8e8a4", //VBif222!
		Date:        "2012-12-12",
		Age:         20,
		Description: "Hahahahaha",
		ImgSrc:      "/img/Yachty-tout.jpg",
		Tags:        []string{"soccer", "anime"},
	}
	db.users[3] = User{
		ID:          3,
		Name:        "Mikhail3",
		Email:       "mumeu222@mail.ru",
		Password:    "af57966e1958f52e41550e822dd8e8a4", //VBif222!
		Date:        "2012-12-12",
		Age:         20,
		Description: "Hahahahaha",
		ImgSrc:      "/img/Yachty-tout.jpg",
		Tags:        []string{"soccer", "anime"},
	}
	db.users[4] = User{
		ID:          4,
		Name:        "Mikhail4",
		Email:       "mumeu222@mail.ru",
		Password:    "af57966e1958f52e41550e822dd8e8a4", //VBif222!
		Date:        "2012-12-12",
		Age:         20,
		Description: "Hahahahaha",
		ImgSrc:      "/img/Yachty-tout.jpg",
		Tags:        []string{"soccer", "anime"},
	}
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
	env := &Env{
		db:        db, // NewMockDB()
		sessionDB: NewSessionDB(),
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/currentuser", env.currentUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/login", env.loginHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/editprofile", env.editProfileHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/signup", env.signupHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/logout", env.logoutHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/nextswipeuser", env.nextUserHandler).Methods("POST", "OPTIONS")
	router.Use(CORSMiddleware)

	router.PathPrefix("/api/documentation/").Handler(httpSwagger.WrapHandler)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":443",
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	log.Fatal(srv.ListenAndServeTLS("api.ijia.me.crt", "api.ijia.me.key"))
}
