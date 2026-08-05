package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/avito-hack/backend/internal/module/user/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/events"
	"github.com/avito-hack/backend/internal/shared/password"
	"github.com/avito-hack/backend/internal/shared/vo"
)

type RegisterUserCommand struct {
	Email    string
	Password string
	FullName string
}

type RegisterUserResult struct {
	User *domain.User
}

type RegisterUserHandler struct {
	users domain.Repository
	tx    TxManager
	clock Clock
	bus   events.Publisher
}

func NewRegisterUserHandler(
	users domain.Repository,
	tx TxManager,
	clock Clock,
	bus events.Publisher,
) *RegisterUserHandler {
	if bus == nil {
		bus = events.NopPublisher{}
	}

	return &RegisterUserHandler{users: users, tx: tx, clock: clock, bus: bus}
}

func (h *RegisterUserHandler) Handle(ctx context.Context, cmd RegisterUserCommand) (RegisterUserResult, error) {
	email, err := vo.NewEmail(cmd.Email)
	if err != nil {
		return RegisterUserResult{}, domainerr.NewInvalid("email", "invalid email address")
	}

	hash, err := password.NewHash(cmd.Password)
	if err != nil {
		return RegisterUserResult{}, mapPasswordError(err)
	}

	user, err := domain.NewUser(domain.NewUserParams{
		Email:        email,
		PasswordHash: hash,
		FullName:     cmd.FullName,
		Role:         auth.RoleUser,
		Now:          h.clock.Now(),
	})
	if err != nil {
		return RegisterUserResult{}, err
	}

	err = events.PublishAfterCommit(ctx, h.bus, h.tx, func(ctx context.Context, out *events.Outbox) error {
		taken, existsErr := h.users.ExistsByEmail(ctx, email)
		if existsErr != nil {
			return fmt.Errorf("check email: %w", existsErr)
		}
		if taken {
			return domain.ErrEmailAlreadyTaken
		}

		if saveErr := h.users.Save(ctx, user); saveErr != nil {
			return saveErr
		}

		out.Add(events.New(events.TypeUserRegistered, user.ID(), user.ID(), h.clock.Now()))

		return nil
	})
	if err != nil {
		return RegisterUserResult{}, err
	}

	return RegisterUserResult{User: user}, nil
}

func mapPasswordError(err error) error {
	switch {
	case errors.Is(err, password.ErrTooShort), errors.Is(err, password.ErrTooLong):
		return domainerr.NewInvalid("password", "field is required")
	default:
		return fmt.Errorf("hash password: %w", err)
	}
}
