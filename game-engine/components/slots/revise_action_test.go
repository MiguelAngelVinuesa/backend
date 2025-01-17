package slots

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewMorphSymbolAction(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	rb := NewSymbol(10, WithKind(WildScatter), MorphInto(11))
	yb := NewSymbol(11, WithKind(WildScatter))
	symbols := NewSymbolSet(sf1, sf2, sf3, sf4, sf5, sf6, rb, yb)
	slots := NewSlots(Grid(5, 3), WithSymbols(symbols))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name    string
		reels   utils.UInt8s
		chances []float64
		symbols utils.Indexes
		indexes utils.Indexes
		want    bool
	}{
		{
			name:    "no symbol",
			chances: []float64{25, 10},
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "all reels, wrong symbol",
			chances: []float64{25, 10},
			symbols: utils.Indexes{8, 9},
			indexes: utils.Indexes{11, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 6, 10, 2, 3},
		},
		{
			name:    "all symbols, wrong reel",
			reels:   utils.UInt8s{2, 3, 4},
			chances: []float64{25, 10},
			indexes: utils.Indexes{10, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 10, 2, 3},
		},
		{
			name:    "good reel, wrong symbol",
			reels:   utils.UInt8s{2, 3, 4},
			chances: []float64{25, 10},
			symbols: utils.Indexes{8, 9},
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "good symbol, wrong reel",
			reels:   utils.UInt8s{2, 3, 4},
			chances: []float64{25, 10},
			symbols: utils.Indexes{10},
			indexes: utils.Indexes{10, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 10, 2, 3},
		},
		{
			name:    "good symbol, good reel",
			reels:   utils.UInt8s{2, 3, 4},
			chances: []float64{25, 10},
			symbols: utils.Indexes{10},
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 10, 1, 2, 3},
			want:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewMorphSymbolAction(tc.reels, tc.chances, tc.symbols...)
			require.NotNil(t, a)

			assert.True(t, a.morphSymbols)
			assert.EqualValues(t, tc.reels, a.morphReels)
			assert.EqualValues(t, tc.chances, a.morphChances)
			assert.EqualValues(t, tc.symbols, a.morphFor)

			var count int
			for ix := 0; ix < 1000; ix++ {
				copy(spin.indexes, tc.indexes)
				if m2 := a.Triggered(spin); m2 != nil {
					count++
				}
			}

			if tc.want {
				assert.NotZero(t, count)
				assert.LessOrEqual(t, count, 300)
			} else {
				assert.Zero(t, count)
			}
		})
	}
}

func TestMorphAction_WithSpinKinds(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name  string
		kinds []SpinKind
		kind  SpinKind
		want  bool
	}{
		{
			name: "none, triggered",
			kind: RegularSpin,
			want: true,
		},
		{
			name:  "regular, not triggered",
			kind:  FirstSpin,
			kinds: []SpinKind{RegularSpin},
		},
		{
			name:  "regular, triggered",
			kind:  RegularSpin,
			kinds: []SpinKind{RegularSpin},
			want:  true,
		},
		{
			name:  "first, not triggered",
			kind:  SecondSpin,
			kinds: []SpinKind{FirstSpin},
		},
		{
			name:  "first, triggered",
			kind:  FirstSpin,
			kinds: []SpinKind{FirstSpin},
			want:  true,
		},
		{
			name:  "second, not triggered",
			kind:  FirstSpin,
			kinds: []SpinKind{SecondSpin},
		},
		{
			name:  "second, triggered",
			kind:  SecondSpin,
			kinds: []SpinKind{SecondSpin},
			want:  true,
		},
		{
			name:  "bonus, not triggered",
			kind:  RegularSpin,
			kinds: []SpinKind{FreeSpin, SuperSpin, RefillSpin},
		},
		{
			name:  "bonus, triggered",
			kind:  RefillSpin,
			kinds: []SpinKind{FreeSpin, SuperSpin, RefillSpin},
			want:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewGenerateSymbolAction(13, []float64{80})
			require.NotNil(t, a)

			if len(tc.kinds) > 0 {
				a.WithSpinKinds(tc.kinds)
				assert.EqualValues(t, tc.kinds, a.spinKinds)
			} else {
				assert.Nil(t, a.spinKinds)
			}

			assert.True(t, a.generateSymbol)

			counts := make(map[uint8]int)
			for ix := 0; ix < 100; ix++ {
				copy(spin.indexes, utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3})
				spin.SetKind(tc.kind)

				a.Triggered(spin)
				n := spin.CountSymbol(13)
				counts[n] = counts[n] + 1
			}

			require.NotZero(t, len(counts))
			assert.NotZero(t, counts[0])

			if tc.want {
				assert.Equal(t, 2, len(counts))
				assert.LessOrEqual(t, counts[0], 35)
				assert.GreaterOrEqual(t, counts[1], 65)
			} else {
				assert.Equal(t, 1, len(counts))
				assert.Equal(t, 100, counts[0])
			}
		})
	}
}

func TestNewGenerateSymbolAction(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name        string
		symbol      utils.Index
		chances     []float64
		indexes     utils.Indexes
		multipliers utils.WeightedGenerator
		tries       int
		want        bool
		min         []int
		max         []int
	}{
		{
			name:    "zero chance",
			symbol:  13,
			chances: []float64{0.0},
			indexes: utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3},
			tries:   1000,
			want:    false,
		},
		{
			name:    "small chance of 1",
			symbol:  13,
			chances: []float64{5.0},
			indexes: utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3},
			tries:   1000,
			want:    true,
			min:     []int{10},
			max:     []int{100},
		},
		{
			name:    "large chance of 1",
			symbol:  13,
			chances: []float64{75.0},
			indexes: utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3},
			tries:   1000,
			want:    true,
			min:     []int{650},
			max:     []int{850},
		},
		{
			name:    "small chance of 1 or 2",
			symbol:  13,
			chances: []float64{10.0, 10.0},
			indexes: utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3},
			tries:   10000,
			want:    true,
			min:     []int{800, 70},
			max:     []int{1800, 170},
		},
		{
			name:    "large chance of 1, small chance of 2",
			symbol:  13,
			chances: []float64{75.0, 5.0},
			indexes: utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3},
			tries:   10000,
			want:    true,
			min:     []int{6500, 325},
			max:     []int{8500, 475},
		},
		{
			name:    "large chance of 2",
			symbol:  13,
			chances: []float64{75.0, 50.0},
			indexes: utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3},
			tries:   10000,
			want:    true,
			min:     []int{2500, 3250},
			max:     []int{4000, 4750},
		},
		{
			name:        "large chance with multipliers",
			symbol:      13,
			chances:     []float64{75.0, 50.0},
			indexes:     utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3},
			multipliers: utils.AcquireWeighting().AddWeights(utils.Indexes{1, 2, 3, 4, 5}, []float64{50, 20, 10, 5, 2}),
			tries:       10000,
			want:        true,
			min:         []int{2500, 3250},
			max:         []int{4000, 4750},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewGenerateSymbolAction(tc.symbol, tc.chances)
			require.NotNil(t, a)

			assert.True(t, a.generateSymbol)

			if tc.multipliers != nil {
				a.WithMultipliers(tc.multipliers)
			}

			counts := make(map[uint8]int)
			countMults := make(map[uint16]int)

			for ix := 0; ix < tc.tries; ix++ {
				copy(spin.indexes, tc.indexes)
				spin.resetMultipliers()

				a.Triggered(spin)

				n := spin.CountSymbol(tc.symbol)
				counts[n] = counts[n] + 1

				if tc.multipliers != nil && len(spin.multipliers) > 0 {
					for iy := 0; iy < len(spin.indexes); iy++ {
						if m := spin.multipliers[iy]; m > 0 {
							assert.Equal(t, tc.symbol, spin.indexes[iy])
							assert.GreaterOrEqual(t, m, uint16(2))
							countMults[m] = countMults[m] + 1
						}
					}
				}
			}

			require.NotZero(t, len(counts))
			assert.NotZero(t, counts[0])

			if tc.want {
				assert.Equal(t, len(tc.chances)+1, len(counts))
				for ix := 1; ix < len(counts); ix++ {
					assert.GreaterOrEqual(t, counts[uint8(ix)], tc.min[ix-1], ix)
					assert.LessOrEqual(t, counts[uint8(ix)], tc.max[ix-1], ix)
				}
				if tc.multipliers != nil {
					assert.NotZero(t, len(countMults))
				}
			} else {
				assert.Equal(t, 1, len(counts))
				assert.Equal(t, tc.tries, counts[0])
			}
		})
	}
}

func TestNewGenOrMorphSymbolAction(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name    string
		symbol  utils.Index
		from    utils.Index
		chances []float64
		reels   []uint8
		dupes   bool
		remain  uint64
		indexes utils.Indexes
		free    uint64
		want    bool
	}{
		{
			name:    "wrong free spin",
			symbol:  11,
			from:    10,
			chances: []float64{30},
			reels:   []uint8{2, 3, 4},
			remain:  2,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			free:    0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewGenOrMorphSymbolAction(tc.symbol, tc.from, tc.chances, tc.reels...)
			require.NotNil(t, a)

			if !tc.dupes {
				a.GenerateNoDupes()
			}
			if tc.remain >= 0 {
				a.WithTriggerFilters(OnRemainingFreeSpins(tc.remain))
			}

			assert.True(t, a.morphSymbols)
			assert.EqualValues(t, tc.chances, a.morphChances)
			assert.Equal(t, 1, len(a.morphFor))
			assert.EqualValues(t, tc.reels, a.morphReels)
			assert.True(t, a.generateSymbol)
			assert.Equal(t, tc.symbol, a.symbol)
			assert.EqualValues(t, tc.chances, a.symbolChances)
			assert.EqualValues(t, tc.reels, a.generateReels)
			assert.Equal(t, tc.dupes, a.genAllowDupes)

			var got bool
			for ix := 0; ix < 1000; ix++ {
				copy(spin.indexes, tc.indexes)
				spin.freeSpins = tc.free

				if a.CanTrigger(spin) {
					if a.Triggered(spin) != nil {
						got = true
					}
				}
			}

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestNewGenerateShapeAction(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	centers := GridOffsets{{1, 1}, {2, 1}, {3, 1}}
	weights := utils.AcquireWeighting().AddWeights([]utils.Index{4, 5, 6, 7, 8, 9}, []float64{80, 60, 40, 20, 10, 5})

	testCases := []struct {
		name    string
		shape   GridOffsets
		chance  float64
		indexes utils.Indexes
		tries   int
		want    bool
		min     int
		max     int
	}{
		{
			name:    "no chance",
			shape:   superX,
			chance:  0.0,
			indexes: utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3},
		},
		{
			name:    "small chance",
			shape:   superX,
			chance:  0.5,
			indexes: utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3},
			tries:   100000,
			want:    true,
			min:     400,
			max:     600,
		},
		{
			name:    "reasonable chance",
			shape:   superX,
			chance:  4.5,
			indexes: utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3},
			tries:   100000,
			want:    true,
			min:     4000,
			max:     5000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewGenerateShapeAction(tc.chance, tc.shape, centers, weights)
			require.NotNil(t, a)

			assert.True(t, a.generateShape)

			var count int

			for ix := 0; ix < tc.tries; ix++ {
				copy(spin.indexes, tc.indexes)
				if a.Triggered(spin) != nil {
					count++

					// TODO: check shape exists

				} else {

					// TODO: check shape does not exist

				}
			}

			if tc.want {
				assert.NotZero(t, count)
				assert.GreaterOrEqual(t, count, tc.min)
				assert.LessOrEqual(t, count, tc.max)
			} else {
				assert.Zero(t, count)
			}
		})
	}
}

func TestNewDeduplicationAction(t *testing.T) {
	symbols := NewSymbolSet(
		NewSymbol(1, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(2, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(3, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(4, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(5, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(6, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(7, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(8, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(9, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(10, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(11, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(12, WithWeights(20, 20, 20, 20, 20)),
	)
	slots := NewSlots(Grid(5, 3), WithSymbols(symbols))

	weights1 := utils.AcquireWeighting().AddWeights(utils.Indexes{2, 3, 4, 5, 6}, []float64{85, 14.5, 0.4, 0.09, 0.01})
	defer weights1.Release()

	testCases := []struct {
		name    string
		weights utils.WeightedGenerator
		indexes utils.Indexes
		want    bool
		test2   bool
	}{
		{
			name:    "no dedup",
			weights: weights1,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 1, 5, 9},
			want:    false,
		},
		{
			name:    "dedup 3x1",
			weights: utils.AcquireWeighting().AddWeights(utils.Indexes{2}, []float64{100}),
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 8, 9, 10, 11, 12, 1, 5, 9},
			want:    true,
			test2:   true,
		},
		{
			name:    "dedup 3x1,3,9",
			weights: utils.AcquireWeighting().AddWeights(utils.Indexes{2}, []float64{100}),
			indexes: utils.Indexes{1, 2, 3, 4, 3, 6, 1, 8, 9, 3, 11, 9, 1, 5, 9},
			want:    true,
			test2:   true,
		},
		{
			name:    "dedup 7x6",
			weights: weights1,
			indexes: utils.Indexes{1, 6, 3, 6, 3, 6, 6, 8, 9, 1, 11, 6, 6, 5, 6},
			want:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prng := rng.NewRNG()
			defer prng.ReturnToPool()

			m := NewDeduplicationAction(tc.weights)
			require.NotNil(t, m)

			spin := AcquireSpin(slots, prng)
			defer spin.Release()

			copy(spin.indexes, tc.indexes)

			m2 := m.Triggered(spin)
			if tc.want {
				require.NotNil(t, m2)
			} else {
				require.Nil(t, m2)
			}

			if tc.test2 {
				m3 := m.Triggered(spin)
				require.Nil(t, m3)
			}
		})
	}
}

func TestNewDeduplicationAction_Fail(t *testing.T) {
	t.Run("dedup fail", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		slots := NewSlots(Grid(5, 5), WithSymbols(setF1))

		weights1 := utils.AcquireWeighting().AddWeights(utils.Indexes{1}, []float64{100})
		defer weights1.Release()

		m := NewDeduplicationAction(weights1)

		spin := AcquireSpin(slots, prng)
		defer spin.Release()

		defer func() {
			e := recover()
			require.NotNil(t, e)
		}()

		m2 := m.Triggered(spin)
		require.Nil(t, m2)
	})
}

func TestNewDedupSymbolAction(t *testing.T) {
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
		want    bool
		count   uint8
	}{
		{
			name:    "no hits, all reels",
			symbol:  1,
			reels:   nil,
			indexes: utils.Indexes{2, 3, 4, 2, 3, 4, 2, 3, 4, 2, 3, 4, 2, 3, 4},
		},
		{
			name:    "no hits, reels 2,3,4",
			symbol:  1,
			reels:   utils.UInt8s{2, 3, 4},
			indexes: utils.Indexes{2, 3, 4, 2, 3, 4, 2, 3, 4, 2, 3, 4, 2, 3, 4},
		},
		{
			name:    "no dupes, all reels",
			symbol:  1,
			reels:   nil,
			indexes: utils.Indexes{2, 3, 1, 2, 1, 4, 1, 3, 4, 1, 3, 4, 2, 1, 4},
		},
		{
			name:    "no dupes, reels 2,3,4",
			symbol:  1,
			reels:   utils.UInt8s{2, 3, 4},
			indexes: utils.Indexes{1, 1, 1, 2, 1, 4, 1, 3, 4, 2, 3, 1, 1, 1, 1},
		},
		{
			name:    "1 dupe, all reels",
			symbol:  1,
			reels:   nil,
			indexes: utils.Indexes{2, 3, 1, 2, 1, 4, 1, 1, 4, 1, 3, 4, 2, 1, 4},
			want:    true,
			count:   5,
		},
		{
			name:    "1 dupe, reels 2,3,4",
			symbol:  1,
			reels:   utils.UInt8s{2, 3, 4},
			indexes: utils.Indexes{1, 1, 1, 2, 1, 4, 1, 3, 4, 1, 3, 1, 1, 1, 1},
			want:    true,
			count:   9,
		},
		{
			name:    "3 dupes, all reels",
			symbol:  1,
			reels:   nil,
			indexes: utils.Indexes{1, 3, 1, 2, 1, 4, 1, 1, 4, 1, 3, 4, 2, 1, 1},
			want:    true,
			count:   5,
		},
		{
			name:    "3 dupes, reels 1,3,5",
			symbol:  1,
			reels:   utils.UInt8s{1, 3, 5},
			indexes: utils.Indexes{1, 1, 1, 2, 3, 4, 1, 3, 1, 1, 3, 1, 1, 1, 1},
			want:    true,
			count:   5,
		},
		{
			name:    "all dupes, all reels",
			symbol:  1,
			reels:   nil,
			indexes: utils.Indexes{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			want:    true,
			count:   5,
		},
		{
			name:    "all dupes, reels 1,3,5",
			symbol:  1,
			reels:   utils.UInt8s{1, 3, 5},
			indexes: utils.Indexes{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			want:    true,
			count:   9,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewDedupSymbolAction(tc.symbol, tc.reels)
			require.NotNil(t, a)

			copy(spin.indexes, tc.indexes)

			got := a.Triggered(spin)
			if tc.want {
				require.NotNil(t, got)
				assert.Equal(t, tc.count, spin.CountSymbol(tc.symbol))
			} else {
				require.Nil(t, got)
			}
		})
	}
}

func TestNewRemoveSymbolsAction(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name        string
		symbol      utils.Index
		reels       utils.UInt8s
		dedupe      bool
		replSymbols utils.Indexes
		replChances []float64
		replReels   utils.UInt8s
		indexes     utils.Indexes
		want        bool
		check       utils.UInt8s
	}{
		{
			name:        "no trigger (1)",
			symbol:      10,
			reels:       utils.UInt8s{1},
			dedupe:      false,
			replSymbols: utils.Indexes{6, 7, 8, 9},
			replChances: []float64{75, 75, 75, 75},
			replReels:   utils.UInt8s{2},
			indexes:     utils.Indexes{1, 2, 3, 7, 8, 9, 1, 2, 3, 7, 8, 9, 1, 2, 3},
		},
		{
			name:        "no trigger (2)",
			symbol:      10,
			reels:       utils.UInt8s{2},
			dedupe:      false,
			replSymbols: utils.Indexes{6, 7, 8, 9},
			replChances: []float64{75, 75, 75, 75},
			replReels:   utils.UInt8s{1},
			indexes:     utils.Indexes{7, 8, 9, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name:        "almost trigger (1)",
			symbol:      10,
			reels:       utils.UInt8s{1},
			dedupe:      false,
			replSymbols: utils.Indexes{6, 7, 8, 9},
			replChances: []float64{75, 75, 75, 75},
			replReels:   utils.UInt8s{2},
			indexes:     utils.Indexes{1, 2, 3, 10, 7, 9, 1, 2, 3, 7, 8, 9, 1, 2, 3},
		},
		{
			name:        "almost trigger (2)",
			symbol:      10,
			reels:       utils.UInt8s{2},
			dedupe:      false,
			replSymbols: utils.Indexes{6, 7, 8, 9},
			replChances: []float64{75, 75, 75, 75},
			replReels:   utils.UInt8s{1},
			indexes:     utils.Indexes{10, 7, 9, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name:        "trigger, no replace (1)",
			symbol:      10,
			reels:       utils.UInt8s{1},
			dedupe:      false,
			replSymbols: utils.Indexes{6, 7, 8, 9},
			replChances: []float64{75, 75, 75, 75},
			replReels:   utils.UInt8s{2},
			indexes:     utils.Indexes{10, 5, 6, 1, 2, 3, 7, 8, 9, 4, 5, 6, 1, 2, 3},
		},
		{
			name:        "trigger, no replace (2)",
			symbol:      10,
			reels:       utils.UInt8s{2},
			dedupe:      false,
			replSymbols: utils.Indexes{6, 7, 8, 9},
			replChances: []float64{75, 75, 75, 75},
			replReels:   utils.UInt8s{1},
			indexes:     utils.Indexes{1, 2, 3, 10, 6, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name:        "trigger, replace 1",
			symbol:      10,
			reels:       utils.UInt8s{1},
			dedupe:      false,
			replSymbols: utils.Indexes{6, 7, 8, 9},
			replChances: []float64{100, 75, 75, 75},
			replReels:   utils.UInt8s{2},
			indexes:     utils.Indexes{10, 5, 6, 1, 6, 3, 7, 8, 9, 4, 5, 6, 1, 2, 3},
			want:        true,
		},
		{
			name:        "trigger, replace 2",
			symbol:      10,
			reels:       utils.UInt8s{2},
			dedupe:      false,
			replSymbols: utils.Indexes{6, 7, 8, 9},
			replChances: []float64{75, 75, 100, 75},
			replReels:   utils.UInt8s{1},
			indexes:     utils.Indexes{1, 8, 6, 10, 6, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			want:        true,
		},
		{
			name:        "trigger, replace 3",
			symbol:      10,
			reels:       utils.UInt8s{2},
			dedupe:      false,
			replSymbols: utils.Indexes{6, 7, 8, 9},
			replChances: []float64{25, 75, 50, 100},
			replReels:   utils.UInt8s{1},
			indexes:     utils.Indexes{9, 8, 6, 10, 6, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			want:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewReplaceSymbolsAction(tc.symbol, tc.reels, tc.replSymbols, tc.replChances, tc.replReels)
			require.NotNil(t, a)

			assert.True(t, a.replaceSymbols)
			assert.True(t, a.genAllowDupes)
			assert.Equal(t, tc.symbol, a.symbol)
			assert.Equal(t, tc.reels, a.detectReels)
			assert.Equal(t, tc.replSymbols, a.replSymbols)
			assert.Equal(t, tc.replChances, a.replChances)
			assert.Equal(t, tc.replReels, a.replReels)

			if tc.dedupe {
				a.GenerateNoDupes()
				assert.False(t, a.genAllowDupes)
			}

			copy(spin.indexes, tc.indexes)

			got := a.Triggered(spin)
			if tc.want {
				assert.NotNil(t, got)
				assert.NotEqualValues(t, tc.indexes, spin.indexes)
			} else {
				assert.Nil(t, got)
				assert.EqualValues(t, tc.indexes, spin.indexes)
			}
		})
	}
}

func TestNewRemovePayoutsAction_Mech1(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name      string
		direction PayDirection
		remWilds  bool
		reelDupes bool
		indexes   utils.Indexes
		want      bool
	}{
		{
			name:      "ltr, none",
			direction: PayLTR,
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:      "ltr, one",
			direction: PayLTR,
			indexes:   utils.Indexes{1, 4, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			want:      true,
		},
		{
			name:      "ltr, three",
			direction: PayLTR,
			indexes:   utils.Indexes{4, 6, 5, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			want:      true,
		},
		{
			name:      "ltr, wild",
			direction: PayLTR,
			remWilds:  true,
			indexes:   utils.Indexes{1, 9, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			want:      true,
		},
		{
			name:      "ltr, wilds",
			direction: PayLTR,
			remWilds:  true,
			indexes:   utils.Indexes{9, 9, 9, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			want:      true,
		},
		{
			name:      "rtl, none",
			direction: PayRTL,
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:      "rtl, one",
			direction: PayRTL,
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 6, 2, 3},
			want:      true,
		},
		{
			name:      "rtl, three",
			direction: PayRTL,
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 4, 5, 6},
			want:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewPreventPayoutsAction(100.0, tc.direction, 1, tc.remWilds)
			require.NotNil(t, a)
			assert.True(t, a.preventPayouts)
			assert.Equal(t, tc.direction, a.prevPayoutDir)
			assert.Equal(t, 1, a.prevPayoutMech)

			if !tc.reelDupes {
				a.GenerateNoDupes()
			}

			copy(spin.indexes, tc.indexes)

			got := a.Triggered(spin)
			if tc.want {
				assert.NotNil(t, got)
				assert.NotEqualValues(t, tc.indexes, spin.indexes)
			} else {
				assert.Nil(t, got)
				assert.EqualValues(t, tc.indexes, spin.indexes)
			}
		})
	}
}

func TestNewGenerateBonusAction(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1), HotReelsAsBonusSymbol())

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name    string
		count   uint8
		chances []float64
		kind    SpinKind
		indexes utils.Indexes
		hot     []bool
		bonus   utils.Index
		want    bool
	}{
		{
			name:    "not a free spin",
			count:   5,
			chances: []float64{0, 100, 80, 70, 60, 50, 40, 30, 20},
			kind:    FirstSpin,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			hot:     []bool{false, false, false, false, false},
			bonus:   4,
		},
		{
			name:    "no bonus symbol (1)",
			count:   5,
			chances: []float64{0, 100, 80, 70, 60, 50, 40, 30, 20},
			kind:    FreeSpin,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			hot:     []bool{false, false, false, false, false},
			bonus:   0,
		},
		{
			name:    "no bonus symbol (2)",
			count:   5,
			chances: []float64{0, 100, 80, 70, 60, 50, 40, 30, 20},
			kind:    FreeSpin,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			hot:     []bool{false, false, false, false, false},
			bonus:   math.MaxUint16,
		},
		{
			name:    "already enough",
			count:   3,
			chances: []float64{0, 100, 80, 70, 60, 50, 40, 30, 20},
			kind:    FreeSpin,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 4, 2, 3, 4, 5, 6, 1, 2, 3},
			hot:     []bool{false, false, false, false, false},
			bonus:   4,
		},
		{
			name:    "generate - no hot (1)",
			count:   4,
			chances: []float64{0, 100, 80, 70, 60, 50, 40, 30, 20},
			kind:    FreeSpin,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			hot:     []bool{false, false, false, false, false},
			bonus:   4,
			want:    true,
		},
		{
			name:    "generate - no hot (2)",
			count:   5,
			chances: []float64{0, 100, 80, 70, 60, 50, 40, 30, 20},
			kind:    FreeSpin,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			hot:     []bool{false, false, false, false, false},
			bonus:   6,
			want:    true,
		},
		{
			name:    "generate - hot (1)",
			count:   4,
			chances: []float64{0, 100, 80, 70, 60, 50, 40, 30, 20},
			kind:    FreeSpin,
			indexes: utils.Indexes{1, 2, 3, 2, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 4},
			hot:     []bool{false, true, false, false, false},
			bonus:   4,
			want:    true,
		},
		{
			name:    "generate - hot (2)",
			count:   5,
			chances: []float64{0, 100, 80, 70, 60, 50, 40, 30, 20},
			kind:    FreeSpin,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 3, 5, 6, 4, 2, 3},
			hot:     []bool{false, false, true, true, false},
			bonus:   6,
			want:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewGenerateBonusAction(tc.count, tc.chances)
			require.NotNil(t, a)

			assert.True(t, a.generateBonus)
			assert.Equal(t, tc.count, a.bonusCount)
			assert.Equal(t, tc.chances, a.bonusChances)

			spin.kind = tc.kind
			spin.bonusSymbol = tc.bonus
			copy(spin.indexes, tc.indexes)
			copy(spin.hot, tc.hot)

			if tc.want {
				var a2 SpinActioner
				for ix := 0; ix < 25; ix++ {
					if a2 = a.Triggered(spin); a2 != nil {
						break
					}
				}
				require.NotNil(t, a2)
				assert.GreaterOrEqual(t, spin.CountBonusSymbol(), tc.count)
			} else {
				a2 := a.Triggered(spin)
				require.Nil(t, a2)
			}
		})
	}
}

func TestMath(t *testing.T) {
	testCases := []struct {
		offs int
		rows int
		want int
	}{
		{offs: 0, rows: 4, want: 0},
		{offs: 1, rows: 4, want: 0},
		{offs: 2, rows: 4, want: 0},
		{offs: 3, rows: 4, want: 0},
		{offs: 4, rows: 4, want: 4},
		{offs: 5, rows: 4, want: 4},
		{offs: 6, rows: 4, want: 4},
		{offs: 7, rows: 4, want: 4},
		{offs: 8, rows: 4, want: 8},
		{offs: 9, rows: 4, want: 8},
		{offs: 10, rows: 4, want: 8},
		{offs: 11, rows: 4, want: 8},
		{offs: 12, rows: 4, want: 12},
		{offs: 13, rows: 4, want: 12},
		{offs: 14, rows: 4, want: 12},
		{offs: 15, rows: 4, want: 12},
		{offs: 16, rows: 4, want: 16},
		{offs: 17, rows: 4, want: 16},
		{offs: 18, rows: 4, want: 16},
		{offs: 19, rows: 4, want: 16},
		{offs: 20, rows: 4, want: 20},
		{offs: 21, rows: 4, want: 20},
		{offs: 22, rows: 4, want: 20},
		{offs: 23, rows: 4, want: 20},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%d", tc.offs, tc.rows), func(t *testing.T) {
			got := (tc.offs / tc.rows) * tc.rows
			assert.Equal(t, tc.want, got)
		})
	}
}

func BenchmarkGenerateSymbolAction80pct(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()
	spin.SetKind(FreeSpin)

	a := NewGenerateSymbolAction(13, []float64{80}).WithSpinKinds([]SpinKind{FreeSpin})

	indexes := utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}

	for i := 0; i < b.N; i++ {
		copy(spin.indexes, indexes)
		a.Triggered(spin)
	}
}

func BenchmarkGenerateSymbolAction100pct(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()
	spin.SetKind(FreeSpin)

	a := NewGenerateSymbolAction(13, []float64{100}).WithSpinKinds([]SpinKind{FreeSpin})

	indexes := utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}

	for i := 0; i < b.N; i++ {
		copy(spin.indexes, indexes)
		a.Triggered(spin)
	}
}

func BenchmarkGenerateShapeAction10pct(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()
	spin.SetKind(FreeSpin)

	grid := GridOffsets{{-1, -1}, {-1, 1}, {0, 0}, {1, -1}, {1, 1}}
	centers := GridOffsets{{1, 1}, {2, 1}, {3, 1}}

	weights := utils.AcquireWeighting().AddWeights(utils.Indexes{4, 5, 6}, []float64{1, 1, 1})
	defer weights.Release()

	a := NewGenerateShapeAction(10, grid, centers, weights).WithSpinKinds([]SpinKind{SecondSpin, FreeSpin})

	indexes := utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}

	for i := 0; i < b.N; i++ {
		copy(spin.indexes, indexes)
		a.Triggered(spin)
	}
}

func BenchmarkGenerateShapeAction100pct(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	symbols := NewSymbolSet(
		NewSymbol(1, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(2, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(3, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(4, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(5, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(6, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(7, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(8, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(9, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(10, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(11, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(12, WithWeights(20, 20, 20, 20, 20)),
	)
	slots := NewSlots(Grid(5, 3), WithSymbols(symbols))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()
	spin.SetKind(FreeSpin)

	grid := GridOffsets{{-1, -1}, {-1, 1}, {0, 0}, {1, -1}, {1, 1}}
	centers := GridOffsets{{1, 1}, {2, 1}, {3, 1}}

	weights := utils.AcquireWeighting().AddWeights(utils.Indexes{4, 5, 6}, []float64{1, 1, 1})
	defer weights.Release()

	a := NewGenerateShapeAction(100, grid, centers, weights).WithSpinKinds([]SpinKind{SecondSpin, FreeSpin})

	indexes := utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}

	for i := 0; i < b.N; i++ {
		copy(spin.indexes, indexes)
		a.Triggered(spin)
	}
}

func BenchmarkDeduplicationAction85(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	symbols := NewSymbolSet(
		NewSymbol(1, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(2, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(3, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(4, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(5, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(6, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(7, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(8, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(9, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(10, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(11, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(12, WithWeights(20, 20, 20, 20, 20)),
	)
	slots := NewSlots(Grid(5, 3), WithSymbols(symbols))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()
	spin.SetKind(FirstSpin)

	weights := utils.AcquireWeighting().AddWeights(utils.Indexes{2, 3, 4}, []float64{85, 14.5, 0.5})
	defer weights.Release()

	a := NewDeduplicationAction(weights).WithSpinKinds([]SpinKind{FirstSpin})

	indexes := utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}

	for i := 0; i < b.N; i++ {
		copy(spin.indexes, indexes)
		a.Triggered(spin)
	}
}

func BenchmarkDeduplicationAction100(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	symbols := NewSymbolSet(
		NewSymbol(1, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(2, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(3, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(4, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(5, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(6, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(7, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(8, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(9, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(10, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(11, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(12, WithWeights(20, 20, 20, 20, 20)),
	)
	slots := NewSlots(Grid(5, 3), WithSymbols(symbols))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()
	spin.SetKind(FirstSpin)

	weights := utils.AcquireWeighting().AddWeights(utils.Indexes{2}, []float64{100})
	defer weights.Release()

	a := NewDeduplicationAction(weights).WithSpinKinds([]SpinKind{FirstSpin})

	indexes := utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}

	for i := 0; i < b.N; i++ {
		copy(spin.indexes, indexes)
		a.Triggered(spin)
	}
}

func BenchmarkDeduplicationAction15(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	symbols := NewSymbolSet(
		NewSymbol(1, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(2, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(3, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(4, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(5, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(6, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(7, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(8, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(9, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(10, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(11, WithWeights(20, 20, 20, 20, 20)),
		NewSymbol(12, WithWeights(20, 20, 20, 20, 20)),
	)
	slots := NewSlots(Grid(5, 3), WithSymbols(symbols))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()
	spin.SetKind(FirstSpin)

	weights := utils.AcquireWeighting().AddWeights(utils.Indexes{3}, []float64{15})
	defer weights.Release()

	a := NewDeduplicationAction(weights).WithSpinKinds([]SpinKind{FirstSpin})

	indexes := utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}

	for i := 0; i < b.N; i++ {
		copy(spin.indexes, indexes)
		a.Triggered(spin)
	}
}
