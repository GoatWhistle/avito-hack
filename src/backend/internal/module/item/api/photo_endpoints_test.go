package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/domain"
)

func jpegBytes(size int) []byte {
	raw := make([]byte, size)
	copy(raw, []byte{0xFF, 0xD8, 0xFF, 0xE0})

	return raw
}

func multipartBody(
	t *testing.T, field, filename string, content []byte,
) (body *bytes.Buffer, contentType string) {
	t.Helper()

	body = &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(field, filename)
	require.NoError(t, err)

	_, err = part.Write(content)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	return body, writer.FormDataContentType()
}

func (f *fixture) upload(t *testing.T, path, field, filename string, content []byte) *httptest.ResponseRecorder {
	t.Helper()

	body, contentType := multipartBody(t, field, filename, content)

	req := httptest.NewRequest(http.MethodPost, path, body)
	req.Header.Set("Content-Type", contentType)

	rec := httptest.NewRecorder()
	f.router.ServeHTTP(rec, req)

	return rec
}

func TestAddPhotoStoresUploadedImage(t *testing.T) {
	t.Parallel()

	actor := userActor()
	f := newFixture(t, &actor)
	item := f.seedItem(t, actor.ID, domain.StatusDraft)

	rec := f.upload(t, "/items/"+item.ID().String()+"/photos", "photo", "cat.jpg", jpegBytes(64))

	require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "/uploads/photo.jpg", body["url"])

	stored, err := f.photos.ByItemID(context.Background(), item.ID())
	require.NoError(t, err)
	assert.Len(t, stored, 1)
}

func TestAddPhotoRequiresMultipartField(t *testing.T) {
	t.Parallel()

	actor := userActor()
	f := newFixture(t, &actor)
	item := f.seedItem(t, actor.ID, domain.StatusDraft)

	rec := f.upload(t, "/items/"+item.ID().String()+"/photos", "attachment", "cat.jpg", jpegBytes(64))

	require.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "multipart file field is required")
}

func TestAddPhotoRejectsNonMultipartBody(t *testing.T) {
	t.Parallel()

	actor := userActor()
	f := newFixture(t, &actor)
	item := f.seedItem(t, actor.ID, domain.StatusDraft)

	rec := f.do(t, http.MethodPost, "/items/"+item.ID().String()+"/photos", `{"photo":"x"}`)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAddPhotoRejectsDisallowedMIME(t *testing.T) {
	t.Parallel()

	actor := userActor()
	f := newFixture(t, &actor)
	item := f.seedItem(t, actor.ID, domain.StatusDraft)

	gif := append([]byte("GIF89a"), make([]byte, 64)...)

	rec := f.upload(t, "/items/"+item.ID().String()+"/photos", "photo", "cat.gif", gif)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "only jpeg, png and webp images are allowed")
}

func TestAddPhotoRejectsOversizedUpload(t *testing.T) {
	t.Parallel()

	actor := userActor()
	f := newFixture(t, &actor)
	item := f.seedItem(t, actor.ID, domain.StatusDraft)

	rec := f.upload(t, "/items/"+item.ID().String()+"/photos", "photo", "huge.jpg", jpegBytes((1<<20)+1024))

	require.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), photoFieldName)
}

func TestAddPhotoRejectsBadItemID(t *testing.T) {
	t.Parallel()

	actor := userActor()
	f := newFixture(t, &actor)

	rec := f.upload(t, "/items/not-a-uuid/photos", "photo", "cat.jpg", jpegBytes(64))

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAddPhotoRequiresAuthentication(t *testing.T) {
	t.Parallel()

	f := newFixture(t, nil)

	rec := f.upload(t, "/items/"+uuid.New().String()+"/photos", "photo", "cat.jpg", jpegBytes(64))

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAddPhotoForbiddenForNonOwner(t *testing.T) {
	t.Parallel()

	actor := userActor()
	f := newFixture(t, &actor)
	item := f.seedItem(t, uuid.New(), domain.StatusDraft)

	rec := f.upload(t, "/items/"+item.ID().String()+"/photos", "photo", "cat.jpg", jpegBytes(64))

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAddPhotoUnknownItemIsNotFound(t *testing.T) {
	t.Parallel()

	actor := userActor()
	f := newFixture(t, &actor)

	rec := f.upload(t, "/items/"+uuid.New().String()+"/photos", "photo", "cat.jpg", jpegBytes(64))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

const photoFieldName = "photo"
