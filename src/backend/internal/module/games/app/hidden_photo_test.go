package app_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/games/app"
	"github.com/avito-hack/backend/internal/module/games/domain"
	"github.com/avito-hack/backend/internal/module/games/domain/moreless"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

const photoSecret = "hidden-photo-test-secret"

func seedMorelessRound(t *testing.T, rounds *stubRounds, userID uuid.UUID) *domain.Round {
	t.Helper()

	round := domain.NewRound(userID, moreless.Slug, testNow)
	round.SetPayload(json.RawMessage(
		`{"left":{"DisplayID":"leftitem0001","PriceKopeks":100},` +
			`"right":{"DisplayID":"rightitem001","PhotoURL":"/uploads/rightitem001/p.jpg","PriceKopeks":900},` +
			`"seen":["leftitem0001","rightitem001"]}`))
	require.NoError(t, rounds.Save(context.Background(), round))

	return round
}

func TestHiddenPhotoResolvesForTheOwner(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	userID := uuid.New()
	round := seedMorelessRound(t, rounds, userID)

	signer := moreless.NewPhotoSigner(photoSecret)
	handler := app.NewHiddenPhotoHandler(rounds, signer)

	url, err := handler.Handle(context.Background(), app.HiddenPhotoQuery{
		UserID: userID, Token: signer.Sign(round.DisplayID(), "rightitem001"),
	})

	require.NoError(t, err)
	assert.Equal(t, "/uploads/rightitem001/p.jpg", url)
}

func TestHiddenPhotoOfAnotherUserIsNotFound(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	round := seedMorelessRound(t, rounds, uuid.New())

	signer := moreless.NewPhotoSigner(photoSecret)
	handler := app.NewHiddenPhotoHandler(rounds, signer)

	_, err := handler.Handle(context.Background(), app.HiddenPhotoQuery{
		UserID: uuid.New(), Token: signer.Sign(round.DisplayID(), "rightitem001"),
	})

	require.ErrorIs(t, err, domainerr.ErrNotFound)
}

func TestHiddenPhotoRejectsForgedToken(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	userID := uuid.New()
	round := seedMorelessRound(t, rounds, userID)

	handler := app.NewHiddenPhotoHandler(rounds, moreless.NewPhotoSigner(photoSecret))

	for name, token := range map[string]string{
		"foreign key": moreless.NewPhotoSigner("other-secret").Sign(round.DisplayID(), "rightitem001"),
		"garbage":     "mlp.abc.def",
		"empty":       "",
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := handler.Handle(context.Background(), app.HiddenPhotoQuery{UserID: userID, Token: token})
			require.ErrorIs(t, err, domainerr.ErrNotFound)
		})
	}
}

func TestHiddenPhotoRefusesItemOutsideTheRound(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	userID := uuid.New()
	round := seedMorelessRound(t, rounds, userID)

	signer := moreless.NewPhotoSigner(photoSecret)
	handler := app.NewHiddenPhotoHandler(rounds, signer)

	_, err := handler.Handle(context.Background(), app.HiddenPhotoQuery{
		UserID: userID, Token: signer.Sign(round.DisplayID(), "leftitem0001"),
	})

	require.ErrorIs(t, err, domainerr.ErrNotFound,
		"a token must not become a lookup oracle for arbitrary listings")
}
