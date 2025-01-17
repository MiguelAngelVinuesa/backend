package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestAcquireStickyEvent(t *testing.T) {
	testCases := []struct {
		name         string
		spin         *SpinResult
		wantStickies utils.UInt8s
		wantJ        string
	}{
		{
			name: "few",
			spin: &SpinResult{
				initial: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
				sticky:  utils.UInt8s{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
			},
			wantStickies: utils.UInt8s{3, 12},
			wantJ:        `{"kind":9,"stickies":[3,12]}`,
		},
		{
			name: "many",
			spin: &SpinResult{
				initial: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
				sticky:  utils.UInt8s{0, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 1, 0, 1},
			},
			wantStickies: utils.UInt8s{3, 4, 7, 9, 12, 14},
			wantJ:        `{"kind":9,"stickies":[3,4,7,9,12,14]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := AcquireStickyEvent(tc.spin).(*StickyEvent)
			require.NotNil(t, e)
			defer e.Release()

			assert.Equal(t, tc.wantStickies, e.stickies)

			enc := zjson.AcquireEncoder(1024)
			defer enc.Release()

			e.Encode(enc)
			got := enc.Bytes()
			assert.Equal(t, tc.wantJ, string(got))
		})
	}
}
