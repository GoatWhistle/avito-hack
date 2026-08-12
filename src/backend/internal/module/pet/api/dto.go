package api

import (
	"time"

	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/module/pet/domain"
)

type RewardCatalogResponse struct {
	Items      []rewardCatalogItem `json:"items"`
	NextCursor string              `json:"next_cursor,omitempty"`
}

type UserRewardListResponse struct {
	Items      []myRewardItem `json:"items"`
	NextCursor string         `json:"next_cursor,omitempty"`
}

type QuestListResponse struct {
	Items      []questItem `json:"items"`
	NextCursor string      `json:"next_cursor,omitempty"`
}

type questItem struct {
	ID        string `json:"id"`
	Action    string `json:"action"`
	Target    int    `json:"target"`
	Reward    int    `json:"reward_xp"`
	Current   int    `json:"progress_current"`
	Completed bool   `json:"completed"`
	Claimed   bool   `json:"claimed"`
}

type rewardCatalogItem struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Kind           string `json:"kind"`
	ConditionType  string `json:"condition_type"`
	ConditionValue int    `json:"condition_value"`
	Unlocked       bool   `json:"unlocked"`
	Claimed        bool   `json:"claimed"`
	Status         string `json:"status,omitempty"`
	Current        int    `json:"progress_current"`
	Target         int    `json:"progress_target"`
}

type myRewardItem struct {
	RewardID    string     `json:"reward_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Kind        string     `json:"kind"`
	Status      string     `json:"status"`
	Code        string     `json:"code,omitempty"`
	GrantedAt   time.Time  `json:"granted_at"`
	ActivatedAt *time.Time `json:"activated_at,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type activateRewardResponse struct {
	RewardID string `json:"reward_id"`
	Code     string `json:"code"`
	Status   string `json:"status"`
}

type clientMessage struct {
	Type      string `json:"type"`
	RequestID string `json:"request_id,omitempty"`
}

type serverMessage struct {
	Type      string `json:"type"`
	RequestID string `json:"request_id,omitempty"`
	Payload   any    `json:"payload,omitempty"`
}

type errorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type petPayload struct {
	ID              string       `json:"id"`
	UserID          string       `json:"user_id"`
	Name            string       `json:"name"`
	Stage           domain.Stage `json:"stage"`
	State           domain.State `json:"state"`
	Level           int          `json:"level"`
	XP              int          `json:"xp"`
	NextLevelXP     int          `json:"next_level_xp"`
	Satiety         int          `json:"satiety"`
	Happiness       int          `json:"happiness"`
	Energy          int          `json:"energy"`
	StreakDays      int          `json:"streak_days"`
	Freezes         int          `json:"freezes"`
	IsHatched       bool         `json:"is_hatched"`
	HatchedAt       *time.Time   `json:"hatched_at,omitempty"`
	LastCheckInDate *time.Time   `json:"last_checkin_date,omitempty"`
	LastDecayTime   time.Time    `json:"last_decay_time"`
	UpdatedAt       time.Time    `json:"updated_at"`
	FeedAvailableAt *time.Time   `json:"feed_available_at"`
	CheckInApplied  bool         `json:"checkin_applied"`
	CheckIn         *checkInInfo `json:"checkin,omitempty"`
}

type checkInInfo struct {
	XPGranted       int           `json:"xp_granted"`
	Level           int           `json:"level"`
	PreviousLevel   int           `json:"previous_level"`
	NextLevelXP     int           `json:"next_level_xp"`
	UnlockedRewards []string      `json:"unlocked_rewards"`
	Streak          streakPayload `json:"streak"`
}

func toPetPayload(p *domain.Pet) petPayload {
	return petPayload{
		ID:     domain.PublicToken(p.ID()),
		UserID: domain.PublicToken(p.UserID()),
		Name:   p.Name(), Stage: p.Stage(), State: p.State(),
		Level: p.Level(), XP: p.XP(), NextLevelXP: p.NextLevelXP(),
		Satiety: p.Satiety(), Happiness: p.Happiness(), Energy: p.Energy(),
		StreakDays: p.StreakDays(), Freezes: p.Freezes(),
		IsHatched: p.IsHatched(), HatchedAt: p.HatchedAt(), LastCheckInDate: p.LastCheckInDate(),
		LastDecayTime: p.LastDecayTime(), UpdatedAt: p.UpdatedAt(),
	}
}

func toStatePayload(view app.StateView) petPayload {
	payload := toPetPayload(view.Pet)
	payload.FeedAvailableAt = view.FeedAvailableAt
	payload.CheckInApplied = view.CheckInApplied

	if view.CheckInApplied {
		payload.CheckIn = &checkInInfo{
			XPGranted:       view.Progress.XPGranted,
			Level:           view.Progress.Level,
			PreviousLevel:   view.Progress.PreviousLevel,
			NextLevelXP:     view.Progress.NextLevelXP,
			UnlockedRewards: unlockedOrEmpty(view.Progress.UnlockedRewards),
			Streak:          toStreakPayload(view.Streak),
		}
	}

	return payload
}

type progressPayload struct {
	Level       int          `json:"level"`
	XP          int          `json:"xp"`
	NextLevelXP int          `json:"next_level_xp"`
	XPToNext    int          `json:"xp_to_next_level"`
	IsMaxLevel  bool         `json:"is_max_level"`
	Stage       domain.Stage `json:"stage"`
	StreakDays  int          `json:"streak_days"`
	Freezes     int          `json:"freezes"`
}

func toProgressPayload(p app.ProgressView) progressPayload {
	return progressPayload{
		Level: p.Level, XP: p.XP, NextLevelXP: p.NextLevelXP, XPToNext: p.XPToNext,
		IsMaxLevel: p.IsMaxLevel, Stage: p.Stage, StreakDays: p.StreakDays, Freezes: p.Freezes,
	}
}

type streakPayload struct {
	Days             int  `json:"days"`
	Continued        bool `json:"continued"`
	FreezeUsed       bool `json:"freeze_used"`
	Reset            bool `json:"reset"`
	MilestoneBonus   int  `json:"milestone_bonus"`
	MilestoneReached int  `json:"milestone_reached"`
	FreezesLeft      int  `json:"freezes_left"`
}

type checkInPayload struct {
	Pet             petPayload    `json:"pet"`
	XPGranted       int           `json:"xp_granted"`
	Level           int           `json:"level"`
	PreviousLevel   int           `json:"previous_level"`
	NextLevelXP     int           `json:"next_level_xp"`
	UnlockedRewards []string      `json:"unlocked_rewards"`
	Streak          streakPayload `json:"streak"`
}

func toCheckInPayload(result app.ActionResult) checkInPayload {
	return checkInPayload{
		Pet:             toPetPayload(result.Pet),
		XPGranted:       result.Progress.XPGranted,
		Level:           result.Progress.Level,
		PreviousLevel:   result.Progress.PreviousLevel,
		NextLevelXP:     result.Progress.NextLevelXP,
		UnlockedRewards: unlockedOrEmpty(result.Progress.UnlockedRewards),
		Streak:          toStreakPayload(result.Streak),
	}
}

func toStreakPayload(outcome domain.StreakOutcome) streakPayload {
	return streakPayload{
		Days: outcome.Days, Continued: outcome.Continued,
		FreezeUsed: outcome.FreezeUsed, Reset: outcome.Reset,
		MilestoneBonus:   outcome.MilestoneBonus,
		MilestoneReached: outcome.MilestoneReached,
		FreezesLeft:      outcome.FreezesLeft,
	}
}

func unlockedOrEmpty(rewards []string) []string {
	if rewards == nil {
		return []string{}
	}

	return rewards
}
