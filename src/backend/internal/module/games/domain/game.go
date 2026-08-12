package domain

import (
	"context"
	"encoding/json"
)

type View struct {
	Prompt json.RawMessage
}

type GuessOutcome struct {
	Correct bool
	Reveal  json.RawMessage
	Next    *View
}

type Game interface {
	Slug() string
	TargetStreak() int
	Start(ctx context.Context, r *Round) (View, error)
	Guess(ctx context.Context, r *Round, move json.RawMessage) (GuessOutcome, error)
	Resume(ctx context.Context, r *Round) (View, error)
}

type Registry struct {
	games map[string]Game
	order []string
}

func NewRegistry(games ...Game) *Registry {
	registry := &Registry{games: make(map[string]Game, len(games))}

	for _, game := range games {
		registry.Register(game)
	}

	return registry
}

func (r *Registry) Register(game Game) {
	if game == nil {
		return
	}

	slug := game.Slug()
	if slug == "" {
		return
	}

	if _, exists := r.games[slug]; !exists {
		r.order = append(r.order, slug)
	}

	r.games[slug] = game
}

func (r *Registry) Get(slug string) (Game, error) {
	game, ok := r.games[slug]
	if !ok {
		return nil, ErrGameNotFound
	}

	return game, nil
}

func (r *Registry) All() []Game {
	out := make([]Game, 0, len(r.order))

	for _, slug := range r.order {
		out = append(out, r.games[slug])
	}

	return out
}

func (r *Registry) Slugs() []string {
	out := make([]string, len(r.order))
	copy(out, r.order)

	return out
}
