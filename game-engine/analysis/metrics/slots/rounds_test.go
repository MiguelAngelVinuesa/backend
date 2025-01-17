package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRounds(t *testing.T) {
	type result struct {
		bet         int64
		win         int64
		freeSpins   uint64
		refillSpins uint64
		superSpins  uint64
	}

	testCases := []struct {
		name    string
		results []result
		want    *Rounds
	}{
		{
			name: "empty",
			want: &Rounds{
				Bets:           &MinMaxInt64{Counts: map[int64]uint64{}, first: true},
				BetsNoFree:     &MinMaxInt64{Counts: map[int64]uint64{}, first: true},
				BetsFree:       &MinMaxInt64{Counts: map[int64]uint64{}, first: true},
				Wins:           &MinMaxInt64{Counts: map[int64]uint64{}, first: true},
				WinsNoFree:     &MinMaxInt64{Counts: map[int64]uint64{}, first: true},
				WinsFree:       &MinMaxInt64{Counts: map[int64]uint64{}, first: true},
				FreeSpins:      &MinMaxUInt64{Counts: map[uint64]uint64{}, first: true},
				RefillSpins:    &MinMaxUInt64{Counts: map[uint64]uint64{}, first: true},
				SuperSpins:     &MinMaxUInt64{Counts: map[uint64]uint64{}, first: true},
				FreeSpinRounds: []uint64{},
				RefillRounds:   []uint64{},
				SuperRounds:    []uint64{},
				SymbolsUsed:    []uint64{},
				SymbolsNoFree:  []uint64{},
				SymbolsFree:    []uint64{},
			},
		},
		{
			name: "one, no free spins",
			results: []result{
				{bet: 100, win: 500},
			},
			want: &Rounds{
				Count:          1,
				Bets:           &MinMaxInt64{Count: 1, Total: 100, Min: 100, Max: 100, Counts: map[int64]uint64{100: 1}},
				BetsNoFree:     &MinMaxInt64{Count: 1, Total: 100, Min: 100, Max: 100, Counts: map[int64]uint64{100: 1}},
				BetsFree:       &MinMaxInt64{Counts: map[int64]uint64{}, first: true},
				Wins:           &MinMaxInt64{Count: 1, Total: 500, Min: 500, Max: 500, Counts: map[int64]uint64{500: 1}},
				WinsNoFree:     &MinMaxInt64{Count: 1, Total: 500, Min: 500, Max: 500, Counts: map[int64]uint64{500: 1}},
				WinsFree:       &MinMaxInt64{Counts: map[int64]uint64{}, first: true},
				FreeSpins:      &MinMaxUInt64{Count: 1, Counts: map[uint64]uint64{0: 1}},
				RefillSpins:    &MinMaxUInt64{Count: 1, Counts: map[uint64]uint64{0: 1}},
				SuperSpins:     &MinMaxUInt64{Count: 1, Counts: map[uint64]uint64{0: 1}},
				FreeSpinRounds: []uint64{1},
				RefillRounds:   []uint64{1},
				SuperRounds:    []uint64{1},
				SymbolsUsed:    []uint64{},
				SymbolsNoFree:  []uint64{},
				SymbolsFree:    []uint64{},
			},
		},
		{
			name: "few, no free spins",
			results: []result{
				{bet: 100, win: 500},
				{bet: 100},
				{bet: 50, win: 0},
				{bet: 100, win: 50},
			},
			want: &Rounds{
				Count:          4,
				Bets:           &MinMaxInt64{Count: 4, Total: 350, Min: 50, Max: 100, Counts: map[int64]uint64{50: 1, 100: 3}},
				BetsNoFree:     &MinMaxInt64{Count: 4, Total: 350, Min: 50, Max: 100, Counts: map[int64]uint64{50: 1, 100: 3}},
				BetsFree:       &MinMaxInt64{Counts: map[int64]uint64{}, first: true},
				Wins:           &MinMaxInt64{Count: 4, Total: 550, Min: 50, Max: 500, Counts: map[int64]uint64{0: 2, 50: 1, 500: 1}},
				WinsNoFree:     &MinMaxInt64{Count: 4, Total: 550, Min: 50, Max: 500, Counts: map[int64]uint64{0: 2, 50: 1, 500: 1}},
				WinsFree:       &MinMaxInt64{Counts: map[int64]uint64{}, first: true},
				FreeSpins:      &MinMaxUInt64{Count: 4, Counts: map[uint64]uint64{0: 4}},
				RefillSpins:    &MinMaxUInt64{Count: 4, Counts: map[uint64]uint64{0: 4}},
				SuperSpins:     &MinMaxUInt64{Count: 4, Counts: map[uint64]uint64{0: 4}},
				FreeSpinRounds: []uint64{4},
				RefillRounds:   []uint64{4},
				SuperRounds:    []uint64{4},
				SymbolsUsed:    []uint64{},
				SymbolsNoFree:  []uint64{},
				SymbolsFree:    []uint64{},
			},
		},
		{
			name: "one, free spins",
			results: []result{
				{bet: 100, win: 1500, freeSpins: 10},
			},
			want: &Rounds{
				Count:          1,
				Bets:           &MinMaxInt64{Count: 1, Total: 100, Min: 100, Max: 100, Counts: map[int64]uint64{100: 1}},
				BetsNoFree:     &MinMaxInt64{Counts: map[int64]uint64{}, first: true},
				BetsFree:       &MinMaxInt64{Count: 1, Total: 100, Min: 100, Max: 100, Counts: map[int64]uint64{100: 1}},
				Wins:           &MinMaxInt64{Count: 1, Total: 1500, Min: 1500, Max: 1500, Counts: map[int64]uint64{1500: 1}},
				WinsNoFree:     &MinMaxInt64{Counts: map[int64]uint64{}, first: true},
				WinsFree:       &MinMaxInt64{Count: 1, Total: 1500, Min: 1500, Max: 1500, Counts: map[int64]uint64{1500: 1}},
				FreeSpins:      &MinMaxUInt64{Count: 1, Total: 10, Min: 10, Max: 10, Counts: map[uint64]uint64{10: 1}},
				RefillSpins:    &MinMaxUInt64{Count: 1, Counts: map[uint64]uint64{0: 1}},
				SuperSpins:     &MinMaxUInt64{Count: 1, Counts: map[uint64]uint64{0: 1}},
				FreeSpinRounds: []uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				RefillRounds:   []uint64{1},
				SuperRounds:    []uint64{1},
				SymbolsUsed:    []uint64{},
				SymbolsNoFree:  []uint64{},
				SymbolsFree:    []uint64{},
			},
		},
		{
			name: "few, free spins",
			results: []result{
				{bet: 100, win: 5000, freeSpins: 18},
				{bet: 100},
				{bet: 50},
				{bet: 100, win: 7550, freeSpins: 8},
			},
			want: &Rounds{
				Count:          4,
				Bets:           &MinMaxInt64{Count: 4, Total: 350, Min: 50, Max: 100, Counts: map[int64]uint64{50: 1, 100: 3}},
				BetsNoFree:     &MinMaxInt64{Count: 2, Total: 150, Min: 50, Max: 100, Counts: map[int64]uint64{50: 1, 100: 1}},
				BetsFree:       &MinMaxInt64{Count: 2, Total: 200, Min: 100, Max: 100, Counts: map[int64]uint64{100: 2}},
				Wins:           &MinMaxInt64{Count: 4, Total: 12550, Min: 5000, Max: 7550, Counts: map[int64]uint64{0: 2, 5000: 1, 7550: 1}},
				WinsNoFree:     &MinMaxInt64{Count: 2, Counts: map[int64]uint64{0: 2}},
				WinsFree:       &MinMaxInt64{Count: 2, Total: 12550, Min: 5000, Max: 7550, Counts: map[int64]uint64{5000: 1, 7550: 1}},
				FreeSpins:      &MinMaxUInt64{Count: 4, Total: 26, Min: 8, Max: 18, Counts: map[uint64]uint64{0: 2, 8: 1, 18: 1}},
				RefillSpins:    &MinMaxUInt64{Count: 4, Counts: map[uint64]uint64{0: 4}},
				SuperSpins:     &MinMaxUInt64{Count: 4, Counts: map[uint64]uint64{0: 4}},
				FreeSpinRounds: []uint64{2, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				RefillRounds:   []uint64{4},
				SuperRounds:    []uint64{4},
				SymbolsUsed:    []uint64{},
				SymbolsNoFree:  []uint64{},
				SymbolsFree:    []uint64{},
			},
		},
		{
			name: "few, refill spins",
			results: []result{
				{bet: 100, win: 4000, refillSpins: 13},
				{bet: 100},
				{bet: 50},
				{bet: 100, win: 6550, refillSpins: 7},
			},
			want: &Rounds{
				Count:          4,
				Bets:           &MinMaxInt64{Count: 4, Total: 350, Min: 50, Max: 100, Counts: map[int64]uint64{50: 1, 100: 3}},
				BetsNoFree:     &MinMaxInt64{Count: 4, Total: 350, Min: 50, Max: 100, Counts: map[int64]uint64{50: 1, 100: 3}},
				BetsFree:       &MinMaxInt64{Counts: map[int64]uint64{}, first: true},
				Wins:           &MinMaxInt64{Count: 4, Total: 10550, Min: 4000, Max: 6550, Counts: map[int64]uint64{0: 2, 4000: 1, 6550: 1}},
				WinsNoFree:     &MinMaxInt64{Count: 4, Total: 10550, Min: 4000, Max: 6550, Counts: map[int64]uint64{0: 2, 4000: 1, 6550: 1}},
				WinsFree:       &MinMaxInt64{Counts: map[int64]uint64{}, first: true},
				FreeSpins:      &MinMaxUInt64{Count: 4, Counts: map[uint64]uint64{0: 4}},
				RefillSpins:    &MinMaxUInt64{Count: 4, Total: 20, Min: 7, Max: 13, Counts: map[uint64]uint64{0: 2, 7: 1, 13: 1}},
				SuperSpins:     &MinMaxUInt64{Count: 4, Counts: map[uint64]uint64{0: 4}},
				FreeSpinRounds: []uint64{4},
				RefillRounds:   []uint64{2, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1},
				SuperRounds:    []uint64{4},
				SymbolsUsed:    []uint64{},
				SymbolsNoFree:  []uint64{},
				SymbolsFree:    []uint64{},
			},
		},
		{
			name: "few, super spins",
			results: []result{
				{bet: 100, win: 1000, superSpins: 11},
				{bet: 100},
				{bet: 50},
				{bet: 100, win: 700, superSpins: 4},
			},
			want: &Rounds{
				Count:          4,
				Bets:           &MinMaxInt64{Count: 4, Total: 350, Min: 50, Max: 100, Counts: map[int64]uint64{50: 1, 100: 3}},
				BetsNoFree:     &MinMaxInt64{Count: 4, Total: 350, Min: 50, Max: 100, Counts: map[int64]uint64{50: 1, 100: 3}},
				BetsFree:       &MinMaxInt64{Counts: map[int64]uint64{}, first: true},
				Wins:           &MinMaxInt64{Count: 4, Total: 1700, Min: 700, Max: 1000, Counts: map[int64]uint64{0: 2, 700: 1, 1000: 1}},
				WinsNoFree:     &MinMaxInt64{Count: 4, Total: 1700, Min: 700, Max: 1000, Counts: map[int64]uint64{0: 2, 700: 1, 1000: 1}},
				WinsFree:       &MinMaxInt64{Counts: map[int64]uint64{}, first: true},
				FreeSpins:      &MinMaxUInt64{Count: 4, Counts: map[uint64]uint64{0: 4}},
				RefillSpins:    &MinMaxUInt64{Count: 4, Counts: map[uint64]uint64{0: 4}},
				SuperSpins:     &MinMaxUInt64{Count: 4, Total: 15, Min: 4, Max: 11, Counts: map[uint64]uint64{0: 2, 4: 1, 11: 1}},
				FreeSpinRounds: []uint64{4},
				RefillRounds:   []uint64{4},
				SuperRounds:    []uint64{2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1},
				SymbolsUsed:    []uint64{},
				SymbolsNoFree:  []uint64{},
				SymbolsFree:    []uint64{},
			},
		},
		{
			name: "many",
			results: []result{
				{bet: 100, win: 5000, freeSpins: 18},
				{bet: 100},
				{bet: 50, refillSpins: 2},
				{bet: 100, win: 7000, freeSpins: 18},
				{bet: 100},
				{bet: 100, win: 7550, freeSpins: 8, refillSpins: 3},
				{bet: 100, win: 3500, freeSpins: 16},
				{bet: 100, superSpins: 5},
				{bet: 50, win: 2050, freeSpins: 8},
				{bet: 100, win: 50, superSpins: 3},
			},
			want: &Rounds{
				Count:          10,
				Bets:           &MinMaxInt64{Count: 10, Total: 900, Min: 50, Max: 100, Counts: map[int64]uint64{50: 2, 100: 8}},
				BetsNoFree:     &MinMaxInt64{Count: 5, Total: 450, Min: 50, Max: 100, Counts: map[int64]uint64{50: 1, 100: 4}},
				BetsFree:       &MinMaxInt64{Count: 5, Total: 450, Min: 50, Max: 100, Counts: map[int64]uint64{50: 1, 100: 4}},
				Wins:           &MinMaxInt64{Count: 10, Total: 25150, Min: 50, Max: 7550, Counts: map[int64]uint64{0: 4, 50: 1, 2050: 1, 3500: 1, 5000: 1, 7000: 1, 7550: 1}},
				WinsNoFree:     &MinMaxInt64{Count: 5, Total: 50, Min: 50, Max: 50, Counts: map[int64]uint64{0: 4, 50: 1}},
				WinsFree:       &MinMaxInt64{Count: 5, Total: 25100, Min: 2050, Max: 7550, Counts: map[int64]uint64{2050: 1, 3500: 1, 5000: 1, 7000: 1, 7550: 1}},
				FreeSpins:      &MinMaxUInt64{Count: 10, Total: 68, Min: 8, Max: 18, Counts: map[uint64]uint64{0: 5, 8: 2, 16: 1, 18: 2}},
				RefillSpins:    &MinMaxUInt64{Count: 10, Total: 5, Min: 2, Max: 3, Counts: map[uint64]uint64{0: 8, 2: 1, 3: 1}},
				SuperSpins:     &MinMaxUInt64{Count: 10, Total: 8, Min: 3, Max: 5, Counts: map[uint64]uint64{0: 8, 3: 1, 5: 1}},
				FreeSpinRounds: []uint64{5, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 1, 0, 2},
				RefillRounds:   []uint64{8, 0, 1, 1},
				SuperRounds:    []uint64{8, 0, 0, 1, 0, 1},
				SymbolsUsed:    []uint64{},
				SymbolsNoFree:  []uint64{},
				SymbolsFree:    []uint64{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRounds(nil)
			require.NotNil(t, r)
			defer r.Release()

			for _, res := range tc.results {
				r.NewRound(res.bet, res.win, res.freeSpins, res.refillSpins, res.superSpins, nil)
			}

			if !tc.want.Equals(r) {
				assert.EqualValues(t, tc.want, r)
			}

			n := r.Clone().(*Rounds)
			require.NotNil(t, n)
			defer n.Release()

			if !tc.want.Equals(n) {
				assert.EqualValues(t, tc.want, n)
			}

			n.ResetData()
			assert.Zero(t, n.Count)
			assert.Zero(t, n.Bets.Total)
			assert.Equal(t, int64(0), n.Bets.Min)
			assert.Equal(t, int64(0), n.Bets.Max)
			assert.Zero(t, len(n.Bets.Counts))
			assert.Zero(t, n.Wins.Total)
			assert.Equal(t, int64(0), n.Wins.Min)
			assert.Equal(t, int64(0), n.Wins.Max)
			assert.Zero(t, len(n.Wins.Counts))
			assert.Zero(t, n.FreeSpins.Total)
			assert.Equal(t, uint64(0), n.FreeSpins.Min)
			assert.Zero(t, n.FreeSpins.Max)
			assert.Zero(t, len(n.FreeSpins.Counts))
		})
	}
}
