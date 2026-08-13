package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/shared/publicid"
)

type State string

const (
	StateActive State = "active"
	StateWon    State = "won"
	StateLost   State = "lost"
)

type Run struct {
	ID              uuid.UUID
	DisplayID       string
	UserID          uuid.UUID
	WeekStart       time.Time
	State           State
	Board           [BoardSize]Symbol
	Opened          []int
	PrizeID         string
	RewardCode      string
	RewardExpiresAt *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewRun(userID uuid.UUID, week Week, board [BoardSize]Symbol, prize *Prize, now time.Time) Run {
	run := Run{
		ID: uuid.New(), DisplayID: publicid.New(), UserID: userID, WeekStart: week.Key,
		State: StateActive, Board: board, Opened: []int{}, CreatedAt: now, UpdatedAt: now,
	}
	if prize != nil {
		run.PrizeID = prize.ID
	}

	return run
}

func (r Run) IsOpened(slot int) bool {
	for _, opened := range r.Opened {
		if opened == slot {
			return true
		}
	}

	return false
}

func (r *Run) Reveal(slot int, now time.Time) (Symbol, error) {
	if r.State != StateActive {
		return "", ErrRunFinished()
	}
	if slot < 0 || slot >= BoardSize {
		return "", ErrInvalidSlot()
	}
	if r.IsOpened(slot) {
		return "", ErrSlotOpened()
	}

	symbol := r.Board[slot]
	r.Opened = append(r.Opened, slot)
	r.UpdatedAt = now

	if r.openCount(symbol) >= 3 {
		r.State = StateWon
	} else if len(r.Opened) == BoardSize {
		r.State = StateLost
	}

	return symbol, nil
}

func (r *Run) AttachReward(code string, expiresAt time.Time) {
	r.RewardCode = code
	r.RewardExpiresAt = &expiresAt
}

func (r Run) openCount(symbol Symbol) int {
	count := 0
	for _, slot := range r.Opened {
		if r.Board[slot] == symbol {
			count++
		}
	}

	return count
}

func ValidateRun(r Run) error {
	if err := ValidateBoard(r.Board, r.PrizeID); err != nil {
		return err
	}

	seen := make(map[int]struct{}, len(r.Opened))
	openedCounts := make(map[Symbol]int)
	for _, slot := range r.Opened {
		if slot < 0 || slot >= BoardSize {
			return fmt.Errorf("weekly lottery run contains invalid opened slot %d", slot)
		}
		if _, exists := seen[slot]; exists {
			return fmt.Errorf("weekly lottery run contains duplicate opened slot %d", slot)
		}
		seen[slot] = struct{}{}
		openedCounts[r.Board[slot]]++
	}

	hasOpenTriple := false
	for _, count := range openedCounts {
		if count >= 3 {
			hasOpenTriple = true
		}
	}

	switch r.State {
	case StateActive:
		if len(r.Opened) >= BoardSize || hasOpenTriple || r.RewardCode != "" || r.RewardExpiresAt != nil {
			return fmt.Errorf("active weekly lottery run has finished state data")
		}
	case StateWon:
		if !hasOpenTriple || r.PrizeID == "" || r.RewardCode == "" || r.RewardExpiresAt == nil {
			return fmt.Errorf("won weekly lottery run is missing result data")
		}
	case StateLost:
		if len(r.Opened) != BoardSize || hasOpenTriple || r.PrizeID != "" || r.RewardCode != "" {
			return fmt.Errorf("lost weekly lottery run has invalid result data")
		}
	default:
		return fmt.Errorf("weekly lottery run contains unknown state %q", r.State)
	}

	return nil
}
