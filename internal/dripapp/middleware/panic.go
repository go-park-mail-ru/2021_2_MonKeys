package middleware

import (
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/logger"
	"dripapp/internal/pkg/responses"
	"fmt"
	"net/http"
)

func PanicRecovery(l logger.Logger) (mw func(http.Handler) http.Handler) {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("Recovered from panic with err: %s on Method: [%s] %s\n", err, r.Method, r.RequestURI)
					responses.SendErrorResponse(w, models.InternalServerError500, l.ErrorLogging)
					return
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
