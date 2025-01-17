package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestAcquireBombEvent(t *testing.T) {
	testCases := []struct {
		name        string
		spin        *SpinResult
		reels       int
		rows        int
		center      Offsets
		grid        GridOffsets
		wantCenter  uint8
		wantGrid    utils.UInt8s
		wantChanged utils.UInt8s
		wantJ       string
	}{
		{
			name: "top-left",
			spin: &SpinResult{
				initial:     utils.Indexes{13, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
				afterExpand: utils.Indexes{4, 4, 3, 4, 4, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			},
			reels:       5,
			rows:        3,
			center:      Offsets{0, 0},
			grid:        GridOffsets{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 0}, {0, 1}, {1, -1}, {1, 0}, {1, 1}},
			wantCenter:  0,
			wantGrid:    utils.UInt8s{0, 1, 3, 4},
			wantChanged: utils.UInt8s{0, 1, 4},
			wantJ:       `{"kind":2,"center":0,"grid":[0,1,3,4],"changed":[0,1,4]}`,
		},
		{
			name: "center-top & bottom-right",
			spin: &SpinResult{
				initial:     utils.Indexes{1, 2, 3, 4, 5, 6, 13, 2, 3, 4, 5, 6, 1, 2, 13},
				afterExpand: utils.Indexes{1, 2, 3, 6, 6, 6, 6, 6, 3, 6, 6, 6, 1, 6, 6},
			},
			reels:       5,
			rows:        3,
			center:      Offsets{2, 0},
			grid:        GridOffsets{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 0}, {0, 1}, {1, -1}, {1, 0}, {1, 1}},
			wantCenter:  6,
			wantGrid:    utils.UInt8s{3, 4, 6, 7, 9, 10},
			wantChanged: utils.UInt8s{3, 4, 6, 7, 9, 10},
			wantJ:       `{"kind":2,"center":6,"grid":[3,4,6,7,9,10],"changed":[3,4,6,7,9,10]}`,
		},
		{
			name: "middle",
			spin: &SpinResult{
				initial:     utils.Indexes{1, 2, 3, 4, 13, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
				afterExpand: utils.Indexes{2, 2, 2, 2, 2, 2, 2, 2, 2, 4, 5, 6, 1, 2, 3},
			},
			reels:       5,
			rows:        3,
			center:      Offsets{1, 1},
			grid:        GridOffsets{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 0}, {0, 1}, {1, -1}, {1, 0}, {1, 1}},
			wantCenter:  4,
			wantGrid:    utils.UInt8s{0, 1, 2, 3, 4, 5, 6, 7, 8},
			wantChanged: utils.UInt8s{0, 2, 3, 4, 5, 6, 8},
			wantJ:       `{"kind":2,"center":4,"grid":[0,1,2,3,4,5,6,7,8],"changed":[0,2,3,4,5,6,8]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := AcquireBombEvent(tc.spin, tc.reels, tc.rows, tc.center, tc.grid).(*BombEvent)
			require.NotNil(t, e)
			defer e.Release()

			assert.Equal(t, results.BombEvent, e.Kind())
			assert.Equal(t, tc.wantCenter, e.center)
			assert.EqualValues(t, tc.wantGrid, e.grid)
			assert.EqualValues(t, tc.wantChanged, e.changed)

			enc := zjson.AcquireEncoder(1024)
			defer enc.Release()

			e.Encode(enc)
			got := enc.Bytes()
			assert.Equal(t, tc.wantJ, string(got))
		})
	}
}
