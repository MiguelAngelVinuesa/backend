package slots

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

var (
	reelShort = utils.Indexes{1, 2, 3, 4, 1, 5, 6, 2, 3, 1, 4, 2, 3}
	reel1     = utils.Indexes{4, 5, 1, 7, 9, 8, 6, 1, 7, 8, 3, 6, 4, 0, 9, 5, 1, 6, 8, 3, 9, 5, 1, 6, 4, 2, 8, 7, 3, 4, 5, 2, 6, 7, 3, 8, 9, 0, 7, 4, 3, 9, 5, 1, 6, 8, 2, 4, 5, 1, 6, 7, 3, 5, 9, 1, 8, 9, 2, 6, 4, 0, 7, 8, 2, 5, 4, 1, 9, 7, 3, 6, 9, 2, 4, 5, 3, 6, 8, 9, 7, 2, 8, 6, 1, 5, 4, 3, 8, 7, 1, 9, 8, 2, 7, 9, 0, 6, 7, 3, 4, 5, 1, 9, 8, 2, 4, 7, 1, 5, 8, 2, 9, 6, 3, 1, 5, 9, 3, 7, 6, 1, 4, 8, 1, 2, 9, 4, 1, 6, 8, 2, 7, 5, 1, 3, 9, 8, 3, 7, 4, 1, 5, 6, 1, 2, 4, 9, 3, 8, 6, 2, 5, 7, 0, 1, 8, 9, 1, 7, 4, 1, 6, 5, 2, 3, 4, 7, 1, 5, 8, 2, 9, 6, 1, 3, 5, 9, 3, 7, 6, 1, 4, 8, 2, 1, 9, 4, 1, 6, 8, 2, 7, 5, 3, 9, 8, 3, 7, 4, 1, 5, 6, 2, 4, 9, 3, 8, 6, 2, 5, 7, 1, 0, 8, 9, 1, 7, 4, 1, 6, 5, 3, 2}
	reel2     = utils.Indexes{4, 5, 0, 7, 9, 8, 6, 3, 7, 8, 0, 6, 4, 2, 9, 5, 1, 6, 8, 3, 9, 5, 0, 6, 4, 3, 8, 7, 2, 4, 5, 6, 7, 3, 8, 9, 0, 7, 4, 3, 9, 5, 1, 6, 8, 2, 4, 5, 0, 6, 7, 2, 5, 9, 0, 8, 9, 3, 6, 4, 7, 8, 2, 5, 4, 0, 9, 7, 3, 6, 9, 2, 4, 5, 3, 6, 8, 7, 9, 2, 8, 6, 1, 5, 4, 3, 8, 7, 9, 6, 2, 7, 9, 0, 8, 7, 3, 4, 5, 0, 9, 8, 2, 4, 7, 1, 5, 8, 2, 9, 6, 3, 0, 5, 9, 3, 7, 6, 0, 4, 8, 0, 2, 9, 4, 0, 6, 8, 2, 7, 5, 0, 3, 9, 8, 3, 7, 4, 0, 5, 6, 0, 2, 4, 9, 3, 8, 6, 2, 5, 7, 0, 1, 8, 9, 0, 7, 4, 0, 6, 5, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 0, 3, 5, 9, 3, 7, 6, 0, 4, 8, 2, 9, 4, 0, 6, 8, 2, 7, 5, 3, 0, 9, 8, 3, 7, 4, 0, 5, 6, 2, 4, 9, 3, 8, 6, 2, 5, 7, 1, 0, 8, 9, 0, 7, 4, 0, 6, 5, 2, 3}
	reel3     = utils.Indexes{4, 5, 3, 7, 9, 1, 0, 8, 6, 0, 7, 8, 1, 6, 4, 0, 1, 9, 5, 1, 6, 8, 3, 9, 5, 1, 6, 4, 2, 0, 8, 7, 3, 4, 5, 2, 6, 7, 3, 8, 9, 0, 1, 7, 4, 3, 9, 5, 1, 6, 8, 2, 4, 5, 0, 2, 6, 7, 3, 5, 9, 1, 8, 9, 2, 6, 4, 0, 1, 7, 8, 2, 5, 4, 1, 9, 7, 3, 6, 9, 2, 0, 4, 5, 3, 6, 8, 0, 9, 7, 2, 8, 6, 0, 5, 4, 3, 0, 8, 7, 1, 9, 8, 2, 7, 9, 0, 1, 6, 7, 3, 4, 5, 1, 9, 8, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 3, 2, 5, 9, 3, 7, 6, 1, 4, 8, 0, 2, 9, 4, 1, 6, 8, 2, 7, 5, 0, 3, 9, 8, 3, 7, 4, 0, 5, 6, 1, 2, 4, 9, 3, 8, 6, 2, 5, 7, 0, 1, 8, 9, 0, 7, 4, 1, 6, 5, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 1, 3, 5, 9, 3, 7, 6, 1, 4, 8, 2, 0, 9, 4, 1, 6, 8, 2, 7, 5, 3, 0, 9, 8, 3, 7, 4, 0, 5, 6, 2, 1, 4, 9, 3, 8, 6, 2, 5, 7, 1, 0, 8, 9, 0, 7, 4, 1, 6, 5, 3, 2}
	reel4     = utils.Indexes{4, 5, 3, 7, 9, 1, 0, 8, 6, 0, 7, 8, 1, 6, 4, 0, 1, 9, 5, 1, 6, 8, 3, 9, 5, 1, 6, 4, 2, 0, 8, 7, 3, 4, 5, 2, 6, 7, 3, 8, 9, 0, 1, 7, 4, 3, 9, 5, 1, 6, 8, 2, 4, 5, 0, 2, 6, 7, 3, 5, 9, 1, 8, 9, 2, 6, 4, 0, 1, 7, 8, 2, 5, 4, 1, 9, 7, 3, 6, 9, 2, 0, 4, 5, 3, 6, 8, 0, 9, 7, 2, 8, 6, 0, 5, 4, 3, 0, 8, 7, 1, 9, 8, 2, 7, 9, 0, 1, 6, 7, 3, 4, 5, 1, 9, 8, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 3, 2, 5, 9, 3, 7, 6, 1, 4, 8, 0, 2, 9, 4, 1, 6, 8, 2, 7, 5, 0, 3, 9, 8, 3, 7, 4, 0, 5, 6, 1, 2, 4, 9, 3, 8, 6, 2, 5, 7, 0, 1, 8, 9, 0, 7, 4, 1, 6, 5, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 1, 3, 5, 9, 3, 7, 6, 1, 4, 8, 2, 0, 9, 4, 1, 6, 8, 2, 7, 5, 3, 0, 9, 8, 3, 7, 4, 0, 5, 6, 2, 1, 4, 9, 3, 8, 6, 2, 5, 7, 1, 0, 8, 9, 0, 7, 4, 1, 6, 5, 3, 2}
	reel5     = utils.Indexes{4, 5, 3, 7, 9, 1, 0, 8, 6, 0, 7, 8, 1, 6, 4, 0, 1, 9, 5, 1, 6, 8, 3, 9, 5, 1, 6, 4, 2, 0, 8, 7, 3, 4, 5, 2, 6, 7, 3, 8, 9, 0, 1, 7, 4, 3, 9, 5, 1, 6, 8, 2, 4, 5, 0, 2, 6, 7, 3, 5, 9, 1, 8, 9, 2, 6, 4, 0, 1, 7, 8, 2, 5, 4, 1, 9, 7, 3, 6, 9, 2, 0, 4, 5, 3, 6, 8, 0, 9, 7, 2, 8, 6, 0, 5, 4, 3, 0, 8, 7, 1, 9, 8, 2, 7, 9, 0, 1, 6, 7, 3, 4, 5, 1, 9, 8, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 3, 2, 5, 9, 3, 7, 6, 1, 4, 8, 0, 2, 9, 4, 1, 6, 8, 2, 7, 5, 0, 3, 9, 8, 3, 7, 4, 0, 5, 6, 1, 2, 4, 9, 3, 8, 6, 2, 5, 7, 0, 1, 8, 9, 0, 7, 4, 1, 6, 5, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 1, 3, 5, 9, 3, 7, 6, 1, 4, 8, 2, 0, 9, 4, 1, 6, 8, 2, 7, 5, 3, 0, 9, 8, 3, 7, 4, 0, 5, 6, 2, 1, 4, 9, 3, 8, 6, 2, 5, 7, 1, 0, 8, 9, 0, 7, 4, 1, 6, 5, 3, 2}
)

func TestNewSymbolReel(t *testing.T) {
	testCases := []struct {
		name string
		rows uint8
		reel utils.Indexes
		offs int
		want utils.Indexes
	}{
		{name: "short - 0", rows: 3, reel: reelShort, offs: 0, want: utils.Indexes{1, 2, 3}},
		{name: "short - 1", rows: 3, reel: reelShort, offs: 1, want: utils.Indexes{2, 3, 4}},
		{name: "short - 2", rows: 3, reel: reelShort, offs: 2, want: utils.Indexes{3, 4, 1}},
		{name: "short - 3", rows: 3, reel: reelShort, offs: 3, want: utils.Indexes{4, 1, 5}},
		{name: "short - last-2", rows: 3, reel: reelShort, offs: len(reelShort) - 3, want: utils.Indexes{4, 2, 3}},
		{name: "short - last-1", rows: 3, reel: reelShort, offs: len(reelShort) - 2, want: utils.Indexes{2, 3, 1}},
		{name: "short - last", rows: 3, reel: reelShort, offs: len(reelShort) - 1, want: utils.Indexes{3, 1, 2}},
		{name: "reel1 - 0", rows: 3, reel: reel1, offs: 0, want: utils.Indexes{4, 5, 1}},
		{name: "reel1 - 1", rows: 3, reel: reel1, offs: 1, want: utils.Indexes{5, 1, 7}},
		{name: "reel1 - 2", rows: 3, reel: reel1, offs: 2, want: utils.Indexes{1, 7, 9}},
		{name: "reel1 - 3", rows: 3, reel: reel1, offs: 3, want: utils.Indexes{7, 9, 8}},
		{name: "reel1 - last-2", rows: 3, reel: reel1, offs: len(reel1) - 3, want: utils.Indexes{5, 3, 2}},
		{name: "reel1 - last-1", rows: 3, reel: reel1, offs: len(reel1) - 2, want: utils.Indexes{3, 2, 4}},
		{name: "reel1 - last", rows: 3, reel: reel1, offs: len(reel1) - 1, want: utils.Indexes{2, 4, 5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sr := NewSymbolReel(tc.rows, tc.reel...).WithFlag(1)
			require.NotNil(t, sr)

			assert.Equal(t, sr.flag, 1)
			assert.Equal(t, tc.rows, sr.rows)
			assert.Equal(t, len(tc.reel)+int(tc.rows)-1, len(sr.stops))
			assert.EqualValues(t, tc.reel, sr.Reel())

			got := make(utils.Indexes, tc.rows)
			copy(got, sr.stops[tc.offs:])
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestSymbolReel_Fill(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1), PayDirections(PayLTR))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	t.Run("symbol reel fill", func(t *testing.T) {
		sr := NewSymbolReel(3, reel1...)
		require.NotNil(t, sr)

		combos := len(reel1)
		counts := make(map[int]int, combos)
		reel := make(utils.Indexes, 3)
		runs := 1000000

		for ix := 0; ix < runs; ix++ {
			sr.Fill(spin, reel)
			iy := int(reel[0])*100 + int(reel[1])*10 + int(reel[2])
			counts[iy] = counts[iy] + 1
		}

		assert.LessOrEqual(t, len(counts), combos)

		avg := runs / len(counts)
		low1 := avg * 6 / 10
		low2 := avg * 12 / 10
		low3 := avg * 18 / 10
		low4 := avg * 24 / 10
		high1 := avg * 13 / 10
		high2 := avg * 23 / 10
		high3 := avg * 33 / 10
		high4 := avg * 43 / 10

		for ix := range counts {
			c := counts[ix]

			switch {
			case c > 3*avg:
				assert.GreaterOrEqualf(t, c, low4, fmt.Sprintf("index %d - 4x", ix))
				assert.LessOrEqual(t, c, high4, fmt.Sprintf("index %d", ix))
			case c > 2*avg:
				assert.GreaterOrEqualf(t, c, low3, fmt.Sprintf("index %d - 3x", ix))
				assert.LessOrEqual(t, c, high3, fmt.Sprintf("index %d", ix))
			case c > avg:
				assert.GreaterOrEqualf(t, c, low2, fmt.Sprintf("index %d - 2x", ix))
				assert.LessOrEqual(t, c, high2, fmt.Sprintf("index %d", ix))
			default:
				assert.GreaterOrEqualf(t, c, low1, fmt.Sprintf("index %d - 1x", ix))
				assert.LessOrEqual(t, c, high1, fmt.Sprintf("index %d", ix))
			}
		}
	})
}

func TestNewSymbolReels(t *testing.T) {
	t.Run("new symbol reels", func(t *testing.T) {
		r1 := NewSymbolReel(3, reel1...)
		r2 := NewSymbolReel(3, reel2...)
		r3 := NewSymbolReel(3, reel3...)
		r4 := NewSymbolReel(3, reel4...)
		r5 := NewSymbolReel(3, reel5...)

		sr := NewSymbolReels(r1, r2, r3, r4, r5).WithFlag(2)
		require.NotNil(t, sr)

		assert.Equal(t, sr.flag, 2)
		assert.EqualValues(t, reel1, sr.Reel(1))
		assert.EqualValues(t, reel2, sr.Reel(2))
		assert.EqualValues(t, reel3, sr.Reel(3))
		assert.EqualValues(t, reel4, sr.Reel(4))
		assert.EqualValues(t, reel5, sr.Reel(5))
	})
}
