package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func ContextWithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return withTx(ctx, tx)
}
