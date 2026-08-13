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
	Slug         string `json:"slug"`
	TargetStreak int    `json:"target_streak"`
	MaxAttempts  int    `json:"max_attempts,omitempty"`
	DailyDone    bool   `json:"daily_done"`
}

type GameListResponse struct {
	Games     []GameResponse `json:"games"`
	Streak    StreakResponse `json:"streak"`
	DailyDone bool           `json:"daily_done"`
}

type ActiveRoundResponse struct {
	RoundID      string          `json:"round_id"`
	Streak       int             `json:"streak"`
	AttemptsUsed int             `json:"attempts_used,omitempty"`
	MaxAttempts  int             `json:"max_attempts,omitempty"`
	Prompt       json.RawMessage `json:"prompt"`
}

type GameStateResponse struct {
	Slug         string               `json:"slug"`
	TargetStreak int                  `json:"target_streak"`
	MaxAttempts  int                  `json:"max_attempts,omitempty"`
	Streak       StreakResponse       `json:"streak"`
	Daily        DailyResponse        `json:"daily"`
	ActiveRound  *ActiveRoundResponse `json:"active_round,omitempty"`
	BestScore    int                  `json:"best_score,omitempty"`
}

type StartRoundResponse struct {
	RoundID      string          `json:"round_id"`
	Streak       int             `json:"streak"`
	TargetStreak int             `json:"target_streak"`
	AttemptsUsed int             `json:"attempts_used,omitempty"`
	MaxAttempts  int             `json:"max_attempts,omitempty"`
	Prompt       json.RawMessage `json:"prompt"`
}

type GuessRequest struct {
	Move json.RawMessage `json:"move" validate:"required"`
}

type GuessResponse struct {
	Correct          bool            `json:"correct"`
	Progress         string          `json:"progress"`
	Reveal           json.RawMessage `json:"reveal,omitempty"`
	Streak           int             `json:"streak"`
	AttemptsUsed     int             `json:"attempts_used,omitempty"`
	MaxAttempts      int             `json:"max_attempts,omitempty"`
	State            string          `json:"state"`
	Prompt           json.RawMessage `json:"prompt,omitempty"`
	AttemptCompleted bool            `json:"attempt_completed"`
	StreakAfter      *StreakResponse `json:"streak_after,omitempty"`
	BestScore        int             `json:"best_score,omitempty"`
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

func toGameListResponse(view app.GameListView) GameListResponse {
	games := make([]GameResponse, 0, len(view.Games))

	for _, game := range view.Games {
		games = append(games, GameResponse{
			Slug:         game.Slug,
			TargetStreak: game.TargetStreak,
			MaxAttempts:  game.MaxAttempts,
			DailyDone:    game.DailyDone,
		})
	}

	return GameListResponse{
		Games:     games,
		Streak:    toStreakResponse(view.Streak),
		DailyDone: view.DailyDone,
	}
}

func toStateResponse(view app.StateView) GameStateResponse {
	response := GameStateResponse{
		Slug:         view.Slug,
		TargetStreak: view.TargetStreak,
		MaxAttempts:  view.MaxAttempts,
		Streak:       toStreakResponse(view.Streak),
		Daily:        DailyResponse{Attempts: view.Daily.Attempts, BestStreak: view.Daily.BestStreak},
		BestScore:    view.BestScore,
	}

	if view.ActiveRound != nil {
		response.ActiveRound = &ActiveRoundResponse{
			RoundID:      view.ActiveRound.RoundID,
			Streak:       view.ActiveRound.Streak,
			AttemptsUsed: view.ActiveRound.AttemptsUsed,
			MaxAttempts:  view.ActiveRound.MaxAttempts,
			Prompt:       view.ActiveRound.Prompt,
		}
	}

	return response
}

func toGuessResponse(result app.GuessResult) GuessResponse {
	response := GuessResponse{
		Correct:          result.Correct,
		Progress:         string(result.Progress),
		Reveal:           result.Reveal,
		Streak:           result.Streak,
		AttemptsUsed:     result.AttemptsUsed,
		MaxAttempts:      result.MaxAttempts,
		State:            result.State.String(),
		Prompt:           result.Prompt,
		AttemptCompleted: result.AttemptCompleted,
		BestScore:        result.BestScore,
	}

	if result.StreakAfter != nil {
		streak := toStreakResponse(*result.StreakAfter)
		response.StreakAfter = &streak
	}

	return response
}
