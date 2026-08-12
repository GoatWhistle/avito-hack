package infra_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/infra"
)

func TestLocalPhotoStorageSaveWritesFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		contentType string
		wantExt     string
	}{
		{name: "jpeg", contentType: "image/jpeg", wantExt: ".jpg"},
		{name: "png", contentType: "image/png", wantExt: ".png"},
		{name: "webp", contentType: "image/webp", wantExt: ".webp"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			root := t.TempDir()
			itemID := uuid.New()
			storage := infra.NewLocalPhotoStorage(root, "/media")

			stored, err := storage.Save(context.Background(), itemID, strings.NewReader("payload"), tc.contentType)

			require.NoError(t, err)
			assert.Equal(t, tc.contentType, stored.ContentType)
			assert.Equal(t, int64(len("payload")), stored.Size)
			assert.True(t, strings.HasSuffix(stored.Name, tc.wantExt))

			written, err := os.ReadFile(filepath.Join(root, itemID.String(), stored.Name))

			require.NoError(t, err)
			assert.Equal(t, "payload", string(written))
		})
	}
}

func TestLocalPhotoStorageSaveRejectsUnsupportedContentType(t *testing.T) {
	t.Parallel()

	storage := infra.NewLocalPhotoStorage(t.TempDir(), "/media")

	_, err := storage.Save(context.Background(), uuid.New(), strings.NewReader("x"), "application/pdf")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported content type")
}

func TestLocalPhotoStorageSaveFailsWhenDirNotCreatable(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	blocker := filepath.Join(root, "blocked")
	require.NoError(t, os.WriteFile(blocker, []byte("not a dir"), 0o600))

	storage := infra.NewLocalPhotoStorage(blocker, "/media")

	_, err := storage.Save(context.Background(), uuid.New(), strings.NewReader("x"), "image/png")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "create photo dir")
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) {
	return 0, errRead
}

var errRead = errors.New("read failed")

func TestLocalPhotoStorageSaveCleansUpOnCopyError(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	itemID := uuid.New()
	storage := infra.NewLocalPhotoStorage(root, "/media")

	_, err := storage.Save(context.Background(), itemID, failingReader{}, "image/jpeg")

	require.ErrorIs(t, err, errRead)
	assert.Contains(t, err.Error(), "write photo file")

	entries, readErr := os.ReadDir(filepath.Join(root, itemID.String()))

	require.NoError(t, readErr)
	assert.Empty(t, entries)
}

func TestLocalPhotoStorageDeleteRejectsSeparatorName(t *testing.T) {
	t.Parallel()

	storage := infra.NewLocalPhotoStorage(t.TempDir(), "/media")

	err := storage.Delete(context.Background(), uuid.New(), string(filepath.Separator))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid photo name")
}

func TestLocalPhotoStorageDeleteWrapsUnexpectedError(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	itemID := uuid.New()

	dir := filepath.Join(root, itemID.String(), "photo.jpg")
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "child"), 0o750))

	storage := infra.NewLocalPhotoStorage(root, "/media")

	err := storage.Delete(context.Background(), itemID, "photo.jpg")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "remove photo file")
}

func TestLocalPhotoStorageURL(t *testing.T) {
	t.Parallel()

	displayID := "abc123def456"
	storage := infra.NewLocalPhotoStorage(t.TempDir(), "/media/photos")

	assert.Equal(t, "/media/photos/"+displayID+"/a.jpg", storage.URL(displayID, "a.jpg"))
}

func TestLocalPhotoStorageDeleteRemovesExistingFile(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	itemID := uuid.New()
	dir := filepath.Join(root, itemID.String())
	require.NoError(t, os.MkdirAll(dir, 0o750))

	target := filepath.Join(dir, "photo.jpg")
	require.NoError(t, os.WriteFile(target, []byte("data"), 0o600))

	storage := infra.NewLocalPhotoStorage(root, "/media")

	require.NoError(t, storage.Delete(context.Background(), itemID, "photo.jpg"))

	_, err := os.Stat(target)

	assert.True(t, os.IsNotExist(err))
}

func TestLocalPhotoStorageDeleteIgnoresMissingFile(t *testing.T) {
	t.Parallel()

	storage := infra.NewLocalPhotoStorage(t.TempDir(), "/media")

	require.NoError(t, storage.Delete(context.Background(), uuid.New(), "absent.jpg"))
}

func TestLocalPhotoStorageDeleteStripsTraversal(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	itemID := uuid.New()
	dir := filepath.Join(root, itemID.String())
	require.NoError(t, os.MkdirAll(dir, 0o750))

	outside := filepath.Join(root, "secret.txt")
	require.NoError(t, os.WriteFile(outside, []byte("secret"), 0o600))

	inside := filepath.Join(dir, "secret.txt")
	require.NoError(t, os.WriteFile(inside, []byte("copy"), 0o600))

	storage := infra.NewLocalPhotoStorage(root, "/media")

	require.NoError(t, storage.Delete(context.Background(), itemID, "../secret.txt"))

	_, err := os.Stat(outside)

	require.NoError(t, err)

	_, err = os.Stat(inside)

	assert.True(t, os.IsNotExist(err))
}

func TestLocalPhotoStorageDeleteRejectsEmptyName(t *testing.T) {
	t.Parallel()

	storage := infra.NewLocalPhotoStorage(t.TempDir(), "/media")

	err := storage.Delete(context.Background(), uuid.New(), "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid photo name")
}
