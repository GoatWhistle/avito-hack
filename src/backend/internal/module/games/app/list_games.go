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

func (h *ListGamesHandler) Handle(ctx context.Context, q ListGamesQuery) ([]GameView, error) {
	games := h.registry.All()
	out := make([]GameView, 0, len(games))

	for _, game := range games {
		streak, err := h.progress.Streak(ctx, q.UserID, game.Slug())
		if err != nil {
			return nil, err
		}

		daily, err := h.progress.Daily(ctx, q.UserID, game.Slug(), q.ClientDay)
		if err != nil {
			return nil, err
		}

		out = append(out, GameView{
			Slug:         game.Slug(),
			TargetStreak: game.TargetStreak(),
			DailyDone:    daily.Attempts > 0,
			Streak:       toStreakView(streak),
		})
	}

	return out, nil
}
