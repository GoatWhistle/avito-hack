package app

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/weeklylottery/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

type StartHandler struct {
	repository Repository
	tx         TxManager
	clock      Clock
	generator  *domain.Generator
}

func NewStartHandler(repository Repository, tx TxManager, clock Clock, generator *domain.Generator) *StartHandler {
	return &StartHandler{repository: repository, tx: tx, clock: clock, generator: generator}
}

func (h *StartHandler) Handle(ctx context.Context, userID uuid.UUID) (RunView, bool, error) {
	var (
		result  RunView
		created bool
	)

	err := h.tx.WithTx(ctx, func(ctx context.Context) error {
		now := h.clock.Now()
		week := domain.WeekAt(now)
		existing, err := h.repository.ByUserWeek(ctx, userID, week.Key, true)
		if err == nil {
			result = toRunView(existing)
			return nil
		}
		if !errors.Is(err, domainerr.ErrNotFound) {
			return err
		}

		board, prize, err := h.generator.Generate()
		if err != nil {
			return err
		}

		run := domain.NewRun(userID, week, board, prize, now)
		inserted, err := h.repository.Create(ctx, run)
		if err != nil {
			return err
		}
		if !inserted {
			run, err = h.repository.ByUserWeek(ctx, userID, week.Key, true)
			if err != nil {
				return err
			}
		} else {
			created = true
		}

		result = toRunView(run)
		return nil
	})
	if err != nil {
		return RunView{}, false, err
	}

	return result, created, nil
}
