package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/ws"
)

func TestWebSocketRejectsMissingToken(t *testing.T) {
	t.Parallel()

	handler := NewWebSocketHandler(scriptedService{}, stubTokens{}, ws.NewHub(), nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ws/pet", http.NoBody))

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestWebSocketRejectsUnparseableToken(t *testing.T) {
	t.Parallel()

	handler := NewWebSocketHandler(scriptedService{}, failingTokens{}, ws.NewHub(), nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ws/pet?token=forged", http.NoBody))

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestWebSocketRejectsNonUpgradeRequest(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
	handler := NewWebSocketHandler(scriptedService{}, stubTokens{actor: actor}, ws.NewHub(), nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ws/pet?token=valid", http.NoBody))

	assert.NotEqual(t, http.StatusUnauthorized, rec.Code)
}

func TestWebSocketReturnsStateForGetMessage(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	handler := NewWebSocketHandler(scriptedService{pet: wsPet(userID)}, stubTokens{actor: actor}, ws.NewHub(), nil)

	conn := dialHandler(t, handler)
	response := exchange(t, conn, clientMessage{Type: messageGetState, RequestID: "req-1"})

	assert.Equal(t, messageState, response.Type)
	assert.Equal(t, "req-1", response.RequestID)
	assert.NotNil(t, response.Payload)
}

func TestWebSocketReturnsStateForPetMessage(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	handler := NewWebSocketHandler(scriptedService{pet: wsPet(userID)}, stubTokens{actor: actor}, ws.NewHub(), nil)

	conn := dialHandler(t, handler)
	response := exchange(t, conn, clientMessage{Type: messagePet, RequestID: "req-2"})

	assert.Equal(t, messageState, response.Type)
	assert.Equal(t, "req-2", response.RequestID)
}

func TestWebSocketRejectsUnknownMessageType(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
	handler := NewWebSocketHandler(scriptedService{}, stubTokens{actor: actor}, ws.NewHub(), nil)

	conn := dialHandler(t, handler)
	response := exchange(t, conn, clientMessage{Type: "pet.teleport", RequestID: "req-3"})

	require.Equal(t, messageError, response.Type)
	assert.Equal(t, "req-3", response.RequestID)

	payload, ok := response.Payload.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "unknown_message", payload["code"])
}

func TestWebSocketReportsServiceFailures(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		service scriptedService
		message string
	}{
		{
			name:    "state failure",
			service: scriptedService{stateErr: errors.New("pg down")},
			message: messageGetState,
		},
		{
			name:    "stroke failure",
			service: scriptedService{strokeErr: errors.New("redis down")},
			message: messagePet,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
			handler := NewWebSocketHandler(tc.service, stubTokens{actor: actor}, ws.NewHub(), nil)

			conn := dialHandler(t, handler)
			response := exchange(t, conn, clientMessage{Type: tc.message, RequestID: "req-4"})

			require.Equal(t, messageError, response.Type)

			payload, ok := response.Payload.(map[string]any)
			require.True(t, ok)
			assert.Equal(t, "internal_error", payload["code"])
		})
	}
}

func TestWebSocketHandlesSequentialMessages(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	handler := NewWebSocketHandler(scriptedService{pet: wsPet(userID)}, stubTokens{actor: actor}, ws.NewHub(), nil)

	conn := dialHandler(t, handler)

	assert.Equal(t, messagePong, exchange(t, conn, clientMessage{Type: messagePing, RequestID: "1"}).Type)
	assert.Equal(t, messageState, exchange(t, conn, clientMessage{Type: messageGetState, RequestID: "2"}).Type)
	assert.Equal(t, messageError, exchange(t, conn, clientMessage{Type: "nope", RequestID: "3"}).Type)
}

func TestWebSocketWorksWithoutHub(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	handler := NewWebSocketHandler(scriptedService{pet: wsPet(userID)}, stubTokens{actor: actor}, nil, nil)

	conn := dialHandler(t, handler)
	response := exchange(t, conn, clientMessage{Type: messagePing, RequestID: "req-5"})

	assert.Equal(t, messagePong, response.Type)
}

func TestWebSocketBroadcastReachesConnectedClient(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	hub := ws.NewHub()
	handler := NewWebSocketHandler(scriptedService{pet: wsPet(userID)}, stubTokens{actor: actor}, hub, nil)

	conn := dialHandler(t, handler)

	require.Equal(t, messagePong, exchange(t, conn, clientMessage{Type: messagePing, RequestID: "warmup"}).Type)

	hub.Broadcast(userID, ws.Message{Type: "level.up", Payload: map[string]any{"level": 3}})

	var pushed serverMessage
	require.NoError(t, wsjson.Read(t.Context(), conn, &pushed))
	assert.Equal(t, "level.up", pushed.Type)
	assert.Empty(t, pushed.RequestID)
}
