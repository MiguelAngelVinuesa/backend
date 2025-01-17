package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewSymbolInject(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name    string
		symbol  utils.Index
		reels   utils.UInt8s
		indexes utils.Indexes
		sticky  []bool
		want    bool
	}{
		{
			name:    "no stickies, all reels, injected",
			symbol:  10,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			sticky:  []bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
			want:    true,
		},
		{
			name:    "few stickies, all reels, injected",
			symbol:  10,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			sticky:  []bool{false, false, false, false, false, true, false, false, false, false, false, true, false, false, false},
			want:    true,
		},
		{
			name:    "many stickies, all reels, injected",
			symbol:  10,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			sticky:  []bool{false, true, true, true, false, true, false, true, false, true, false, true, false, false, true},
			want:    true,
		},
		{
			name:    "all stickies, all reels, not injected",
			symbol:  10,
			reels:   utils.UInt8s{2, 3, 4},
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			sticky:  []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true},
		},
		{
			name:    "no stickies, reel 2+3+4, injected",
			symbol:  10,
			reels:   utils.UInt8s{2, 3, 4},
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			sticky:  []bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
			want:    true,
		},
		{
			name:    "few stickies, reel 2+3+4, injected",
			symbol:  10,
			reels:   utils.UInt8s{2, 3, 4},
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			sticky:  []bool{false, false, false, false, false, true, false, false, false, false, false, true, false, false, false},
			want:    true,
		},
		{
			name:    "many stickies, reel 2+3+4, injected",
			symbol:  10,
			reels:   utils.UInt8s{2, 3, 4},
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			sticky:  []bool{false, true, true, true, false, true, false, true, false, true, false, true, false, false, true},
			want:    true,
		},
		{
			name:    "all stickies, reel 2+3+4, not injected",
			symbol:  10,
			reels:   utils.UInt8s{2, 3, 4},
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			sticky:  []bool{false, true, true, true, true, true, true, true, true, true, true, true, false, false, true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewSymbolInject(tc.symbol, tc.reels...)
			require.NotNil(t, a)
			assert.True(t, a.single)
			assert.False(t, a.singleFromEdge)
			assert.Equal(t, tc.symbol, a.symbol)
			assert.EqualValues(t, tc.reels, a.singleReels)

			counts := make(map[uint8]int)
			for ix := 0; ix < 1000; ix++ {
				copy(spin.indexes, tc.indexes)
				copy(spin.sticky, tc.sticky)
				spin.ResetEffects()

				got := a.Triggered(spin)
				if tc.want {
					assert.NotNil(t, got)

					var count int

					for offs := range spin.injections {
						if symbol := spin.injections[offs]; symbol != utils.NullIndex {
							count++
							counts[uint8(offs)] = counts[uint8(offs)] + 1

							assert.Equal(t, tc.symbol, symbol)

							if len(tc.reels) > 0 {
								reel, _ := spin.gridDef.ReelRowFromOffset(offs)
								reel++
								assert.Contains(t, tc.reels, uint8(reel))
							}
						}
					}

					assert.Equal(t, 1, count)
				} else {
					assert.Nil(t, got)
				}
			}

			if len(tc.reels) > 0 {
				options := len(tc.reels) * 3
				for _, reel := range tc.reels {
					for row := 0; row < 3; row++ {
						if tc.sticky[int(reel-1)*3+row] {
							options--
						}
					}
				}
				assert.Equal(t, options, len(counts))
			} else {
				options := 15
				for ix := range tc.sticky {
					if tc.sticky[ix] {
						options--
					}
				}
				assert.Equal(t, options, len(counts))
			}

			for k, v := range counts {
				assert.NotZero(t, v)
				assert.False(t, tc.sticky[k])
			}
		})
	}
}

func TestNewSymbolInjectFromEdge(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(7, 7), WithSymbols(setF1), WithMask(4, 5, 6, 7, 6, 5, 4))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name    string
		symbol  utils.Index
		steps   uint8
		indexes utils.Indexes
		sticky  []bool
		want    bool
	}{
		{
			name:   "no stickies, outer ring, injected",
			symbol: 10,
			steps:  1,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 6, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false},
			want: true,
		},
		{
			name:   "few stickies, outer ring, injected",
			symbol: 10,
			steps:  1,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 6, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, true, false, false, false, false, false,
				false, false, false, false, false, false, false,
				true, false, false, true, false, false, false,
				false, false, false, true, false, false, false,
				false, false, false, false, true, false, false,
				false, false, false, false, false, false, false,
				true, false, false, false, false, false, false},
			want: true,
		},
		{
			name:   "many stickies, outer ring, injected",
			symbol: 10,
			steps:  1,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 6, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, true, true, true, false, false, false,
				true, false, true, false, true, false, false,
				true, false, false, true, true, false, false,
				false, false, false, true, false, false, true,
				false, false, false, false, true, true, false,
				true, false, false, false, false, false, false,
				false, true, false, true, false, false, false},
			want: true,
		},
		{
			name:   "all stickies, outer ring, not injected",
			symbol: 10,
			steps:  1,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 6, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				true, true, true, true, false, false, false,
				true, false, true, false, true, false, false,
				true, false, false, true, true, true, false,
				true, false, false, true, false, false, true,
				true, false, false, false, true, true, false,
				true, false, false, false, true, false, false,
				true, true, true, true, false, false, false},
		},
		{
			name:   "no stickies, 2nd ring, injected",
			symbol: 10,
			steps:  2,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 6, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false},
			want: true,
		},
		{
			name:   "few stickies, 2nd ring, injected",
			symbol: 10,
			steps:  2,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 6, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, false, false, false, false, false, false,
				false, true, false, false, true, false, false,
				false, false, false, false, true, false, false,
				true, false, false, false, false, false, false,
				false, true, false, false, false, true, false,
				false, false, true, false, false, false, false,
				false, false, false, true, false, false, false},
			want: true,
		},
		{
			name:   "many stickies, 2nd ring, injected",
			symbol: 10,
			steps:  2,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 6, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, false, false, false, false, false, false,
				false, true, false, true, true, false, false,
				false, true, false, false, true, false, false,
				true, false, false, false, false, false, false,
				false, true, false, false, false, true, false,
				false, false, true, true, false, false, false,
				false, false, false, true, false, false, false},
			want: true,
		},
		{
			name:   "all stickies, 2nd ring, not injected",
			symbol: 10,
			steps:  2,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 6, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, false, false, false, false, false, false,
				false, true, true, true, true, false, false,
				false, true, false, false, true, false, false,
				true, true, false, false, false, true, false,
				false, true, false, false, true, true, false,
				false, true, true, true, false, false, false,
				false, false, false, true, false, false, false},
		},
		{
			name:   "no stickies, 3rd ring, injected",
			symbol: 10,
			steps:  3,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 6, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false},
			want: true,
		},
		{
			name:   "few stickies, 3rd ring, injected",
			symbol: 10,
			steps:  3,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 6, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, true, false, false, false,
				false, false, false, false, true, false, false,
				false, false, true, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false},
			want: true,
		},
		{
			name:   "all stickies, 3rd ring, not injected",
			symbol: 10,
			steps:  3,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 6, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, true, false, false, false, false, false,
				true, false, false, false, false, false, false,
				false, false, true, true, false, false, false,
				false, false, true, false, true, false, false,
				false, true, true, true, false, false, false,
				false, false, false, false, false, false, false,
				false, false, true, false, false, false, false},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewSymbolInjectFromEdge(tc.symbol, tc.steps)
			require.NotNil(t, a)
			assert.True(t, a.single)
			assert.True(t, a.singleFromEdge)
			assert.Equal(t, tc.symbol, a.symbol)
			assert.Equal(t, tc.steps, a.singleEdgeSteps)

			counts := make(map[uint8]int)
			for ix := 0; ix < 1000; ix++ {
				copy(spin.indexes, tc.indexes)
				copy(spin.sticky, tc.sticky)
				spin.ResetEffects()

				got := a.Triggered(spin)
				if tc.want {
					assert.NotNil(t, got)

					var count int

					for offs := range spin.injections {
						if symbol := spin.injections[offs]; symbol != utils.NullIndex {
							count++
							counts[uint8(offs)] = counts[uint8(offs)] + 1

							assert.Equal(t, tc.symbol, symbol)
							assert.Equal(t, tc.steps, spin.gridDef.stepsOffGrid[offs])
						}
					}

					assert.Equal(t, 1, count)
				} else {
					assert.Nil(t, got)
				}
			}

			options := spin.gridDef.TilesFromEdge(tc.steps)
			for offs := range tc.sticky {
				if tc.sticky[offs] && spin.gridDef.stepsOffGrid[offs] == tc.steps {
					options--
				}
			}
			assert.Equal(t, options, uint8(len(counts)))
		})
	}
}

func TestNewClusterOffsetInject(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(7, 7), WithSymbols(setF1))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name    string
		min     uint8
		max     uint8
		offset  uint8
		indexes utils.Indexes
		sticky  []bool
		want    utils.Index
	}{
		{
			name:   "no stickies, center, cluster 5",
			min:    5,
			max:    5,
			offset: 24,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 6, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false},
			want: 6,
		},
		{
			name:   "few stickies, center, cluster 5",
			min:    5,
			max:    5,
			offset: 24,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 6, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, true, false, true, false, false,
				false, false, true, false, false, false, false,
				false, false, false, false, true, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false},
			want: 6,
		},
		{
			name:   "many stickies, center, cluster 5",
			min:    5,
			max:    5,
			offset: 24,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 6, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, false, true, false, false, false, false,
				false, false, true, false, false, false, false,
				false, true, true, false, true, false, false,
				true, false, true, false, true, false, false,
				false, false, false, true, true, false, false,
				false, true, true, true, false, false, false,
				false, false, false, false, false, false, false},
			want: 6,
		},
		{
			name:   "no stickies, center, cluster 5-6",
			min:    5,
			max:    6,
			offset: 24,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 6, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false},
			want: 6,
		},
		{
			name:   "many stickies, center, cluster 5-6",
			min:    5,
			max:    6,
			offset: 24,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 6, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, false, true, false, false, false, false,
				true, false, false, true, false, false, false,
				false, false, true, true, true, false, false,
				false, false, true, false, false, true, false,
				false, true, false, true, false, true, false,
				false, false, true, false, false, false, false,
				false, true, false, false, false, false, false},
			want: 6,
		},
		{
			name:   "no stickies, center, cluster 5-10",
			min:    5,
			max:    10,
			offset: 24,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 6, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false},
			want: 6,
		},
		{
			name:   "few stickies, center, cluster 10-20",
			min:    10,
			max:    20,
			offset: 24,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 6, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, true, false, true, false, false,
				false, false, false, false, false, true, false,
				false, false, false, true, false, false, false,
				false, true, false, false, false, false, false,
				false, false, false, false, false, false, false},
			want: 6,
		},
		{
			name:   "boxed in, center, cluster 10-15",
			min:    10,
			max:    15,
			offset: 24,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 6, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, true, true, false, false, false,
				false, false, true, false, true, false, false,
				false, false, true, true, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false},
			want: 6,
		},
		{
			name:   "all stickies, center, no cluster",
			min:    5,
			max:    5,
			offset: 24,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 6, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				true, true, true, true, false, false, false,
				true, true, true, true, true, false, false,
				true, true, true, true, true, true, false,
				true, true, true, true, true, true, true,
				true, true, true, true, true, true, false,
				true, true, true, true, true, false, false,
				true, true, true, true, false, false, false},
		},
		{
			name:   "no stickies, center wild, cluster 5-10",
			min:    5,
			max:    10,
			offset: 24,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 12, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false},
			want: 12,
		},
		{
			name:   "few stickies, center wild, cluster 5-10",
			min:    5,
			max:    10,
			offset: 24,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 12, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, true, false, false, false,
				false, false, true, false, false, false, false,
				false, false, true, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false},
			want: 12,
		},
		{
			name:   "boxed in, center wild, cluster 5-10",
			min:    5,
			max:    10,
			offset: 24,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 12, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, true, true, false, false, false,
				false, false, true, false, true, false, false,
				false, false, true, true, false, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false},
			want: 12,
		},
		{
			name:   "all stickies, center wild, no cluster",
			min:    5,
			max:    10,
			offset: 24,
			indexes: utils.Indexes{
				1, 2, 3, 4, 0, 0, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 5, 6, 0,
				3, 4, 5, 12, 7, 8, 9,
				1, 2, 3, 4, 5, 6, 0,
				5, 6, 7, 8, 9, 0, 0,
				1, 2, 3, 4, 0, 0, 0},
			sticky: []bool{
				true, true, true, true, false, false, false,
				true, true, true, true, true, false, false,
				true, true, true, true, true, true, false,
				true, true, true, true, true, true, true,
				true, true, true, true, true, true, false,
				true, true, true, true, true, false, false,
				true, true, true, true, false, false, false},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewClusterOffsetInject(tc.min, tc.max, tc.offset)
			require.NotNil(t, a)

			for ix := 0; ix < 25; ix++ {
				copy(spin.indexes, tc.indexes)
				copy(spin.sticky, tc.sticky)
				spin.ResetEffects()

				got := a.Triggered(spin)
				if tc.want > 0 {
					require.NotNil(t, got)

					var count uint8
					for offset := range spin.injections {
						if spin.injections[offset] != utils.NullIndex {
							count++
							assert.Equal(t, tc.want, spin.injections[offset])

							var n int
							list := spin.gridDef.withoutSelf[offset]
							for k := range list {
								if spin.indexes[k] == tc.want {
									n++
								}
							}
							assert.NotZero(t, n)
						}
					}

					assert.GreaterOrEqual(t, count, tc.min)
					assert.LessOrEqual(t, count, tc.max)
				} else {
					assert.Nil(t, got)
				}
			}
		})
	}
}
