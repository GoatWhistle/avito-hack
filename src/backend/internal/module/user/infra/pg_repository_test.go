package infra_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/user/domain"
	"github.com/avito-hack/backend/internal/module/user/infra"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/password"
	"github.com/avito-hack/backend/internal/shared/postgres"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
	"github.com/avito-hack/backend/internal/shared/vo"
)

var errDB = errors.New("db down")

var fixedTime = time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)

const storedHash = "$2a$12$abcdefghijklmnopqrstuv"

func ctxWith(tx *pgtest.Tx) context.Context {
	return postgres.ContextWithTx(context.Background(), tx)
}

func newUser(t *testing.T) *domain.User {
	t.Helper()

	email, err := vo.NewEmail("owner@example.com")
	require.NoError(t, err)

	hash, err := password.RestoreHash(storedHash)
	require.NoError(t, err)

	return domain.RestoreUser(domain.RestoreUserParams{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hash,
		FullName:     "Owner Name",
		Role:         auth.RoleUser,
		CreatedAt:    fixedTime,
		UpdatedAt:    fixedTime,
	})
}

func userValues(t *testing.T, id uuid.UUID) []any {
	t.Helper()

	return []any{
		id, "owner1234567", "owner@example.com", storedHash, "Owner Name",
		string(auth.RoleUser), fixedTime, fixedTime,
	}
}

func TestPgRepositorySaveInsertsUser(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{}
	user := newUser(t)

	err := infra.NewPgRepository(nil).Save(ctxWith(tx), user)

	require.NoError(t, err)
	require.Len(t, tx.ExecCalls, 1)
	assert.Contains(t, tx.ExecCalls[0].SQL, "INSERT INTO users")
	assert.Contains(t, tx.ExecCalls[0].SQL, "ON CONFLICT (id) DO UPDATE")
	assert.Equal(t, []any{
		user.ID(), user.DisplayID(), "owner@example.com", storedHash, "Owner Name",
		string(auth.RoleUser), fixedTime, fixedTime,
	}, tx.ExecCalls[0].Args)
}

func TestPgRepositorySaveMapsUniqueViolation(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{ExecErrs: []error{&pgconn.PgError{Code: "23505"}}}

	err := infra.NewPgRepository(nil).Save(ctxWith(tx), newUser(t))

	require.ErrorIs(t, err, domain.ErrEmailAlreadyTaken)
}

func TestPgRepositorySaveWrapsOtherPgError(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{ExecErrs: []error{&pgconn.PgError{Code: "23503"}}}

	err := infra.NewPgRepository(nil).Save(ctxWith(tx), newUser(t))

	require.Error(t, err)
	require.NotErrorIs(t, err, domain.ErrEmailAlreadyTaken)
	assert.Contains(t, err.Error(), "save user")
}

func TestPgRepositorySaveWrapsError(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{ExecErrs: []error{errDB}}

	err := infra.NewPgRepository(nil).Save(ctxWith(tx), newUser(t))

	require.ErrorIs(t, err, errDB)
}

func TestPgRepositoryByIDReturnsUser(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Values: userValues(t, id)}}}

	user, err := infra.NewPgRepository(nil).ByID(ctxWith(tx), id)

	require.NoError(t, err)
	assert.Equal(t, id, user.ID())
	assert.Equal(t, "owner@example.com", user.Email().String())
	assert.Equal(t, auth.RoleUser, user.Role())
	require.Len(t, tx.QueryRowCalls, 1)
	assert.Contains(t, tx.QueryRowCalls[0].SQL, "FROM users WHERE id = $1 AND deleted_at IS NULL")
	assert.Equal(t, []any{id}, tx.QueryRowCalls[0].Args)
}

func TestPgRepositoryByIDReturnsNotFound(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Err: pgx.ErrNoRows}}}

	_, err := infra.NewPgRepository(nil).ByID(ctxWith(tx), uuid.New())

	require.ErrorIs(t, err, domain.ErrUserNotFound)
}

func TestPgRepositoryByIDWrapsError(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Err: errDB}}}

	_, err := infra.NewPgRepository(nil).ByID(ctxWith(tx), uuid.New())

	require.ErrorIs(t, err, errDB)
	assert.Contains(t, err.Error(), "query user")
}

func TestPgRepositoryByEmailReturnsUser(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Values: userValues(t, id)}}}

	email, err := vo.NewEmail("owner@example.com")
	require.NoError(t, err)

	user, err := infra.NewPgRepository(nil).ByEmail(ctxWith(tx), email)

	require.NoError(t, err)
	assert.Equal(t, id, user.ID())
	assert.Contains(t, tx.QueryRowCalls[0].SQL, "WHERE email = $1")
	assert.Equal(t, []any{"owner@example.com"}, tx.QueryRowCalls[0].Args)
}

func TestPgRepositoryByEmailReturnsNotFound(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Err: pgx.ErrNoRows}}}

	_, err := infra.NewPgRepository(nil).ByEmail(ctxWith(tx), vo.Email{})

	require.ErrorIs(t, err, domain.ErrUserNotFound)
}

func TestPgRepositoryExistsByEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		row  pgtest.Row
		want bool
		err  error
	}{
		{name: "exists", row: pgtest.Row{Values: []any{true}}, want: true},
		{name: "missing", row: pgtest.Row{Values: []any{false}}, want: false},
		{name: "error", row: pgtest.Row{Err: errDB}, err: errDB},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tx := &pgtest.Tx{RowResults: []pgtest.Row{tc.row}}

			email, err := vo.NewEmail("check@example.com")
			require.NoError(t, err)

			got, err := infra.NewPgRepository(nil).ExistsByEmail(ctxWith(tx), email)

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				assert.Contains(t, err.Error(), "check user email")

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
			assert.Contains(t, tx.QueryRowCalls[0].SQL, "SELECT EXISTS")
			assert.Equal(t, []any{"check@example.com"}, tx.QueryRowCalls[0].Args)
		})
	}
}
