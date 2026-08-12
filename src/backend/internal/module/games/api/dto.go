package api

import (
	"encoding/json"

	"github.com/avito-hack/backend/internal/module/games/app"
)

type StreakResponse struct {
	CurrentDays int  `json:"current_days"`
	BestDays    int  `json:"best_days"`
	RewardReady bool `json:"reward_ready"`
}

type DailyResponse struct {
	Attempts   int `json:"attempts"`
	BestStreak int `json:"best_streak"`
}

type GameResponse struct {
	Slug         string         `json:"slug"`
	TargetStreak int            `json:"target_streak"`
	DailyDone    bool           `json:"daily_done"`
	Streak       StreakResponse `json:"streak"`
}

type ActiveRoundResponse struct {
	RoundID string          `json:"round_id"`
	Streak  int             `json:"streak"`
	Prompt  json.RawMessage `json:"prompt"`
}

type GameStateResponse struct {
	Slug         string               `json:"slug"`
	TargetStreak int                  `json:"target_streak"`
	Streak       StreakResponse       `json:"streak"`
	Daily        DailyResponse        `json:"daily"`
	ActiveRound  *ActiveRoundResponse `json:"active_round,omitempty"`
}

type StartRoundResponse struct {
	RoundID      string          `json:"round_id"`
	Streak       int             `json:"streak"`
	TargetStreak int             `json:"target_streak"`
	Prompt       json.RawMessage `json:"prompt"`
}

type GuessRequest struct {
	Move json.RawMessage `json:"move" validate:"required"`
}

type GuessResponse struct {
	Correct          bool            `json:"correct"`
	Reveal           json.RawMessage `json:"reveal,omitempty"`
	Streak           int             `json:"streak"`
	State            string          `json:"state"`
	Prompt           json.RawMessage `json:"prompt,omitempty"`
	AttemptCompleted bool            `json:"attempt_completed"`
	StreakAfter      *StreakResponse `json:"streak_after,omitempty"`
}

type ClaimGameRewardResponse struct {
	Code string `json:"code"`
}

func toStreakResponse(v app.StreakView) StreakResponse {
	return StreakResponse{
		CurrentDays: v.CurrentDays,
		BestDays:    v.BestDays,
		RewardReady: v.RewardReady,
	}
}

func toGameResponses(views []app.GameView) []GameResponse {
	out := make([]GameResponse, 0, len(views))

	for _, view := range views {
		out = append(out, GameResponse{
			Slug:         view.Slug,
			TargetStreak: view.TargetStreak,
			DailyDone:    view.DailyDone,
			Streak:       toStreakResponse(view.Streak),
		})
	}

	return out
}

func toStateResponse(view app.StateView) GameStateResponse {
	response := GameStateResponse{
		Slug:         view.Slug,
		TargetStreak: view.TargetStreak,
		Streak:       toStreakResponse(view.Streak),
		Daily:        DailyResponse{Attempts: view.Daily.Attempts, BestStreak: view.Daily.BestStreak},
	}

	if view.ActiveRound != nil {
		response.ActiveRound = &ActiveRoundResponse{
			RoundID: view.ActiveRound.RoundID,
			Streak:  view.ActiveRound.Streak,
			Prompt:  view.ActiveRound.Prompt,
		}
	}

	return response
}

func toGuessResponse(result app.GuessResult) GuessResponse {
	response := GuessResponse{
		Correct:          result.Correct,
		Reveal:           result.Reveal,
		Streak:           result.Streak,
		State:            result.State.String(),
		Prompt:           result.Prompt,
		AttemptCompleted: result.AttemptCompleted,
	}

	if result.StreakAfter != nil {
		streak := toStreakResponse(*result.StreakAfter)
		response.StreakAfter = &streak
	}

	return response
}
