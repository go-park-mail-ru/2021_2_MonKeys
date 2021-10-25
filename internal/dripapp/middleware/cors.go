package middleware

import (
	"log"
	"net/http"
	"strings"
)

var allowedOrigins = map[string]struct{}{
	"http://127.0.0.1": {},
	"http://localhost": {},
	"http://ijia.me":   {},

	"https://127.0.0.1": {},
	"https://localhost": {},
	"https://ijia.me":   {},
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		_, isIn := allowedOrigins[origin]
		if isIn {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			log.Println("unknown origin", `"`+origin+`"`)
			http.Error(w, "Access denied", http.StatusForbidden)
		}
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
		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
