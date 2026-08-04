package api

import (
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/pet/domain"
)

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
	ID              uuid.UUID    `json:"id"`
	UserID          uuid.UUID    `json:"user_id"`
	Name            string       `json:"name"`
	Stage           domain.Stage `json:"stage"`
	Level           int          `json:"level"`
	XP              int          `json:"xp"`
	NextLevelXP     int          `json:"next_level_xp"`
	Satiety         int          `json:"satiety"`
	Happiness       int          `json:"happiness"`
	StreakDays      int          `json:"streak_days"`
	LastCheckInDate *time.Time   `json:"last_checkin_date,omitempty"`
	LastDecayTime   time.Time    `json:"last_decay_time"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

func toPetPayload(p *domain.Pet) petPayload {
	return petPayload{
		ID: p.ID(), UserID: p.UserID(), Name: p.Name(), Stage: p.Stage(), Level: p.Level(), XP: p.XP(),
		NextLevelXP: p.NextLevelXP(), Satiety: p.Satiety(), Happiness: p.Happiness(),
		StreakDays: p.StreakDays(), LastCheckInDate: p.LastCheckInDate(),
		LastDecayTime: p.LastDecayTime(), UpdatedAt: p.UpdatedAt(),
	}
}
