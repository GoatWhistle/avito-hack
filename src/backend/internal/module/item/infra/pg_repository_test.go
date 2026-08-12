package infra_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/module/item/infra"
	"github.com/avito-hack/backend/internal/shared/postgres"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
	"github.com/avito-hack/backend/internal/shared/vo"
)

var errDB = errors.New("db down")

var fixedTime = time.Date(2024, 5, 4, 9, 0, 0, 0, time.UTC)

func ctxWith(tx *pgtest.Tx) context.Context {
	return postgres.ContextWithTx(context.Background(), tx)
}

func newItem(t *testing.T) *domain.Item {
	t.Helper()

	return domain.RestoreItem(domain.RestoreItemParams{
		ID:          uuid.New(),
		DisplayID:   domain.NewDisplayID(),
		OwnerID:     uuid.New(),
		Title:       "Table",
		Description: "Oak table",
		Price:       vo.MustMoney(9900),
		Status:      domain.StatusPublished,
		Attributes:  domain.NewAttributes(map[string]string{"material": "oak"}),
		CreatedAt:   fixedTime,
		UpdatedAt:   fixedTime,
	})
}

func itemValues(t *testing.T, id, ownerID uuid.UUID) []any {
	t.Helper()

	return []any{
		id, domain.NewDisplayID(), ownerID, "Table", "Oak table", int64(9900),
		string(domain.StatusPublished), []byte(`{"material":"oak"}`), fixedTime, fixedTime, false, false,
	}
}

func TestItemRepositorySaveInsertsRow(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{}
	item := newItem(t)

	err := infra.NewPgRepository(nil).Save(ctxWith(tx), item)

	require.NoError(t, err)
	require.Len(t, tx.ExecCalls, 1)
	assert.Contains(t, tx.ExecCalls[0].SQL, "INSERT INTO items")
	assert.Contains(t, tx.ExecCalls[0].SQL, "ON CONFLICT (id) DO UPDATE")

	args := tx.ExecCalls[0].Args
	require.Len(t, args, 12)
	assert.Equal(t, item.ID(), args[0])
	assert.Equal(t, item.DisplayID(), args[1])
	assert.Equal(t, item.OwnerID(), args[2])
	assert.Equal(t, "Table", args[3])
	assert.Equal(t, "Oak table", args[4])
	assert.Equal(t, int64(9900), args[5])
	assert.Equal(t, "published", args[6])
	attributes, ok := args[7].([]byte)
	require.True(t, ok)
	assert.JSONEq(t, `{"material":"oak"}`, string(attributes))
	assert.Equal(t, fixedTime, args[8])
	assert.Equal(t, fixedTime, args[9])
	assert.Equal(t, false, args[10])
	assert.Equal(t, false, args[11])
}

func TestItemRepositorySaveWrapsError(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{ExecErrs: []error{errDB}}

	err := infra.NewPgRepository(nil).Save(ctxWith(tx), newItem(t))

	require.ErrorIs(t, err, errDB)
	assert.Contains(t, err.Error(), "save item")
}

func TestItemRepositoryByIDReturnsItem(t *testing.T) {
	t.Parallel()

	id, ownerID := uuid.New(), uuid.New()
	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Values: itemValues(t, id, ownerID)}}}

	item, err := infra.NewPgRepository(nil).ByID(ctxWith(tx), id)

	require.NoError(t, err)
	assert.Equal(t, id, item.ID())
	assert.Equal(t, ownerID, item.OwnerID())
	assert.Equal(t, domain.StatusPublished, item.Status())
	require.Len(t, tx.QueryRowCalls, 1)
	assert.Contains(t, tx.QueryRowCalls[0].SQL, "FROM items WHERE id = $1 AND deleted_at IS NULL")
	assert.NotContains(t, tx.QueryRowCalls[0].SQL, "FOR UPDATE")
	assert.Equal(t, []any{id}, tx.QueryRowCalls[0].Args)
}

func TestItemRepositoryByIDReturnsNotFound(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Err: pgx.ErrNoRows}}}

	_, err := infra.NewPgRepository(nil).ByID(ctxWith(tx), uuid.New())

	require.ErrorIs(t, err, domain.ErrItemNotFound)
}

func TestItemRepositoryByIDWrapsError(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Err: errDB}}}

	_, err := infra.NewPgRepository(nil).ByID(ctxWith(tx), uuid.New())

	require.ErrorIs(t, err, errDB)
	assert.Contains(t, err.Error(), "query item")
}

func TestItemRepositoryByDisplayIDReturnsItem(t *testing.T) {
	t.Parallel()

	id, ownerID := uuid.New(), uuid.New()
	displayID := domain.NewDisplayID()
	values := itemValues(t, id, ownerID)
	values[1] = displayID
	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Values: values}}}

	item, err := infra.NewPgRepository(nil).ByDisplayID(ctxWith(tx), displayID)

	require.NoError(t, err)
	assert.Equal(t, id, item.ID())
	assert.Equal(t, displayID, item.DisplayID())
	require.Len(t, tx.QueryRowCalls, 1)
	assert.Contains(t, tx.QueryRowCalls[0].SQL, "FROM items WHERE display_id = $1 AND deleted_at IS NULL")
	assert.Equal(t, []any{displayID}, tx.QueryRowCalls[0].Args)
}

func TestItemRepositoryByDisplayIDReturnsNotFound(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Err: pgx.ErrNoRows}}}

	_, err := infra.NewPgRepository(nil).ByDisplayID(ctxWith(tx), "unknown12345")

	require.ErrorIs(t, err, domain.ErrItemNotFound)
}

func TestItemRepositoryByIDForUpdateLocksRow(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Values: itemValues(t, id, uuid.New())}}}

	item, err := infra.NewPgRepository(nil).ByIDForUpdate(ctxWith(tx), id)

	require.NoError(t, err)
	assert.Equal(t, id, item.ID())
	assert.Contains(t, tx.QueryRowCalls[0].SQL, "FOR UPDATE")
}

func TestItemRepositoryByIDForUpdateReturnsNotFound(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Err: pgx.ErrNoRows}}}

	_, err := infra.NewPgRepository(nil).ByIDForUpdate(ctxWith(tx), uuid.New())

	require.ErrorIs(t, err, domain.ErrItemNotFound)
}

func TestItemRepositoryByIDRejectsCorruptRow(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	values := itemValues(t, id, uuid.New())
	values[5] = int64(-1)

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Values: values}}}

	_, err := infra.NewPgRepository(nil).ByID(ctxWith(tx), id)

	require.ErrorIs(t, err, vo.ErrNegativeMoney)
	assert.Contains(t, err.Error(), "query item")
}
