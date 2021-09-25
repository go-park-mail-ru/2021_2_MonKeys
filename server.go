package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

const (
	StatusOK         = 200
	StatusBadRequest = 400
	StatusNotFound   = 404
	StatusInternalServerError = 500
)

func sendResp(resp JSON, w *http.ResponseWriter) {
	byteResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
	}
	(*w).WriteHeader(http.StatusOK)
	(*w).Write(byteResp)
}

func (env *Env) cookieHandler(w http.ResponseWriter, r *http.Request) {
	var resp JSON

	session, err := r.Cookie("sessionId")
	if err == http.ErrNoCookie {
		resp.Status = StatusNotFound
		sendResp(resp, &w)
		return
	}

	currentUser, err := env.sessionDB.getUserByCookie(session.Value)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, &w)
		return
	}

	resp.Status = StatusOK
	resp.Body = currentUser

	sendResp(resp, &w)
}

func (env *Env) loginHandler(w http.ResponseWriter, r *http.Request) {
	var resp JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, &w)
		return
	}

	var logUserData LoginUser
	err = json.Unmarshal(byteReq, &logUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, &w)
		return
	}

	identifiableUser, err := env.db.getUserModel(logUserData.Email)
	if err != nil {
		resp.Status = StatusInternalServerError
		sendResp(resp, &w)
		return
	}

	status := StatusOK
	if identifiableUser.isCorrectPassword(logUserData.Password) {
		expiration := time.Now().Add(10 * time.Hour)

		data := logUserData.Password + time.Now().String()
		md5CookieValue := fmt.Sprintf("%x", md5.Sum([]byte(data)))

		cookie := http.Cookie{
			Name:     "sessionId",
			Value:    md5CookieValue,
			Expires:  expiration,
			Secure:   true,
			HttpOnly: true,
		}

		err = env.sessionDB.newSessionCookie(md5CookieValue, identifiableUser.ID)
		if err != nil {
			resp.Status = StatusInternalServerError
			sendResp(resp, &w)
			return
		}

		http.SetCookie(w, &cookie)
	} else {
		status = StatusNotFound
	}

	resp.Status = status
	sendResp(resp, &w)
}

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

type Env struct {
	db interface {
		getUserModel(string) (User, error)
	}
	sessionDB interface {
		getUserByCookie(string) (User, error)
		newSessionCookie(string, uint64) error
	}
}

func init() {
	marvin := User{
		ID:          1,
		Name:        "Mikhail",
		Email:       "mumeu222@mail.ru",
		Password:    "VBif222!",
		Age:         20,
		Description: "Hahahahaha",
		ImgSrc:      "/img/Yachty-tout.jpg",
		Tags:        []string{"haha", "hihi"},
	}
	users[1] = marvin
}

func main() {
	/*db, err := sql.Open("postgres", "postgres://user:pass@localhost/bookstore")
	if err != nil {
		log.Fatal(err)
	}

	env := &Env{
		db: ModelsDB{DB: db},
	}
	*/

	env := &Env{
		db:        MockDB{},
		sessionDB: MockSessionDB{},
	}

	mux := mux.NewRouter()

	mux.HandleFunc("/api/v1/cookie", env.cookieHandler).Methods("GET")
	mux.HandleFunc("/api/v1/login", env.loginHandler).Methods("POST")

	spa := spaHandler{staticPath: "static", indexPath: "index.html"}
	mux.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler:      mux,
		Addr:         ":8080",
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	log.Fatal(srv.ListenAndServe())
}
