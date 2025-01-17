package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewPayline(t *testing.T) {
	type increase struct {
		symbol utils.Index
		count  uint8
		payout float64
	}

	testCases := []struct {
		name      string
		id        int
		maxSymbol utils.Index
		rowMap    utils.UInt8s
		increases []increase
		want      *Payline
	}{
		{
			name:      "empty",
			id:        1,
			maxSymbol: 11,
			rowMap:    utils.UInt8s{1, 1, 1, 1, 1},
			want: &Payline{
				ID:            1,
				RowMap:        utils.UInt8s{1, 1, 1, 1, 1},
				Symbols:       []uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				SymbolLengths: [][]uint64{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}},
				Lengths:       []uint64{},
				Payouts:       &MinMaxFloat64{Counts: map[int64]uint64{}, first: true, Decimals: 1, Factor: 10},
			},
		},
		{
			name:      "one",
			id:        2,
			maxSymbol: 11,
			rowMap:    utils.UInt8s{0, 0, 0, 0, 0},
			increases: []increase{
				{4, 3, 0.5},
			},
			want: &Payline{
				ID:            2,
				RowMap:        utils.UInt8s{0, 0, 0, 0, 0},
				Count:         1,
				Symbols:       []uint64{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
				SymbolLengths: [][]uint64{{}, {}, {}, {}, {0, 0, 0, 1}, {}, {}, {}, {}, {}, {}, {}},
				Lengths:       []uint64{0, 0, 0, 1},
				Payouts:       &MinMaxFloat64{Count: 1, Total: 0.5, Min: 0.5, Max: 0.5, Counts: map[int64]uint64{5: 1}, Decimals: 1, Factor: 10},
			},
		},
		{
			name:      "few, no dups",
			id:        9,
			maxSymbol: 11,
			rowMap:    utils.UInt8s{2, 1, 1, 1, 2},
			increases: []increase{
				{4, 4, 2.5},
				{2, 5, 20},
				{7, 2, 0.5},
			},
			want: &Payline{
				ID:            9,
				RowMap:        utils.UInt8s{2, 1, 1, 1, 2},
				Count:         3,
				Symbols:       []uint64{0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 0, 0},
				SymbolLengths: [][]uint64{{}, {}, {0, 0, 0, 0, 0, 1}, {}, {0, 0, 0, 0, 1}, {}, {}, {0, 0, 1}, {}, {}, {}, {}},
				Lengths:       []uint64{0, 0, 1, 0, 1, 1},
				Payouts:       &MinMaxFloat64{Count: 3, Total: 23.0, Min: 0.5, Max: 20.0, Counts: map[int64]uint64{5: 1, 25: 1, 200: 1}, Decimals: 1, Factor: 10},
			},
		},
		{
			name:      "many, dups",
			id:        9,
			maxSymbol: 11,
			rowMap:    utils.UInt8s{2, 1, 1, 1, 2},
			increases: []increase{
				{4, 4, 2.5},
				{2, 5, 20},
				{7, 2, 0.5},
				{2, 4, 2.5},
				{4, 5, 20},
				{1, 2, 0.5},
				{7, 4, 2.5},
				{1, 5, 20},
				{4, 2, 0.5},
				{1, 4, 2.5},
				{2, 5, 20},
				{7, 2, 0.5},
			},
			want: &Payline{
				ID:            9,
				RowMap:        utils.UInt8s{2, 1, 1, 1, 2},
				Count:         12,
				Symbols:       []uint64{0, 3, 3, 0, 3, 0, 0, 3, 0, 0, 0, 0},
				SymbolLengths: [][]uint64{{}, {0, 0, 1, 0, 1, 1}, {0, 0, 0, 0, 1, 2}, {}, {0, 0, 1, 0, 1, 1}, {}, {}, {0, 0, 2, 0, 1}, {}, {}, {}, {}},
				Lengths:       []uint64{0, 0, 4, 0, 4, 4},
				Payouts:       &MinMaxFloat64{Count: 12, Total: 92.0, Min: 0.5, Max: 20.0, Counts: map[int64]uint64{5: 4, 25: 4, 200: 4}, Decimals: 1, Factor: 10},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPayline(tc.id, tc.maxSymbol, tc.rowMap)
			require.NotNil(t, p)
			defer p.Release()

			for _, i := range tc.increases {
				p.Increase(i.symbol, i.count, i.payout)
			}

			if !tc.want.Equals(p) {
				assert.EqualValues(t, tc.want, p)
			}

			n := p.Clone().(*Payline)
			require.NotNil(t, n)
			defer n.Release()

			if !tc.want.Equals(n) {
				assert.EqualValues(t, tc.want, n)
			}

			n.ResetData()
			assert.Zero(t, n.Count)
			assert.Zero(t, len(n.Lengths))

			for _, v := range n.Symbols {
				assert.Zero(t, v)
			}

			assert.Zero(t, n.Payouts.Total)
			assert.Zero(t, n.Payouts.Min)
			assert.Zero(t, n.Payouts.Max)
			assert.Zero(t, len(n.Payouts.Counts))
		})
	}
}
