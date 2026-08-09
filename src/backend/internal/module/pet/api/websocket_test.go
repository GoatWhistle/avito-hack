package api

import (
	"context"
	"net"
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

type stubTokens struct {
	actor auth.Actor
}

func (s stubTokens) Parse(string) (auth.Actor, error) {
	return s.actor, nil
}

type stubService struct{}

func (stubService) State(context.Context, uuid.UUID) (*domain.Pet, error)  { return nil, nil }
func (stubService) Stroke(context.Context, uuid.UUID) (*domain.Pet, error) { return nil, nil }
func (stubService) Feed(context.Context, uuid.UUID) (*domain.Pet, error)   { return nil, nil }

func TestWebSocketPing(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
	var config net.ListenConfig
	listener, err := config.Listen(t.Context(), "tcp4", "127.0.0.1:0")
	if err != nil {
		t.Skipf("local sockets are unavailable: %v", err)
	}
	handler := NewWebSocketHandler(stubService{}, stubTokens{actor: actor}, ws.NewHub(), nil)
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

	require.NoError(t, wsjson.Write(t.Context(), conn, clientMessage{Type: messagePing, RequestID: "42"}))

	var response serverMessage
	require.NoError(t, wsjson.Read(t.Context(), conn, &response))
	assert.Equal(t, messagePong, response.Type)
	assert.Equal(t, "42", response.RequestID)
}
