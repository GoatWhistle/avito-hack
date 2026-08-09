package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func uploadsRouter(t *testing.T, guard uploadGuard) (chi.Router, uuid.UUID) {
	t.Helper()

	dir := t.TempDir()
	itemID := uuid.New()

	require.NoError(t, os.MkdirAll(filepath.Join(dir, itemID.String()), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, itemID.String(), "photo.jpg"), []byte("binary"), 0o600))

	r := chi.NewRouter()
	mountUploads(r, "/uploads", dir, guard)

	return r, itemID
}

func TestMountUploadsServesVisibleItem(t *testing.T) {
	t.Parallel()

	r, itemID := uploadsRouter(t, func(context.Context, uuid.UUID) (bool, error) { return true, nil })

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/uploads/"+itemID.String()+"/photo.jpg", http.NoBody))

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "binary", rec.Body.String())
	assert.Equal(t, "private, max-age=86400, must-revalidate", rec.Header().Get("Cache-Control"))
	assert.Equal(t, "nosniff", rec.Header().Get("X-Content-Type-Options"))
}

func TestMountUploadsHidesInvisibleItem(t *testing.T) {
	t.Parallel()

	var asked uuid.UUID

	r, itemID := uploadsRouter(t, func(_ context.Context, id uuid.UUID) (bool, error) {
		asked = id

		return false, nil
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/uploads/"+itemID.String()+"/photo.jpg", http.NoBody))

	require.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, itemID, asked)
	assert.NotContains(t, rec.Body.String(), "binary")
}

func TestMountUploadsGuardFailureIsNotLeaking(t *testing.T) {
	t.Parallel()

	r, itemID := uploadsRouter(t, func(context.Context, uuid.UUID) (bool, error) {
		return false, errors.New("database is down")
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/uploads/"+itemID.String()+"/photo.jpg", http.NoBody))

	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.NotContains(t, rec.Body.String(), "database is down")
}

func TestMountUploadsRejectsMalformedPaths(t *testing.T) {
	t.Parallel()

	called := 0
	r, _ := uploadsRouter(t, func(context.Context, uuid.UUID) (bool, error) {
		called++

		return true, nil
	})

	paths := []string{
		"/uploads/not-a-uuid/photo.jpg",
		"/uploads/photo.jpg",
		"/uploads/" + uuid.New().String() + "/nested/photo.jpg",
	}

	for _, p := range paths {
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, p, http.NoBody))

		assert.Equal(t, http.StatusNotFound, rec.Code, p)
	}

	assert.Equal(t, 0, called)
}

func TestUploadGuardWithoutPoolAllows(t *testing.T) {
	t.Parallel()

	visible, err := newUploadGuard(nil)(context.Background(), uuid.New())

	require.NoError(t, err)
	assert.True(t, visible)
}
