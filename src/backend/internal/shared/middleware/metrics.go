package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request latency distribution",
		Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
	}, []string{"method", "route", "status"})

	requestsInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "Number of HTTP requests currently being served",
	})
)

func Metrics() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := chimw.NewWrapResponseWriter(w, r.ProtoMajor)

			requestsInFlight.Inc()
			defer requestsInFlight.Dec()

			next.ServeHTTP(wrapped, r)

			requestDuration.WithLabelValues(
				r.Method,
				routePattern(r),
				strconv.Itoa(wrapped.Status()),
			).Observe(time.Since(start).Seconds())
		})
	}
}

func routePattern(r *http.Request) string {
	ctx := chi.RouteContext(r.Context())
	if ctx == nil {
		return "unknown"
	}

	pattern := ctx.RoutePattern()
	if pattern == "" {
		return "unknown"
	}

	return pattern
}
