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
	ByDisplayIDForUpdate(ctx context.Context, displayID string) (*domain.Round, error)
	ActiveByUser(ctx context.Context, userID uuid.UUID, gameSlug string) (*domain.Round, error)
	LockUserGame(ctx context.Context, userID uuid.UUID, gameSlug string) error
}

type ScoreRepository interface {
	BestScore(ctx context.Context, userID uuid.UUID, gameSlug string) (int, error)
	SaveBestScore(ctx context.Context, userID uuid.UUID, gameSlug string, score int) (int, error)
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
	DailyAny(ctx context.Context, userID uuid.UUID, day domain.Day) (domain.DailyProgress, error)
	SaveDailyAny(ctx context.Context, userID uuid.UUID, day domain.Day, progress domain.DailyProgress) error
	Streak(ctx context.Context, userID uuid.UUID) (domain.Streak, error)
	StreakForUpdate(ctx context.Context, userID uuid.UUID) (domain.Streak, error)
	SaveStreak(ctx context.Context, userID uuid.UUID, streak domain.Streak) error
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
	RoundID      string
	Streak       int
	AttemptsUsed int
	MaxAttempts  int
	Prompt       []byte
}

type GameView struct {
	Slug         string
	TargetStreak int
	MaxAttempts  int
	DailyDone    bool
}

type GameListView struct {
	Games     []GameView
	Streak    StreakView
	DailyDone bool
}

type StateView struct {
	Slug         string
	TargetStreak int
	MaxAttempts  int
	Streak       StreakView
	Daily        DailyView
	ActiveRound  *ActiveRoundView
	BestScore    int
}

type StartRoundResult struct {
	RoundID      string
	Streak       int
	TargetStreak int
	AttemptsUsed int
	MaxAttempts  int
	Prompt       []byte
}

type GuessResult struct {
	Correct          bool
	Progress         domain.Progress
	Reveal           []byte
	Streak           int
	AttemptsUsed     int
	MaxAttempts      int
	State            domain.State
	Prompt           []byte
	AttemptCompleted bool
	StreakAfter      *StreakView
	BestScore        int
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
