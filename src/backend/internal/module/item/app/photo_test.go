package app_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

func TestAddPhotoStoresAtNextPosition(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusDraft)
	photos := &countingPhotos{stubPhotos: stubPhotos{count: 3}}
	storage := &stubStorage{}

	handler := app.NewAddPhotoHandler(&stubRepository{item: item}, photos, storage, passthroughTx{}, fakeClock{}, nil)

	photo, err := handler.Handle(t.Context(), app.AddPhotoCommand{
		ItemID: item.ID(), ActorID: ownerID,
		Content: strings.NewReader("binary"), ContentType: "image/jpeg",
	})

	require.NoError(t, err)
	assert.Equal(t, 3, photo.Position())
	assert.Equal(t, "/media/photo.jpg", photo.URL())
	assert.Equal(t, item.ID(), photo.ItemID())
	assert.Equal(t, fixedTime, photo.CreatedAt())
	assert.Equal(t, 1, storage.saved)
	assert.Equal(t, "image/jpeg", storage.lastType)
	assert.Len(t, photos.added, 1)
}

func TestAddPhotoRejectsNonOwner(t *testing.T) {
	t.Parallel()

	item := itemWithStatus(t, uuid.New(), domain.StatusDraft)
	storage := &stubStorage{}

	handler := app.NewAddPhotoHandler(
		&stubRepository{item: item}, &countingPhotos{}, storage, passthroughTx{}, fakeClock{}, nil)

	_, err := handler.Handle(t.Context(), app.AddPhotoCommand{
		ItemID: item.ID(), ActorID: uuid.New(),
		Content: strings.NewReader("binary"), ContentType: "image/jpeg",
	})

	require.ErrorIs(t, err, domainerr.ErrForbidden)
	assert.Equal(t, storage.saved, storage.removed)
}

func TestAddPhotoEnforcesLimit(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusDraft)
	storage := &stubStorage{}
	photos := &countingPhotos{stubPhotos: stubPhotos{count: domain.MaxPhotosPerItem}}

	handler := app.NewAddPhotoHandler(&stubRepository{item: item}, photos, storage, passthroughTx{}, fakeClock{}, nil)

	_, err := handler.Handle(t.Context(), app.AddPhotoCommand{
		ItemID: item.ID(), ActorID: ownerID,
		Content: strings.NewReader("binary"), ContentType: "image/jpeg",
	})

	require.ErrorIs(t, err, domainerr.ErrConflict)
	assert.Equal(t, storage.saved, storage.removed)
	assert.Empty(t, photos.added)
}

func TestAddPhotoPropagatesFailures(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	sentinel := errors.New("infrastructure failure")

	tests := []struct {
		name    string
		photos  *countingPhotos
		storage *stubStorage
		repo    domain.Repository
	}{
		{
			name:    "item load failure",
			photos:  &countingPhotos{},
			storage: &stubStorage{},
			repo:    &failingRepository{err: sentinel},
		},
		{
			name:    "count failure",
			photos:  &countingPhotos{countErr: sentinel},
			storage: &stubStorage{},
		},
		{
			name:    "storage failure",
			photos:  &countingPhotos{},
			storage: &stubStorage{saveErr: sentinel},
		},
		{
			name:    "persist failure",
			photos:  &countingPhotos{addErr: sentinel},
			storage: &stubStorage{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			item := itemWithStatus(t, ownerID, domain.StatusDraft)
			repo := tc.repo
			if repo == nil {
				repo = &stubRepository{item: item}
			}

			handler := app.NewAddPhotoHandler(repo, tc.photos, tc.storage, passthroughTx{}, fakeClock{}, nil)

			_, err := handler.Handle(t.Context(), app.AddPhotoCommand{
				ItemID: item.ID(), ActorID: ownerID,
				Content: strings.NewReader("binary"), ContentType: "image/jpeg",
			})

			require.ErrorIs(t, err, sentinel)
		})
	}
}

func TestDeletePhotoRemovesFileAfterCommit(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusDraft)
	photos := &countingPhotos{}
	storage := &stubStorage{}

	handler := app.NewDeletePhotoHandler(&stubRepository{item: item}, photos, storage, passthroughTx{}, fakeClock{}, nil)

	err := handler.Handle(t.Context(), app.DeletePhotoCommand{
		ItemID: item.ID(), PhotoDisplayID: domain.NewDisplayID(), ActorID: ownerID,
	})

	require.NoError(t, err)
	assert.Equal(t, 1, photos.deleted)
	assert.Equal(t, 1, storage.removed)
}

func TestDeletePhotoKeepsFileOnRollback(t *testing.T) {
	t.Parallel()

	item := itemWithStatus(t, uuid.New(), domain.StatusDraft)
	photos := &countingPhotos{}
	storage := &stubStorage{}

	handler := app.NewDeletePhotoHandler(&stubRepository{item: item}, photos, storage, passthroughTx{}, fakeClock{}, nil)

	err := handler.Handle(t.Context(), app.DeletePhotoCommand{
		ItemID: item.ID(), PhotoDisplayID: domain.NewDisplayID(), ActorID: uuid.New(),
	})

	require.ErrorIs(t, err, domainerr.ErrForbidden)
	assert.Equal(t, 0, photos.deleted)
	assert.Equal(t, 0, storage.removed)
}

func TestAddPhotoRejectsEmptyStorageURL(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusDraft)

	handler := app.NewAddPhotoHandler(
		&stubRepository{item: item}, &countingPhotos{}, &stubStorage{returnedURL: "   "},
		passthroughTx{}, fakeClock{}, nil)

	_, err := handler.Handle(t.Context(), app.AddPhotoCommand{
		ItemID: item.ID(), ActorID: ownerID,
		Content: strings.NewReader("binary"), ContentType: "image/jpeg",
	})

	var invalid *domainerr.InvalidError
	require.ErrorAs(t, err, &invalid)
	assert.Equal(t, "url", invalid.Field)
}
