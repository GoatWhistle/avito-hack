//go:build integration

package integration_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	itemapp "github.com/avito-hack/backend/internal/module/item/app"
	itemdomain "github.com/avito-hack/backend/internal/module/item/domain"
	iteminfra "github.com/avito-hack/backend/internal/module/item/infra"
	petapp "github.com/avito-hack/backend/internal/module/pet/app"
	petdomain "github.com/avito-hack/backend/internal/module/pet/domain"
	petinfra "github.com/avito-hack/backend/internal/module/pet/infra"
	userapp "github.com/avito-hack/backend/internal/module/user/app"
	userinfra "github.com/avito-hack/backend/internal/module/user/infra"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/events"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

const (
	testJWTSecret  = "integration-jwt-secret-value-32-bytes"
	testHMACSecret = "integration-hmac-secret-value-32-bytes"
	testJWTTTL     = 15 * time.Minute
)

type stubClock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *stubClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.now
}

func (c *stubClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

type memoryCache struct {
	mu   sync.Mutex
	pets map[uuid.UUID]*petdomain.Pet
}

func newMemoryCache() *memoryCache {
	return &memoryCache{pets: map[uuid.UUID]*petdomain.Pet{}}
}

func (c *memoryCache) Get(_ context.Context, userID uuid.UUID) (*petdomain.Pet, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.pets[userID], nil
}

func (c *memoryCache) Set(_ context.Context, pet *petdomain.Pet) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.pets[pet.UserID()] = pet

	return nil
}

func (c *memoryCache) Delete(_ context.Context, userID uuid.UUID) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.pets, userID)

	return nil
}

type env struct {
	clock    *stubClock
	tx       *postgres.TxManager
	tokens   *auth.TokenService
	register *userapp.RegisterUserHandler
	login    *userapp.LoginUserHandler
	profile  *userapp.GetProfileHandler
	create   *itemapp.CreateItemHandler
	status   *itemapp.ChangeStatusHandler
	items    itemdomain.Repository
	pets     *petapp.Service
	rewards  *petapp.RewardService
}

func newEnv(t *testing.T) *env {
	t.Helper()

	clk := &stubClock{now: time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC)}
	tx := postgres.NewTxManager(pool)
	tokens := auth.NewTokenService(testJWTSecret, testJWTTTL, clk.Now)

	users := userinfra.NewPgRepository(pool)
	items := iteminfra.NewPgRepository(pool)
	photos := iteminfra.NewPgPhotoRepository(pool)

	petRepo := petinfra.NewPgRepository(pool)
	journal := petinfra.NewPgXPEventRepository(pool)

	signer, err := petdomain.NewRewardSigner(testHMACSecret)
	require.NoError(t, err)

	pets := petapp.NewService(petRepo, journal, newMemoryCache(), tx, clk)

	rewards := petapp.NewRewardService(petapp.RewardServiceDeps{
		Rewards: petinfra.NewPgRewardRepository(pool),
		Pets:    petRepo,
		Signer:  signer,
		Tx:      tx,
		Clock:   clk,
	})

	return &env{
		clock:    clk,
		tx:       tx,
		tokens:   tokens,
		register: userapp.NewRegisterUserHandler(users, tx, clk, nil, tokens),
		login:    userapp.NewLoginUserHandler(users, tokens),
		profile:  userapp.NewGetProfileHandler(users),
		create:   itemapp.NewCreateItemHandler(items, tx, clk),
		status:   itemapp.NewChangeStatusHandler(items, photos, tx, clk, events.NopPublisher{}, nil),
		items:    items,
		pets:     pets,
		rewards:  rewards,
	}
}

func uniqueEmail() string {
	return "user-" + uuid.NewString() + "@example.com"
}
