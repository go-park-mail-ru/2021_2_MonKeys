package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"server/MockDB"
	"server/Models"
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

func sendResp(resp Models.JSON, w http.ResponseWriter) {
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

func createSessionCookie(user Models.LoginUser) http.Cookie {
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
	var resp Models.JSON
	session, err := r.Cookie("sessionId")
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		return
	}

	currentUser, err := env.getUserByCookie(session.Value)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
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
	var resp Models.JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		return
	}

	var logUserData Models.LoginUser
	err = json.Unmarshal(byteReq, &logUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		return
	}

	identifiableUser, err := env.db.GetUser(logUserData.Email)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		return
	}

	status := StatusOK
	if identifiableUser.IsCorrectPassword(logUserData.Password) {
		cookie := createSessionCookie(logUserData)
		err = env.sessionDB.NewSessionCookie(cookie.Value, identifiableUser.ID)
		if err != nil {
			resp.Status = StatusInternalServerError
			sendResp(resp, w)
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
	var resp Models.JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		return
	}

	var logUserData Models.LoginUser
	err = json.Unmarshal(byteReq, &logUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		return
	}

	identifiableUser, _ := env.db.GetUser(logUserData.Email)
	if !identifiableUser.IsEmpty() {
		resp.Status = StatusEmailAlreadyExists
		sendResp(resp, w)
		return
	}

	user, err := env.db.CreateUser(logUserData)
	if err != nil {
		resp.Status = StatusInternalServerError
		sendResp(resp, w)
		return
	}

	cookie := createSessionCookie(logUserData)
	err = env.sessionDB.NewSessionCookie(cookie.Value, user.ID)
	if err != nil {
		resp.Status = StatusInternalServerError
		sendResp(resp, w)
		return
	}

	http.SetCookie(w, &cookie)

	resp.Status = StatusOK
	sendResp(resp, w)
}

func (env *Env) editProfileHandler(w http.ResponseWriter, r *http.Request) {
	var resp Models.JSON
	session, err := r.Cookie("sessionId")
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		return
	}

	currentUser, err := env.getUserByCookie(session.Value)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		return
	}

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		return
	}

	var newUserData Models.User
	err = json.Unmarshal(byteReq, &newUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		return
	}

	err = currentUser.FillProfile(newUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		return
	}

	err = env.db.UpdateUser(currentUser)
	if err != nil {
		resp.Status = StatusInternalServerError
		sendResp(resp, w)
		return
	}

	resp.Status = StatusOK
	resp.Body = currentUser

	sendResp(resp, w)
}

func (env *Env) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("sessionId")
	if err != nil {
		sendResp(Models.JSON{Status: StatusNotFound}, w)
		return
	}

	err = env.sessionDB.DeleteSessionCookie(session.Value)
	if err != nil {
		sendResp(Models.JSON{Status: StatusInternalServerError}, w)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
}

func (env *Env) nextUserHandler(w http.ResponseWriter, r *http.Request) {
	var resp Models.JSON

	// get current user by cookie
	session, err := r.Cookie("sessionId")
	if err == http.ErrNoCookie {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		return
	}
	currentUser, err := env.getUserByCookie(session.Value)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		return
	}

	// get swiped usedata for registrationr id from json
	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		return
	}
	var swipedUserData Models.SwipedUser
	err = json.Unmarshal(byteReq, &swipedUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		return
	}

	// add in swaped users map for current user
	err = env.db.AddSwipedUsers(currentUser.ID, swipedUserData.Id)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		return
	}
	// find next user for swipe
	nextUser, err := env.db.GetNextUserForSwipe(currentUser.ID)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, w)
		return
	}

	resp.Status = StatusOK
	resp.Body = nextUser

	sendResp(resp, w)
}

type Env struct {
	db interface {
		GetUser(email string) (Models.User, error)
		GetUserByID(userID uint64) (Models.User, error)
		CreateUser(logUserData Models.LoginUser) (Models.User, error)
		AddSwipedUsers(currentUserId, swipedUserId uint64) error
		GetNextUserForSwipe(currentUserId uint64) (Models.User, error)
		UpdateUser(newUserData Models.User) error
	}
	sessionDB interface {
		GetUserIDByCookie(sessionCookie string) (uint64, error)
		NewSessionCookie(sessionCookie string, userId uint64) error
		DeleteSessionCookie(sessionCookie string) error
	}
}

func (env Env) getUserByCookie(sessionCookie string) (Models.User, error) {
	userID, err := env.sessionDB.GetUserIDByCookie(sessionCookie)
	if err != nil {
		return Models.User{}, errors.New("error sessionDB: GetUserIDByCookie")
	}

	user, err := env.db.GetUserByID(userID)
	if err != nil {
		return Models.User{}, errors.New("error db: getUserByID")
	}

	return user, nil
}

var (
	db = MockDB.NewMockDB()
)

/*func init() {
	db.
	db.users[1] = Models.User{
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
	db.users[2] = Models.User{
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
	db.users[3] = Models.User{
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
	db.users[4] = Models.User{
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
}*/

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
		sessionDB: MockDB.NewSessionDB(),
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
