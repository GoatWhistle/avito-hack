package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/ws"
)

const (
	messageGetState = "pet.get"
	messagePet      = "pet.pet"
	messageState    = "pet.state"
	messagePing     = "ping"
	messagePong     = "pong"
	messageError    = "error"
	readLimit       = 4096
)

type tokenParser interface {
	Parse(raw string) (auth.Actor, error)
}

type petService interface {
	State(ctx context.Context, userID uuid.UUID) (*domain.Pet, error)
	Stroke(ctx context.Context, userID uuid.UUID) (*domain.Pet, error)
}

type connectionHub interface {
	Register(userID uuid.UUID, sink ws.Sink) func()
}

type WebSocketHandler struct {
	service        petService
	tokens         tokenParser
	hub            connectionHub
	originPatterns []string
}

func NewWebSocketHandler(
	service petService,
	tokens tokenParser,
	hub connectionHub,
	originPatterns []string,
) *WebSocketHandler {
	return &WebSocketHandler{service: service, tokens: tokens, hub: hub, originPatterns: originPatterns}
}

func (h *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	actor, err := h.authenticate(r)
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{OriginPatterns: h.originPatterns})
	if err != nil {
		return
	}
	//nolint:errcheck // closing an already-closed connection is expected here
	defer func() { _ = conn.Close(websocket.StatusNormalClosure, "connection closed") }()
	conn.SetReadLimit(readLimit)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	client := newConnection(conn)
	defer client.shutdown()

	if h.hub != nil {
		release := h.hub.Register(actor.ID, client)
		defer release()
	}

	go client.pump(ctx)

	h.serve(ctx, client, actor)
}

func (h *WebSocketHandler) authenticate(r *http.Request) (auth.Actor, error) {
	raw := r.URL.Query().Get("token")
	if raw == "" {
		return auth.Actor{}, domainerr.ErrUnauthorized
	}
	actor, err := h.tokens.Parse(raw)
	if err != nil {
		return auth.Actor{}, domainerr.ErrUnauthorized
	}

	return actor, nil
}

func (h *WebSocketHandler) serve(ctx context.Context, client *connection, actor auth.Actor) {
	for {
		message, err := readMessage(ctx, client.conn)
		if err != nil {
			if websocket.CloseStatus(err) == -1 && !errors.Is(err, context.Canceled) {
				//nolint:errcheck // the read already failed, the close status is best-effort
				_ = client.conn.Close(websocket.StatusInvalidFramePayloadData, "invalid message")
			}

			return
		}

		select {
		case <-client.closed:
			return
		default:
		}

		h.handle(ctx, client, actor, message)
	}
}

func readMessage(ctx context.Context, conn *websocket.Conn) (clientMessage, error) {
	readCtx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	var message clientMessage
	if err := wsjson.Read(readCtx, conn, &message); err != nil {
		return clientMessage{}, err
	}

	return message, nil
}

func (h *WebSocketHandler) handle(
	ctx context.Context,
	client *connection,
	actor auth.Actor,
	message clientMessage,
) {
	var (
		pet *domain.Pet
		err error
	)

	switch message.Type {
	case messagePing:
		client.Reply(message.RequestID, messagePong, nil)

		return
	case messageGetState:
		pet, err = h.service.State(ctx, actor.ID)
	case messagePet:
		pet, err = h.service.Stroke(ctx, actor.ID)
	default:
		client.Reply(message.RequestID, messageError,
			errorPayload{Code: "unknown_message", Message: "unsupported message type"})

		return
	}

	if err != nil {
		client.Reply(message.RequestID, messageError,
			errorPayload{Code: "internal_error", Message: "failed to update pet"})

		return
	}

	client.Reply(message.RequestID, messageState, toPetPayload(pet))
}
