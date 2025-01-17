package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewSlots(t *testing.T) {
	testCases := []struct {
		name           string
		reelCount      uint8
		rowCount       uint8
		symbols        *SymbolSet
		noRepeat       uint8
		directions     PayDirection
		highest        bool
		cascadingReels bool
		mask           utils.UInt8s
		maxPayout      float64
		targetRTP      float64
	}{
		{
			name:       "set1 5x3 RTL",
			reelCount:  5,
			rowCount:   3,
			symbols:    set1,
			directions: PayLTR,
		},
		{
			name:       "set1 3x3 LTR",
			reelCount:  3,
			rowCount:   3,
			symbols:    set1,
			directions: PayRTL,
		},
		{
			name:           "set1 5x3 cascading reels",
			reelCount:      5,
			rowCount:       3,
			symbols:        set1,
			directions:     PayRTL,
			cascadingReels: true,
		},
		{
			name:       "set2 5x5 Both max 20000x",
			reelCount:  5,
			rowCount:   5,
			symbols:    set2,
			noRepeat:   4,
			directions: PayBoth,
			highest:    true,
			maxPayout:  20000,
		},
		{
			name:       "set2 5x5 all 3-4-5-4-3",
			reelCount:  5,
			rowCount:   5,
			symbols:    set2,
			noRepeat:   4,
			directions: PayLTR,
			mask:       utils.UInt8s{3, 4, 5, 4, 3},
		},
		{
			name:       "with rtp",
			reelCount:  5,
			rowCount:   3,
			symbols:    set2,
			noRepeat:   2,
			directions: PayLTR,
			targetRTP:  96.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := make([]SlotOption, 0, 8)
			if tc.reelCount != defaultReelCount || tc.rowCount != defaultRowCount {
				opts = append(opts, Grid(tc.reelCount, tc.rowCount))
			}
			opts = append(opts, WithSymbols(tc.symbols))
			if tc.noRepeat != defaultNoRepeat {
				opts = append(opts, NoRepeat(tc.noRepeat))
			}
			opts = append(opts, PayDirections(tc.directions))
			if tc.highest != defaultHighestPayout {
				opts = append(opts, HighestPayout())
			}
			if tc.cascadingReels != defaultCascadingReels {
				opts = append(opts, CascadingReels(true))
			}
			if len(tc.mask) > 0 {
				opts = append(opts, WithMask(tc.mask...))
			}
			if tc.maxPayout > 0.0 {
				opts = append(opts, MaxPayout(tc.maxPayout))
			}
			if tc.targetRTP > 0.0 {
				opts = append(opts, WithRTP(tc.targetRTP))
			}

			s := NewSlots(opts...)
			require.NotNil(t, s)

			assert.Equal(t, int(tc.reelCount), s.ReelCount())
			assert.Equal(t, int(tc.rowCount), s.RowCount())
			assert.Equal(t, tc.symbols, s.Symbols())

			assert.Equal(t, tc.directions, s.directions)
			assert.Equal(t, tc.highest, s.highestPayout)
			assert.Equal(t, tc.cascadingReels, s.cascadingReels)
			assert.Equal(t, tc.maxPayout, s.maxPayout)

			target := s.RTP()
			assert.Equal(t, tc.targetRTP, target)

			if len(tc.mask) > 0 {
				assert.EqualValues(t, tc.mask, s.mask)
			} else {
				assert.Zero(t, len(s.mask))
			}
		})
	}
}

func TestNewSlotsFail(t *testing.T) {
	t.Run("new slots fail grid size", func(t *testing.T) {
		defer func() {
			e := recover()
			require.NotNil(t, e)
		}()

		s := NewSlots(Grid(0, 1))
		require.Nil(t, s)
	})
}

func TestSlots_GetBonusSymbol(t *testing.T) {
	t.Run("get bonus symbol", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		symbols := NewSymbolSet(sf1, sf2, sf3, sf4, sf5, sf6, wwf1, hf1, scf1)
		require.NotNil(t, symbols)

		w := utils.AcquireWeighting()
		require.NotNil(t, w)
		defer w.Release()

		w.AddWeight(sf1.id, 100)
		w.AddWeight(sf2.id, 80)
		w.AddWeight(sf3.id, 60)
		w.AddWeight(sf4.id, 40)
		w.AddWeight(sf5.id, 20)
		w.AddWeight(sf6.id, 20)
		w.AddWeight(hf1.id, 5)
		symbols.SetBonusWeights(w)

		slots := NewSlots(Grid(5, 3), NoRepeat(2), WithSymbols(symbols))
		require.NotNil(t, slots)

		counts := make(map[utils.Index]int)
		for ix := 0; ix < 10000; ix++ {
			n := slots.symbols.GetBonusSymbol(prng)
			counts[n] = counts[n] + 1
		}

		assert.Zero(t, counts[0])
		assert.NotZero(t, counts[sf1.id])
		assert.NotZero(t, counts[sf2.id])
		assert.NotZero(t, counts[sf3.id])
		assert.NotZero(t, counts[sf4.id])
		assert.NotZero(t, counts[sf5.id])
		assert.NotZero(t, counts[sf6.id])
		assert.Zero(t, counts[wwf1.id])
		assert.NotZero(t, counts[hf1.id])
		assert.Zero(t, counts[scf1.id])
	})
}

func TestSlots_WithSpinner(t *testing.T) {
	t.Run("slots with spinner", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		symbols := NewSymbolSet(sf1, sf2, sf3, sf4, sf5, sf6, wwf1, hf1, scf1)
		require.NotNil(t, symbols)

		r1 := NewSymbolReel(3, reel1...)
		r2 := NewSymbolReel(3, reel2...)
		r3 := NewSymbolReel(3, reel3...)
		r4 := NewSymbolReel(3, reel4...)
		r5 := NewSymbolReel(3, reel5...)

		sr := NewSymbolReels(r1, r2, r3, r4, r5)
		require.NotNil(t, sr)

		slots := NewSlots(Grid(5, 3), WithSymbols(symbols), WithSpinner(sr))
		require.NotNil(t, slots)

		spin := AcquireSpin(slots, prng)
		require.NotNil(t, spin)

		spin.Spin()

		temp := make(utils.Indexes, 3)
		var offs int

		for r := uint8(1); r <= 5; r++ {
			copy(temp, spin.indexes[offs:offs+3])

			reel := sr.Reel(r)
			reelX := append(reel, reel[0], reel[1])

			var found bool
			for ix := 0; ix < len(reel)-2; ix++ {
				found = true
				for iy := range temp {
					if temp[iy] != reelX[ix+iy] {
						found = false
						break
					}
				}
				if found {
					break
				}
			}

			assert.True(t, found)

			offs += 3
		}
	})
}
