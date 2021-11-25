package middleware

import (
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/logger"
	"dripapp/internal/pkg/monitoring"
	"dripapp/internal/pkg/responses"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func PanicRecovery(l logger.Logger, metrics *monitoring.PromMetrics) (mw func(http.Handler) http.Handler) {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqTime := time.Now()
			defer func() {
				if err := recover(); err != nil {
					respTime := time.Since(reqTime)
					metrics.Hits.WithLabelValues(
						strconv.Itoa(http.StatusInternalServerError), r.URL.Path, r.Method).Inc()
					fmt.Printf("Recovered from panic with err: %s on Method: [%s] %s\n", err, r.Method, r.RequestURI)
					metrics.Timings.WithLabelValues(
						strconv.Itoa(http.StatusInternalServerError), r.URL.String(),
						r.Method).Observe(respTime.Seconds())
					responses.SendError(w, models.InternalServerError500, l.ErrorLogging)
					return
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
