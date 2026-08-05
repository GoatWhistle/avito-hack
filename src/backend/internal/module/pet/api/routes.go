package api

import (
	"github.com/go-chi/chi/v5"
)

func (h *LeaderboardHandlers) RegisterRoutes(r chi.Router) {
	r.With(h.deps.Authenticate).Get("/leaderboard", h.List)
}

func (h *PetHandlers) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(h.deps.Authenticate)

		r.Get("/pet", h.Get)
		r.Post("/pet/actions/stroke", h.Stroke)
		r.Post("/checkin", h.CheckIn)
		r.Get("/progress", h.Progress)
	})
}

func (h *RewardHandlers) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(h.deps.Authenticate)

		r.Get("/rewards", h.Catalog)
		r.Get("/rewards/my", h.Mine)
		r.Post("/rewards/{id}/activate", h.Activate)
	})
}
