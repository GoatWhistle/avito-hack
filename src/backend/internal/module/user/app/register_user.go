package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/avito-hack/backend/internal/module/user/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/password"
	"github.com/avito-hack/backend/internal/shared/vo"
)

type RegisterUserCommand struct {
	Email       string
	Password    string
	DisplayName string
}

type RegisterUserResult struct {
	User *domain.User
}

type RegisterUserHandler struct {
	users domain.Repository
	tx    TxManager
	clock Clock
}

func NewRegisterUserHandler(users domain.Repository, tx TxManager, clock Clock) *RegisterUserHandler {
	return &RegisterUserHandler{users: users, tx: tx, clock: clock}
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
		DisplayName:  cmd.DisplayName,
		Role:         auth.RoleUser,
		Now:          h.clock.Now(),
	})
	if err != nil {
		return RegisterUserResult{}, err
	}

	err = h.tx.WithTx(ctx, func(ctx context.Context) error {
		taken, existsErr := h.users.ExistsByEmail(ctx, email)
		if existsErr != nil {
			return fmt.Errorf("check email: %w", existsErr)
		}
		if taken {
			return domain.ErrEmailAlreadyTaken
		}

		return h.users.Save(ctx, user)
	})
	if err != nil {
		return RegisterUserResult{}, err
	}

	return RegisterUserResult{User: user}, nil
}

func mapPasswordError(err error) error {
	switch {
	case errors.Is(err, password.ErrTooShort):
		return domainerr.NewInvalid("password", fmt.Sprintf("value is shorter than minimum: %d", password.MinLength))
	case errors.Is(err, password.ErrTooLong):
		return domainerr.NewInvalid("password", fmt.Sprintf("value is longer than maximum: %d", password.MaxLength))
	default:
		return fmt.Errorf("hash password: %w", err)
	}
}
