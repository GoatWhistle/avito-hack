package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/middleware"
)

func TestMetricsPassesRequestThroughUnmodified(t *testing.T) {
	t.Parallel()

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("payload"))
	})

	rec := httptest.NewRecorder()
	middleware.Metrics()(next).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", http.NoBody))

	assert.Equal(t, http.StatusAccepted, rec.Code)
	assert.Equal(t, "payload", rec.Body.String())
}

func TestMetricsWithMatchedChiRoutePattern(t *testing.T) {
	t.Parallel()

	var seenPattern string

	router := chi.NewRouter()
	router.Use(middleware.Metrics())
	router.Get("/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		seenPattern = chi.RouteContext(r.Context()).RoutePattern()
		w.WriteHeader(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/items/abc", http.NoBody))

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "/items/{id}", seenPattern)
}

func TestMetricsWithoutRouteContextDoesNotPanic(t *testing.T) {
	t.Parallel()

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })

	rec := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		middleware.Metrics()(next).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/no-router", http.NoBody))
	})

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestMetricsWithEmptyRoutePattern(t *testing.T) {
	t.Parallel()

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))

	rec := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		middleware.Metrics()(next).ServeHTTP(rec, req)
	})

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestMetricsObservesUnmatchedRoutes(t *testing.T) {
	t.Parallel()

	router := chi.NewRouter()
	router.Use(middleware.Metrics())
	router.Get("/known", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/unknown", http.NoBody))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestMetricsHandlesConcurrentRequests(t *testing.T) {
	t.Parallel()

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	handler := middleware.Metrics()(next)

	done := make(chan struct{})

	for range 20 {
		go func() {
			defer func() { done <- struct{}{} }()

			handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", http.NoBody))
		}()
	}

	for range 20 {
		<-done
	}
}
