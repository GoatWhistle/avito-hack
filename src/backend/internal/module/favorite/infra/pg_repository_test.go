package infra_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/favorite/domain"
	"github.com/avito-hack/backend/internal/module/favorite/infra"
	"github.com/avito-hack/backend/internal/shared/postgres"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
)

var errDB = errors.New("db down")

var fixedTime = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func ctxWith(tx *pgtest.Tx) context.Context {
	return postgres.ContextWithTx(context.Background(), tx)
}

func newFavorite(t *testing.T) *domain.Favorite {
	t.Helper()

	fav, err := domain.New(uuid.New(), uuid.New(), fixedTime)
	require.NoError(t, err)

	return fav
}

func TestPgRepositoryAddReportsInsertedWhenRowAffected(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{ExecTags: []pgconn.CommandTag{pgtest.Affected(1)}}
	repo := infra.NewPgRepository(nil)

	added, err := repo.Add(ctxWith(tx), newFavorite(t))

	require.NoError(t, err)
	assert.True(t, added)
	require.Len(t, tx.ExecCalls, 1)
	assert.Contains(t, tx.ExecCalls[0].SQL, "INSERT INTO favorites")
}

func TestPgRepositoryAddReportsNotInsertedOnConflict(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{ExecTags: []pgconn.CommandTag{pgtest.Affected(0)}}
	repo := infra.NewPgRepository(nil)

	added, err := repo.Add(ctxWith(tx), newFavorite(t))

	require.NoError(t, err)
	assert.False(t, added)
}

func TestPgRepositoryAddWrapsError(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{ExecErrs: []error{errDB}}
	repo := infra.NewPgRepository(nil)

	_, err := repo.Add(ctxWith(tx), newFavorite(t))

	require.ErrorIs(t, err, errDB)
}

func TestPgRepositoryRemoveDeletesRow(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{ExecTags: []pgconn.CommandTag{pgtest.Affected(1)}}
	repo := infra.NewPgRepository(nil)

	err := repo.Remove(ctxWith(tx), uuid.New(), uuid.New())

	require.NoError(t, err)
	assert.Contains(t, tx.ExecCalls[0].SQL, "DELETE FROM favorites")
}

func TestPgRepositoryRemoveReturnsNotFoundWhenNoRows(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{ExecTags: []pgconn.CommandTag{pgtest.Affected(0)}}
	repo := infra.NewPgRepository(nil)

	err := repo.Remove(ctxWith(tx), uuid.New(), uuid.New())

	require.ErrorIs(t, err, domain.ErrFavoriteNotFound)
}

func TestPgRepositoryRemoveWrapsError(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{ExecErrs: []error{errDB}}
	repo := infra.NewPgRepository(nil)

	err := repo.Remove(ctxWith(tx), uuid.New(), uuid.New())

	require.ErrorIs(t, err, errDB)
}

func TestPgItemCheckerExists(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		row  pgtest.Row
		want bool
		err  error
	}{
		{name: "exists", row: pgtest.Row{Values: []any{true}}, want: true},
		{name: "missing", row: pgtest.Row{Values: []any{false}}, want: false},
		{name: "error", row: pgtest.Row{Err: errDB}, err: errDB},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tx := &pgtest.Tx{RowResults: []pgtest.Row{tc.row}}
			checker := infra.NewPgItemChecker(nil)

			got, err := checker.Exists(ctxWith(tx), uuid.New())

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
