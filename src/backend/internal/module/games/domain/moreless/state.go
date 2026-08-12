package moreless

import (
	"encoding/json"
	"fmt"

	"github.com/avito-hack/backend/internal/module/games/domain"
)

type state struct {
	Left  Item     `json:"left"`
	Right Item     `json:"right"`
	Seen  []string `json:"seen"`
}

type movePayload struct {
	Choice string `json:"choice"`
}

type revealPayload struct {
	RightPrice int64 `json:"right_price"`
}

type knownSide struct {
	ItemID      string `json:"item_id"`
	Title       string `json:"title"`
	PhotoURL    string `json:"photo_url"`
	PriceKopeks int64  `json:"price"`
}

type hiddenSide struct {
	ItemID   string `json:"item_id"`
	Title    string `json:"title"`
	PhotoURL string `json:"photo_url"`
}

type promptPayload struct {
	Left  knownSide  `json:"left"`
	Right hiddenSide `json:"right"`
}

func decodeState(raw json.RawMessage) (state, error) {
	var s state
	if err := json.Unmarshal(raw, &s); err != nil {
		return state{}, fmt.Errorf("decode moreless state: %w", err)
	}

	if s.Left.DisplayID == "" || s.Right.DisplayID == "" {
		return state{}, domain.ErrRoundNotFound
	}

	return s, nil
}

func (s state) prompt() (json.RawMessage, error) {
	payload := promptPayload{
		Left: knownSide{
			ItemID:      s.Left.DisplayID,
			Title:       s.Left.Title,
			PhotoURL:    s.Left.PhotoURL,
			PriceKopeks: s.Left.PriceKopeks,
		},
		Right: hiddenSide{
			ItemID:   s.Right.DisplayID,
			Title:    s.Right.Title,
			PhotoURL: s.Right.PhotoURL,
		},
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal moreless prompt: %w", err)
	}

	return encoded, nil
}
