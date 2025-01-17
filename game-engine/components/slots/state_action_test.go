package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewSymbolsStateAction(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name       string
		count      uint8
		altSymbols bool
		indexes    utils.Indexes
		want       []bool
	}{
		{
			name:    "limit 3, none flagged",
			count:   3,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			want:    []bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		},
		{
			name:    "limit 3, 1 flagged",
			count:   3,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 8, 1},
			want:    []bool{false, true, false, false, false, false, false, false, false, false, false, false, false, false, false},
		},
		{
			name:    "limit 3, 3 flagged",
			count:   3,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 6, 5, 4},
			want:    []bool{false, false, false, false, true, true, true, false, false, false, false, false, false, false, false},
		},
		{
			name:       "limit 5, none flagged",
			count:      5,
			altSymbols: true,
			indexes:    utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 4, 4, 4, 1, 2, 3},
			want:       []bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		},
		{
			name:       "limit 5, 1 flagged",
			count:      5,
			altSymbols: true,
			indexes:    utils.Indexes{1, 2, 3, 4, 2, 2, 1, 2, 3, 4, 5, 2, 2, 8, 2},
			want:       []bool{false, false, true, false, false, false, false, false, false, false, false, false, false, false, false},
		},
		{
			name:       "limit 5, 3 flagged",
			count:      5,
			altSymbols: true,
			indexes:    utils.Indexes{11, 2, 3, 11, 11, 2, 2, 3, 3, 3, 2, 11, 2, 3, 11},
			want:       []bool{false, false, true, true, false, false, false, false, false, false, false, true, false, false, false},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewFlagSymbolsAction(tc.count, tc.altSymbols)
			require.NotNil(t, a)

			assert.Equal(t, tc.count, a.flagCount)
			assert.Equal(t, tc.altSymbols, a.altSymbols)

			state := AcquireSymbolsState(setF1)
			require.NotNil(t, state)
			defer state.Release()

			spin.Debug(tc.indexes)

			a.StateUpdate(spin, state)
			assert.EqualValues(t, tc.want, state.flagged)
		})
	}
}

func TestStateAction_StateTriggered(t *testing.T) {
	t.Run("state triggered", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		a := NewFlagSymbolsAction(5, true, WithFreeSpins(0, 6, 15))
		require.NotNil(t, a)

		state := AcquireSymbolsState(setF1, 14)
		require.NotNil(t, state)

		a2 := a.StateTriggered(state)
		require.Nil(t, a2)

		state.flagged = []bool{false, true, true, true, true, true, true, true, false, true, true, true, true, true, true}
		a2 = a.StateTriggered(state)
		require.Nil(t, a2)

		state.flagged = []bool{false, true, true, true, true, true, true, true, true, true, true, true, true, true, false}
		a2 = a.StateTriggered(state)
		require.NotNil(t, a2)

		counts := make(map[uint8]int, 16)
		for ix := 0; ix < 1000; ix++ {
			n := a.NrOfSpins(prng)
			counts[n] = counts[n] + 1
		}

		assert.Equal(t, 10, len(counts))
		for k, v := range counts {
			assert.GreaterOrEqual(t, v, 50, k)
		}
	})
}
