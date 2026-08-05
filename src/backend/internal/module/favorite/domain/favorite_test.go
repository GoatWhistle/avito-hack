package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/favorite/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

var favoriteTime = time.Date(2026, time.September, 15, 18, 30, 0, 0, time.UTC)

func TestNewFavorite(t *testing.T) {
	t.Parallel()

	userID, itemID := uuid.New(), uuid.New()

	favorite, err := domain.New(userID, itemID, favoriteTime)

	require.NoError(t, err)
	assert.Equal(t, userID, favorite.UserID())
	assert.Equal(t, itemID, favorite.ItemID())
	assert.Equal(t, favoriteTime, favorite.CreatedAt())
}

func TestNewFavoriteValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		userID    uuid.UUID
		itemID    uuid.UUID
		wantField string
	}{
		{name: "missing user", itemID: uuid.New(), wantField: "user_id"},
		{name: "missing item", userID: uuid.New(), wantField: "item_id"},
		{name: "both missing", wantField: "user_id"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			favorite, err := domain.New(tc.userID, tc.itemID, favoriteTime)

			require.Nil(t, favorite)

			var invalid *domainerr.InvalidError
			require.ErrorAs(t, err, &invalid)
			assert.Equal(t, tc.wantField, invalid.Field)
		})
	}
}

func TestRestoreFavorite(t *testing.T) {
	t.Parallel()

	userID, itemID := uuid.New(), uuid.New()

	favorite := domain.Restore(userID, itemID, favoriteTime)

	assert.Equal(t, userID, favorite.UserID())
	assert.Equal(t, itemID, favorite.ItemID())
	assert.Equal(t, favoriteTime, favorite.CreatedAt())
}

func TestFavoriteErrorsMapToSentinels(t *testing.T) {
	t.Parallel()

	assert.ErrorIs(t, domain.ErrFavoriteNotFound, domainerr.ErrNotFound)
	assert.ErrorIs(t, domain.ErrItemNotFound, domainerr.ErrNotFound)
	assert.Contains(t, domain.ErrFavoriteNotFound.Error(), "favorite not found")
	assert.Contains(t, domain.ErrItemNotFound.Error(), "item not found")
}
