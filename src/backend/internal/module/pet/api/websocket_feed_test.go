package api

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/ws"
)

var errFeedFailed = errors.New("feed failed")

func TestWebSocketFeedReturnsState(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	handler := NewWebSocketHandler(
		scriptedService{pet: wsPet(userID)}, stubTokens{actor: actor}, ws.NewHub(), nil)

	conn := dialHandler(t, handler)
	response := exchange(t, conn, clientMessage{Type: messageFeed, RequestID: "req-feed"})

	assert.Equal(t, messageState, response.Type)
	assert.Equal(t, "req-feed", response.RequestID)
}

func TestWebSocketFeedReportsServiceFailure(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
	handler := NewWebSocketHandler(
		scriptedService{feedErr: errFeedFailed}, stubTokens{actor: actor}, ws.NewHub(), nil)

	conn := dialHandler(t, handler)
	response := exchange(t, conn, clientMessage{Type: messageFeed, RequestID: "req-fail"})

	assert.Equal(t, messageError, response.Type)

	payload, ok := response.Payload.(map[string]any)
	if assert.True(t, ok) {
		assert.Equal(t, "internal_error", payload["code"])
	}
}
