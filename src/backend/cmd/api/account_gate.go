package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/user/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

type accountGate struct {
	users domain.Repository
}

func newAccountGate(users domain.Repository) *accountGate {
	return &accountGate{users: users}
}

func (g *accountGate) IsActive(ctx context.Context, userID uuid.UUID) (bool, error) {
	_, err := g.users.ByID(ctx, userID)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, domainerr.ErrNotFound) {
		return false, nil
	}

	return false, fmt.Errorf("check account state: %w", err)
}
