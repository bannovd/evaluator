package middleware

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// DefaultBuckets prometheus buckets (in seconds).
	DefaultBuckets = []float64{0.0001, 0.001, 0.01, 0.1, 0.5, 1, 1.5, 2}
)

const (
	reqsName    = "requests_total"
	latencyName = "request_duration_milliseconds"
)

// Prometheus is a handler that exposes prometheus metrics for the number of requests,
// the latency and the response size, partitioned by status code, method and HTTP path.
// Usage: pass its `ServeHTTP` to a route or globally.
type Prometheus struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

// NewPrometheus returns a new prometheus middleware.
// If buckets are empty then `DefaultBuckets` are setted.
func NewPrometheus(name string, buckets ...float64) *Prometheus {
	p := Prometheus{}

	p.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        reqsName,
			Help:        "How many HTTP requests processed",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{},
	)
	prometheus.MustRegister(p.reqs)

	if len(buckets) == 0 {
		buckets = DefaultBuckets
	}

	p.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        latencyName,
		Help:        "How long it took to process the request",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     buckets,
	},
		[]string{},
	)
	prometheus.MustRegister(p.latency)

	return &p
}

// PrometheusMiddleware handler
func (p *Prometheus) PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		defer func() {
			p.reqs.WithLabelValues().Add(1)
			p.latency.WithLabelValues().
				Observe(time.Since(now).Seconds())
		}()
		next.ServeHTTP(w, r)
	})
}
