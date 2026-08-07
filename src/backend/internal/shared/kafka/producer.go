//go:build kafka

package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/avito-hack/backend/internal/shared/events"
)

type Producer struct {
	client  *kgo.Client
	topic   string
	timeout time.Duration
}

func NewProducer(cfg Config) (*Producer, error) {
	if !cfg.ValidProducer() {
		return nil, ErrInvalidConfig
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ClientID(cfg.clientID()),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.AllowAutoTopicCreation(),
		kgo.ProduceRequestTimeout(cfg.timeout()),
		kgo.RecordDeliveryTimeout(cfg.timeout()),
	)
	if err != nil {
		return nil, fmt.Errorf("create kafka producer: %w", err)
	}

	return &Producer{client: client, topic: cfg.Topic, timeout: cfg.timeout()}, nil
}

func (p *Producer) Publish(ctx context.Context, e events.Event) error {
	m := toMessage(e)
	if err := validateMessage(m); err != nil {
		return err
	}

	raw, err := encodeMessage(m)
	if err != nil {
		return err
	}

	record := &kgo.Record{
		Topic: p.topic,
		Key:   []byte(e.UserID.String()),
		Value: raw,
		Headers: []kgo.RecordHeader{
			{Key: "content-type", Value: []byte("application/json")},
			{Key: "event-type", Value: []byte(e.Type)},
		},
	}

	produceCtx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	if err := p.client.ProduceSync(produceCtx, record).FirstErr(); err != nil {
		return fmt.Errorf("publish kafka event: %w", err)
	}

	return nil
}

func (p *Producer) Ping(ctx context.Context) error {
	pingCtx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	if err := p.client.Ping(pingCtx); err != nil {
		return fmt.Errorf("ping kafka: %w", err)
	}

	return nil
}

func (p *Producer) Close() {
	p.client.Close()
}
