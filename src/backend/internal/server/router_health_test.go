package server_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/config"
	"github.com/avito-hack/backend/internal/server"
)

type stubModule struct {
	path   string
	called *bool
}

func (m stubModule) RegisterRoutes(r chi.Router) {
	r.Get(m.path, func(w http.ResponseWriter, _ *http.Request) {
		*m.called = true
		w.WriteHeader(http.StatusOK)
	})
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func minimalRouter(t *testing.T, modules ...server.ModuleRegistrar) http.Handler {
	t.Helper()

	upgrade := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusSwitchingProtocols)
	})

	return server.NewRouter(server.RouterDeps{
		Config:    config.Config{AllowedOrigins: []string{"https://app.example.com"}},
		Logger:    discardLogger(),
		Modules:   modules,
		WebSocket: upgrade,
	})
}

func TestNewRouterHealthzIsOK(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	minimalRouter(t).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", http.NoBody))

	require.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
}

func TestNewRouterUnknownRouteIs404(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	minimalRouter(t).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/nope", http.NoBody))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestNewRouterExposesMetrics(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	minimalRouter(t).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", http.NoBody))

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "go_goroutines")
}

func TestNewRouterMountsModulesUnderAPIV1(t *testing.T) {
	t.Parallel()

	called := false

	rec := httptest.NewRecorder()
	minimalRouter(t, stubModule{path: "/things", called: &called}).
		ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/things", http.NoBody))

	require.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, called)
}

func TestNewRouterMountsWebSocketHandler(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	minimalRouter(t).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/ws", http.NoBody))

	assert.Equal(t, http.StatusSwitchingProtocols, rec.Code)
}

func TestNewRouterAppliesCORS(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/healthz", http.NoBody)
	req.Header.Set("Origin", "https://app.example.com")

	rec := httptest.NewRecorder()
	minimalRouter(t).ServeHTTP(rec, req)

	assert.Equal(t, "https://app.example.com", rec.Header().Get("Access-Control-Allow-Origin"))
}

func TestNewRouterAssignsRequestID(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/healthz", http.NoBody)
	req.Header.Set("X-Request-Id", "req-777")

	rec := httptest.NewRecorder()
	minimalRouter(t).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestNewRouterRecoversPanickingModule(t *testing.T) {
	t.Parallel()

	handler := server.NewRouter(server.RouterDeps{
		Config:  config.Config{},
		Logger:  discardLogger(),
		Modules: []server.ModuleRegistrar{panicModule{}},
	})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/boom", http.NoBody))

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "internal_error")
}

type panicModule struct{}

func (panicModule) RegisterRoutes(r chi.Router) {
	r.Get("/boom", func(http.ResponseWriter, *http.Request) { panic("module exploded") })
}

func TestNewRouterWithoutUploadConfigSkipsUploads(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	minimalRouter(t).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/uploads/anything.jpg", http.NoBody))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestReadyzWithoutPoolIsUnavailable(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		minimalRouter(t).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/readyz", http.NoBody))
	})

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.JSONEq(t, `{"status":"unavailable"}`, rec.Body.String())
}
