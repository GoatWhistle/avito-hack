package domain_test

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/user/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/password"
	"github.com/avito-hack/backend/internal/shared/vo"
)

const defaultPassword = "correct horse battery"

var now = time.Date(2026, time.March, 10, 9, 0, 0, 0, time.UTC)

func mustEmail(t *testing.T, raw string) vo.Email {
	t.Helper()

	email, err := vo.NewEmail(raw)
	require.NoError(t, err)

	return email
}

var (
	hashOnce   sync.Once
	cachedHash password.Hash
)

func mustHash(t *testing.T, plain string) password.Hash {
	t.Helper()

	if plain == defaultPassword {
		hashOnce.Do(func() {
			hash, err := password.NewHash(defaultPassword)
			require.NoError(t, err)
			cachedHash = hash
		})

		return cachedHash
	}

	hash, err := password.NewHash(plain)
	require.NoError(t, err)

	return hash
}

func TestNewUserValidation(t *testing.T) {
	t.Parallel()

	email := mustEmail(t, "user@example.com")
	hash := mustHash(t, defaultPassword)

	tests := []struct {
		name      string
		params    domain.NewUserParams
		wantErr   bool
		wantRole  auth.Role
		wantName  string
		wantField string
	}{
		{
			name:     "defaults role to user",
			params:   domain.NewUserParams{Email: email, PasswordHash: hash, FullName: "Ivan", Now: now},
			wantRole: auth.RoleUser,
			wantName: "Ivan",
		},
		{
			name: "keeps explicit moderator role",
			params: domain.NewUserParams{
				Email: email, PasswordHash: hash, FullName: "Mod", Role: auth.RoleModerator, Now: now,
			},
			wantRole: auth.RoleModerator,
			wantName: "Mod",
		},
		{
			name: "trims display name",
			params: domain.NewUserParams{
				Email: email, PasswordHash: hash, FullName: "   Anna   ", Now: now,
			},
			wantRole: auth.RoleUser,
			wantName: "Anna",
		},
		{
			name: "rejects blank full name",
			params: domain.NewUserParams{
				Email: email, PasswordHash: hash, FullName: "     ", Now: now,
			},
			wantErr:   true,
			wantField: "full_name",
		},
		{
			name:      "rejects empty full name",
			params:    domain.NewUserParams{Email: email, PasswordHash: hash, Now: now},
			wantErr:   true,
			wantField: "full_name",
		},
		{
			name:      "rejects zero email",
			params:    domain.NewUserParams{PasswordHash: hash, FullName: "Ivan", Now: now},
			wantErr:   true,
			wantField: "email",
		},
		{
			name: "rejects unknown role",
			params: domain.NewUserParams{
				Email: email, PasswordHash: hash, FullName: "Ivan", Role: auth.Role("root"), Now: now,
			},
			wantErr:   true,
			wantField: "role",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			user, err := domain.NewUser(tc.params)

			if tc.wantErr {
				require.Error(t, err)

				var invalid *domainerr.InvalidError
				require.ErrorAs(t, err, &invalid)
				assert.Equal(t, tc.wantField, invalid.Field)

				return
			}

			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, user.ID())
			assert.Equal(t, tc.wantRole, user.Role())
			assert.Equal(t, tc.wantName, user.FullName())
			assert.Equal(t, now, user.CreatedAt())
			assert.Equal(t, now, user.UpdatedAt())
			assert.Equal(t, user.ID(), user.Actor().ID)
			assert.Equal(t, tc.wantRole, user.Actor().Role)
		})
	}
}

func TestUserAcceptsLongFullName(t *testing.T) {
	t.Parallel()

	email := mustEmail(t, "edge@example.com")
	hash := mustHash(t, defaultPassword)

	for _, length := range []int{1, 100, 500} {
		user, err := domain.NewUser(domain.NewUserParams{
			Email: email, PasswordHash: hash, FullName: strings.Repeat("a", length), Now: now,
		})

		require.NoError(t, err)
		assert.Len(t, user.FullName(), length)
	}
}

func TestUserAuthenticate(t *testing.T) {
	t.Parallel()

	user, err := domain.NewUser(domain.NewUserParams{
		Email:        mustEmail(t, "auth@example.com"),
		PasswordHash: mustHash(t, defaultPassword),
		FullName:  "Auth",
		Now:          now,
	})
	require.NoError(t, err)

	require.NoError(t, user.Authenticate(defaultPassword))
	assert.Error(t, user.Authenticate("wrong pass value"))
}

func TestUserRename(t *testing.T) {
	t.Parallel()

	user, err := domain.NewUser(domain.NewUserParams{
		Email:        mustEmail(t, "rename@example.com"),
		PasswordHash: mustHash(t, defaultPassword),
		FullName:  "Before",
		Now:          now,
	})
	require.NoError(t, err)

	later := now.Add(time.Hour)

	require.NoError(t, user.Rename("  After  ", later))
	assert.Equal(t, "After", user.FullName())
	assert.Equal(t, later, user.UpdatedAt())

	require.Error(t, user.Rename("   ", later.Add(time.Hour)))
	assert.Equal(t, "After", user.FullName())
	assert.Equal(t, later, user.UpdatedAt())
}

func TestUserChangePassword(t *testing.T) {
	t.Parallel()

	user, err := domain.NewUser(domain.NewUserParams{
		Email:        mustEmail(t, "pwd@example.com"),
		PasswordHash: mustHash(t, defaultPassword),
		FullName:  "Pwd",
		Now:          now,
	})
	require.NoError(t, err)

	later := now.Add(2 * time.Hour)
	user.ChangePassword(mustHash(t, "another secret pass"), later)

	require.NoError(t, user.Authenticate("another secret pass"))
	assert.Error(t, user.Authenticate(defaultPassword))
	assert.Equal(t, later, user.UpdatedAt())
}

func TestRestoreUserKeepsAllFields(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	email := mustEmail(t, "restore@example.com")
	hash := mustHash(t, defaultPassword)
	updated := now.Add(time.Hour)

	user := domain.RestoreUser(domain.RestoreUserParams{
		ID: id, Email: email, PasswordHash: hash, FullName: "Restored",
		Role: auth.RoleAdmin, CreatedAt: now, UpdatedAt: updated,
	})

	assert.Equal(t, id, user.ID())
	assert.Equal(t, email, user.Email())
	assert.Equal(t, hash, user.PasswordHash())
	assert.Equal(t, "Restored", user.FullName())
	assert.Equal(t, auth.RoleAdmin, user.Role())
	assert.Equal(t, now, user.CreatedAt())
	assert.Equal(t, updated, user.UpdatedAt())
}
