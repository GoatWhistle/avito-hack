package api

import (
	"github.com/go-chi/chi/v5"
)

func (h *SummaryHandlers) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(h.deps.Authenticate)

		r.Get("/summary/today", h.Today)
		r.Get("/summary/history", h.History)
	})
}
