package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestGridJumps_TestGridJump_5x3(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name    string
		indexes utils.Indexes
		sticky  []bool
		params  JumpParams
		offsets utils.UInt8s
		want    int
	}{
		{
			name:    "empties, any, none",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 6, 1, 2, 3},
			sticky:  []bool{false, false, false, false, false, false, true, false, false, false, false, false, false, false, false},
			params:  JumpParams{direction: GridAny, minJump: 1, maxJump: 1},
			offsets: utils.UInt8s{6},
		},
		{
			name:    "empties, any, one",
			indexes: utils.Indexes{1, 2, 3, 4, 0, 6, 10, 2, 3, 4, 5, 6, 1, 2, 3},
			sticky:  []bool{false, false, false, false, false, false, true, false, false, false, false, false, false, false, false},
			params:  JumpParams{direction: GridAny, minJump: 1, maxJump: 1},
			offsets: utils.UInt8s{6},
			want:    1,
		},
		{
			name:    "empties, any, three",
			indexes: utils.Indexes{10, 0, 3, 4, 5, 6, 10, 0, 3, 4, 5, 6, 1, 0, 10},
			sticky:  []bool{true, false, false, false, false, false, true, false, false, false, false, false, false, false, true},
			params:  JumpParams{direction: GridAny, minJump: 1, maxJump: 1},
			offsets: utils.UInt8s{0, 6, 14},
			want:    3,
		},
		{
			name:    "empties, left, none",
			indexes: utils.Indexes{0, 0, 0, 4, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0},
			sticky:  []bool{false, false, false, false, false, false, true, false, false, false, false, false, false, false, false},
			params:  JumpParams{direction: GridLeft, minJump: 1, maxJump: 1},
			offsets: utils.UInt8s{6},
		},
		{
			name:    "empties, right, none",
			indexes: utils.Indexes{0, 0, 0, 0, 0, 0, 10, 0, 0, 1, 0, 0, 0, 0, 0},
			sticky:  []bool{false, false, false, false, false, false, true, false, false, false, false, false, false, false, false},
			params:  JumpParams{direction: GridRight, minJump: 1, maxJump: 1},
			offsets: utils.UInt8s{6},
		},
		{
			name:    "empties, left, one",
			indexes: utils.Indexes{0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0},
			sticky:  []bool{false, false, false, false, false, false, true, false, false, false, false, false, false, false, false},
			params:  JumpParams{direction: GridLeft, minJump: 1, maxJump: 1},
			offsets: utils.UInt8s{6},
			want:    1,
		},
		{
			name:    "empties, right, one",
			indexes: utils.Indexes{0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0},
			sticky:  []bool{false, false, false, false, false, false, true, false, false, false, false, false, false, false, false},
			params:  JumpParams{direction: GridRight, minJump: 1, maxJump: 1},
			offsets: utils.UInt8s{6},
			want:    1,
		},
		{
			name:    "on, any, one",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 6, 1, 2, 3},
			sticky:  []bool{false, false, false, false, false, false, true, false, false, false, false, false, false, false, false},
			params:  JumpParams{direction: GridAny, onSymbols: true, refill: true, minJump: 1, maxJump: 1},
			offsets: utils.UInt8s{6},
			want:    1,
		},
		{
			name:    "on, any, three",
			indexes: utils.Indexes{1, 2, 11, 4, 5, 6, 10, 2, 3, 4, 5, 6, 1, 12, 3},
			sticky:  []bool{false, false, true, false, false, false, true, false, false, false, false, false, false, true, false},
			params:  JumpParams{direction: GridAny, onSymbols: true, refill: true, minJump: 1, maxJump: 1},
			offsets: utils.UInt8s{2, 6, 13},
			want:    3,
		},
		{
			name:    "on + off, any, one",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 6, 1, 2, 3},
			sticky:  []bool{false, false, false, false, false, false, true, false, false, false, false, false, false, false, false},
			params:  JumpParams{direction: GridAny, onSymbols: true, refill: true, minJump: 1, maxJump: 1, offGrid: 25},
			offsets: utils.UInt8s{6},
			want:    1,
		},
		{
			name:    "on + off, any, three",
			indexes: utils.Indexes{1, 2, 11, 4, 5, 6, 10, 2, 3, 4, 5, 6, 1, 12, 3},
			sticky:  []bool{false, false, true, false, false, false, true, false, false, false, false, false, false, true, false},
			params:  JumpParams{direction: GridAny, onSymbols: true, refill: true, minJump: 1, maxJump: 1, offGrid: 25},
			offsets: utils.UInt8s{2, 6, 13},
			want:    3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.params.offGrid > 0 {
				var off int

				for ix := 0; ix < 100; ix++ {
					copy(spin.indexes, tc.indexes)
					copy(spin.sticky, tc.sticky)

					var jumps GridJumps
					for _, offs := range tc.offsets {
						jumps = jumps.TestGridJump(&tc.params, offs, spin)
					}

					require.Equal(t, tc.want, len(jumps))

					for iy := range jumps {
						j := jumps[iy]
						if j.offGrid {
							off++
						}
					}
				}

				assert.NotZero(t, off)
			} else {
				copy(spin.indexes, tc.indexes)
				copy(spin.sticky, tc.sticky)

				var jumps GridJumps
				for _, offs := range tc.offsets {
					jumps = jumps.TestGridJump(&tc.params, offs, spin)
				}

				require.Equal(t, tc.want, len(jumps))

				for iy := range jumps {
					j := jumps[iy]
					assert.False(t, j.offGrid)
				}
			}
		})
	}
}
