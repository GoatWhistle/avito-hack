package moreless

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/avito-hack/backend/internal/module/games/domain"
)

const (
	Slug         = "moreless"
	targetStreak = 7

	choiceHigher = "higher"
	choiceLower  = "lower"

	nearTieRatio    = 0.02
	maxDrawTries    = 3
	maxOpeningTries = 3
)

type priceBand struct {
	min float64
	max float64
}

var priceBands = []priceBand{
	{min: 0.5, max: 2.0},
	{min: 0.33, max: 3.0},
	{min: 0.2, max: 5.0},
}

type Item struct {
	DisplayID   string
	Title       string
	PhotoURL    string
	PriceKopeks int64
}

type ItemPool interface {
	Random(ctx context.Context, exclude []string) (Item, error)
	RandomNear(ctx context.Context, seen []string, referenceKopeks int64, minRatio, maxRatio float64) (Item, error)
}

type Game struct {
	pool ItemPool
}

func New(pool ItemPool) *Game {
	return &Game{pool: pool}
}

func (g *Game) Slug() string { return Slug }

func (g *Game) TargetStreak() int { return targetStreak }

func (g *Game) Start(ctx context.Context, r *domain.Round) (domain.View, error) {
	var left, right Item

	for attempt := 0; attempt < maxOpeningTries; attempt++ {
		candidate, err := g.pool.Random(ctx, nil)
		if err != nil {
			return domain.View{}, err
		}

		partner, near, err := g.drawNear(ctx, []string{candidate.DisplayID}, candidate.PriceKopeks)
		if err != nil {
			return domain.View{}, err
		}

		left, right = candidate, partner

		if near {
			break
		}
	}

	if left.DisplayID == "" || right.DisplayID == "" {
		return domain.View{}, domain.ErrNoItemsInPool
	}

	state := state{
		Left:  left,
		Right: right,
		Seen:  []string{left.DisplayID, right.DisplayID},
	}

	return g.commit(r, state)
}

func (g *Game) Resume(_ context.Context, r *domain.Round) (domain.View, error) {
	current, err := decodeState(r.Payload())
	if err != nil {
		return domain.View{}, err
	}

	prompt, err := current.prompt()
	if err != nil {
		return domain.View{}, err
	}

	return domain.View{Prompt: prompt}, nil
}

func (g *Game) Guess(
	ctx context.Context,
	r *domain.Round,
	move json.RawMessage,
) (domain.GuessOutcome, error) {
	choice, err := parseMove(move)
	if err != nil {
		return domain.GuessOutcome{}, err
	}

	current, err := decodeState(r.Payload())
	if err != nil {
		return domain.GuessOutcome{}, err
	}

	correct := isCorrect(choice, current.Left.PriceKopeks, current.Right.PriceKopeks)

	reveal, err := json.Marshal(revealPayload{RightPrice: current.Right.PriceKopeks})
	if err != nil {
		return domain.GuessOutcome{}, fmt.Errorf("marshal reveal: %w", err)
	}

	if !correct {
		return domain.GuessOutcome{Correct: false, Reveal: reveal}, nil
	}

	if r.Streak()+1 >= targetStreak {
		return domain.GuessOutcome{Correct: true, Reveal: reveal}, nil
	}

	next, err := g.advance(ctx, r, current)
	if err != nil {
		return domain.GuessOutcome{}, err
	}

	return domain.GuessOutcome{Correct: true, Reveal: reveal, Next: &next}, nil
}

func (g *Game) advance(ctx context.Context, r *domain.Round, current state) (domain.View, error) {
	fresh, _, err := g.drawNear(ctx, current.Seen, current.Right.PriceKopeks)
	if err != nil {
		return domain.View{}, err
	}

	next := state{
		Left:  current.Right,
		Right: fresh,
		Seen:  append(append([]string(nil), current.Seen...), fresh.DisplayID),
	}

	return g.commit(r, next)
}

func (g *Game) drawNear(ctx context.Context, seen []string, reference int64) (Item, bool, error) {
	if reference <= 0 {
		item, err := g.pool.Random(ctx, seen)

		return item, false, err
	}

	var nearTie Item

	tries := 0

	for _, band := range priceBands {
		if tries >= maxDrawTries {
			break
		}

		tries++

		candidate, err := g.pool.RandomNear(ctx, seen, reference, band.min, band.max)
		if errors.Is(err, domain.ErrNoItemsInPool) {
			continue
		}
		if err != nil {
			return Item{}, false, err
		}

		if !isNearTie(candidate.PriceKopeks, reference) {
			return candidate, true, nil
		}

		if nearTie.DisplayID == "" {
			nearTie = candidate
		}
	}

	if nearTie.DisplayID != "" {
		return nearTie, true, nil
	}

	item, err := g.pool.Random(ctx, seen)

	return item, false, err
}

func isNearTie(price, reference int64) bool {
	if price == reference {
		return false
	}

	delta := price - reference
	if delta < 0 {
		delta = -delta
	}

	return float64(delta) <= float64(reference)*nearTieRatio
}

func (g *Game) commit(r *domain.Round, s state) (domain.View, error) {
	payload, err := json.Marshal(s)
	if err != nil {
		return domain.View{}, fmt.Errorf("marshal moreless state: %w", err)
	}

	prompt, err := s.prompt()
	if err != nil {
		return domain.View{}, err
	}

	r.SetPayload(payload)

	return domain.View{Prompt: prompt}, nil
}

func isCorrect(choice string, leftPrice, rightPrice int64) bool {
	if rightPrice == leftPrice {
		return true
	}

	if choice == choiceHigher {
		return rightPrice > leftPrice
	}

	return rightPrice < leftPrice
}

func parseMove(raw json.RawMessage) (string, error) {
	if len(raw) == 0 {
		return "", domain.ErrInvalidMove("move is required")
	}

	var move movePayload
	if err := json.Unmarshal(raw, &move); err != nil {
		return "", domain.ErrInvalidMove("move must be an object with a choice field")
	}

	if move.Choice != choiceHigher && move.Choice != choiceLower {
		return "", domain.ErrInvalidMove("choice must be higher or lower")
	}

	return move.Choice, nil
}
