package infra_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/module/item/infra"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
)

func TestPhotoRepositoryAddInsertsRow(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{}
	photo := domain.RestorePhoto(uuid.New(), uuid.New(), "http://cdn/a.jpg", 3, fixedTime)

	err := infra.NewPgPhotoRepository(nil).Add(ctxWith(tx), photo)

	require.NoError(t, err)
	require.Len(t, tx.ExecCalls, 1)
	assert.Contains(t, tx.ExecCalls[0].SQL, "INSERT INTO item_photos")
	assert.Equal(t,
		[]any{photo.ID(), photo.ItemID(), "http://cdn/a.jpg", 3, fixedTime},
		tx.ExecCalls[0].Args)
}

func TestPhotoRepositoryAddWrapsError(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{ExecErrs: []error{errDB}}
	photo := domain.RestorePhoto(uuid.New(), uuid.New(), "http://cdn/a.jpg", 0, fixedTime)

	err := infra.NewPgPhotoRepository(nil).Add(ctxWith(tx), photo)

	require.ErrorIs(t, err, errDB)
	assert.Contains(t, err.Error(), "insert photo")
}

func TestPhotoRepositoryByItemIDReturnsPhotos(t *testing.T) {
	t.Parallel()

	itemID := uuid.New()
	first, second := uuid.New(), uuid.New()

	rows := &pgtest.Rows{Records: [][]any{
		{first, itemID, "http://cdn/1.jpg", 0, fixedTime},
		{second, itemID, "http://cdn/2.jpg", 1, fixedTime},
	}}
	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rows}}

	photos, err := infra.NewPgPhotoRepository(nil).ByItemID(ctxWith(tx), itemID)

	require.NoError(t, err)
	require.Len(t, photos, 2)
	assert.Equal(t, first, photos[0].ID())
	assert.Equal(t, 1, photos[1].Position())
	assert.True(t, rows.Closed)
	require.Len(t, tx.QueryCalls, 1)
	assert.Contains(t, tx.QueryCalls[0].SQL, "FROM item_photos")
	assert.Contains(t, tx.QueryCalls[0].SQL, "ORDER BY position")
	assert.Equal(t, []any{itemID}, tx.QueryCalls[0].Args)
}

func TestPhotoRepositoryByItemIDReturnsEmpty(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{QueryRows: []pgx.Rows{&pgtest.Rows{}}}

	photos, err := infra.NewPgPhotoRepository(nil).ByItemID(ctxWith(tx), uuid.New())

	require.NoError(t, err)
	assert.Empty(t, photos)
}

func TestPhotoRepositoryByItemIDWrapsQueryError(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{QueryErrs: []error{errDB}}

	_, err := infra.NewPgPhotoRepository(nil).ByItemID(ctxWith(tx), uuid.New())

	require.ErrorIs(t, err, errDB)
	assert.Contains(t, err.Error(), "query rows")
}

func TestPhotoRepositoryByItemIDPropagatesScanError(t *testing.T) {
	t.Parallel()

	rows := &pgtest.Rows{Records: [][]any{{nil}}, ScanErr: errDB}
	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rows}}

	_, err := infra.NewPgPhotoRepository(nil).ByItemID(ctxWith(tx), uuid.New())

	require.ErrorIs(t, err, errDB)
	assert.Contains(t, err.Error(), "scan photo")
}

func TestPhotoRepositoryByItemIDPropagatesIterationError(t *testing.T) {
	t.Parallel()

	rows := &pgtest.Rows{IterErr: errDB}
	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rows}}

	_, err := infra.NewPgPhotoRepository(nil).ByItemID(ctxWith(tx), uuid.New())

	require.ErrorIs(t, err, errDB)
	assert.Contains(t, err.Error(), "iterate rows")
}

func TestPhotoRepositoryCountByItemID(t *testing.T) {
	t.Parallel()

	itemID := uuid.New()
	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Values: []any{4}}}}

	count, err := infra.NewPgPhotoRepository(nil).CountByItemID(ctxWith(tx), itemID)

	require.NoError(t, err)
	assert.Equal(t, 4, count)
	assert.Contains(t, tx.QueryRowCalls[0].SQL, "count(*) FROM item_photos")
	assert.Equal(t, []any{itemID}, tx.QueryRowCalls[0].Args)
}

func TestPhotoRepositoryCountByItemIDWrapsError(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Err: errDB}}}

	_, err := infra.NewPgPhotoRepository(nil).CountByItemID(ctxWith(tx), uuid.New())

	require.ErrorIs(t, err, errDB)
	assert.Contains(t, err.Error(), "count photos")
}

func TestPhotoRepositoryDeleteByIDReturnsURL(t *testing.T) {
	t.Parallel()

	itemID, photoID := uuid.New(), uuid.New()
	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Values: []any{"http://cdn/a.jpg"}}}}

	url, err := infra.NewPgPhotoRepository(nil).DeleteByID(ctxWith(tx), itemID, photoID)

	require.NoError(t, err)
	assert.Equal(t, "http://cdn/a.jpg", url)
	assert.Contains(t, tx.QueryRowCalls[0].SQL, "DELETE FROM item_photos")
	assert.Contains(t, tx.QueryRowCalls[0].SQL, "RETURNING url")
	assert.Equal(t, []any{itemID, photoID}, tx.QueryRowCalls[0].Args)
}

func TestPhotoRepositoryDeleteByIDReturnsNotFound(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Err: pgx.ErrNoRows}}}

	_, err := infra.NewPgPhotoRepository(nil).DeleteByID(ctxWith(tx), uuid.New(), uuid.New())

	require.ErrorIs(t, err, domain.ErrPhotoNotFound)
}

func TestPhotoRepositoryDeleteByIDWrapsError(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Err: errDB}}}

	_, err := infra.NewPgPhotoRepository(nil).DeleteByID(ctxWith(tx), uuid.New(), uuid.New())

	require.ErrorIs(t, err, errDB)
	assert.Contains(t, err.Error(), "delete photo")
}
