package infra

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/password"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
	"github.com/avito-hack/backend/internal/shared/vo"
)

var mapperTime = time.Date(2023, 7, 1, 8, 0, 0, 0, time.UTC)

func validRow(id uuid.UUID) userRow {
	return userRow{
		id:           id,
		email:        "user@example.com",
		passwordHash: "$2a$12$hash",
		fullName:     "Full Name",
		role:         string(auth.RoleAdmin),
		createdAt:    mapperTime,
		updatedAt:    mapperTime,
	}
}

func TestToDomainRestoresUser(t *testing.T) {
	t.Parallel()

	id := uuid.New()

	user, err := toDomain(validRow(id))

	require.NoError(t, err)
	assert.Equal(t, id, user.ID())
	assert.Equal(t, "user@example.com", user.Email().String())
	assert.Equal(t, "$2a$12$hash", user.PasswordHash().String())
	assert.Equal(t, "Full Name", user.FullName())
	assert.Equal(t, auth.RoleAdmin, user.Role())
	assert.Equal(t, mapperTime, user.CreatedAt())
	assert.Equal(t, mapperTime, user.UpdatedAt())
}

func TestToDomainRejectsInvalidEmail(t *testing.T) {
	t.Parallel()

	row := validRow(uuid.New())
	row.email = "not-an-email"

	_, err := toDomain(row)

	require.ErrorIs(t, err, vo.ErrInvalidEmail)
	assert.Contains(t, err.Error(), "restore email for user")
}

func TestToDomainRejectsEmptyPasswordHash(t *testing.T) {
	t.Parallel()

	row := validRow(uuid.New())
	row.passwordHash = ""

	_, err := toDomain(row)

	require.ErrorIs(t, err, password.ErrEmptyHash)
	assert.Contains(t, err.Error(), "restore password hash for user")
}

func TestScanUserReadsAllColumns(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	row := pgtest.Row{Values: []any{
		id, "scan@example.com", "$2a$12$scan", "Scan User",
		string(auth.RoleModerator), mapperTime, mapperTime,
	}}

	user, err := scanUser(row)

	require.NoError(t, err)
	assert.Equal(t, id, user.ID())
	assert.Equal(t, "scan@example.com", user.Email().String())
	assert.Equal(t, auth.RoleModerator, user.Role())
}

func TestScanUserPropagatesScanError(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("scan failed")

	_, err := scanUser(pgtest.Row{Err: sentinel})

	require.ErrorIs(t, err, sentinel)
}

func TestScanUserPropagatesMappingError(t *testing.T) {
	t.Parallel()

	row := pgtest.Row{Values: []any{
		uuid.New(), "broken", "$2a$12$scan", "Scan User",
		string(auth.RoleUser), mapperTime, mapperTime,
	}}

	_, err := scanUser(row)

	require.ErrorIs(t, err, vo.ErrInvalidEmail)
}

func TestIsUniqueViolation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{name: "plain error", err: errors.New("boom"), want: false},
		{name: "other pg code", err: &pgconn.PgError{Code: "23503"}, want: false},
		{name: "unique violation", err: &pgconn.PgError{Code: uniqueViolationCode}, want: true},
		{
			name: "wrapped unique violation",
			err:  errors.Join(errors.New("ctx"), &pgconn.PgError{Code: uniqueViolationCode}),
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, isUniqueViolation(tc.err))
		})
	}
}
