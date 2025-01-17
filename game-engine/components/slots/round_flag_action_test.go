package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/rng"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewRoundFlagShapeDetect(t *testing.T) {
	testCases := []struct {
		name    string
		flag    int
		symbol  utils.Index
		grid    GridOffsets
		indexes utils.Indexes
		want    bool
	}{
		{
			name:    "no shape",
			flag:    1,
			grid:    GridOffsets{{0, 0}, {0, 1}, {5, 0}, {5, 1}},
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 5, 0, 1, 2, 3, 4, 4, 3, 2, 1, 5, 4, 3, 0, 2, 1, 0, 0},
		},
		{
			name:    "almost shape",
			flag:    1,
			grid:    GridOffsets{{0, 0}, {0, 1}, {5, 0}, {5, 1}},
			indexes: utils.Indexes{1, 1, 0, 0, 3, 4, 5, 0, 1, 2, 3, 4, 4, 3, 2, 1, 5, 4, 3, 0, 1, 2, 0, 0},
		},
		{
			name:    "have shape",
			flag:    1,
			grid:    GridOffsets{{0, 0}, {0, 1}, {5, 0}, {5, 1}},
			indexes: utils.Indexes{1, 1, 0, 0, 3, 4, 5, 0, 1, 2, 3, 4, 4, 3, 2, 1, 5, 4, 3, 0, 1, 1, 0, 0},
			want:    true,
		},
		{
			name:    "no shape - symbol 5",
			flag:    1,
			symbol:  5,
			grid:    GridOffsets{{0, 0}, {0, 1}, {5, 0}, {5, 1}},
			indexes: utils.Indexes{1, 5, 0, 0, 3, 4, 5, 0, 1, 2, 3, 4, 4, 3, 2, 1, 5, 4, 3, 0, 2, 1, 0, 0},
		},
		{
			name:    "almost shape - symbol 5",
			flag:    1,
			symbol:  5,
			grid:    GridOffsets{{0, 0}, {0, 1}, {5, 0}, {5, 1}},
			indexes: utils.Indexes{5, 5, 0, 0, 3, 4, 5, 0, 1, 2, 3, 4, 4, 3, 2, 1, 5, 4, 3, 0, 5, 2, 0, 0},
		},
		{
			name:    "have shape - symbol 5",
			flag:    1,
			symbol:  5,
			grid:    GridOffsets{{0, 0}, {0, 1}, {5, 0}, {5, 1}},
			indexes: utils.Indexes{5, 5, 0, 0, 3, 4, 5, 0, 1, 2, 3, 4, 4, 3, 2, 1, 5, 4, 3, 0, 5, 5, 0, 0},
			want:    true,
		},
	}

	prng := rng.AcquireRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(6, 4), WithSymbols(setF1))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spin.ResetSpin()
			spin.indexes = tc.indexes

			a := NewRoundFlagShapeDetect(tc.flag, tc.symbol, tc.grid)
			require.NotNil(t, a)
			assert.True(t, a.shapeDetect)
			assert.Equal(t, tc.flag, a.flag)
			assert.Equal(t, tc.symbol, a.symbol)
			assert.Equal(t, tc.grid, a.shapeGrid)

			got := a.Triggered(spin)
			if tc.want {
				assert.NotNil(t, got)
				assert.Equal(t, 1, spin.roundFlags[tc.flag])
			} else {
				assert.Nil(t, got)
				assert.Zero(t, spin.roundFlags[tc.flag])
			}
		})
	}
}
