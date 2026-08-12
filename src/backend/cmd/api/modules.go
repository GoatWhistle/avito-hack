package main

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/avito-hack/backend/internal/config"
	"github.com/avito-hack/backend/internal/module/favorite"
	"github.com/avito-hack/backend/internal/module/games"
	"github.com/avito-hack/backend/internal/module/item"
	"github.com/avito-hack/backend/internal/module/pet"
	"github.com/avito-hack/backend/internal/module/raccoon"
	"github.com/avito-hack/backend/internal/module/user"
	"github.com/avito-hack/backend/internal/server"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/clock"
	"github.com/avito-hack/backend/internal/shared/events"
	"github.com/avito-hack/backend/internal/shared/kafka"
	"github.com/avito-hack/backend/internal/shared/middleware"
	"github.com/avito-hack/backend/internal/shared/postgres"
	"github.com/avito-hack/backend/internal/shared/validate"
)

type builtModules struct {
	registrars []server.ModuleRegistrar
	webSocket  http.Handler
	background *background
}

func startKafka(
	cfg config.Config,
	petModule *pet.Module,
	log *slog.Logger,
) (*kafka.Producer, *kafka.Consumer) {
	if !cfg.KafkaEnabled() {
		log.Info("kafka disabled, using in-process event delivery")

		return nil, nil
	}

	kafkaCfg := kafka.Config{
		Brokers:       cfg.KafkaBrokers,
		Topic:         cfg.KafkaTopic,
		ConsumerGroup: cfg.KafkaGroup,
		ClientID:      cfg.KafkaClientID,
		Timeout:       cfg.KafkaTimeout,
	}

	producer, err := kafka.NewProducer(kafkaCfg)
	if err != nil {
		log.Warn("kafka producer unavailable, falling back to in-process delivery",
			slog.Any("error", err))

		return nil, nil
	}

	consumer, err := kafka.NewConsumer(kafkaCfg, petModule.Async, log)
	if err != nil {
		log.Warn("kafka consumer unavailable, falling back to in-process delivery",
			slog.Any("error", err))
		producer.Close()

		return nil, nil
	}

	return producer, consumer
}

func buildModules(
	cfg config.Config,
	pool *pgxpool.Pool,
	redisClient *redis.Client,
	log *slog.Logger,
) (*builtModules, error) {
	appClock := clock.New()
	tx := postgres.NewTxManager(pool)
	validator := validate.New()
	tokens := auth.NewTokenService(cfg.JWTSecret, cfg.JWTTTL, appClock.Now)

	authenticate := middleware.Authenticate(tokens)
	optionalAuth := middleware.OptionalAuthenticate(tokens)

	bus := events.NewBus(log)

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
		Logger:           log,
		FlushInterval:    cfg.PetFlushInterval,
		FlushBatchSize:   cfg.PetFlushBatchSize,
	})
	if err != nil {
		return nil, err
	}

	producer, consumer := startKafka(cfg, petModule, log)

	var publisher events.Publisher = bus
	if producer != nil {
		publisher = events.NewAsyncPublisher(producer, bus, log)
	}

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

	petModule.Service.WithAccountGate(newAccountGate(userModule.Repository))

	itemModule := item.New(item.Options{
		Pool:             pool,
		Tx:               tx,
		Clock:            appClock,
		Bus:              publisher,
		Validator:        validator,
		Authenticate:     authenticate,
		OptionalAuth:     optionalAuth,
		MaxBodyBytes:     cfg.MaxBodyBytes,
		MaxPhotoBytes:    cfg.MaxPhotoBytes,
		UploadDir:        cfg.UploadDir,
		UploadURL:        cfg.UploadURL,
		OpenRouterAPIKey: cfg.OpenRouterAPIKey,
		OpenRouterModel:  cfg.OpenRouterModel,
	})

	favoriteModule := favorite.New(favorite.Options{
		Pool:         pool,
		Tx:           tx,
		Clock:        appClock,
		Bus:          publisher,
		Items:        itemModule.Repository,
		Authenticate: authenticate,
	})

	raccoonModule := raccoon.New(raccoon.Options{
		Pool:         pool,
		Pets:         petModule.Service,
		Rewards:      petModule.Rewards,
		Badges:       petModule.Badges,
		Validator:    validator,
		Authenticate: authenticate,
		MaxBodyBytes: cfg.MaxBodyBytes,
	})

	gamesModule := games.New(games.Options{
		Pool:         pool,
		Tx:           tx,
		Clock:        appClock,
		Validator:    validator,
		Authenticate: authenticate,
		MaxBodyBytes: cfg.MaxBodyBytes,
	})

	return &builtModules{
		registrars: []server.ModuleRegistrar{
			userModule.Handlers,
			itemModule.Handlers,
			favoriteModule.Handlers,
			raccoonModule.Handlers,
			gamesModule.Handlers,
			petModule.Leaderboard,
			petModule.Pet,
			petModule.RewardAPI,
			petModule.QuestAPI,
			petModule.SummaryAPI,
		},
		webSocket: petModule.WebSocket,
		background: &background{
			consumer: consumer,
			producer: producer,
			petModul: petModule,
			log:      log,
		},
	}, nil
}
