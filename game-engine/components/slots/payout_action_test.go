package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewPaylinesAction(t *testing.T) {
	p1 := &Payline{1, nil, nil}
	p2 := &Payline{2, nil, nil}
	p3 := &Payline{3, nil, nil}
	p4 := &Payline{4, nil, nil}
	p5 := &Payline{5, nil, nil}
	p6 := &Payline{6, nil, nil}
	p7 := &Payline{7, nil, nil}
	p8 := &Payline{8, nil, nil}
	p9 := &Payline{9, nil, nil}
	p10 := &Payline{10, nil, nil}
	p11 := &Payline{11, nil, nil}
	p12 := &Payline{12, nil, nil}

	testCases := []struct {
		name     string
		paylines Paylines
	}{
		{"1", Paylines{p1}},
		{"2", Paylines{p1, p2}},
		{"5", Paylines{p1, p2, p3, p4, p5}},
		{"8", Paylines{p1, p2, p3, p4, p5, p6, p7, p8}},
		{"12", Paylines{p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewPaylinesAction()
			require.NotNil(t, a)

			assert.True(t, a.paylines)
			assert.False(t, a.allPaylines)
			assert.Nil(t, a.cluster)
		})
	}
}

func TestNewAllPaylinesAction(t *testing.T) {
	t.Run("new all paylines", func(t *testing.T) {
		a := NewAllPaylinesAction(true)
		require.NotNil(t, a)

		assert.False(t, a.paylines)
		assert.True(t, a.allPaylines)
		assert.Nil(t, a.cluster)
	})
}

func TestNewClusterPaysAction(t *testing.T) {
	t.Run("new cluster pays", func(t *testing.T) {
		a := NewClusterPayoutsAction(5, 5)
		require.NotNil(t, a)

		assert.False(t, a.paylines)
		assert.False(t, a.allPaylines)

		require.NotNil(t, a.cluster)
		assert.Equal(t, uint8(5), a.cluster.reels)
		assert.Equal(t, uint8(5), a.cluster.rows)
	})
}

func TestPayoutAction_Payout(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	p1 := NewPayline(1, 3, 1, 1, 1, 1, 1)
	p2 := NewPayline(2, 3, 0, 0, 0, 0, 0)
	p3 := NewPayline(3, 3, 2, 2, 2, 2, 2)
	p4 := NewPayline(4, 3, 1, 2, 2, 2, 1)
	p5 := NewPayline(5, 3, 1, 0, 0, 0, 1)

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1), WithPaylines(PayLTR, true, p1, p2, p3, p4, p5))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name        string
		a           *PayoutAction
		indexes     utils.Indexes
		multipliers []uint16
		want        results.Payouts
	}{
		{
			name:    "paylines - no win",
			a:       NewPaylinesAction(),
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			want:    results.Payouts{},
		},
		{
			name:    "paylines - 1 win",
			a:       NewPaylinesAction(),
			indexes: utils.Indexes{1, 2, 3, 1, 5, 6, 1, 2, 3, 1, 5, 6, 1, 2, 3},
			want: results.Payouts{
				WinlinePayoutFromData(4, 1, 1, 5, PayLTR, 2, utils.UInt8s{0, 0, 0, 0, 0}),
			},
		},
		{
			name:    "paylines - 3 wins",
			a:       NewPaylinesAction(),
			indexes: utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 6, 1, 5, 3},
			want: results.Payouts{
				WinlinePayoutFromData(2, 1, 2, 4, PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
				WinlinePayoutFromData(4, 1, 1, 5, PayLTR, 2, utils.UInt8s{0, 0, 0, 0, 0}),
				WinlinePayoutFromData(1, 1, 3, 3, PayLTR, 3, utils.UInt8s{2, 2, 2, 2, 2}),
			},
		},
		{
			name:    "paylines - 1 win with wilds (no highest payout)",
			a:       NewPaylinesAction(),
			indexes: utils.Indexes{1, 2, 9, 4, 5, 1, 1, 2, 9, 3, 4, 5, 1, 2, 3},
			want: results.Payouts{
				WinlinePayoutFromData(0.5, 1, 1, 3, PayLTR, 3, utils.UInt8s{2, 2, 2, 2, 2}),
			},
		},
		{
			name:    "paylines - 2 wins with wilds (no highest payout)",
			a:       NewPaylinesAction(),
			indexes: utils.Indexes{1, 2, 3, 4, 5, 9, 1, 2, 9, 3, 4, 5, 1, 2, 3},
			want: results.Payouts{
				WinlinePayoutFromData(0.5, 1, 2, 3, PayLTR, 4, utils.UInt8s{1, 2, 2, 2, 1}),
				WinlinePayoutFromData(1, 1, 3, 3, PayLTR, 3, utils.UInt8s{2, 2, 2, 2, 2}),
			},
		},
		{
			name:    "paylines - 2 wins with wilds (highest payout)",
			a:       NewPaylinesAction(),
			indexes: utils.Indexes{1, 9, 3, 4, 5, 9, 1, 2, 9, 3, 4, 1, 1, 2, 3},
			want: results.Payouts{
				WinlinePayoutFromData(4, 1, 12, 3, PayLTR, 4, utils.UInt8s{1, 2, 2, 2, 1}),
				WinlinePayoutFromData(1, 1, 3, 3, PayLTR, 3, utils.UInt8s{2, 2, 2, 2, 2}),
			},
		},
		{
			name:    "all paylines - no win",
			a:       NewAllPaylinesAction(true),
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			want:    results.Payouts{},
		},
		{
			name:    "all paylines - 1 win",
			a:       NewAllPaylinesAction(true),
			indexes: utils.Indexes{1, 2, 3, 4, 5, 1, 1, 2, 3, 4, 5, 1, 2, 5, 3},
			want: results.Payouts{
				AllPaylinePayout(2, 1, 1, 4, utils.UInt8s{0, 2, 0, 2, 0}).Clone().(*SpinPayout),
			},
		},
		{
			name:    "all paylines - 2 wins",
			a:       NewAllPaylinesAction(true),
			indexes: utils.Indexes{9, 2, 3, 1, 5, 1, 1, 2, 3, 4, 5, 1, 2, 1, 3},
			want: results.Payouts{
				AllPaylinePayout(4, 1, 1, 5, utils.UInt8s{0, 0, 0, 2, 1}).Clone().(*SpinPayout),
				AllPaylinePayout(4, 1, 1, 5, utils.UInt8s{0, 2, 0, 2, 1}).Clone().(*SpinPayout),
			},
		},
		{
			name:    "all paylines - 2 wins + wilds",
			a:       NewAllPaylinesAction(true),
			indexes: utils.Indexes{1, 2, 3, 1, 5, 1, 12, 2, 3, 4, 5, 9, 2, 4, 3},
			want: results.Payouts{
				AllPaylinePayout(2, 1, 1, 4, utils.UInt8s{0, 0, 0, 2, 0}).Clone().(*SpinPayout),
				AllPaylinePayout(2, 1, 1, 4, utils.UInt8s{0, 2, 0, 2, 0}).Clone().(*SpinPayout),
			},
		},
		{
			name:        "all paylines - 2 wins + multiplier wilds",
			a:           NewAllPaylinesAction(true),
			indexes:     utils.Indexes{1, 2, 3, 1, 5, 1, 12, 2, 3, 4, 5, 9, 2, 4, 3},
			multipliers: []uint16{0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want: results.Payouts{
				AllPaylinePayout(2, 3, 1, 4, utils.UInt8s{0, 0, 0, 2, 0}).Clone().(*SpinPayout),
				AllPaylinePayout(2, 3, 1, 4, utils.UInt8s{0, 2, 0, 2, 0}).Clone().(*SpinPayout),
			},
		},
		{
			name:    "all paylines - many wins + multiplier wilds",
			a:       NewAllPaylinesAction(true),
			indexes: utils.Indexes{9, 3, 2, 9, 5, 4, 12, 12, 9, 1, 1, 1, 5, 4, 9},
			//  9  9 12  1  5
			//  3  5 12  1  4
			//  2  4  9  1  9
			multipliers: []uint16{0, 0, 0, 0, 0, 0, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0},
			want: results.Payouts{
				AllPaylinePayout(4, 2, 1, 5, utils.UInt8s{0, 0, 0, 0, 2}),
				AllPaylinePayout(4, 2, 1, 5, utils.UInt8s{0, 0, 0, 1, 2}),
				AllPaylinePayout(4, 2, 1, 5, utils.UInt8s{0, 0, 0, 2, 2}),
				AllPaylinePayout(4, 2, 1, 5, utils.UInt8s{0, 0, 1, 0, 2}),
				AllPaylinePayout(4, 2, 1, 5, utils.UInt8s{0, 0, 1, 1, 2}),
				AllPaylinePayout(4, 2, 1, 5, utils.UInt8s{0, 0, 1, 2, 2}),
				AllPaylinePayout(4, 1, 1, 5, utils.UInt8s{0, 0, 2, 0, 2}),
				AllPaylinePayout(4, 1, 1, 5, utils.UInt8s{0, 0, 2, 1, 2}),
				AllPaylinePayout(4, 1, 1, 5, utils.UInt8s{0, 0, 2, 2, 2}),
				AllPaylinePayout(1, 2, 5, 3, utils.UInt8s{0, 1, 0, 0, 0}),
				AllPaylinePayout(1, 2, 5, 3, utils.UInt8s{0, 1, 1, 0, 0}),
				AllPaylinePayout(1, 1, 5, 3, utils.UInt8s{0, 1, 2, 0, 0}),
				AllPaylinePayout(1, 2, 4, 3, utils.UInt8s{0, 2, 0, 0, 0}),
				AllPaylinePayout(1, 2, 4, 3, utils.UInt8s{0, 2, 1, 0, 0}),
				AllPaylinePayout(1, 1, 4, 3, utils.UInt8s{0, 2, 2, 0, 0}),
				AllPaylinePayout(1, 2, 3, 3, utils.UInt8s{1, 0, 0, 0, 0}),
				AllPaylinePayout(1, 2, 3, 3, utils.UInt8s{1, 0, 1, 0, 0}),
				AllPaylinePayout(1, 1, 3, 3, utils.UInt8s{1, 0, 2, 0, 0}),
				AllPaylinePayout(0.5, 2, 2, 3, utils.UInt8s{2, 0, 0, 0, 0}),
				AllPaylinePayout(0.5, 2, 2, 3, utils.UInt8s{2, 0, 1, 0, 0}),
				AllPaylinePayout(0.5, 1, 2, 3, utils.UInt8s{2, 0, 2, 0, 0}),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spin.ResetSpin()
			copy(spin.indexes, tc.indexes)

			if l := len(tc.multipliers); l > 0 {
				spin.multipliers = spin.multipliers[:l]
				copy(spin.multipliers, tc.multipliers)
			} else {
				spin.multipliers = spin.multipliers[:0]
			}

			res := results.AcquireResult(nil, 0)
			defer res.Release()

			tc.a.Payout(spin, res)
			require.Equal(t, len(tc.want), len(res.Payouts))

			for ix := range tc.want {
				res1, ok1 := tc.want[ix].(*SpinPayout)
				res2, ok2 := res.Payouts[ix].(*SpinPayout)

				require.True(t, ok1)
				require.True(t, ok2)

				if !res1.Equals(res2) {
					assert.EqualValues(t, res1, res2)
				}
			}
		})
	}
}

func TestMGD_Payouts(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	symbols := NewSymbolSet(
		NewSymbol(1, WithPayouts(0, 0, 0.5, 0.5, 0.8, 1.5)),
		NewSymbol(2, WithPayouts(0, 0, 0.5, 0.5, 0.8, 1.5)),
		NewSymbol(3, WithPayouts(0, 0, 0.5, 0.8, 1, 2)),
		NewSymbol(4, WithPayouts(0, 0, 0.5, 0.8, 1, 2)),
		NewSymbol(5, WithPayouts(0, 0, 0.5, 1, 1.5, 2.5)),
		NewSymbol(6, WithPayouts(0, 0, 0.8, 1, 2, 3)),
		NewSymbol(7, WithPayouts(0, 0, 1, 1.5, 3, 5)),
		NewSymbol(8, WithPayouts(0, 0, 1.5, 3, 4, 8)),
		NewSymbol(9, WithPayouts(0, 0, 0, 0, 0, 0), WithKind(Scatter)),
		NewSymbol(10, WithPayouts(0, 0, 0, 0, 0, 0), WithKind(Wild)),
		NewSymbol(11, WithPayouts(0, 0, 0, 0, 0, 0), WithKind(Wild), WithMultiplier(2)),
		NewSymbol(12, WithPayouts(0, 0, 0, 0, 0, 0), WithKind(Wild), WithMultiplier(3)),
		NewSymbol(13, WithPayouts(0, 0, 0, 0, 0, 0), WithKind(Wild), WithMultiplier(4)),
	)

	slots := NewSlots(
		Grid(6, 4),
		WithMask(2, 3, 4, 4, 3, 2),
		WithSymbols(symbols),
	)

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name    string
		indexes utils.Indexes
		want    results.Payouts
	}{
		{
			name:    "test1",
			indexes: utils.Indexes{10, 1, 0, 0, 10, 2, 3, 0, 10, 5, 6, 7, 1, 2, 3, 4, 5, 6, 7, 0, 1, 2, 0, 0},
			want: results.Payouts{
				AllPaylinePayout(1.5, 1, 8, 3, utils.UInt8s{0, 0, 0, 0, 0, 0}),
				AllPaylinePayout(1.5, 1, 8, 3, utils.UInt8s{0, 0, 0, 1, 0, 0}),
				AllPaylinePayout(1.5, 1, 8, 3, utils.UInt8s{0, 0, 0, 2, 0, 0}),
				AllPaylinePayout(1.5, 1, 8, 3, utils.UInt8s{0, 0, 0, 3, 0, 0}),
				AllPaylinePayout(0.5, 1, 5, 3, utils.UInt8s{0, 0, 1, 0, 0, 0}),
				AllPaylinePayout(0.8, 1, 6, 3, utils.UInt8s{0, 0, 2, 0, 0, 0}),
				AllPaylinePayout(1.0, 1, 7, 3, utils.UInt8s{0, 0, 3, 0, 0, 0}),
				AllPaylinePayout(0.5, 1, 2, 4, utils.UInt8s{0, 1, 0, 1, 0, 0}),
				AllPaylinePayout(0.8, 1, 3, 4, utils.UInt8s{0, 2, 0, 2, 0, 0}),
				AllPaylinePayout(0.5, 1, 1, 4, utils.UInt8s{1, 0, 0, 0, 0, 0}),
			},
		},
		{
			name:    "test2",
			indexes: utils.Indexes{10, 10, 0, 0, 10, 5, 3, 0, 11, 11, 10, 9, 6, 1, 1, 4, 1, 2, 10, 0, 10, 10, 0, 0},
			want: results.Payouts{
				AllPaylinePayout(3, 2, 6, 6, utils.UInt8s{0, 0, 0, 0, 2, 0}),
				AllPaylinePayout(3, 2, 6, 6, utils.UInt8s{0, 0, 0, 0, 2, 1}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{0, 0, 0, 1, 0, 0}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{0, 0, 0, 1, 0, 1}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{0, 0, 0, 1, 2, 0}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{0, 0, 0, 1, 2, 1}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{0, 0, 0, 2, 0, 0}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{0, 0, 0, 2, 0, 1}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{0, 0, 0, 2, 2, 0}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{0, 0, 0, 2, 2, 1}),
				AllPaylinePayout(2, 2, 4, 6, utils.UInt8s{0, 0, 0, 3, 2, 0}),
				AllPaylinePayout(2, 2, 4, 6, utils.UInt8s{0, 0, 0, 3, 2, 1}),
				AllPaylinePayout(3, 2, 6, 6, utils.UInt8s{0, 0, 1, 0, 2, 0}),
				AllPaylinePayout(3, 2, 6, 6, utils.UInt8s{0, 0, 1, 0, 2, 1}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{0, 0, 1, 1, 0, 0}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{0, 0, 1, 1, 0, 1}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{0, 0, 1, 1, 2, 0}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{0, 0, 1, 1, 2, 1}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{0, 0, 1, 2, 0, 0}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{0, 0, 1, 2, 0, 1}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{0, 0, 1, 2, 2, 0}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{0, 0, 1, 2, 2, 1}),
				AllPaylinePayout(2, 2, 4, 6, utils.UInt8s{0, 0, 1, 3, 2, 0}),
				AllPaylinePayout(2, 2, 4, 6, utils.UInt8s{0, 0, 1, 3, 2, 1}),
				AllPaylinePayout(3, 1, 6, 6, utils.UInt8s{0, 0, 2, 0, 2, 0}),
				AllPaylinePayout(3, 1, 6, 6, utils.UInt8s{0, 0, 2, 0, 2, 1}),
				AllPaylinePayout(1.5, 1, 1, 6, utils.UInt8s{0, 0, 2, 1, 0, 0}),
				AllPaylinePayout(1.5, 1, 1, 6, utils.UInt8s{0, 0, 2, 1, 0, 1}),
				AllPaylinePayout(1.5, 1, 1, 6, utils.UInt8s{0, 0, 2, 1, 2, 0}),
				AllPaylinePayout(1.5, 1, 1, 6, utils.UInt8s{0, 0, 2, 1, 2, 1}),
				AllPaylinePayout(1.5, 1, 1, 6, utils.UInt8s{0, 0, 2, 2, 0, 0}),
				AllPaylinePayout(1.5, 1, 1, 6, utils.UInt8s{0, 0, 2, 2, 0, 1}),
				AllPaylinePayout(1.5, 1, 1, 6, utils.UInt8s{0, 0, 2, 2, 2, 0}),
				AllPaylinePayout(1.5, 1, 1, 6, utils.UInt8s{0, 0, 2, 2, 2, 1}),
				AllPaylinePayout(2, 1, 4, 6, utils.UInt8s{0, 0, 2, 3, 2, 0}),
				AllPaylinePayout(2, 1, 4, 6, utils.UInt8s{0, 0, 2, 3, 2, 1}),
				AllPaylinePayout(0.5, 2, 5, 3, utils.UInt8s{0, 1, 0, 0, 0, 0}),
				AllPaylinePayout(0.5, 2, 5, 3, utils.UInt8s{0, 1, 1, 0, 0, 0}),
				AllPaylinePayout(0.5, 1, 5, 3, utils.UInt8s{0, 1, 2, 0, 0, 0}),
				AllPaylinePayout(0.5, 2, 3, 3, utils.UInt8s{0, 2, 0, 0, 0, 0}),
				AllPaylinePayout(0.5, 2, 3, 3, utils.UInt8s{0, 2, 1, 0, 0, 0}),
				AllPaylinePayout(0.5, 1, 3, 3, utils.UInt8s{0, 2, 2, 0, 0, 0}),
				AllPaylinePayout(3, 2, 6, 6, utils.UInt8s{1, 0, 0, 0, 2, 0}),
				AllPaylinePayout(3, 2, 6, 6, utils.UInt8s{1, 0, 0, 0, 2, 1}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{1, 0, 0, 1, 0, 0}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{1, 0, 0, 1, 0, 1}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{1, 0, 0, 1, 2, 0}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{1, 0, 0, 1, 2, 1}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{1, 0, 0, 2, 0, 0}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{1, 0, 0, 2, 0, 1}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{1, 0, 0, 2, 2, 0}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{1, 0, 0, 2, 2, 1}),
				AllPaylinePayout(2, 2, 4, 6, utils.UInt8s{1, 0, 0, 3, 2, 0}),
				AllPaylinePayout(2, 2, 4, 6, utils.UInt8s{1, 0, 0, 3, 2, 1}),
				AllPaylinePayout(3, 2, 6, 6, utils.UInt8s{1, 0, 1, 0, 2, 0}),
				AllPaylinePayout(3, 2, 6, 6, utils.UInt8s{1, 0, 1, 0, 2, 1}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{1, 0, 1, 1, 0, 0}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{1, 0, 1, 1, 0, 1}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{1, 0, 1, 1, 2, 0}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{1, 0, 1, 1, 2, 1}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{1, 0, 1, 2, 0, 0}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{1, 0, 1, 2, 0, 1}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{1, 0, 1, 2, 2, 0}),
				AllPaylinePayout(1.5, 2, 1, 6, utils.UInt8s{1, 0, 1, 2, 2, 1}),
				AllPaylinePayout(2, 2, 4, 6, utils.UInt8s{1, 0, 1, 3, 2, 0}),
				AllPaylinePayout(2, 2, 4, 6, utils.UInt8s{1, 0, 1, 3, 2, 1}),
				AllPaylinePayout(3, 1, 6, 6, utils.UInt8s{1, 0, 2, 0, 2, 0}),
				AllPaylinePayout(3, 1, 6, 6, utils.UInt8s{1, 0, 2, 0, 2, 1}),
				AllPaylinePayout(1.5, 1, 1, 6, utils.UInt8s{1, 0, 2, 1, 0, 0}),
				AllPaylinePayout(1.5, 1, 1, 6, utils.UInt8s{1, 0, 2, 1, 0, 1}),
				AllPaylinePayout(1.5, 1, 1, 6, utils.UInt8s{1, 0, 2, 1, 2, 0}),
				AllPaylinePayout(1.5, 1, 1, 6, utils.UInt8s{1, 0, 2, 1, 2, 1}),
				AllPaylinePayout(1.5, 1, 1, 6, utils.UInt8s{1, 0, 2, 2, 0, 0}),
				AllPaylinePayout(1.5, 1, 1, 6, utils.UInt8s{1, 0, 2, 2, 0, 1}),
				AllPaylinePayout(1.5, 1, 1, 6, utils.UInt8s{1, 0, 2, 2, 2, 0}),
				AllPaylinePayout(1.5, 1, 1, 6, utils.UInt8s{1, 0, 2, 2, 2, 1}),
				AllPaylinePayout(2, 1, 4, 6, utils.UInt8s{1, 0, 2, 3, 2, 0}),
				AllPaylinePayout(2, 1, 4, 6, utils.UInt8s{1, 0, 2, 3, 2, 1}),
				AllPaylinePayout(0.5, 2, 5, 3, utils.UInt8s{1, 1, 0, 0, 0, 0}),
				AllPaylinePayout(0.5, 2, 5, 3, utils.UInt8s{1, 1, 1, 0, 0, 0}),
				AllPaylinePayout(0.5, 1, 5, 3, utils.UInt8s{1, 1, 2, 0, 0, 0}),
				AllPaylinePayout(0.5, 2, 3, 3, utils.UInt8s{1, 2, 0, 0, 0, 0}),
				AllPaylinePayout(0.5, 2, 3, 3, utils.UInt8s{1, 2, 1, 0, 0, 0}),
				AllPaylinePayout(0.5, 1, 3, 3, utils.UInt8s{1, 2, 2, 0, 0, 0}),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spin.ResetSpin()
			copy(spin.indexes, tc.indexes)

			res := results.AcquireResult(nil, 0)
			defer res.Release()

			a := NewAllPaylinesAction(true)
			a.Payout(spin, res)
			require.Equal(t, len(tc.want), len(res.Payouts))

			for ix := range tc.want {
				res1, ok1 := tc.want[ix].(*SpinPayout)
				res2, ok2 := res.Payouts[ix].(*SpinPayout)

				require.True(t, ok1)
				require.True(t, ok2)

				if !res1.Equals(res2) {
					require.EqualValues(t, res1, res2)
				}
			}
		})
	}
}

func TestPurgePayoutActions(t *testing.T) {
	testCases := []struct {
		name     string
		actions  PayoutActions
		capacity int
		want     int
	}{
		{"nil, 2", nil, 2, 2},
		{"nil, 5", nil, 5, 5},
		{"2, 1", PayoutActions{&PayoutAction{}, &PayoutAction{}}, 1, 2},
		{"2, 2", PayoutActions{&PayoutAction{}, &PayoutAction{}}, 2, 2},
		{"2, 5", PayoutActions{&PayoutAction{}, &PayoutAction{}}, 5, 5},
		{"4, 1", PayoutActions{&PayoutAction{}, &PayoutAction{}, &PayoutAction{}, &PayoutAction{}}, 1, 4},
		{"4, 2", PayoutActions{&PayoutAction{}, &PayoutAction{}, &PayoutAction{}, &PayoutAction{}}, 2, 4},
		{"4, 5", PayoutActions{&PayoutAction{}, &PayoutAction{}, &PayoutAction{}, &PayoutAction{}}, 5, 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := PurgePayoutActions(tc.actions, tc.capacity)
			require.NotNil(t, p)
			assert.Zero(t, len(p))
			assert.Equal(t, tc.want, cap(p))
		})
	}
}

func TestPayoutAction_ClusterPayout(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	mask := utils.UInt8s{4, 5, 6, 5, 4}
	slots := NewSlots(Grid(5, 6), WithSymbols(setF1), PayDirections(PayLTR), WithMask(mask...))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name    string
		a       *PayoutAction
		indexes utils.Indexes
		want    results.Payouts
	}{
		{
			name:    "no hits",
			a:       NewClusterPayoutsAction(5, 6, ClusterGridMask(mask, Hexagonal)),
			indexes: utils.Indexes{1, 2, 3, 4, 0, 0, 3, 4, 5, 6, 7, 0, 1, 2, 3, 4, 5, 6, 3, 4, 5, 6, 7, 0, 1, 2, 3, 4, 0, 0},
			want:    results.Payouts{},
		},
		{
			name:    "1 hit",
			a:       NewClusterPayoutsAction(5, 6, ClusterGridMask(mask, Hexagonal)),
			indexes: utils.Indexes{1, 2, 3, 4, 0, 0, 3, 4, 5, 5, 7, 0, 1, 2, 3, 5, 5, 6, 3, 4, 5, 6, 7, 0, 1, 2, 3, 4, 0, 0},
			want: results.Payouts{
				ClusterPayout(7, 0, 5, 5, utils.UInt8s{8, 9, 15, 16, 20}).Clone().(*SpinPayout),
			},
		},
		{
			name:    "3 hits",
			a:       NewClusterPayoutsAction(5, 6, ClusterGridMask(mask, Hexagonal)),
			indexes: utils.Indexes{1, 2, 3, 4, 0, 0, 3, 2, 5, 5, 7, 0, 1, 2, 3, 5, 5, 6, 3, 4, 5, 6, 7, 0, 1, 2, 2, 2, 0, 0},
			want: results.Payouts{
				ClusterPayout(0.5, 0, 2, 3, utils.UInt8s{1, 7, 13}).Clone().(*SpinPayout),
				ClusterPayout(7, 0, 5, 5, utils.UInt8s{8, 9, 15, 16, 20}).Clone().(*SpinPayout),
				ClusterPayout(0.5, 0, 2, 3, utils.UInt8s{25, 26, 27}).Clone().(*SpinPayout),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spin.ResetSpin()
			spin.indexes = tc.indexes

			res := results.AcquireResult(nil, 0)
			defer res.Release()

			tc.a.Payout(spin, res)
			assert.Equal(t, len(tc.want), len(res.Payouts))
		})
	}
}

func TestNewRemovePayoutsAction(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1), PayDirections(PayLTR))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name    string
		min     float64
		max     float64
		chance  float64
		wilds   bool
		dupes   bool
		indexes utils.Indexes
		res     *results.Result
		want    bool
	}{
		{
			name:    "not in range, low",
			min:     5.0,
			max:     10.0,
			chance:  100.0,
			wilds:   true,
			indexes: utils.Indexes{1, 2, 3, 1, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			res: &results.Result{
				Total: 4.99,
				Payouts: results.Payouts{
					WinlinePayoutFromData(4.99, 1.0, 1, 3, PayLTR, 1, utils.UInt8s{0, 0, 0, 0, 0}),
				},
				Animations: results.Animations{},
			},
		},
		{
			name:    "not in range, high",
			min:     5.0,
			max:     10.0,
			chance:  100.0,
			wilds:   true,
			indexes: utils.Indexes{1, 2, 3, 1, 5, 3, 1, 2, 6, 1, 5, 6, 1, 2, 3},
			res: &results.Result{
				Total: 10.1,
				Payouts: results.Payouts{
					WinlinePayoutFromData(10, 1.0, 1, 5, PayLTR, 1, utils.UInt8s{0, 0, 0, 0, 0}),
					WinlinePayoutFromData(0.1, 1.0, 3, 2, PayLTR, 2, utils.UInt8s{2, 2, 2, 2, 2}),
				},
				Animations: results.Animations{},
			},
		},
		{
			name:    "in range",
			min:     5.0,
			max:     10.0,
			chance:  100.0,
			wilds:   true,
			indexes: utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 6, 1, 2, 6, 1, 5, 3},
			res: &results.Result{
				Total: 6.5,
				Payouts: results.Payouts{
					WinlinePayoutFromData(5, 1.0, 1, 5, PayLTR, 1, utils.UInt8s{0, 0, 0, 0, 0}),
					WinlinePayoutFromData(1, 1.0, 2, 4, PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
					WinlinePayoutFromData(0.5, 1.0, 3, 3, PayLTR, 1, utils.UInt8s{2, 2, 2, 2, 2}),
				},
				Animations: results.Animations{},
			},
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewRemovePayoutsAction(tc.min, tc.max, tc.chance, PayLTR, 1, tc.wilds, tc.dupes)
			require.NotNil(t, a)
			assert.Equal(t, tc.min, a.remMinFactor)
			assert.Equal(t, tc.max, a.remMaxFactor)
			assert.Equal(t, tc.chance, a.remPayoutChance)
			assert.Equal(t, PayLTR, a.remPayoutDir)
			assert.Equal(t, 1, a.remPayoutMech)
			assert.Equal(t, tc.wilds, a.remPayoutWilds)
			assert.Equal(t, tc.dupes, a.remDupes)

			copy(spin.indexes, tc.indexes)

			res := results.AcquireResult(nil, 0, tc.res.Payouts...)
			defer res.Release()

			a.Payout(spin, res)
			if tc.want {
				assert.NotEqualValues(t, tc.indexes, spin.indexes)
				assert.Zero(t, len(res.Payouts))
			} else {
				assert.EqualValues(t, tc.indexes, spin.indexes)
				assert.Equal(t, tc.res.Total, res.Total)
				for ix := range res.Payouts {
					p1 := res.Payouts[ix].(*SpinPayout)
					p2 := tc.res.Payouts[ix].(*SpinPayout)
					assert.Equal(t, p1.Kind(), p2.Kind())
					assert.Equal(t, p1.Count(), p2.Count())
					assert.Equal(t, p1.Symbol(), p2.Symbol())
					assert.Equal(t, p1.Factor(), p2.Factor())
					assert.Equal(t, p1.Multiplier(), p2.Multiplier())
					assert.EqualValues(t, p1.PayRows(), p2.PayRows())
					assert.EqualValues(t, p1.PayMap(), p2.PayMap())
				}
			}
		})
	}
}
