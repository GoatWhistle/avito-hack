package middleware_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/middleware"
)

func bufLogger() (*slog.Logger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	handler := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})

	return slog.New(handler), buf
}

func decodeLog(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))

	return entry
}

func TestAccessLogRecordsRequestFacts(t *testing.T) {
	t.Parallel()

	log, buf := bufLogger()

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("hello"))
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/items", http.NoBody)
	req.RemoteAddr = "10.0.0.7:54321"

	middleware.AccessLog(log)(next).ServeHTTP(httptest.NewRecorder(), req)

	entry := decodeLog(t, buf)
	assert.Equal(t, "http request", entry["msg"])
	assert.Equal(t, http.MethodPost, entry["method"])
	assert.Equal(t, "/api/v1/items", entry["path"])
	assert.InDelta(t, http.StatusCreated, entry["status"], 0)
	assert.InDelta(t, 5, entry["bytes"], 0)
	assert.Equal(t, "10.0.0.7:54321", entry["remote_ip"])
	assert.Contains(t, entry, "duration")
}

func TestAccessLogDefaultsStatusToOKWhenHandlerWritesBodyOnly(t *testing.T) {
	t.Parallel()

	log, buf := bufLogger()

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("body"))
	})

	middleware.AccessLog(log)(next).
		ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", http.NoBody))

	entry := decodeLog(t, buf)
	assert.InDelta(t, http.StatusOK, entry["status"], 0)
	assert.InDelta(t, 4, entry["bytes"], 0)
}

func TestAccessLogCapturesRequestID(t *testing.T) {
	t.Parallel()

	log, buf := bufLogger()

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	req.Header.Set("X-Request-Id", "trace-42")

	handler := chimw.RequestID(middleware.AccessLog(log)(next))
	handler.ServeHTTP(httptest.NewRecorder(), req)

	assert.Equal(t, "trace-42", decodeLog(t, buf)["request_id"])
}

func TestAccessLogWithoutRequestIDLogsEmptyString(t *testing.T) {
	t.Parallel()

	log, buf := bufLogger()

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })

	middleware.AccessLog(log)(next).
		ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", http.NoBody))

	assert.Empty(t, decodeLog(t, buf)["request_id"])
}

func TestRecovererTurnsPanicIntoInternalError(t *testing.T) {
	t.Parallel()

	log, buf := bufLogger()

	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("kaboom") })

	rec := httptest.NewRecorder()
	middleware.Recoverer(log)(next).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/boom", http.NoBody))

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"error":{"code":"internal_error","message":"internal server error"}}`, rec.Body.String())

	entry := decodeLog(t, buf)
	assert.Equal(t, "panic recovered", entry["msg"])
	assert.Equal(t, "kaboom", entry["panic"])
	assert.Contains(t, entry["stack"], "runtime/debug.Stack")
}

func TestRecovererDoesNotLeakPanicTextToClient(t *testing.T) {
	t.Parallel()

	log, _ := bufLogger()

	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("secret db password") })

	rec := httptest.NewRecorder()
	middleware.Recoverer(log)(next).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", http.NoBody))

	assert.NotContains(t, rec.Body.String(), "secret db password")
}

func TestRecovererRepanicsOnErrAbortHandler(t *testing.T) {
	t.Parallel()

	log, _ := bufLogger()

	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic(http.ErrAbortHandler) })

	handler := middleware.Recoverer(log)(next)

	assert.PanicsWithError(t, http.ErrAbortHandler.Error(), func() {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", http.NoBody))
	})
}

func TestRecovererPassesThroughWhenNoPanic(t *testing.T) {
	t.Parallel()

	log, buf := bufLogger()

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	rec := httptest.NewRecorder()
	middleware.Recoverer(log)(next).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", http.NoBody))

	assert.Equal(t, http.StatusTeapot, rec.Code)
	assert.Empty(t, buf.String())
}

func TestRecovererHandlesNonStringPanicValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
	}{
		{name: "int", value: 42},
		{name: "struct", value: struct{ Code string }{Code: "x"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			log, buf := bufLogger()

			next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic(tc.value) })

			rec := httptest.NewRecorder()
			middleware.Recoverer(log)(next).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", http.NoBody))

			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Contains(t, buf.String(), "panic recovered")
		})
	}
}
