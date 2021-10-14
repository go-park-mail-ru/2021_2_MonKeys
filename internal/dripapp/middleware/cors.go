package middleware

import (
	"net/http"
	"strings"
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
