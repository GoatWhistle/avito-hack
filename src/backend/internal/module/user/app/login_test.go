package app_test

import (
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/user/app"
	"github.com/avito-hack/backend/internal/module/user/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

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
			require.ErrorIs(t, err, domainerr.ErrUnauthorized)
		})
	}
}

func medianLoginDuration(t *testing.T, handler *app.LoginUserHandler, cmd app.LoginUserCommand) time.Duration {
	t.Helper()

	const samples = 5

	durations := make([]time.Duration, 0, samples)

	for range samples {
		start := time.Now()
		_, err := handler.Handle(t.Context(), cmd)
		durations = append(durations, time.Since(start))

		require.ErrorIs(t, err, domain.ErrInvalidCredential)
	}

	slices.Sort(durations)

	return durations[samples/2]
}

func TestLoginUserDoesNotLeakAccountExistenceByTiming(t *testing.T) {
	t.Parallel()

	repo := newStubUsers()
	seedUser(t, repo, "timing@example.com")
	handler := app.NewLoginUserHandler(repo, stubTokens{token: "jwt"})

	_, _ = handler.Handle(t.Context(), app.LoginUserCommand{Email: "warmup@example.com", Password: testPassword})

	existing := medianLoginDuration(t, handler, app.LoginUserCommand{
		Email: "timing@example.com", Password: "wrong password",
	})
	missing := medianLoginDuration(t, handler, app.LoginUserCommand{
		Email: "ghost@example.com", Password: "wrong password",
	})

	require.Positive(t, existing)
	require.Positive(t, missing)

	ratio := float64(existing) / float64(missing)
	if ratio < 1 {
		ratio = 1 / ratio
	}

	assert.Lessf(t, ratio, 5.0,
		"login timing must not reveal account existence: existing=%s missing=%s", existing, missing)
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
