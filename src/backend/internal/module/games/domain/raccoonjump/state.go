package raccoonjump

import (
	"encoding/json"
	"fmt"

	"github.com/avito-hack/backend/internal/module/games/domain"
)

type collectible struct {
	Index       int    `json:"index"`
	DisplayID   string `json:"display_id"`
	Title       string `json:"title"`
	PriceKopeks int64  `json:"price_kopeks"`
	PhotoURL    string `json:"photo_url"`
}

type state struct {
	StartedAtUnixMilli int64         `json:"started_at_unix_milli"`
	Seed               uint32        `json:"seed"`
	Collectibles       []collectible `json:"collectibles"`
	Finished           bool          `json:"finished"`
	Score              int           `json:"score"`
	BestScore          int           `json:"best_score"`
	Collected          []int         `json:"collected"`
}

type movePayload struct {
	Score     *int  `json:"score"`
	Collected []int `json:"collected"`
}

type promptCollectible struct {
	Index    int    `json:"index"`
	Title    string `json:"title"`
	PhotoURL string `json:"photo_url"`
}

type promptPayload struct {
	Seed            uint32              `json:"seed"`
	MaxScore        int                 `json:"max_score"`
	MinStreakScore  int                 `json:"min_streak_score"`
	ScorePerSecond  int                 `json:"max_score_per_second"`
	Collectibles    []promptCollectible `json:"collectibles"`
	BestScore       int                 `json:"best_score"`
	StartedAtMillis int64               `json:"started_at_unix_milli"`
}

type listingPayload struct {
	DisplayID   string `json:"display_id"`
	Title       string `json:"title"`
	PriceKopeks int64  `json:"price_kopeks"`
	PhotoURL    string `json:"photo_url"`
}

type revealPayload struct {
	Score         int              `json:"score"`
	BestScore     int              `json:"best_score"`
	NewBest       bool             `json:"new_best"`
	CountsToStrek bool             `json:"counts_toward_streak"`
	Listings      []listingPayload `json:"listings,omitempty"`
}

func decodeState(raw json.RawMessage) (state, error) {
	var s state
	if err := json.Unmarshal(raw, &s); err != nil {
		return state{}, fmt.Errorf("decode raccoonjump state: %w", err)
	}

	if s.StartedAtUnixMilli == 0 {
		return state{}, domain.ErrRoundNotFound
	}

	return s, nil
}

func (s state) prompt() (json.RawMessage, error) {
	collectibles := make([]promptCollectible, 0, len(s.Collectibles))

	for _, item := range s.Collectibles {
		collectibles = append(collectibles, promptCollectible{
			Index:    item.Index,
			Title:    item.Title,
			PhotoURL: item.PhotoURL,
		})
	}

	payload := promptPayload{
		Seed:            s.Seed,
		MaxScore:        MaxScore,
		MinStreakScore:  MinStreakScore,
		ScorePerSecond:  maxScorePerSecond,
		Collectibles:    collectibles,
		BestScore:       s.BestScore,
		StartedAtMillis: s.StartedAtUnixMilli,
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal raccoonjump prompt: %w", err)
	}

	return encoded, nil
}

func (s state) collectedListings(indexes []int) []listingPayload {
	byIndex := make(map[int]collectible, len(s.Collectibles))
	for _, item := range s.Collectibles {
		byIndex[item.Index] = item
	}

	seen := make(map[int]bool, len(indexes))
	out := make([]listingPayload, 0, len(indexes))

	for _, index := range indexes {
		if seen[index] {
			continue
		}

		seen[index] = true

		item, ok := byIndex[index]
		if !ok {
			continue
		}

		out = append(out, listingPayload{
			DisplayID:   item.DisplayID,
			Title:       item.Title,
			PriceKopeks: item.PriceKopeks,
			PhotoURL:    item.PhotoURL,
		})
	}

	return out
}
