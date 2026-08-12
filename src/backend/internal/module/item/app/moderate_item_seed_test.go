package app_test

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/events"
	"github.com/avito-hack/backend/internal/shared/vo"
)

func seedItemWithStatus(t *testing.T, ownerID uuid.UUID, status domain.Status) *domain.Item {
	t.Helper()

	return domain.RestoreItem(domain.RestoreItemParams{
		ID: uuid.New(), DisplayID: domain.NewDisplayID(), OwnerID: ownerID, Title: "Seed chair",
		Description: "provided by platform",
		Price:       vo.MustMoney(1000), Status: status, Attributes: domain.NewAttributes(nil),
		CreatedAt: fixedTime, UpdatedAt: fixedTime, IsSeed: true, AIVerified: true,
	})
}

func TestAddPhotoRejectsSeedItem(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := seedItemWithStatus(t, ownerID, domain.StatusDraft)
	repo := &stubRepository{item: item}
	photos := &countingPhotos{}
	storage := &stubStorage{}

	handler := app.NewAddPhotoHandler(repo, photos, storage, passthroughTx{}, fakeClock{}, nil)

	_, err := handler.Handle(t.Context(), app.AddPhotoCommand{
		ItemID: item.ID(), ActorID: ownerID,
		Content: strings.NewReader("binary"), ContentType: "image/jpeg",
	})

	require.ErrorIs(t, err, domainerr.ErrConflict)
	assert.Nil(t, repo.saved)
}

func TestAddPhotoTriggersModerationCycleOnPublishedItem(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusPublished)
	item.MarkAIVerified()
	repo := &stubRepository{item: item}
	photos := &countingPhotos{}
	storage := &stubStorage{}

	moderate := app.NewModerateItemHandler(
		repo, photos,
		&stubModerationProvider{result: domain.ModerationResult{Verdict: domain.ModerationApproved, Provider: "test"}},
		&stubModerationLog{}, stubPhotoBytesLoader{}, passthroughTx{}, fakeClock{}, events.NopPublisher{},
	)

	handler := app.NewAddPhotoHandler(repo, photos, storage, passthroughTx{}, fakeClock{}, moderate)

	_, err := handler.Handle(t.Context(), app.AddPhotoCommand{
		ItemID: item.ID(), ActorID: ownerID,
		Content: strings.NewReader("binary"), ContentType: "image/jpeg",
	})

	require.NoError(t, err)
	require.NotNil(t, repo.saved)
	assert.Equal(t, domain.StatusPublished, repo.saved.Status())
	assert.True(t, repo.saved.AIVerified())
}

func TestDeletePhotoTriggersReverificationOnPublishedItem(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusPublished)
	item.MarkAIVerified()
	repo := &stubRepository{item: item}
	photos := &countingPhotos{}
	storage := &stubStorage{}

	handler := app.NewDeletePhotoHandler(repo, photos, storage, passthroughTx{}, fakeClock{}, nil)

	err := handler.Handle(t.Context(), app.DeletePhotoCommand{
		ItemID: item.ID(), PhotoDisplayID: domain.NewDisplayID(), ActorID: ownerID,
	})

	require.NoError(t, err)
	require.NotNil(t, repo.saved)
	assert.Equal(t, domain.StatusModeration, repo.saved.Status())
	assert.False(t, repo.saved.AIVerified())
}

func TestUpdateItemHandlerTriggersSynchronousReverification(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusPublished)
	item.MarkAIVerified()
	repo := &stubRepository{item: item}
	bus, _ := collectEvents(t, events.TypeItemUpdated)

	moderate := app.NewModerateItemHandler(
		repo, &stubPhotos{count: 2},
		&stubModerationProvider{result: domain.ModerationResult{Verdict: domain.ModerationApproved, Provider: "test"}},
		&stubModerationLog{}, stubPhotoBytesLoader{}, passthroughTx{}, fakeClock{}, events.NopPublisher{},
	)

	handler := app.NewUpdateItemHandler(repo, &stubPhotos{count: 2}, passthroughTx{}, fakeClock{}, bus, moderate)

	title := "Updated title triggers moderation"
	updated, err := handler.Handle(context.Background(), app.UpdateItemCommand{
		ItemID: item.ID(), ActorID: ownerID, Title: &title,
	})

	require.NoError(t, err)
	assert.Equal(t, domain.StatusPublished, updated.Status())
	assert.True(t, updated.AIVerified())
}
