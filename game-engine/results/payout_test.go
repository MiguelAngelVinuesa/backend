package results

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPurgePayouts(t *testing.T) {
	testCases := []struct {
		name    string
		in      Payouts
		cap     int
		wantCap int
	}{
		{"nil", nil, 5, 5},
		{"empty", Payouts{}, 5, 5},
		{"short", Payouts{p1.Clone().(Payout), p2.Clone().(Payout), p3.Clone().(Payout)}, 5, 5},
		{"exact", Payouts{p1.Clone().(Payout), p2.Clone().(Payout), p3.Clone().(Payout), p4.Clone().(Payout), p5.Clone().(Payout)}, 5, 5},
		{"long", Payouts{p1.Clone().(Payout), p2.Clone().(Payout), p3.Clone().(Payout), p4.Clone().(Payout), p5.Clone().(Payout), p6.Clone().(Payout), p7.Clone().(Payout)}, 5, 7},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := PurgePayouts(tc.in, tc.cap)
			require.NotNil(t, out)
			assert.Equal(t, tc.wantCap, cap(out))
			assert.Zero(t, len(out))
		})
	}
}
