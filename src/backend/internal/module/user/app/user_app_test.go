package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/user/app"
	"github.com/avito-hack/backend/internal/module/user/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

func TestRegisterUserSuccess(t *testing.T) {
	t.Parallel()

	repo := newStubUsers()
	handler := newRegisterHandler(repo)

	result, err := handler.Handle(t.Context(), app.RegisterUserCommand{
		Email:    "  NewUser@Example.COM ",
		Password: testPassword,
		FullName: "  Ivan  ",
	})

	require.NoError(t, err)
	require.NotNil(t, result.User)
	assert.Equal(t, "newuser@example.com", result.User.Email().String())
	assert.Equal(t, "Ivan", result.User.FullName())
	assert.Equal(t, auth.RoleUser, result.User.Role())
	assert.Equal(t, fixedNow, result.User.CreatedAt())
	assert.Len(t, repo.byID, 1)
	assert.Equal(t, "signed.jwt.token", result.Token)
}

func TestRegisterUserValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		cmd       app.RegisterUserCommand
		wantField string
	}{
		{
			name:      "invalid email",
			cmd:       app.RegisterUserCommand{Email: "not-an-email", Password: testPassword, FullName: "Ivan"},
			wantField: "email",
		},
		{
			name:      "empty email",
			cmd:       app.RegisterUserCommand{Email: "", Password: testPassword, FullName: "Ivan"},
			wantField: "email",
		},
		{
			name:      "empty password",
			cmd:       app.RegisterUserCommand{Email: "a@example.com", Password: "", FullName: "Ivan"},
			wantField: "password",
		},
		{
			name:      "blank full name",
			cmd:       app.RegisterUserCommand{Email: "a@example.com", Password: testPassword, FullName: "  "},
			wantField: "full_name",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := newStubUsers()
			handler := newRegisterHandler(repo)

			_, err := handler.Handle(t.Context(), tc.cmd)

			var invalid *domainerr.InvalidError
			require.ErrorAs(t, err, &invalid)
			assert.Equal(t, tc.wantField, invalid.Field)
			assert.Empty(t, repo.byID)
		})
	}
}

func TestRegisterUserRejectsTakenEmail(t *testing.T) {
	t.Parallel()

	repo := newStubUsers()
	handler := newRegisterHandler(repo)
	cmd := app.RegisterUserCommand{Email: "dup@example.com", Password: testPassword, FullName: "Ivan"}

	_, err := handler.Handle(t.Context(), cmd)
	require.NoError(t, err)

	_, err = handler.Handle(t.Context(), cmd)
	require.ErrorIs(t, err, domain.ErrEmailAlreadyTaken)
	require.ErrorIs(t, err, domainerr.ErrConflict)
	assert.Len(t, repo.byID, 1)
}

func TestRegisterUserPropagatesRepositoryFailures(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("db is down")

	t.Run("exists check fails", func(t *testing.T) {
		t.Parallel()

		repo := newStubUsers()
		repo.existsErr = sentinel
		handler := newRegisterHandler(repo)

		_, err := handler.Handle(t.Context(), app.RegisterUserCommand{
			Email: "err@example.com", Password: testPassword, FullName: "Ivan",
		})

		require.ErrorIs(t, err, sentinel)
	})

	t.Run("save fails", func(t *testing.T) {
		t.Parallel()

		repo := newStubUsers()
		repo.saveErr = sentinel
		handler := newRegisterHandler(repo)

		_, err := handler.Handle(t.Context(), app.RegisterUserCommand{
			Email: "err2@example.com", Password: testPassword, FullName: "Ivan",
		})

		require.ErrorIs(t, err, sentinel)
	})
}

func seedUser(t *testing.T, repo *stubUsers, email string) *domain.User {
	t.Helper()

	handler := newRegisterHandler(repo)

	result, err := handler.Handle(context.Background(), app.RegisterUserCommand{
		Email: email, Password: testPassword, FullName: "Seeded",
	})
	require.NoError(t, err)

	return result.User
}
