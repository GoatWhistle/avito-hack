package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/api"
	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/pagination"
)

type stubLeaderboardRead struct {
	entries    []app.LeaderboardEntry
	lastFilter app.LeaderboardFilter
	rank       int
	hasRank    bool
	rankCalls  int
}

func (s *stubLeaderboardRead) Page(_ context.Context, f app.LeaderboardFilter) ([]app.LeaderboardEntry, error) {
	s.lastFilter = f

	return s.entries, nil
}

func (s *stubLeaderboardRead) RankOf(_ context.Context, _ uuid.UUID) (rank int, found bool, err error) {
	s.rankCalls++

	return s.rank, s.hasRank, nil
}

func newLeaderboardRouter(read app.LeaderboardReadModel, actor auth.Actor) http.Handler {
	handlers := api.NewLeaderboardHandlers(api.LeaderboardDeps{
		Leaderboard: app.NewLeaderboardHandler(read),
		Authenticate: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if actor.IsZero() {
					w.WriteHeader(http.StatusUnauthorized)

					return
				}

				next.ServeHTTP(w, r.WithContext(auth.WithActor(r.Context(), actor)))
			})
		},
	})

	router := chi.NewRouter()
	handlers.RegisterRoutes(router)

	return router
}

type leaderboardBody struct {
	Items []struct {
		UserID     string `json:"user_id"`
		Name       string `json:"name"`
		Level      int    `json:"level"`
		XP         int    `json:"xp"`
		StreakDays int    `json:"streak_days"`
		Rank       int    `json:"rank"`
	} `json:"items"`
	MyRank     *int   `json:"my_rank"`
	NextCursor string `json:"next_cursor"`
}

func doLeaderboard(t *testing.T, router http.Handler, target string) (*httptest.ResponseRecorder, leaderboardBody) {
	t.Helper()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, http.NoBody))

	var body leaderboardBody
	if rec.Code == http.StatusOK {
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	}

	return rec, body
}

func TestLeaderboardContract(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	read := &stubLeaderboardRead{
		entries: []app.LeaderboardEntry{
			{UserID: userID, Name: "Alice", Level: 9, XP: 900, StreakDays: 12, Rank: 1},
			{UserID: uuid.New(), Name: "Bob", Level: 8, XP: 800, StreakDays: 3, Rank: 2},
		},
		rank:    1,
		hasRank: true,
	}

	router := newLeaderboardRouter(read, auth.Actor{ID: userID, Role: auth.RoleUser})
	rec, body := doLeaderboard(t, router, "/leaderboard?with_my_rank=true")

	require.Equal(t, http.StatusOK, rec.Code)
	require.Len(t, body.Items, 2)

	assert.Equal(t, userID.String(), body.Items[0].UserID)
	assert.Equal(t, "Alice", body.Items[0].Name)
	assert.Equal(t, 9, body.Items[0].Level)
	assert.Equal(t, 900, body.Items[0].XP)
	assert.Equal(t, 12, body.Items[0].StreakDays)
	assert.Equal(t, 1, body.Items[0].Rank)

	require.NotNil(t, body.MyRank)
	assert.Equal(t, 1, *body.MyRank)
}

func TestLeaderboardSkipsRankQueryByDefault(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	read := &stubLeaderboardRead{
		entries: []app.LeaderboardEntry{{UserID: userID, Name: "Alice", Level: 9, XP: 900, Rank: 1}},
		rank:    1,
		hasRank: true,
	}

	router := newLeaderboardRouter(read, auth.Actor{ID: userID, Role: auth.RoleUser})
	rec, body := doLeaderboard(t, router, "/leaderboard")

	require.Equal(t, http.StatusOK, rec.Code)
	require.Len(t, body.Items, 1)
	assert.Nil(t, body.MyRank)
	assert.Zero(t, read.rankCalls)
}

func TestLeaderboardMyRankNullWithoutPet(t *testing.T) {
	t.Parallel()

	router := newLeaderboardRouter(&stubLeaderboardRead{}, auth.Actor{ID: uuid.New(), Role: auth.RoleUser})
	rec, body := doLeaderboard(t, router, "/leaderboard")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Nil(t, body.MyRank)
	assert.Empty(t, body.Items)
	assert.Contains(t, rec.Body.String(), `"items":[]`)
	assert.Contains(t, rec.Body.String(), `"my_rank":null`)
}

func TestLeaderboardRequiresAuth(t *testing.T) {
	t.Parallel()

	router := newLeaderboardRouter(&stubLeaderboardRead{}, auth.Actor{})
	rec, _ := doLeaderboard(t, router, "/leaderboard")

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLeaderboardQueryParams(t *testing.T) {
	t.Parallel()

	cursor := pagination.LeaderboardCursor{Level: 5, XP: 50, UserID: uuid.New()}

	tests := []struct {
		name       string
		target     string
		wantStatus int
		wantLimit  int
		wantOffset int
		wantCursor pagination.LeaderboardCursor
	}{
		{name: "default limit", target: "/leaderboard", wantStatus: http.StatusOK, wantLimit: 21},
		{name: "explicit limit", target: "/leaderboard?limit=5", wantStatus: http.StatusOK, wantLimit: 6},
		{name: "limit clamped to max", target: "/leaderboard?limit=500", wantStatus: http.StatusOK, wantLimit: 101},
		{name: "negative limit falls back", target: "/leaderboard?limit=-3", wantStatus: http.StatusOK, wantLimit: 21},
		{name: "bad limit", target: "/leaderboard?limit=abc", wantStatus: http.StatusBadRequest},
		{name: "bad cursor", target: "/leaderboard?cursor=%21%21%21", wantStatus: http.StatusBadRequest},
		{
			name: "valid cursor", target: "/leaderboard?cursor=" + cursor.Encode(),
			wantStatus: http.StatusOK, wantLimit: 21, wantCursor: cursor,
		},
		{
			name: "around me", target: "/leaderboard?around=me",
			wantStatus: http.StatusOK, wantLimit: 7, wantOffset: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			read := &stubLeaderboardRead{rank: 10, hasRank: true}
			router := newLeaderboardRouter(read, auth.Actor{ID: uuid.New(), Role: auth.RoleUser})
			rec, _ := doLeaderboard(t, router, tt.target)

			require.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus != http.StatusOK {
				return
			}

			assert.Equal(t, tt.wantLimit, read.lastFilter.Limit)
			assert.Equal(t, tt.wantOffset, read.lastFilter.Offset)
			assert.Equal(t, tt.wantCursor, read.lastFilter.Cursor)
		})
	}
}
