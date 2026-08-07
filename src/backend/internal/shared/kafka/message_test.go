//go:build kafka

package kafka

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/events"
)

func sampleEvent() events.Event {
	return events.New(
		events.TypeItemPublished, uuid.New(), uuid.New(),
		time.Now().UTC().Truncate(time.Microsecond),
	).WithPayload(events.Payload{Title: "Chair", PriceKopeks: 150000})
}

func TestMessageRoundTripPreservesEvent(t *testing.T) {
	t.Parallel()

	original := sampleEvent()

	raw, err := encodeMessage(toMessage(original))
	require.NoError(t, err)

	decoded, err := decodeMessage(raw)
	require.NoError(t, err)

	restored := toEvent(decoded)
	assert.Equal(t, original.Type, restored.Type)
	assert.Equal(t, original.UserID, restored.UserID)
	assert.Equal(t, original.SubjectID, restored.SubjectID)
	assert.Equal(t, original.Payload.Title, restored.Payload.Title)
	assert.Equal(t, original.Payload.PriceKopeks, restored.Payload.PriceKopeks)
	assert.True(t, original.OccurredAt.Equal(restored.OccurredAt))
}

func TestEventIDIsDeterministic(t *testing.T) {
	t.Parallel()

	event := sampleEvent()

	assert.Equal(t, EventID(event), EventID(event),
		"redelivery of the same event must produce the same dedup id")
}

func TestEventIDDiffersPerSubject(t *testing.T) {
	t.Parallel()

	first := sampleEvent()
	second := first
	second.SubjectID = uuid.New()

	assert.NotEqual(t, EventID(first), EventID(second))
}

func TestDecodeRejectsGarbage(t *testing.T) {
	t.Parallel()

	_, err := decodeMessage([]byte("{not json"))
	require.Error(t, err)
}

func TestDecodeRejectsUnsupportedVersion(t *testing.T) {
	t.Parallel()

	m := toMessage(sampleEvent())
	m.Version = 99

	raw, err := encodeMessage(m)
	require.NoError(t, err)

	_, err = decodeMessage(raw)
	require.ErrorIs(t, err, ErrInvalidMessage)
}

func TestDecodeRejectsMissingUser(t *testing.T) {
	t.Parallel()

	m := toMessage(sampleEvent())
	m.UserID = uuid.Nil

	raw, err := encodeMessage(m)
	require.NoError(t, err)

	_, err = decodeMessage(raw)
	require.ErrorIs(t, err, ErrInvalidMessage)
}

func TestConfigValidation(t *testing.T) {
	t.Parallel()

	assert.False(t, Config{}.ValidProducer())
	assert.False(t, Config{Brokers: []string{"kafka:9092"}}.ValidProducer())
	assert.True(t, Config{Brokers: []string{"kafka:9092"}, Topic: "pet"}.ValidProducer())

	producer := Config{Brokers: []string{"kafka:9092"}, Topic: "pet"}
	assert.False(t, producer.ValidConsumer())

	producer.ConsumerGroup = "pet-service"
	assert.True(t, producer.ValidConsumer())
}
