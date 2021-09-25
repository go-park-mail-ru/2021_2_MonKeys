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
	StatusBadRequest = 400
	StatusNotFound   = 404
	StatusOk         = 200
)

type JSON struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

type CurrentUserBody struct {
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Age         uint      `json:"age"`
	Description string   `json:"description"`
	ImgSrc      string   `json:"imgSrc"`
	Tags        []string `json:"tags"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (env *Env) cookieHandler(w http.ResponseWriter, r *http.Request) {
	currentStatus := StatusNotFound
	var resp JSON

	session, err := r.Cookie("sessionId")
	if err == http.ErrNoCookie {
		currentStatus = StatusNotFound
		return
	}

	currentUser, err := env.sessionDB.getUserByCookie(session.Value)
	if err != nil {
		currentStatus = StatusNotFound
		return
	}

	userBody := CurrentUserBody{
		currentUser.Name,
		currentUser.Email,
		currentUser.Age,
		currentUser.Description,
		currentUser.ImgSrc,
		currentUser.Tags,
	}

	currentStatus = StatusOk
	resp.Body = userBody

	resp.Status = currentStatus

	byteResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteResp)
}

func (env *Env) loginHandler(w http.ResponseWriter, r *http.Request) {
	currentStatus := StatusNotFound
	var resp JSON

	byteReq, _ := ioutil.ReadAll(r.Body)
	strReq := string(byteReq)

	var logUserData LoginUser
	err := json.Unmarshal([]byte(strReq), &logUserData)
	if err != nil {
		// no valid json data
		currentStatus = StatusBadRequest
	}

	identifiableUser, _ := env.db.getUserModel(logUserData.Email)

	if identifiableUser.Password == logUserData.Password {
		currentStatus = StatusOk

		// create cookie
		expiration := time.Now().Add(10 * time.Hour)
		//md5CookieValue := md5.Sum([]byte(logUserData.Email))
		md5CookieValue := fmt.Sprintf("%x", md5.Sum([]byte(logUserData.Email + logUserData.Password)))
		cookie := http.Cookie{
			Name:     "sessionId",
			Value:    md5CookieValue,
			Expires:  expiration,
			Secure:   true,
			HttpOnly: true,
		}

		env.sessionDB.newSessionCookie(identifiableUser.ID, md5CookieValue)

		http.SetCookie(w, &cookie)
	}

	resp.Status = currentStatus
	byteResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteResp)
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
		newSessionCookie(uint64, string) error
	}
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
