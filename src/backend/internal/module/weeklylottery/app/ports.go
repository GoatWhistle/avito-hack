package app

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/weeklylottery/domain"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(context.Context) error) error
}

type Clock interface {
	Now() time.Time
}

type Repository interface {
	ByUserWeek(ctx context.Context, userID uuid.UUID, weekStart time.Time, lock bool) (domain.Run, error)
	ByDisplayID(ctx context.Context, displayID string, lock bool) (domain.Run, error)
	Create(ctx context.Context, run domain.Run) (bool, error)
	Save(ctx context.Context, run domain.Run) error
}

type CodeSigner interface {
	Issue(userID uuid.UUID, rewardID string) (string, error)
	Verify(userID uuid.UUID, rewardID, code string) error
}

type SlotView struct {
	Index  int
	Opened bool
	Symbol domain.Symbol
}

type PrizeView struct {
	ID           string
	Title        string
	Description  string
	BenefitType  string
	BenefitValue int
	ScopeType    string
	ScopeValue   string
	Code         string
	ExpiresAt    time.Time
}

type PrizeCatalogView struct {
	ID           string
	Symbol       domain.Symbol
	Title        string
	Description  string
	BenefitType  string
	BenefitValue int
	ScopeType    string
	ScopeValue   string
}

type RunView struct {
	ID        string
	State     domain.State
	Slots     []SlotView
	Prize     *PrizeView
	CreatedAt time.Time
}

type StateView struct {
	Available     bool
	WeekStart     time.Time
	NextAvailable time.Time
	Run           *RunView
}

type RevealResult struct {
	Index  int
	Symbol domain.Symbol
	Run    RunView
}

func toRunView(run domain.Run) RunView {
	slots := make([]SlotView, domain.BoardSize)
	for index := range domain.BoardSize {
		slots[index] = SlotView{Index: index, Opened: run.IsOpened(index)}
		if slots[index].Opened {
			slots[index].Symbol = run.Board[index]
		}
	}

	view := RunView{ID: run.DisplayID, State: run.State, Slots: slots, CreatedAt: run.CreatedAt}
	if run.State == domain.StateWon {
		if prize, ok := domain.PrizeByID(run.PrizeID); ok && run.RewardExpiresAt != nil {
			view.Prize = &PrizeView{
				ID: prize.ID, Title: prize.Title, Description: prize.Description,
				BenefitType: prize.BenefitType, BenefitValue: prize.BenefitValue,
				ScopeType: prize.ScopeType, ScopeValue: prize.ScopeValue,
				Code: run.RewardCode, ExpiresAt: *run.RewardExpiresAt,
			}
		}
	}

	return view
}
