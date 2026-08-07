//go:build !kafka

package kafka_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/events"
	"github.com/avito-hack/backend/internal/shared/kafka"
)

type noopHandler struct{}

func (noopHandler) Handle(context.Context, events.Event, string) error { return nil }

func validConfig() kafka.Config {
	return kafka.Config{
		Brokers:       []string{"127.0.0.1:9092"},
		Topic:         "pet-events",
		ConsumerGroup: "pet-service",
	}
}

func TestStubNewProducerDisabled(t *testing.T) {
	t.Parallel()

	producer, err := kafka.NewProducer(validConfig())
	require.ErrorIs(t, err, kafka.ErrDisabled)
	assert.Nil(t, producer)
}

func TestStubProducerMethodsDisabled(t *testing.T) {
	t.Parallel()

	producer := &kafka.Producer{}
	event := events.Event{
		Type:       events.Type("pet.stroked"),
		UserID:     uuid.MustParse("44444444-4444-4444-4444-444444444444"),
		OccurredAt: time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC),
	}

	require.ErrorIs(t, producer.Publish(context.Background(), event), kafka.ErrDisabled)
	require.ErrorIs(t, producer.Ping(context.Background()), kafka.ErrDisabled)
	assert.NotPanics(t, producer.Close)
}

func TestStubNewConsumerDisabled(t *testing.T) {
	t.Parallel()

	consumer, err := kafka.NewConsumer(validConfig(), noopHandler{}, slog.Default())
	require.ErrorIs(t, err, kafka.ErrDisabled)
	assert.Nil(t, consumer)
}

func TestStubNewConsumerDisabledWithNilDeps(t *testing.T) {
	t.Parallel()

	consumer, err := kafka.NewConsumer(kafka.Config{}, nil, nil)
	require.ErrorIs(t, err, kafka.ErrDisabled)
	assert.Nil(t, consumer)
}

func TestStubConsumerRunReturnsOnContextDone(t *testing.T) {
	t.Parallel()

	consumer := &kafka.Consumer{}
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() { done <- consumer.Run(ctx) }()

	cancel()

	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("consumer run did not return after context cancellation")
	}

	assert.NotPanics(t, consumer.Close)
}

func TestStubConsumerRunWithAlreadyCancelledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	require.NoError(t, (&kafka.Consumer{}).Run(ctx))
}

func TestErrorsAreDistinct(t *testing.T) {
	t.Parallel()

	require.NotErrorIs(t, kafka.ErrDisabled, kafka.ErrInvalidConfig)
	require.NotErrorIs(t, kafka.ErrDisabled, kafka.ErrInvalidMessage)
	require.NotErrorIs(t, kafka.ErrInvalidConfig, kafka.ErrInvalidMessage)
	require.EqualError(t, kafka.ErrDisabled, "kafka support is not built into this binary")
	require.EqualError(t, kafka.ErrInvalidConfig, "invalid kafka config")
	require.EqualError(t, kafka.ErrInvalidMessage, "invalid kafka message")
}
