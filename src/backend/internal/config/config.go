package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/caarlos0/env/v11"

	"github.com/avito-hack/backend/internal/shared/postgres"
)

type Config struct {
	HTTPAddr    string        `env:"HTTP_ADDR"       envDefault:":8080"`
	DatabaseURL string        `env:"DATABASE_URL,required"`
	RedisAddr   string        `env:"REDIS_ADDR"      envDefault:"localhost:6379"`
	JWTSecret   string        `env:"JWT_SECRET,required"`
	JWTTTL      time.Duration `env:"JWT_TTL"         envDefault:"1h"`

	RewardHMACSecret string `env:"REWARD_HMAC_SECRET,required"`

	AllowedOrigins []string `env:"ALLOWED_ORIGINS" envSeparator:","`
	LogLevel       string   `env:"LOG_LEVEL"       envDefault:"info"`
	LogFormat      string   `env:"LOG_FORMAT"      envDefault:"pretty"`
	LogColor       string   `env:"LOG_COLOR"       envDefault:"auto"`

	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"5s"`
	ReadTimeout       time.Duration `env:"READ_TIMEOUT"        envDefault:"15s"`
	WriteTimeout      time.Duration `env:"WRITE_TIMEOUT"       envDefault:"30s"`
	IdleTimeout       time.Duration `env:"IDLE_TIMEOUT"        envDefault:"60s"`
	ShutdownTimeout   time.Duration `env:"SHUTDOWN_TIMEOUT"    envDefault:"15s"`
	RequestTimeout    time.Duration `env:"REQUEST_TIMEOUT"     envDefault:"30s"`

	RateLimitEnabled    bool    `env:"RATE_LIMIT_ENABLED"     envDefault:"true"`
	RateLimitRPS        float64 `env:"RATE_LIMIT_RPS"         envDefault:"20"`
	RateLimitBurst      float64 `env:"RATE_LIMIT_BURST"       envDefault:"40"`
	AuthRateLimitRPS    float64 `env:"AUTH_RATE_LIMIT_RPS"    envDefault:"5"`
	AuthRateLimitBurst  float64 `env:"AUTH_RATE_LIMIT_BURST"  envDefault:"10"`
	RateLimitTrustProxy bool    `env:"RATE_LIMIT_TRUST_PROXY" envDefault:"false"`

	MaxBodyBytes  int64  `env:"MAX_BODY_BYTES"  envDefault:"1048576"`
	MaxPhotoBytes int64  `env:"MAX_PHOTO_BYTES" envDefault:"5242880"`
	UploadDir     string `env:"UPLOAD_DIR"      envDefault:"/data/uploads"`
	UploadURL     string `env:"UPLOAD_URL"      envDefault:"/uploads"`

	DBMaxConns          int32         `env:"DB_MAX_CONNS"            envDefault:"20"`
	DBMinConns          int32         `env:"DB_MIN_CONNS"            envDefault:"2"`
	DBMaxConnLifetime   time.Duration `env:"DB_MAX_CONN_LIFETIME"    envDefault:"1h"`
	DBMaxConnIdleTime   time.Duration `env:"DB_MAX_CONN_IDLE_TIME"   envDefault:"30m"`
	DBHealthCheckPeriod time.Duration `env:"DB_HEALTH_CHECK_PERIOD"  envDefault:"1m"`
	DBConnectTimeout    time.Duration `env:"DB_CONNECT_TIMEOUT"      envDefault:"5s"`

	KafkaBrokers  []string      `env:"KAFKA_BROKERS"   envSeparator:","`
	KafkaTopic    string        `env:"KAFKA_TOPIC"     envDefault:"pet.activity"`
	KafkaGroup    string        `env:"KAFKA_GROUP"     envDefault:"pet-service"`
	KafkaClientID string        `env:"KAFKA_CLIENT_ID" envDefault:"pet-service"`
	KafkaTimeout  time.Duration `env:"KAFKA_TIMEOUT"   envDefault:"3s"`

	PetFlushInterval  time.Duration `env:"PET_FLUSH_INTERVAL"   envDefault:"15s"`
	PetFlushBatchSize int64         `env:"PET_FLUSH_BATCH_SIZE" envDefault:"100"`

	OpenRouterAPIKey string `env:"OPENROUTER_API_KEY"`
	OpenRouterModel  string `env:"OPENROUTER_MODEL" envDefault:"google/gemini-2.0-flash-exp:free"`

	MetricsToken string `env:"METRICS_TOKEN"`
}

func (c Config) KafkaEnabled() bool {
	return len(c.KafkaBrokers) > 0 && c.KafkaTopic != ""
}

func (c Config) PoolOptions() postgres.PoolOptions {
	return postgres.PoolOptions{
		MaxConns:          c.DBMaxConns,
		MinConns:          c.DBMinConns,
		MaxConnLifetime:   c.DBMaxConnLifetime,
		MaxConnIdleTime:   c.DBMaxConnIdleTime,
		HealthCheckPeriod: c.DBHealthCheckPeriod,
		ConnectTimeout:    c.DBConnectTimeout,
	}
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
