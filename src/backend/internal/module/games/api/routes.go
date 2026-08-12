package api

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) RegisterRoutes(r chi.Router) {
	r.Route("/games", func(r chi.Router) {
		r.Use(h.deps.Authenticate)

		r.Get("/", h.List)
		r.Get("/{slug}/state", h.GetState)
		r.Post("/{slug}/rounds", h.StartRound)
		r.Post("/{slug}/rounds/{rid}/guess", h.Guess)
		r.Post("/{slug}/reward/claim", h.ClaimReward)
	})
}
