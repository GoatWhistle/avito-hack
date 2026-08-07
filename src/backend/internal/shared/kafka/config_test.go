package kafka

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigClientID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		config Config
		want   string
	}{
		{name: "empty falls back to default", config: Config{}, want: defaultClientID},
		{name: "explicit value wins", config: Config{ClientID: "custom"}, want: "custom"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.config.clientID())
		})
	}
}

func TestConfigTimeout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		config Config
		want   time.Duration
	}{
		{name: "zero falls back to default", config: Config{}, want: defaultTimeout},
		{name: "negative falls back to default", config: Config{Timeout: -time.Second}, want: defaultTimeout},
		{name: "positive value wins", config: Config{Timeout: 9 * time.Second}, want: 9 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.config.timeout())
		})
	}
}

func TestConfigValidProducer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		config Config
		want   bool
	}{
		{name: "brokers and topic", config: Config{Brokers: []string{"b:9092"}, Topic: "events"}, want: true},
		{name: "no brokers", config: Config{Topic: "events"}, want: false},
		{name: "empty broker slice", config: Config{Brokers: []string{}, Topic: "events"}, want: false},
		{name: "no topic", config: Config{Brokers: []string{"b:9092"}}, want: false},
		{name: "empty config", config: Config{}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.config.ValidProducer())
		})
	}
}

func TestConfigValidConsumer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		config Config
		want   bool
	}{
		{
			name:   "full config",
			config: Config{Brokers: []string{"b:9092"}, Topic: "events", ConsumerGroup: "pets"},
			want:   true,
		},
		{
			name:   "missing consumer group",
			config: Config{Brokers: []string{"b:9092"}, Topic: "events"},
			want:   false,
		},
		{
			name:   "consumer group without topic",
			config: Config{Brokers: []string{"b:9092"}, ConsumerGroup: "pets"},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.config.ValidConsumer())
		})
	}
}
