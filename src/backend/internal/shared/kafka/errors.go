package kafka

import (
	"context"
	"errors"

	"github.com/avito-hack/backend/internal/shared/events"
)

var (
	ErrInvalidConfig  = errors.New("invalid kafka config")
	ErrInvalidMessage = errors.New("invalid kafka message")
	ErrDisabled       = errors.New("kafka support is not built into this binary")
)

type Handler interface {
	Handle(ctx context.Context, e events.Event, id string) error
}
