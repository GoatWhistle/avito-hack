package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/user/domain"
)

type UpdateProfileCommand struct {
	UserID      uuid.UUID
	DisplayName string
}

type UpdateProfileHandler struct {
	users domain.Repository
	tx    TxManager
	clock Clock
}

func NewUpdateProfileHandler(users domain.Repository, tx TxManager, clock Clock) *UpdateProfileHandler {
	return &UpdateProfileHandler{users: users, tx: tx, clock: clock}
}

func (h *UpdateProfileHandler) Handle(ctx context.Context, cmd UpdateProfileCommand) (*domain.User, error) {
	var updated *domain.User

	err := h.tx.WithTx(ctx, func(ctx context.Context) error {
		user, err := h.users.ByID(ctx, cmd.UserID)
		if err != nil {
			return fmt.Errorf("load user: %w", err)
		}

		if err := user.Rename(cmd.DisplayName, h.clock.Now()); err != nil {
			return err
		}

		if err := h.users.Save(ctx, user); err != nil {
			return err
		}

		updated = user

		return nil
	})
	if err != nil {
		return nil, err
	}

	return updated, nil
}
