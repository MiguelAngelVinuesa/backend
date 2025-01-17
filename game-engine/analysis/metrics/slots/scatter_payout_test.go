package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewScatterPayout(t *testing.T) {
	type increase struct {
		symbol utils.Index
		count  uint8
		payout float64
	}

	testCases := []struct {
		name      string
		maxSymbol utils.Index
		increases []increase
		want      *ScatterPayout
	}{
		{
			name:      "empty",
			maxSymbol: 11,
			want: &ScatterPayout{
				Symbols:       []uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Lengths:       []uint64{},
				SymbolLengths: [][]uint64{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}},
				Payouts:       AcquireMinMaxFloat64(2),
			},
		},
		{
			name:      "single",
			maxSymbol: 11,
			increases: []increase{
				{1, 3, 5},
			},
			want: &ScatterPayout{
				Count:         1,
				Symbols:       []uint64{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Lengths:       []uint64{0, 0, 0, 1},
				SymbolLengths: [][]uint64{{}, {0, 0, 0, 1}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}},
				Payouts: &MinMaxFloat64{
					Count:    1,
					Total:    5,
					Min:      5,
					Max:      5,
					Counts:   map[int64]uint64{500: 1},
					Decimals: 2,
					Factor:   100,
				},
			},
		},
		{
			name:      "three without dups",
			maxSymbol: 11,
			increases: []increase{
				{2, 5, 50},
				{7, 2, 10},
				{4, 4, 20},
			},
			want: &ScatterPayout{
				Count:         3,
				Symbols:       []uint64{0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 0, 0},
				Lengths:       []uint64{0, 0, 1, 0, 1, 1},
				SymbolLengths: [][]uint64{{}, {}, {0, 0, 0, 0, 0, 1}, {}, {0, 0, 0, 0, 1}, {}, {}, {0, 0, 1}, {}, {}, {}, {}},
				Payouts: &MinMaxFloat64{
					Count:    3,
					Total:    80.0,
					Min:      10.0,
					Max:      50.0,
					Counts:   map[int64]uint64{1000: 1, 2000: 1, 5000: 1},
					Decimals: 2,
					Factor:   100,
				},
			},
		},
		{
			name:      "many with dups",
			maxSymbol: 11,
			increases: []increase{
				{7, 5, 500},
				{2, 5, 50},
				{7, 2, 10},
				{4, 5, 200},
				{2, 3, 20},
				{7, 3, 100},
				{4, 4, 20},
			},
			want: &ScatterPayout{
				Count:         7,
				Symbols:       []uint64{0, 0, 2, 0, 2, 0, 0, 3, 0, 0, 0, 0},
				Lengths:       []uint64{0, 0, 1, 2, 1, 3},
				SymbolLengths: [][]uint64{{}, {}, {0, 0, 0, 1, 0, 1}, {}, {0, 0, 0, 0, 1, 1}, {}, {}, {0, 0, 1, 1, 0, 1}, {}, {}, {}, {}},
				Payouts: &MinMaxFloat64{
					Count:    7,
					Total:    900.0,
					Min:      10.0,
					Max:      500.0,
					Counts:   map[int64]uint64{1000: 1, 2000: 2, 5000: 1, 10000: 1, 20000: 1, 50000: 1},
					Decimals: 2,
					Factor:   100,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewScatterPayout(tc.maxSymbol)
			require.NotNil(t, s)
			defer s.Release()

			for _, i := range tc.increases {
				s.Increase(i.symbol, i.count, i.payout)
			}

			if !tc.want.Equals(s) {
				assert.EqualValues(t, tc.want, s)
			}

			n := s.Clone().(*ScatterPayout)
			require.NotNil(t, n)
			defer n.Release()

			if !tc.want.Equals(n) {
				assert.EqualValues(t, tc.want, n)
			}

			n.ResetData()
			assert.Zero(t, n.Count)
			defer n.Release()

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
