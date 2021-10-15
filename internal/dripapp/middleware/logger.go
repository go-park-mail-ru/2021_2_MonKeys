package middleware

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func Logger(logFile *os.File) (mw func(http.Handler) http.Handler) {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			mw := io.MultiWriter(os.Stdout, logFile)
			log.SetOutput(mw)
			log.Printf("LOG [%s] %s, %s %s",
				r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))

			next.ServeHTTP(w, r)
		})
	}
}
