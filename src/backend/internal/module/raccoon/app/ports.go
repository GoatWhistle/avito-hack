package app

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type BadgeView struct {
	ID          string
	Name        string
	Description string
	IconURL     string
	EarnedAt    *time.Time
}

type RaccoonProfileView struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	Name          string
	Level         int
	XP            int
	XPToNextLevel int
	CurrentStreak int
	Stage         string
	State         string
	Badges        []BadgeView
}

type PetStateReader interface {
	Profile(ctx context.Context, userID uuid.UUID) (RaccoonProfileView, error)
}

type BadgeReader interface {
	Earned(ctx context.Context, userID uuid.UUID) ([]BadgeView, error)
}

type RewardActivator interface {
	Activate(ctx context.Context, userID uuid.UUID, rewardID string) (string, error)
}
