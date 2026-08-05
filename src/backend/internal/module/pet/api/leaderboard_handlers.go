package api

import (
	"net/http"

	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/httpx"
	"github.com/avito-hack/backend/internal/shared/pagination"
)

type leaderboardEntryResponse struct {
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	Level      int    `json:"level"`
	XP         int    `json:"xp"`
	StreakDays int    `json:"streak_days"`
	Rank       int    `json:"rank"`
}

type leaderboardResponse struct {
	Items      []leaderboardEntryResponse `json:"items"`
	MyRank     *int                       `json:"my_rank"`
	NextCursor string                     `json:"next_cursor,omitempty"`
}

type LeaderboardDeps struct {
	Leaderboard  *app.LeaderboardHandler
	Authenticate func(http.Handler) http.Handler
}

type LeaderboardHandlers struct {
	deps LeaderboardDeps
}

func NewLeaderboardHandlers(deps LeaderboardDeps) *LeaderboardHandlers {
	return &LeaderboardHandlers{deps: deps}
}

func (h *LeaderboardHandlers) List(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	query, err := leaderboardQueryFromRequest(r)
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	query.UserID = actor.ID

	result, err := h.deps.Leaderboard.Handle(r.Context(), query)
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, toLeaderboardResponse(result))
}

func leaderboardQueryFromRequest(r *http.Request) (app.LeaderboardQuery, error) {
	limit, err := httpx.IntQuery(r, "limit", pagination.DefaultLimit)
	if err != nil {
		return app.LeaderboardQuery{}, err
	}

	cursor, err := pagination.DecodeLeaderboardCursor(r.URL.Query().Get("cursor"))
	if err != nil {
		return app.LeaderboardQuery{}, domainerr.NewInvalid("cursor", "invalid cursor")
	}

	return app.LeaderboardQuery{
		Cursor: cursor,
		Limit:  pagination.NormalizeLimit(limit),
		Around: r.URL.Query().Get("around") == "me",
	}, nil
}

func toLeaderboardResponse(result app.LeaderboardResult) leaderboardResponse {
	items := make([]leaderboardEntryResponse, 0, len(result.Items))

	for _, entry := range result.Items {
		items = append(items, leaderboardEntryResponse{
			UserID:     entry.UserID.String(),
			Name:       entry.Name,
			Level:      entry.Level,
			XP:         entry.XP,
			StreakDays: entry.StreakDays,
			Rank:       entry.Rank,
		})
	}

	return leaderboardResponse{Items: items, MyRank: result.MyRank, NextCursor: result.NextCursor}
}
