package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/middleware"
)

func okHandler(called *bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		*called = true
		w.WriteHeader(http.StatusOK)
	})
}

func TestCORSAllowsConfiguredOrigin(t *testing.T) {
	t.Parallel()

	var called bool

	req := httptest.NewRequest(http.MethodGet, "/api/v1/items", http.NoBody)
	req.Header.Set("Origin", "https://app.example.com")

	rec := httptest.NewRecorder()
	middleware.CORS([]string{"https://app.example.com"})(okHandler(&called)).ServeHTTP(rec, req)

	require.True(t, called)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "https://app.example.com", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))
	assert.Contains(t, rec.Header().Get("Access-Control-Expose-Headers"), "X-Request-Id")
}

func TestCORSDisallowedOriginGetsNoAllowHeader(t *testing.T) {
	t.Parallel()

	var called bool

	req := httptest.NewRequest(http.MethodGet, "/api/v1/items", http.NoBody)
	req.Header.Set("Origin", "https://evil.example.com")

	rec := httptest.NewRecorder()
	middleware.CORS([]string{"https://app.example.com"})(okHandler(&called)).ServeHTTP(rec, req)

	require.True(t, called)
	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSPreflightShortCircuits(t *testing.T) {
	t.Parallel()

	var called bool

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/items", http.NoBody)
	req.Header.Set("Origin", "https://app.example.com")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	req.Header.Set("Access-Control-Request-Headers", "Authorization,Content-Type")

	rec := httptest.NewRecorder()
	middleware.CORS([]string{"https://app.example.com"})(okHandler(&called)).ServeHTTP(rec, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "300", rec.Header().Get("Access-Control-Max-Age"))
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), http.MethodPost)
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Headers"), "Authorization")
}

func TestCORSPreflightAllowsEveryConfiguredMethod(t *testing.T) {
	t.Parallel()

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPatch,
		http.MethodPut,
		http.MethodDelete,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			t.Parallel()

			var called bool

			req := httptest.NewRequest(http.MethodOptions, "/api/v1/items", http.NoBody)
			req.Header.Set("Origin", "https://app.example.com")
			req.Header.Set("Access-Control-Request-Method", method)

			rec := httptest.NewRecorder()
			middleware.CORS([]string{"https://app.example.com"})(okHandler(&called)).ServeHTTP(rec, req)

			assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), method)
		})
	}
}

func TestCORSEmptyOriginsFallsBackToLocalhost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		origins []string
	}{
		{name: "nil slice", origins: nil},
		{name: "empty slice", origins: []string{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var called bool

			req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			req.Header.Set("Origin", "http://localhost:3000")

			rec := httptest.NewRecorder()
			middleware.CORS(tc.origins)(okHandler(&called)).ServeHTTP(rec, req)

			require.True(t, called)
			assert.Equal(t, "http://localhost:3000", rec.Header().Get("Access-Control-Allow-Origin"))
		})
	}
}

func TestCORSWithoutOriginHeaderPassesThrough(t *testing.T) {
	t.Parallel()

	var called bool

	rec := httptest.NewRecorder()
	middleware.CORS(nil)(okHandler(&called)).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", http.NoBody))

	require.True(t, called)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"))
}
