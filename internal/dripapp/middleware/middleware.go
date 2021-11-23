package middleware

import (
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/logger"
	"dripapp/internal/pkg/monitoring"
	"os"

	"github.com/gorilla/mux"
)

func NewMiddleware(r *mux.Router, sp models.SessionRepository, logFile *os.File, l logger.Logger) {
	sm := sessionMiddleware{
		sessionRepo: sp,
	}
	metrics := monitoring.RegisterMetrics(r)

	r.Use(Logger(logFile, l,metrics))
	//r.Use(CORS(l))
	r.Use(PanicRecovery(l,metrics))
	r.Use(sm.SessionMiddleware(l))
}
