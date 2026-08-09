package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/user/domain"
)

type RefreshSessionCommand struct {
	UserID uuid.UUID
}

type RefreshSessionResult struct {
	User      *domain.User
	Token     string
	ExpiresAt time.Time
}

type RefreshSessionHandler struct {
	users  domain.Repository
	tokens TokenIssuer
}

func NewRefreshSessionHandler(users domain.Repository, tokens TokenIssuer) *RefreshSessionHandler {
	return &RefreshSessionHandler{users: users, tokens: tokens}
}

func (h *RefreshSessionHandler) Handle(
	ctx context.Context,
	cmd RefreshSessionCommand,
) (RefreshSessionResult, error) {
	user, err := h.users.ByID(ctx, cmd.UserID)
	if err != nil {
		return RefreshSessionResult{}, fmt.Errorf("load user: %w", err)
	}

	token, expiresAt, err := h.tokens.Issue(user.Actor())
	if err != nil {
		return RefreshSessionResult{}, fmt.Errorf("issue token: %w", err)
	}

	return RefreshSessionResult{User: user, Token: token, ExpiresAt: expiresAt}, nil
}
