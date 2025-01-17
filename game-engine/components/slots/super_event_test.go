package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestAcquireSuperEvent(t *testing.T) {
	testCases := []struct {
		name      string
		spin      *SpinResult
		wantShape utils.UInt8s
		wantJ     string
	}{
		{
			name: "first",
			spin: &SpinResult{
				initial: utils.Indexes{1, 2, 3, 4, 5, 4, 1, 4, 3, 4, 5, 4, 1, 2, 3},
				sticky:  utils.UInt8s{0, 0, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 0, 0},
			},
			wantShape: utils.UInt8s{3, 5, 7, 9, 11},
			wantJ:     `{"kind":3,"shape":[3,5,7,9,11],"first":1}`,
		},
		{
			name: "next",
			spin: &SpinResult{
				initial: utils.Indexes{1, 2, 3, 4, 5, 4, 1, 4, 3, 4, 5, 4, 1, 2, 4},
				sticky:  utils.UInt8s{0, 0, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 0, 1},
			},
			wantShape: utils.UInt8s{3, 5, 7, 9, 11},
			wantJ:     `{"kind":3,"shape":[3,5,7,9,11]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := AcquireSuperEvent(tc.spin).(*SuperEvent)
			require.NotNil(t, e)
			defer e.Release()

			assert.Equal(t, tc.wantShape, e.shape)

			enc := zjson.AcquireEncoder(1024)
			defer enc.Release()

			e.Encode(enc)
			got := enc.Bytes()
			assert.Equal(t, tc.wantJ, string(got))
		})
	}
}
