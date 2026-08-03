package api

import (
	"net/http"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

func (h *Handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.UUIDParam(r, "id")
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	actor, _ := auth.ActorFrom(r.Context()) //nolint:errcheck // anonymous access is allowed, the actor is optional

	item, err := h.deps.GetItem.Handle(r.Context(), app.GetItemQuery{ItemID: id, Actor: actor})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, toItemResponse(item))
}

func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	page, err := httpx.PageFromRequest(r)
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	status, err := parseStatusQuery(r)
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	ownerID, _, err := httpx.OptionalUUIDQuery(r, "owner_id")
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	result, err := h.deps.ListItems.Handle(r.Context(), app.ListItemsQuery{
		Status:  status,
		OwnerID: ownerID,
		Search:  r.URL.Query().Get("search"),
		Cursor:  page.Cursor,
		Limit:   page.Limit,
	})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, httpx.NewListResponse(toListResponse(result.Items), result.NextCursor))
}

func (h *Handlers) ListMine(w http.ResponseWriter, r *http.Request) {
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

	status, err := parseStatusQuery(r)
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	result, err := h.deps.ListItems.Handle(r.Context(), app.ListItemsQuery{
		Status:  status,
		OwnerID: actor.ID,
		Cursor:  page.Cursor,
		Limit:   page.Limit,
	})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, httpx.NewListResponse(toListResponse(result.Items), result.NextCursor))
}

func parseStatusQuery(r *http.Request) (domain.Status, error) {
	raw := r.URL.Query().Get("status")
	if raw == "" {
		return "", nil
	}

	status := domain.Status(raw)
	if !status.Valid() {
		return "", domainerr.NewInvalid("status", "allowed values: draft moderation published archived")
	}

	return status, nil
}
