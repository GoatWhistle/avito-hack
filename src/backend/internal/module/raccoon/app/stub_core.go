package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type StubCoreRaccoonService struct{}

func NewStubCoreRaccoonService() *StubCoreRaccoonService {
	return &StubCoreRaccoonService{}
}

func (s *StubCoreRaccoonService) ProcessUserAction(ctx context.Context, event ActionEvent) (*ProcessActionResult, error) {
	now := time.Now()
	return &ProcessActionResult{
		XPAdded:       50,
		CurrentXP:     350,
		CurrentLevel:  3,
		LevelUp:       false,
		CurrentStreak: 5,
		NewBadges: []BadgeView{
			{
				ID:          "badge-first-step",
				Name:        "First Step",
				Description: "Completed your first action",
				IconURL:     "/assets/badges/first-step.png",
				EarnedAt:    &now,
			},
		},
	}, nil
}

func (s *StubCoreRaccoonService) EvaluateBadges(ctx context.Context, userID uuid.UUID) ([]BadgeView, error) {
	now := time.Now()
	return []BadgeView{
		{
			ID:          "badge-first-step",
			Name:        "First Step",
			Description: "Completed your first action",
			IconURL:     "/assets/badges/first-step.png",
			EarnedAt:    &now,
		},
		{
			ID:          "badge-streak-5",
			Name:        "Streak Master",
			Description: "Maintained a 5-day streak",
			IconURL:     "/assets/badges/streak-5.png",
			EarnedAt:    &now,
		},
	}, nil
}

func (s *StubCoreRaccoonService) CheckRateLimit(ctx context.Context, userID uuid.UUID, actionType string) (bool, error) {
	return true, nil
}

type StubRaccoonStateReader struct{}

func NewStubRaccoonStateReader() *StubRaccoonStateReader {
	return &StubRaccoonStateReader{}
}

func (r *StubRaccoonStateReader) GetRaccoonProfile(ctx context.Context, userID uuid.UUID) (*RaccoonProfileView, error) {
	now := time.Now()
	return &RaccoonProfileView{
		ID:            uuid.New(),
		UserID:        userID,
		Name:          "Rocky the Raccoon",
		Level:         3,
		XP:            350,
		XPToNextLevel: 500,
		CurrentStreak: 5,
		Badges: []BadgeView{
			{
				ID:          "badge-first-step",
				Name:        "First Step",
				Description: "Completed your first action",
				IconURL:     "/assets/badges/first-step.png",
				EarnedAt:    &now,
			},
		},
	}, nil
}

type StubPromocodeGeneratorAdapter struct{}

func NewStubPromocodeGeneratorAdapter() *StubPromocodeGeneratorAdapter {
	return &StubPromocodeGeneratorAdapter{}
}

func (p *StubPromocodeGeneratorAdapter) Generate(ctx context.Context, userID uuid.UUID, rewardID string) (string, error) {
	return fmt.Sprintf("PROMO-%s-%d", rewardID, time.Now().Unix()), nil
}

type StubNotificationAdapter struct{}

func NewStubNotificationAdapter() *StubNotificationAdapter {
	return &StubNotificationAdapter{}
}

func (n *StubNotificationAdapter) SendInAppPush(ctx context.Context, userID uuid.UUID, title string, body string) error {
	return nil
}
