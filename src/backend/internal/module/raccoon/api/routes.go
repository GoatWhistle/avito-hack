package api

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) RegisterRoutes(r chi.Router) {
	r.Route("/raccoon", func(r chi.Router) {
		r.Use(h.deps.Authenticate)
		r.Get("/profile", h.GetProfile)
	})

	r.Route("/badges", func(r chi.Router) {
		r.Use(h.deps.Authenticate)
		r.Get("/", h.GetBadges)
	})

	r.Route("/rewards", func(r chi.Router) {
		r.Use(h.deps.Authenticate)
		r.Post("/claim", h.ClaimReward)
	})
}
