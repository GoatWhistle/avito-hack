package app

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/weeklylottery/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

type GetStateHandler struct {
	repository Repository
	clock      Clock
}

func NewGetStateHandler(repository Repository, clock Clock) *GetStateHandler {
	return &GetStateHandler{repository: repository, clock: clock}
}

func (h *GetStateHandler) Handle(ctx context.Context, userID uuid.UUID) (StateView, error) {
	week := domain.WeekAt(h.clock.Now())
	view := StateView{Available: true, WeekStart: week.Start, NextAvailable: week.End}

	run, err := h.repository.ByUserWeek(ctx, userID, week.Key, false)
	if errors.Is(err, domainerr.ErrNotFound) {
		return view, nil
	}
	if err != nil {
		return StateView{}, err
	}

	runView := toRunView(run)
	view.Available = run.State == domain.StateActive
	view.Run = &runView

	return view, nil
}
