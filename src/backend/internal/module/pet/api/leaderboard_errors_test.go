package api_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

type failingLeaderboardRead struct {
	pageErr error
	rankErr error
}

func (r failingLeaderboardRead) Page(
	context.Context, app.LeaderboardFilter,
) ([]app.LeaderboardEntry, error) {
	return nil, r.pageErr
}

func (r failingLeaderboardRead) RankOf(context.Context, uuid.UUID) (rank int, found bool, err error) {
	return 0, false, r.rankErr
}

func TestLeaderboardPropagatesPageFailure(t *testing.T) {
	t.Parallel()

	router := newLeaderboardRouter(
		failingLeaderboardRead{pageErr: errors.New("leaderboard view is being rebuilt")},
		auth.Actor{ID: uuid.New(), Role: auth.RoleUser},
	)

	rec, _ := doLeaderboard(t, router, "/leaderboard")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.NotContains(t, rec.Body.String(), "rebuilt")
}

func TestLeaderboardPropagatesRankFailure(t *testing.T) {
	t.Parallel()

	router := newLeaderboardRouter(
		failingLeaderboardRead{rankErr: errors.New("rank index missing")},
		auth.Actor{ID: uuid.New(), Role: auth.RoleUser},
	)

	rec, _ := doLeaderboard(t, router, "/leaderboard?around=me")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestLeaderboardMapsDomainErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "not found", err: domainerr.ErrNotFound, wantStatus: http.StatusNotFound},
		{name: "forbidden", err: domainerr.ErrForbidden, wantStatus: http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			router := newLeaderboardRouter(
				failingLeaderboardRead{pageErr: tt.err},
				auth.Actor{ID: uuid.New(), Role: auth.RoleUser},
			)

			rec, _ := doLeaderboard(t, router, "/leaderboard")

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
