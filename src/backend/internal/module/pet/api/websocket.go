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
	Pet(ctx context.Context, userID uuid.UUID) (*domain.Pet, error)
}

type WebSocketHandler struct {
	service        petService
	tokens         tokenParser
	originPatterns []string
}

func NewWebSocketHandler(service petService, tokens tokenParser, originPatterns []string) *WebSocketHandler {
	return &WebSocketHandler{service: service, tokens: tokens, originPatterns: originPatterns}
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
	defer func() { _ = conn.Close(websocket.StatusNormalClosure, "connection closed") }()
	conn.SetReadLimit(readLimit)

	h.serve(r.Context(), conn, actor)
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

func (h *WebSocketHandler) serve(ctx context.Context, conn *websocket.Conn, actor auth.Actor) {
	for {
		var message clientMessage
		if err := wsjson.Read(ctx, conn, &message); err != nil {
			if websocket.CloseStatus(err) == -1 && !errors.Is(err, context.Canceled) {
				_ = conn.Close(websocket.StatusInvalidFramePayloadData, "invalid message")
			}
			return
		}

		if err := h.handle(ctx, conn, actor, message); err != nil {
			return
		}
	}
}

func (h *WebSocketHandler) handle(
	ctx context.Context,
	conn *websocket.Conn,
	actor auth.Actor,
	message clientMessage,
) error {
	var (
		pet *domain.Pet
		err error
	)

	switch message.Type {
	case messagePing:
		return wsjson.Write(ctx, conn, serverMessage{Type: messagePong, RequestID: message.RequestID})
	case messageGetState:
		pet, err = h.service.State(ctx, actor.ID)
	case messagePet:
		pet, err = h.service.Pet(ctx, actor.ID)
	default:
		return wsjson.Write(ctx, conn, serverMessage{
			Type: messageError, RequestID: message.RequestID,
			Payload: errorPayload{Code: "unknown_message", Message: "unsupported message type"},
		})
	}
	if err != nil {
		return wsjson.Write(ctx, conn, serverMessage{
			Type: messageError, RequestID: message.RequestID,
			Payload: errorPayload{Code: "internal_error", Message: "failed to update pet"},
		})
	}

	return wsjson.Write(ctx, conn, serverMessage{
		Type: messageState, RequestID: message.RequestID, Payload: toPetPayload(pet),
	})
}
