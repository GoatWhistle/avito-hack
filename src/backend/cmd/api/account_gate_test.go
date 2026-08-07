package main

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/user/domain"
	"github.com/avito-hack/backend/internal/shared/vo"
)

type stubUserRepository struct {
	user *domain.User
	err  error
}

func (s *stubUserRepository) Save(context.Context, *domain.User) error { return nil }

func (s *stubUserRepository) ByID(context.Context, uuid.UUID) (*domain.User, error) {
	return s.user, s.err
}

func (s *stubUserRepository) ByEmail(context.Context, vo.Email) (*domain.User, error) {
	return s.user, s.err
}

func (s *stubUserRepository) ExistsByEmail(context.Context, vo.Email) (bool, error) {
	return s.user != nil, s.err
}

func TestAccountGateReportsActiveUser(t *testing.T) {
	t.Parallel()

	gate := newAccountGate(&stubUserRepository{user: &domain.User{}})

	active, err := gate.IsActive(context.Background(), uuid.New())

	require.NoError(t, err)
	assert.True(t, active)
}

func TestAccountGateReportsSoftDeletedUserAsInactive(t *testing.T) {
	t.Parallel()

	gate := newAccountGate(&stubUserRepository{err: domain.ErrUserNotFound})

	active, err := gate.IsActive(context.Background(), uuid.New())

	require.NoError(t, err)
	assert.False(t, active)
}

func TestAccountGatePropagatesUnexpectedError(t *testing.T) {
	t.Parallel()

	failure := errors.New("connection refused")
	gate := newAccountGate(&stubUserRepository{err: failure})

	active, err := gate.IsActive(context.Background(), uuid.New())

	require.ErrorIs(t, err, failure)
	assert.False(t, active)
}
