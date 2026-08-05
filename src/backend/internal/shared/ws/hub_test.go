package ws_test

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/ws"
)

type recordingSink struct {
	mu       sync.Mutex
	messages []ws.Message
}

func (s *recordingSink) Send(m ws.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.messages = append(s.messages, m)
}

func (s *recordingSink) count() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.messages)
}

func TestHub_BroadcastReachesAllUserConnections(t *testing.T) {
	t.Parallel()

	hub := ws.NewHub()
	userID := uuid.New()
	other := uuid.New()

	first, second, stranger := &recordingSink{}, &recordingSink{}, &recordingSink{}

	hub.Register(userID, first)
	hub.Register(userID, second)
	hub.Register(other, stranger)

	hub.Broadcast(userID, ws.Message{Type: "pet.updated"})

	require.Equal(t, 1, first.count())
	require.Equal(t, 1, second.count())
	require.Equal(t, 0, stranger.count())
}

func TestHub_UnregisterStopsDelivery(t *testing.T) {
	t.Parallel()

	hub := ws.NewHub()
	userID := uuid.New()
	sink := &recordingSink{}

	release := hub.Register(userID, sink)
	require.True(t, hub.Online(userID))

	release()
	require.False(t, hub.Online(userID))

	hub.Broadcast(userID, ws.Message{Type: "xp.gained"})
	require.Equal(t, 0, sink.count())
}

func TestHub_BroadcastToUnknownUserIsNoop(t *testing.T) {
	t.Parallel()

	hub := ws.NewHub()

	require.NotPanics(t, func() { hub.Broadcast(uuid.New(), ws.Message{Type: "level.up"}) })
}

func TestHub_ConcurrentRegisterAndBroadcast(t *testing.T) {
	t.Parallel()

	hub := ws.NewHub()
	userID := uuid.New()

	var wg sync.WaitGroup
	for range 32 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			release := hub.Register(userID, &recordingSink{})
			hub.Broadcast(userID, ws.Message{Type: "streak.updated"})
			release()
		}()
	}
	wg.Wait()

	require.False(t, hub.Online(userID))
}
