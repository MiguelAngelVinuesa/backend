package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRoundResult(t *testing.T) {
	testCases := []struct {
		name    string
		spinSeq int
		result  *results.Result
	}{
		{
			name:    "r0",
			spinSeq: 1,
			result:  r0,
		},
		{
			name:    "r1",
			spinSeq: 2,
			result:  r1,
		},
		{
			name:    "r2",
			spinSeq: 3,
			result:  r2,
		},
		{
			name:    "r3",
			spinSeq: 4,
			result:  r3,
		},
		{
			name:    "r4",
			spinSeq: 5,
			result:  r4,
		},
		{
			name:    "r5",
			spinSeq: 6,
			result:  r5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := AcquireRoundResult(tc.spinSeq, tc.result)
			require.NotNil(t, r)
			defer r.Release()

			assert.Equal(t, tc.spinSeq, r.SpinSeq)
			assert.Zero(t, r.BalanceBefore)
			assert.Zero(t, r.BalanceAfter)
			assert.Zero(t, r.Bet)
			assert.Zero(t, r.TotalWin)
			assert.Zero(t, r.Win)
			assert.Zero(t, r.ProgressiveWin)
			assert.Equal(t, tc.result.Payouts, r.Payouts)
			assert.Equal(t, tc.result.Total, r.TotalPayout)
			assert.Equal(t, tc.result.AwardedFreeGames, r.AwardedFreeGames)
			assert.Equal(t, tc.result.FreeGames, r.FreeGames)
		})
	}
}
