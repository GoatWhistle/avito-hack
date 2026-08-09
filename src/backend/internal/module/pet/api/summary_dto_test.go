package api

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/domain"
)

func TestOrEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{name: "nil becomes empty slice", input: nil, want: []string{}},
		{name: "empty stays empty", input: []string{}, want: []string{}},
		{name: "values are preserved", input: []string{"a", "b"}, want: []string{"a", "b"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := orEmpty(tc.input)

			require.NotNil(t, got)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestToAdvicePayloadNil(t *testing.T) {
	t.Parallel()

	assert.Nil(t, toAdvicePayload(nil))
}

func TestToAdvicePayloadPopulated(t *testing.T) {
	t.Parallel()

	itemID := uuid.New()
	advice := &domain.Advice{
		Text: "Добавь фото", ItemID: &itemID,
		ItemTitle: "Велосипед", Action: domain.AdviceAddPhoto,
	}

	got := toAdvicePayload(advice)

	require.NotNil(t, got)
	assert.Equal(t, "Добавь фото", got.Text)
	assert.Equal(t, "add_photo", got.Action)
	assert.Equal(t, "Велосипед", got.ItemTitle)
	require.NotNil(t, got.ItemID)
	assert.Equal(t, itemID, *got.ItemID)
}

func TestToAdvicePayloadWithoutItemID(t *testing.T) {
	t.Parallel()

	got := toAdvicePayload(&domain.Advice{Text: "Заходи завтра", Action: domain.AdviceCheckIn})

	require.NotNil(t, got)
	assert.Nil(t, got.ItemID)
	assert.Equal(t, "check_in", got.Action)
}

func TestToFactsPayloadZeroValue(t *testing.T) {
	t.Parallel()

	got := toFactsPayload(domain.DayFacts{})

	assert.Equal(t, 0, got.TotalXP)
	assert.Equal(t, []summaryActionPayload{}, got.Actions)
	assert.Equal(t, []string{}, got.Rewards)
	assert.Equal(t, []string{}, got.Badges)
	assert.Equal(t, 0, got.IssuesCount)
	assert.Nil(t, got.Rank)
}

func TestToFactsPayloadPopulated(t *testing.T) {
	t.Parallel()

	facts := domain.DayFacts{
		TotalXP: 80,
		Actions: []domain.XPByAction{
			{Action: domain.ActionDailyCheckIn, Count: 1, Amount: 25},
			{Action: domain.ActionFavorite, Count: 3, Amount: 15},
		},
		Level:     domain.LevelFacts{Previous: 2, Current: 3, XP: 80, NextLevelXP: 100, XPToNext: 20},
		LeveledUp: true,
		Rewards:   []string{"r1"},
		Badges:    []string{"b1", "b2"},
		Pet: domain.PetSnapshot{
			Stage: domain.StageTeen, State: domain.StateHappy,
			Satiety: 60, Happiness: 75, Energy: 85,
		},
		Streak:      domain.StreakFacts{Days: 5, Broken: true},
		Leaderboard: domain.LeaderboardFacts{Rank: 9, Previous: 14, Known: true},
		Issues: []domain.ListingIssue{
			{ItemID: uuid.New(), Title: "Bike", Kind: domain.IssueNoPhoto},
			{ItemID: uuid.New(), Title: "Chair", Kind: domain.IssueStale},
		},
	}

	got := toFactsPayload(facts)

	assert.Equal(t, 80, got.TotalXP)
	require.Len(t, got.Actions, 2)
	assert.Equal(t, "daily_checkin", got.Actions[0].Action)
	assert.Equal(t, 3, got.Actions[1].Count)
	assert.Equal(t, 15, got.Actions[1].Amount)
	assert.Equal(t, 3, got.Level)
	assert.Equal(t, 2, got.PrevLevel)
	assert.True(t, got.LeveledUp)
	assert.Equal(t, 20, got.XPToNext)
	assert.Equal(t, []string{"r1"}, got.Rewards)
	assert.Equal(t, []string{"b1", "b2"}, got.Badges)
	assert.Equal(t, domain.StageTeen, got.Stage)
	assert.Equal(t, 5, got.StreakDays)
	assert.True(t, got.StreakBroke)
	assert.Equal(t, 2, got.IssuesCount)

	require.NotNil(t, got.Rank)
	assert.Equal(t, 9, *got.Rank)
}

func TestToFactsPayloadUnknownLeaderboardRankIsNil(t *testing.T) {
	t.Parallel()

	facts := domain.DayFacts{Leaderboard: domain.LeaderboardFacts{Rank: 3, Known: false}}

	assert.Nil(t, toFactsPayload(facts).Rank)
}

func TestToSummaryPayload(t *testing.T) {
	t.Parallel()

	created := time.Date(2026, time.March, 3, 9, 0, 0, 0, time.UTC)
	summary := domain.RestoreDailySummary(domain.RestoreSummaryParams{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		Date:        time.Date(2026, time.March, 2, 0, 0, 0, 0, time.UTC),
		Facts:       domain.DayFacts{TotalXP: 10},
		Message:     "Хороший день",
		GeneratedBy: domain.SummarySourceLLM,
		CreatedAt:   created,
	})

	got := toSummaryPayload(summary)

	assert.Equal(t, summary.ID().String(), got.ID)
	assert.Equal(t, "2026-03-02", got.Date)
	assert.Equal(t, "Хороший день", got.Message)
	assert.Equal(t, domain.SummarySourceLLM, got.GeneratedBy)
	assert.Equal(t, created, got.CreatedAt)
	assert.Nil(t, got.Advice)
	assert.Equal(t, 10, got.Facts.TotalXP)
}
