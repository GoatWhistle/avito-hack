package ws

import (
	"sync"

	"github.com/google/uuid"
)

type Message struct {
	Type    string `json:"type"`
	Payload any    `json:"payload,omitempty"`
}

type Sink interface {
	Send(m Message)
}

type Hub struct {
	mu    sync.RWMutex
	conns map[uuid.UUID]map[Sink]struct{}
}

func NewHub() *Hub {
	return &Hub{conns: make(map[uuid.UUID]map[Sink]struct{})}
}

func (h *Hub) Register(userID uuid.UUID, s Sink) func() {
	if s == nil || userID == uuid.Nil {
		return func() {}
	}

	h.mu.Lock()
	if h.conns[userID] == nil {
		h.conns[userID] = make(map[Sink]struct{})
	}
	h.conns[userID][s] = struct{}{}
	h.mu.Unlock()

	return func() { h.unregister(userID, s) }
}

func (h *Hub) unregister(userID uuid.UUID, s Sink) {
	h.mu.Lock()
	defer h.mu.Unlock()

	sinks, ok := h.conns[userID]
	if !ok {
		return
	}

	delete(sinks, s)

	if len(sinks) == 0 {
		delete(h.conns, userID)
	}
}

func (h *Hub) Broadcast(userID uuid.UUID, m Message) {
	for _, s := range h.sinks(userID) {
		s.Send(m)
	}
}

func (h *Hub) sinks(userID uuid.UUID) []Sink {
	h.mu.RLock()
	defer h.mu.RUnlock()

	sinks := make([]Sink, 0, len(h.conns[userID]))
	for s := range h.conns[userID] {
		sinks = append(sinks, s)
	}

	return sinks
}
