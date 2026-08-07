package infra_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"

	"github.com/avito-hack/backend/internal/shared/postgres"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
)

var errDB = errors.New("db down")

var (
	fixedTime = time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)
	dayStart  = time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	dayEnd    = time.Date(2024, 3, 16, 0, 0, 0, 0, time.UTC)
)

var (
	userA = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	userB = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	itemA = uuid.MustParse("33333333-3333-3333-3333-333333333333")
)

func ctxWith(tx pgx.Tx) context.Context {
	return postgres.ContextWithTx(context.Background(), tx)
}

func rowsOf(records ...[]any) *pgtest.Rows {
	return &pgtest.Rows{Records: records}
}

func requireSQL(t *testing.T, call pgtest.Call, fragments ...string) {
	t.Helper()

	for _, fragment := range fragments {
		assert.Contains(t, call.SQL, fragment)
	}
}
