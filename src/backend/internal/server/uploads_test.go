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

	"github.com/avito-hack/backend/internal/module/item/domain"
)

func uploadsRouter(t *testing.T, itemID uuid.UUID, guard uploadGuard) chi.Router {
	t.Helper()

	dir := t.TempDir()

	require.NoError(t, os.MkdirAll(filepath.Join(dir, itemID.String()), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, itemID.String(), "photo.jpg"), []byte("binary"), 0o600))

	r := chi.NewRouter()
	mountUploads(r, "/uploads", dir, guard)

	return r
}

func TestMountUploadsServesVisibleItem(t *testing.T) {
	t.Parallel()

	itemID := uuid.New()
	displayID := domain.NewDisplayID()

	r := uploadsRouter(t, itemID, func(context.Context, string) (uuid.UUID, bool, error) {
		return itemID, true, nil
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/uploads/"+displayID+"/photo.jpg", http.NoBody))

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "binary", rec.Body.String())
	assert.Equal(t, "private, max-age=86400, must-revalidate", rec.Header().Get("Cache-Control"))
	assert.Equal(t, "nosniff", rec.Header().Get("X-Content-Type-Options"))
}

func TestMountUploadsHidesInvisibleItem(t *testing.T) {
	t.Parallel()

	var asked string

	displayID := domain.NewDisplayID()

	r := uploadsRouter(t, uuid.New(), func(_ context.Context, id string) (uuid.UUID, bool, error) {
		asked = id

		return uuid.Nil, false, nil
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/uploads/"+displayID+"/photo.jpg", http.NoBody))

	require.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, displayID, asked)
	assert.NotContains(t, rec.Body.String(), "binary")
}

func TestMountUploadsGuardFailureIsNotLeaking(t *testing.T) {
	t.Parallel()

	displayID := domain.NewDisplayID()

	r := uploadsRouter(t, uuid.New(), func(context.Context, string) (uuid.UUID, bool, error) {
		return uuid.Nil, false, errors.New("database is down")
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/uploads/"+displayID+"/photo.jpg", http.NoBody))

	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.NotContains(t, rec.Body.String(), "database is down")
}

func TestMountUploadsDoesNotExposeItemUUID(t *testing.T) {
	t.Parallel()

	itemID := uuid.New()

	r := uploadsRouter(t, itemID, func(context.Context, string) (uuid.UUID, bool, error) {
		return itemID, true, nil
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/uploads/"+itemID.String()+"/photo.jpg", http.NoBody))

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.NotContains(t, rec.Body.String(), "binary")
}

func TestMountUploadsRejectsMalformedPaths(t *testing.T) {
	t.Parallel()

	called := 0
	displayID := domain.NewDisplayID()

	r := uploadsRouter(t, uuid.New(), func(context.Context, string) (uuid.UUID, bool, error) {
		called++

		return uuid.New(), true, nil
	})

	paths := []string{
		"/uploads/not-a-display-id/photo.jpg",
		"/uploads/photo.jpg",
		"/uploads/" + displayID + "/nested/photo.jpg",
	}

	for _, p := range paths {
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, p, http.NoBody))

		assert.Equal(t, http.StatusNotFound, rec.Code, p)
	}

	assert.Equal(t, 0, called)
}

func TestUploadGuardWithoutPoolDenies(t *testing.T) {
	t.Parallel()

	itemID, visible, err := newUploadGuard(nil)(context.Background(), domain.NewDisplayID())

	require.NoError(t, err)
	assert.False(t, visible)
	assert.Equal(t, uuid.Nil, itemID)
}
