package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/avito-hack/backend/internal/module/user/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/vo"
)

type LoginUserCommand struct {
	Email    string
	Password string
}

type LoginUserResult struct {
	User      *domain.User
	Token     string
	ExpiresAt time.Time
}

type LoginUserHandler struct {
	users  domain.Repository
	tokens TokenIssuer
}

func NewLoginUserHandler(users domain.Repository, tokens TokenIssuer) *LoginUserHandler {
	return &LoginUserHandler{users: users, tokens: tokens}
}

func (h *LoginUserHandler) Handle(ctx context.Context, cmd LoginUserCommand) (LoginUserResult, error) {
	email, err := vo.NewEmail(cmd.Email)
	if err != nil {
		return LoginUserResult{}, domain.ErrInvalidCredential
	}

	user, err := h.users.ByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			return LoginUserResult{}, domain.ErrInvalidCredential
		}

		return LoginUserResult{}, fmt.Errorf("load user: %w", err)
	}

	if err := user.Authenticate(cmd.Password); err != nil {
		return LoginUserResult{}, domain.ErrInvalidCredential
	}

	token, expiresAt, err := h.tokens.Issue(user.Actor())
	if err != nil {
		return LoginUserResult{}, fmt.Errorf("issue token: %w", err)
	}

	return LoginUserResult{User: user, Token: token, ExpiresAt: expiresAt}, nil
}
