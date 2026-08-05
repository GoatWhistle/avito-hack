package config_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/config"
)

const (
	validSecret = "0123456789abcdef0123"
	databaseURL = "postgres://user:pass@localhost:5432/db?sslmode=disable"
)

func setRequired(t *testing.T) {
	t.Helper()

	t.Setenv("DATABASE_URL", databaseURL)
	t.Setenv("JWT_SECRET", validSecret)
	t.Setenv("REWARD_HMAC_SECRET", validSecret)
}

func TestLoadAppliesDefaults(t *testing.T) {
	setRequired(t)

	cfg, err := config.Load()

	require.NoError(t, err)
	assert.Equal(t, ":8080", cfg.HTTPAddr)
	assert.Equal(t, databaseURL, cfg.DatabaseURL)
	assert.Equal(t, "localhost:6379", cfg.RedisAddr)
	assert.Equal(t, 15*time.Minute, cfg.JWTTTL)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "pretty", cfg.LogFormat)
	assert.Equal(t, "auto", cfg.LogColor)
	assert.Equal(t, int64(1048576), cfg.MaxBodyBytes)
	assert.Equal(t, int64(5242880), cfg.MaxPhotoBytes)
	assert.Equal(t, "/data/uploads", cfg.UploadDir)
	assert.Equal(t, 30*time.Second, cfg.RequestTimeout)
}

func TestLoadReadsOverrides(t *testing.T) {
	setRequired(t)
	t.Setenv("HTTP_ADDR", ":9090")
	t.Setenv("JWT_TTL", "1h")
	t.Setenv("ALLOWED_ORIGINS", "https://a.example,https://b.example")
	t.Setenv("MAX_BODY_BYTES", "2048")
	t.Setenv("LOG_LEVEL", "debug")

	cfg, err := config.Load()

	require.NoError(t, err)
	assert.Equal(t, ":9090", cfg.HTTPAddr)
	assert.Equal(t, time.Hour, cfg.JWTTTL)
	assert.Equal(t, []string{"https://a.example", "https://b.example"}, cfg.AllowedOrigins)
	assert.Equal(t, int64(2048), cfg.MaxBodyBytes)
}

func TestLoadRequiresMandatorySecrets(t *testing.T) {
	for _, missing := range []string{"JWT_SECRET", "REWARD_HMAC_SECRET"} {
		t.Run("missing "+missing, func(t *testing.T) {
			setRequired(t)
			t.Setenv(missing, "")

			_, err := config.Load()

			require.Error(t, err)
		})
	}
}

func TestLoadValidatesValues(t *testing.T) {
	tests := []struct {
		name string
		key  string
		val  string
	}{
		{name: "short jwt secret", key: "JWT_SECRET", val: "tooshort"},
		{name: "short reward secret", key: "REWARD_HMAC_SECRET", val: "tooshort"},
		{name: "zero jwt ttl", key: "JWT_TTL", val: "0s"},
		{name: "negative jwt ttl", key: "JWT_TTL", val: "-1m"},
		{name: "zero body limit", key: "MAX_BODY_BYTES", val: "0"},
		{name: "negative body limit", key: "MAX_BODY_BYTES", val: "-10"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			setRequired(t)
			t.Setenv(tc.key, tc.val)

			_, err := config.Load()

			require.Error(t, err)
		})
	}
}

func TestLoadRejectsMalformedDuration(t *testing.T) {
	setRequired(t)
	t.Setenv("JWT_TTL", "forever")

	_, err := config.Load()

	require.Error(t, err)
}

func TestSlogLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level string
		want  slog.Level
	}{
		{level: "debug", want: slog.LevelDebug},
		{level: "info", want: slog.LevelInfo},
		{level: "warn", want: slog.LevelWarn},
		{level: "error", want: slog.LevelError},
		{level: "", want: slog.LevelInfo},
		{level: "verbose", want: slog.LevelInfo},
	}

	for _, tc := range tests {
		t.Run("level "+tc.level, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, config.Config{LogLevel: tc.level}.SlogLevel())
		})
	}
}
