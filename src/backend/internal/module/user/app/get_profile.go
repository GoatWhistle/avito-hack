package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/user/domain"
)

type GetProfileQuery struct {
	UserID uuid.UUID
}

type GetProfileHandler struct {
	users domain.Repository
}

func NewGetProfileHandler(users domain.Repository) *GetProfileHandler {
	return &GetProfileHandler{users: users}
}

func (h *GetProfileHandler) Handle(ctx context.Context, q GetProfileQuery) (*domain.User, error) {
	user, err := h.users.ByID(ctx, q.UserID)
	if err != nil {
		return nil, fmt.Errorf("load user: %w", err)
	}

	return user, nil
}
