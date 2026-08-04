package app

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/pet/domain"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(context.Context) error) error
}

type Clock interface {
	Now() time.Time
}

type Cache interface {
	Get(ctx context.Context, userID uuid.UUID) (*domain.Pet, error)
	Set(ctx context.Context, pet *domain.Pet) error
}

type Publisher interface {
	
}
