package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewSpinState(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	testCases := []struct {
		name      string
		indexes   utils.Indexes
		symbol    utils.Index
		sticky    []bool
		newSymbol utils.Index
		newSticky []bool
	}{
		{
			name:      "2, changed to 1",
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 7, 8, 9, 10, 11, 12},
			symbol:    2,
			sticky:    []bool{false, true, false, false, false, false, false, true, false, false, false, false, false, false, false},
			newSymbol: 1,
			newSticky: []bool{true, false, false, false, false, false, true, false, false, false, false, false, false, false, false},
		},
		{
			name:      "2, changed to 2",
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 7, 8, 9, 10, 11, 12},
			symbol:    2,
			sticky:    []bool{false, true, false, false, false, false, false, true, false, false, false, false, false, false, false},
			newSymbol: 2,
			newSticky: []bool{false, true, false, false, false, false, false, true, false, false, false, false, false, false, false},
		},
		{
			name:      "7, changed to 9",
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 7, 8, 9, 10, 11, 12},
			symbol:    7,
			sticky:    []bool{false, false, false, false, false, false, true, false, false, true, false, false, false, false, false},
			newSymbol: 9,
			newSticky: []bool{false, false, false, false, false, false, false, false, true, false, false, true, false, false, false},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spin := AcquireSpin(slots, prng)
			defer spin.Release()

			spin.indexes = tc.indexes
			spin.sticky = tc.sticky
			spin.stickySymbol = tc.symbol

			s := AcquireSpinState(spin)
			require.NotNil(t, s)
			defer s.Release()

			assert.EqualValues(t, tc.indexes, s.indexes)
			assert.EqualValues(t, tc.sticky, s.sticky)
			assert.Equal(t, tc.symbol, s.stickySymbol)

			c := s.Clone().(*SpinState)
			require.NotNil(t, c)
			assert.EqualValues(t, c, s)

			if tc.newSymbol > 0 {
				c.SetStickySymbol(tc.newSymbol)
				assert.Equal(t, tc.newSymbol, c.stickySymbol)
				assert.EqualValues(t, tc.newSticky, c.sticky)
			}
		})
	}
}
