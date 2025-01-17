package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
)

func TestNewReelAnticipation(t *testing.T) {
	testCases := []struct {
		name      string
		startReel uint8
		stopReel  uint8
		j         string
	}{
		{
			name: "0-0",
			j:    `{"kind":1}`,
		},
		{
			name:      "1-0",
			startReel: 1,
			j:         `{"kind":1,"startReel":1}`,
		},
		{
			name:     "0-1",
			stopReel: 1,
			j:        `{"kind":1,"stopReel":1}`,
		},
		{
			name:      "1-2",
			startReel: 1,
			stopReel:  2,
			j:         `{"kind":1,"startReel":1,"stopReel":2}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := AcquireReelAnticipation(tc.startReel, tc.stopReel).(*ReelAnticipation)
			require.NotNil(t, e)
			defer e.Release()

			assert.Equal(t, results.ReelAnticipationEvent, e.Kind())
			assert.Equal(t, tc.startReel, e.startReel)
			assert.Equal(t, tc.stopReel, e.stopReel)

			enc := zjson.AcquireEncoder(128)
			defer enc.Release()

			enc.Object(e)
			j := enc.Bytes()
			assert.Equal(t, tc.j, string(j))
		})
	}
}
