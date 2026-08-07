package infra

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
)

func summaryValues(facts, advice []byte) []any {
	return []any{
		uuid.New(), uuid.New(), scanTime, facts, "Message", advice, "llm", scanTime,
	}
}

func TestScanSummaryDecodesFactsAndAdvice(t *testing.T) {
	t.Parallel()

	itemID := uuid.New()
	facts, err := json.Marshal(domain.DayFacts{
		Date: scanTime, PetName: "Enot", TotalXP: 90,
		Level: domain.LevelFacts{Current: 4, XP: 90},
	})
	require.NoError(t, err)

	advice, err := json.Marshal(domain.Advice{
		Text: "Set a price", ItemID: &itemID, Action: domain.AdviceSetPrice,
	})
	require.NoError(t, err)

	summary, err := scanSummary(pgtest.Row{Values: summaryValues(facts, advice)})

	require.NoError(t, err)
	assert.Equal(t, "Enot", summary.Facts().PetName)
	assert.Equal(t, 90, summary.Facts().TotalXP)
	assert.Equal(t, 4, summary.Facts().Level.Current)
	assert.Equal(t, domain.SummarySourceLLM, summary.GeneratedBy())
	require.NotNil(t, summary.Advice())
	assert.Equal(t, domain.AdviceSetPrice, summary.Advice().Action)
	require.NotNil(t, summary.Advice().ItemID)
	assert.Equal(t, itemID, *summary.Advice().ItemID)
}

func TestScanSummaryEmptyAdviceStaysNil(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		advice []byte
	}{
		{name: "nil", advice: nil},
		{name: "empty slice", advice: []byte{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			summary, err := scanSummary(pgtest.Row{Values: summaryValues([]byte(`{}`), tc.advice)})

			require.NoError(t, err)
			assert.Nil(t, summary.Advice())
		})
	}
}

func TestScanSummaryRejectsMalformedJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		facts  []byte
		advice []byte
		want   string
	}{
		{name: "facts", facts: []byte(`{`), advice: nil, want: "decode summary facts"},
		{name: "empty facts", facts: nil, advice: nil, want: "decode summary facts"},
		{
			name:   "advice",
			facts:  []byte(`{}`),
			advice: []byte(`{"text":`),
			want:   "decode summary advice",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := scanSummary(pgtest.Row{Values: summaryValues(tc.facts, tc.advice)})

			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.want)
		})
	}
}

func TestScanSummaryPropagatesScanError(t *testing.T) {
	t.Parallel()

	_, err := scanSummary(pgtest.Row{Err: errScan})

	require.ErrorIs(t, err, errScan)
}

func TestEncodeAdvice(t *testing.T) {
	t.Parallel()

	encoded, err := encodeAdvice(nil)

	require.NoError(t, err)
	assert.Nil(t, encoded)

	encoded, err = encodeAdvice(&domain.Advice{Text: "Check in", Action: domain.AdviceCheckIn})

	require.NoError(t, err)
	assert.JSONEq(t, `{"text":"Check in","action":"check_in"}`, string(encoded))
}

func TestEncodeAdviceRoundTripsThroughScanSummary(t *testing.T) {
	t.Parallel()

	advice := &domain.Advice{Text: "Refresh", Action: domain.AdviceRefreshListing}

	encoded, err := encodeAdvice(advice)
	require.NoError(t, err)

	summary, err := scanSummary(pgtest.Row{Values: summaryValues([]byte(`{}`), encoded)})

	require.NoError(t, err)
	require.NotNil(t, summary.Advice())
	assert.Equal(t, *advice, *summary.Advice())
	assert.Equal(t, time.Duration(0), summary.CreatedAt().Sub(scanTime))
}
