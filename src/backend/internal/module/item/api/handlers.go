package api

import (
	"net/http"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

type Deps struct {
	CreateItem    *app.CreateItemHandler
	UpdateItem    *app.UpdateItemHandler
	ChangeStatus  *app.ChangeStatusHandler
	GetItem       *app.GetItemHandler
	ListItems     *app.ListItemsHandler
	AddPhoto      *app.AddPhotoHandler
	ListPhotos    *app.ListPhotosHandler
	DeletePhoto   *app.DeletePhotoHandler
	Validator     httpx.Validator
	Authenticate  func(http.Handler) http.Handler
	OptionalAuth  func(http.Handler) http.Handler
	MaxBodyBytes  int64
	MaxPhotoBytes int64
}

type Handlers struct {
	deps    Deps
	decoder httpx.Decoder
}

func NewHandlers(deps Deps) *Handlers {
	return &Handlers{deps: deps, decoder: httpx.NewDecoder(deps.Validator, deps.MaxBodyBytes)}
}

func (h *Handlers) decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	return h.decoder.Decode(w, r, dst)
}
