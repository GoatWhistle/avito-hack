package api

import (
	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/ws"
)

const (
	eventPetUpdated    = "pet.updated"
	eventPetHatched    = "pet.hatched"
	eventXPGained      = "xp.gained"
	eventLevelUp       = "level.up"
	eventRewardGranted = "reward.granted"
	eventStreakUpdated = "streak.updated"
)

type broadcaster interface {
	Broadcast(userID uuid.UUID, m ws.Message)
}

type xpGainedPayload struct {
	Amount int    `json:"amount"`
	Reason string `json:"reason"`
	Total  int    `json:"total"`
}

type levelUpPayload struct {
	Level int `json:"level"`
}

type rewardGrantedPayload struct {
	RewardID string `json:"reward_id"`
	Title    string `json:"title"`
}

type streakUpdatedPayload struct {
	Days      int  `json:"days"`
	Milestone bool `json:"milestone"`
}

type Notifier struct {
	hub broadcaster
}

func NewNotifier(hub broadcaster) *Notifier {
	return &Notifier{hub: hub}
}

func (n *Notifier) PetUpdated(userID uuid.UUID, pet *domain.Pet) {
	if n.hub == nil || pet == nil {
		return
	}

	n.hub.Broadcast(userID, ws.Message{Type: eventPetUpdated, Payload: toPetPayload(pet)})
}

func (n *Notifier) PetHatched(userID uuid.UUID, pet *domain.Pet) {
	if n.hub == nil || pet == nil {
		return
	}

	n.hub.Broadcast(userID, ws.Message{Type: eventPetHatched, Payload: toPetPayload(pet)})
}

func (n *Notifier) XPGained(userID uuid.UUID, amount int, reason string, total int) {
	if n.hub == nil {
		return
	}

	n.hub.Broadcast(userID, ws.Message{
		Type:    eventXPGained,
		Payload: xpGainedPayload{Amount: amount, Reason: reason, Total: total},
	})
}

func (n *Notifier) LevelUp(userID uuid.UUID, level int) {
	if n.hub == nil {
		return
	}

	n.hub.Broadcast(userID, ws.Message{Type: eventLevelUp, Payload: levelUpPayload{Level: level}})
}

func (n *Notifier) RewardGranted(userID uuid.UUID, rewardID, title string) {
	if n.hub == nil {
		return
	}

	n.hub.Broadcast(userID, ws.Message{
		Type:    eventRewardGranted,
		Payload: rewardGrantedPayload{RewardID: rewardID, Title: title},
	})
}

func (n *Notifier) StreakUpdated(userID uuid.UUID, days int, milestone bool) {
	if n.hub == nil {
		return
	}

	n.hub.Broadcast(userID, ws.Message{
		Type:    eventStreakUpdated,
		Payload: streakUpdatedPayload{Days: days, Milestone: milestone},
	})
}
