package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/config"
	"github.com/avito-hack/backend/internal/module/item"
	"github.com/avito-hack/backend/internal/module/user"
	"github.com/avito-hack/backend/internal/server"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/clock"
	"github.com/avito-hack/backend/internal/shared/logger"
	"github.com/avito-hack/backend/internal/shared/middleware"
	"github.com/avito-hack/backend/internal/shared/postgres"
	"github.com/avito-hack/backend/internal/shared/validate"
)

func main() {
	if isHealthcheckMode() {
		if err := runHealthcheck(healthcheckAddr()); err != nil {
			slog.Error("healthcheck failed", slog.Any("error", err))
			os.Exit(1)
		}

		return
	}

	if err := run(); err != nil {
		slog.Error("fatal error", slog.Any("error", err))
		os.Exit(1)
	}
}

func healthcheckAddr() string {
	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		return addr
	}

	return ":8080"
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	log := logger.New(logger.Options{
		Level:  cfg.SlogLevel(),
		Format: cfg.LogFormat,
		Color:  cfg.LogColor,
	})
	slog.SetDefault(log)

	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close()

	handler := server.NewRouter(server.RouterDeps{
		Config:  cfg,
		Pool:    pool,
		Logger:  log,
		Modules: buildModules(cfg, pool),
	})

	return server.New(cfg, handler, log).Run(ctx)
}

func buildModules(cfg config.Config, pool *pgxpool.Pool) []server.ModuleRegistrar {
	appClock := clock.New()
	tx := postgres.NewTxManager(pool)
	validator := validate.New()
	tokens := auth.NewTokenService(cfg.JWTSecret, cfg.JWTTTL, appClock.Now)

	authenticate := middleware.Authenticate(tokens)
	optionalAuth := middleware.OptionalAuthenticate(tokens)

	userModule := user.New(user.Options{
		Pool:         pool,
		Tx:           tx,
		Clock:        appClock,
		Tokens:       tokens,
		Validator:    validator,
		Authenticate: authenticate,
		MaxBodyBytes: cfg.MaxBodyBytes,
	})

	itemModule := item.New(item.Options{
		Pool:         pool,
		Tx:           tx,
		Clock:        appClock,
		Validator:    validator,
		Authenticate: authenticate,
		OptionalAuth: optionalAuth,
		MaxBodyBytes: cfg.MaxBodyBytes,
	})

	return []server.ModuleRegistrar{
		userModule.Handlers,
		itemModule.Handlers,
	}
}
