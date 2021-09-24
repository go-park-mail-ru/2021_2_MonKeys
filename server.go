package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

// type User struct {
// 	ID          int
// 	Name        string
// 	Age         int
// 	Description string
// 	Img         string
// 	Tags        []string
// }

var (
	users   = make(map[string]string)
	cookies = make(map[string]string)
)

type StatusLogedInJSON struct {
	Status string `json:"status"` // status 400 200
	// Body   interface{}
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (statusJSON *StatusLogedInJSON) ChangeStatus() {
	// 200 404 not found
	if statusJSON.Status == "error" {
		statusJSON.Status = "ok"
	} else {
		statusJSON.Status = "error"
	}
}

func homePageHandler(rw http.ResponseWriter, r *http.Request) {
	// _, err := r.Cookie("session_id")
	// loggedIn := (err != http.ErrNoCookie)

	t, _ := template.ParseFiles("static/main.html")
	_ = t.Execute(rw, nil)
}

// func cookieHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "GET" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// }

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// перенести в роутер
	// if r.Method != "POST" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	fmt.Println(cookies["session_id"])

	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		fmt.Println(11111)
		return
	}
	if cookies[session.Name] == session.Value {
		m := StatusLogedInJSON{"ok"}
		b, err := json.Marshal(m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	} else if r.Method == "POST" {
		b, _ := ioutil.ReadAll(r.Body)
		jsn := string(b)
		var logUserData LoginUser
		err = json.Unmarshal([]byte(jsn), &logUserData)
		if err != nil {
			// error 400 in json
			fmt.Println(12494)
			fmt.Println(err)
		}
		m := StatusLogedInJSON{"error"}
		for key, value := range users {
			if key == logUserData.Email && value == logUserData.Password {
				m.ChangeStatus()
			}
		}
		b, err = json.Marshal(m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if m.Status == "ok" {
			// cookie
			expiration := time.Now().Add(10 * time.Hour)
			cookie := http.Cookie{
				Name:     "session_id",
				Value:    logUserData.Email,
				Expires:  expiration,
				Secure:   true,
				HttpOnly: true,
			}
			// r.Cookie()
			// fmt.Println(r.Cookie("session_id"))
			cookies["session_id"] = logUserData.Email

			http.SetCookie(w, &cookie)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}

}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	http.Redirect(w, r, "/", http.StatusFound)
}

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
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
	users["mumeu222@mail.ru"] = "VBif222!"

	mux := mux.NewRouter()

	mux.HandleFunc("/api/v1/login", loginHandler)
	// mux.HandleFunc("/logout", logoutHandler)

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
