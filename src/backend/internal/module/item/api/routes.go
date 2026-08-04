package api

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) RegisterRoutes(r chi.Router) {
	r.Route("/items", func(r chi.Router) {
		r.With(h.deps.Authenticate).Get("/mine", h.ListMine)
		r.With(h.deps.Authenticate).Post("/", h.Create)
		r.With(h.deps.Authenticate).Patch("/{id}", h.Update)
		r.With(h.deps.Authenticate).Post("/{id}/status", h.ChangeStatus)

		r.With(h.deps.OptionalAuth).Get("/", h.List)
		r.With(h.deps.OptionalAuth).Get("/{id}", h.GetByID)
	})
}
