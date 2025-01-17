package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestAcquireReelNudge(t *testing.T) {
	testCases := []struct {
		name     string
		teaser   bool
		reel     uint8
		size     uint8
		symbol   utils.Index
		location NudgeLocation
	}{
		{name: "teaser top", teaser: true, reel: 3, size: 1, symbol: 11, location: NudgeTop},
		{name: "teaser bottom", teaser: true, reel: 1, size: 1, symbol: 11, location: NudgeBottom},
		{name: "nudge 1 top", reel: 2, size: 1, symbol: 11, location: NudgeTop},
		{name: "nudge 2 top", reel: 4, size: 2, symbol: 11, location: NudgeTop},
		{name: "nudge 1 bottom", reel: 1, size: 1, symbol: 11, location: NudgeBottom},
		{name: "nudge 2 bottom", reel: 5, size: 2, symbol: 11, location: NudgeBottom},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := AcquireReelNudge(tc.teaser, tc.reel, tc.size, tc.symbol, tc.location)
			require.NotNil(t, a)
			defer a.Release()

			assert.Equal(t, tc.teaser, a.teaser)
			assert.Equal(t, tc.reel, a.reel)
			assert.Equal(t, tc.size, a.size)
			assert.Equal(t, tc.symbol, a.symbol)
			assert.Equal(t, tc.location, a.location)
		})
	}
}
