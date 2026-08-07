package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/avito-hack/backend/internal/config"
	"github.com/avito-hack/backend/internal/server"
	"github.com/avito-hack/backend/internal/shared/logger"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

// @title Avito Hack API
// @version 1.0.0
// @description REST и WebSocket API маркетплейса с геймификацией (питомец-енот)
//
// @contact.name Avito Hack Team
// @license.name MIT
//
// @servers.url http://localhost:8080
// @servers.description Локальная разработка
//
// @securityDefinitions.apikey bearerAuth
// @in header
// @name Authorization
//
// @tag.name Auth
// @tag.name Users
// @tag.name Items
// @tag.name Photos
// @tag.name Favorites
// @tag.name Pet
// @tag.name Rewards
// @tag.name Leaderboard
// @tag.name Raccoon
// @tag.name WebSocket
// @tag.name System
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

	pool, err := postgres.NewPoolWithOptions(ctx, cfg.DatabaseURL, cfg.PoolOptions())
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr,
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	})
	defer func() {
		if closeErr := redisClient.Close(); closeErr != nil {
			log.Warn("close redis", slog.Any("error", closeErr))
		}
	}()

	built, err := buildModules(cfg, pool, redisClient, log)
	if err != nil {
		return fmt.Errorf("build modules: %w", err)
	}

	stopBackground := built.background.run(ctx)
	defer stopBackground()

	handler := server.NewRouter(server.RouterDeps{
		Config:    cfg,
		Pool:      pool,
		Logger:    log,
		Modules:   built.registrars,
		WebSocket: built.webSocket,
	})

	return server.New(cfg, handler, log).Run(ctx)
}
