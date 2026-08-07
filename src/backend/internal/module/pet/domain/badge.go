package domain

import (
	"time"

	"github.com/google/uuid"
)

type Badge struct {
	id          string
	name        string
	description string
	iconURL     string
}

func NewBadge(id, name, description, iconURL string) Badge {
	return Badge{id: id, name: name, description: description, iconURL: iconURL}
}

func (b Badge) ID() string          { return b.id }
func (b Badge) Name() string        { return b.name }
func (b Badge) Description() string { return b.description }
func (b Badge) IconURL() string     { return b.iconURL }

type EarnedBadge struct {
	badge    Badge
	userID   uuid.UUID
	earnedAt time.Time
}

func NewEarnedBadge(badge Badge, userID uuid.UUID, earnedAt time.Time) EarnedBadge {
	return EarnedBadge{badge: badge, userID: userID, earnedAt: earnedAt}
}

func (e EarnedBadge) Badge() Badge        { return e.badge }
func (e EarnedBadge) UserID() uuid.UUID   { return e.userID }
func (e EarnedBadge) EarnedAt() time.Time { return e.earnedAt }
func (e EarnedBadge) ID() string          { return e.badge.id }
func (e EarnedBadge) Name() string        { return e.badge.name }
func (e EarnedBadge) Description() string { return e.badge.description }
func (e EarnedBadge) IconURL() string     { return e.badge.iconURL }

const badgeRaccoonFriend = "raccoon_friend"

const raccoonFriendLevel = 10

func BadgesEarnedBy(pet *Pet) []string {
	if pet == nil {
		return nil
	}

	earned := make([]string, 0, 1)
	if pet.Level() >= raccoonFriendLevel {
		earned = append(earned, badgeRaccoonFriend)
	}

	return earned
}
