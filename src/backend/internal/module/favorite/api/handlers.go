package api

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/favorite/app"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

type Deps struct {
	AddFavorite    *app.AddFavoriteHandler
	RemoveFavorite *app.RemoveFavoriteHandler
	ListFavorites  *app.ListFavoritesHandler
	Authenticate   func(http.Handler) http.Handler
}

type Handlers struct {
	deps Deps
}

func NewHandlers(deps Deps) *Handlers {
	return &Handlers{deps: deps}
}

func (h *Handlers) Add(w http.ResponseWriter, r *http.Request) {
	actor, itemID, ok := h.actorAndItem(w, r)
	if !ok {
		return
	}

	if err := h.deps.AddFavorite.Handle(r.Context(), app.AddFavoriteCommand{
		UserID: actor.ID,
		ItemID: itemID,
	}); err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.NoContent(w)
}

func (h *Handlers) Remove(w http.ResponseWriter, r *http.Request) {
	actor, itemID, ok := h.actorAndItem(w, r)
	if !ok {
		return
	}

	if err := h.deps.RemoveFavorite.Handle(r.Context(), app.RemoveFavoriteCommand{
		UserID: actor.ID,
		ItemID: itemID,
	}); err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.NoContent(w)
}

func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	page, err := httpx.PageFromRequest(r)
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	result, err := h.deps.ListFavorites.Handle(r.Context(), app.ListFavoritesQuery{
		UserID: actor.ID,
		Cursor: page.Cursor,
		Limit:  page.Limit,
	})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, httpx.NewListResponse(toFavoriteListResponse(result.Items), result.NextCursor))
}

func (h *Handlers) actorAndItem(w http.ResponseWriter, r *http.Request) (auth.Actor, uuid.UUID, bool) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)

		return auth.Actor{}, uuid.Nil, false
	}

	itemID, err := httpx.UUIDParam(r, "id")
	if err != nil {
		apierr.Write(w, r, err)

		return auth.Actor{}, uuid.Nil, false
	}

	return actor, itemID, true
}
