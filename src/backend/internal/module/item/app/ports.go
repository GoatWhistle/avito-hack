package app

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(context.Context) error) error
}

type Clock interface {
	Now() time.Time
}

type OwnerView struct {
	ID          uuid.UUID
	DisplayID   string
	DisplayName string
}

type OwnerProvider interface {
	ByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]OwnerView, error)
}

type StoredFile struct {
	Name        string
	ContentType string
	Size        int64
}

type PhotoStorage interface {
	Save(ctx context.Context, itemID uuid.UUID, content io.Reader, contentType string) (StoredFile, error)
	URL(itemDisplayID, name string) string
	Delete(ctx context.Context, itemID uuid.UUID, name string) error
}

type PhotoBytesLoader interface {
	DataURL(ctx context.Context, itemID uuid.UUID, publicURL string) (string, error)
}
