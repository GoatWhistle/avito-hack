package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/favorite/app"
	"github.com/avito-hack/backend/internal/shared/pagination"
)

type stubReadModel struct {
	rows       []app.FavoriteItem
	lastFilter app.ListFilter
	err        error
}

func (s *stubReadModel) List(_ context.Context, f app.ListFilter) ([]app.FavoriteItem, error) {
	s.lastFilter = f
	if s.err != nil {
		return nil, s.err
	}

	return s.rows, nil
}

func favoriteRows(n int) []app.FavoriteItem {
	rows := make([]app.FavoriteItem, 0, n)
	for i := range n {
		rows = append(rows, app.FavoriteItem{
			ItemID: uuid.New(), OwnerID: uuid.New(), Title: "item",
			Status: "published", CreatedAt: fixedTime.Add(-time.Duration(i) * time.Minute),
		})
	}

	return rows
}

func TestListFavoritesRequestsOneExtraRow(t *testing.T) {
	t.Parallel()

	read := &stubReadModel{}
	userID := uuid.New()

	_, err := app.NewListFavoritesHandler(read).Handle(t.Context(), app.ListFavoritesQuery{
		UserID: userID, Limit: 5,
	})

	require.NoError(t, err)
	assert.Equal(t, 6, read.lastFilter.Limit)
	assert.Equal(t, userID, read.lastFilter.UserID)
}

func TestListFavoritesPagination(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		rows       int
		limit      int
		wantItems  int
		wantCursor bool
	}{
		{name: "empty result", rows: 0, limit: 10, wantItems: 0},
		{name: "fewer than limit", rows: 4, limit: 10, wantItems: 4},
		{name: "exactly limit", rows: 10, limit: 10, wantItems: 10},
		{name: "one extra row", rows: 11, limit: 10, wantItems: 10, wantCursor: true},
		{name: "default limit applied", rows: 3, limit: 0, wantItems: 3},
		{name: "limit clamped", rows: 3, limit: 5000, wantItems: 3},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			read := &stubReadModel{rows: favoriteRows(tc.rows)}

			result, err := app.NewListFavoritesHandler(read).Handle(t.Context(), app.ListFavoritesQuery{
				UserID: uuid.New(), Limit: tc.limit,
			})

			require.NoError(t, err)
			assert.Len(t, result.Items, tc.wantItems)

			if !tc.wantCursor {
				assert.Empty(t, result.NextCursor)

				return
			}

			decoded, decodeErr := pagination.DecodeCursor(result.NextCursor)
			require.NoError(t, decodeErr)
			assert.Equal(t, result.Items[len(result.Items)-1].ItemID, decoded.ID)
		})
	}
}

func TestListFavoritesPassesCursorThrough(t *testing.T) {
	t.Parallel()

	cursor := pagination.Cursor{CreatedAt: fixedTime, ID: uuid.New()}
	read := &stubReadModel{}

	_, err := app.NewListFavoritesHandler(read).Handle(t.Context(), app.ListFavoritesQuery{
		UserID: uuid.New(), Cursor: cursor, Limit: 10,
	})

	require.NoError(t, err)
	assert.Equal(t, cursor.ID, read.lastFilter.Cursor.ID)
	assert.True(t, cursor.CreatedAt.Equal(read.lastFilter.Cursor.CreatedAt))
}

func TestListFavoritesPropagatesReadError(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("read model is down")

	_, err := app.NewListFavoritesHandler(&stubReadModel{err: sentinel}).Handle(t.Context(), app.ListFavoritesQuery{
		UserID: uuid.New(), Limit: 10,
	})

	require.ErrorIs(t, err, sentinel)
}

func TestAddFavoritePropagatesItemCheckFailure(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("items table unreachable")
	handler := app.NewAddFavoriteHandler(
		newStubRepository(), failingItems{err: sentinel}, passthroughTx{}, fakeClock{}, nil)

	err := handler.Handle(t.Context(), app.AddFavoriteCommand{UserID: uuid.New(), ItemID: uuid.New()})

	require.ErrorIs(t, err, sentinel)
}

func TestAddFavoriteValidatesIDs(t *testing.T) {
	t.Parallel()

	handler := app.NewAddFavoriteHandler(
		newStubRepository(), stubItems{exists: true}, passthroughTx{}, fakeClock{}, nil)

	require.Error(t, handler.Handle(t.Context(), app.AddFavoriteCommand{ItemID: uuid.New()}))
	require.Error(t, handler.Handle(t.Context(), app.AddFavoriteCommand{UserID: uuid.New()}))
}

type failingItems struct{ err error }

func (f failingItems) Exists(context.Context, uuid.UUID) (bool, error) { return false, f.err }
