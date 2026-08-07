package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/config"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

func TestPoolOptionDefaultsMatchPostgresDefaults(t *testing.T) {
	setRequired(t)

	cfg, err := config.Load()

	require.NoError(t, err)
	assert.Equal(t, postgres.DefaultPoolOptions(), cfg.PoolOptions())
}

func TestPoolOptionsMapFromConfig(t *testing.T) {
	setRequired(t)
	t.Setenv("DB_MAX_CONNS", "50")
	t.Setenv("DB_MIN_CONNS", "5")
	t.Setenv("DB_MAX_CONN_LIFETIME", "2h")
	t.Setenv("DB_MAX_CONN_IDLE_TIME", "10m")
	t.Setenv("DB_HEALTH_CHECK_PERIOD", "30s")
	t.Setenv("DB_CONNECT_TIMEOUT", "7s")

	cfg, err := config.Load()

	require.NoError(t, err)
	assert.Equal(t, postgres.PoolOptions{
		MaxConns:          50,
		MinConns:          5,
		MaxConnLifetime:   2 * time.Hour,
		MaxConnIdleTime:   10 * time.Minute,
		HealthCheckPeriod: 30 * time.Second,
		ConnectTimeout:    7 * time.Second,
	}, cfg.PoolOptions())
}

func TestPoolConfigValidation(t *testing.T) {
	tests := []struct {
		name string
		key  string
		val  string
	}{
		{name: "zero max conns", key: "DB_MAX_CONNS", val: "0"},
		{name: "negative max conns", key: "DB_MAX_CONNS", val: "-1"},
		{name: "excessive max conns", key: "DB_MAX_CONNS", val: "5000"},
		{name: "negative min conns", key: "DB_MIN_CONNS", val: "-1"},
		{name: "min above max", key: "DB_MIN_CONNS", val: "999"},
		{name: "zero lifetime", key: "DB_MAX_CONN_LIFETIME", val: "0s"},
		{name: "zero idle time", key: "DB_MAX_CONN_IDLE_TIME", val: "0s"},
		{name: "zero health check", key: "DB_HEALTH_CHECK_PERIOD", val: "0s"},
		{name: "zero connect timeout", key: "DB_CONNECT_TIMEOUT", val: "0s"},
		{name: "excessive connect timeout", key: "DB_CONNECT_TIMEOUT", val: "30m"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			setRequired(t)
			t.Setenv(tc.key, tc.val)

			_, err := config.Load()

			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.key)
		})
	}
}
