package app_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/user/app"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

func TestRefreshSessionIssuesTokenForExistingUser(t *testing.T) {
	t.Parallel()

	repo := newStubUsers()
	registered := seedUser(t, repo, "refresh@example.com")

	handler := app.NewRefreshSessionHandler(repo, stubTokens{token: "fresh-token"})

	result, err := handler.Handle(t.Context(), app.RefreshSessionCommand{UserID: registered.ID()})

	require.NoError(t, err)
	assert.Equal(t, "fresh-token", result.Token)
	assert.Equal(t, registered.ID(), result.User.ID())
	assert.Equal(t, fixedNow.Add(time.Hour), result.ExpiresAt)
}

func TestRefreshSessionKeepsActorIdentity(t *testing.T) {
	t.Parallel()

	repo := newStubUsers()
	registered := seedUser(t, repo, "identity@example.com")

	tokens := &actorCapturingTokens{token: "fresh-token"}
	handler := app.NewRefreshSessionHandler(repo, tokens)

	_, err := handler.Handle(t.Context(), app.RefreshSessionCommand{UserID: registered.ID()})

	require.NoError(t, err)
	assert.Equal(t, registered.ID(), tokens.actor.ID)
	assert.Equal(t, registered.Role(), tokens.actor.Role)
}

func TestRefreshSessionFailsForMissingUser(t *testing.T) {
	t.Parallel()

	handler := app.NewRefreshSessionHandler(newStubUsers(), stubTokens{token: "fresh-token"})

	_, err := handler.Handle(t.Context(), app.RefreshSessionCommand{UserID: uuid.New()})

	require.Error(t, err)
	assert.ErrorIs(t, err, domainerr.ErrNotFound)
}

func TestRefreshSessionPropagatesIssuerFailure(t *testing.T) {
	t.Parallel()

	repo := newStubUsers()
	registered := seedUser(t, repo, "issuer@example.com")

	handler := app.NewRefreshSessionHandler(repo, stubTokens{err: errors.New("sign failed")})

	_, err := handler.Handle(t.Context(), app.RefreshSessionCommand{UserID: registered.ID()})

	require.Error(t, err)
}
