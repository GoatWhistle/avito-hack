package api

import (
	"time"

	"github.com/avito-hack/backend/internal/module/pet/domain"
)

type advicePayload struct {
	Text      string  `json:"text"`
	ItemID    *string `json:"item_id,omitempty"`
	ItemTitle string  `json:"item_title,omitempty"`
	Action    string  `json:"action"`
}

type summaryActionPayload struct {
	Action string `json:"action"`
	Count  int    `json:"count"`
	Amount int    `json:"amount"`
}

type summaryFactsPayload struct {
	TotalXP     int                    `json:"total_xp"`
	Actions     []summaryActionPayload `json:"actions"`
	Level       int                    `json:"level"`
	PrevLevel   int                    `json:"previous_level"`
	LeveledUp   bool                   `json:"leveled_up"`
	XP          int                    `json:"xp"`
	NextLevelXP int                    `json:"next_level_xp"`
	XPToNext    int                    `json:"xp_to_next_level"`
	Rewards     []string               `json:"rewards"`
	Badges      []string               `json:"badges"`
	Stage       domain.Stage           `json:"stage"`
	State       domain.State           `json:"state"`
	Satiety     int                    `json:"satiety"`
	Happiness   int                    `json:"happiness"`
	Energy      int                    `json:"energy"`
	StreakDays  int                    `json:"streak_days"`
	StreakBroke bool                   `json:"streak_broken"`
	Rank        *int                   `json:"leaderboard_rank,omitempty"`
	IssuesCount int                    `json:"issues_count"`
}

type summaryPayload struct {
	ID          string               `json:"id"`
	Date        string               `json:"date"`
	Message     string               `json:"message"`
	Advice      *advicePayload       `json:"advice,omitempty"`
	GeneratedBy domain.SummarySource `json:"generated_by"`
	Facts       summaryFactsPayload  `json:"facts"`
	CreatedAt   time.Time            `json:"created_at"`
}

func toSummaryPayload(summary *domain.DailySummary) summaryPayload {
	facts := summary.Facts()

	return summaryPayload{
		ID:          domain.PublicToken(summary.ID()),
		Date:        summary.Date().Format(time.DateOnly),
		Message:     summary.Message(),
		Advice:      toAdvicePayload(summary.Advice()),
		GeneratedBy: summary.GeneratedBy(),
		Facts:       toFactsPayload(facts),
		CreatedAt:   summary.CreatedAt(),
	}
}

func toAdvicePayload(advice *domain.Advice) *advicePayload {
	if advice == nil {
		return nil
	}

	return &advicePayload{
		Text: advice.Text, ItemID: advice.ItemID, ItemTitle: advice.ItemTitle,
		Action: string(advice.Action),
	}
}

func toFactsPayload(facts domain.DayFacts) summaryFactsPayload {
	actions := make([]summaryActionPayload, 0, len(facts.Actions))
	for _, action := range facts.Actions {
		actions = append(actions, summaryActionPayload{
			Action: string(action.Action), Count: action.Count, Amount: action.Amount,
		})
	}

	payload := summaryFactsPayload{
		TotalXP: facts.TotalXP, Actions: actions,
		Level: facts.Level.Current, PrevLevel: facts.Level.Previous, LeveledUp: facts.LeveledUp,
		XP: facts.Level.XP, NextLevelXP: facts.Level.NextLevelXP, XPToNext: facts.Level.XPToNext,
		Rewards: orEmpty(facts.Rewards), Badges: orEmpty(facts.Badges),
		Stage: facts.Pet.Stage, State: facts.Pet.State, Satiety: facts.Pet.Satiety,
		Happiness: facts.Pet.Happiness, Energy: facts.Pet.Energy,
		StreakDays: facts.Streak.Days, StreakBroke: facts.Streak.Broken,
		IssuesCount: len(facts.Issues),
	}

	if facts.Leaderboard.Known {
		rank := facts.Leaderboard.Rank
		payload.Rank = &rank
	}

	return payload
}

func orEmpty(values []string) []string {
	if values == nil {
		return []string{}
	}

	return values
}
