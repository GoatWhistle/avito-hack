package app

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/games/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

type GetStateQuery struct {
	UserID    uuid.UUID
	GameSlug  string
	ClientDay domain.Day
}

type GetStateHandler struct {
	registry *domain.Registry
	rounds   RoundRepository
	progress ProgressRepository
}

func NewGetStateHandler(
	registry *domain.Registry,
	rounds RoundRepository,
	progress ProgressRepository,
) *GetStateHandler {
	return &GetStateHandler{registry: registry, rounds: rounds, progress: progress}
}

func (h *GetStateHandler) Handle(ctx context.Context, q GetStateQuery) (StateView, error) {
	game, err := h.registry.Get(q.GameSlug)
	if err != nil {
		return StateView{}, err
	}

	streak, err := h.progress.Streak(ctx, q.UserID)
	if err != nil {
		return StateView{}, err
	}

	daily, err := h.progress.Daily(ctx, q.UserID, q.GameSlug, q.ClientDay)
	if err != nil {
		return StateView{}, err
	}

	view := StateView{
		Slug:         game.Slug(),
		TargetStreak: game.TargetStreak(),
		MaxAttempts:  domain.MaxAttemptsOf(game),
		Streak:       toStreakView(streak),
		Daily:        DailyView{Attempts: daily.Attempts, BestStreak: daily.BestStreak},
	}

	round, err := h.rounds.ActiveByUser(ctx, q.UserID, q.GameSlug)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			return view, nil
		}

		return StateView{}, err
	}

	resumed, err := game.Resume(ctx, round)
	if err != nil {
		return view, nil //nolint:nilerr // an unreadable payload must not block the state screen
	}

	view.ActiveRound = &ActiveRoundView{
		RoundID:      round.DisplayID(),
		Streak:       round.Streak(),
		AttemptsUsed: domain.AttemptsUsedOf(game, round),
		MaxAttempts:  domain.MaxAttemptsOf(game),
		Prompt:       resumed.Prompt,
	}

	return view, nil
}
