package api

import (
	"net/http"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

type validator interface {
	Struct(dst any) error
}

type Deps struct {
	CreateItem    *app.CreateItemHandler
	UpdateItem    *app.UpdateItemHandler
	ChangeStatus  *app.ChangeStatusHandler
	GetItem       *app.GetItemHandler
	ListItems     *app.ListItemsHandler
	AddPhoto      *app.AddPhotoHandler
	ListPhotos    *app.ListPhotosHandler
	DeletePhoto   *app.DeletePhotoHandler
	Validator     validator
	Authenticate  func(http.Handler) http.Handler
	OptionalAuth  func(http.Handler) http.Handler
	MaxBodyBytes  int64
	MaxPhotoBytes int64
}

type Handlers struct {
	deps Deps
}

func NewHandlers(deps Deps) *Handlers {
	return &Handlers{deps: deps}
}

func (h *Handlers) decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := httpx.DecodeJSON(w, r, h.deps.MaxBodyBytes, dst); err != nil {
		apierr.WriteBadRequest(w, r, err.Error())
		return false
	}

	if err := h.deps.Validator.Struct(dst); err != nil {
		apierr.Write(w, r, err)
		return false
	}

	return true
}
