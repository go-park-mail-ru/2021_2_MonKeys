package middleware

import (
	_sessionModels "dripapp/internal/microservices/auth/models"
	"dripapp/internal/pkg/logger"
	"dripapp/internal/pkg/monitoring"
	"os"

	"github.com/gorilla/mux"
)

func NewMiddleware(r *mux.Router, sp _sessionModels.SessionRepository, logFile *os.File, l logger.Logger) {

	metrics := monitoring.RegisterMetrics(r)

	r.Use(Logger(logFile, l, metrics))
	//r.Use(CORS(l))
	r.Use(PanicRecovery(l, metrics))
}
