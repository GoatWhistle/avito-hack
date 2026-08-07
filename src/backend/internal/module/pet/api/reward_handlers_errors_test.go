package api_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/domain"
)

func TestRewardMinePropagatesServiceError(t *testing.T) {
	t.Parallel()

	actor := testActor()
	svc := &stubRewardService{mineErr: errors.New("grants table unavailable")}

	rec := doRequest(t, rewardRouter(t, svc, &actor), http.MethodGet, "/rewards/my")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.NotContains(t, rec.Body.String(), "grants table")
}

func TestRewardMineMapsDomainErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "duplicate action", err: domain.ErrDuplicateAction, wantStatus: http.StatusConflict},
		{name: "limit reached", err: domain.ErrLimitReached, wantStatus: http.StatusConflict},
		{name: "condition not met", err: domain.ErrConditionNotMet, wantStatus: http.StatusBadRequest},
		{name: "invalid action", err: domain.ErrInvalidAction, wantStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actor := testActor()
			svc := &stubRewardService{mineErr: tt.err}

			rec := doRequest(t, rewardRouter(t, svc, &actor), http.MethodGet, "/rewards/my")

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestRewardMineReturnsEmptyList(t *testing.T) {
	t.Parallel()

	actor := testActor()

	rec := doRequest(t, rewardRouter(t, &stubRewardService{}, &actor), http.MethodGet, "/rewards/my")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"items":[]`)
}
