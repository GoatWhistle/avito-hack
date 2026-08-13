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
	RightPrice  int64  `json:"right_price"`
	RightItemID string `json:"right_item_id"`
	RightTitle  string `json:"right_title"`
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

func HiddenPhotoURL(payload json.RawMessage, displayID string) (string, bool) {
	current, err := decodeState(payload)
	if err != nil {
		return "", false
	}

	if current.Right.DisplayID != displayID || current.Right.PhotoURL == "" {
		return "", false
	}

	return current.Right.PhotoURL, true
}

func hiddenPhotoURL(token, original string) string {
	if original == "" {
		return ""
	}

	return "/api/v1/games/photo/" + token
}

func (s state) prompt(roundID string, signer PhotoSigner) (json.RawMessage, error) {
	token := signer.Sign(roundID, s.Right.DisplayID)

	payload := promptPayload{
		Left: knownSide{
			ItemID:      s.Left.DisplayID,
			Title:       s.Left.Title,
			PhotoURL:    s.Left.PhotoURL,
			PriceKopeks: s.Left.PriceKopeks,
		},
		Right: hiddenSide{
			ItemID:   token,
			Title:    "",
			PhotoURL: hiddenPhotoURL(token, s.Right.PhotoURL),
		},
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal moreless prompt: %w", err)
	}

	return encoded, nil
}
