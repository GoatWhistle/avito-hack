package api

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/ws"
)

func closedConnection(t *testing.T) *websocket.Conn {
	t.Helper()

	var config net.ListenConfig
	listener, err := config.Listen(t.Context(), "tcp4", "127.0.0.1:0")
	if err != nil {
		t.Skipf("local sockets are unavailable: %v", err)
	}

	accepted := make(chan struct{}, 1)
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, acceptErr := websocket.Accept(w, r, nil)
		if acceptErr != nil {
			return
		}
		defer func() { _ = conn.CloseNow() }()
		accepted <- struct{}{}
		<-r.Context().Done()
	}))
	server.Listener = listener
	server.Start()

	client, resp, err := websocket.Dial(t.Context(), "ws"+strings.TrimPrefix(server.URL, "http"), nil)
	require.NoError(t, err)
	if resp != nil && resp.Body != nil {
		require.NoError(t, resp.Body.Close())
	}

	<-accepted
	require.NoError(t, client.CloseNow())
	server.Close()

	return client
}

func TestConnectionSendEnqueuesMessage(t *testing.T) {
	t.Parallel()

	c := newConnection(nil)

	c.Send(ws.Message{Type: "pet.updated", Payload: map[string]int{"level": 3}})

	select {
	case m := <-c.out:
		assert.Equal(t, "pet.updated", m.Type)
		assert.Empty(t, m.RequestID)
		assert.Equal(t, map[string]int{"level": 3}, m.Payload)
	default:
		t.Fatal("expected the message to be queued")
	}
}

func TestConnectionReplyCarriesRequestID(t *testing.T) {
	t.Parallel()

	c := newConnection(nil)

	c.Reply("req-7", messagePong, nil)

	m := <-c.out
	assert.Equal(t, messagePong, m.Type)
	assert.Equal(t, "req-7", m.RequestID)
}

func TestConnectionSendAfterShutdownNeverBlocks(t *testing.T) {
	t.Parallel()

	c := newConnection(nil)
	c.shutdown()

	done := make(chan struct{})
	go func() {
		for range sendBuffer * 2 {
			c.Send(ws.Message{Type: "pet.updated"})
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Send blocked after shutdown")
	}

	assert.LessOrEqual(t, len(c.out), sendBuffer)
}

func TestConnectionShutdownIsIdempotent(t *testing.T) {
	t.Parallel()

	c := newConnection(nil)

	assert.NotPanics(t, func() {
		c.shutdown()
		c.shutdown()
		c.shutdown()
	})

	select {
	case <-c.closed:
	default:
		t.Fatal("expected the closed channel to be closed")
	}
}

func TestConnectionDropsMessagesWhenBufferIsFull(t *testing.T) {
	t.Parallel()

	c := newConnection(nil)

	for range sendBuffer {
		c.Send(ws.Message{Type: "xp.gained"})
	}
	require.Len(t, c.out, sendBuffer)

	c.Send(ws.Message{Type: "overflow"})

	assert.Len(t, c.out, sendBuffer)
	for range sendBuffer {
		assert.Equal(t, "xp.gained", (<-c.out).Type)
	}
}

func TestConnectionPumpStopsOnContextCancel(t *testing.T) {
	t.Parallel()

	c := newConnection(nil)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		c.pump(ctx)
		close(done)
	}()

	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("pump did not stop after context cancellation")
	}
}

func TestConnectionPumpStopsOnShutdown(t *testing.T) {
	t.Parallel()

	c := newConnection(nil)

	done := make(chan struct{})
	go func() {
		c.pump(context.Background())
		close(done)
	}()

	c.shutdown()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("pump did not stop after shutdown")
	}
}

func TestConnectionWriteFailureShutsDownPump(t *testing.T) {
	t.Parallel()

	conn := closedConnection(t)
	c := newConnection(conn)

	done := make(chan struct{})
	go func() {
		c.pump(context.Background())
		close(done)
	}()

	c.Send(ws.Message{Type: "pet.updated"})

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("pump did not stop after a write failure")
	}

	select {
	case <-c.closed:
	default:
		t.Fatal("expected the connection to be marked closed")
	}
}

func TestConnectionWriteReportsFailure(t *testing.T) {
	t.Parallel()

	c := newConnection(closedConnection(t))

	assert.False(t, c.write(context.Background(), serverMessage{Type: "pet.updated"}))
}

func TestConnectionPingReportsFailure(t *testing.T) {
	t.Parallel()

	c := newConnection(closedConnection(t))

	assert.False(t, c.ping(context.Background()))

	select {
	case <-c.closed:
	default:
		t.Fatal("expected a failed ping to close the connection")
	}
}
