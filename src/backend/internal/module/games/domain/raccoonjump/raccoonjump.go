package raccoonjump

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/avito-hack/backend/internal/module/games/domain"
)

const Slug = "raccoonjump"

const (
	maxCollectibles   = 4
	collectibleStride = 40
	firstCollectible  = 30
)

type Game struct {
	listings ListingPool
	clock    Clock
}

func New(listings ListingPool, clock Clock) *Game {
	return &Game{listings: listings, clock: clock}
}

func (g *Game) Slug() string { return Slug }

func (g *Game) TargetStreak() int { return 0 }

func (g *Game) Start(ctx context.Context, r *domain.Round) (domain.View, error) {
	s := state{
		StartedAtUnixMilli: g.now().UnixMilli(),
		Seed:               newSeed(),
		Collectibles:       g.drawCollectibles(ctx),
		Collected:          make([]int, 0),
	}

	return g.commit(r, s)
}

func (g *Game) Resume(_ context.Context, r *domain.Round) (domain.View, error) {
	s, err := decodeState(r.Payload())
	if err != nil {
		return domain.View{}, err
	}

	prompt, err := s.prompt()
	if err != nil {
		return domain.View{}, err
	}

	return domain.View{Prompt: prompt}, nil
}

func (g *Game) Guess(_ context.Context, r *domain.Round, move json.RawMessage) (domain.GuessOutcome, error) {
	s, err := decodeState(r.Payload())
	if err != nil {
		return domain.GuessOutcome{}, domain.ErrInvalidMove("failed to decode state")
	}

	score, collected, err := parseMove(move)
	if err != nil {
		return domain.GuessOutcome{}, err
	}

	if s.Finished {
		return g.replay(s)
	}

	elapsed := g.now().Sub(time.UnixMilli(s.StartedAtUnixMilli))

	if score > MaxPlausibleScore(elapsed) {
		return domain.GuessOutcome{}, ErrImplausibleScore
	}

	s.Finished = true
	s.Score = score
	s.Collected = validCollected(s, collected)

	if score > s.BestScore {
		s.BestScore = score
	}

	if _, err := g.commit(r, s); err != nil {
		return domain.GuessOutcome{}, err
	}

	reveal, err := s.reveal(true)
	if err != nil {
		return domain.GuessOutcome{}, err
	}

	return domain.GuessOutcome{Correct: true, Progress: domain.ProgressWin, Reveal: reveal}, nil
}

func (g *Game) replay(s state) (domain.GuessOutcome, error) {
	reveal, err := s.reveal(false)
	if err != nil {
		return domain.GuessOutcome{}, err
	}

	return domain.GuessOutcome{Correct: true, Progress: domain.ProgressWin, Reveal: reveal}, nil
}

func (s state) reveal(newBest bool) (json.RawMessage, error) {
	payload := revealPayload{
		Score:         s.Score,
		BestScore:     s.BestScore,
		NewBest:       newBest && s.Score >= s.BestScore && s.Score > 0,
		CountsToStrek: CountsTowardStreak(s.Score),
		Listings:      s.collectedListings(s.Collected),
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal raccoonjump reveal: %w", err)
	}

	return encoded, nil
}

func validCollected(s state, collected []int) []int {
	byIndex := make(map[int]bool, len(s.Collectibles))
	for _, item := range s.Collectibles {
		byIndex[item.Index] = true
	}

	seen := make(map[int]bool, len(collected))
	out := make([]int, 0, len(collected))

	for _, index := range collected {
		if seen[index] || !byIndex[index] {
			continue
		}

		if index > s.Score {
			continue
		}

		seen[index] = true

		out = append(out, index)
	}

	return out
}

func parseMove(raw json.RawMessage) (int, []int, error) {
	if len(raw) == 0 {
		return 0, nil, domain.ErrInvalidMove("move is required")
	}

	var move movePayload
	if err := json.Unmarshal(raw, &move); err != nil {
		return 0, nil, domain.ErrInvalidMove("move must be an object with a score field")
	}

	if move.Score == nil {
		return 0, nil, domain.ErrInvalidMove("score is required")
	}

	score := *move.Score

	if score < 0 {
		return 0, nil, domain.ErrInvalidMove("score must not be negative")
	}

	if score > MaxScore {
		return 0, nil, ErrImplausibleScore
	}

	if len(move.Collected) > maxCollectibles {
		return 0, nil, domain.ErrInvalidMove("too many collectibles reported")
	}

	return score, move.Collected, nil
}

func (g *Game) drawCollectibles(ctx context.Context) []collectible {
	if g.listings == nil {
		return nil
	}

	listings, err := g.listings.RandomListings(ctx, maxCollectibles)
	if err != nil {
		return nil
	}

	out := make([]collectible, 0, len(listings))

	for i, listing := range listings {
		out = append(out, collectible{
			Index:       firstCollectible + i*collectibleStride,
			DisplayID:   listing.DisplayID,
			Title:       listing.Title,
			PriceKopeks: listing.PriceKopeks,
			PhotoURL:    listing.PhotoURL,
		})
	}

	return out
}

func (g *Game) commit(r *domain.Round, s state) (domain.View, error) {
	payload, err := json.Marshal(s)
	if err != nil {
		return domain.View{}, fmt.Errorf("marshal raccoonjump state: %w", err)
	}

	prompt, err := s.prompt()
	if err != nil {
		return domain.View{}, err
	}

	r.SetPayload(payload)

	return domain.View{Prompt: prompt}, nil
}

func newSeed() uint32 {
	var buf [4]byte

	if _, err := rand.Read(buf[:]); err != nil {
		return 1
	}

	seed := binary.BigEndian.Uint32(buf[:])
	if seed == 0 {
		return 1
	}

	return seed
}

func (g *Game) now() time.Time {
	if g.clock == nil {
		return time.Now()
	}

	return g.clock.Now()
}

func (g *Game) RoundScore(r *domain.Round) int {
	s, err := decodeState(r.Payload())
	if err != nil {
		return 0
	}

	if !s.Finished {
		return 0
	}

	return s.Score
}

func (g *Game) CountsTowardStreak(score int) bool {
	return CountsTowardStreak(score)
}
