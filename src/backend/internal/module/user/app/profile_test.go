package app_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/user/app"
	"github.com/avito-hack/backend/internal/module/user/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

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
	require.ErrorIs(t, err, domainerr.ErrNotFound)
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
