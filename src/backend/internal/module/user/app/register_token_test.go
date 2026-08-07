package app_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/user/app"
	"github.com/avito-hack/backend/internal/shared/auth"
)

func TestRegisterUserIssuesTokenForCreatedUser(t *testing.T) {
	t.Parallel()

	repo := newStubUsers()
	tokens := &actorCapturingTokens{token: "issued.jwt"}
	handler := app.NewRegisterUserHandler(repo, passthroughTx{}, fixedClock{}, nil, tokens)

	result, err := handler.Handle(t.Context(), app.RegisterUserCommand{
		Email:    "fresh@example.com",
		Password: testPassword,
		FullName: "Ivan",
	})

	require.NoError(t, err)
	assert.Equal(t, "issued.jwt", result.Token)
	assert.Equal(t, result.User.ID(), tokens.actor.ID)
	assert.Equal(t, auth.RoleUser, tokens.actor.Role)
}

func TestRegisterUserFailsWhenTokenIssuerFails(t *testing.T) {
	t.Parallel()

	repo := newStubUsers()
	sentinel := errors.New("signer down")
	handler := app.NewRegisterUserHandler(repo, passthroughTx{}, fixedClock{}, nil, stubTokens{err: sentinel})

	_, err := handler.Handle(t.Context(), app.RegisterUserCommand{
		Email:    "fresh@example.com",
		Password: testPassword,
		FullName: "Ivan",
	})

	require.ErrorIs(t, err, sentinel)
}
