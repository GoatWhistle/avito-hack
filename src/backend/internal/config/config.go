package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPAddr       string        `env:"HTTP_ADDR"       envDefault:":8080"`
	DatabaseURL    string        `env:"DATABASE_URL,required"`
	RedisAddr      string        `env:"REDIS_ADDR"      envDefault:"localhost:6379"`
	JWTSecret      string        `env:"JWT_SECRET,required"`
	JWTTTL         time.Duration `env:"JWT_TTL"         envDefault:"15m"`
	RefreshTTL     time.Duration `env:"REFRESH_TTL"     envDefault:"720h"`
	AllowedOrigins []string      `env:"ALLOWED_ORIGINS" envSeparator:","`
	LogLevel       string        `env:"LOG_LEVEL"       envDefault:"info"`
	LogFormat      string        `env:"LOG_FORMAT"      envDefault:"pretty"`
	LogColor       string        `env:"LOG_COLOR"       envDefault:"auto"`

	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"5s"`
	ReadTimeout       time.Duration `env:"READ_TIMEOUT"        envDefault:"15s"`
	WriteTimeout      time.Duration `env:"WRITE_TIMEOUT"       envDefault:"30s"`
	IdleTimeout       time.Duration `env:"IDLE_TIMEOUT"        envDefault:"60s"`
	ShutdownTimeout   time.Duration `env:"SHUTDOWN_TIMEOUT"    envDefault:"15s"`
	RequestTimeout    time.Duration `env:"REQUEST_TIMEOUT"     envDefault:"30s"`

	MaxBodyBytes int64 `env:"MAX_BODY_BYTES" envDefault:"1048576"`
}

func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse env: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) SlogLevel() slog.Level {
	switch c.LogLevel {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

const minJWTSecretLen = 16

func (c Config) validate() error {
	if len(c.JWTSecret) < minJWTSecretLen {
		return fmt.Errorf("JWT_SECRET must be at least %d characters long", minJWTSecretLen)
	}
	if c.JWTTTL <= 0 {
		return fmt.Errorf("JWT_TTL must be positive, got %s", c.JWTTTL)
	}
	if c.MaxBodyBytes <= 0 {
		return fmt.Errorf("MAX_BODY_BYTES must be positive, got %d", c.MaxBodyBytes)
	}

	return nil
}
