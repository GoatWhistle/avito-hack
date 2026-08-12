package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ModerationVerdict string

const (
	ModerationApproved    ModerationVerdict = "approved"
	ModerationRejected    ModerationVerdict = "rejected"
	ModerationUnavailable ModerationVerdict = "unavailable"
)

type ModerationSubject struct {
	ItemID      uuid.UUID
	Title       string
	Description string
	PhotoURLs   []string
}

type ModerationResult struct {
	Verdict  ModerationVerdict
	Reason   string
	Provider string
}

type ModerationProvider interface {
	Review(ctx context.Context, subject ModerationSubject) ModerationResult
}

type ModerationLogEntry struct {
	ID        uuid.UUID
	ItemID    uuid.UUID
	Verdict   ModerationVerdict
	Reason    string
	Provider  string
	CheckedAt time.Time
}

func NewModerationLogEntry(itemID uuid.UUID, result ModerationResult, now time.Time) ModerationLogEntry {
	return ModerationLogEntry{
		ID:        uuid.New(),
		ItemID:    itemID,
		Verdict:   result.Verdict,
		Reason:    result.Reason,
		Provider:  result.Provider,
		CheckedAt: now,
	}
}

type ModerationLogRepository interface {
	Add(ctx context.Context, entry ModerationLogEntry) error
	LatestByItemID(ctx context.Context, itemID uuid.UUID) (*ModerationLogEntry, error)
}
