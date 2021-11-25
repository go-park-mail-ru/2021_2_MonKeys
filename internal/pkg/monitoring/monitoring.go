package monitoring

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PromMetrics struct {
	Hits    *prometheus.CounterVec
	Timings *prometheus.HistogramVec
}

func RegisterMetrics(r *mux.Router) *PromMetrics {
	var metrics PromMetrics

	metrics.Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hits",
	}, []string{"status", "path", "method"})

	metrics.Timings = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "timings",
		},
		[]string{"status", "path", "method"},
	)

	prometheus.MustRegister(metrics.Hits, metrics.Timings)

	r.Handle("/metrics", promhttp.Handler())

	return &metrics
}
