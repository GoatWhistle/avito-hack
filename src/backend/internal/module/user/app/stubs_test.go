package app_test

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/user/app"
	"github.com/avito-hack/backend/internal/module/user/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/vo"
)

const testPassword = "correct horse battery"

var fixedNow = time.Date(2026, time.April, 1, 10, 0, 0, 0, time.UTC)

type fixedClock struct{}

func (fixedClock) Now() time.Time { return fixedNow }

type passthroughTx struct{}

func (passthroughTx) WithTx(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type stubTokens struct {
	token string
	err   error
}

func (s stubTokens) Issue(auth.Actor) (string, time.Time, error) {
	if s.err != nil {
		return "", time.Time{}, s.err
	}

	return s.token, fixedNow.Add(time.Hour), nil
}

type actorCapturingTokens struct {
	token string
	actor auth.Actor
}

func (s *actorCapturingTokens) Issue(actor auth.Actor) (string, time.Time, error) {
	s.actor = actor

	return s.token, fixedNow.Add(time.Hour), nil
}

func newRegisterHandler(users domain.Repository) *app.RegisterUserHandler {
	return app.NewRegisterUserHandler(
		users, passthroughTx{}, fixedClock{}, nil, stubTokens{token: "signed.jwt.token"})
}

type stubUsers struct {
	byID       map[uuid.UUID]*domain.User
	byEmail    map[string]*domain.User
	saveErr    error
	existsErr  error
	byEmailErr error
	byIDErr    error
}

func newStubUsers() *stubUsers {
	return &stubUsers{
		byID:    make(map[uuid.UUID]*domain.User),
		byEmail: make(map[string]*domain.User),
	}
}

func (s *stubUsers) Save(_ context.Context, user *domain.User) error {
	if s.saveErr != nil {
		return s.saveErr
	}

	s.byID[user.ID()] = user
	s.byEmail[user.Email().String()] = user

	return nil
}

func (s *stubUsers) ByID(_ context.Context, id uuid.UUID) (*domain.User, error) {
	if s.byIDErr != nil {
		return nil, s.byIDErr
	}

	user, ok := s.byID[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (s *stubUsers) ByEmail(_ context.Context, email vo.Email) (*domain.User, error) {
	if s.byEmailErr != nil {
		return nil, s.byEmailErr
	}

	user, ok := s.byEmail[email.String()]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (s *stubUsers) ExistsByEmail(_ context.Context, email vo.Email) (bool, error) {
	if s.existsErr != nil {
		return false, s.existsErr
	}

	_, ok := s.byEmail[email.String()]

	return ok, nil
}
