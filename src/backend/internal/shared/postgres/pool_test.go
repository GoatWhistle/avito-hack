package postgres_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/avito-hack/backend/internal/shared/postgres"
)

func TestDefaultPoolOptionsAreSane(t *testing.T) {
	t.Parallel()

	opts := postgres.DefaultPoolOptions()

	assert.Positive(t, opts.MaxConns)
	assert.GreaterOrEqual(t, opts.MinConns, int32(0))
	assert.LessOrEqual(t, opts.MinConns, opts.MaxConns)
	assert.Positive(t, opts.MaxConnLifetime)
	assert.Positive(t, opts.MaxConnIdleTime)
	assert.Positive(t, opts.HealthCheckPeriod)
	assert.Positive(t, opts.ConnectTimeout)
}

func TestNewPoolWithOptionsRejectsBadURL(t *testing.T) {
	t.Parallel()

	_, err := postgres.NewPoolWithOptions(t.Context(), "://nonsense", postgres.PoolOptions{
		MaxConns:       4,
		ConnectTimeout: time.Second,
	})

	assert.Error(t, err)
}
