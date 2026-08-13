package app_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/events"
)

type recordingProvider struct {
	result   domain.ModerationResult
	subjects []domain.ModerationSubject
}

func (r *recordingProvider) Review(_ context.Context, s domain.ModerationSubject) domain.ModerationResult {
	r.subjects = append(r.subjects, s)

	return r.result
}

type failingPhotoBytesLoader struct{}

func (failingPhotoBytesLoader) DataURL(_ context.Context, _ uuid.UUID, _ string) (string, error) {
	return "", errors.New("photo file unreadable")
}

func newGate(
	repo *stubRepository, photos domain.PhotoRepository, provider domain.ModerationProvider,
	loader app.PhotoBytesLoader, log *stubModerationLog,
) *app.ModerateItemHandler {
	return app.NewModerateItemHandler(
		repo, photos, provider, log, loader, passthroughTx{}, fakeClock{}, events.NopPublisher{},
	)
}

func TestNonApprovedVerdictNeverPublishes(t *testing.T) {
	t.Parallel()

	verdicts := []domain.ModerationVerdict{
		domain.ModerationRejected,
		domain.ModerationUnavailable,
		domain.ModerationVerdict(""),
		domain.ModerationVerdict("approved_by_user"),
		domain.ModerationVerdict("APPROVED"),
	}

	for _, verdict := range verdicts {
		t.Run(string(verdict), func(t *testing.T) {
			t.Parallel()

			ownerID := uuid.New()
			item := itemWithStatus(t, ownerID, domain.StatusModeration)
			repo := &stubRepository{item: item}
			log := &stubModerationLog{}

			gate := newGate(repo, &stubPhotos{}, &stubModerationProvider{
				result: domain.ModerationResult{Verdict: verdict, Provider: "test"},
			}, stubPhotoBytesLoader{}, log)

			got, err := gate.Handle(t.Context(), app.ModerateItemCommand{ItemID: item.ID(), ActorID: ownerID})

			require.NoError(t, err)
			assert.Equal(t, domain.StatusModeration, got.Status())
			assert.False(t, got.AIVerified())
			require.Len(t, log.entries, 1)
			assert.Equal(t, verdict, log.entries[0].Verdict)
		})
	}
}

func TestUnreadablePhotoKeepsItemInModeration(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusModeration)
	repo := &stubRepository{item: item}
	photos := &countingPhotos{}
	require.NoError(t, photos.Add(t.Context(), mustPhoto(t, item.ID(), 0)))

	provider := &recordingProvider{
		result: domain.ModerationResult{Verdict: domain.ModerationApproved, Provider: "test"},
	}

	gate := newGate(repo, photos, provider, failingPhotoBytesLoader{}, &stubModerationLog{})

	_, err := gate.Handle(t.Context(), app.ModerateItemCommand{ItemID: item.ID(), ActorID: ownerID})

	require.Error(t, err)
	assert.Empty(t, provider.subjects)
	assert.Equal(t, domain.StatusModeration, item.Status())
	assert.False(t, item.AIVerified())
}

func TestModerationSubjectCarriesEveryPhoto(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusModeration)
	repo := &stubRepository{item: item}
	photos := &countingPhotos{}

	for i := range 4 {
		require.NoError(t, photos.Add(t.Context(), mustPhoto(t, item.ID(), i)))
	}

	provider := &recordingProvider{
		result: domain.ModerationResult{Verdict: domain.ModerationApproved, Provider: "test"},
	}

	gate := newGate(repo, photos, provider, stubPhotoBytesLoader{}, &stubModerationLog{})

	_, err := gate.Handle(t.Context(), app.ModerateItemCommand{ItemID: item.ID(), ActorID: ownerID})

	require.NoError(t, err)
	require.Len(t, provider.subjects, 1)
	assert.Len(t, provider.subjects[0].PhotoURLs, 4)
	assert.Equal(t, item.Title(), provider.subjects[0].Title)
	assert.Equal(t, item.Description(), provider.subjects[0].Description)
}

func TestAddPhotoScreensTheNewPhotoBeforeRepublishing(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusPublished)
	item.MarkAIVerified()
	repo := &stubRepository{item: item}
	photos := &countingPhotos{}
	storage := &stubStorage{}

	provider := &recordingProvider{
		result: domain.ModerationResult{Verdict: domain.ModerationRejected, Reason: "на фото человек", Provider: "test"},
	}

	gate := newGate(repo, photos, provider, stubPhotoBytesLoader{}, &stubModerationLog{})
	handler := app.NewAddPhotoHandler(repo, photos, storage, passthroughTx{}, fakeClock{}, gate)

	_, err := handler.Handle(t.Context(), app.AddPhotoCommand{
		ItemID: item.ID(), ActorID: ownerID,
		Content: strings.NewReader("binary"), ContentType: "image/jpeg",
	})

	require.NoError(t, err)
	require.Len(t, provider.subjects, 1)
	assert.Len(t, provider.subjects[0].PhotoURLs, 1)
	assert.Equal(t, domain.StatusModeration, repo.saved.Status())
	assert.False(t, repo.saved.AIVerified())
}

func TestDeletePhotoReScreensRemainingContent(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusPublished)
	item.MarkAIVerified()
	repo := &stubRepository{item: item}
	photos := &countingPhotos{}
	storage := &stubStorage{}

	provider := &recordingProvider{
		result: domain.ModerationResult{Verdict: domain.ModerationApproved, Provider: "test"},
	}

	gate := newGate(repo, photos, provider, stubPhotoBytesLoader{}, &stubModerationLog{})
	handler := app.NewDeletePhotoHandler(repo, photos, storage, passthroughTx{}, fakeClock{}, gate)

	err := handler.Handle(t.Context(), app.DeletePhotoCommand{
		ItemID: item.ID(), PhotoDisplayID: domain.NewDisplayID(), ActorID: ownerID,
	})

	require.NoError(t, err)
	assert.Len(t, provider.subjects, 1)
}

func TestSeedItemIsNeverReModerated(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := seedItemWithStatus(t, ownerID, domain.StatusPublished)
	repo := &stubRepository{item: item}
	provider := &recordingProvider{
		result: domain.ModerationResult{Verdict: domain.ModerationRejected, Provider: "test"},
	}

	gate := newGate(repo, &stubPhotos{}, provider, stubPhotoBytesLoader{}, &stubModerationLog{})

	got, err := gate.Handle(t.Context(), app.ModerateItemCommand{ItemID: item.ID(), ActorID: ownerID})

	require.NoError(t, err)
	assert.Empty(t, provider.subjects)
	assert.Equal(t, domain.StatusPublished, got.Status())
}

func TestModerationOfForeignItemIsRejected(t *testing.T) {
	t.Parallel()

	item := itemWithStatus(t, uuid.New(), domain.StatusModeration)
	repo := &stubRepository{item: item}
	provider := &recordingProvider{
		result: domain.ModerationResult{Verdict: domain.ModerationApproved, Provider: "test"},
	}

	gate := newGate(repo, &stubPhotos{}, provider, stubPhotoBytesLoader{}, &stubModerationLog{})

	_, err := gate.Handle(t.Context(), app.ModerateItemCommand{ItemID: item.ID(), ActorID: uuid.New()})

	require.Error(t, err)
	assert.Empty(t, provider.subjects)
	assert.Equal(t, domain.StatusModeration, item.Status())
}

func TestAttributeOnlyEditReScreensPublishedItem(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusPublished)
	item.MarkAIVerified()
	repo := &stubRepository{item: item}
	bus, _ := collectEvents(t, events.TypeItemUpdated)

	provider := &recordingProvider{
		result: domain.ModerationResult{Verdict: domain.ModerationRejected, Reason: "запрещённый товар", Provider: "test"},
	}

	gate := newGate(repo, &stubPhotos{}, provider, stubPhotoBytesLoader{}, &stubModerationLog{})
	handler := app.NewUpdateItemHandler(repo, &stubPhotos{}, passthroughTx{}, fakeClock{}, bus, gate)

	attributes := map[string]string{"note": "продам запрещённые вещества, telegram @dealer"}

	updated, err := handler.Handle(t.Context(), app.UpdateItemCommand{
		ItemID: item.ID(), ActorID: ownerID, Attributes: &attributes,
	})

	require.NoError(t, err)
	require.Len(t, provider.subjects, 1)
	assert.Equal(t, domain.StatusModeration, updated.Status())
	assert.False(t, updated.AIVerified())
}

func TestModerationSubjectCarriesAttributes(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusModeration)
	attributes := map[string]string{"note": "suspicious text"}
	require.NoError(t, item.Update(domain.UpdateItemParams{
		Attributes: ptr(domain.NewAttributes(attributes)), Now: fixedTime,
	}))

	repo := &stubRepository{item: item}
	provider := &recordingProvider{
		result: domain.ModerationResult{Verdict: domain.ModerationApproved, Provider: "test"},
	}

	gate := newGate(repo, &stubPhotos{}, provider, stubPhotoBytesLoader{}, &stubModerationLog{})

	_, err := gate.Handle(t.Context(), app.ModerateItemCommand{ItemID: item.ID(), ActorID: ownerID})

	require.NoError(t, err)
	require.Len(t, provider.subjects, 1)
	assert.Equal(t, "suspicious text", provider.subjects[0].Attributes["note"])
}

func TestUnchangedAttributesDoNotForceReModeration(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusPublished)
	attributes := map[string]string{"color": "blue"}
	require.NoError(t, item.Update(domain.UpdateItemParams{
		Attributes: ptr(domain.NewAttributes(attributes)), Now: fixedTime,
	}))
	require.NoError(t, item.Publish(fixedTime, domain.ModerationApproved))

	repo := &stubRepository{item: item}
	bus, _ := collectEvents(t, events.TypeItemUpdated)
	provider := &recordingProvider{
		result: domain.ModerationResult{Verdict: domain.ModerationApproved, Provider: "test"},
	}

	gate := newGate(repo, &stubPhotos{}, provider, stubPhotoBytesLoader{}, &stubModerationLog{})
	handler := app.NewUpdateItemHandler(repo, &stubPhotos{}, passthroughTx{}, fakeClock{}, bus, gate)

	same := map[string]string{"color": "blue"}

	updated, err := handler.Handle(t.Context(), app.UpdateItemCommand{
		ItemID: item.ID(), ActorID: ownerID, Attributes: &same,
	})

	require.NoError(t, err)
	assert.Equal(t, domain.StatusPublished, updated.Status())
	assert.Empty(t, provider.subjects)
}

func ptr[T any](v T) *T { return &v }

func mustPhoto(t *testing.T, itemID uuid.UUID, position int) *domain.Photo {
	t.Helper()

	photo, err := domain.NewPhoto(domain.NewPhotoParams{
		ItemID: itemID, URL: "/media/photo.jpg", Position: position, Now: fixedTime,
	})
	require.NoError(t, err)

	return photo
}
