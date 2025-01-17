package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	sf1   = NewSymbol(1, WithPayouts(0, 0, 0.5, 2, 4), WithWeights(90, 70, 90, 70, 90, 60, 90))
	sf2   = NewSymbol(2, WithPayouts(0, 0, 0.5, 2, 5), WithWeights(90, 70, 90, 70, 90, 70, 90))
	sf3   = NewSymbol(3, WithPayouts(0, 0, 1, 2.5, 5), WithWeights(70, 90, 70, 90, 70, 50, 70))
	sf4   = NewSymbol(4, WithPayouts(0, 0, 1, 4, 7), WithWeights(70, 90, 70, 90, 70, 30, 70))
	sf5   = NewSymbol(5, WithPayouts(0, 0, 1, 4, 7), WithWeights(50, 70, 50, 70, 50, 60, 50))
	sf6   = NewSymbol(6, WithPayouts(0, 0, 1, 4, 7), WithWeights(50, 70, 50, 70, 50, 70, 50))
	wff1  = NewSymbol(7, WildFor(1, 2), WithPayouts(0, 0, 0.5, 2, 5), WithWeights(20, 10, 20, 10, 20, 10))
	wff2  = NewSymbol(8, WildFor(3, 4), WithPayouts(0, 0, 1, 4, 7), WithWeights(10, 20, 10, 20, 10, 20))
	wwf1  = NewSymbol(9, WithKind(Wild), WithPayouts(0, 0, 1, 4, 10), WithWeights(0, 8, 8, 8, 0, 0, 0))
	hf1   = NewSymbol(10, WithKind(Hero), WithWeights(2, 2, 2, 2, 2, 0, 0))
	scf1  = NewSymbol(11, WithKind(Scatter), WithPayouts(0, 1, 2, 5, 12, 0, 0), WithWeights(2, 2, 2, 2, 2, 0, 0))
	wwf2  = NewSymbol(12, WithKind(Wild), WithPayouts(0, 2, 4, 6, 15), WithWeights(4, 4, 0, 4, 4, 0, 0))
	hf2   = NewSymbol(13, WithKind(Hero), WithWeights(2, 2, 2, 2, 2, 0, 0))
	scf2  = NewSymbol(14, WithKind(Scatter), WithPayouts(0, 1.5, 3, 5, 12), WithWeights(2, 2, 2, 2, 2, 0, 0))
	setF1 = NewSymbolSet(sf1, sf2, sf3, sf4, sf5, sf6, wff1, wff2, wwf1, hf1, scf1, wwf2, hf2, scf2)
)

func TestNewSpin(t *testing.T) {
	testCases := []struct {
		name       string
		reelCount  uint8
		rowCount   uint8
		symbols    *SymbolSet
		noRepeat   uint8
		directions PayDirection
		mask       utils.UInt8s
		zeroes     int
	}{
		{
			name:       "3x3 LTR",
			reelCount:  3,
			rowCount:   3,
			symbols:    setF1,
			noRepeat:   2,
			directions: PayLTR,
		},
		{
			name:       "5x3 RTL",
			reelCount:  5,
			rowCount:   3,
			symbols:    setF1,
			noRepeat:   2,
			directions: PayRTL,
		},
		{
			name:       "5x5 Both",
			reelCount:  5,
			rowCount:   5,
			symbols:    setF1,
			noRepeat:   4,
			directions: PayBoth,
		},
		{
			name:       "5x5 all 3-4-5-4-3",
			reelCount:  5,
			rowCount:   5,
			symbols:    setF1,
			directions: PayLTR,
			mask:       utils.UInt8s{3, 4, 5, 4, 3},
			zeroes:     6,
		},
		{
			name:       "5x5 all 1-2-3-4-5",
			reelCount:  5,
			rowCount:   5,
			symbols:    setF1,
			directions: PayLTR,
			mask:       utils.UInt8s{1, 2, 3, 4, 5},
			zeroes:     10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prng := rng.NewRNG()
			defer prng.ReturnToPool()

			opts := make([]SlotOption, 0)
			opts = append(opts, Grid(tc.reelCount, tc.rowCount))
			opts = append(opts, WithSymbols(tc.symbols))
			opts = append(opts, PayDirections(tc.directions))

			if tc.noRepeat > 0 {
				opts = append(opts, NoRepeat(tc.noRepeat))
			}
			if len(tc.mask) > 0 {
				opts = append(opts, WithMask(tc.mask...))
			}

			slots := NewSlots(opts...)
			require.NotNil(t, slots)

			assert.Equal(t, tc.directions, slots.directions)

			s := AcquireSpin(slots, prng)
			require.NotNil(t, s)
			defer s.Release()

			assert.Equal(t, slots, s.slots)
			assert.Less(t, s.newWilds, uint8(25))
			assert.Less(t, s.newHeroes, uint8(25))
			assert.Less(t, s.newScatters, uint8(25))

			if len(tc.mask) > 0 {
				for reel := range s.reels {
					offset := reel * s.rowCount
					for row := 0; row < s.rowCount; row++ {
						if row < int(s.slots.mask[reel]) {
							assert.NotZero(t, s.indexes[offset+row], row)
						} else {
							assert.Zero(t, s.indexes[offset+row], row)
						}
					}
				}
			} else {
				for _, n := range s.indexes {
					assert.NotZero(t, n)
				}
			}

			counts := make(map[utils.Index]int)
			max := 100000
			for ix := 0; ix < max; ix++ {
				if slots.reelCount > 3 {
					s.Spin(1, 4)
				} else {
					s.Spin(1)
				}
				for _, n := range s.indexes {
					counts[n] = counts[n] + 1
				}
			}

			if len(tc.mask) > 0 {
				assert.Equal(t, max*tc.zeroes, counts[0])
			} else {
				assert.Zero(t, counts[0])
			}

			assert.NotZero(t, counts[1])
			assert.NotZero(t, counts[2])
			assert.NotZero(t, counts[3])
			assert.NotZero(t, counts[4])
			assert.NotZero(t, counts[5])
			assert.NotZero(t, counts[6])
			assert.NotZero(t, counts[7])
			assert.NotZero(t, counts[8])
			assert.NotZero(t, counts[9])
			assert.NotZero(t, counts[10])
			assert.NotZero(t, counts[11])
			assert.NotZero(t, counts[12])
			assert.NotZero(t, counts[13])
			assert.NotZero(t, counts[14])

			for ix := utils.Index(15); ix < utils.MaxIndex; ix++ {
				assert.Zero(t, counts[ix], ix)
			}
			assert.Zero(t, counts[utils.MaxIndex], utils.MaxIndex)

			s.LockReels()

			for ix := 0; ix < 10000; ix++ {
				s.spin(s.reels)
				offset := 0
				for reel := 0; reel < s.reelCount; reel++ {
					m := int(s.mask[reel])
					for row := 0; row < s.rowCount; row++ {
						if row < m {
							assert.NotZero(t, s.indexes[offset+row])
							s.indexes[offset+row] = 0
						} else {
							assert.Zero(t, s.indexes[offset+row])
						}
					}
					offset += s.rowCount
				}
			}
		})
	}
}

func TestSpin_LockReels(t *testing.T) {
	t.Run("lock reels", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		slots := NewSlots(Grid(5, 5), WithSymbols(setF1), NoRepeat(4))
		require.NotNil(t, slots)

		s := AcquireSpin(slots, prng)
		require.NotNil(t, s)
		defer s.Release()

		for ix := range s.reels {
			assert.False(t, s.locked[ix])
		}

		s.LockReels(3, 4)
		assert.False(t, s.locked[0])
		assert.False(t, s.locked[1])
		assert.True(t, s.locked[2])
		assert.True(t, s.locked[3])
		assert.False(t, s.locked[4])
		assert.EqualValues(t, utils.UInt8s{3, 4}, s.Locked(nil))

		s.LockReels(2, 4)
		assert.False(t, s.locked[0])
		assert.True(t, s.locked[1])
		assert.False(t, s.locked[2])
		assert.True(t, s.locked[3])
		assert.False(t, s.locked[4])
		assert.EqualValues(t, utils.UInt8s{2, 4}, s.Locked(nil))

		s.LockReels(3)
		assert.False(t, s.locked[0])
		assert.False(t, s.locked[1])
		assert.True(t, s.locked[2])
		assert.False(t, s.locked[3])
		assert.False(t, s.locked[4])
		assert.EqualValues(t, utils.UInt8s{3}, s.Locked(nil))
	})
}

func TestSpin_CountSymbol(t *testing.T) {
	t.Run("spin count symbol", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		slots := NewSlots(Grid(5, 5), WithSymbols(setF1), NoRepeat(4))
		require.NotNil(t, slots)

		s := AcquireSpin(slots, prng)
		require.NotNil(t, s)
		defer s.Release()

		s.indexes = utils.Indexes{1, 2, 3, 4, 2, 1, 2, 3, 1, 6, 1, 2, 3, 6, 7, 1, 2, 6, 7, 8, 1, 6, 7, 8, 9}

		assert.Equal(t, uint8(6), s.CountSymbol(1))
		assert.Equal(t, uint8(5), s.CountSymbol(2))
		assert.Equal(t, uint8(3), s.CountSymbol(3))
		assert.Equal(t, uint8(1), s.CountSymbol(4))
		assert.Equal(t, uint8(0), s.CountSymbol(5))
		assert.Equal(t, uint8(4), s.CountSymbol(6))
		assert.Equal(t, uint8(3), s.CountSymbol(7))
		assert.Equal(t, uint8(2), s.CountSymbol(8))
		assert.Equal(t, uint8(1), s.CountSymbol(9))
		assert.Zero(t, s.CountSymbol(0))
		assert.Zero(t, s.CountSymbol(10))
	})
}

func TestSpin_CountBonusSymbol(t *testing.T) {
	t.Run("spin count bonus symbol", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		slots := NewSlots(Grid(5, 5), WithSymbols(setF1), NoRepeat(4))
		require.NotNil(t, slots)

		s := AcquireSpin(slots, prng)
		require.NotNil(t, s)
		defer s.Release()

		s.indexes = utils.Indexes{1, 2, 3, 4, 2, 1, 2, 3, 1, 6, 1, 2, 3, 6, 7, 1, 2, 6, 7, 8, 1, 6, 7, 8, 9}

		s.bonusSymbol = 1
		assert.Equal(t, uint8(5), s.CountBonusSymbol())
		s.bonusSymbol = 2
		assert.Equal(t, uint8(4), s.CountBonusSymbol())
		s.bonusSymbol = 3
		assert.Equal(t, uint8(3), s.CountBonusSymbol())
		s.bonusSymbol = 4
		assert.Equal(t, uint8(1), s.CountBonusSymbol())
		s.bonusSymbol = 5
		assert.Equal(t, uint8(0), s.CountBonusSymbol())
		s.bonusSymbol = 6
		assert.Equal(t, uint8(4), s.CountBonusSymbol())
		s.bonusSymbol = 7
		assert.Equal(t, uint8(3), s.CountBonusSymbol())
		s.bonusSymbol = 8
		assert.Equal(t, uint8(2), s.CountBonusSymbol())
		s.bonusSymbol = 9
		assert.Equal(t, uint8(1), s.CountBonusSymbol())
		s.bonusSymbol = 0
		assert.Zero(t, s.CountBonusSymbol())
		s.bonusSymbol = 10
		assert.Zero(t, s.CountBonusSymbol())
	})
}

func TestSpin_MismatchPaylines(t *testing.T) {
	t.Run("spin mismatch paylines", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		slots := NewSlots(Grid(5, 5), WithSymbols(setF1), NoRepeat(4), PayDirections(PayBoth))
		require.NotNil(t, slots)

		s := AcquireSpin(slots, prng)
		require.NotNil(t, s)
		defer s.Release()

		s.indexes = utils.Indexes{1, 2, 3, 4, 2, 1, 2, 3, 4, 5, 9, 12, 9, 12, 9, 5, 6, 7, 8, 9, 5, 6, 7, 8, 9}
		s.MismatchPaylines()

		assert.Zero(t, s.CountSymbol(7))
		assert.Zero(t, s.CountSymbol(8))
		assert.Zero(t, s.CountSymbol(9))
		assert.Zero(t, s.CountSymbol(12))

		assert.NotEqual(t, s.indexes[0], s.indexes[5])
		assert.NotEqual(t, s.indexes[0], s.indexes[6])
		assert.NotEqual(t, s.indexes[0], s.indexes[7])
		assert.NotEqual(t, s.indexes[0], s.indexes[8])
		assert.NotEqual(t, s.indexes[0], s.indexes[9])

		assert.NotEqual(t, s.indexes[1], s.indexes[5])
		assert.NotEqual(t, s.indexes[1], s.indexes[6])
		assert.NotEqual(t, s.indexes[1], s.indexes[7])
		assert.NotEqual(t, s.indexes[1], s.indexes[8])
		assert.NotEqual(t, s.indexes[1], s.indexes[9])

		assert.NotEqual(t, s.indexes[2], s.indexes[5])
		assert.NotEqual(t, s.indexes[2], s.indexes[6])
		assert.NotEqual(t, s.indexes[2], s.indexes[7])
		assert.NotEqual(t, s.indexes[2], s.indexes[8])
		assert.NotEqual(t, s.indexes[2], s.indexes[9])

		assert.NotEqual(t, s.indexes[3], s.indexes[5])
		assert.NotEqual(t, s.indexes[3], s.indexes[6])
		assert.NotEqual(t, s.indexes[3], s.indexes[7])
		assert.NotEqual(t, s.indexes[3], s.indexes[8])
		assert.NotEqual(t, s.indexes[3], s.indexes[9])

		assert.NotEqual(t, s.indexes[4], s.indexes[5])
		assert.NotEqual(t, s.indexes[4], s.indexes[6])
		assert.NotEqual(t, s.indexes[4], s.indexes[7])
		assert.NotEqual(t, s.indexes[4], s.indexes[8])
		assert.NotEqual(t, s.indexes[4], s.indexes[9])

		assert.NotEqual(t, s.indexes[15], s.indexes[20])
		assert.NotEqual(t, s.indexes[15], s.indexes[21])
		assert.NotEqual(t, s.indexes[15], s.indexes[22])
		assert.NotEqual(t, s.indexes[15], s.indexes[23])
		assert.NotEqual(t, s.indexes[15], s.indexes[24])

		assert.NotEqual(t, s.indexes[16], s.indexes[20])
		assert.NotEqual(t, s.indexes[16], s.indexes[21])
		assert.NotEqual(t, s.indexes[16], s.indexes[22])
		assert.NotEqual(t, s.indexes[16], s.indexes[23])
		assert.NotEqual(t, s.indexes[16], s.indexes[24])

		assert.NotEqual(t, s.indexes[17], s.indexes[20])
		assert.NotEqual(t, s.indexes[17], s.indexes[21])
		assert.NotEqual(t, s.indexes[17], s.indexes[22])
		assert.NotEqual(t, s.indexes[17], s.indexes[23])
		assert.NotEqual(t, s.indexes[17], s.indexes[24])

		assert.NotEqual(t, s.indexes[18], s.indexes[20])
		assert.NotEqual(t, s.indexes[18], s.indexes[21])
		assert.NotEqual(t, s.indexes[18], s.indexes[22])
		assert.NotEqual(t, s.indexes[18], s.indexes[23])
		assert.NotEqual(t, s.indexes[18], s.indexes[24])

		assert.NotEqual(t, s.indexes[19], s.indexes[20])
		assert.NotEqual(t, s.indexes[19], s.indexes[21])
		assert.NotEqual(t, s.indexes[19], s.indexes[22])
		assert.NotEqual(t, s.indexes[19], s.indexes[23])
		assert.NotEqual(t, s.indexes[19], s.indexes[24])
	})
}

func TestSpin_ForcePaidTrigger(t *testing.T) {
	t.Run("spin force paid trigger", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		slots := NewSlots(Grid(5, 5), WithSymbols(setF1), NoRepeat(4), PayDirections(PayBoth))
		require.NotNil(t, slots)

		s := AcquireSpin(slots, prng)
		require.NotNil(t, s)
		defer s.Release()

		tr := NewPaidAction(BonusGame, 0, 100, 11, 3)
		require.NotNil(t, tr)

		s.indexes = utils.Indexes{1, 2, 3, 4, 2, 5, 6, 7, 5, 11, 3, 4, 6, 7, 7, 5, 6, 7, 8, 1, 2, 3, 3, 4, 2}

		s.ForcePaidAction(tr)
		assert.Equal(t, tr.triggerCount, s.CountSymbol(tr.symbol))

		s.ForcePaidAction(tr)
		assert.Equal(t, tr.triggerCount, s.CountSymbol(tr.symbol))
	})
}

func TestSpin_Results(t *testing.T) {
	t.Run("spin results", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		slots := NewSlots(Grid(5, 3), WithSymbols(setF1), NoRepeat(2))
		require.NotNil(t, slots)

		spin := AcquireSpin(slots, prng)
		require.NotNil(t, spin)
		defer spin.Release()

		spin.indexes = utils.Indexes{1, 2, 3, 1, 4, 5, 2, 3, 1, 9, 2, 3, 4, 5, 6}
	})
}

func TestSpin_landFloatingSymbols5x3(t *testing.T) {
	testCases := []struct {
		name    string
		indexes utils.Indexes
		want    utils.Indexes
	}{
		{"none", utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}, utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}},
		{"one (1)", utils.Indexes{0, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}, utils.Indexes{0, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}},
		{"one (2)", utils.Indexes{1, 2, 3, 1, 2, 0, 1, 2, 3, 1, 2, 3, 1, 2, 3}, utils.Indexes{1, 2, 3, 0, 1, 2, 1, 2, 3, 1, 2, 3, 1, 2, 3}},
		{"one (3)", utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 0, 2, 3, 1, 2, 3}, utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 0, 2, 3, 1, 2, 3}},
		{"one (4)", utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 0}, utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 0, 1, 2}},
		{"few (1)", utils.Indexes{0, 2, 3, 0, 2, 3, 0, 2, 3, 1, 2, 3, 1, 2, 3}, utils.Indexes{0, 2, 3, 0, 2, 3, 0, 2, 3, 1, 2, 3, 1, 2, 3}},
		{"few (2)", utils.Indexes{1, 2, 3, 1, 2, 0, 0, 0, 3, 1, 2, 0, 1, 2, 3}, utils.Indexes{1, 2, 3, 0, 1, 2, 0, 0, 3, 0, 1, 2, 1, 2, 3}},
		{"few (3)", utils.Indexes{1, 2, 3, 0, 0, 3, 1, 0, 3, 1, 0, 3, 0, 2, 3}, utils.Indexes{1, 2, 3, 0, 0, 3, 0, 1, 3, 0, 1, 3, 0, 2, 3}},
		{"few (4)", utils.Indexes{0, 2, 0, 1, 0, 3, 0, 2, 3, 1, 2, 3, 1, 2, 0}, utils.Indexes{0, 0, 2, 0, 1, 3, 0, 2, 3, 1, 2, 3, 0, 1, 2}},
		{"many (1)", utils.Indexes{0, 0, 0, 0, 2, 3, 0, 2, 3, 0, 0, 0, 1, 2, 3}, utils.Indexes{0, 0, 0, 0, 2, 3, 0, 2, 3, 0, 0, 0, 1, 2, 3}},
		{"many (2)", utils.Indexes{1, 2, 0, 0, 2, 0, 0, 0, 3, 1, 2, 0, 0, 0, 3}, utils.Indexes{0, 1, 2, 0, 0, 2, 0, 0, 3, 0, 1, 2, 0, 0, 3}},
		{"many (3)", utils.Indexes{1, 0, 3, 0, 0, 3, 1, 0, 0, 0, 2, 3, 0, 0, 3}, utils.Indexes{0, 1, 3, 0, 0, 3, 0, 0, 1, 0, 2, 3, 0, 0, 3}},
		{"many (4)", utils.Indexes{0, 0, 0, 1, 0, 3, 0, 2, 0, 1, 0, 0, 1, 0, 0}, utils.Indexes{0, 0, 0, 0, 1, 3, 0, 0, 2, 0, 0, 1, 0, 0, 1}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prng := rng.NewRNG()
			defer prng.ReturnToPool()

			slots := NewSlots(Grid(5, 3), WithSymbols(setF1))
			require.NotNil(t, slots)

			spin := AcquireSpin(slots, prng)
			require.NotNil(t, spin)
			defer spin.Release()

			spin.indexes = tc.indexes
			spin.CascadeFloatingSymbols()
			assert.EqualValues(t, tc.want, spin.indexes)
		})
	}
}

func TestSpin_landFloatingSymbols3x7(t *testing.T) {
	testCases := []struct {
		name    string
		indexes utils.Indexes
		want    utils.Indexes
	}{
		{
			name:    "none",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7},
			want:    utils.Indexes{1, 2, 3, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7},
		},
		{
			name:    "one (1)",
			indexes: utils.Indexes{0, 2, 3, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7},
			want:    utils.Indexes{0, 2, 3, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7},
		},
		{
			name:    "one (2)",
			indexes: utils.Indexes{1, 2, 0, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7},
			want:    utils.Indexes{0, 1, 2, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7},
		},
		{
			name:    "one (3)",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 0, 1, 2, 3, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7},
			want:    utils.Indexes{0, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7},
		},
		{
			name:    "few (1)",
			indexes: utils.Indexes{1, 2, 0, 4, 0, 6, 7, 1, 2, 3, 4, 0, 6, 7, 1, 2, 3, 4, 5, 6, 7},
			want:    utils.Indexes{0, 0, 1, 2, 4, 6, 7, 0, 1, 2, 3, 4, 6, 7, 1, 2, 3, 4, 5, 6, 7},
		},
		{
			name:    "few (2)",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7, 1, 2, 0, 4, 0, 6, 0},
			want:    utils.Indexes{1, 2, 3, 4, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7, 0, 0, 0, 1, 2, 4, 6},
		},
		{
			name:    "few (3)",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 0, 0, 1, 2, 0, 0, 5, 0, 0, 1, 0, 0, 4, 5, 6, 7},
			want:    utils.Indexes{0, 0, 1, 2, 3, 4, 5, 0, 0, 0, 0, 1, 2, 5, 0, 0, 1, 4, 5, 6, 7},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prng := rng.NewRNG()
			defer prng.ReturnToPool()

			slots := NewSlots(Grid(3, 7), WithSymbols(setF1))
			require.NotNil(t, slots)

			spin := AcquireSpin(slots, prng)
			require.NotNil(t, spin)
			defer spin.Release()

			spin.indexes = tc.indexes
			spin.CascadeFloatingSymbols()
			assert.EqualValues(t, tc.want, spin.indexes)
		})
	}
}

func TestSpin_refillEmptyPositions(t *testing.T) {
	testCases := []struct {
		name    string
		indexes utils.Indexes
	}{
		{"none", utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}},
		{"one (1)", utils.Indexes{0, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}},
		{"one (2)", utils.Indexes{1, 2, 3, 1, 2, 0, 1, 2, 3, 1, 2, 3, 1, 2, 3}},
		{"one (3)", utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 0, 2, 3, 1, 2, 3}},
		{"one (4)", utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 0}},
		{"few (1)", utils.Indexes{0, 2, 3, 0, 2, 3, 0, 2, 3, 1, 2, 3, 1, 2, 3}},
		{"few (2)", utils.Indexes{1, 2, 3, 1, 2, 0, 0, 0, 3, 1, 2, 0, 1, 2, 3}},
		{"few (3)", utils.Indexes{1, 2, 3, 0, 0, 3, 1, 0, 3, 0, 2, 3, 0, 2, 3}},
		{"few (4)", utils.Indexes{0, 2, 0, 1, 0, 3, 0, 2, 3, 1, 2, 3, 1, 2, 0}},
		{"many (1)", utils.Indexes{0, 0, 0, 0, 2, 3, 0, 2, 3, 0, 0, 0, 1, 2, 3}},
		{"many (2)", utils.Indexes{1, 2, 0, 0, 2, 0, 0, 0, 3, 1, 2, 0, 0, 0, 3}},
		{"many (3)", utils.Indexes{1, 0, 3, 0, 0, 3, 1, 0, 0, 0, 2, 3, 0, 0, 3}},
		{"many (4)", utils.Indexes{0, 0, 0, 1, 0, 3, 0, 2, 3, 0, 0, 0, 1, 0, 0}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prng := rng.NewRNG()
			defer prng.ReturnToPool()

			slots := NewSlots(Grid(5, 3), WithSymbols(setF1))
			require.NotNil(t, slots)

			spin := AcquireSpin(slots, prng)
			require.NotNil(t, spin)
			defer spin.Release()

			spin.indexes = tc.indexes
			spin.refill(spin.reels)

			for ix := range tc.indexes {
				if tc.indexes[ix] == 0 {
					assert.NotZero(t, spin.indexes[ix])
				} else {
					assert.Equal(t, tc.indexes[ix], spin.indexes[ix])
				}
			}
		})
	}
}

func TestSpin_Refill(t *testing.T) {
	testCases := []struct {
		name      string
		indexes   utils.Indexes
		cascading bool
		want      utils.Indexes
	}{
		{"none", utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}, true, utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}},
		{"one (1)", utils.Indexes{0, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}, true, utils.Indexes{0, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}},
		{"one (2)", utils.Indexes{1, 2, 3, 1, 2, 0, 1, 2, 3, 1, 2, 3, 1, 2, 3}, true, utils.Indexes{1, 2, 3, 0, 1, 2, 1, 2, 3, 1, 2, 3, 1, 2, 3}},
		{"one (3)", utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 0, 2, 3, 1, 2, 3}, true, utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 0, 2, 3, 1, 2, 3}},
		{"one (4)", utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 0}, false, utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 0}},
		{"few (1)", utils.Indexes{0, 2, 3, 0, 2, 3, 0, 2, 3, 1, 2, 3, 1, 2, 3}, true, utils.Indexes{0, 2, 3, 0, 2, 3, 0, 2, 3, 1, 2, 3, 1, 2, 3}},
		{"few (2)", utils.Indexes{1, 2, 3, 1, 2, 0, 0, 0, 3, 1, 2, 0, 1, 2, 3}, true, utils.Indexes{1, 2, 3, 0, 1, 2, 0, 0, 3, 0, 1, 2, 1, 2, 3}},
		{"few (3)", utils.Indexes{1, 2, 3, 0, 0, 3, 1, 0, 3, 1, 0, 3, 0, 2, 3}, true, utils.Indexes{1, 2, 3, 0, 0, 3, 0, 1, 3, 0, 1, 3, 0, 2, 3}},
		{"few (4)", utils.Indexes{0, 2, 0, 1, 0, 3, 0, 2, 3, 1, 2, 3, 1, 2, 0}, false, utils.Indexes{0, 2, 0, 1, 0, 3, 0, 2, 3, 1, 2, 3, 1, 2, 0}},
		{"many (1)", utils.Indexes{0, 0, 0, 0, 2, 3, 0, 2, 3, 0, 0, 0, 1, 2, 3}, true, utils.Indexes{0, 0, 0, 0, 2, 3, 0, 2, 3, 0, 0, 0, 1, 2, 3}},
		{"many (2)", utils.Indexes{1, 2, 0, 0, 2, 0, 0, 0, 3, 1, 2, 0, 0, 0, 3}, true, utils.Indexes{0, 1, 2, 0, 0, 2, 0, 0, 3, 0, 1, 2, 0, 0, 3}},
		{"many (3)", utils.Indexes{1, 0, 3, 0, 0, 3, 1, 0, 0, 0, 2, 3, 0, 0, 3}, true, utils.Indexes{0, 1, 3, 0, 0, 3, 0, 0, 1, 0, 2, 3, 0, 0, 3}},
		{"many (4)", utils.Indexes{0, 0, 0, 1, 0, 3, 0, 2, 0, 1, 0, 0, 1, 0, 0}, false, utils.Indexes{0, 0, 0, 1, 0, 3, 0, 2, 0, 1, 0, 0, 1, 0, 0}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prng := rng.NewRNG()
			defer prng.ReturnToPool()

			var slots *Slots
			if tc.cascading {
				slots = NewSlots(Grid(5, 3), WithSymbols(setF1), CascadingReels(true))
			} else {
				slots = NewSlots(Grid(5, 3), WithSymbols(setF1))
			}
			require.NotNil(t, slots)

			spin := AcquireSpin(slots, prng)
			require.NotNil(t, spin)
			defer spin.Release()

			spin.indexes = tc.indexes

			if slots.CascadingReels() {
				spin.CascadeFloatingSymbols()
			}

			spin.Refill()

			for ix := range tc.indexes {
				if tc.want[ix] == 0 {
					assert.NotZero(t, spin.indexes[ix])
				} else {
					assert.Equal(t, tc.want[ix], spin.indexes[ix])
				}
			}
		})
	}
}

func TestSpin_Hot(t *testing.T) {
	slots := NewSlots(Grid(5, 3), WithSymbols(set3))

	testCases := []struct {
		name  string
		reels utils.UInt8s
		want  utils.UInt8s
	}{
		{"none", utils.UInt8s{}, utils.UInt8s{}},
		{"1", utils.UInt8s{1}, utils.UInt8s{2}},
		{"2", utils.UInt8s{2}, utils.UInt8s{3}},
		{"4", utils.UInt8s{4}, utils.UInt8s{5}},
		{"1+2", utils.UInt8s{2, 1}, utils.UInt8s{2, 3}},
		{"1+4", utils.UInt8s{4, 1}, utils.UInt8s{2, 5}},
		{"1+2+4", utils.UInt8s{4, 2, 1}, utils.UInt8s{2, 3, 5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prng := rng.NewRNG()
			defer prng.ReturnToPool()

			spin := AcquireSpin(slots, prng)
			require.NotNil(t, spin)
			defer spin.Release()

			spin.Spin()

			for _, reel := range tc.reels {
				spin.HotReel(reel)
			}

			got := spin.Hot(make(utils.UInt8s, 5))
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestSpin_Debug(t *testing.T) {
	t.Run("spin debug", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

		spin := AcquireSpin(slots, prng)
		defer spin.Release()

		spin.LockReels(1, 3)

		indexes := utils.Indexes{10, 2, 3, 4, 5, 6, 10, 2, 3, 7, 8, 9, 3, 2, 11}
		spin.Debug(indexes)

		assert.Equal(t, []bool{false, false, false, false, false}, spin.locked)
		assert.Equal(t, indexes, spin.indexes)
		assert.Equal(t, uint8(1), spin.newWilds)
		assert.Equal(t, uint8(2), spin.newHeroes)
		assert.Equal(t, uint8(1), spin.newScatters)
	})
}

func TestSpin_BonusSymbol(t *testing.T) {
	t.Run("bonus symbol", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		slots := NewSlots(Grid(5, 3), WithSymbols(set3))

		spin := AcquireSpin(slots, prng)
		require.NotNil(t, spin)
		defer spin.Release()

		assert.Equal(t, utils.MaxIndex, spin.BonusSymbol())
		assert.False(t, spin.altActive)

		spin.SetBonusSymbol(2, false)
		assert.Equal(t, utils.Index(2), spin.BonusSymbol())
		assert.False(t, spin.altActive)

		spin.SetBonusSymbol(5, true)
		assert.Equal(t, utils.Index(5), spin.BonusSymbol())
		assert.True(t, spin.altActive)
	})
}

func TestSpin_AltSymbols(t *testing.T) {
	t.Run("spin with alternate symbols", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		slots := NewSlots(Grid(5, 3), WithSymbols(set1), WithAltSymbols(set3))

		spin := AcquireSpin(slots, prng)
		require.NotNil(t, spin)
		defer spin.Release()

		spin.SetBonusSymbol(2, true)

		counts := make(map[utils.Index]int)
		for ix := 0; ix < 10000; ix++ {
			spin.Spin()

			for _, id := range spin.indexes {
				counts[id] = counts[id] + 1
			}
		}

		for _, symbol := range set3.symbols {
			assert.NotZero(t, counts[symbol.id])
		}
	})
}

func TestSpin_TestChance2(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name   string
		chance float64
		want   bool
	}{
		{name: "no chance", chance: 0},
		{name: "always", chance: 100, want: true},
		{name: "small chance", chance: 5, want: true},
		{name: "big chance", chance: 95, want: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var count int
			for ix := 0; ix < 10000; ix++ {
				if spin.TestChance2(tc.chance) {
					count++
				}
			}
			assert.Equal(t, tc.want, count > 0)
		})
	}
}

func TestSpin_TestChance4(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name   string
		chance float64
		want   bool
	}{
		{name: "no chance", chance: 0},
		{name: "always", chance: 100, want: true},
		{name: "small chance", chance: 1, want: true},
		{name: "big chance", chance: 99, want: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var count int
			for ix := 0; ix < 10000; ix++ {
				if spin.TestChance4(tc.chance) {
					count++
				}
			}
			assert.Equal(t, tc.want, count > 0)
		})
	}
}

func BenchmarkSpin_Spin5x3Repeat(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(set3))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	for i := 0; i < b.N; i++ {
		s.spin(s.reels)
	}
}

func BenchmarkSpin_Spin5x3NoRepeat(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(set3), NoRepeat(2))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	for i := 0; i < b.N; i++ {
		s.spin(s.reels)
	}
}

func BenchmarkSpin_Spin5x5Repeat(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 5), WithSymbols(set3))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	for i := 0; i < b.N; i++ {
		s.spin(s.reels)
	}
}

func BenchmarkSpin_Spin5x5NoRepeat(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 5), WithSymbols(set3), NoRepeat(4))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	for i := 0; i < b.N; i++ {
		s.spin(s.reels)
	}
}

func BenchmarkSpin_CountSingleSpecial5g5x5(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 5), WithSymbols(setF1), PayDirections(PayBoth))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	s.indexes = utils.Indexes{1, 2, 3, 4, 5, 3, 4, 5, 1, 2, 1, 2, 3, 4, 5, 3, 4, 5, 1, 2, 1, 2, 3, 4, 5}

	for i := 0; i < b.N; i++ {
		if sym := s.symbols.GetSymbol(1); sym != nil {
			if sym.IsWild() {
				s.newWilds++
			} else if sym.IsHero() {
				s.newHeroes++
			}
			if sym.IsScatter() {
				s.newScatters++
			}
		}
	}
}

func BenchmarkSpin_CountSpecialRow5g5x5Both(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 5), WithSymbols(setF1), PayDirections(PayBoth))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	s.indexes = utils.Indexes{1, 2, 3, 4, 5, 3, 4, 5, 1, 2, 1, 2, 3, 4, 5, 3, 4, 5, 1, 2, 1, 2, 3, 4, 5}

	for i := 0; i < b.N; i++ {
		s.testSpecialsOnReel(0, 5)
	}
}

func BenchmarkSpin_CountNewSpecials5g5x5Both(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 5), WithSymbols(setF1), PayDirections(PayBoth))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	s.indexes = utils.Indexes{1, 2, 3, 4, 5, 3, 4, 5, 1, 2, 1, 2, 3, 4, 5, 3, 4, 5, 1, 2, 1, 2, 3, 4, 5}

	for i := 0; i < b.N; i++ {
		s.CountSpecials()
	}
}

func BenchmarkSpin_Result3x3LTR(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(3, 3), WithSymbols(setF1), PayDirections(PayLTR))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	for i := 0; i < b.N; i++ {
		s.Spin()
	}
}

func BenchmarkSpin_Result5x3LTR(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1), PayDirections(PayLTR))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	for i := 0; i < b.N; i++ {
		s.Spin()
	}
}

func BenchmarkSpin_Result5x3LTRnoRepeat(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1), PayDirections(PayLTR), NoRepeat(2))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	for i := 0; i < b.N; i++ {
		s.Spin()
	}
}

func BenchmarkSpin_Result5x3Both(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1), PayDirections(PayBoth))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	for i := 0; i < b.N; i++ {
		s.Spin()
	}
}

func BenchmarkSpin_Result6x3LTR(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(6, 3), WithSymbols(setF1), PayDirections(PayLTR))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	for i := 0; i < b.N; i++ {
		s.Spin()
	}
}

func BenchmarkSpin_Result6x3Both(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(6, 3), WithSymbols(setF1), PayDirections(PayBoth))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	for i := 0; i < b.N; i++ {
		s.Spin()
	}
}

func BenchmarkSpin_Result5x5LTR(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 5), WithSymbols(setF1), PayDirections(PayLTR))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	for i := 0; i < b.N; i++ {
		s.Spin()
	}
}

func BenchmarkSpin_Result5x5LTRnoRepeat(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 5), WithSymbols(setF1), PayDirections(PayLTR), NoRepeat(4))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	for i := 0; i < b.N; i++ {
		s.Spin()
	}
}

func BenchmarkSpin_Result5x5Both(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 5), WithSymbols(setF1), PayDirections(PayBoth))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	for i := 0; i < b.N; i++ {
		s.Spin()
	}
}

func BenchmarkSpin_Result6x4LTR(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(6, 4), WithSymbols(setF1), PayDirections(PayLTR))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	for i := 0; i < b.N; i++ {
		s.Spin()
	}
}

func BenchmarkSpin_Result6x4LTRnoRepeat(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(6, 4), WithSymbols(setF1), PayDirections(PayLTR), NoRepeat(3))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	for i := 0; i < b.N; i++ {
		s.Spin()
	}
}

func BenchmarkSpin_Result6x4Both(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(6, 4), WithSymbols(setF1), PayDirections(PayBoth))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	for i := 0; i < b.N; i++ {
		s.Spin()
	}
}
