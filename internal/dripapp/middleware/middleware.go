package middleware

import (
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/logger"
	"os"

	"github.com/gorilla/mux"
)

func NewMiddleware(r *mux.Router, sp models.SessionRepository, logFile *os.File, l logger.Logger) {
	sm := sessionMiddleware{
		sessionRepo: sp,
	}
	r.Use(Logger(logFile, l))
	//r.Use(CORS(Logger))
	r.Use(PanicRecovery(l))
	r.Use(sm.SessionMiddleware(l))
}
