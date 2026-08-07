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

// @Id addItemPhoto
// @Summary Загрузка фотографии
// @Description Загрузка одного файла через `multipart/form-data`, поле формы — `photo`.
// @Description
// @Description Ограничения:
// @Description
// @Description - Тип определяется по содержимому (сниффинг первых 512 байт),
// @Description а не по заголовку `Content-Type` части. Допустимы только
// @Description `image/jpeg`, `image/png`, `image/webp`.
// @Description - Размер ограничен настройкой `MaxPhotoBytes` (по умолчанию 5 МиБ);
// @Description превышение даёт 400 с `field: photo`.
// @Description - Не более 10 фотографий на объявление (`position` от 0 до 9,
// @Description CHECK-ограничение); одиннадцатая загрузка даёт 409.
// @Description
// @Description Не идемпотентна: каждый вызов создаёт новую фотографию с новым `position`.
// @Description
// @Description Побочный эффект: доменное событие о добавлении фотографии
// @Description асинхронно приносит питомцу опыт.
// @Tags Photos
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Идентификатор объявления" format(uuid)
// @Param photo formData file true "Файл изображения — jpeg, png или webp."
// @Success 201 {object} photoResponse "Фотография загружена"
// @Failure 400 {object} apierr.ErrorEnvelope "Файл отсутствует, слишком велик или имеет недопустимый тип."
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 403 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 409 {object} apierr.ErrorEnvelope "Достигнут лимит в 10 фотографий"
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/items/{id}/photos [post]
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

// @Id listItemPhotos
// @Summary Фотографии объявления
// @Description Отдаёт плоский массив фотографий, отсортированный по `position`.
// @Description Ответ НЕ обёрнут в `{items, next_cursor}` — это голый массив,
// @Description в отличие от прочих списков. Анонимный доступ разрешён.
// @Tags Photos
// @Produce json
// @Param id path string true "Идентификатор объявления" format(uuid)
// @Success 200 {array} photoResponse "Массив фотографий"
// @Failure 400 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Security []
// @Router /api/v1/items/{id}/photos [get]
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

// @Id deleteItemPhoto
// @Summary Удаление фотографии
// @Description Удаляет фотографию объявления. Доступно только владельцу объявления.
// @Description Идемпотентна по результату для клиента, но повторный вызов вернёт 404,
// @Description так как записи уже нет. Тело ответа пустое.
// @Tags Photos
// @Produce json
// @Param id path string true "Идентификатор объявления" format(uuid)
// @Param photoID path string true "Идентификатор фотографии." format(uuid)
// @Success 204 "Фотография удалена"
// @Failure 400 {object} apierr.ErrorEnvelope
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 403 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/items/{id}/photos/{photoID} [delete]
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
