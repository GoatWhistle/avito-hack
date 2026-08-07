package config

import (
	"fmt"
	"strings"
	"time"
)

func (c Config) validateTimeouts() error {
	timeouts := []struct {
		name  string
		value time.Duration
	}{
		{name: "READ_HEADER_TIMEOUT", value: c.ReadHeaderTimeout},
		{name: "READ_TIMEOUT", value: c.ReadTimeout},
		{name: "WRITE_TIMEOUT", value: c.WriteTimeout},
		{name: "IDLE_TIMEOUT", value: c.IdleTimeout},
		{name: "SHUTDOWN_TIMEOUT", value: c.ShutdownTimeout},
		{name: "REQUEST_TIMEOUT", value: c.RequestTimeout},
	}

	for _, timeout := range timeouts {
		if err := positiveDuration(timeout.name, timeout.value, maxTimeout); err != nil {
			return err
		}
	}

	if c.ReadHeaderTimeout > c.ReadTimeout {
		return fmt.Errorf("READ_HEADER_TIMEOUT (%s) must not exceed READ_TIMEOUT (%s)",
			c.ReadHeaderTimeout, c.ReadTimeout)
	}

	return nil
}

func positiveDuration(name string, value, limit time.Duration) error {
	if value <= 0 {
		return fmt.Errorf("%s must be positive, got %s", name, value)
	}
	if limit > 0 && value > limit {
		return fmt.Errorf("%s must not exceed %s, got %s", name, limit, value)
	}

	return nil
}

func (c Config) validateLimits() error {
	if c.MaxBodyBytes <= 0 {
		return fmt.Errorf("MAX_BODY_BYTES must be positive, got %d", c.MaxBodyBytes)
	}
	if c.MaxBodyBytes > maxBodyBytesLimit {
		return fmt.Errorf("MAX_BODY_BYTES must not exceed %d, got %d",
			int64(maxBodyBytesLimit), c.MaxBodyBytes)
	}
	if c.MaxPhotoBytes <= 0 {
		return fmt.Errorf("MAX_PHOTO_BYTES must be positive, got %d", c.MaxPhotoBytes)
	}
	if c.MaxPhotoBytes > maxBodyBytesLimit {
		return fmt.Errorf("MAX_PHOTO_BYTES must not exceed %d, got %d",
			int64(maxBodyBytesLimit), c.MaxPhotoBytes)
	}
	if strings.TrimSpace(c.UploadDir) == "" {
		return fmt.Errorf("UPLOAD_DIR must not be empty")
	}
	if !strings.HasPrefix(c.UploadURL, "/") {
		return fmt.Errorf("UPLOAD_URL must start with '/', got %q", c.UploadURL)
	}

	return nil
}

func (c Config) validatePool() error {
	if c.DBMaxConns <= 0 {
		return fmt.Errorf("DB_MAX_CONNS must be positive, got %d", c.DBMaxConns)
	}
	if c.DBMaxConns > maxDBConns {
		return fmt.Errorf("DB_MAX_CONNS must not exceed %d, got %d", maxDBConns, c.DBMaxConns)
	}
	if c.DBMinConns < 0 {
		return fmt.Errorf("DB_MIN_CONNS must not be negative, got %d", c.DBMinConns)
	}
	if c.DBMinConns > c.DBMaxConns {
		return fmt.Errorf("DB_MIN_CONNS (%d) must not exceed DB_MAX_CONNS (%d)",
			c.DBMinConns, c.DBMaxConns)
	}

	durations := []struct {
		name  string
		value time.Duration
	}{
		{name: "DB_MAX_CONN_LIFETIME", value: c.DBMaxConnLifetime},
		{name: "DB_MAX_CONN_IDLE_TIME", value: c.DBMaxConnIdleTime},
		{name: "DB_HEALTH_CHECK_PERIOD", value: c.DBHealthCheckPeriod},
	}
	for _, duration := range durations {
		if err := positiveDuration(duration.name, duration.value, 0); err != nil {
			return err
		}
	}

	return positiveDuration("DB_CONNECT_TIMEOUT", c.DBConnectTimeout, maxTimeout)
}

func (c Config) validateKafka() error {
	if err := positiveDuration("KAFKA_TIMEOUT", c.KafkaTimeout, maxTimeout); err != nil {
		return err
	}
	if !c.KafkaEnabled() {
		return nil
	}

	for _, broker := range c.KafkaBrokers {
		if err := validateHostPort("KAFKA_BROKERS", strings.TrimSpace(broker), false); err != nil {
			return err
		}
	}
	if strings.TrimSpace(c.KafkaGroup) == "" {
		return fmt.Errorf("KAFKA_GROUP must not be empty when Kafka is enabled")
	}
	if strings.TrimSpace(c.KafkaClientID) == "" {
		return fmt.Errorf("KAFKA_CLIENT_ID must not be empty when Kafka is enabled")
	}

	return nil
}

func (c Config) validatePet() error {
	if err := positiveDuration("PET_FLUSH_INTERVAL", c.PetFlushInterval, 0); err != nil {
		return err
	}
	if c.PetFlushBatchSize <= 0 {
		return fmt.Errorf("PET_FLUSH_BATCH_SIZE must be positive, got %d", c.PetFlushBatchSize)
	}

	return nil
}
