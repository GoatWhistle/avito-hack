package api

import (
	"net/http"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	var req createItemRequest
	if !h.decode(w, r, &req) {
		return
	}

	item, err := h.deps.CreateItem.Handle(r.Context(), app.CreateItemCommand{
		OwnerID:     actor.ID,
		Title:       req.Title,
		Description: req.Description,
		PriceKopeks: req.PriceKopeks,
		Attributes:  req.Attributes,
	})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.Created(w, toItemResponse(item))
}

func (h *Handlers) Update(w http.ResponseWriter, r *http.Request) {
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

	var req updateItemRequest
	if !h.decode(w, r, &req) {
		return
	}

	item, err := h.deps.UpdateItem.Handle(r.Context(), app.UpdateItemCommand{
		ItemID:      id,
		ActorID:     actor.ID,
		Title:       req.Title,
		Description: req.Description,
		PriceKopeks: req.PriceKopeks,
		Attributes:  req.Attributes,
	})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, toItemResponse(item))
}

func (h *Handlers) ChangeStatus(w http.ResponseWriter, r *http.Request) {
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

	var req changeStatusRequest
	if !h.decode(w, r, &req) {
		return
	}

	item, err := h.deps.ChangeStatus.Handle(r.Context(), app.ChangeStatusCommand{
		ItemID: id,
		Actor:  actor,
		Action: app.StatusAction(req.Action),
	})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, toItemResponse(item))
}
