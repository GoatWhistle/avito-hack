package api_test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/api"
	"github.com/avito-hack/backend/internal/shared/ws"
)

type recordingHub struct {
	users    []uuid.UUID
	messages []ws.Message
}

func (h *recordingHub) Broadcast(userID uuid.UUID, m ws.Message) {
	h.users = append(h.users, userID)
	h.messages = append(h.messages, m)
}

func (h *recordingHub) last() ws.Message {
	return h.messages[len(h.messages)-1]
}

func payloadOf(t *testing.T, m ws.Message) map[string]any {
	t.Helper()

	raw, err := json.Marshal(m.Payload)
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(raw, &decoded))

	return decoded
}

func TestNotifierPetUpdatedBroadcastsPetPayload(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	hub := &recordingHub{}

	api.NewNotifier(hub).PetUpdated(userID, testPet(userID))

	require.Len(t, hub.messages, 1)
	assert.Equal(t, "pet.updated", hub.last().Type)
	assert.Equal(t, userID, hub.users[0])

	payload := payloadOf(t, hub.last())
	assert.Equal(t, "Enot", payload["name"])
	assert.Equal(t, userID.String(), payload["user_id"])
}

func TestNotifierPetHatchedBroadcastsPetPayload(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	hub := &recordingHub{}

	api.NewNotifier(hub).PetHatched(userID, testPet(userID))

	require.Len(t, hub.messages, 1)
	assert.Equal(t, "pet.hatched", hub.last().Type)
	assert.Equal(t, userID, hub.users[0])
	assert.Equal(t, true, payloadOf(t, hub.last())["is_hatched"])
}

func TestNotifierXPGained(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	hub := &recordingHub{}

	api.NewNotifier(hub).XPGained(userID, 25, "check_in", 140)

	require.Len(t, hub.messages, 1)
	assert.Equal(t, "xp.gained", hub.last().Type)
	assert.Equal(t, userID, hub.users[0])

	payload := payloadOf(t, hub.last())
	assert.InDelta(t, 25, payload["amount"], 0)
	assert.Equal(t, "check_in", payload["reason"])
	assert.InDelta(t, 140, payload["total"], 0)
}

func TestNotifierLevelUp(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	hub := &recordingHub{}

	api.NewNotifier(hub).LevelUp(userID, 7)

	require.Len(t, hub.messages, 1)
	assert.Equal(t, "level.up", hub.last().Type)
	assert.InDelta(t, 7, payloadOf(t, hub.last())["level"], 0)
}

func TestNotifierRewardGranted(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	hub := &recordingHub{}

	api.NewNotifier(hub).RewardGranted(userID, "reward-9", "Золотая лапа")

	require.Len(t, hub.messages, 1)
	assert.Equal(t, "reward.granted", hub.last().Type)

	payload := payloadOf(t, hub.last())
	assert.Equal(t, "reward-9", payload["reward_id"])
	assert.Equal(t, "Золотая лапа", payload["title"])
}

func TestNotifierStreakUpdated(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		days      int
		milestone bool
	}{
		{name: "plain day", days: 3, milestone: false},
		{name: "milestone day", days: 7, milestone: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userID := uuid.New()
			hub := &recordingHub{}

			api.NewNotifier(hub).StreakUpdated(userID, tc.days, tc.milestone)

			require.Len(t, hub.messages, 1)
			assert.Equal(t, "streak.updated", hub.last().Type)

			payload := payloadOf(t, hub.last())
			assert.InDelta(t, tc.days, payload["days"], 0)
			assert.Equal(t, tc.milestone, payload["milestone"])
		})
	}
}

func TestNotifierRoutesToTheRequestedUserOnly(t *testing.T) {
	t.Parallel()

	first, second := uuid.New(), uuid.New()
	hub := &recordingHub{}
	notifier := api.NewNotifier(hub)

	notifier.LevelUp(first, 2)
	notifier.LevelUp(second, 3)

	require.Len(t, hub.users, 2)
	assert.Equal(t, first, hub.users[0])
	assert.Equal(t, second, hub.users[1])
}

func TestNotifierWithNilHubIsNoop(t *testing.T) {
	t.Parallel()

	notifier := api.NewNotifier(nil)
	userID := uuid.New()

	assert.NotPanics(t, func() {
		notifier.PetUpdated(userID, testPet(userID))
		notifier.PetHatched(userID, testPet(userID))
		notifier.XPGained(userID, 1, "r", 2)
		notifier.LevelUp(userID, 1)
		notifier.RewardGranted(userID, "r", "t")
		notifier.StreakUpdated(userID, 1, false)
	})
}

func TestNotifierWithNilPetDoesNotBroadcast(t *testing.T) {
	t.Parallel()

	hub := &recordingHub{}
	notifier := api.NewNotifier(hub)
	userID := uuid.New()

	notifier.PetUpdated(userID, nil)
	notifier.PetHatched(userID, nil)

	assert.Empty(t, hub.messages)
}
