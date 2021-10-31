package middleware

import (
	"dripapp/internal/dripapp/models"
	"os"

	"github.com/gorilla/mux"
)

func NewMiddleware(r *mux.Router, sp models.SessionRepository, logFile *os.File) {
	sm := sessionMiddleware{
		sessionRepo: sp,
	}
	r.Use(Logger(logFile))
	r.Use(CORS)
	r.Use(PanicRecovery)
	r.Use(sm.SessionMiddleware)
}
