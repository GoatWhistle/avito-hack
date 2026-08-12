package app

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/games/domain"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(context.Context) error) error
}

type Clock interface {
	Now() time.Time
}

type RoundRepository interface {
	Save(ctx context.Context, round *domain.Round) error
	ByDisplayID(ctx context.Context, displayID string) (*domain.Round, error)
	ActiveByUser(ctx context.Context, userID uuid.UUID, gameSlug string) (*domain.Round, error)
}

type ProgressRepository interface {
	Daily(ctx context.Context, userID uuid.UUID, gameSlug string, day domain.Day) (domain.DailyProgress, error)
	SaveDaily(
		ctx context.Context,
		userID uuid.UUID,
		gameSlug string,
		day domain.Day,
		progress domain.DailyProgress,
	) error
	Streak(ctx context.Context, userID uuid.UUID, gameSlug string) (domain.Streak, error)
	SaveStreak(ctx context.Context, userID uuid.UUID, gameSlug string, streak domain.Streak) error
}

type StreakView struct {
	CurrentDays int
	BestDays    int
	RewardReady bool
}

type DailyView struct {
	Attempts   int
	BestStreak int
}

type ActiveRoundView struct {
	RoundID string
	Streak  int
	Prompt  []byte
}

type GameView struct {
	Slug         string
	TargetStreak int
	DailyDone    bool
	Streak       StreakView
}

type StateView struct {
	Slug         string
	TargetStreak int
	Streak       StreakView
	Daily        DailyView
	ActiveRound  *ActiveRoundView
}

type StartRoundResult struct {
	RoundID      string
	Streak       int
	TargetStreak int
	Prompt       []byte
}

type GuessResult struct {
	Correct          bool
	Reveal           []byte
	Streak           int
	State            domain.State
	Prompt           []byte
	AttemptCompleted bool
	StreakAfter      *StreakView
}

type ClaimRewardResult struct {
	Code string
}

func toStreakView(s domain.Streak) StreakView {
	return StreakView{
		CurrentDays: s.CurrentDays,
		BestDays:    s.BestDays,
		RewardReady: domain.RewardReady(s),
	}
}
