package infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/avito-hack/backend/internal/module/pet/domain"
)

const petCacheTTL = 5 * time.Minute

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (c *RedisCache) Get(ctx context.Context, userID uuid.UUID) (*domain.Pet, error) {
	raw, err := c.client.Get(ctx, cacheKey(userID)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get cached pet: %w", err)
	}

	var value cachedPet
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, fmt.Errorf("decode cached pet: %w", err)
	}

	return value.toDomain(), nil
}

func (c *RedisCache) Set(ctx context.Context, pet *domain.Pet) error {
	raw, err := json.Marshal(toCachedPet(pet))
	if err != nil {
		return fmt.Errorf("encode cached pet: %w", err)
	}
	if err := c.client.Set(ctx, cacheKey(pet.UserID()), raw, petCacheTTL).Err(); err != nil {
		return fmt.Errorf("set cached pet: %w", err)
	}

	return nil
}

func (c *RedisCache) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := c.client.Del(ctx, cacheKey(userID)).Err(); err != nil {
		return fmt.Errorf("delete cached pet: %w", err)
	}

	return nil
}

func cacheKey(userID uuid.UUID) string {
	return "pet:user:" + userID.String()
}

type cachedPet struct {
	ID              uuid.UUID    `json:"id"`
	UserID          uuid.UUID    `json:"user_id"`
	Name            string       `json:"name"`
	Stage           domain.Stage `json:"stage"`
	Level           int          `json:"level"`
	XP              int          `json:"xp"`
	NextLevelXP     int          `json:"next_level_xp"`
	Satiety         int          `json:"satiety"`
	Happiness       int          `json:"happiness"`
	Energy          int          `json:"energy"`
	StreakDays      int          `json:"streak_days"`
	Freezes         int          `json:"freezes"`
	LastCheckInDate *time.Time   `json:"last_checkin_date,omitempty"`
	HatchedAt       *time.Time   `json:"hatched_at,omitempty"`
	LastDecayTime   time.Time    `json:"last_decay_time"`
	UpdatedAt       time.Time    `json:"updated_at"`

	InteractionVersion int64 `json:"interaction_version"`
}

func toCachedPet(p *domain.Pet) cachedPet {
	return cachedPet{
		ID: p.ID(), UserID: p.UserID(), Name: p.Name(), Stage: p.Stage(), Level: p.Level(), XP: p.XP(),
		NextLevelXP: p.NextLevelXP(), Satiety: p.Satiety(), Happiness: p.Happiness(),
		Energy: p.Energy(), StreakDays: p.StreakDays(), Freezes: p.Freezes(),
		LastCheckInDate: p.LastCheckInDate(), HatchedAt: p.HatchedAt(),
		LastDecayTime: p.LastDecayTime(), UpdatedAt: p.UpdatedAt(),
		InteractionVersion: p.InteractionVersion(),
	}
}

func (p cachedPet) toDomain() *domain.Pet {
	return domain.Restore(domain.RestoreParams{
		ID: p.ID, UserID: p.UserID, Name: p.Name, Stage: p.Stage, Level: p.Level, XP: p.XP,
		NextLevelXP: p.NextLevelXP, Satiety: p.Satiety, Happiness: p.Happiness,
		Energy: p.Energy, StreakDays: p.StreakDays, Freezes: p.Freezes,
		LastCheckInDate: p.LastCheckInDate, HatchedAt: p.HatchedAt,
		LastDecayTime: p.LastDecayTime, UpdatedAt: p.UpdatedAt,
		InteractionVersion: p.InteractionVersion,
	})
}
