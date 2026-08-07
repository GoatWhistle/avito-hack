package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

type summaryService interface {
	Today(ctx context.Context, userID uuid.UUID) (*domain.DailySummary, error)
	History(ctx context.Context, q app.SummaryHistoryQuery) (app.SummaryHistoryResult, error)
}

type SummaryDeps struct {
	Summaries    summaryService
	Authenticate func(http.Handler) http.Handler
}

type SummaryHandlers struct {
	deps SummaryDeps
}

func NewSummaryHandlers(deps SummaryDeps) *SummaryHandlers {
	return &SummaryHandlers{deps: deps}
}

func (h *SummaryHandlers) Today(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorOf(w, r)
	if !ok {
		return
	}

	summary, err := h.deps.Summaries.Today(r.Context(), actor.ID)
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	httpx.OK(w, toSummaryPayload(summary))
}

func (h *SummaryHandlers) History(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorOf(w, r)
	if !ok {
		return
	}

	page, err := httpx.PageFromRequest(r)
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	result, err := h.deps.Summaries.History(r.Context(), app.SummaryHistoryQuery{
		UserID: actor.ID, Cursor: page.Cursor, Limit: page.Limit,
	})
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	items := make([]summaryPayload, 0, len(result.Items))
	for _, summary := range result.Items {
		items = append(items, toSummaryPayload(summary))
	}

	httpx.OK(w, httpx.NewListResponse(items, result.NextCursor))
}
