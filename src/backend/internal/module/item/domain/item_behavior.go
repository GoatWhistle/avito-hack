package domain

import (
	"time"

	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/vo"
)

func errSeedImmutable() error {
	return domainerr.NewConflict("platform-provided listing cannot be modified")
}

func (i *Item) SubmitForModeration(now time.Time) error {
	if i.isSeed {
		return errSeedImmutable()
	}
	if !i.status.CanTransitionTo(StatusModeration) {
		return errTransition(i.status, StatusModeration)
	}
	if i.description == "" {
		return domainerr.NewInvalid("description", "field is required before moderation")
	}

	i.status = StatusModeration
	i.updatedAt = now

	return nil
}

func (i *Item) Publish(now time.Time, verdict ModerationVerdict) error {
	if !i.status.CanTransitionTo(StatusPublished) {
		return errTransition(i.status, StatusPublished)
	}

	if verdict != ModerationApproved {
		return errNotApproved()
	}

	i.status = StatusPublished
	i.aiVerified = true
	i.updatedAt = now

	return nil
}

func (i *Item) MarkSold(now time.Time) error {
	if i.isSeed {
		return errSeedImmutable()
	}
	if !i.status.CanTransitionTo(StatusSold) {
		return errTransition(i.status, StatusSold)
	}

	i.status = StatusSold
	i.updatedAt = now

	return nil
}

func (i *Item) Archive(now time.Time) error {
	if i.isSeed {
		return errSeedImmutable()
	}
	if !i.status.CanTransitionTo(StatusArchived) {
		return errTransition(i.status, StatusArchived)
	}

	i.status = StatusArchived
	i.updatedAt = now

	return nil
}

func (i *Item) Restore(now time.Time) error {
	if i.isSeed {
		return errSeedImmutable()
	}
	if !i.status.CanTransitionTo(StatusDraft) {
		return errTransition(i.status, StatusDraft)
	}

	i.status = StatusDraft
	i.updatedAt = now

	return nil
}

type UpdateItemParams struct {
	Title       *string
	Description *string
	Price       *vo.Money
	Attributes  *Attributes
	Now         time.Time
}

func (i *Item) RequireReverification(now time.Time) error {
	if i.isSeed {
		return errSeedImmutable()
	}

	i.aiVerified = false
	if i.status == StatusPublished {
		i.status = StatusModeration
	}
	i.updatedAt = now

	return nil
}

func (i *Item) Update(p UpdateItemParams) error {
	if i.isSeed {
		return errSeedImmutable()
	}

	if i.status == StatusArchived {
		return domainerr.NewConflict("archived item cannot be modified")
	}

	if i.status == StatusSold {
		return domainerr.NewConflict("sold item cannot be modified")
	}

	contentChanged := false

	if p.Title != nil {
		title, err := normalizeTitle(*p.Title)
		if err != nil {
			return err
		}
		if title != i.title {
			contentChanged = true
		}
		i.title = title
	}

	if p.Description != nil {
		description, err := normalizeDescription(*p.Description)
		if err != nil {
			return err
		}
		if description != i.description {
			contentChanged = true
		}
		i.description = description
	}

	if p.Price != nil {
		i.price = *p.Price
	}

	if p.Attributes != nil {
		if !i.attributes.Equal(*p.Attributes) {
			contentChanged = true
		}
		i.attributes = p.Attributes.Clone()
	}

	if contentChanged {
		i.aiVerified = false
		if i.status == StatusPublished {
			i.status = StatusModeration
		}
	}

	i.updatedAt = p.Now

	return nil
}
