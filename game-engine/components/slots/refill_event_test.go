package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestAcquireRefillEvent(t *testing.T) {
	testCases := []struct {
		name string
		spin *SpinResult
		want string
	}{
		{
			name: "none",
			spin: &SpinResult{
				initial: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			},
			want: `{"kind":12}`,
		},
		{
			name: "few",
			spin: &SpinResult{
				initial:    utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
				afterClear: utils.Indexes{1, 2, 3, 4, 0, 0, 1, 2, 0, 4, 5, 6, 1, 2, 3},
			},
			want: `{"kind":12,"refill":[4,5,8]}`,
		},
		{
			name: "many",
			spin: &SpinResult{
				initial:    utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
				afterClear: utils.Indexes{1, 0, 0, 4, 0, 0, 0, 0, 3, 0, 0, 0, 0, 2, 0},
			},
			want: `{"kind":12,"refill":[1,2,4,5,6,7,9,10,11,12,14]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := AcquireRefillEvent(tc.spin).(*RefillEvent)
			require.NotNil(t, e)
			defer e.Release()

			enc := zjson.AcquireEncoder(1024)
			defer enc.Release()

			e.Encode(enc)
			got := enc.Bytes()
			assert.Equal(t, tc.want, string(got))
		})
	}
}
