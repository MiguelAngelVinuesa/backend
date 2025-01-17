package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewStickyChoice(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	testCases := []struct {
		name    string
		indexes utils.Indexes
		symbol  utils.Index
		want    utils.UInt8s
	}{
		{
			name:    "two",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6},
			symbol:  3,
			want:    utils.UInt8s{2, 11},
		},
		{
			name:    "three",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 4, 9, 1, 2, 3, 4, 5, 6},
			symbol:  4,
			want:    utils.UInt8s{3, 7, 12},
		},
		{
			name:    "four",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 1, 1, 8, 9, 1, 2, 3, 4, 5, 6},
			symbol:  1,
			want:    utils.UInt8s{0, 5, 6, 9},
		},
		{
			name:    "five",
			indexes: utils.Indexes{9, 2, 3, 4, 9, 6, 7, 9, 9, 1, 2, 3, 4, 5, 9},
			symbol:  9,
			want:    utils.UInt8s{0, 4, 7, 8, 14},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spin := AcquireSpin(slots, prng)
			defer spin.Release()

			spin.indexes = tc.indexes

			c := AcquireStickyChoice(tc.symbol, spin)
			require.NotNil(t, c)
			defer c.Release()

			assert.Equal(t, tc.symbol, c.Symbol)
			assert.EqualValues(t, tc.want, c.Positions)
		})
	}
}
