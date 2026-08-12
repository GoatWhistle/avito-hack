package infra_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/infra"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
)

func TestOwnerProviderByIDsShortCircuitsOnEmptyInput(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{}

	owners, err := infra.NewPgOwnerProvider(nil).ByIDs(ctxWith(tx), nil)

	require.NoError(t, err)
	assert.Empty(t, owners)
	assert.Empty(t, tx.QueryCalls)
}

func TestOwnerProviderByIDsReturnsMap(t *testing.T) {
	t.Parallel()

	first, second := uuid.New(), uuid.New()
	ids := []uuid.UUID{first, second}

	rows := &pgtest.Rows{Records: [][]any{
		{first, "alice1234567", "Alice"},
		{second, "bob123456789", "Bob"},
	}}
	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rows}}

	owners, err := infra.NewPgOwnerProvider(nil).ByIDs(ctxWith(tx), ids)

	require.NoError(t, err)
	require.Len(t, owners, 2)
	assert.Equal(t, "Alice", owners[first].DisplayName)
	assert.Equal(t, "alice1234567", owners[first].DisplayID)
	assert.Equal(t, second, owners[second].ID)
	require.Len(t, tx.QueryCalls, 1)
	assert.Contains(t, tx.QueryCalls[0].SQL, "SELECT id, display_id, full_name FROM users WHERE id = ANY($1)")
	assert.Equal(t, []any{ids}, tx.QueryCalls[0].Args)
}

func TestOwnerProviderByIDsExcludesSoftDeletedUsers(t *testing.T) {
	t.Parallel()

	rows := &pgtest.Rows{Records: [][]any{}}
	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rows}}

	_, err := infra.NewPgOwnerProvider(nil).ByIDs(ctxWith(tx), []uuid.UUID{uuid.New()})

	require.NoError(t, err)
	require.Len(t, tx.QueryCalls, 1)
	assert.Contains(t, tx.QueryCalls[0].SQL, "deleted_at IS NULL")
}

func TestOwnerProviderByIDsWrapsQueryError(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{QueryErrs: []error{errDB}}

	_, err := infra.NewPgOwnerProvider(nil).ByIDs(ctxWith(tx), []uuid.UUID{uuid.New()})

	require.ErrorIs(t, err, errDB)
	assert.Contains(t, err.Error(), "query rows")
}

func TestOwnerProviderByIDsPropagatesScanError(t *testing.T) {
	t.Parallel()

	rows := &pgtest.Rows{Records: [][]any{{nil}}, ScanErr: errDB}
	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rows}}

	_, err := infra.NewPgOwnerProvider(nil).ByIDs(ctxWith(tx), []uuid.UUID{uuid.New()})

	require.ErrorIs(t, err, errDB)
	assert.Contains(t, err.Error(), "scan owner")
}
