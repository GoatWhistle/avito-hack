package app

import (
	"context"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/shared/domainerr"
)

type GetRaccoonProfileUseCase struct {
	pets   PetStateReader
	badges BadgeReader
}

func NewGetRaccoonProfileUseCase(pets PetStateReader, badges BadgeReader) *GetRaccoonProfileUseCase {
	return &GetRaccoonProfileUseCase{pets: pets, badges: badges}
}

func (uc *GetRaccoonProfileUseCase) Execute(ctx context.Context, userID uuid.UUID) (RaccoonProfileView, error) {
	if userID == uuid.Nil {
		return RaccoonProfileView{}, domainerr.NewInvalid("user_id", "user id is required")
	}

	profile, err := uc.pets.Profile(ctx, userID)
	if err != nil {
		return RaccoonProfileView{}, err
	}

	badges, err := uc.ListBadges(ctx, userID)
	if err != nil {
		return RaccoonProfileView{}, err
	}
	profile.Badges = badges

	return profile, nil
}

func (uc *GetRaccoonProfileUseCase) ListBadges(ctx context.Context, userID uuid.UUID) ([]BadgeView, error) {
	if userID == uuid.Nil {
		return nil, domainerr.NewInvalid("user_id", "user id is required")
	}
	if uc.badges == nil {
		return []BadgeView{}, nil
	}

	return uc.badges.Earned(ctx, userID)
}
