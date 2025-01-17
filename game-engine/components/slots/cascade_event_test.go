package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestAcquireCascadeEvent(t *testing.T) {
	testCases := []struct {
		name     string
		spin     *SpinResult
		reels    int
		rows     int
		wantFrom utils.UInt8s
		wantTo   utils.UInt8s
		wantJ    string
	}{
		{
			name: "3x9 single",
			spin: &SpinResult{
				initial:    utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9},
				afterClear: utils.Indexes{1, 2, 0, 4, 5, 6, 7, 8, 9, 1, 2, 3, 0, 5, 6, 7, 8, 9, 1, 0, 3, 4, 5, 6, 7, 8, 9},
			},
			reels:    3,
			rows:     9,
			wantFrom: utils.UInt8s{1, 0, 11, 10, 9, 18},
			wantTo:   utils.UInt8s{2, 1, 12, 11, 10, 19},
			wantJ:    `{"kind":11,"from":[1,0,11,10,9,18],"to":[2,1,12,11,10,19]}`,
		},
		{
			name: "3x9 single repeated",
			spin: &SpinResult{
				initial:    utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9},
				afterClear: utils.Indexes{1, 2, 0, 4, 5, 0, 7, 8, 9, 1, 0, 3, 0, 5, 0, 7, 0, 9, 1, 0, 3, 4, 5, 6, 0, 8, 9},
			},
			reels:    3,
			rows:     9,
			wantFrom: utils.UInt8s{4, 3, 1, 0, 15, 13, 11, 9, 23, 22, 21, 20, 18},
			wantTo:   utils.UInt8s{5, 4, 3, 2, 16, 15, 14, 13, 24, 23, 22, 21, 20},
			wantJ:    `{"kind":11,"from":[4,3,1,0,15,13,11,9,23,22,21,20,18],"to":[5,4,3,2,16,15,14,13,24,23,22,21,20]}`,
		},
		{
			name: "3x9 multi",
			spin: &SpinResult{
				initial:    utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9},
				afterClear: utils.Indexes{1, 2, 0, 0, 0, 6, 7, 8, 9, 1, 2, 3, 0, 0, 0, 0, 8, 9, 1, 0, 0, 4, 5, 6, 7, 8, 9},
			},
			reels:    3,
			rows:     9,
			wantFrom: utils.UInt8s{1, 0, 11, 10, 9, 18},
			wantTo:   utils.UInt8s{4, 3, 15, 14, 13, 20},
			wantJ:    `{"kind":11,"from":[1,0,11,10,9,18],"to":[4,3,15,14,13,20]}`,
		},
		{
			name: "3x9 repeat",
			spin: &SpinResult{
				initial:    utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9},
				afterClear: utils.Indexes{1, 2, 0, 0, 0, 6, 0, 0, 9, 1, 0, 3, 0, 0, 0, 0, 8, 9, 1, 0, 0, 4, 0, 6, 0, 0, 9},
			},
			reels:    3,
			rows:     9,
			wantFrom: utils.UInt8s{5, 1, 0, 11, 9, 23, 21, 18},
			wantTo:   utils.UInt8s{7, 6, 5, 15, 14, 25, 24, 23},
			wantJ:    `{"kind":11,"from":[5,1,0,11,9,23,21,18],"to":[7,6,5,15,14,25,24,23]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := AcquireCascadeEvent(tc.spin, tc.reels, tc.rows).(*CascadeEvent)
			require.NotNil(t, e)
			defer e.Release()

			assert.EqualValues(t, tc.wantFrom, e.from)
			assert.EqualValues(t, tc.wantTo, e.to)

			enc := zjson.AcquireEncoder(1024)
			defer enc.Release()

			e.Encode(enc)
			got := enc.Bytes()
			assert.Equal(t, tc.wantJ, string(got))
		})
	}
}
