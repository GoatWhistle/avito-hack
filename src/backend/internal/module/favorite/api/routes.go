package api

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) RegisterRoutes(r chi.Router) {
	r.With(h.deps.Authenticate).Post("/items/{id}/favorite", h.Add)
	r.With(h.deps.Authenticate).Delete("/items/{id}/favorite", h.Remove)
	r.With(h.deps.Authenticate).Get("/favorites", h.List)
}
