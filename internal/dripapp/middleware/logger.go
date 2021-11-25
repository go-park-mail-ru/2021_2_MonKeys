package middleware

import (
	"dripapp/internal/pkg/logger"
	"dripapp/internal/pkg/monitoring"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Logger(logFile *os.File, l logger.Logger, metrics *monitoring.PromMetrics) (mw func(http.Handler) http.Handler) {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			multiWriter := io.MultiWriter(os.Stdout, logFile)
			log.SetOutput(multiWriter)
			l.InfoLogging(r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
			reqTime := time.Now()
			next.ServeHTTP(w, r)
			respTime := time.Since(reqTime)
			if r.URL.Path != "/metrics" {
				metrics.Hits.WithLabelValues(strconv.Itoa(http.StatusOK), r.URL.String(), r.Method).Inc()
				metrics.Timings.WithLabelValues(
					strconv.Itoa(http.StatusOK), r.URL.String(), r.Method).Observe(respTime.Seconds())
			}
		})
	}
}
