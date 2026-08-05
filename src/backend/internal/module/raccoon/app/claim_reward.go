package app

import (
	"context"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/shared/domainerr"
)

type ClaimRewardResult struct {
	RewardID  string
	Promocode string
}

type ClaimRewardUseCase struct {
	rewards RewardActivator
}

func NewClaimRewardUseCase(rewards RewardActivator) *ClaimRewardUseCase {
	return &ClaimRewardUseCase{rewards: rewards}
}

func (uc *ClaimRewardUseCase) Execute(
	ctx context.Context,
	userID uuid.UUID,
	rewardID string,
) (ClaimRewardResult, error) {
	if userID == uuid.Nil {
		return ClaimRewardResult{}, domainerr.NewInvalid("user_id", "user id is required")
	}
	if rewardID == "" {
		return ClaimRewardResult{}, domainerr.NewInvalid("reward_id", "reward id is required")
	}

	code, err := uc.rewards.Activate(ctx, userID, rewardID)
	if err != nil {
		return ClaimRewardResult{}, err
	}

	return ClaimRewardResult{RewardID: rewardID, Promocode: code}, nil
}
