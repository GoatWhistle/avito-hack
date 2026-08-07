package infra

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/pet/app"
)

var errInvalidHotState = errors.New("invalid redis hot state response")

func hotKey(userID uuid.UUID) string      { return hotKeyString(userID.String()) }
func hotKeyString(userID string) string   { return "pet:hot:" + userID }
func cooldownKey(userID uuid.UUID) string { return "pet:cooldown:stroke:" + userID.String() }
func dirtyKey() string                    { return "pet:hot:dirty" }

func hotStateFromValues(userID uuid.UUID, values []any) (app.HotState, error) {
	if len(values) < 4 {
		return app.HotState{}, errInvalidHotState
	}

	happiness, err := valueInt64(values[0])
	if err != nil {
		return app.HotState{}, err
	}
	satiety, err := valueInt64(values[1])
	if err != nil {
		return app.HotState{}, err
	}
	version, err := valueInt64(values[2])
	if err != nil {
		return app.HotState{}, err
	}
	updatedAt, err := valueInt64(values[3])
	if err != nil {
		return app.HotState{}, err
	}

	return app.HotState{
		UserID: userID, Happiness: int(happiness), Satiety: int(satiety),
		Version: version, UpdatedAt: time.Unix(0, updatedAt).UTC(),
	}, nil
}

func hotStateFromMap(userID uuid.UUID, values map[string]string) (app.HotState, error) {
	return hotStateFromValues(userID, []any{
		values["happiness"], values["satiety"], values["version"], values["updated_at"],
	})
}

func valueInt64(value any) (int64, error) {
	switch typed := value.(type) {
	case int64:
		return typed, nil
	case string:
		parsed, err := strconv.ParseInt(typed, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("parse redis integer: %w", err)
		}

		return parsed, nil
	case nil:
		return 0, errInvalidHotState
	default:
		return 0, fmt.Errorf("unexpected redis value type %T", value)
	}
}
