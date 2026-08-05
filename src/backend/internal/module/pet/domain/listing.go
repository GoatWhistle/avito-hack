package domain

import (
	"github.com/google/uuid"
)

const (
	ActionItemPublished Action = "item_published"
	ActionItemSold      Action = "item_sold"
)

const (
	itemPublishedXP    = 50
	itemSoldXP         = 100
	maxPublishesPerDay = 3
	publishSatiety     = 20
	soldMood           = 40
)

func (p *Pet) RewardItemPublished(action LimitedAction) (Progress, error) {
	progress, err := p.rewardLimited(action, maxPublishesPerDay, itemPublishedXP)
	if err != nil {
		return Progress{}, err
	}

	p.Feed(publishSatiety, action.Now)

	return progress, nil
}

func (p *Pet) RewardItemSold(action LimitedAction) (Progress, error) {
	if action.SubjectID == uuid.Nil {
		return Progress{}, ErrInvalidAction
	}

	if !action.Unique {
		return Progress{}, ErrDuplicateAction
	}

	progress := p.applyAward(itemSoldXP, action.Now)
	p.Cheer(soldMood, action.Now)

	return progress, nil
}
