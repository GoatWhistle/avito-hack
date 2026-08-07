package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/api"
	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

type stubSummaryService struct {
	summary    *domain.DailySummary
	result     app.SummaryHistoryResult
	todayErr   error
	historyErr error
	lastQuery  app.SummaryHistoryQuery
}

func (s *stubSummaryService) Today(context.Context, uuid.UUID) (*domain.DailySummary, error) {
	if s.todayErr != nil {
		return nil, s.todayErr
	}

	return s.summary, nil
}

func (s *stubSummaryService) History(
	_ context.Context,
	q app.SummaryHistoryQuery,
) (app.SummaryHistoryResult, error) {
	s.lastQuery = q

	if s.historyErr != nil {
		return app.SummaryHistoryResult{}, s.historyErr
	}

	return s.result, nil
}

func summaryRouter(t *testing.T, svc *stubSummaryService, actor *auth.Actor) http.Handler {
	t.Helper()

	router := chi.NewRouter()
	api.NewSummaryHandlers(api.SummaryDeps{
		Summaries:    svc,
		Authenticate: injectActor(actor),
	}).RegisterRoutes(router)

	return router
}

func testSummary(t *testing.T, userID uuid.UUID) *domain.DailySummary {
	t.Helper()

	itemID := uuid.New()

	return domain.RestoreDailySummary(domain.RestoreSummaryParams{
		ID:     uuid.New(),
		UserID: userID,
		Date:   time.Date(2026, time.March, 2, 0, 0, 0, 0, time.UTC),
		Facts: domain.DayFacts{
			TotalXP: 55,
			Actions: []domain.XPByAction{{Action: domain.ActionDailyCheckIn, Count: 1, Amount: 25}},
			Level:   domain.LevelFacts{Previous: 2, Current: 3, XP: 55, NextLevelXP: 100, XPToNext: 45},
			Rewards: []string{"reward-1"},
			Badges:  []string{"badge-1"},
			Pet: domain.PetSnapshot{
				Stage: domain.StageBaby, State: domain.StateHappy,
				Satiety: 70, Happiness: 80, Energy: 90,
			},
			Streak:      domain.StreakFacts{Days: 4, Broken: false, CheckedIn: true},
			Leaderboard: domain.LeaderboardFacts{Rank: 12, Previous: 20, Known: true},
			Issues:      []domain.ListingIssue{{ItemID: itemID, Title: "Bike", Kind: domain.IssueNoPhoto}},
		},
		Message:     "Отличный день",
		Advice:      &domain.Advice{Text: "Добавь фото", ItemID: &itemID, Action: domain.AdviceAddPhoto},
		GeneratedBy: domain.SummarySourceTemplate,
		CreatedAt:   petTestTime,
	})
}

func TestSummaryTodayReturnsPayload(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	svc := &stubSummaryService{summary: testSummary(t, userID)}

	rec := doRequest(t, summaryRouter(t, svc, &actor), http.MethodGet, "/summary/today")

	require.Equal(t, http.StatusOK, rec.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	assert.Equal(t, "2026-03-02", body["date"])
	assert.Equal(t, "Отличный день", body["message"])
	assert.Equal(t, "template", body["generated_by"])

	advice, ok := body["advice"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "Добавь фото", advice["text"])
	assert.Equal(t, "add_photo", advice["action"])

	facts, ok := body["facts"].(map[string]any)
	require.True(t, ok)
	assert.InDelta(t, 55, facts["total_xp"], 0)
	assert.InDelta(t, 12, facts["leaderboard_rank"], 0)
	assert.InDelta(t, 1, facts["issues_count"], 0)
}

func TestSummaryTodayRequiresActor(t *testing.T) {
	t.Parallel()

	svc := &stubSummaryService{summary: testSummary(t, uuid.New())}

	rec := doRequest(t, summaryRouter(t, svc, nil), http.MethodGet, "/summary/today")

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "unauthorized")
}

func TestSummaryTodayMapsServiceErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "not found", err: domainerr.ErrNotFound, wantStatus: http.StatusNotFound},
		{name: "unexpected", err: errors.New("db is down"), wantStatus: http.StatusInternalServerError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
			svc := &stubSummaryService{todayErr: tc.err}

			rec := doRequest(t, summaryRouter(t, svc, &actor), http.MethodGet, "/summary/today")

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.NotContains(t, rec.Body.String(), "db is down")
		})
	}
}

func TestSummaryHistoryReturnsList(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	svc := &stubSummaryService{result: app.SummaryHistoryResult{
		Items:      []*domain.DailySummary{testSummary(t, userID), testSummary(t, userID)},
		NextCursor: "next-page-token",
	}}

	rec := doRequest(t, summaryRouter(t, svc, &actor), http.MethodGet, "/summary/history?limit=5")

	require.Equal(t, http.StatusOK, rec.Code)

	var body struct {
		Items []struct {
			Date    string `json:"date"`
			Message string `json:"message"`
		} `json:"items"`
		NextCursor string `json:"next_cursor"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	require.Len(t, body.Items, 2)
	assert.Equal(t, "2026-03-02", body.Items[0].Date)
	assert.Equal(t, "next-page-token", body.NextCursor)
	assert.Equal(t, userID, svc.lastQuery.UserID)
	assert.Equal(t, 5, svc.lastQuery.Limit)
}

func TestSummaryHistoryEmptyReturnsEmptyArray(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
	svc := &stubSummaryService{}

	rec := doRequest(t, summaryRouter(t, svc, &actor), http.MethodGet, "/summary/history")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"items":[]`)
}

func TestSummaryHistoryRejectsBadQueryParams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		target string
	}{
		{name: "non numeric limit", target: "/summary/history?limit=many"},
		{name: "malformed cursor", target: "/summary/history?cursor=not-base64!!"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}

			rec := doRequest(t, summaryRouter(t, &stubSummaryService{}, &actor), http.MethodGet, tc.target)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "validation_error")
		})
	}
}

func TestSummaryHistoryRequiresActor(t *testing.T) {
	t.Parallel()

	rec := doRequest(t, summaryRouter(t, &stubSummaryService{}, nil), http.MethodGet, "/summary/history")

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestSummaryHistoryPropagatesServiceFailure(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
	svc := &stubSummaryService{historyErr: errors.New("read model exploded")}

	rec := doRequest(t, summaryRouter(t, svc, &actor), http.MethodGet, "/summary/history")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.NotContains(t, rec.Body.String(), "read model exploded")
}
