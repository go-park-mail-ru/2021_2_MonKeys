package main

import (
	"fmt"
	"net/http"
	"time"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_id")
	loggedIn := (err != http.ErrNoCookie)

	if loggedIn {
		w.Write(page)
	} else {
		fmt.Fprintln(w, `<a href="/login">login</a>`)
		fmt.Fprintln(w, "You need to login")
	}
}

var loginFormTmpl = []byte(`
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link href="dist\css\main.css" rel="stylesheet">
    <title>Feed</title>
</head>
	<body>
	<form method="post">
		Login: <input type="text" name="login">
		Password: <input type="password" name="password">
		<input type="submit" value="Login">
	</form>
	<script type="text/javascript" src="data/sample.js"></script>
	</body>
</html>
`)
var page = []byte(`
<html>
	<body>
	<h1>Заголовок</h1>
	  <!-- Комментарий -->
	  <p>Первый абзац.</p>
	  <p>Второй абзац.</p>
	</body>
</html>
`)

func loginPage(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_id")
	loggedIn := (err != http.ErrNoCookie) // TODO проверка куки, запрос в бд

	if loggedIn {
		//http.Redirect(w, r, "/", http.StatusFound)
	}

	if r.Method != http.MethodPost {
		w.Write(loginFormTmpl)
		return
	}

	r.ParseForm()
	inputLogin := r.FormValue("login")
	inputPassword := r.FormValue("password")

	if inputLogin == "qwe" && inputPassword == "1" {
		expiration := time.Now().Add(10 * time.Hour)
		cookie := http.Cookie{
			Name:    "session_id",
			Value:   inputLogin,
			Expires: expiration,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/", http.StatusFound)
	}

	fmt.Fprintln(w, "<h1>jopa!!!</h1> tvoi dannie gavno")
}

func logoutPage(w http.ResponseWriter, r *http.Request) {
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
	mux := http.NewServeMux()
	mux.HandleFunc("/login", loginPage)
	mux.HandleFunc("/logout", logoutPage)
	mux.HandleFunc("/", mainPage)
	staticHandler := http.StripPrefix(
		"/data/",
		http.FileServer(http.Dir("./static")),
	)
	mux.Handle("/data/", staticHandler)

	server := http.Server{
		Addr:         ":8001",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("starting server at :8001")
	server.ListenAndServe()
}
