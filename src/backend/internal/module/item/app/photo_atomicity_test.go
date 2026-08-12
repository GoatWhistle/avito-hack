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
)

type lockRecordingRepository struct {
	item   *domain.Item
	events *[]string
}

func (r *lockRecordingRepository) Save(_ context.Context, _ *domain.Item) error { return nil }

func (r *lockRecordingRepository) ByID(_ context.Context, _ uuid.UUID) (*domain.Item, error) {
	*r.events = append(*r.events, "by-id")

	return r.item, nil
}

func (r *lockRecordingRepository) ByIDForUpdate(_ context.Context, _ uuid.UUID) (*domain.Item, error) {
	*r.events = append(*r.events, "lock")

	return r.item, nil
}

func (r *lockRecordingRepository) ByDisplayID(_ context.Context, _ string) (*domain.Item, error) {
	return r.item, nil
}

type recordingPhotos struct {
	countingPhotos
	events *[]string
}

func (p *recordingPhotos) CountByItemID(ctx context.Context, id uuid.UUID) (int, error) {
	*p.events = append(*p.events, "count")

	return p.countingPhotos.CountByItemID(ctx, id)
}

func (p *recordingPhotos) Add(ctx context.Context, photo *domain.Photo) error {
	*p.events = append(*p.events, "add")

	return p.countingPhotos.Add(ctx, photo)
}

type recordingTx struct{ events *[]string }

func (t recordingTx) WithTx(ctx context.Context, fn func(context.Context) error) error {
	*t.events = append(*t.events, "tx-begin")

	if err := fn(ctx); err != nil {
		*t.events = append(*t.events, "tx-rollback")

		return err
	}

	*t.events = append(*t.events, "tx-commit")

	return nil
}

func TestAddPhotoCountsAndInsertsUnderRowLockInOneTx(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusDraft)

	var events []string
	repo := &lockRecordingRepository{item: item, events: &events}
	photos := &recordingPhotos{countingPhotos: countingPhotos{stubPhotos: stubPhotos{count: 1}}, events: &events}

	handler := app.NewAddPhotoHandler(repo, photos, &stubStorage{}, recordingTx{events: &events}, fakeClock{}, nil)

	_, err := handler.Handle(t.Context(), app.AddPhotoCommand{
		ItemID: item.ID(), ActorID: ownerID,
		Content: strings.NewReader("binary"), ContentType: "image/jpeg",
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"tx-begin", "lock", "count", "add", "tx-commit"}, events)
}

func TestAddPhotoRemovesFileWhenCommitFails(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusDraft)
	commitErr := errors.New("commit failed")
	storage := &stubStorage{}

	handler := app.NewAddPhotoHandler(
		&stubRepository{item: item}, &countingPhotos{}, storage,
		failingCommitTx{err: commitErr}, fakeClock{}, nil)

	_, err := handler.Handle(t.Context(), app.AddPhotoCommand{
		ItemID: item.ID(), ActorID: ownerID,
		Content: strings.NewReader("binary"), ContentType: "image/jpeg",
	})

	require.ErrorIs(t, err, commitErr)
	assert.Equal(t, 1, storage.saved)
	assert.Equal(t, 1, storage.removed)
}

type failingCommitTx struct{ err error }

func (t failingCommitTx) WithTx(ctx context.Context, fn func(context.Context) error) error {
	if err := fn(ctx); err != nil {
		return err
	}

	return t.err
}
