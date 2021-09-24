package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
	"time"
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
	// body
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

func main() {
	users["mumeu222@mail.ru"] = "VBif222!"

	mux := http.NewServeMux()
	mux.HandleFunc("/", homePageHandler)
	mux.HandleFunc("/api/v1/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)
	staticHandler := http.StripPrefix(
		"/data/",
		http.FileServer(http.Dir("./static")),
	)
	mux.Handle("/data/", staticHandler)

	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  http.DefaultClient.Timeout,
		WriteTimeout: http.DefaultClient.Timeout,
	}

	fmt.Println("starting server at :8080")
	server.ListenAndServe()
}
