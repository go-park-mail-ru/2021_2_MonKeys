package middleware

import (
	"dripapp/internal/pkg/logger"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func Logger(logFile *os.File, l logger.Logger) (mw func(http.Handler) http.Handler) {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			multiWriter := io.MultiWriter(os.Stdout, logFile)
			log.SetOutput(multiWriter)
			l.InfoLogging(r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))

			next.ServeHTTP(w, r)
		})
	}
}
