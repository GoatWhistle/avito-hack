package api

import (
	"bytes"
	"io"
	"net/http"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

const (
	photoFormField = "photo"
	sniffLen       = 512
)

var allowedPhotoTypes = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/webp": {},
}

func (h *Handlers) AddPhoto(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	id, err := httpx.UUIDParam(r, "id")
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	content, contentType, err := h.readPhoto(w, r)
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	photo, err := h.deps.AddPhoto.Handle(r.Context(), app.AddPhotoCommand{
		ItemID:      id,
		ActorID:     actor.ID,
		Content:     content,
		ContentType: contentType,
	})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.Created(w, toPhotoResponse(photo))
}

func (h *Handlers) readPhoto(w http.ResponseWriter, r *http.Request) (io.Reader, string, error) {
	r.Body = http.MaxBytesReader(w, r.Body, h.deps.MaxPhotoBytes)

	file, header, err := r.FormFile(photoFormField)
	if err != nil {
		return nil, "", domainerr.NewInvalid(photoFormField, "multipart file field is required")
	}
	defer func() { _ = file.Close() }() //nolint:errcheck // the upload is read fully before closing

	if header.Size > h.deps.MaxPhotoBytes {
		return nil, "", domainerr.NewInvalid(photoFormField, "file is too large")
	}

	raw, err := io.ReadAll(io.LimitReader(file, h.deps.MaxPhotoBytes+1))
	if err != nil {
		return nil, "", domainerr.NewInvalid(photoFormField, "cannot read uploaded file")
	}

	if int64(len(raw)) > h.deps.MaxPhotoBytes {
		return nil, "", domainerr.NewInvalid(photoFormField, "file is too large")
	}

	contentType, err := detectPhotoType(raw)
	if err != nil {
		return nil, "", err
	}

	return bytes.NewReader(raw), contentType, nil
}

func detectPhotoType(raw []byte) (string, error) {
	head := raw
	if len(head) > sniffLen {
		head = head[:sniffLen]
	}

	contentType := http.DetectContentType(head)
	if _, ok := allowedPhotoTypes[contentType]; !ok {
		return "", domainerr.NewInvalid(photoFormField, "only jpeg, png and webp images are allowed")
	}

	return contentType, nil
}

func (h *Handlers) ListPhotos(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.UUIDParam(r, "id")
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	photos, err := h.deps.ListPhotos.Handle(r.Context(), id)
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, toPhotoListResponse(photos))
}

func (h *Handlers) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	itemID, err := httpx.UUIDParam(r, "id")
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	photoID, err := httpx.UUIDParam(r, "photoID")
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	if err := h.deps.DeletePhoto.Handle(r.Context(), app.DeletePhotoCommand{
		ItemID:  itemID,
		PhotoID: photoID,
		ActorID: actor.ID,
	}); err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.NoContent(w)
}
