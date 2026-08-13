package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/weeklylottery/domain"
)

const rewardLifetime = 30 * 24 * time.Hour

type RevealCommand struct {
	UserID uuid.UUID
	RunID  string
	Slot   int
}

type RevealHandler struct {
	repository Repository
	tx         TxManager
	clock      Clock
	signer     CodeSigner
}

func NewRevealHandler(repository Repository, tx TxManager, clock Clock, signer CodeSigner) *RevealHandler {
	return &RevealHandler{repository: repository, tx: tx, clock: clock, signer: signer}
}

func (h *RevealHandler) Handle(ctx context.Context, cmd RevealCommand) (RevealResult, error) {
	if cmd.Slot < 0 || cmd.Slot >= domain.BoardSize {
		return RevealResult{}, domain.ErrInvalidSlot()
	}

	var result RevealResult
	err := h.tx.WithTx(ctx, func(ctx context.Context) error {
		run, err := h.repository.ByDisplayID(ctx, cmd.RunID, true)
		if err != nil {
			return err
		}
		if run.UserID != cmd.UserID {
			return domain.ErrRunNotFound
		}

		now := h.clock.Now()
		if !run.WeekStart.Equal(domain.WeekAt(now).Key) {
			return domain.ErrRunExpired()
		}
		symbol, err := run.Reveal(cmd.Slot, now)
		if err != nil {
			return err
		}

		if run.State == domain.StateWon {
			if h.signer == nil {
				return fmt.Errorf("weekly lottery reward signer is not configured")
			}
			rewardSubject := run.PrizeID + ":" + run.DisplayID
			code, issueErr := h.signer.Issue(run.UserID, rewardSubject)
			if issueErr != nil {
				return fmt.Errorf("issue weekly lottery reward: %w", issueErr)
			}
			if verifyErr := h.signer.Verify(run.UserID, rewardSubject, code); verifyErr != nil {
				return fmt.Errorf("verify weekly lottery reward: %w", verifyErr)
			}
			run.AttachReward(code, now.Add(rewardLifetime))
		}

		if err := h.repository.Save(ctx, run); err != nil {
			return err
		}

		result = RevealResult{Index: cmd.Slot, Symbol: symbol, Run: toRunView(run)}
		return nil
	})
	if err != nil {
		return RevealResult{}, err
	}

	return result, nil
}
