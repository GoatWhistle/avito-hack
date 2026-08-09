package infra

import (
	"context"
	"time"

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

type badgeProgressReader interface {
	Progress(ctx context.Context, userID uuid.UUID) ([]petapp.BadgeStatus, error)
}

type BadgeAdapter struct {
	badges   badgeRepository
	progress badgeProgressReader
}

func NewBadgeAdapter(badges badgeRepository) *BadgeAdapter {
	return &BadgeAdapter{badges: badges}
}

func (a *BadgeAdapter) WithProgress(progress badgeProgressReader) *BadgeAdapter {
	a.progress = progress

	return a
}

func (a *BadgeAdapter) Earned(ctx context.Context, userID uuid.UUID) ([]app.BadgeView, error) {
	earned, err := a.badges.EarnedBy(ctx, userID)
	if err != nil {
		return nil, err
	}

	earnedAt := make(map[string]time.Time, len(earned))
	for _, badge := range earned {
		earnedAt[badge.ID()] = badge.EarnedAt()
	}

	if a.progress == nil {
		return earnedViews(earned), nil
	}

	statuses, err := a.progress.Progress(ctx, userID)
	if err != nil {
		return nil, err
	}

	views := make([]app.BadgeView, 0, len(statuses))
	for _, status := range statuses {
		view := app.BadgeView{
			ID:          status.Badge.ID(),
			Name:        status.Badge.Name(),
			Description: status.Badge.Description(),
			IconURL:     status.Badge.IconURL(),
			Current:     status.Current,
			Target:      status.Target,
		}
		if at, ok := earnedAt[status.Badge.ID()]; ok {
			view.EarnedAt = &at
		}
		views = append(views, view)
	}

	return views, nil
}

func earnedViews(earned []petdomain.EarnedBadge) []app.BadgeView {
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

	return views
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
