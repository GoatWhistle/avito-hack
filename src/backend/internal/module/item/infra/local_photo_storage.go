package infra

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/item/app"
)

const (
	dirPerm  = 0o755
	filePerm = 0o644
)

var extensions = map[string]string{
	contentTypeJPEG: ".jpg",
	contentTypePNG:  ".png",
	contentTypeWebP: ".webp",
}

type LocalPhotoStorage struct {
	root      string
	publicURL string
}

func NewLocalPhotoStorage(root, publicURL string) *LocalPhotoStorage {
	return &LocalPhotoStorage{root: root, publicURL: publicURL}
}

func (s *LocalPhotoStorage) Save(
	_ context.Context,
	itemID uuid.UUID,
	content io.Reader,
	contentType string,
) (app.StoredFile, error) {
	ext, ok := extensions[contentType]
	if !ok {
		return app.StoredFile{}, fmt.Errorf("unsupported content type %q", contentType)
	}

	dir := filepath.Join(s.root, itemID.String())
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return app.StoredFile{}, fmt.Errorf("create photo dir: %w", err)
	}

	name := uuid.New().String() + ext

	target := filepath.Join(dir, name)

	//nolint:gosec // the file name is a freshly generated uuid, not user input
	file, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, filePerm)
	if err != nil {
		return app.StoredFile{}, fmt.Errorf("create photo file: %w", err)
	}
	defer func() { _ = file.Close() }() //nolint:errcheck // the file is fsync-free and errors surface on Copy

	written, err := io.Copy(file, content)
	if err != nil {
		_ = os.Remove(target) //nolint:errcheck // best-effort cleanup of a partial upload

		return app.StoredFile{}, fmt.Errorf("write photo file: %w", err)
	}

	return app.StoredFile{Name: name, ContentType: contentType, Size: written}, nil
}

func (s *LocalPhotoStorage) URL(itemDisplayID, name string) string {
	return path.Join(s.publicURL, itemDisplayID, name)
}

func (s *LocalPhotoStorage) Delete(_ context.Context, itemID uuid.UUID, name string) error {
	base := filepath.Base(name)
	if base == "." || base == string(filepath.Separator) {
		return fmt.Errorf("invalid photo name %q", name)
	}

	target := filepath.Join(s.root, itemID.String(), base)
	if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove photo file: %w", err)
	}

	return nil
}
