package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/vo"
)

var (
	ErrUserNotFound      = fmt.Errorf("user not found: %w", domainerr.ErrNotFound)
	ErrEmailAlreadyTaken = fmt.Errorf("email already taken: %w", domainerr.ErrConflict)
	ErrInvalidCredential = fmt.Errorf("invalid email or password: %w", domainerr.ErrUnauthorized)
)

type Repository interface {
	Save(ctx context.Context, user *User) error
	ByID(ctx context.Context, id uuid.UUID) (*User, error)
	ByEmail(ctx context.Context, email vo.Email) (*User, error)
	ExistsByEmail(ctx context.Context, email vo.Email) (bool, error)
}
