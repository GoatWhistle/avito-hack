package config

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

const (
	minJWTSecretLen    = 16
	minRewardSecretLen = 16
	maxJWTTTL          = 720 * time.Hour
	maxTimeout         = 10 * time.Minute
	maxDBConns         = 1000
	maxBodyBytesLimit  = 1 << 30
)

var (
	validLogLevels  = []string{"debug", "info", "warn", "error"}
	validLogFormats = []string{"pretty", "json", "text"}
	validLogColors  = []string{"auto", "always", "never"}
)

func (c Config) validate() error {
	validators := []func() error{
		c.validateSecrets,
		c.validateAddresses,
		c.validateLogging,
		c.validateTimeouts,
		c.validateLimits,
		c.validatePool,
		c.validateKafka,
		c.validatePet,
	}

	for _, validate := range validators {
		if err := validate(); err != nil {
			return err
		}
	}

	return nil
}

func (c Config) validateSecrets() error {
	if strings.TrimSpace(c.DatabaseURL) == "" {
		return fmt.Errorf("DATABASE_URL must not be empty")
	}
	if err := validateDatabaseURL(c.DatabaseURL); err != nil {
		return err
	}
	if len(c.JWTSecret) < minJWTSecretLen {
		return fmt.Errorf("JWT_SECRET must be at least %d characters long, got %d",
			minJWTSecretLen, len(c.JWTSecret))
	}
	if len(c.RewardHMACSecret) < minRewardSecretLen {
		return fmt.Errorf("REWARD_HMAC_SECRET must be at least %d characters long, got %d",
			minRewardSecretLen, len(c.RewardHMACSecret))
	}
	if c.JWTTTL <= 0 {
		return fmt.Errorf("JWT_TTL must be positive, got %s", c.JWTTTL)
	}
	if c.JWTTTL > maxJWTTTL {
		return fmt.Errorf("JWT_TTL must not exceed %s, got %s", maxJWTTTL, c.JWTTTL)
	}

	return nil
}

func validateDatabaseURL(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("DATABASE_URL is not a valid URL: %w", err)
	}
	if parsed.Scheme != "postgres" && parsed.Scheme != "postgresql" {
		return fmt.Errorf(
			"DATABASE_URL must use postgres:// or postgresql:// scheme, got %q", parsed.Scheme)
	}
	if parsed.Host == "" {
		return fmt.Errorf("DATABASE_URL must contain a host")
	}

	return nil
}

func (c Config) validateAddresses() error {
	if err := validateHostPort("HTTP_ADDR", c.HTTPAddr, true); err != nil {
		return err
	}

	return validateHostPort("REDIS_ADDR", c.RedisAddr, false)
}

func validateHostPort(name, value string, allowEmptyHost bool) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s must not be empty", name)
	}

	host, port, err := net.SplitHostPort(value)
	if err != nil {
		return fmt.Errorf("%s must be in host:port form, got %q", name, value)
	}
	if host == "" && !allowEmptyHost {
		return fmt.Errorf("%s must contain a host, got %q", name, value)
	}
	if port == "" {
		return fmt.Errorf("%s must contain a port, got %q", name, value)
	}

	return nil
}

func (c Config) validateLogging() error {
	if err := oneOf("LOG_LEVEL", c.LogLevel, validLogLevels); err != nil {
		return err
	}
	if err := oneOf("LOG_FORMAT", c.LogFormat, validLogFormats); err != nil {
		return err
	}

	return oneOf("LOG_COLOR", c.LogColor, validLogColors)
}

func oneOf(name, value string, allowed []string) error {
	for _, candidate := range allowed {
		if value == candidate {
			return nil
		}
	}

	return fmt.Errorf("%s must be one of [%s], got %q", name, strings.Join(allowed, ", "), value)
}
