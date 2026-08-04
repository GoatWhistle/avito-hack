package infra

import (
	"context"
	"encoding/json"
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
	StreakDays      int          `json:"streak_days"`
	LastCheckInDate *time.Time   `json:"last_checkin_date,omitempty"`
	LastDecayTime   time.Time    `json:"last_decay_time"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

func toCachedPet(p *domain.Pet) cachedPet {
	return cachedPet{
		ID: p.ID(), UserID: p.UserID(), Name: p.Name(), Stage: p.Stage(), Level: p.Level(), XP: p.XP(),
		NextLevelXP: p.NextLevelXP(), Satiety: p.Satiety(), Happiness: p.Happiness(),
		StreakDays: p.StreakDays(), LastCheckInDate: p.LastCheckInDate(),
		LastDecayTime: p.LastDecayTime(), UpdatedAt: p.UpdatedAt(),
	}
}

func (p cachedPet) toDomain() *domain.Pet {
	return domain.Restore(domain.RestoreParams{
		ID: p.ID, UserID: p.UserID, Name: p.Name, Stage: p.Stage, Level: p.Level, XP: p.XP,
		NextLevelXP: p.NextLevelXP, Satiety: p.Satiety, Happiness: p.Happiness,
		StreakDays: p.StreakDays, LastCheckInDate: p.LastCheckInDate,
		LastDecayTime: p.LastDecayTime, UpdatedAt: p.UpdatedAt,
	})
}
