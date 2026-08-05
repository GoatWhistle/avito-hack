package infra

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/user/domain"
	"github.com/avito-hack/backend/internal/shared/postgres"
	"github.com/avito-hack/backend/internal/shared/vo"
)

const userColumns = `id, email, password_hash, full_name, role, created_at, updated_at`

type PgRepository struct {
	pool *pgxpool.Pool
}

func NewPgRepository(pool *pgxpool.Pool) *PgRepository {
	return &PgRepository{pool: pool}
}

func (r *PgRepository) Save(ctx context.Context, user *domain.User) error {
	const query = `
		INSERT INTO users (id, email, password_hash, full_name, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			password_hash = EXCLUDED.password_hash,
			full_name = EXCLUDED.full_name,
			role = EXCLUDED.role,
			updated_at = EXCLUDED.updated_at`

	_, err := postgres.QuerierFrom(ctx, r.pool).Exec(ctx, query,
		user.ID(),
		user.Email().String(),
		user.PasswordHash().String(),
		user.FullName(),
		string(user.Role()),
		user.CreatedAt(),
		user.UpdatedAt(),
	)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrEmailAlreadyTaken
		}

		return fmt.Errorf("save user: %w", err)
	}

	return nil
}

func (r *PgRepository) ByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT ` + userColumns + ` FROM users WHERE id = $1 AND deleted_at IS NULL`

	return r.queryOne(ctx, query, id)
}

func (r *PgRepository) ByEmail(ctx context.Context, email vo.Email) (*domain.User, error) {
	query := `SELECT ` + userColumns + ` FROM users WHERE email = $1 AND deleted_at IS NULL`

	return r.queryOne(ctx, query, email.String())
}

func (r *PgRepository) ExistsByEmail(ctx context.Context, email vo.Email) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)`

	var exists bool
	if err := postgres.QuerierFrom(ctx, r.pool).QueryRow(ctx, query, email.String()).Scan(&exists); err != nil {
		return false, fmt.Errorf("check user email: %w", err)
	}

	return exists, nil
}

func (r *PgRepository) queryOne(ctx context.Context, query string, args ...any) (*domain.User, error) {
	row := postgres.QuerierFrom(ctx, r.pool).QueryRow(ctx, query, args...)

	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("query user: %w", err)
	}

	return user, nil
}
