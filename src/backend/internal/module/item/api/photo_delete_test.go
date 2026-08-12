package api_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/domain"
)

func (f *fixture) seedPhoto(t *testing.T, itemID uuid.UUID) *domain.Photo {
	t.Helper()

	photo := domain.RestorePhoto(uuid.New(), domain.NewDisplayID(), itemID, "/uploads/photo.jpg", 0, fixedTime)
	require.NoError(t, f.photos.Add(context.Background(), photo))

	return photo
}

func TestDeletePhotoRemovesIt(t *testing.T) {
	t.Parallel()

	actor := userActor()
	f := newFixture(t, &actor)
	item := f.seedItem(t, actor.ID, domain.StatusDraft)
	photo := f.seedPhoto(t, item.ID())

	path := "/items/" + item.DisplayID() + "/photos/" + photo.DisplayID()
	rec := f.do(t, http.MethodDelete, path, "")

	require.Equal(t, http.StatusNoContent, rec.Code, rec.Body.String())
	assert.Empty(t, rec.Body.String())

	remaining, err := f.photos.ByItemID(context.Background(), item.ID())
	require.NoError(t, err)
	assert.Empty(t, remaining)
}

func TestDeletePhotoRequiresAuthentication(t *testing.T) {
	t.Parallel()

	f := newFixture(t, nil)

	path := "/items/" + uuid.New().String() + "/photos/" + uuid.New().String()
	rec := f.do(t, http.MethodDelete, path, "")

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDeletePhotoForbiddenForNonOwner(t *testing.T) {
	t.Parallel()

	actor := userActor()
	f := newFixture(t, &actor)
	item := f.seedItem(t, uuid.New(), domain.StatusDraft)
	photo := f.seedPhoto(t, item.ID())

	path := "/items/" + item.DisplayID() + "/photos/" + photo.DisplayID()
	rec := f.do(t, http.MethodDelete, path, "")

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestDeletePhotoUnknownItemIsNotFound(t *testing.T) {
	t.Parallel()

	actor := userActor()
	f := newFixture(t, &actor)

	path := "/items/" + uuid.New().String() + "/photos/" + uuid.New().String()
	rec := f.do(t, http.MethodDelete, path, "")

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeletePhotoRejectsMalformedIDs(t *testing.T) {
	t.Parallel()

	actor := userActor()
	f := newFixture(t, &actor)
	item := f.seedItem(t, actor.ID, domain.StatusDraft)

	tests := []struct {
		name string
		path string
	}{
		{name: "bad photo id", path: "/items/" + item.DisplayID() + "/photos/nope"},
		{name: "photo uuid is not accepted", path: "/items/" + item.DisplayID() + "/photos/" + uuid.New().String()},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := f.do(t, http.MethodDelete, tc.path, "")

			assert.Equal(t, http.StatusNotFound, rec.Code)
		})
	}
}

func TestDeletePhotoRejectsUnknownItemID(t *testing.T) {
	t.Parallel()

	actor := userActor()
	f := newFixture(t, &actor)

	rec := f.do(t, http.MethodDelete, "/items/nope/photos/"+uuid.New().String(), "")

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
