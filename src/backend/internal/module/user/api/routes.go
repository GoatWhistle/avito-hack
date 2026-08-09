package api

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)

		r.Group(func(r chi.Router) {
			r.Use(h.deps.Authenticate)
			r.Post("/refresh", h.Refresh)
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.deps.Authenticate)
			r.Get("/me", h.Me)
			r.Patch("/me", h.UpdateMe)
		})
	})
}
