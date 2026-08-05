package api

const badgeTimeLayout = "2006-01-02T15:04:05Z"

type BadgeResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon_url"`
	EarnedAt    string `json:"earned_at,omitempty"`
}

type RaccoonProfileResponse struct {
	ID            string          `json:"id"`
	UserID        string          `json:"user_id"`
	Name          string          `json:"name"`
	Level         int             `json:"level"`
	XP            int             `json:"xp"`
	XPToNextLevel int             `json:"xp_to_next_level"`
	CurrentStreak int             `json:"current_streak"`
	Stage         string          `json:"stage"`
	State         string          `json:"state"`
	Badges        []BadgeResponse `json:"badges"`
}

type ClaimRewardRequest struct {
	RewardID string `json:"reward_id" validate:"required"`
}

type ClaimRewardResponse struct {
	RewardID  string `json:"reward_id"`
	Promocode string `json:"promocode"`
}
