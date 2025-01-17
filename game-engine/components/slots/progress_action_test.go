package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewResetProgressAction(t *testing.T) {
	t.Run("reset progress", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

		spin := AcquireSpin(slots, prng)
		defer spin.Release()

		a := NewResetProgressAction(1)
		require.NotNil(t, a)
		assert.True(t, a.reset)
		assert.False(t, a.payoutSymbols)
		assert.Equal(t, 1, a.startLevel)

		spin.progressLevel = 17
		got := a.Triggered(spin)
		assert.NotNil(t, got)
		assert.Equal(t, 1, spin.progressLevel)
	})
}

func TestNewPayoutSymbolProgress(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	testCases := []struct {
		name    string
		start   int
		max     int
		indexes utils.Indexes
		payouts utils.UInt8s
		trigger bool
		want    int
	}{
		{
			name:    "no payouts - 0",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			payouts: utils.UInt8s{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:    "no payouts - 13",
			start:   13,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			payouts: utils.UInt8s{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:    "3 payouts - 0",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			payouts: utils.UInt8s{0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0},
			trigger: true,
			want:    3,
		},
		{
			name:    "3 payouts - 13",
			start:   13,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			payouts: utils.UInt8s{0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			trigger: true,
			want:    16,
		},
		{
			name:    "7 payouts - 0",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			payouts: utils.UInt8s{1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1},
			trigger: true,
			want:    7,
		},
		{
			name:    "7 payouts - 13 - max 20",
			start:   13,
			max:     20,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			payouts: utils.UInt8s{0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 1, 1, 0, 1, 1},
			trigger: true,
			want:    20,
		},
		{
			name:    "7 payouts - start 14 - max 20",
			start:   14,
			max:     20,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			payouts: utils.UInt8s{0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 1, 1, 0, 1, 1},
			trigger: true,
			want:    20,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewPayoutSymbolProgress(tc.start, tc.max)
			require.NotNil(t, a)
			assert.True(t, a.payoutSymbols)
			assert.False(t, a.reset)

			spin := AcquireSpin(slots, prng)
			require.NotNil(t, spin)
			defer spin.Release()

			spin.indexes = tc.indexes
			spin.payouts = tc.payouts

			a2 := a.Triggered(spin)
			if tc.trigger {
				assert.NotNil(t, a2)
				assert.Equal(t, tc.want, spin.progressLevel)
			} else {
				assert.Nil(t, a2)
				assert.Zero(t, spin.progressLevel)
			}
		})
	}
}
