package app

import (
	"context"
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
	DisplayName string
}

type OwnerProvider interface {
	ByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]OwnerView, error)
}
