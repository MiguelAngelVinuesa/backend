package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewScatterAction(t *testing.T) {
	alt := NewScatterPayoutAction(8, 2, 6)

	wheelWeights := utils.AcquireWeighting().AddWeights(utils.Indexes{1, 2, 3}, []float64{65, 25, 10})

	testCases := []struct {
		name          string
		kind          SpinActionResult
		symbol        utils.Index
		nrOfSpins     uint8
		scatterCount  uint8
		scatterPayout float64
		bonusSymbol   bool
		multi         utils.Indexes
	}{
		{
			name:          "payout",
			kind:          Payout,
			symbol:        5,
			scatterCount:  3,
			scatterPayout: 7.5,
		},
		{
			name:         "1 scatter",
			kind:         FreeSpins,
			symbol:       8,
			nrOfSpins:    1,
			scatterCount: 1,
		},
		{
			name:         "3 scatters",
			kind:         FreeSpins,
			symbol:       9,
			nrOfSpins:    10,
			scatterCount: 3,
			bonusSymbol:  true,
		},
		{
			name:         "3 scatters / multi symbol",
			kind:         FreeSpins,
			symbol:       9,
			nrOfSpins:    10,
			scatterCount: 3,
			bonusSymbol:  true,
			multi:        utils.Indexes{10},
		},
		{
			name:         "bonus game",
			kind:         BonusGame,
			symbol:       4,
			scatterCount: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var a *ScatterAction
			switch tc.kind {
			case Payout:
				a = NewScatterPayoutAction(tc.symbol, tc.scatterCount, tc.scatterPayout)
			case FreeSpins:
				a = NewScatterFreeSpinsAction(tc.nrOfSpins, false, tc.symbol, tc.scatterCount, tc.bonusSymbol)
			case BonusGame:
				a = NewScatterBonusWheelAction(tc.symbol, 3, 1, wheelWeights)
			}
			require.NotNil(t, a)

			assert.Equal(t, tc.kind, a.result)
			assert.Equal(t, tc.nrOfSpins, a.nrOfSpins)
			assert.Equal(t, tc.symbol, a.symbol)
			assert.Equal(t, tc.scatterCount, a.scatterCount)
			assert.Equal(t, tc.scatterPayout, a.scatterPayout)
			assert.Equal(t, tc.bonusSymbol, a.bonusSymbol)

			if tc.multi != nil {
				a2 := a.WithMultiSymbols(tc.multi...)
				require.Equal(t, a, a2)
				assert.Equal(t, tc.multi, a.multiSymbols)
			}

			n := a.WithAlternate(alt)
			assert.Equal(t, a, n)
			assert.Equal(t, alt, n.alternate)
		})
	}
}

func TestNewAllScatterAction(t *testing.T) {
	t.Run("new all scatter action", func(t *testing.T) {
		a := NewAllScatterAction(3)
		require.NotNil(t, a)

		assert.Equal(t, ExtraPayouts, a.stage)
		assert.Equal(t, Payout, a.result)
		assert.True(t, a.allScatters)
		assert.Equal(t, uint8(3), a.allMinimum)
	})
}

func TestNewBonusScatterAction(t *testing.T) {
	t.Run("new bonus scatter action", func(t *testing.T) {
		a := NewBonusScatterAction(10)
		require.NotNil(t, a)

		assert.Equal(t, AwardBonuses, a.stage)
		assert.Equal(t, Payout, a.result)
		assert.True(t, a.bonusScatter)
		assert.Equal(t, uint8(10), a.bonusLines)
	})
}

func TestScatterAction_WithMultiSymbolsFail(t *testing.T) {
	t.Run("new scatter with multi symbols fail", func(t *testing.T) {
		defer func() {
			e := recover()
			require.NotNil(t, e)
		}()

		a := NewBonusScatterAction(10).WithMultiSymbols(11)
		require.Nil(t, a)
	})
}

func TestScatterAction_Triggered(t *testing.T) {
	a1 := NewScatterFreeSpinsAction(10, false, 8, 5, false)
	a2 := NewScatterFreeSpinsAction(3, false, 8, 3, false)
	a3 := NewScatterPayoutAction(8, 2, 6)
	a4 := NewScatterFreeSpinsAction(10, true, 10, 3, true).WithMultiSymbols(11)
	a5 := NewScatterFreeSpinsAction(10, true, 10, 3, true).WithMultiSymbols(11, 12)

	a1.WithAlternate(a2)
	a2.WithAlternate(a3)

	testCases := []struct {
		name    string
		scatter *ScatterAction
		spin    *Spin
		want    SpinActioner
	}{
		{"a1 - no hit", a1, &Spin{reelCount: 3, indexes: utils.Indexes{1, 2, 3, 1, 8, 3, 1, 2, 3}, newScatters: 1}, nil},
		{"a2 - no hit", a2, &Spin{reelCount: 3, indexes: utils.Indexes{1, 2, 8, 1, 2, 3, 1, 2, 3}, newScatters: 1}, nil},
		{"a3 - no hit", a3, &Spin{reelCount: 3, indexes: utils.Indexes{1, 2, 3, 1, 2, 3, 8, 2, 3}, newScatters: 1}, nil},
		{"a4 - no hit", a4, &Spin{reelCount: 3, indexes: utils.Indexes{1, 2, 3, 10, 2, 3, 8, 11, 3}, newScatters: 2}, nil},
		{"a5 - no hit", a5, &Spin{reelCount: 3, indexes: utils.Indexes{1, 2, 3, 10, 2, 3, 8, 11, 3}, newScatters: 2}, nil},
		{"a1 - hit", a1, &Spin{reelCount: 5, indexes: utils.Indexes{1, 8, 2, 3, 8, 1, 8, 1, 8, 3, 1, 8, 3, 2, 8}, newScatters: 6}, a1},
		{"a2 - hit", a2, &Spin{reelCount: 3, indexes: utils.Indexes{1, 8, 3, 8, 2, 3, 1, 8, 3}, newScatters: 3}, a2},
		{"a3 - hit", a3, &Spin{reelCount: 3, indexes: utils.Indexes{8, 2, 3, 1, 2, 3, 8, 2, 3}, newScatters: 2}, a3},
		{"a1 - alt a2", a1, &Spin{reelCount: 3, indexes: utils.Indexes{1, 8, 3, 1, 8, 8, 1, 2, 3}, newScatters: 3}, a2},
		{"a1 - alt a3", a1, &Spin{reelCount: 3, indexes: utils.Indexes{1, 8, 3, 1, 2, 3, 1, 2, 8}, newScatters: 2}, a3},
		{"a2 - alt a3", a2, &Spin{reelCount: 3, indexes: utils.Indexes{1, 2, 8, 1, 2, 8, 1, 2, 3}, newScatters: 2}, a3},
		{"a4 - hit 3x10", a4, &Spin{reelCount: 3, indexes: utils.Indexes{1, 10, 3, 10, 2, 3, 8, 10, 3}, newScatters: 3}, a4},
		{"a4 - hit 3x11", a4, &Spin{reelCount: 3, indexes: utils.Indexes{11, 1, 3, 1, 11, 3, 8, 11, 3}, newScatters: 3}, a4},
		{"a4 - hit 2x10 + 2x11", a4, &Spin{reelCount: 3, indexes: utils.Indexes{11, 10, 3, 11, 2, 3, 8, 10, 3}, newScatters: 4}, a4},
		{"a4 - hit 2x10 + 3x11", a4, &Spin{reelCount: 3, indexes: utils.Indexes{11, 10, 3, 10, 2, 3, 8, 11, 11}, newScatters: 5}, a4},
		{"a5 - hit 1x10 + 1x11 + 1x12", a5, &Spin{reelCount: 3, indexes: utils.Indexes{1, 10, 3, 12, 2, 3, 8, 11, 4}, newScatters: 3}, a5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.scatter.Triggered(tc.spin)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestScatterAction_Payout(t *testing.T) {
	testCases := []struct {
		name    string
		indexes utils.Indexes
		payout  results.Payout
	}{
		{
			name:    "no scatters",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "1 scatter",
			indexes: utils.Indexes{1, 2, 14, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "2 scatters",
			indexes: utils.Indexes{1, 2, 14, 4, 5, 6, 1, 2, 3, 4, 14, 6, 1, 2, 3},
			payout:  ScatterSymbolPayoutWithMap(1.5, 0, 14, 2, utils.UInt8s{2, 10}),
		},
		{
			name:    "3 scatters",
			indexes: utils.Indexes{1, 2, 14, 4, 5, 14, 1, 2, 3, 4, 14, 6, 1, 2, 3},
			payout:  ScatterSymbolPayoutWithMap(3, 0, 14, 3, utils.UInt8s{2, 5, 10}),
		},
		{
			name:    "4 scatters",
			indexes: utils.Indexes{1, 2, 14, 4, 5, 6, 14, 2, 3, 4, 14, 6, 1, 2, 14},
			payout:  ScatterSymbolPayoutWithMap(10, 0, 14, 4, utils.UInt8s{2, 6, 10, 14}),
		},
		{
			name:    "5 scatters",
			indexes: utils.Indexes{1, 2, 14, 4, 5, 6, 1, 14, 3, 4, 14, 6, 1, 14, 14},
			payout:  ScatterSymbolPayoutWithMap(25, 0, 14, 5, utils.UInt8s{2, 7, 10, 13, 14}),
		},
		{
			name:    "7 scatters",
			indexes: utils.Indexes{1, 2, 14, 14, 5, 14, 1, 14, 3, 4, 14, 6, 1, 14, 14},
			payout:  ScatterSymbolPayoutWithMap(25, 0, 14, 7, utils.UInt8s{2, 3, 5, 7, 10, 13, 14}),
		},
	}

	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	a1 := NewScatterPayoutAction(14, 2, 1.5)
	a2 := NewScatterPayoutAction(14, 3, 3).WithAlternate(a1)
	a3 := NewScatterPayoutAction(14, 4, 10).WithAlternate(a2)
	a4 := NewScatterPayoutAction(14, 5, 25).WithAlternate(a3)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spin := AcquireSpin(slots, prng)
			defer spin.Release()

			spin.indexes = tc.indexes
			spin.CountSpecials()

			res := &results.Result{}

			a := a4.Payout(spin, res)

			if tc.payout == nil {
				assert.Nil(t, a)
			} else {
				require.NotNil(t, a)
				assert.Equal(t, tc.payout.(*SpinPayout).Total(), res.Total)

				require.Equal(t, 1, len(res.Payouts))

				p := res.Payouts[0]
				require.NotNil(t, p)

				if !tc.payout.(*SpinPayout).Equals(p.(*SpinPayout)) {
					assert.EqualValues(t, tc.payout, p)
				}
			}
		})
	}
}

func TestNewAllScatterAction_Triggered(t *testing.T) {
	testCases := []struct {
		name    string
		indexes utils.Indexes
		minimum uint8
		want    bool
	}{
		{"none (1)", utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 3, false},
		{"none (2)", utils.Indexes{1, 2, 3, 1, 2, 3, 4, 5, 6, 7, 8, 9, 7, 8, 9}, 3, false},
		{"none (3)", utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 7, 8, 9, 7, 8, 9}, 4, false},
		{"none (4)", utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 7, 8, 9}, 5, false},
		{"none (5)", utils.Indexes{7, 8, 9, 7, 8, 9, 1, 1, 1, 7, 8, 9, 7, 8, 9}, 5, false},
		{"3x 1", utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 1, 9}, 3, true},
		{"4x 1", utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 1, 6, 7, 1, 9}, 3, true},
		{"8x 1", utils.Indexes{1, 2, 3, 1, 1, 6, 1, 2, 1, 4, 1, 1, 7, 1, 9}, 4, true},
		{"3x 3", utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 8, 3}, 3, true},
		{"4x 3", utils.Indexes{1, 2, 3, 3, 5, 6, 1, 2, 3, 4, 5, 6, 7, 8, 3}, 3, true},
		{"8x 3", utils.Indexes{3, 2, 3, 3, 5, 6, 3, 3, 3, 4, 5, 6, 3, 8, 3}, 5, true},
		{"3x 9", utils.Indexes{9, 9, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 3, true},
		{"4x 9", utils.Indexes{9, 9, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 9, 9}, 3, true},
		{"8x 9", utils.Indexes{9, 9, 3, 4, 5, 6, 9, 9, 9, 4, 9, 6, 7, 9, 9}, 6, true},
	}

	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spin := AcquireSpin(slots, prng)
			defer spin.Release()
			spin.indexes = tc.indexes

			a := NewAllScatterAction(tc.minimum)

			got := a.Triggered(spin) == a
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestNewAllScatterAction_Payout(t *testing.T) {
	testCases := []struct {
		name    string
		indexes utils.Indexes
		minimum uint8
		payout  results.Payout
		payout2 results.Payout
	}{
		{
			name:    "none (1)",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			minimum: 3,
		},
		{
			name:    "none (2)",
			indexes: utils.Indexes{1, 2, 3, 1, 2, 3, 4, 5, 6, 7, 8, 9, 7, 8, 9},
			minimum: 3,
		},
		{
			name:    "none (3)",
			indexes: utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 7, 8, 9, 7, 8, 9},
			minimum: 4,
		},
		{
			name:    "none (4)",
			indexes: utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 7, 8, 9},
			minimum: 5,
		},
		{
			name:    "none (5)",
			indexes: utils.Indexes{7, 8, 9, 7, 8, 9, 1, 1, 1, 7, 8, 9, 7, 8, 9},
			minimum: 5,
		},
		{
			name:    "3x 1",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 1, 9},
			minimum: 3,
			payout:  ScatterSymbolPayoutWithMap(0.5, 0, 1, 3, utils.UInt8s{0, 6, 13}),
		},
		{
			name:    "4x 1",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 1, 6, 7, 1, 9},
			minimum: 3,
			payout:  ScatterSymbolPayoutWithMap(2, 0, 1, 4, utils.UInt8s{0, 6, 10, 13}),
		},
		{
			name:    "8x 1",
			indexes: utils.Indexes{1, 2, 3, 1, 1, 6, 1, 2, 1, 4, 1, 1, 7, 1, 9},
			minimum: 4,
			payout:  ScatterSymbolPayoutWithMap(4, 0, 1, 8, utils.UInt8s{0, 3, 4, 6, 8, 10, 11, 13}),
		},
		{
			name:    "3x 3",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 8, 3},
			minimum: 3,
			payout:  ScatterSymbolPayoutWithMap(1, 0, 3, 3, utils.UInt8s{2, 8, 14}),
		},
		{
			name:    "4x 3",
			indexes: utils.Indexes{1, 2, 3, 3, 5, 6, 1, 2, 3, 4, 5, 6, 7, 8, 3},
			minimum: 3,
			payout:  ScatterSymbolPayoutWithMap(2.5, 0, 3, 4, utils.UInt8s{2, 3, 8, 14}),
		},
		{
			name:    "8x 3",
			indexes: utils.Indexes{3, 2, 3, 3, 5, 6, 3, 3, 3, 4, 5, 6, 3, 8, 3},
			minimum: 5,
			payout:  ScatterSymbolPayoutWithMap(5, 0, 3, 8, utils.UInt8s{0, 2, 3, 6, 7, 8, 12, 14}),
		},
		{
			name:    "3x 9",
			indexes: utils.Indexes{9, 9, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			minimum: 3,
			payout:  ScatterSymbolPayoutWithMap(1, 0, 9, 3, utils.UInt8s{0, 1, 14}),
		},
		{
			name:    "4x 9",
			indexes: utils.Indexes{9, 9, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 9, 9},
			minimum: 3,
			payout:  ScatterSymbolPayoutWithMap(4, 0, 9, 4, utils.UInt8s{0, 1, 13, 14}),
		},
		{
			name:    "8x 9",
			indexes: utils.Indexes{9, 9, 3, 4, 5, 6, 9, 9, 9, 4, 9, 6, 7, 9, 9},
			minimum: 6,
			payout:  ScatterSymbolPayoutWithMap(10, 0, 9, 8, utils.UInt8s{0, 1, 6, 7, 8, 10, 13, 14}),
		},
		{
			name:    "3x 1 + 3x 4",
			indexes: utils.Indexes{1, 2, 4, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 1, 9},
			minimum: 3,
			payout:  ScatterSymbolPayoutWithMap(0.5, 0, 1, 3, utils.UInt8s{0, 6, 13}),
			payout2: ScatterSymbolPayoutWithMap(1, 0, 4, 3, utils.UInt8s{2, 3, 9}),
		},
	}

	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spin := AcquireSpin(slots, prng)
			defer spin.Release()
			spin.indexes = tc.indexes

			a := NewAllScatterAction(tc.minimum)

			res := &results.Result{}
			a.Payout(spin, res)

			if tc.payout == nil {
				assert.Zero(t, res.Total)
				assert.Zero(t, len(res.Payouts))
			} else {
				assert.GreaterOrEqual(t, res.Total, tc.payout.(*SpinPayout).Total())
				require.GreaterOrEqual(t, len(res.Payouts), 1)

				p := res.Payouts[0]
				require.NotNil(t, p)

				if !tc.payout.(*SpinPayout).Equals(p.(*SpinPayout)) {
					assert.EqualValues(t, tc.payout, p)
				}

				if tc.payout2 != nil {
					p = res.Payouts[1]
					require.NotNil(t, p)

					if !tc.payout2.(*SpinPayout).Equals(p.(*SpinPayout)) {
						assert.EqualValues(t, tc.payout2, p)
					}
				}
			}
		})
	}
}

func TestNewBonusScatterAction_Triggered(t *testing.T) {
	testCases := []struct {
		name    string
		indexes utils.Indexes
		symbol  utils.Index
		hot     []bool
		want    bool
	}{
		{"no - no hot", utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 7, 4, 5, 6, 7, 8, 9}, 3, []bool{false, false, false, false, false}, false},
		{"no same reel - no hot", utils.Indexes{1, 2, 3, 6, 6, 6, 1, 2, 7, 4, 5, 9, 7, 8, 9}, 6, []bool{false, false, false, false, false}, false},
		{"no - 1 hot", utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 7, 4, 6, 9, 7, 8, 9}, 5, []bool{false, true, false, false, false}, false},
		{"no same reel - 1 hot", utils.Indexes{1, 2, 3, 10, 10, 10, 1, 2, 7, 4, 5, 9, 7, 8, 9}, 10, []bool{false, true, false, false, false}, false},
		{"yes 3x - no hot", utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 1, 9}, 1, []bool{false, false, false, false, false}, true},
		{"yes 4x - no hot", utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 3, 5, 6, 3, 1, 9}, 3, []bool{false, false, false, false, false}, true},
		{"yes 5x - no hot", utils.Indexes{6, 6, 3, 4, 5, 6, 1, 6, 3, 4, 5, 6, 7, 6, 6}, 6, []bool{false, false, false, false, false}, true},
		{"yes 3x - 1 hot", utils.Indexes{1, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 6, 7, 1, 9}, 1, []bool{false, false, true, false, false}, true},
		{"yes 4x - 1 hot", utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 3, 1, 9}, 3, []bool{false, false, false, true, false}, true},
		{"yes 5x - 1 hot", utils.Indexes{6, 6, 3, 4, 5, 6, 1, 6, 3, 4, 5, 6, 7, 9, 9}, 6, []bool{false, false, false, false, true}, true},
		{"yes 3x - 2 hot", utils.Indexes{1, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 6, 7, 1, 9}, 1, []bool{true, false, true, false, false}, true},
		{"yes 4x - 2 hot", utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 3, 1, 9}, 3, []bool{false, false, true, true, false}, true},
		{"yes 5x - 2 hot", utils.Indexes{6, 6, 3, 4, 5, 6, 1, 6, 3, 4, 5, 6, 7, 9, 9}, 6, []bool{false, true, false, false, true}, true},
	}

	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1), HotReelsAsBonusSymbol())

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spin := AcquireSpin(slots, prng)
			defer spin.Release()
			spin.indexes = tc.indexes
			spin.bonusSymbol = tc.symbol
			spin.hot = tc.hot

			a := NewBonusScatterAction(10)

			got := a.Triggered(spin) == a
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestBonusScatterAction_Payout(t *testing.T) {
	testCases := []struct {
		name     string
		indexes  utils.Indexes
		symbol   utils.Index
		hot      []bool
		payout   results.Payout
		expanded utils.Indexes
	}{
		{
			name:    "no - no hot",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 7, 4, 5, 6, 7, 8, 9},
			symbol:  3,
			hot:     []bool{false, false, false, false, false},
		},
		{
			name:    "no same reel - no hot",
			indexes: utils.Indexes{1, 2, 3, 6, 6, 6, 1, 2, 7, 4, 5, 9, 7, 8, 9},
			symbol:  6,
			hot:     []bool{false, false, false, false, false},
		},
		{
			name:    "no - 1 hot",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 7, 4, 6, 9, 7, 8, 9},
			symbol:  5,
			hot:     []bool{false, true, false, false, false},
		},
		{
			name:    "no same reel - 1 hot",
			indexes: utils.Indexes{1, 2, 3, 10, 10, 10, 1, 2, 7, 4, 5, 9, 7, 8, 9},
			symbol:  10,
			hot:     []bool{false, true, false, false, false},
		},
		{
			name:     "yes 3x - no hot",
			indexes:  utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 1, 9},
			symbol:   1,
			hot:      []bool{false, false, false, false, false},
			payout:   BonusSymbolPayout(5, 1, 1, 3),
			expanded: utils.Indexes{1, 1, 1, 4, 5, 6, 1, 1, 1, 4, 5, 6, 1, 1, 1},
		},
		{
			name:     "yes 2x - no hot",
			indexes:  utils.Indexes{1, 2, 3, 4, 5, 6, 11, 2, 3, 4, 5, 6, 11, 10, 9},
			symbol:   11,
			hot:      []bool{false, false, false, false, false},
			payout:   BonusSymbolPayout(10, 1, 11, 2),
			expanded: utils.Indexes{1, 2, 3, 4, 5, 6, 11, 11, 11, 4, 5, 6, 11, 11, 11},
		},
		{
			name:     "yes 4x - no hot",
			indexes:  utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 3, 5, 6, 3, 1, 9},
			symbol:   3,
			hot:      []bool{false, false, false, false, false},
			payout:   BonusSymbolPayout(25, 1, 3, 4),
			expanded: utils.Indexes{3, 3, 3, 4, 5, 6, 3, 3, 3, 3, 3, 3, 3, 3, 3},
		},
		{
			name:     "yes 5x - no hot",
			indexes:  utils.Indexes{6, 6, 3, 4, 5, 6, 1, 6, 3, 4, 5, 6, 7, 6, 6},
			symbol:   6,
			hot:      []bool{false, false, false, false, false},
			payout:   BonusSymbolPayout(70, 1, 6, 5),
			expanded: utils.Indexes{6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6},
		},
		{
			name:     "yes 3x - 1 hot",
			indexes:  utils.Indexes{1, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 6, 7, 1, 9},
			symbol:   1,
			hot:      []bool{false, false, true, false, false},
			payout:   BonusSymbolPayout(5, 1, 1, 3),
			expanded: utils.Indexes{1, 1, 1, 4, 5, 6, 1, 1, 1, 4, 5, 6, 1, 1, 1},
		},
		{
			name:     "yes 2x - 1 hot",
			indexes:  utils.Indexes{1, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 6, 7, 1, 11},
			symbol:   11,
			hot:      []bool{false, false, true, false, false},
			payout:   BonusSymbolPayout(10, 1, 11, 2),
			expanded: utils.Indexes{1, 2, 3, 4, 5, 6, 11, 11, 11, 4, 5, 6, 11, 11, 11},
		},
		{
			name:     "yes 4x - 1 hot",
			indexes:  utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 3, 1, 9},
			symbol:   3,
			hot:      []bool{false, false, false, true, false},
			payout:   BonusSymbolPayout(25, 1, 3, 4),
			expanded: utils.Indexes{3, 3, 3, 4, 5, 6, 3, 3, 3, 3, 3, 3, 3, 3, 3},
		},
		{
			name:     "yes 5x - 1 hot",
			indexes:  utils.Indexes{6, 6, 3, 4, 5, 6, 1, 6, 3, 4, 5, 6, 7, 9, 9},
			symbol:   6,
			hot:      []bool{false, false, false, false, true},
			payout:   BonusSymbolPayout(70, 1, 6, 5),
			expanded: utils.Indexes{6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6},
		},
		{
			name:     "yes 2x - 2 hot",
			indexes:  utils.Indexes{1, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 6, 7, 1, 9},
			symbol:   11,
			hot:      []bool{true, false, true, false, false},
			payout:   BonusSymbolPayout(10, 1, 11, 2),
			expanded: utils.Indexes{11, 11, 11, 4, 5, 6, 11, 11, 11, 4, 5, 6, 7, 1, 9},
		},
		{
			name:     "yes 3x - 2 hot",
			indexes:  utils.Indexes{1, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 6, 7, 1, 9},
			symbol:   1,
			hot:      []bool{true, false, true, false, false},
			payout:   BonusSymbolPayout(5, 1, 1, 3),
			expanded: utils.Indexes{1, 1, 1, 4, 5, 6, 1, 1, 1, 4, 5, 6, 1, 1, 1},
		},
		{
			name:     "yes 4x - 2 hot",
			indexes:  utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 3, 1, 9},
			symbol:   3,
			hot:      []bool{false, false, true, true, false},
			payout:   BonusSymbolPayout(25, 1, 3, 4),
			expanded: utils.Indexes{3, 3, 3, 4, 5, 6, 3, 3, 3, 3, 3, 3, 3, 3, 3},
		},
		{
			name:     "yes 5x - 2 hot",
			indexes:  utils.Indexes{6, 6, 3, 4, 5, 6, 1, 6, 3, 4, 5, 6, 7, 9, 9},
			symbol:   6,
			hot:      []bool{false, true, false, false, true},
			payout:   BonusSymbolPayout(70, 1, 6, 5),
			expanded: utils.Indexes{6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6},
		},
	}

	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1), HotReelsAsBonusSymbol())

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spin := AcquireSpin(slots, prng)
			defer spin.Release()
			copy(spin.indexes, tc.indexes)
			spin.bonusSymbol = tc.symbol
			spin.hot = tc.hot

			a := NewBonusScatterAction(10)

			res := &results.Result{}
			a.Payout(spin, res)

			if tc.payout == nil {
				assert.Zero(t, res.Total)
				assert.Zero(t, len(res.Payouts))
			} else {
				assert.Equal(t, tc.payout.(*SpinPayout).Total(), res.Total)
				require.Equal(t, 1, len(res.Payouts))
			}

			if tc.expanded == nil {
				assert.Equal(t, tc.indexes, spin.indexes)
			} else {
				assert.NotEqual(t, tc.indexes, spin.indexes)
				assert.Equal(t, tc.expanded, spin.indexes)
			}
		})
	}
}

func TestBonusScatterAction_PayoutFail(t *testing.T) {
	t.Run("scatter action payout fail", func(t *testing.T) {
		defer func() {
			e := recover()
			require.NotNilf(t, e, "should have failed")
		}()

		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		slots := NewSlots(Grid(5, 3), WithSymbols(setF1), HotReelsAsBonusSymbol())

		spin := AcquireSpin(slots, prng)
		defer spin.Release()

		spin.indexes = utils.Indexes{1, 0, 0, 2, 0, 0, 3, 0, 0, 4, 0, 0, 5, 0, 0}
		spin.SetBonusSymbol(0, false)

		a := NewBonusScatterAction(10)

		a.Payout(spin, nil)
		require.Nil(t, a, "should have failed")
	})
}
