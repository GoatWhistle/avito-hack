package domain

import (
	"fmt"

	"github.com/avito-hack/backend/internal/shared/domainerr"
)

var ErrRunNotFound = fmt.Errorf("weekly lottery run not found: %w", domainerr.ErrNotFound)

func ErrInvalidSlot() error {
	return domainerr.NewInvalid("slot", "slot must be between 0 and 8")
}

func ErrSlotOpened() error {
	return domainerr.NewConflict("slot is already opened")
}

func ErrRunFinished() error {
	return domainerr.NewConflict("weekly lottery run is already finished")
}

func ErrRunExpired() error {
	return domainerr.NewConflict("weekly lottery run belongs to a previous week")
}
