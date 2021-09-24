package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
	ID          uint64
	Name        string
	Email       string
	Password    string
	Age         uint64
	Description string
	ImgSrc      string
	Tags        []string
}

var (
	users   = make(map[uint64]User)
	cookies = make(map[string]uint64)
)

const (
	StatusBadRequest = 400
	StatusNotFound   = 404
	StatusOk         = 200
)

type JSON struct {
	Status uint64      `json:"status"`
	Body   interface{} `json:"body"`
}

type CurrentUserBody struct {
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Age         uint64   `json:"age"`
	Description string   `json:"description"`
	ImgSrc      string   `json:"imgSrc"`
	Tags        []string `json:"tags"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func cookieHandler(w http.ResponseWriter, r *http.Request) {
	var currentStatus uint64
	currentStatus = StatusNotFound
	var resp JSON

	session, err := r.Cookie("sessionId")
	if err == http.ErrNoCookie {
		currentStatus = StatusNotFound
	}
	if len(cookies) == 0 {
		currentStatus = StatusNotFound
	} else {
		currentUserId, okCookie := cookies[session.Value]
		if okCookie {
			currentUser, okUser := users[currentUserId]
			if !okUser {
				currentStatus = StatusNotFound
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
		}
	}

	resp.Status = currentStatus

	byteResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteResp)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var currentStatus uint64
	currentStatus = StatusNotFound
	var resp JSON

	byteReq, _ := ioutil.ReadAll(r.Body)
	strReq := string(byteReq)

	var logUserData LoginUser
	err := json.Unmarshal([]byte(strReq), &logUserData)
	if err != nil {
		// no valid json data
		currentStatus = StatusBadRequest
	}

	for _, value := range users {
		if value.Email == logUserData.Email && value.Password == logUserData.Password {
			currentStatus = StatusOk

			// create cookie
			expiration := time.Now().Add(10 * time.Hour)
			md5CookieValue := md5.Sum([]byte(logUserData.Email))
			cookie := http.Cookie{
				Name:     "sessionId",
				Value:    hex.EncodeToString(md5CookieValue[:]),
				Expires:  expiration,
				Secure:   true,
				HttpOnly: true,
			}

			cookies[hex.EncodeToString(md5CookieValue[:])] = value.ID

			http.SetCookie(w, &cookie)
		}
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

func main() {
	marvin := User{
		ID:          1,
		Name:        "Mikhail",
		Email:       "mumeu222@mail.ru",
		Password:    "VBif222!",
		Age:         20,
		Description: "Hahahahaha",
		ImgSrc:      "/static/users/user1",
		Tags:        []string{"haha", "hihi"},
	}
	users[1] = marvin

	mux := mux.NewRouter()

	mux.HandleFunc("/api/v1/cookie", cookieHandler).Methods("GET")
	mux.HandleFunc("/api/v1/login", loginHandler).Methods("POST")

	spa := spaHandler{staticPath: "static", indexPath: "index.html"}
	mux.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler:      mux,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	log.Fatal(srv.ListenAndServe())
}
