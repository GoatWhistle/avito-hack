package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/pagination"
)

type stubReadModel struct {
	rows       []app.ListItem
	lastFilter app.ListFilter
	err        error
}

func (s *stubReadModel) List(_ context.Context, f app.ListFilter) ([]app.ListItem, error) {
	s.lastFilter = f
	if s.err != nil {
		return nil, s.err
	}

	return s.rows, nil
}

type stubOwners struct {
	names   map[uuid.UUID]string
	lastIDs []uuid.UUID
	err     error
}

func (s *stubOwners) ByIDs(_ context.Context, ids []uuid.UUID) (map[uuid.UUID]app.OwnerView, error) {
	s.lastIDs = ids
	if s.err != nil {
		return nil, s.err
	}

	result := make(map[uuid.UUID]app.OwnerView, len(ids))
	for _, id := range ids {
		if name, ok := s.names[id]; ok {
			result[id] = app.OwnerView{ID: id, DisplayName: name}
		}
	}

	return result, nil
}

func makeRows(n int, ownerID uuid.UUID) []app.ListItem {
	rows := make([]app.ListItem, 0, n)
	for i := range n {
		rows = append(rows, app.ListItem{
			ID: uuid.New(), OwnerID: ownerID, Title: "item",
			Status: domain.StatusPublished, CreatedAt: fixedTime.Add(-time.Duration(i) * time.Minute),
		})
	}

	return rows
}

func TestListItemsRequestsOneExtraRow(t *testing.T) {
	t.Parallel()

	read := &stubReadModel{rows: makeRows(3, uuid.New())}
	handler := app.NewListItemsHandler(read, nil)

	_, err := handler.Handle(t.Context(), app.ListItemsQuery{Limit: 5})

	require.NoError(t, err)
	assert.Equal(t, 6, read.lastFilter.Limit)
}

func TestListItemsPagination(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		rows       int
		limit      int
		wantItems  int
		wantCursor bool
	}{
		{name: "fewer rows than limit", rows: 3, limit: 10, wantItems: 3},
		{name: "exactly limit rows", rows: 10, limit: 10, wantItems: 10},
		{name: "one extra row signals next page", rows: 11, limit: 10, wantItems: 10, wantCursor: true},
		{name: "zero limit falls back to default", rows: 2, limit: 0, wantItems: 2},
		{name: "limit clamped to max", rows: 2, limit: 1000, wantItems: 2},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			read := &stubReadModel{rows: makeRows(tc.rows, uuid.New())}
			handler := app.NewListItemsHandler(read, nil)

			result, err := handler.Handle(t.Context(), app.ListItemsQuery{Limit: tc.limit})

			require.NoError(t, err)
			assert.Len(t, result.Items, tc.wantItems)

			if !tc.wantCursor {
				assert.Empty(t, result.NextCursor)

				return
			}

			decoded, decodeErr := pagination.DecodeCursor(result.NextCursor)
			require.NoError(t, decodeErr)
			assert.Equal(t, result.Items[len(result.Items)-1].ID, decoded.ID)
		})
	}
}

func TestListItemsPassesFiltersThrough(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	cursor := pagination.Cursor{CreatedAt: fixedTime, ID: uuid.New()}
	read := &stubReadModel{}
	handler := app.NewListItemsHandler(read, nil)

	_, err := handler.Handle(t.Context(), app.ListItemsQuery{
		Status: domain.StatusPublished, OwnerID: ownerID, Search: "bike", Cursor: cursor, Limit: 7,
	})

	require.NoError(t, err)
	assert.Equal(t, domain.StatusPublished, read.lastFilter.Status)
	assert.Equal(t, ownerID, read.lastFilter.OwnerID)
	assert.Equal(t, "bike", read.lastFilter.Search)
	assert.Equal(t, cursor.ID, read.lastFilter.Cursor.ID)
}

func TestListItemsRejectsCursorFromAnotherSort(t *testing.T) {
	t.Parallel()

	cursor := pagination.Cursor{CreatedAt: fixedTime, ID: uuid.New(), Sort: "newest"}
	read := &stubReadModel{}
	handler := app.NewListItemsHandler(read, nil)

	_, err := handler.Handle(t.Context(), app.ListItemsQuery{
		Sort: "price_asc", Cursor: cursor, Limit: 7,
	})

	var invalid *domainerr.InvalidError
	require.ErrorAs(t, err, &invalid)
	assert.Equal(t, "cursor", invalid.Field)
}

func TestListItemsAcceptsCursorMatchingSort(t *testing.T) {
	t.Parallel()

	cursor := pagination.Cursor{
		CreatedAt: fixedTime, ID: uuid.New(), PriceKopeks: 500, Sort: "price_asc",
	}
	read := &stubReadModel{}
	handler := app.NewListItemsHandler(read, nil)

	_, err := handler.Handle(t.Context(), app.ListItemsQuery{
		Sort: "price_asc", Cursor: cursor, Limit: 7,
	})

	require.NoError(t, err)
	assert.Equal(t, app.ListSortPriceAsc, read.lastFilter.Sort)
	assert.Equal(t, int64(500), read.lastFilter.Cursor.PriceKopeks)
}

func TestListItemsAcceptsLegacyCursorOnDefaultSort(t *testing.T) {
	t.Parallel()

	cursor := pagination.Cursor{CreatedAt: fixedTime, ID: uuid.New()}
	read := &stubReadModel{}
	handler := app.NewListItemsHandler(read, nil)

	_, err := handler.Handle(t.Context(), app.ListItemsQuery{Cursor: cursor, Limit: 7})

	require.NoError(t, err)
	assert.Equal(t, app.ListSortNewest, read.lastFilter.Sort)
}

func TestListItemsFillsOwnerNamesWithoutDuplicates(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	otherID := uuid.New()
	rows := makeRows(3, ownerID)
	rows[2].OwnerID = otherID

	owners := &stubOwners{names: map[uuid.UUID]string{ownerID: "Alice"}}
	handler := app.NewListItemsHandler(&stubReadModel{rows: rows}, owners)

	result, err := handler.Handle(t.Context(), app.ListItemsQuery{Limit: 10})

	require.NoError(t, err)
	assert.Len(t, owners.lastIDs, 2)
	assert.Equal(t, "Alice", result.Items[0].OwnerName)
	assert.Equal(t, "Alice", result.Items[1].OwnerName)
	assert.Empty(t, result.Items[2].OwnerName)
}

func TestListItemsSkipsOwnerLookupWhenEmpty(t *testing.T) {
	t.Parallel()

	owners := &stubOwners{}
	handler := app.NewListItemsHandler(&stubReadModel{}, owners)

	result, err := handler.Handle(t.Context(), app.ListItemsQuery{Limit: 10})

	require.NoError(t, err)
	assert.Empty(t, result.Items)
	assert.Nil(t, owners.lastIDs)
}

func TestListItemsPropagatesErrors(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("query failed")

	t.Run("read model failure", func(t *testing.T) {
		t.Parallel()

		handler := app.NewListItemsHandler(&stubReadModel{err: sentinel}, nil)

		_, err := handler.Handle(t.Context(), app.ListItemsQuery{Limit: 10})

		require.ErrorIs(t, err, sentinel)
	})

	t.Run("owner lookup failure", func(t *testing.T) {
		t.Parallel()

		read := &stubReadModel{rows: makeRows(2, uuid.New())}
		handler := app.NewListItemsHandler(read, &stubOwners{err: sentinel})

		_, err := handler.Handle(t.Context(), app.ListItemsQuery{Limit: 10})

		require.ErrorIs(t, err, sentinel)
	})
}

func TestListPhotosHandler(t *testing.T) {
	t.Parallel()

	itemID := uuid.New()
	photo, err := domain.NewPhoto(domain.NewPhotoParams{
		ItemID: itemID, URL: "/media/a.jpg", Position: 0, Now: fixedTime,
	})
	require.NoError(t, err)

	photos := &countingPhotos{added: []*domain.Photo{photo}}

	listed, err := app.NewListPhotosHandler(photos).Handle(t.Context(), itemID)
	require.NoError(t, err)
	require.Len(t, listed, 1)
	assert.Equal(t, "/media/a.jpg", listed[0].URL())

	sentinel := errors.New("query failed")
	_, err = app.NewListPhotosHandler(&countingPhotos{countErr: sentinel}).Handle(t.Context(), itemID)
	require.ErrorIs(t, err, sentinel)
}

func TestDeletePhotoHandler(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusDraft)
	photos := &countingPhotos{}

	handler := app.NewDeletePhotoHandler(&stubRepository{item: item}, photos, &stubStorage{}, passthroughTx{}, fakeClock{}, nil)

	require.NoError(t, handler.Handle(t.Context(), app.DeletePhotoCommand{
		ItemID: item.ID(), PhotoDisplayID: domain.NewDisplayID(), ActorID: ownerID,
	}))
	assert.Equal(t, 1, photos.deleted)

	err := handler.Handle(t.Context(), app.DeletePhotoCommand{
		ItemID: item.ID(), PhotoDisplayID: domain.NewDisplayID(), ActorID: uuid.New(),
	})
	require.ErrorIs(t, err, domainerr.ErrForbidden)
	assert.Equal(t, 1, photos.deleted)

	sentinel := errors.New("load failed")
	failing := app.NewDeletePhotoHandler(&failingRepository{err: sentinel}, photos, &stubStorage{}, passthroughTx{}, fakeClock{}, nil)
	require.ErrorIs(t, failing.Handle(t.Context(), app.DeletePhotoCommand{}), sentinel)
}
