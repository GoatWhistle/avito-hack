package app

import (
	"context"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/games/domain"
)

type ListGamesQuery struct {
	UserID    uuid.UUID
	ClientDay domain.Day
}

type ListGamesHandler struct {
	registry *domain.Registry
	progress ProgressRepository
}

func NewListGamesHandler(registry *domain.Registry, progress ProgressRepository) *ListGamesHandler {
	return &ListGamesHandler{registry: registry, progress: progress}
}

func (h *ListGamesHandler) Handle(ctx context.Context, q ListGamesQuery) (GameListView, error) {
	streak, err := h.progress.Streak(ctx, q.UserID)
	if err != nil {
		return GameListView{}, err
	}

	anyDaily, err := h.progress.DailyAny(ctx, q.UserID, q.ClientDay)
	if err != nil {
		return GameListView{}, err
	}

	games := h.registry.All()
	out := make([]GameView, 0, len(games))

	for _, game := range games {
		daily, dailyErr := h.progress.Daily(ctx, q.UserID, game.Slug(), q.ClientDay)
		if dailyErr != nil {
			return GameListView{}, dailyErr
		}

		out = append(out, GameView{
			Slug:         game.Slug(),
			TargetStreak: game.TargetStreak(),
			MaxAttempts:  domain.MaxAttemptsOf(game),
			DailyDone:    daily.Attempts > 0,
		})
	}

	return GameListView{
		Games:     out,
		Streak:    toStreakView(streak),
		DailyDone: anyDaily.Attempts > 0,
	}, nil
}
