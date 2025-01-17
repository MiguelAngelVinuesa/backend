package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpinAction(t *testing.T) {
	testCases := []struct {
		name   string
		stage  SpinActionStage
		result SpinActionResult
		spins  uint8
	}{
		{"payout / 1", ExtraPayouts, Payout, 1},
		{"wild / 2", AwardBonuses, FreeSpins, 2},
		{"bonus / 5", AwardBonuses, BonusGame, 5},
		{"clear payouts", TestClearance, Refill, 0},
		{"sticky symbols", TestStickiness, Processed, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prng := rng.NewRNG()
			defer prng.ReturnToPool()

			a := &SpinAction{stage: tc.stage, result: tc.result, nrOfSpins: tc.spins}
			assert.Equal(t, tc.stage, a.Stage())
			assert.Equal(t, tc.result, a.Result())
			assert.Equal(t, tc.spins, a.NrOfSpins(prng))
			assert.False(t, a.AltSymbols())
			assert.Nil(t, a.Triggered(nil))
			assert.Equal(t, false, a.BonusSymbol())

			a.Payout(nil, nil) // just for code coverage as the function is a no-op
		})
	}
}

func TestPurgeSpinActions(t *testing.T) {
	a1 := &SpinAction{}
	a2 := &SpinAction{}
	a3 := &SpinAction{}
	a4 := &SpinAction{}
	a5 := &SpinAction{}
	a6 := &SpinAction{}
	a7 := &SpinAction{}

	testCases := []struct {
		name    string
		in      SpinActions
		cap     int
		wantCap int
	}{
		{"nil", nil, 5, 5},
		{"empty", SpinActions{}, 5, 5},
		{"short", SpinActions{a1, a2, a3}, 5, 5},
		{"exact", SpinActions{a1, a2, a3, a4, a5}, 5, 5},
		{"long", SpinActions{a1, a2, a3, a4, a5, a6, a7}, 5, 7},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := PurgeSpinActions(tc.in, tc.cap)
			require.NotNil(t, out)
			assert.Equal(t, tc.wantCap, cap(out))
			assert.Zero(t, len(out))
		})
	}
}
