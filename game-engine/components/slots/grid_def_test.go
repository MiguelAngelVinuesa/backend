package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewGridDefinition(t *testing.T) {
	testCases := []struct {
		name       string
		reels      int
		rows       int
		mask       utils.UInt8s
		wantMask   utils.UInt8s
		testPairs  [][2]uint8
		testFail   []bool
		testWant   []GridDirection
		edge1count uint8
		edge2count uint8
		edge3count uint8
		edge4count uint8
	}{
		{
			name:      "3x3 no mask",
			reels:     3,
			rows:      3,
			wantMask:  utils.UInt8s{3, 3, 3},
			testPairs: [][2]uint8{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5}, {0, 6}, {1, 3}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}, {6, 0}, {3, 1}},
			testFail:  []bool{false, false, true, false, false, true, true, false, false, true, false, false, true, true, false},
			testWant:  []GridDirection{0, GridDown, 0, GridRight, GridRightDown, 0, 0, GridRightUp, GridUp, 0, GridLeft, GridLeftUp, 0, 0, GridLeftDown},
		},
		{
			name:     "5x3 no mask",
			reels:    5,
			rows:     3,
			wantMask: utils.UInt8s{3, 3, 3, 3, 3},
		},
		{
			name:      "6x4 no mask",
			reels:     6,
			rows:      4,
			wantMask:  utils.UInt8s{4, 4, 4, 4, 4, 4},
			testPairs: [][2]uint8{{5, 0}, {5, 1}, {5, 2}, {5, 3}, {5, 4}, {5, 5}, {5, 6}, {5, 7}, {5, 8}, {5, 9}, {5, 10}, {5, 11}, {5, 12}},
			testFail:  []bool{false, false, false, true, false, false, false, true, false, false, false, true, true},
			testWant:  []GridDirection{GridLeftUp, GridLeft, GridLeftDown, 0, GridUp, 0, GridDown, 0, GridRightUp, GridRight, GridRightDown, 0, 0},
		},
		{
			name:       "5x4 mask descending",
			reels:      5,
			rows:       4,
			mask:       utils.UInt8s{4, 3, 3, 3, 2},
			wantMask:   utils.UInt8s{4, 3, 3, 3, 2},
			testPairs:  [][2]uint8{{5, 0}, {5, 1}, {5, 2}, {5, 3}, {5, 4}, {5, 5}, {5, 6}, {5, 8}, {5, 9}, {5, 10}, {5, 12}},
			testFail:   []bool{true, false, false, true, false, false, false, false, false, false, true},
			testWant:   []GridDirection{0, GridLeftUp, GridLeftDown, 0, GridUp, 0, GridDown, GridRightUp, GridRight, GridRightDown, 0},
			edge1count: 12,
			edge2count: 3,
		},
		{
			name:       "5x4 mask ascending",
			reels:      5,
			rows:       4,
			mask:       utils.UInt8s{2, 3, 3, 3, 4},
			wantMask:   utils.UInt8s{2, 3, 3, 3, 4},
			testPairs:  [][2]uint8{{5, 0}, {5, 1}, {5, 4}, {5, 5}, {5, 6}, {5, 8}, {5, 9}, {5, 10}, {5, 12}},
			testFail:   []bool{false, false, false, false, false, false, false, false, true},
			testWant:   []GridDirection{GridLeftUp, GridLeftDown, GridUp, 0, GridDown, GridRightUp, GridRight, GridRightDown, 0},
			edge1count: 12,
			edge2count: 3,
		},
		{
			name:       "6x4 mask hexagonal",
			reels:      6,
			rows:       4,
			mask:       utils.UInt8s{2, 3, 4, 4, 3, 2},
			wantMask:   utils.UInt8s{2, 3, 4, 4, 3, 2},
			edge1count: 12,
			edge2count: 6,
		},
		{
			name:       "7x7 mask hexagonal",
			reels:      7,
			rows:       7,
			mask:       utils.UInt8s{4, 5, 6, 7, 6, 5, 4},
			wantMask:   utils.UInt8s{4, 5, 6, 7, 6, 5, 4},
			testPairs:  [][2]uint8{{24, 15}, {24, 16}, {24, 17}, {24, 18}, {24, 19}, {24, 22}, {24, 23}, {24, 24}, {24, 25}, {24, 26}, {24, 29}, {24, 30}, {24, 31}, {24, 32}, {24, 33}},
			testFail:   []bool{true, false, false, true, true, true, false, false, false, true, true, false, false, true, true},
			testWant:   []GridDirection{0, GridLeftUp, GridLeftDown, 0, 0, 0, GridUp, 0, GridDown, 0, 0, GridRightUp, GridRightDown, 0, 0},
			edge1count: 18,
			edge2count: 12,
			edge3count: 6,
			edge4count: 1,
		},
		{
			name:       "5x6 mask trapezoid",
			reels:      5,
			rows:       6,
			mask:       utils.UInt8s{2, 4, 6, 4, 2},
			wantMask:   utils.UInt8s{2, 4, 6, 4, 2},
			testPairs:  [][2]uint8{{14, 6}, {14, 7}, {14, 8}, {14, 9}, {14, 12}, {14, 13}, {14, 14}, {14, 15}, {14, 16}, {14, 17}, {14, 18}, {14, 19}, {14, 20}, {14, 21}},
			testFail:   []bool{false, false, false, true, true, false, false, false, true, true, false, false, false, true},
			testWant:   []GridDirection{GridLeftUp, GridLeft, GridLeftDown, 0, 0, GridUp, 0, GridDown, 0, 0, GridRightUp, GridRight, GridRightDown, 0},
			edge1count: 10,
			edge2count: 6,
			edge3count: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGridDefinition(tc.reels, tc.rows, tc.mask)
			require.NotNil(t, g)
			assert.Equal(t, uint8(tc.reels), g.ReelCount())
			assert.Equal(t, uint8(tc.rows), g.RowCount())
			assert.EqualValues(t, tc.wantMask, g.GridMask())

			if tc.edge1count > 0 {
				assert.Equal(t, tc.edge1count, g.TilesFromEdge(1))
			}
			if tc.edge2count > 0 {
				assert.Equal(t, tc.edge2count, g.TilesFromEdge(2))
			}
			if tc.edge3count > 0 {
				assert.Equal(t, tc.edge3count, g.TilesFromEdge(3))
			}
			if tc.edge4count > 0 {
				assert.Equal(t, tc.edge4count, g.TilesFromEdge(4))
			}

			if len(tc.mask) > 0 {
				assert.True(t, g.HaveMask())
			} else {
				assert.False(t, g.HaveMask())
			}

			for reel := 0; reel < tc.reels; reel++ {
				max := tc.rows
				if len(tc.mask) > 0 {
					max = int(tc.mask[reel])
				}
				for row := 0; row < max; row++ {
					assert.NotEmpty(t, g.Neighbors(uint8(reel*tc.rows+row)))
				}
			}

			for ix := range tc.testPairs {
				p := tc.testPairs[ix]
				ok, dir := g.IsNeighbor(p[0], p[1])
				assert.Equal(t, tc.testFail[ix], !ok)
				assert.Equal(t, tc.testWant[ix], dir)
			}
		})
	}
}

func TestNewGridDefinition_Fail(t *testing.T) {
	testCases := []struct {
		name  string
		reels int
		rows  int
		mask  utils.UInt8s
	}{
		{name: "empty reels", rows: 3},
		{name: "reels too large", reels: 13, rows: 3},
		{name: "empty rows", reels: 3},
		{name: "rows too large", reels: 5, rows: 13},
		{name: "grid too large", reels: 12, rows: 10},
		{name: "short mask", reels: 5, rows: 3, mask: utils.UInt8s{1, 2, 3}},
		{name: "long mask", reels: 5, rows: 3, mask: utils.UInt8s{1, 2, 3, 4, 5, 6}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				e := recover()
				require.NotNil(t, e)
			}()

			g := NewGridDefinition(tc.reels, tc.rows, tc.mask)
			require.Nil(t, g)
		})
	}
}

func TestGridDefinition_IsOnTheEdge(t *testing.T) {
	testCases := []struct {
		name   string
		reels  int
		rows   int
		mask   utils.UInt8s
		offset int
		want   bool
	}{
		{name: "5x3 - no mask - left edge", reels: 5, rows: 3, offset: 1, want: true},
		{name: "5x3 - no mask - right edge", reels: 5, rows: 3, offset: 13, want: true},
		{name: "5x3 - no mask - top edge", reels: 5, rows: 3, offset: 6, want: true},
		{name: "5x3 - no mask - bottom edge", reels: 5, rows: 3, offset: 8, want: true},
		{name: "5x3 - no mask - middle", reels: 5, rows: 3, offset: 7},
		{name: "6x4 - mask - left edge", reels: 6, rows: 4, mask: utils.UInt8s{2, 3, 4, 4, 3, 2}, offset: 1, want: true},
		{name: "6x4 - mask - right edge", reels: 6, rows: 4, mask: utils.UInt8s{2, 3, 4, 4, 3, 2}, offset: 21, want: true},
		{name: "6x4 - mask - top edge", reels: 6, rows: 4, mask: utils.UInt8s{2, 3, 4, 4, 3, 2}, offset: 8, want: true},
		{name: "6x4 - mask - bottom edge", reels: 6, rows: 4, mask: utils.UInt8s{2, 3, 4, 4, 3, 2}, offset: 6, want: true},
		{name: "6x4 - mask - middle (1)", reels: 6, rows: 4, mask: utils.UInt8s{2, 3, 4, 4, 3, 2}, offset: 9},
		{name: "6x4 - mask - middle (2)", reels: 6, rows: 4, mask: utils.UInt8s{2, 3, 4, 4, 3, 2}, offset: 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGridDefinition(tc.reels, tc.rows, tc.mask)
			require.NotNil(t, g)

			got := g.IsOnTheEdge(tc.offset)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGridNeighbors_RandomNeighbor(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	t.Run("get random", func(t *testing.T) {
		g := NewGridDefinition(7, 7, utils.UInt8s{4, 5, 6, 7, 6, 5, 4})
		require.NotNil(t, g)

		list := g.withoutSelf[24]

		counts := make(map[uint8]int)
		for ix := 0; ix < 6000; ix++ {
			offs := list.RandomNeighbor(prng)
			counts[offs] = counts[offs] + 1
		}

		assert.Equal(t, 6, len(counts))

		for _, v := range counts {
			assert.GreaterOrEqual(t, v, 900)
			assert.LessOrEqual(t, v, 1100)
		}
	})
}
