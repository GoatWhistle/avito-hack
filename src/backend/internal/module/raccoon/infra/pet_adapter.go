package infra

import (
	"context"

	"github.com/google/uuid"

	petapp "github.com/avito-hack/backend/internal/module/pet/app"
	petdomain "github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/module/raccoon/app"
)

type petStateService interface {
	State(ctx context.Context, userID uuid.UUID) (*petdomain.Pet, error)
}

type PetAdapter struct {
	pets petStateService
}

func NewPetAdapter(pets petStateService) *PetAdapter {
	return &PetAdapter{pets: pets}
}

func (a *PetAdapter) Profile(ctx context.Context, userID uuid.UUID) (app.RaccoonProfileView, error) {
	pet, err := a.pets.State(ctx, userID)
	if err != nil {
		return app.RaccoonProfileView{}, err
	}

	return app.RaccoonProfileView{
		ID:            pet.ID(),
		UserID:        pet.UserID(),
		Name:          pet.Name(),
		Level:         pet.Level(),
		XP:            pet.XP(),
		XPToNextLevel: pet.NextLevelXP(),
		CurrentStreak: pet.StreakDays(),
		Stage:         string(pet.Stage()),
		State:         string(pet.State()),
	}, nil
}

type badgeRepository interface {
	EarnedBy(ctx context.Context, userID uuid.UUID) ([]petdomain.EarnedBadge, error)
}

type BadgeAdapter struct {
	badges badgeRepository
}

func NewBadgeAdapter(badges badgeRepository) *BadgeAdapter {
	return &BadgeAdapter{badges: badges}
}

func (a *BadgeAdapter) Earned(ctx context.Context, userID uuid.UUID) ([]app.BadgeView, error) {
	earned, err := a.badges.EarnedBy(ctx, userID)
	if err != nil {
		return nil, err
	}

	views := make([]app.BadgeView, 0, len(earned))
	for _, badge := range earned {
		earnedAt := badge.EarnedAt()
		views = append(views, app.BadgeView{
			ID:          badge.ID(),
			Name:        badge.Name(),
			Description: badge.Description(),
			IconURL:     badge.IconURL(),
			EarnedAt:    &earnedAt,
		})
	}

	return views, nil
}

type RewardAdapter struct {
	rewards *petapp.RewardService
}

func NewRewardAdapter(rewards *petapp.RewardService) *RewardAdapter {
	return &RewardAdapter{rewards: rewards}
}

func (a *RewardAdapter) Activate(ctx context.Context, userID uuid.UUID, rewardID string) (string, error) {
	return a.rewards.Activate(ctx, userID, rewardID)
}
