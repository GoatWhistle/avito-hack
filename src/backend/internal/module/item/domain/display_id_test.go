package domain_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/avito-hack/backend/internal/module/item/domain"
)

var displayIDPattern = regexp.MustCompile(`^[a-z0-9]{12}$`)

func TestNewDisplayIDFormat(t *testing.T) {
	t.Parallel()

	id := domain.NewDisplayID()

	assert.Len(t, id, 12)
	assert.Regexp(t, displayIDPattern, id)
	assert.NotContains(t, id, "i")
	assert.NotContains(t, id, "l")
	assert.NotContains(t, id, "o")
	assert.NotContains(t, id, "u")
}

func TestNewDisplayIDIsUnpredictable(t *testing.T) {
	t.Parallel()

	seen := make(map[string]struct{}, 1000)

	for range 1000 {
		id := domain.NewDisplayID()

		_, exists := seen[id]
		assert.False(t, exists, "collision detected for %s", id)
		seen[id] = struct{}{}
	}
}
