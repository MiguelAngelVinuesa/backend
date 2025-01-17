package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWildPayout(t *testing.T) {
	testCases := []struct {
		name   string
		symbol utils.Index
		count  uint8
		payout float64
	}{
		{"symbol 9, 3", 9, 3, 5},
		{"symbol 10, 5", 10, 5, 27.5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := NewWildPayoutAction(tc.symbol, tc.count, tc.payout)
			require.NotNil(t, w)

			assert.Equal(t, Payout, w.result)
			assert.Zero(t, w.nrOfSpins)
			assert.Equal(t, tc.symbol, w.symbol)
			assert.Equal(t, tc.count, w.wildCount)
			assert.Equal(t, tc.payout, w.wildPayout)
			assert.False(t, w.needHero)
			assert.False(t, w.expandToReel)
			assert.False(t, w.expandToSymbol)
			assert.Zero(t, len(w.expandGrid))
		})
	}
}

func TestNewWildExpansion(t *testing.T) {
	testCases := []struct {
		name           string
		nrOfSpins      uint8
		wildCount      uint8
		needHero       bool
		expandingWilds bool
		expandFirst    bool
	}{
		{"1 wild", 1, 1, false, false, false},
		{"2 wilds", 3, 2, false, false, false},
		{"need hero", 1, 1, true, false, false},
		{"expanding", 1, 1, true, true, false},
		{"expanding first", 1, 1, true, true, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := NewWildExpansion(tc.nrOfSpins, false, 9, tc.wildCount, tc.needHero, tc.expandingWilds, tc.expandFirst)
			require.NotNil(t, w)

			assert.Equal(t, FreeSpins, w.result)
			assert.Equal(t, tc.nrOfSpins, w.nrOfSpins)
			assert.Equal(t, utils.Index(9), w.symbol)
			assert.Equal(t, tc.wildCount, w.wildCount)
			assert.Equal(t, tc.needHero, w.needHero)
			assert.Equal(t, tc.expandingWilds, w.expandToReel)
			assert.False(t, w.expandToSymbol)
			assert.Zero(t, len(w.expandGrid))
		})
	}
}

func TestNewWildTransform(t *testing.T) {
	testCases := []struct {
		name       string
		symbol     utils.Index
		before     bool
		expandGrid GridOffsets
		wantGrid   GridOffsets
	}{
		{
			name:     "before, no grid",
			symbol:   9,
			before:   true,
			wantGrid: GridOffsets{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 0}, {0, 1}, {1, -1}, {1, 0}, {1, 1}},
		},
		{
			name:       "before, with grid",
			symbol:     10,
			before:     true,
			expandGrid: GridOffsets{{-1, 0}, {0, -1}, {0, 0}, {1, -1}, {1, 1}},
			wantGrid:   GridOffsets{{-1, 0}, {0, -1}, {0, 0}, {1, -1}, {1, 1}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := NewWildTransform(tc.symbol, tc.before, tc.expandGrid)
			require.NotNil(t, w)

			if tc.before {
				assert.Equal(t, ExpandBefore, w.stage)
			} else {
				assert.Equal(t, ExpandAfter, w.stage)
			}

			assert.Equal(t, Processed, w.result)
			assert.Zero(t, w.nrOfSpins)
			assert.Equal(t, tc.symbol, w.symbol)
			assert.Equal(t, uint8(1), w.wildCount)
			assert.False(t, w.needHero)
			assert.False(t, w.expandToReel)
			assert.True(t, w.expandToSymbol)
			assert.EqualValues(t, tc.wantGrid, w.expandGrid)
		})
	}
}

func TestWildTrigger_Expandable(t *testing.T) {
	testCases := []struct {
		name    string
		indexes utils.Indexes
		wilds   uint8
		heroes  uint8
		action  *WildAction
		want    bool
	}{
		{
			name:    "no wilds",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			action:  NewWildExpansion(1, false, 9, 1, false, true, false),
		},
		{
			name:    "1 wild",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 9, 5, 6, 1, 2, 3},
			wilds:   1,
			action:  NewWildExpansion(1, false, 9, 1, false, true, false),
			want:    true,
		},
		{
			name:    "2 wilds",
			indexes: utils.Indexes{1, 2, 9, 4, 5, 9, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			wilds:   2,
			action:  NewWildExpansion(1, false, 9, 1, false, true, false),
			want:    true,
		},
		{
			name:    "1 wild no hero",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 9, 1, 2, 3},
			wilds:   1,
			action:  NewWildExpansion(1, false, 9, 1, true, true, false),
		},
		{
			name:    "1 wild 1 hero",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 9},
			wilds:   1,
			heroes:  1,
			action:  NewWildExpansion(1, false, 9, 1, true, true, false),
			want:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spin := &Spin{reelCount: 5, newWilds: tc.wilds, newHeroes: tc.heroes}
			spin.indexes = tc.indexes
			got := tc.action.Triggered(spin) != nil
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestWildTrigger_Expand(t *testing.T) {
	s := NewSlots(Grid(5, 3), WithSymbols(setF1))

	testCases := []struct {
		name    string
		indexes utils.Indexes
		trigger *WildAction
		want    utils.Indexes
		locked  []bool
	}{
		{
			name:    "no wilds",
			indexes: utils.Indexes{0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2},
			trigger: NewWildExpansion(1, false, 9, 1, false, true, false),
			want:    utils.Indexes{0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2},
			locked:  []bool{false, false, false, false, false},
		},
		{
			name:    "1 wild",
			indexes: utils.Indexes{0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 9, 0, 1, 2},
			trigger: NewWildExpansion(1, false, 9, 1, false, true, false),
			want:    utils.Indexes{0, 1, 2, 0, 1, 2, 0, 1, 2, 9, 9, 9, 0, 1, 2},
			locked:  []bool{false, false, false, true, false},
		},
		{
			name:    "2 wilds",
			indexes: utils.Indexes{0, 1, 2, 0, 9, 2, 0, 1, 2, 0, 1, 9, 0, 1, 2},
			trigger: NewWildExpansion(1, false, 9, 1, true, true, false),
			want:    utils.Indexes{0, 1, 2, 9, 9, 9, 0, 1, 2, 9, 9, 9, 0, 1, 2},
			locked:  []bool{false, true, false, true, false},
		},
		{
			name:    "5 wilds",
			indexes: utils.Indexes{0, 1, 9, 0, 9, 2, 9, 1, 2, 0, 1, 9, 9, 1, 2},
			trigger: NewWildExpansion(1, false, 9, 1, true, true, false),
			want:    utils.Indexes{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
			locked:  []bool{true, true, true, true, true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			indexes := make(utils.Indexes, len(tc.indexes))
			copy(indexes, tc.indexes)

			spin := &Spin{slots: s, indexes: indexes, locked: make([]bool, 5)}

			tr := tc.trigger
			tr.WithLockReels()

			tc.trigger.Expand(spin)
			assert.EqualValues(t, tc.want, indexes)
			assert.EqualValues(t, tc.locked, spin.locked)
		})
	}
}

func TestWildTrigger_Payout(t *testing.T) {
	testCases := []struct {
		name    string
		indexes utils.Indexes
		payout  results.Payout
	}{
		{
			name:    "no wilds",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "2 wilds",
			indexes: utils.Indexes{1, 2, 12, 4, 5, 6, 1, 2, 3, 4, 12, 6, 1, 2, 3},
		},
		{
			name:    "3 wilds",
			indexes: utils.Indexes{1, 2, 12, 4, 5, 12, 1, 2, 3, 4, 12, 6, 1, 2, 3},
			payout:  WildSymbolPayoutWithMap(2, 0, 12, 3, utils.UInt8s{2, 5, 10}),
		},
		{
			name:    "4 wilds",
			indexes: utils.Indexes{1, 2, 12, 4, 5, 6, 12, 2, 3, 4, 12, 6, 1, 2, 12},
			payout:  WildSymbolPayoutWithMap(8, 0, 12, 4, utils.UInt8s{2, 6, 10, 14}),
		},
		{
			name:    "5 wilds",
			indexes: utils.Indexes{1, 2, 12, 4, 5, 6, 1, 12, 3, 4, 12, 6, 1, 12, 12},
			payout:  WildSymbolPayoutWithMap(15, 0, 12, 5, utils.UInt8s{2, 7, 10, 13, 14}),
		},
		{
			name:    "7 wilds",
			indexes: utils.Indexes{1, 2, 12, 4, 12, 6, 12, 12, 3, 4, 12, 6, 1, 12, 12},
			payout:  WildSymbolPayoutWithMap(15, 0, 12, 7, utils.UInt8s{2, 4, 6, 7, 10, 13, 14}),
		},
	}

	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	a3 := NewWildPayoutAction(12, 3, 2)
	a4 := NewWildPayoutAction(12, 4, 8).WithAlternate(a3)
	a5 := NewWildPayoutAction(12, 5, 15).WithAlternate(a4)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spin := AcquireSpin(slots, prng)
			defer spin.Release()

			spin.indexes = tc.indexes

			res := results.AcquireResult(nil, 0)
			require.NotNil(t, res)
			defer res.Release()

			a := a5.Payout(spin, res)

			if tc.payout == nil {
				require.Nil(t, a)
				assert.Zero(t, res.Total)
				assert.Zero(t, len(res.Payouts))
			} else {
				require.NotNil(t, a)
				assert.Equal(t, tc.payout.(*SpinPayout).Total(), res.Total)
				require.Equal(t, 1, len(res.Payouts))

				p := res.Payouts[0]
				require.NotNil(t, p)

				if !tc.payout.(*SpinPayout).Equals(p.(*SpinPayout)) {
					assert.Equal(t, tc.payout, p)
				}
			}
		})
	}
}

func TestWildTrigger_PayoutFail(t *testing.T) {
	t.Run("wild trigger payout fail", func(t *testing.T) {
		defer func() {
			e := recover()
			require.NotNilf(t, e, "should have failed")
		}()

		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

		spin := AcquireSpin(slots, prng)
		defer spin.Release()

		spin.indexes = utils.Indexes{1, 0, 0, 2, 0, 0, 3, 0, 0, 4, 0, 0, 5, 0, 0}

		tr := NewWildPayoutAction(0, 2, 5)

		tr.Payout(spin, nil)
		require.Nil(t, tr, "should have failed")
	})
}

func TestPurgeWildTriggers(t *testing.T) {
	t1 := NewWildExpansion(1, false, 9, 1, false, true, false)
	t2 := NewWildExpansion(1, false, 9, 1, false, true, false)
	t3 := NewWildExpansion(1, false, 9, 1, false, true, false)
	t4 := NewWildExpansion(1, false, 9, 1, false, true, false)
	t5 := NewWildExpansion(1, false, 9, 1, false, true, false)
	t6 := NewWildExpansion(1, false, 9, 1, false, true, false)
	t7 := NewWildExpansion(1, false, 9, 1, false, true, false)

	testCases := []struct {
		name    string
		in      WildActions
		cap     int
		wantCap int
	}{
		{"nil", nil, 5, 5},
		{"empty", WildActions{}, 5, 5},
		{"short", WildActions{t1, t2, t3}, 5, 5},
		{"exact", WildActions{t1, t2, t3, t4, t5}, 5, 5},
		{"long", WildActions{t1, t2, t3, t4, t5, t6, t7}, 5, 7},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := PurgeWildActions(tc.in, tc.cap)
			require.NotNil(t, out)
			assert.Equal(t, tc.wantCap, cap(out))
			assert.Zero(t, len(out))
		})
	}
}
