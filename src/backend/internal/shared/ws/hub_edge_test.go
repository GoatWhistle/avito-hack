package ws_test

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/ws"
)

func TestHub_RegisterRejectsNilSink(t *testing.T) {
	t.Parallel()

	hub := ws.NewHub()
	userID := uuid.New()

	release := hub.Register(userID, nil)
	require.NotNil(t, release)
	require.NotPanics(t, release)

	sink := &recordingSink{}
	hub.Register(userID, sink)
	hub.Broadcast(userID, ws.Message{Type: "pet.updated"})

	assert.Equal(t, 1, sink.count())
}

func TestHub_RegisterRejectsNilUserID(t *testing.T) {
	t.Parallel()

	hub := ws.NewHub()
	sink := &recordingSink{}

	release := hub.Register(uuid.Nil, sink)
	require.NotNil(t, release)

	hub.Broadcast(uuid.Nil, ws.Message{Type: "pet.updated"})
	assert.Equal(t, 0, sink.count())

	require.NotPanics(t, release)
}

func TestHub_ReleaseIsIdempotent(t *testing.T) {
	t.Parallel()

	hub := ws.NewHub()
	userID := uuid.New()
	sink := &recordingSink{}

	release := hub.Register(userID, sink)
	release()
	require.NotPanics(t, release)

	hub.Broadcast(userID, ws.Message{Type: "xp.gained"})
	assert.Equal(t, 0, sink.count())
}

func TestHub_ReleaseOfOneConnectionKeepsOthers(t *testing.T) {
	t.Parallel()

	hub := ws.NewHub()
	userID := uuid.New()
	first, second := &recordingSink{}, &recordingSink{}

	release := hub.Register(userID, first)
	hub.Register(userID, second)

	release()
	hub.Broadcast(userID, ws.Message{Type: "level.up"})

	assert.Equal(t, 0, first.count())
	assert.Equal(t, 1, second.count())
}

func TestHub_RegisteringSameSinkTwiceDeliversOnce(t *testing.T) {
	t.Parallel()

	hub := ws.NewHub()
	userID := uuid.New()
	sink := &recordingSink{}

	hub.Register(userID, sink)
	hub.Register(userID, sink)

	hub.Broadcast(userID, ws.Message{Type: "pet.updated"})

	assert.Equal(t, 1, sink.count())
}

func TestHub_BroadcastCarriesPayload(t *testing.T) {
	t.Parallel()

	hub := ws.NewHub()
	userID := uuid.New()
	sink := &recordingSink{}

	hub.Register(userID, sink)
	hub.Broadcast(userID, ws.Message{Type: "xp.gained", Payload: map[string]int{"amount": 10}})

	sink.mu.Lock()
	defer sink.mu.Unlock()

	require.Len(t, sink.messages, 1)
	assert.Equal(t, "xp.gained", sink.messages[0].Type)
	assert.Equal(t, map[string]int{"amount": 10}, sink.messages[0].Payload)
}

func TestHub_ConcurrentUnregisterOfDistinctUsers(t *testing.T) {
	t.Parallel()

	hub := ws.NewHub()

	var wg sync.WaitGroup
	for range 64 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			userID := uuid.New()
			sink := &recordingSink{}
			release := hub.Register(userID, sink)
			hub.Broadcast(userID, ws.Message{Type: "streak.updated"})
			release()
			hub.Broadcast(userID, ws.Message{Type: "streak.updated"})

			assert.Equal(t, 1, sink.count())
		}()
	}
	wg.Wait()
}
