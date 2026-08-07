package api_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/api"
	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
)

type stubRewardService struct {
	catalog     []app.RewardCatalogEntry
	mine        []app.GrantedReward
	code        string
	catalogErr  error
	mineErr     error
	activateErr error
	lastReward  string
}

func (s *stubRewardService) Catalog(context.Context, uuid.UUID) ([]app.RewardCatalogEntry, error) {
	if s.catalogErr != nil {
		return nil, s.catalogErr
	}

	return s.catalog, nil
}

func (s *stubRewardService) Mine(context.Context, uuid.UUID) ([]app.GrantedReward, error) {
	if s.mineErr != nil {
		return nil, s.mineErr
	}

	return s.mine, nil
}

func (s *stubRewardService) Activate(_ context.Context, _ uuid.UUID, rewardID string) (string, error) {
	s.lastReward = rewardID
	if s.activateErr != nil {
		return "", s.activateErr
	}

	return s.code, nil
}

func mustReward(t *testing.T, id string, value int) domain.Reward {
	t.Helper()

	reward, err := domain.NewReward(domain.RewardParams{
		ID: id, Title: "title " + id, Description: "desc", Kind: domain.RewardKindPromo,
		ConditionType: domain.ConditionLevel, ConditionValue: value,
	})
	require.NoError(t, err)

	return reward
}

func rewardRouter(t *testing.T, svc *stubRewardService, actor *auth.Actor) http.Handler {
	t.Helper()

	router := chi.NewRouter()
	api.NewRewardHandlers(api.RewardDeps{
		Rewards:      svc,
		Authenticate: injectActor(actor),
	}).RegisterRoutes(router)

	return router
}
