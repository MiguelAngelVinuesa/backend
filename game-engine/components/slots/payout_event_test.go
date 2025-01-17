package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestAcquirePayoutEvent(t *testing.T) {
	testCases := []struct {
		name   string
		payout *SpinPayout
		want   string
	}{
		{
			name:   "payline",
			payout: WinlinePayoutFromData(1.5, 1.0, 7, 3, PayLTR, 1, utils.UInt8s{0, 1, 2, 1, 0}).(*SpinPayout),
			want:   `{"kind":6,"payoutKind":1,"count":3,"symbol":7,"factor":1.5,"paylineRows":[0,1,2]}`,
		},
		{
			name:   "all payline",
			payout: AllPaylinePayout(1.2, 1.7, 6, 3, utils.UInt8s{0, 2, 1, 0, 0}).(*SpinPayout),
			want:   `{"kind":6,"payoutKind":1,"count":3,"symbol":6,"factor":2.04,"paylineRows":[0,2,1]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := AcquirePayoutEvent(tc.payout).(*PayoutEvent)
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
