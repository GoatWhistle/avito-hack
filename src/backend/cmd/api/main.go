package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/avito-hack/backend/internal/config"
	"github.com/avito-hack/backend/internal/module/favorite"
	"github.com/avito-hack/backend/internal/module/item"
	"github.com/avito-hack/backend/internal/module/pet"
	"github.com/avito-hack/backend/internal/module/raccoon"
	"github.com/avito-hack/backend/internal/module/user"
	"github.com/avito-hack/backend/internal/server"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/clock"
	"github.com/avito-hack/backend/internal/shared/events"
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

	modules, webSocket, err := buildModules(cfg, pool, redisClient)
	if err != nil {
		return fmt.Errorf("build modules: %w", err)
	}

	handler := server.NewRouter(server.RouterDeps{
		Config:    cfg,
		Pool:      pool,
		Logger:    log,
		Modules:   modules,
		WebSocket: webSocket,
	})

	return server.New(cfg, handler, log).Run(ctx)
}

func buildModules(
	cfg config.Config,
	pool *pgxpool.Pool,
	redisClient *redis.Client,
) ([]server.ModuleRegistrar, http.Handler, error) {
	appClock := clock.New()
	tx := postgres.NewTxManager(pool)
	validator := validate.New()
	tokens := auth.NewTokenService(cfg.JWTSecret, cfg.JWTTTL, appClock.Now)

	authenticate := middleware.Authenticate(tokens)
	optionalAuth := middleware.OptionalAuthenticate(tokens)

	bus := events.NewBus(slog.Default())

	userModule := user.New(user.Options{
		Pool:         pool,
		Tx:           tx,
		Clock:        appClock,
		Tokens:       tokens,
		Bus:          bus,
		Validator:    validator,
		Authenticate: authenticate,
		MaxBodyBytes: cfg.MaxBodyBytes,
	})

	itemModule := item.New(item.Options{
		Pool:          pool,
		Tx:            tx,
		Clock:         appClock,
		Bus:           bus,
		Validator:     validator,
		Authenticate:  authenticate,
		OptionalAuth:  optionalAuth,
		MaxBodyBytes:  cfg.MaxBodyBytes,
		MaxPhotoBytes: cfg.MaxPhotoBytes,
		UploadDir:     cfg.UploadDir,
		UploadURL:     cfg.UploadURL,
	})

	favoriteModule := favorite.New(favorite.Options{
		Pool:         pool,
		Tx:           tx,
		Clock:        appClock,
		Bus:          bus,
		Authenticate: authenticate,
	})

	petModule, err := pet.New(pet.Options{
		Pool:             pool,
		Redis:            redisClient,
		Tx:               tx,
		Clock:            appClock,
		Tokens:           tokens,
		Bus:              bus,
		AllowedOrigins:   cfg.AllowedOrigins,
		RewardHMACSecret: cfg.RewardHMACSecret,
		Authenticate:     authenticate,
	})
	if err != nil {
		return nil, nil, err
	}

	raccoonModule := raccoon.New(raccoon.Options{
		Pool:         pool,
		Pets:         petModule.Service,
		Rewards:      petModule.Rewards,
		Validator:    validator,
		Authenticate: authenticate,
		MaxBodyBytes: cfg.MaxBodyBytes,
	})

	return []server.ModuleRegistrar{
		userModule.Handlers,
		itemModule.Handlers,
		favoriteModule.Handlers,
		raccoonModule.Handlers,
		petModule.Leaderboard,
		petModule.Pet,
		petModule.RewardAPI,
	}, petModule.WebSocket, nil
}
