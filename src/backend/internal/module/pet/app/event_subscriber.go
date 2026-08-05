package app

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/events"
)

type PetNotifier interface {
	PetUpdated(userID uuid.UUID, pet *domain.Pet)
	XPGained(userID uuid.UUID, amount int, reason string, total int)
	LevelUp(userID uuid.UUID, level int)
}

type reaction struct {
	action domain.Action
	award  func(*domain.Pet, domain.LimitedAction) (domain.Progress, error)
}

type Subscriber struct {
	service  *Service
	notifier PetNotifier
	clock    Clock
	rewards  *RewardService
}

func NewSubscriber(service *Service, notifier PetNotifier, clock Clock) *Subscriber {
	return &Subscriber{service: service, notifier: notifier, clock: clock}
}

func (s *Subscriber) WithRewards(rewards *RewardService) *Subscriber {
	s.rewards = rewards

	return s
}

func (s *Subscriber) Register(bus events.Subscriber) {
	if bus == nil {
		return
	}

	bus.Subscribe(events.TypeItemPublished, s.onItemPublished)
	bus.Subscribe(events.TypeItemSold, s.onItemSold)
	bus.Subscribe(events.TypeFavoriteAdded, s.onFavoriteAdded)
	bus.Subscribe(events.TypeUserRegistered, s.onUserRegistered)
}

func (s *Subscriber) onItemPublished(ctx context.Context, e events.Event) error {
	return s.award(ctx, e, reaction{
		action: domain.ActionItemPublished,
		award:  (*domain.Pet).RewardItemPublished,
	})
}

func (s *Subscriber) onItemSold(ctx context.Context, e events.Event) error {
	return s.award(ctx, e, reaction{
		action: domain.ActionItemSold,
		award:  (*domain.Pet).RewardItemSold,
	})
}

func (s *Subscriber) onFavoriteAdded(ctx context.Context, e events.Event) error {
	return s.award(ctx, e, reaction{
		action: domain.ActionFavorite,
		award: func(p *domain.Pet, a domain.LimitedAction) (domain.Progress, error) {
			return p.RewardFavorite(a)
		},
	})
}

func (s *Subscriber) onUserRegistered(ctx context.Context, e events.Event) error {
	pet, err := s.service.Create(ctx, e.UserID)
	if err != nil {
		return err
	}

	s.notify(e.UserID, pet, domain.Progress{}, string(events.TypeUserRegistered))

	return nil
}

func (s *Subscriber) award(ctx context.Context, e events.Event, r reaction) error {
	var progress domain.Progress

	pet, err := s.service.AwardAction(ctx, AwardCommand{
		UserID:    e.UserID,
		Action:    r.action,
		SubjectID: e.SubjectID,
		Apply: func(p *domain.Pet, a domain.LimitedAction) error {
			result, awardErr := r.award(p, a)
			if awardErr != nil {
				return awardErr
			}
			progress = result

			return nil
		},
	})
	if err != nil {
		if isSkippable(err) {
			return nil
		}

		return err
	}

	s.grantRewards(ctx, e.UserID, progress)
	s.notify(e.UserID, pet, progress, string(e.Type))

	return nil
}

func (s *Subscriber) notify(userID uuid.UUID, pet *domain.Pet, progress domain.Progress, reason string) {
	if s.notifier == nil || pet == nil {
		return
	}

	s.notifier.PetUpdated(userID, pet)

	if progress.XPGranted > 0 {
		s.notifier.XPGained(userID, progress.XPGranted, reason, pet.XP())
	}

	if progress.Level > progress.PreviousLevel {
		s.notifier.LevelUp(userID, progress.Level)
	}
}

func (s *Subscriber) grantRewards(ctx context.Context, userID uuid.UUID, progress domain.Progress) {
	if s.rewards == nil || progress.Level <= progress.PreviousLevel {
		return
	}

	if _, err := s.rewards.GrantEligible(ctx, userID); err != nil {
		slog.ErrorContext(ctx, "grant level rewards failed",
			slog.String("user_id", userID.String()),
			slog.Any("error", err),
		)
	}
}

func isSkippable(err error) bool {
	return errors.Is(err, domain.ErrDuplicateAction) ||
		errors.Is(err, domain.ErrLimitReached) ||
		errors.Is(err, domain.ErrConditionNotMet) ||
		errors.Is(err, domain.ErrInvalidAction)
}
