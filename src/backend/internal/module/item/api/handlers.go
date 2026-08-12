package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

type Deps struct {
	Items         domain.Repository
	Owners        app.OwnerProvider
	CreateItem    *app.CreateItemHandler
	UpdateItem    *app.UpdateItemHandler
	ChangeStatus  *app.ChangeStatusHandler
	GetItem       *app.GetItemHandler
	ViewItem      *app.ViewItemHandler
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

func (h *Handlers) resolveDisplayID(w http.ResponseWriter, r *http.Request, name string) (*domain.Item, bool) {
	displayID := chi.URLParam(r, name)

	if !domain.IsValidDisplayID(displayID) {
		apierr.Write(w, r, domain.ErrItemNotFound)
		return nil, false
	}

	item, err := h.deps.Items.ByDisplayID(r.Context(), displayID)
	if err != nil {
		apierr.Write(w, r, err)
		return nil, false
	}

	return item, true
}

func (h *Handlers) ownerDisplayID(ctx context.Context, ownerID uuid.UUID) string {
	if h.deps.Owners == nil {
		return ""
	}

	owners, err := h.deps.Owners.ByIDs(ctx, []uuid.UUID{ownerID})
	if err != nil {
		return ""
	}

	return owners[ownerID].DisplayID
}
