package api

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/ws"
)

type scriptedService struct {
	pet       *domain.Pet
	stateErr  error
	strokeErr error
}

func (s scriptedService) State(context.Context, uuid.UUID) (*domain.Pet, error) {
	return s.pet, s.stateErr
}

func (s scriptedService) Stroke(context.Context, uuid.UUID) (*domain.Pet, error) {
	return s.pet, s.strokeErr
}

type failingTokens struct{}

func (failingTokens) Parse(string) (auth.Actor, error) {
	return auth.Actor{}, errors.New("signature mismatch")
}

func wsPet(userID uuid.UUID) *domain.Pet {
	return domain.Restore(domain.RestoreParams{
		ID: uuid.New(), UserID: userID, Name: "Enot", Stage: domain.StageBaby,
		Level: 2, XP: 20, Satiety: 60, Happiness: 70, Energy: 80,
	})
}

func dialHandler(t *testing.T, handler http.Handler) *websocket.Conn {
	t.Helper()

	var config net.ListenConfig
	listener, err := config.Listen(t.Context(), "tcp4", "127.0.0.1:0")
	if err != nil {
		t.Skipf("local sockets are unavailable: %v", err)
	}

	server := httptest.NewUnstartedServer(handler)
	server.Listener = listener
	server.Start()
	t.Cleanup(server.Close)

	address := "ws" + strings.TrimPrefix(server.URL, "http") + "?token=valid"
	conn, resp, err := websocket.Dial(t.Context(), address, nil)
	require.NoError(t, err)
	if resp != nil && resp.Body != nil {
		require.NoError(t, resp.Body.Close())
	}
	t.Cleanup(func() { _ = conn.Close(websocket.StatusNormalClosure, "test completed") })

	return conn
}

func exchange(t *testing.T, conn *websocket.Conn, sent clientMessage) serverMessage {
	t.Helper()

	require.NoError(t, wsjson.Write(t.Context(), conn, sent))

	var response serverMessage
	require.NoError(t, wsjson.Read(t.Context(), conn, &response))

	return response
}

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
	response := exchange(t, conn, clientMessage{Type: "pet.feed", RequestID: "req-3"})

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
