package infra

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const maxModerationPhotoBytes = 8 * 1024 * 1024

var moderationContentTypeByExt = map[string]string{
	".jpg":  contentTypeJPEG,
	".jpeg": contentTypeJPEG,
	".png":  contentTypePNG,
	".webp": contentTypeWebP,
}

type LocalPhotoBytesLoader struct {
	root      string
	publicURL string
}

func NewLocalPhotoBytesLoader(root, publicURL string) *LocalPhotoBytesLoader {
	return &LocalPhotoBytesLoader{root: root, publicURL: publicURL}
}

func (l *LocalPhotoBytesLoader) DataURL(_ context.Context, itemID uuid.UUID, publicURL string) (string, error) {
	name := filepath.Base(filepath.Clean(strings.TrimPrefix(publicURL, l.publicURL)))
	if name == "." || name == ".." || name == string(filepath.Separator) {
		return "", fmt.Errorf("invalid photo url %q", publicURL)
	}

	path := filepath.Join(l.root, itemID.String(), name)

	//nolint:gosec
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read photo file: %w", err)
	}

	if len(data) > maxModerationPhotoBytes {
		return "", fmt.Errorf("photo file too large for moderation: %d bytes", len(data))
	}

	contentType, ok := moderationContentTypeByExt[strings.ToLower(filepath.Ext(path))]
	if !ok {
		return "", fmt.Errorf("unsupported photo extension %q", filepath.Ext(path))
	}

	encoded := base64.StdEncoding.EncodeToString(data)

	return "data:" + contentType + ";base64," + encoded, nil
}
