package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
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
	handler := app.NewRegisterUserHandler(repo, passthroughTx{}, fixedClock{}, nil)

	result, err := handler.Handle(t.Context(), app.RegisterUserCommand{
		Email:       "  NewUser@Example.COM ",
		Password:    testPassword,
		FullName: "  Ivan  ",
	})

	require.NoError(t, err)
	require.NotNil(t, result.User)
	assert.Equal(t, "newuser@example.com", result.User.Email().String())
	assert.Equal(t, "Ivan", result.User.FullName())
	assert.Equal(t, auth.RoleUser, result.User.Role())
	assert.Equal(t, fixedNow, result.User.CreatedAt())
	assert.Len(t, repo.byID, 1)
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
			handler := app.NewRegisterUserHandler(repo, passthroughTx{}, fixedClock{}, nil)

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
	handler := app.NewRegisterUserHandler(repo, passthroughTx{}, fixedClock{}, nil)
	cmd := app.RegisterUserCommand{Email: "dup@example.com", Password: testPassword, FullName: "Ivan"}

	_, err := handler.Handle(t.Context(), cmd)
	require.NoError(t, err)

	_, err = handler.Handle(t.Context(), cmd)
	require.ErrorIs(t, err, domain.ErrEmailAlreadyTaken)
	assert.ErrorIs(t, err, domainerr.ErrConflict)
	assert.Len(t, repo.byID, 1)
}

func TestRegisterUserPropagatesRepositoryFailures(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("db is down")

	t.Run("exists check fails", func(t *testing.T) {
		t.Parallel()

		repo := newStubUsers()
		repo.existsErr = sentinel
		handler := app.NewRegisterUserHandler(repo, passthroughTx{}, fixedClock{}, nil)

		_, err := handler.Handle(t.Context(), app.RegisterUserCommand{
			Email: "err@example.com", Password: testPassword, FullName: "Ivan",
		})

		require.ErrorIs(t, err, sentinel)
	})

	t.Run("save fails", func(t *testing.T) {
		t.Parallel()

		repo := newStubUsers()
		repo.saveErr = sentinel
		handler := app.NewRegisterUserHandler(repo, passthroughTx{}, fixedClock{}, nil)

		_, err := handler.Handle(t.Context(), app.RegisterUserCommand{
			Email: "err2@example.com", Password: testPassword, FullName: "Ivan",
		})

		require.ErrorIs(t, err, sentinel)
	})
}

func TestLoginUserSuccess(t *testing.T) {
	t.Parallel()

	repo := newStubUsers()
	registered := seedUser(t, repo, "login@example.com")

	handler := app.NewLoginUserHandler(repo, stubTokens{token: "jwt-token"})

	result, err := handler.Handle(t.Context(), app.LoginUserCommand{
		Email: "LOGIN@example.com", Password: testPassword,
	})

	require.NoError(t, err)
	assert.Equal(t, "jwt-token", result.Token)
	assert.Equal(t, registered.ID(), result.User.ID())
	assert.Equal(t, fixedNow.Add(time.Hour), result.ExpiresAt)
}

func TestLoginUserRejectsBadCredentials(t *testing.T) {
	t.Parallel()

	repo := newStubUsers()
	seedUser(t, repo, "creds@example.com")

	handler := app.NewLoginUserHandler(repo, stubTokens{token: "jwt"})

	tests := []struct {
		name string
		cmd  app.LoginUserCommand
	}{
		{name: "malformed email", cmd: app.LoginUserCommand{Email: "bad", Password: testPassword}},
		{name: "unknown email", cmd: app.LoginUserCommand{Email: "ghost@example.com", Password: testPassword}},
		{name: "wrong password", cmd: app.LoginUserCommand{Email: "creds@example.com", Password: "wrong pass!!"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := handler.Handle(t.Context(), tc.cmd)

			require.ErrorIs(t, err, domain.ErrInvalidCredential)
			assert.ErrorIs(t, err, domainerr.ErrUnauthorized)
		})
	}
}

func TestLoginUserPropagatesInfraErrors(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("boom")

	t.Run("repository failure", func(t *testing.T) {
		t.Parallel()

		repo := newStubUsers()
		repo.byEmailErr = sentinel
		handler := app.NewLoginUserHandler(repo, stubTokens{token: "jwt"})

		_, err := handler.Handle(t.Context(), app.LoginUserCommand{
			Email: "any@example.com", Password: testPassword,
		})

		require.ErrorIs(t, err, sentinel)
	})

	t.Run("token issuing failure", func(t *testing.T) {
		t.Parallel()

		repo := newStubUsers()
		seedUser(t, repo, "tok@example.com")
		handler := app.NewLoginUserHandler(repo, stubTokens{err: sentinel})

		_, err := handler.Handle(t.Context(), app.LoginUserCommand{
			Email: "tok@example.com", Password: testPassword,
		})

		require.ErrorIs(t, err, sentinel)
	})
}

func TestGetProfile(t *testing.T) {
	t.Parallel()

	repo := newStubUsers()
	seeded := seedUser(t, repo, "profile@example.com")

	handler := app.NewGetProfileHandler(repo)

	user, err := handler.Handle(t.Context(), app.GetProfileQuery{UserID: seeded.ID()})
	require.NoError(t, err)
	assert.Equal(t, seeded.ID(), user.ID())

	_, err = handler.Handle(t.Context(), app.GetProfileQuery{UserID: uuid.New()})
	require.ErrorIs(t, err, domain.ErrUserNotFound)
	assert.ErrorIs(t, err, domainerr.ErrNotFound)
}

func TestUpdateProfile(t *testing.T) {
	t.Parallel()

	repo := newStubUsers()
	seeded := seedUser(t, repo, "update@example.com")

	handler := app.NewUpdateProfileHandler(repo, passthroughTx{}, fixedClock{})

	updated, err := handler.Handle(t.Context(), app.UpdateProfileCommand{
		UserID: seeded.ID(), FullName: "  Renamed  ",
	})
	require.NoError(t, err)
	assert.Equal(t, "Renamed", updated.FullName())
	assert.Equal(t, fixedNow, updated.UpdatedAt())
}

func TestUpdateProfileFailures(t *testing.T) {
	t.Parallel()

	repo := newStubUsers()
	seeded := seedUser(t, repo, "fail@example.com")
	handler := app.NewUpdateProfileHandler(repo, passthroughTx{}, fixedClock{})

	t.Run("missing user", func(t *testing.T) {
		t.Parallel()

		_, err := handler.Handle(t.Context(), app.UpdateProfileCommand{
			UserID: uuid.New(), FullName: "Whoever",
		})

		require.ErrorIs(t, err, domain.ErrUserNotFound)
	})

	t.Run("blank full name", func(t *testing.T) {
		t.Parallel()

		_, err := handler.Handle(t.Context(), app.UpdateProfileCommand{
			UserID: seeded.ID(), FullName: "   ",
		})

		var invalid *domainerr.InvalidError
		require.ErrorAs(t, err, &invalid)
		assert.Equal(t, "full_name", invalid.Field)
	})
}

func TestUpdateProfileRollsBackOnSaveFailure(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("save failed")
	repo := newStubUsers()
	seeded := seedUser(t, repo, "rollback@example.com")
	repo.saveErr = sentinel

	handler := app.NewUpdateProfileHandler(repo, passthroughTx{}, fixedClock{})

	_, err := handler.Handle(t.Context(), app.UpdateProfileCommand{
		UserID: seeded.ID(), FullName: "Rolled Back",
	})

	require.ErrorIs(t, err, sentinel)
}

func seedUser(t *testing.T, repo *stubUsers, email string) *domain.User {
	t.Helper()

	handler := app.NewRegisterUserHandler(repo, passthroughTx{}, fixedClock{}, nil)

	result, err := handler.Handle(context.Background(), app.RegisterUserCommand{
		Email: email, Password: testPassword, FullName: "Seeded",
	})
	require.NoError(t, err)

	return result.User
}
