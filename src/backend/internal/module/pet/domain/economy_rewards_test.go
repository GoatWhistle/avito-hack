package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRewardQualityListingVideoBonusExceedsBase(t *testing.T) {
	t.Parallel()

	base := QualityListing{
		HasPhoto: true, HasPrice: true,
		Description: strings.Repeat("а", qualityDescriptionLen+1),
		Now:         testTime(),
	}
	withVideo := base
	withVideo.HasVideo = true

	plain, err := petWithParams(70, 70, 100).RewardQualityListing(base)
	require.NoError(t, err)

	video, err := petWithParams(70, 70, 100).RewardQualityListing(withVideo)
	require.NoError(t, err)

	assert.Less(t, plain.XPGranted, video.XPGranted)
}
