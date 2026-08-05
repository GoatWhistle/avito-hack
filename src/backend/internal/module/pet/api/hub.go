package api

import (
	"context"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"

	"github.com/avito-hack/backend/internal/shared/ws"
)

const (
	sendBuffer   = 32
	writeTimeout = 5 * time.Second
	pingInterval = 30 * time.Second
	readTimeout  = 90 * time.Second
)

type connection struct {
	conn   *websocket.Conn
	out    chan serverMessage
	closed chan struct{}
}

func newConnection(conn *websocket.Conn) *connection {
	return &connection{
		conn:   conn,
		out:    make(chan serverMessage, sendBuffer),
		closed: make(chan struct{}),
	}
}

func (c *connection) Send(m ws.Message) {
	c.enqueue(serverMessage{Type: m.Type, Payload: m.Payload})
}

func (c *connection) Reply(requestID, messageType string, payload any) {
	c.enqueue(serverMessage{Type: messageType, RequestID: requestID, Payload: payload})
}

func (c *connection) enqueue(m serverMessage) {
	select {
	case <-c.closed:
	case c.out <- m:
	default:
	}
}

func (c *connection) shutdown() {
	select {
	case <-c.closed:
	default:
		close(c.closed)
	}
}

func (c *connection) pump(ctx context.Context) {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.closed:
			return
		case m := <-c.out:
			if !c.write(ctx, m) {
				return
			}
		case <-ticker.C:
			if !c.ping(ctx) {
				return
			}
		}
	}
}

func (c *connection) write(ctx context.Context, m serverMessage) bool {
	writeCtx, cancel := context.WithTimeout(ctx, writeTimeout)
	defer cancel()

	if err := wsjson.Write(writeCtx, c.conn, m); err != nil {
		_ = c.conn.Close(websocket.StatusInternalError, "write failed") //nolint:errcheck // the connection is already broken
		c.shutdown()

		return false
	}

	return true
}

func (c *connection) ping(ctx context.Context) bool {
	pingCtx, cancel := context.WithTimeout(ctx, writeTimeout)
	defer cancel()

	if err := c.conn.Ping(pingCtx); err != nil {
		_ = c.conn.Close(websocket.StatusGoingAway, "heartbeat timeout") //nolint:errcheck // the peer is already gone
		c.shutdown()

		return false
	}

	return true
}
