package kafka

import "time"

const (
	defaultClientID = "pet-service"
	defaultTimeout  = 3 * time.Second
)

type Config struct {
	Brokers       []string
	Topic         string
	ConsumerGroup string
	ClientID      string
	Timeout       time.Duration
}

func (c Config) clientID() string {
	if c.ClientID == "" {
		return defaultClientID
	}

	return c.ClientID
}

func (c Config) timeout() time.Duration {
	if c.Timeout <= 0 {
		return defaultTimeout
	}

	return c.Timeout
}

func (c Config) ValidProducer() bool {
	return len(c.Brokers) > 0 && c.Topic != ""
}

func (c Config) ValidConsumer() bool {
	return c.ValidProducer() && c.ConsumerGroup != ""
}
