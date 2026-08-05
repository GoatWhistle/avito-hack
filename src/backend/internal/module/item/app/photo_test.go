package app_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

type stubStorage struct {
	saved       int
	removed     int
	lastType    string
	saveErr     error
	returnedURL string
}

func (s *stubStorage) Save(
	_ context.Context, _ uuid.UUID, content io.Reader, contentType string,
) (app.StoredFile, error) {
	if s.saveErr != nil {
		return app.StoredFile{}, s.saveErr
	}

	body, err := io.ReadAll(content)
	if err != nil {
		return app.StoredFile{}, err
	}

	s.saved++
	s.lastType = contentType

	return app.StoredFile{Name: "photo.jpg", ContentType: contentType, Size: int64(len(body))}, nil
}

func (s *stubStorage) Delete(_ context.Context, _ uuid.UUID, _ string) error {
	s.removed++

	return nil
}

func (s *stubStorage) URL(_ uuid.UUID, name string) string {
	if s.returnedURL != "" {
		return s.returnedURL
	}

	return "/media/" + name
}

type countingPhotos struct {
	stubPhotos
	added    []*domain.Photo
	deleted  int
	countErr error
	addErr   error
}

func (c *countingPhotos) Add(_ context.Context, photo *domain.Photo) error {
	if c.addErr != nil {
		return c.addErr
	}

	c.added = append(c.added, photo)

	return nil
}

func (c *countingPhotos) CountByItemID(_ context.Context, _ uuid.UUID) (int, error) {
	if c.countErr != nil {
		return 0, c.countErr
	}

	return c.count, nil
}

func (c *countingPhotos) ByItemID(_ context.Context, itemID uuid.UUID) ([]*domain.Photo, error) {
	if c.countErr != nil {
		return nil, c.countErr
	}

	result := make([]*domain.Photo, 0, len(c.added))
	for _, photo := range c.added {
		if photo.ItemID() == itemID {
			result = append(result, photo)
		}
	}

	return result, nil
}

func (c *countingPhotos) DeleteByID(_ context.Context, _, _ uuid.UUID) (string, error) {
	c.deleted++

	return "/media/photo.jpg", nil
}

func TestAddPhotoStoresAtNextPosition(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusDraft)
	photos := &countingPhotos{stubPhotos: stubPhotos{count: 3}}
	storage := &stubStorage{}

	handler := app.NewAddPhotoHandler(&stubRepository{item: item}, photos, storage, passthroughTx{}, fakeClock{})

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
		&stubRepository{item: item}, &countingPhotos{}, storage, passthroughTx{}, fakeClock{})

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

	handler := app.NewAddPhotoHandler(&stubRepository{item: item}, photos, storage, passthroughTx{}, fakeClock{})

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

			handler := app.NewAddPhotoHandler(repo, tc.photos, tc.storage, passthroughTx{}, fakeClock{})

			_, err := handler.Handle(t.Context(), app.AddPhotoCommand{
				ItemID: item.ID(), ActorID: ownerID,
				Content: strings.NewReader("binary"), ContentType: "image/jpeg",
			})

			require.ErrorIs(t, err, sentinel)
		})
	}
}

func TestAddPhotoRejectsEmptyStorageURL(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusDraft)

	handler := app.NewAddPhotoHandler(
		&stubRepository{item: item}, &countingPhotos{}, &stubStorage{returnedURL: "   "},
		passthroughTx{}, fakeClock{})

	_, err := handler.Handle(t.Context(), app.AddPhotoCommand{
		ItemID: item.ID(), ActorID: ownerID,
		Content: strings.NewReader("binary"), ContentType: "image/jpeg",
	})

	var invalid *domainerr.InvalidError
	require.ErrorAs(t, err, &invalid)
	assert.Equal(t, "url", invalid.Field)
}

