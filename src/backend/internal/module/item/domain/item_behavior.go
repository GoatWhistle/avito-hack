package domain

import (
	"time"

	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/vo"
)

func (i *Item) SubmitForModeration(now time.Time) error {
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

func (i *Item) Publish(now time.Time) error {
	if !i.status.CanTransitionTo(StatusPublished) {
		return errTransition(i.status, StatusPublished)
	}

	i.status = StatusPublished
	i.updatedAt = now

	return nil
}

func (i *Item) MarkSold(now time.Time) error {
	if !i.status.CanTransitionTo(StatusSold) {
		return errTransition(i.status, StatusSold)
	}

	i.status = StatusSold
	i.updatedAt = now

	return nil
}

func (i *Item) Archive(now time.Time) error {
	if !i.status.CanTransitionTo(StatusArchived) {
		return errTransition(i.status, StatusArchived)
	}

	i.status = StatusArchived
	i.updatedAt = now

	return nil
}

func (i *Item) Restore(now time.Time) error {
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

func (i *Item) Update(p UpdateItemParams) error {
	if i.status == StatusArchived {
		return domainerr.NewConflict("archived item cannot be modified")
	}

	if i.status == StatusSold {
		return domainerr.NewConflict("sold item cannot be modified")
	}

	if p.Title != nil {
		title, err := normalizeTitle(*p.Title)
		if err != nil {
			return err
		}
		i.title = title
	}

	if p.Description != nil {
		description, err := normalizeDescription(*p.Description)
		if err != nil {
			return err
		}
		i.description = description
	}

	if p.Price != nil {
		i.price = *p.Price
	}

	if p.Attributes != nil {
		i.attributes = p.Attributes.Clone()
	}

	i.updatedAt = p.Now

	return nil
}
