package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewBestSymbolStickyAction(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	testCases := []struct {
		name    string
		indexes utils.Indexes
		want    utils.Index
		sticky  []bool
	}{
		{
			name:    "2x 9",
			indexes: utils.Indexes{10, 11, 12, 9, 8, 7, 7, 8, 9, 6, 5, 4, 3, 6, 5},
			want:    9,
			sticky:  []bool{false, false, false, true, false, false, false, false, true, false, false, false, false, false, false},
		},
		{
			name:    "3x 7",
			indexes: utils.Indexes{10, 11, 12, 10, 11, 12, 7, 8, 9, 7, 8, 9, 7, 6, 5},
			want:    7,
			sticky:  []bool{false, false, false, false, false, false, true, false, false, true, false, false, true, false, false},
		},
		{
			name:    "4x 1",
			indexes: utils.Indexes{1, 10, 11, 12, 10, 11, 12, 1, 8, 9, 1, 8, 9, 6, 1},
			want:    1,
			sticky:  []bool{true, false, false, false, false, false, false, true, false, false, true, false, false, false, true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewBestSymbolStickyAction()
			require.NotNil(t, a)

			spin := AcquireSpin(slots, prng)
			defer spin.Release()
			spin.indexes = tc.indexes

			a2 := a.TriggeredWithState(spin, nil)
			require.NotNil(t, a2)
			assert.Equal(t, tc.want, spin.stickySymbol)
			assert.EqualValues(t, tc.sticky, spin.sticky)
		})
	}
}
