package middleware

import (
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/logger"
	"dripapp/internal/pkg/responses"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var allowedOrigins = map[string]struct{}{
	"http://127.0.0.1":         {},
	"http://127.0.0.1:8000":    {},
	"http://localhost":         {},
	"http://localhost:8080":    {},
	"http://ijia.me":           {},
	"http://192.168.1.16:8080": {},

	"https://localhost:8080":    {},
	"https://127.0.0.1:8000":    {},
	"https://192.168.1.16:8080": {},
	"https://127.0.0.1":         {},
	"https://localhost":         {},
	"https://ijia.me":           {},
}

func CORS(logger logger.Logger) (mw func(http.Handler) http.Handler) {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			_, isIn := allowedOrigins[origin]
			if isIn {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				mess := fmt.Sprintf("unknown origin %s", origin)
				responses.SendErrorResponse(w, models.HTTPError{
					Code:    http.StatusInternalServerError,
					Message: errors.New(mess),
				}, logger.ErrorLogging)
				return
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
}
