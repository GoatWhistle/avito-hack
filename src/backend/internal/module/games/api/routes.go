package api

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) RegisterRoutes(r chi.Router) {
	r.Route("/games", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.deps.OptionalAuth)

			r.Get("/photo/{token}", h.HiddenPhoto)
		})

		r.Group(func(r chi.Router) {
			r.Use(h.deps.Authenticate)

			r.Get("/", h.List)
			r.Post("/reward/claim", h.ClaimReward)
			r.Get("/{slug}/state", h.GetState)
			r.Post("/{slug}/rounds", h.StartRound)
			r.Post("/{slug}/rounds/{rid}/guess", h.Guess)
		})
	})
}
