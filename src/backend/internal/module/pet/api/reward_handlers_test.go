package api_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
)

func TestRewardCatalogEndpoint(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
	svc := &stubRewardService{catalog: []app.RewardCatalogEntry{
		{
			Reward: mustReward(t, "r1", 1), Unlocked: true, Claimed: true,
			Status: domain.RewardGranted, Current: 3, Target: 1,
		},
		{Reward: mustReward(t, "r9", 9), Current: 3, Target: 9},
	}}

	rec := doRequest(t, rewardRouter(t, svc, &actor), http.MethodGet, "/rewards")

	require.Equal(t, http.StatusOK, rec.Code)

	var body struct {
		Items []struct {
			ID       string `json:"id"`
			Title    string `json:"title"`
			Kind     string `json:"kind"`
			Unlocked bool   `json:"unlocked"`
			Claimed  bool   `json:"claimed"`
			Status   string `json:"status"`
			Current  int    `json:"progress_current"`
			Target   int    `json:"progress_target"`
		} `json:"items"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Len(t, body.Items, 2)

	assert.Equal(t, "r1", body.Items[0].ID)
	assert.Equal(t, "title r1", body.Items[0].Title)
	assert.Equal(t, string(domain.RewardKindPromo), body.Items[0].Kind)
	assert.True(t, body.Items[0].Unlocked)
	assert.True(t, body.Items[0].Claimed)
	assert.Equal(t, string(domain.RewardGranted), body.Items[0].Status)

	assert.False(t, body.Items[1].Unlocked)
	assert.Equal(t, 9, body.Items[1].Target)
	assert.Empty(t, body.Items[1].Status)
}

func TestRewardCatalogEmptyReturnsEmptyArray(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}

	rec := doRequest(t, rewardRouter(t, &stubRewardService{}, &actor), http.MethodGet, "/rewards")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"items":[]`)
}

func TestRewardMineEndpoint(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	grant, err := domain.GrantReward(userID, "r1", petTestTime)
	require.NoError(t, err)

	svc := &stubRewardService{mine: []app.GrantedReward{
		{Reward: mustReward(t, "r1", 1), Grant: grant},
	}}

	rec := doRequest(t, rewardRouter(t, svc, &actor), http.MethodGet, "/rewards/my")

	require.Equal(t, http.StatusOK, rec.Code)

	var body struct {
		Items []struct {
			RewardID string `json:"reward_id"`
			Title    string `json:"title"`
			Status   string `json:"status"`
		} `json:"items"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Len(t, body.Items, 1)
	assert.Equal(t, "r1", body.Items[0].RewardID)
	assert.Equal(t, "title r1", body.Items[0].Title)
	assert.Equal(t, string(domain.RewardGranted), body.Items[0].Status)
}

func TestRewardActivateReturnsCode(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
	svc := &stubRewardService{code: "PROMO-77"}

	rec := doRequest(t, rewardRouter(t, svc, &actor), http.MethodPost, "/rewards/r1/activate")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "r1", svc.lastReward)

	var body struct {
		RewardID string `json:"reward_id"`
		Code     string `json:"code"`
		Status   string `json:"status"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "r1", body.RewardID)
	assert.Equal(t, "PROMO-77", body.Code)
	assert.Equal(t, string(domain.RewardActivated), body.Status)
}

func TestRewardActivateMapsDomainErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{"not granted", domain.ErrRewardNotGranted, http.StatusBadRequest, "validation_error"},
		{"already activated", domain.ErrRewardAlreadyActivated, http.StatusConflict, "conflict"},
		{"not activatable", domain.ErrRewardNotActivatable, http.StatusConflict, "conflict"},
		{"expired", domain.ErrRewardExpired, http.StatusConflict, "conflict"},
		{"locked", domain.ErrRewardLocked, http.StatusBadRequest, "validation_error"},
		{"unexpected", errors.New("boom"), http.StatusInternalServerError, "internal_error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
			svc := &stubRewardService{activateErr: tt.err}

			rec := doRequest(t, rewardRouter(t, svc, &actor), http.MethodPost, "/rewards/r1/activate")

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.wantCode)
		})
	}
}

func TestRewardActivateHidesInternalDetails(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
	svc := &stubRewardService{activateErr: errors.New("pgx: connection refused on 10.0.0.5")}

	rec := doRequest(t, rewardRouter(t, svc, &actor), http.MethodPost, "/rewards/r1/activate")

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.NotContains(t, rec.Body.String(), "10.0.0.5")
	assert.NotContains(t, rec.Body.String(), "pgx")
}

func TestRewardEndpointsRequireActor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/rewards"},
		{http.MethodGet, "/rewards/my"},
		{http.MethodPost, "/rewards/r1/activate"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			t.Parallel()

			rec := doRequest(t, rewardRouter(t, &stubRewardService{}, nil), tt.method, tt.path)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		})
	}
}

func TestRewardCatalogPropagatesFailure(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
	svc := &stubRewardService{catalogErr: errors.New("db down")}

	rec := doRequest(t, rewardRouter(t, svc, &actor), http.MethodGet, "/rewards")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
