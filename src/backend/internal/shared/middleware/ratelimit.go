package middleware

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

const (
	CodeRateLimited    = "rate_limited"
	MessageRateLimited = "too many requests, slow down"
)

const (
	rateLimitSweepFactor = 4
	minSweepInterval     = time.Minute
)

type RateLimitConfig struct {
	Rate       float64
	Burst      float64
	TrustProxy bool
	Now        func() time.Time
}

type bucket struct {
	tokens float64
	seen   time.Time
}

type RateLimiter struct {
	rate       float64
	burst      float64
	trustProxy bool
	now        func() time.Time
	ttl        time.Duration

	mu        sync.Mutex
	buckets   map[string]*bucket
	lastSweep time.Time
}

func NewRateLimiter(cfg RateLimitConfig) *RateLimiter {
	now := cfg.Now
	if now == nil {
		now = time.Now
	}

	ttl := minSweepInterval
	if cfg.Rate > 0 {
		if refill := time.Duration(cfg.Burst / cfg.Rate * float64(time.Second)); refill*rateLimitSweepFactor > ttl {
			ttl = refill * rateLimitSweepFactor
		}
	}

	return &RateLimiter{
		rate:       cfg.Rate,
		burst:      cfg.Burst,
		trustProxy: cfg.TrustProxy,
		now:        now,
		ttl:        ttl,
		buckets:    make(map[string]*bucket),
		lastSweep:  now(),
	}
}

func (l *RateLimiter) Allow(key string) bool {
	if l.rate <= 0 {
		return true
	}

	now := l.now()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.sweepLocked(now)

	b, ok := l.buckets[key]
	if !ok {
		b = &bucket{tokens: l.burst}
		l.buckets[key] = b
	} else {
		b.tokens += now.Sub(b.seen).Seconds() * l.rate
		if b.tokens > l.burst {
			b.tokens = l.burst
		}
	}

	b.seen = now

	if b.tokens < 1 {
		return false
	}

	b.tokens--

	return true
}

func (l *RateLimiter) sweepLocked(now time.Time) {
	if now.Sub(l.lastSweep) < l.ttl {
		return
	}

	for key, b := range l.buckets {
		if now.Sub(b.seen) > l.ttl {
			delete(l.buckets, key)
		}
	}

	l.lastSweep = now
}

func (l *RateLimiter) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if l.Allow(ClientIP(r, l.trustProxy)) {
				next.ServeHTTP(w, r)

				return
			}

			writeRateLimited(w, r, l.retryAfterSeconds())
		})
	}
}

func (l *RateLimiter) retryAfterSeconds() int {
	if l.rate <= 0 {
		return 1
	}

	seconds := int(1/l.rate + 0.999)
	if seconds < 1 {
		return 1
	}

	return seconds
}

func RateLimit(cfg RateLimitConfig) func(http.Handler) http.Handler {
	return NewRateLimiter(cfg).Middleware()
}

func ClientIP(r *http.Request, trustProxy bool) string {
	if trustProxy {
		if ip := firstForwardedIP(r.Header.Get("X-Forwarded-For")); ip != "" {
			return ip
		}

		if ip := normalizeIP(r.Header.Get("X-Real-IP")); ip != "" {
			return ip
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return normalizeIP(r.RemoteAddr)
	}

	return normalizeIP(host)
}

func firstForwardedIP(header string) string {
	for _, part := range strings.Split(header, ",") {
		if ip := normalizeIP(part); ip != "" {
			return ip
		}
	}

	return ""
}

func normalizeIP(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}

	if host, _, err := net.SplitHostPort(trimmed); err == nil {
		trimmed = host
	}

	parsed := net.ParseIP(strings.Trim(trimmed, "[]"))
	if parsed == nil {
		return ""
	}

	return parsed.String()
}

func writeRateLimited(w http.ResponseWriter, r *http.Request, retryAfter int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
	w.WriteHeader(http.StatusTooManyRequests)

	payload := map[string]map[string]string{
		"error": {
			"code":       CodeRateLimited,
			"message":    MessageRateLimited,
			"request_id": middleware.GetReqID(r.Context()),
		},
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.ErrorContext(r.Context(), "encode rate limit response", slog.Any("error", err))
	}
}
