package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

type Service struct {
	pets  domain.Repository
	cache Cache
	tx    TxManager
	clock Clock
}

func NewService(pets domain.Repository, cache Cache, tx TxManager, clock Clock) *Service {
	return &Service{pets: pets, cache: cache, tx: tx, clock: clock}
}

func (s *Service) State(ctx context.Context, userID uuid.UUID) (*domain.Pet, error) {
	if cached, err := s.cache.Get(ctx, userID); err == nil && cached != nil {
		return cached, nil
	}

	pet, err := s.loadOrCreate(ctx, userID, false)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Set(ctx, pet)

	return pet, nil
}

func (s *Service) Pet(ctx context.Context, userID uuid.UUID) (*domain.Pet, error) {
	pet, err := s.loadOrCreate(ctx, userID, true)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Set(ctx, pet)

	return pet, nil
}

func (s *Service) loadOrCreate(ctx context.Context, userID uuid.UUID, applyPet bool) (*domain.Pet, error) {
	var result *domain.Pet
	err := s.tx.WithTx(ctx, func(ctx context.Context) error {
		pet, err := s.pets.ByUserIDForUpdate(ctx, userID)
		if err != nil && !errors.Is(err, domainerr.ErrNotFound) {
			return fmt.Errorf("load pet: %w", err)
		}
		if errors.Is(err, domainerr.ErrNotFound) {
			pet = domain.New(userID, s.clock.Now())
		}
		if applyPet {
			pet.Pet(s.clock.Now())
		}
		if err := s.pets.Save(ctx, pet); err != nil {
			return err
		}
		result = pet
		return nil
	})

	return result, err
}
