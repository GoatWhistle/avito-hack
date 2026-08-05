package app

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(context.Context) error) error
}

type Clock interface {
	Now() time.Time
}

type OwnerView struct {
	ID          uuid.UUID
	DisplayName string
}

type OwnerProvider interface {
	ByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]OwnerView, error)
}

type ActionEvent struct {
	UserID     uuid.UUID
	ActionType string
	Payload    map[string]any
	Timestamp  time.Time
}

type ProcessActionResult struct {
	XPAdded       int
	CurrentXP     int
	CurrentLevel  int
	LevelUp       bool
	CurrentStreak int
	NewBadges     []BadgeView
}

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
	Badges        []BadgeView
}

type CoreRaccoonService interface {
	ProcessUserAction(ctx context.Context, event ActionEvent) (*ProcessActionResult, error)
	EvaluateBadges(ctx context.Context, userID uuid.UUID) ([]BadgeView, error)
	CheckRateLimit(ctx context.Context, userID uuid.UUID, actionType string) (bool, error)
}

type RaccoonStateReader interface {
	GetRaccoonProfile(ctx context.Context, userID uuid.UUID) (*RaccoonProfileView, error)
}

type PromocodeGeneratorAdapter interface {
	Generate(ctx context.Context, userID uuid.UUID, rewardID string) (string, error)
}

type NotificationAdapter interface {
	SendInAppPush(ctx context.Context, userID uuid.UUID, title string, body string) error
}


