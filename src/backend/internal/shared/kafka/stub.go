//go:build !kafka

package kafka

import (
	"context"
	"log/slog"

	"github.com/avito-hack/backend/internal/shared/events"
)

type Producer struct{}

func NewProducer(Config) (*Producer, error) { return nil, ErrDisabled }

func (p *Producer) Publish(context.Context, events.Event) error { return ErrDisabled }

func (p *Producer) Ping(context.Context) error { return ErrDisabled }

func (p *Producer) Close() {}

type Consumer struct{}

func NewConsumer(Config, Handler, *slog.Logger) (*Consumer, error) { return nil, ErrDisabled }

func (c *Consumer) Run(ctx context.Context) error {
	<-ctx.Done()

	return nil
}

func (c *Consumer) Close() {}
