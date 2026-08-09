package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/module/pet/domain"
)

const (
	hotStateTTL    = 24 * time.Hour
	strokeCooldown = time.Second
	feedCooldown   = 5 * time.Hour

	fieldHappiness = "happiness"
	fieldSatiety   = "satiety"
)

type RedisHotStateStore struct {
	client *redis.Client
}

func NewRedisHotStateStore(client *redis.Client) *RedisHotStateStore {
	return &RedisHotStateStore{client: client}
}

func (s *RedisHotStateStore) GetOrInitialize(
	ctx context.Context,
	initial app.HotState,
) (app.HotState, error) {
	values, err := initializeHotState.Run(ctx, s.client,
		[]string{hotKey(initial.UserID), dirtyKey()}, hotArgs(initial)...,
	).Slice()
	if err != nil {
		return app.HotState{}, fmt.Errorf("initialize redis hot state: %w", err)
	}

	return hotStateFromValues(initial.UserID, values)
}

type nudgeRule struct {
	action      string
	field       string
	gain        int
	cooldown    time.Duration
	cooldownKey func(uuid.UUID) string
}

func (s *RedisHotStateStore) Stroke(
	ctx context.Context,
	initial app.HotState,
	now time.Time,
) (app.HotState, bool, error) {
	return s.nudge(ctx, initial, now, nudgeRule{
		action: "stroke", field: fieldHappiness, gain: domain.StrokeHappinessGain,
		cooldown: strokeCooldown, cooldownKey: strokeCooldownKey,
	})
}

func (s *RedisHotStateStore) Feed(
	ctx context.Context,
	initial app.HotState,
	now time.Time,
) (app.HotState, bool, error) {
	return s.nudge(ctx, initial, now, nudgeRule{
		action: "feed", field: fieldSatiety, gain: domain.FeedSatietyGain,
		cooldown: feedCooldown, cooldownKey: feedCooldownKey,
	})
}

func (s *RedisHotStateStore) FeedAvailableAt(
	ctx context.Context,
	userID uuid.UUID,
	now time.Time,
) (*time.Time, error) {
	ttl, err := s.client.TTL(ctx, feedCooldownKey(userID)).Result()
	if err != nil {
		return nil, fmt.Errorf("read feed cooldown ttl: %w", err)
	}
	if ttl <= 0 {
		return nil, nil
	}

	available := now.Add(ttl).UTC()

	return &available, nil
}

func (s *RedisHotStateStore) nudge(
	ctx context.Context,
	initial app.HotState,
	now time.Time,
	rule nudgeRule,
) (app.HotState, bool, error) {
	values, err := nudgeHotState.Run(ctx, s.client,
		[]string{hotKey(initial.UserID), rule.cooldownKey(initial.UserID), dirtyKey()},
		initial.UserID.String(), initial.Happiness, initial.Satiety, initial.Version,
		initial.UpdatedAt.UnixNano(), now.UnixNano(), int(hotStateTTL.Seconds()),
		int(rule.cooldown.Seconds()), rule.gain, domain.MaxParameterValue, rule.field,
	).Slice()
	if err != nil {
		return app.HotState{}, false, fmt.Errorf("%s redis hot state: %w", rule.action, err)
	}
	if len(values) < 5 {
		return app.HotState{}, false, errInvalidHotState
	}

	state, err := hotStateFromValues(initial.UserID, values[:4])
	if err != nil {
		return app.HotState{}, false, err
	}
	applied, err := valueInt64(values[4])
	if err != nil {
		return app.HotState{}, false, err
	}

	return state, applied == 1, nil
}

func (s *RedisHotStateStore) DirtyBatch(ctx context.Context, limit int64) ([]app.HotState, error) {
	ids, err := s.client.SRandMemberN(ctx, dirtyKey(), limit).Result()
	if err != nil {
		return nil, fmt.Errorf("read dirty pet ids: %w", err)
	}
	if len(ids) == 0 {
		return nil, nil
	}

	values, err := s.readStates(ctx, ids)
	if err != nil {
		return nil, err
	}

	states := make([]app.HotState, 0, len(ids))
	for index, id := range ids {
		userID, parseErr := uuid.Parse(id)
		if parseErr != nil {
			return nil, fmt.Errorf("parse dirty pet id: %w", parseErr)
		}
		if len(values[index]) == 0 {
			if remErr := s.client.SRem(ctx, dirtyKey(), id).Err(); remErr != nil {
				return nil, fmt.Errorf("prune dirty pet id: %w", remErr)
			}

			continue
		}
		state, parseErr := hotStateFromMap(userID, values[index])
		if parseErr != nil {
			return nil, parseErr
		}
		states = append(states, state)
	}

	return states, nil
}

func (s *RedisHotStateStore) Acknowledge(ctx context.Context, state app.HotState) error {
	err := acknowledgeHotState.Run(ctx, s.client,
		[]string{hotKey(state.UserID), dirtyKey()}, state.UserID.String(), state.Version,
	).Err()
	if err != nil {
		return fmt.Errorf("acknowledge redis hot state: %w", err)
	}

	return nil
}

func (s *RedisHotStateStore) readStates(ctx context.Context, ids []string) ([]map[string]string, error) {
	pipe := s.client.Pipeline()
	commands := make([]*redis.MapStringStringCmd, 0, len(ids))
	for _, id := range ids {
		commands = append(commands, pipe.HGetAll(ctx, hotKeyString(id)))
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return nil, fmt.Errorf("read dirty pet states: %w", err)
	}

	values := make([]map[string]string, 0, len(commands))
	for _, command := range commands {
		values = append(values, command.Val())
	}

	return values, nil
}

func hotArgs(state app.HotState) []any {
	return []any{
		state.UserID.String(), state.Happiness, state.Satiety,
		state.Version, state.UpdatedAt.UnixNano(), int(hotStateTTL.Seconds()),
	}
}
