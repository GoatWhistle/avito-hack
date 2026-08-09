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
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
)

type scriptedService struct {
	pet       *domain.Pet
	stateErr  error
	strokeErr error
	feedErr   error
}

func (s scriptedService) State(context.Context, uuid.UUID) (*domain.Pet, error) {
	return s.pet, s.stateErr
}

func (s scriptedService) Stroke(context.Context, uuid.UUID) (*domain.Pet, error) {
	return s.pet, s.strokeErr
}

func (s scriptedService) Feed(context.Context, uuid.UUID) (*domain.Pet, error) {
	return s.pet, s.feedErr
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
