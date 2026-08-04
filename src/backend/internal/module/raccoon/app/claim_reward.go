package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ClaimRewardResult struct {
	RewardID  string
	Promocode string
}

type ClaimRewardUseCase struct {
	promocodes PromocodeGeneratorAdapter
	notifs     NotificationAdapter
}

func NewClaimRewardUseCase(promocodes PromocodeGeneratorAdapter, notifs NotificationAdapter) *ClaimRewardUseCase {
	return &ClaimRewardUseCase{
		promocodes: promocodes,
		notifs:     notifs,
	}
}

func (uc *ClaimRewardUseCase) Execute(ctx context.Context, userID uuid.UUID, rewardID string) (*ClaimRewardResult, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("invalid user id")
	}

	if rewardID == "" {
		return nil, fmt.Errorf("invalid reward id")
	}

	promo, err := uc.promocodes.Generate(ctx, userID, rewardID)
	if err != nil {
		return nil, fmt.Errorf("generate promocode: %w", err)
	}

	_ = uc.notifs.SendInAppPush(ctx, userID, "Reward Claimed!", fmt.Sprintf("Your promo code is %s", promo))

	return &ClaimRewardResult{
		RewardID:  rewardID,
		Promocode: promo,
	}, nil
}
