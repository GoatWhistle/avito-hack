package infra

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/avito-hack/backend/internal/module/user/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/password"
	"github.com/avito-hack/backend/internal/shared/vo"
)

const uniqueViolationCode = "23505"

type userRow struct {
	id           uuid.UUID
	email        string
	passwordHash string
	displayName  string
	role         string
	createdAt    time.Time
	updatedAt    time.Time
}

func scanUser(row pgx.Row) (*domain.User, error) {
	var r userRow

	err := row.Scan(&r.id, &r.email, &r.passwordHash, &r.displayName, &r.role, &r.createdAt, &r.updatedAt)
	if err != nil {
		return nil, err
	}

	return toDomain(r)
}

func toDomain(r userRow) (*domain.User, error) {
	email, err := vo.NewEmail(r.email)
	if err != nil {
		return nil, fmt.Errorf("restore email for user %s: %w", r.id, err)
	}

	hash, err := password.RestoreHash(r.passwordHash)
	if err != nil {
		return nil, fmt.Errorf("restore password hash for user %s: %w", r.id, err)
	}

	return domain.RestoreUser(domain.RestoreUserParams{
		ID:           r.id,
		Email:        email,
		PasswordHash: hash,
		DisplayName:  r.displayName,
		Role:         auth.Role(r.role),
		CreatedAt:    r.createdAt,
		UpdatedAt:    r.updatedAt,
	}), nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode
}
