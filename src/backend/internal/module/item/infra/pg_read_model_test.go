package infra_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/module/item/infra"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
)

func TestReadModelListReturnsItems(t *testing.T) {
	t.Parallel()

	first, second := uuid.New(), uuid.New()
	firstDisplayID, secondDisplayID := domain.NewDisplayID(), domain.NewDisplayID()
	ownerID := uuid.New()

	rows := &pgtest.Rows{Records: [][]any{
		{first, firstDisplayID, ownerID, "Bike", int64(1000), "published", fixedTime, false, true, "bikes", "used"},
		{second, secondDisplayID, ownerID, "Chair", int64(2000), "bogus", fixedTime, true, false, "", ""},
	}}
	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rows}}

	items, err := infra.NewPgReadModel(nil).List(ctxWith(tx), app.ListFilter{Limit: 10})

	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, first, items[0].ID)
	assert.Equal(t, firstDisplayID, items[0].DisplayID)
	assert.Equal(t, ownerID, items[0].OwnerID)
	assert.Equal(t, "Bike", items[0].Title)
	assert.Equal(t, domain.StatusPublished, items[0].Status)
	assert.Equal(t, domain.StatusDraft, items[1].Status)
	assert.Equal(t, int64(2000), items[1].PriceKopeks)
	assert.False(t, items[0].IsSeed)
	assert.True(t, items[0].AIVerified)
	assert.True(t, items[1].IsSeed)
	assert.False(t, items[1].AIVerified)
	assert.True(t, rows.Closed)
}

func TestReadModelListPassesFilterArgs(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	tx := &pgtest.Tx{QueryRows: []pgx.Rows{&pgtest.Rows{}}}

	_, err := infra.NewPgReadModel(nil).List(ctxWith(tx), app.ListFilter{
		Status:  domain.StatusPublished,
		OwnerID: ownerID,
		Limit:   15,
	})

	require.NoError(t, err)
	require.Len(t, tx.QueryCalls, 1)
	assert.Contains(t, tx.QueryCalls[0].SQL, "FROM items")
	assert.Equal(t, []any{"published", "sold", "published", ownerID, 15}, tx.QueryCalls[0].Args)
}

func TestReadModelListWrapsQueryError(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{QueryErrs: []error{errDB}}

	_, err := infra.NewPgReadModel(nil).List(ctxWith(tx), app.ListFilter{Limit: 10})

	require.ErrorIs(t, err, errDB)
	assert.Contains(t, err.Error(), "query rows")
}

func TestReadModelListPropagatesScanError(t *testing.T) {
	t.Parallel()

	rows := &pgtest.Rows{Records: [][]any{{nil}}, ScanErr: errDB}
	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rows}}

	_, err := infra.NewPgReadModel(nil).List(ctxWith(tx), app.ListFilter{Limit: 10})

	require.ErrorIs(t, err, errDB)
	assert.Contains(t, err.Error(), "scan item")
}
