//go:build kafka

package kafka

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Consumer struct {
	client  *kgo.Client
	handler Handler
	log     *slog.Logger
}

func NewConsumer(cfg Config, handler Handler, log *slog.Logger) (*Consumer, error) {
	if !cfg.ValidConsumer() || handler == nil {
		return nil, ErrInvalidConfig
	}
	if log == nil {
		log = slog.Default()
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ClientID(cfg.clientID()),
		kgo.ConsumerGroup(cfg.ConsumerGroup),
		kgo.ConsumeTopics(cfg.Topic),
		kgo.DisableAutoCommit(),
	)
	if err != nil {
		return nil, fmt.Errorf("create kafka consumer: %w", err)
	}

	return &Consumer{client: client, handler: handler, log: log}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		fetches := c.client.PollFetches(ctx)
		if ctx.Err() != nil {
			return nil
		}

		c.logFetchErrors(ctx, fetches)

		var records []*kgo.Record
		fetches.EachRecord(func(record *kgo.Record) {
			if c.consumeRecord(ctx, record) {
				records = append(records, record)
			}
		})

		if len(records) == 0 {
			continue
		}
		if err := c.client.CommitRecords(ctx, records...); err != nil && ctx.Err() == nil {
			c.log.WarnContext(ctx, "commit kafka records", slog.Any("error", err))
		}
	}
}

func (c *Consumer) consumeRecord(ctx context.Context, record *kgo.Record) bool {
	m, err := decodeMessage(record.Value)
	if err != nil {
		c.log.WarnContext(ctx, "skip malformed kafka record",
			slog.Int64("offset", record.Offset), slog.Any("error", err))

		return true
	}

	if err := c.handler.Handle(ctx, toEvent(m), m.ID.String()); err != nil {
		if errors.Is(err, context.Canceled) {
			return false
		}
		c.log.ErrorContext(ctx, "handle kafka record",
			slog.String("event", m.Type),
			slog.String("event_id", m.ID.String()),
			slog.Any("error", err),
		)

		return false
	}

	return true
}

func (c *Consumer) logFetchErrors(ctx context.Context, fetches kgo.Fetches) {
	for _, fetchErr := range fetches.Errors() {
		if errors.Is(fetchErr.Err, context.Canceled) {
			continue
		}
		c.log.WarnContext(ctx, "fetch kafka records",
			slog.String("topic", fetchErr.Topic),
			slog.Int("partition", int(fetchErr.Partition)),
			slog.Any("error", fetchErr.Err),
		)
	}
}

func (c *Consumer) Close() {
	c.client.Close()
}
