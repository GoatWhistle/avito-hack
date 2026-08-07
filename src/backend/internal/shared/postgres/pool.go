package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultMaxConns          = 20
	defaultMinConns          = 2
	defaultMaxConnLifetime   = time.Hour
	defaultMaxConnIdleTime   = 30 * time.Minute
	defaultHealthCheckPeriod = time.Minute
	defaultConnectTimeout    = 5 * time.Second
)

type PoolOptions struct {
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
	ConnectTimeout    time.Duration
}

func DefaultPoolOptions() PoolOptions {
	return PoolOptions{
		MaxConns:          defaultMaxConns,
		MinConns:          defaultMinConns,
		MaxConnLifetime:   defaultMaxConnLifetime,
		MaxConnIdleTime:   defaultMaxConnIdleTime,
		HealthCheckPeriod: defaultHealthCheckPeriod,
		ConnectTimeout:    defaultConnectTimeout,
	}
}

func (o PoolOptions) withDefaults() PoolOptions {
	d := DefaultPoolOptions()

	if o.MaxConns <= 0 {
		o.MaxConns = d.MaxConns
	}
	if o.MinConns < 0 {
		o.MinConns = d.MinConns
	}
	if o.MaxConnLifetime <= 0 {
		o.MaxConnLifetime = d.MaxConnLifetime
	}
	if o.MaxConnIdleTime <= 0 {
		o.MaxConnIdleTime = d.MaxConnIdleTime
	}
	if o.HealthCheckPeriod <= 0 {
		o.HealthCheckPeriod = d.HealthCheckPeriod
	}
	if o.ConnectTimeout <= 0 {
		o.ConnectTimeout = d.ConnectTimeout
	}

	return o
}

func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	return NewPoolWithOptions(ctx, databaseURL, DefaultPoolOptions())
}

func NewPoolWithOptions(ctx context.Context, databaseURL string, opts PoolOptions) (*pgxpool.Pool, error) {
	opts = opts.withDefaults()

	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}

	cfg.MaxConns = opts.MaxConns
	cfg.MinConns = opts.MinConns
	cfg.MaxConnLifetime = opts.MaxConnLifetime
	cfg.MaxConnIdleTime = opts.MaxConnIdleTime
	cfg.HealthCheckPeriod = opts.HealthCheckPeriod
	cfg.ConnConfig.ConnectTimeout = opts.ConnectTimeout

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, opts.ConnectTimeout)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()

		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}
