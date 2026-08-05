package api

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(h.deps.Authenticate)

		r.Get("/raccoon/profile", h.GetProfile)
		r.Get("/badges", h.GetBadges)
		r.Post("/rewards/claim", h.ClaimReward)
	})
}
