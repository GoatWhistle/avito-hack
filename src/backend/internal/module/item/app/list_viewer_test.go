package app_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/shared/auth"
)

func TestListItemsPassesViewerToReadModel(t *testing.T) {
	t.Parallel()

	viewer := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
	read := &stubReadModel{}
	handler := app.NewListItemsHandler(read, nil)

	_, err := handler.Handle(t.Context(), app.ListItemsQuery{Viewer: viewer, Limit: 7})

	require.NoError(t, err)
	assert.Equal(t, viewer.ID, read.lastFilter.ViewerID)
}

func TestListItemsLeavesViewerEmptyForAnonymous(t *testing.T) {
	t.Parallel()

	read := &stubReadModel{}
	handler := app.NewListItemsHandler(read, nil)

	_, err := handler.Handle(t.Context(), app.ListItemsQuery{Limit: 7})

	require.NoError(t, err)
	assert.Equal(t, uuid.Nil, read.lastFilter.ViewerID)
}
