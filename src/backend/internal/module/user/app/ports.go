package app

import (
	"context"
	"time"

	"github.com/avito-hack/backend/internal/shared/auth"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(context.Context) error) error
}

type Clock interface {
	Now() time.Time
}

type TokenIssuer interface {
	Issue(actor auth.Actor) (string, time.Time, error)
}
