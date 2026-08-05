package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/domain"
)

func TestNewPhoto_Validation(t *testing.T) {
	t.Parallel()

	itemID := uuid.New()

	tests := []struct {
		name     string
		itemID   uuid.UUID
		url      string
		position int
		wantErr  bool
	}{
		{name: "valid", itemID: itemID, url: "/uploads/a/b.jpg", position: 0},
		{name: "last allowed position", itemID: itemID, url: "/uploads/a/b.jpg", position: 9},
		{name: "blank url", itemID: itemID, url: "   ", position: 0, wantErr: true},
		{name: "missing item", itemID: uuid.Nil, url: "/uploads/a/b.jpg", position: 0, wantErr: true},
		{name: "position out of range", itemID: itemID, url: "/uploads/a/b.jpg", position: 10, wantErr: true},
		{name: "negative position", itemID: itemID, url: "/uploads/a/b.jpg", position: -1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			photo, err := domain.NewPhoto(domain.NewPhotoParams{
				ItemID:   tt.itemID,
				URL:      tt.url,
				Position: tt.position,
				Now:      fixedTime,
			})

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.itemID, photo.ItemID())
			require.Equal(t, tt.position, photo.Position())
			require.NotEqual(t, uuid.Nil, photo.ID())
		})
	}
}

func TestMaxPhotosPerItem(t *testing.T) {
	t.Parallel()

	require.Equal(t, 10, domain.MaxPhotosPerItem)
	require.Error(t, domain.ErrTooManyPhotos())
}
