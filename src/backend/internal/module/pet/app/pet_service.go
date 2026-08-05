package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

type HatchNotifier interface {
	PetHatched(userID uuid.UUID, pet *domain.Pet)
}

type Service struct {
	pets     domain.Repository
	journal  XPJournal
	cache    Cache
	tx       TxManager
	clock    Clock
	notifier HatchNotifier
}

func NewService(pets domain.Repository, journal XPJournal, cache Cache, tx TxManager, clock Clock) *Service {
	return &Service{pets: pets, journal: journal, cache: cache, tx: tx, clock: clock}
}

func (s *Service) WithHatchNotifier(notifier HatchNotifier) *Service {
	s.notifier = notifier

	return s
}

func (s *Service) State(ctx context.Context, userID uuid.UUID) (*domain.Pet, error) {
	cached, err := s.cache.Get(ctx, userID)
	if err == nil && cached != nil {
		cached.ApplyDecay(s.clock.Now())

		return cached, nil
	}

	pet, err := s.pets.ByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	pet.ApplyDecay(s.clock.Now())

	return pet, nil
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID) (*domain.Pet, error) {
	return s.mutate(ctx, userID, func(_ context.Context, pet *domain.Pet) error {
		pet.ApplyDecay(s.clock.Now())

		return nil
	})
}

func (s *Service) Stroke(ctx context.Context, userID uuid.UUID) (*domain.Pet, error) {
	return s.hatching(ctx, userID, func(_ context.Context, pet *domain.Pet) error {
		pet.ApplyDecay(s.clock.Now())
		pet.Stroke(s.clock.Now())

		return nil
	})
}

func (s *Service) hatching(
	ctx context.Context,
	userID uuid.UUID,
	apply func(context.Context, *domain.Pet) error,
) (*domain.Pet, error) {
	hatched := false

	pet, err := s.mutate(ctx, userID, func(ctx context.Context, pet *domain.Pet) error {
		if applyErr := apply(ctx, pet); applyErr != nil {
			return applyErr
		}
		hatched = pet.Hatch(s.clock.Now())

		return nil
	})
	if err != nil {
		return nil, err
	}

	if hatched && s.notifier != nil {
		s.notifier.PetHatched(userID, pet)
	}

	return pet, nil
}

func (s *Service) mutate(
	ctx context.Context,
	userID uuid.UUID,
	apply func(context.Context, *domain.Pet) error,
) (*domain.Pet, error) {
	var result *domain.Pet

	err := s.tx.WithTx(ctx, func(ctx context.Context) error {
		pet, err := s.load(ctx, userID)
		if err != nil {
			return err
		}
		if err := apply(ctx, pet); err != nil {
			return err
		}
		if err := s.pets.Save(ctx, pet); err != nil {
			return err
		}
		result = pet

		return nil
	})
	if err != nil {
		if invalidateErr := s.cache.Delete(ctx, userID); invalidateErr != nil {
			return nil, errors.Join(err, invalidateErr)
		}

		return nil, err
	}

	if err := s.cache.Set(ctx, result); err != nil {
		return nil, fmt.Errorf("cache pet: %w", err)
	}

	return result, nil
}

func (s *Service) load(ctx context.Context, userID uuid.UUID) (*domain.Pet, error) {
	pet, err := s.pets.ByUserIDForUpdate(ctx, userID)
	if err != nil && !errors.Is(err, domainerr.ErrNotFound) {
		return nil, fmt.Errorf("load pet: %w", err)
	}
	if errors.Is(err, domainerr.ErrNotFound) {
		return domain.New(userID, s.clock.Now()), nil
	}

	return pet, nil
}
