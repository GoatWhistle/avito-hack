package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/middleware"
)

func TestRateLimiterAllowsUpToBurst(t *testing.T) {
	t.Parallel()

	limiter := middleware.NewRateLimiter(middleware.RateLimitConfig{
		Rate:  5,
		Burst: 3,
		Now:   func() time.Time { return time.Unix(0, 0) },
	})

	require.True(t, limiter.Allow("a"))
	require.True(t, limiter.Allow("a"))
	require.True(t, limiter.Allow("a"))
	require.False(t, limiter.Allow("a"))
}

func TestRateLimiterIsolatesKeys(t *testing.T) {
	t.Parallel()

	limiter := middleware.NewRateLimiter(middleware.RateLimitConfig{
		Rate:  1,
		Burst: 1,
		Now:   func() time.Time { return time.Unix(0, 0) },
	})

	require.True(t, limiter.Allow("a"))
	require.False(t, limiter.Allow("a"))
	require.True(t, limiter.Allow("b"))
}

func TestRateLimiterRefillsOverTime(t *testing.T) {
	t.Parallel()

	current := time.Unix(0, 0)
	limiter := middleware.NewRateLimiter(middleware.RateLimitConfig{
		Rate:  5,
		Burst: 1,
		Now:   func() time.Time { return current },
	})

	require.True(t, limiter.Allow("a"))
	require.False(t, limiter.Allow("a"))

	current = current.Add(200 * time.Millisecond)
	require.True(t, limiter.Allow("a"))
}

func TestRateLimiterDoesNotExceedBurstAfterIdle(t *testing.T) {
	t.Parallel()

	current := time.Unix(0, 0)
	limiter := middleware.NewRateLimiter(middleware.RateLimitConfig{
		Rate:  10,
		Burst: 2,
		Now:   func() time.Time { return current },
	})

	current = current.Add(time.Hour)

	require.True(t, limiter.Allow("a"))
	require.True(t, limiter.Allow("a"))
	require.False(t, limiter.Allow("a"))
}

func TestRateLimiterDisabledWhenRateZero(t *testing.T) {
	t.Parallel()

	limiter := middleware.NewRateLimiter(middleware.RateLimitConfig{Rate: 0, Burst: 0})

	for range 100 {
		require.True(t, limiter.Allow("a"))
	}
}

func TestRateLimiterConcurrentAccessIsSafe(t *testing.T) {
	t.Parallel()

	limiter := middleware.NewRateLimiter(middleware.RateLimitConfig{Rate: 1000, Burst: 1000})

	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()

			for range 20 {
				limiter.Allow(string(rune('a' + n%5)))
			}
		}(i)
	}

	wg.Wait()
}

func TestRateLimitMiddlewareReturns429WithEnvelope(t *testing.T) {
	t.Parallel()

	handler := middleware.RateLimit(middleware.RateLimitConfig{
		Rate:  1,
		Burst: 1,
		Now:   func() time.Time { return time.Unix(0, 0) },
	})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	first := doRateLimited(handler, "203.0.113.5:1234", nil)
	require.Equal(t, http.StatusOK, first.Code)

	second := doRateLimited(handler, "203.0.113.5:1234", nil)
	require.Equal(t, http.StatusTooManyRequests, second.Code)
	require.NotEmpty(t, second.Header().Get("Retry-After"))

	var payload struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	require.NoError(t, json.Unmarshal(second.Body.Bytes(), &payload))
	require.Equal(t, middleware.CodeRateLimited, payload.Error.Code)
	require.NotEmpty(t, payload.Error.Message)
}

func TestRateLimitMiddlewareSeparatesRemoteAddrs(t *testing.T) {
	t.Parallel()

	handler := middleware.RateLimit(middleware.RateLimitConfig{
		Rate:  1,
		Burst: 1,
		Now:   func() time.Time { return time.Unix(0, 0) },
	})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	require.Equal(t, http.StatusOK, doRateLimited(handler, "203.0.113.5:1111", nil).Code)
	require.Equal(t, http.StatusTooManyRequests, doRateLimited(handler, "203.0.113.5:2222", nil).Code)
	require.Equal(t, http.StatusOK, doRateLimited(handler, "203.0.113.9:1111", nil).Code)
}

func TestRateLimitMiddlewareUsesForwardedIPWhenTrusted(t *testing.T) {
	t.Parallel()

	handler := middleware.RateLimit(middleware.RateLimitConfig{
		Rate:       1,
		Burst:      1,
		TrustProxy: true,
		Now:        func() time.Time { return time.Unix(0, 0) },
	})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	proxied := map[string]string{"X-Forwarded-For": "198.51.100.7, 10.0.0.1"}
	require.Equal(t, http.StatusOK, doRateLimited(handler, "10.0.0.1:9999", proxied).Code)
	require.Equal(t, http.StatusTooManyRequests, doRateLimited(handler, "10.0.0.1:9999", proxied).Code)

	other := map[string]string{"X-Forwarded-For": "198.51.100.8"}
	require.Equal(t, http.StatusOK, doRateLimited(handler, "10.0.0.1:9999", other).Code)
}

func TestRateLimitMiddlewareIgnoresForwardedIPWhenUntrusted(t *testing.T) {
	t.Parallel()

	handler := middleware.RateLimit(middleware.RateLimitConfig{
		Rate:  1,
		Burst: 1,
		Now:   func() time.Time { return time.Unix(0, 0) },
	})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	require.Equal(t, http.StatusOK,
		doRateLimited(handler, "10.0.0.1:9999", map[string]string{"X-Forwarded-For": "198.51.100.7"}).Code)
	require.Equal(t, http.StatusTooManyRequests,
		doRateLimited(handler, "10.0.0.1:9999", map[string]string{"X-Forwarded-For": "198.51.100.8"}).Code)
}

func doRateLimited(handler http.Handler, remoteAddr string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", http.NoBody)
	req.RemoteAddr = remoteAddr

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	return rec
}
