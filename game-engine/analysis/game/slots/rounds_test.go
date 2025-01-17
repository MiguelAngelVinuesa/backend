package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	analyse "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/analysis/metrics/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewRounds(t *testing.T) {
	testCases := []struct {
		name      string
		playerID  string
		reelCount int
		rowCount  int
		bet       int64
		results   results.Results
		bets      *analyse.MinMaxInt64
		wins      *analyse.MinMaxInt64
		symbols   []*analyse.Symbol
		best      []results.Results
		maxPayout float64
	}{
		{
			name:      "no rounds",
			playerID:  "jimmy",
			reelCount: 5,
			rowCount:  3,
			bets:      analyse.AcquireMinMaxInt64(),
			wins:      analyse.AcquireMinMaxInt64(),
		},
		{
			name:      "1 round, no payouts",
			playerID:  "jimmy",
			reelCount: 5,
			rowCount:  3,
			bet:       100,
			results: results.Results{
				results.AcquireResult(
					slots.AcquireSpinResultFromData(utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3}, nil, nil, nil, nil, 0, 0, 0), results.SpinData,
				),
			},
			bets: &analyse.MinMaxInt64{Count: 1, Total: 100, Min: 100, Max: 100, Counts: map[int64]uint64{100: 1}},
			wins: &analyse.MinMaxInt64{Count: 1, Counts: map[int64]uint64{0: 1}},
			symbols: []*analyse.Symbol{nil,
				{ID: 1, Name: "A", TotalCount: 3, TotalReels: []uint64{1, 0, 1, 0, 1}, FirstCount: 3, FirstReels: []uint64{1, 0, 1, 0, 1}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 2, Name: "B", TotalCount: 3, TotalReels: []uint64{1, 0, 1, 0, 1}, FirstCount: 3, FirstReels: []uint64{1, 0, 1, 0, 1}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 3, Name: "C", TotalCount: 3, TotalReels: []uint64{1, 0, 1, 0, 1}, FirstCount: 3, FirstReels: []uint64{1, 0, 1, 0, 1}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 4, Name: "D", TotalCount: 2, TotalReels: []uint64{0, 1, 0, 1, 0}, FirstCount: 2, FirstReels: []uint64{0, 1, 0, 1, 0}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 5, Name: "E", TotalCount: 2, TotalReels: []uint64{0, 1, 0, 1, 0}, FirstCount: 2, FirstReels: []uint64{0, 1, 0, 1, 0}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 6, Name: "F", TotalCount: 2, TotalReels: []uint64{0, 1, 0, 1, 0}, FirstCount: 2, FirstReels: []uint64{0, 1, 0, 1, 0}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
			},
		},
		{
			name:      "1 round, 1 payout",
			playerID:  "jimmy",
			reelCount: 5,
			rowCount:  3,
			bet:       100,
			results: results.Results{
				results.AcquireResult(
					slots.AcquireSpinResultFromData(utils.Indexes{1, 2, 3, 4, 5, 3, 1, 2, 3, 4, 5, 3, 1, 2, 6}, nil, nil, nil, nil, 0, 0, 0),
					results.SpinData,
					slots.WinlinePayoutFromData(2.5, 0, 3, 4, slots.PayLTR, 2, utils.UInt8s{2, 2, 2, 2, 2}),
				),
			},
			bets: &analyse.MinMaxInt64{Count: 1, Total: 100, Min: 100, Max: 100, Counts: map[int64]uint64{100: 1}},
			wins: &analyse.MinMaxInt64{Count: 1, Total: 250, Min: 250, Max: 250, Counts: map[int64]uint64{250: 1}},
			symbols: []*analyse.Symbol{nil,
				{ID: 1, Name: "A", TotalCount: 3, TotalReels: []uint64{1, 0, 1, 0, 1}, FirstCount: 3, FirstReels: []uint64{1, 0, 1, 0, 1}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 2, Name: "B", TotalCount: 3, TotalReels: []uint64{1, 0, 1, 0, 1}, FirstCount: 3, FirstReels: []uint64{1, 0, 1, 0, 1}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 3, Name: "C", TotalCount: 4, TotalReels: []uint64{1, 1, 1, 1, 0}, FirstCount: 4, FirstReels: []uint64{1, 1, 1, 1, 0}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{0, 0, 1}},
				{ID: 4, Name: "D", TotalCount: 2, TotalReels: []uint64{0, 1, 0, 1, 0}, FirstCount: 2, FirstReels: []uint64{0, 1, 0, 1, 0}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 5, Name: "E", TotalCount: 2, TotalReels: []uint64{0, 1, 0, 1, 0}, FirstCount: 2, FirstReels: []uint64{0, 1, 0, 1, 0}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 6, Name: "F", TotalCount: 1, TotalReels: []uint64{0, 0, 0, 0, 1}, FirstCount: 1, FirstReels: []uint64{0, 0, 0, 0, 1}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
			},
		},
		{
			name:      "1 round, 3 payouts",
			playerID:  "jimmy",
			reelCount: 5,
			rowCount:  3,
			bet:       100,
			results: results.Results{
				results.AcquireResult(
					slots.AcquireSpinResultFromData(utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 4, 2, 3, 1, 2, 6}, nil, nil, nil, nil, 0, 0, 0),
					results.SpinData,
					slots.WinlinePayoutFromData(25, 0, 2, 5, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
					slots.WinlinePayoutFromData(0.5, 0, 1, 3, slots.PayLTR, 2, utils.UInt8s{0, 0, 0, 0, 0}),
					slots.WinlinePayoutFromData(0.5, 5, 3, 4, slots.PayLTR, 3, utils.UInt8s{2, 2, 2, 2, 2}),
				),
			},
			bets: &analyse.MinMaxInt64{Count: 1, Total: 100, Min: 100, Max: 100, Counts: map[int64]uint64{100: 1}},
			wins: &analyse.MinMaxInt64{Count: 1, Total: 2800, Min: 2800, Max: 2800, Counts: map[int64]uint64{2800: 1}},
			symbols: []*analyse.Symbol{nil,
				{ID: 1, Name: "A", TotalCount: 4, TotalReels: []uint64{1, 1, 1, 0, 1}, FirstCount: 4, FirstReels: []uint64{1, 1, 1, 0, 1}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{0, 1}},
				{ID: 2, Name: "B", TotalCount: 5, TotalReels: []uint64{1, 1, 1, 1, 1}, FirstCount: 5, FirstReels: []uint64{1, 1, 1, 1, 1}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{0, 0, 0, 1}},
				{ID: 3, Name: "C", TotalCount: 4, TotalReels: []uint64{1, 1, 1, 1, 0}, FirstCount: 4, FirstReels: []uint64{1, 1, 1, 1, 0}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{0, 0, 1}},
				{ID: 4, Name: "D", TotalCount: 1, TotalReels: []uint64{0, 0, 0, 1, 0}, FirstCount: 1, FirstReels: []uint64{0, 0, 0, 1, 0}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 5, Name: "E", TotalCount: 0, TotalReels: []uint64{0, 0, 0, 0, 0}, FirstCount: 0, FirstReels: []uint64{0, 0, 0, 0, 0}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 6, Name: "F", TotalCount: 1, TotalReels: []uint64{0, 0, 0, 0, 1}, FirstCount: 1, FirstReels: []uint64{0, 0, 0, 0, 1}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
			},
		},
		{
			name:      "3 rounds, 4 payouts",
			playerID:  "jimmy",
			reelCount: 5,
			rowCount:  3,
			bet:       100,
			results: results.Results{
				results.AcquireResult(
					slots.AcquireSpinResultFromData(utils.Indexes{1, 2, 3, 1, 4, 3, 1, 2, 3, 4, 2, 3, 1, 2, 6}, nil, nil, nil, nil, 0, 0, 0),
					results.SpinData,
					slots.WinlinePayoutFromData(25, 0, 2, 5, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
					slots.WinlinePayoutFromData(2.5, 0, 3, 4, slots.PayLTR, 3, utils.UInt8s{2, 2, 2, 2, 2}),
				),
				results.AcquireResult(
					slots.AcquireSpinResultFromData(utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 6}, nil, nil, nil, nil, 0, 0, 0),
					results.SpinData,
				),
				results.AcquireResult(
					slots.AcquireSpinResultFromData(utils.Indexes{1, 2, 3, 1, 2, 6, 1, 2, 3, 4, 2, 3, 1, 2, 6}, nil, nil, nil, nil, 0, 0, 0),
					results.SpinData,
					slots.WinlinePayoutFromData(25, 0, 2, 5, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
					slots.WinlinePayoutFromData(0.5, 0, 1, 3, slots.PayLTR, 2, utils.UInt8s{0, 0, 0, 0, 0}),
				),
			},
			bets: &analyse.MinMaxInt64{Count: 3, Total: 300, Min: 100, Max: 100, Counts: map[int64]uint64{100: 3}},
			wins: &analyse.MinMaxInt64{Count: 3, Total: 5300, Min: 2550, Max: 2750, Counts: map[int64]uint64{0: 1, 2550: 1, 2750: 1}},
			symbols: []*analyse.Symbol{nil,
				{ID: 1, Name: "A", TotalCount: 11, TotalReels: []uint64{3, 2, 3, 0, 3}, FirstCount: 11, FirstReels: []uint64{3, 2, 3, 0, 3}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{0, 1}},
				{ID: 2, Name: "B", TotalCount: 12, TotalReels: []uint64{3, 1, 3, 2, 3}, FirstCount: 12, FirstReels: []uint64{3, 1, 3, 2, 3}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{0, 0, 0, 2}},
				{ID: 3, Name: "C", TotalCount: 9, TotalReels: []uint64{3, 1, 3, 2, 0}, FirstCount: 9, FirstReels: []uint64{3, 1, 3, 2, 0}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{0, 0, 1}},
				{ID: 4, Name: "D", TotalCount: 5, TotalReels: []uint64{0, 2, 0, 3, 0}, FirstCount: 5, FirstReels: []uint64{0, 2, 0, 3, 0}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 5, Name: "E", TotalCount: 2, TotalReels: []uint64{0, 1, 0, 1, 0}, FirstCount: 2, FirstReels: []uint64{0, 1, 0, 1, 0}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 6, Name: "F", TotalCount: 6, TotalReels: []uint64{0, 2, 0, 1, 3}, FirstCount: 6, FirstReels: []uint64{0, 2, 0, 1, 3}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
			},
		},
		{
			name:      "1 round, 1 payout, best",
			playerID:  "jimmy",
			reelCount: 5,
			rowCount:  3,
			bet:       100,
			results: results.Results{
				results.AcquireResult(
					slots.AcquireSpinResultFromData(utils.Indexes{1, 2, 6, 4, 5, 6, 1, 2, 6, 4, 5, 6, 1, 2, 6}, nil, nil, nil, nil, 0, 0, 0),
					results.SpinData,
					slots.WinlinePayoutFromData(10000.0, 2, 6, 5, slots.PayLTR, 2, utils.UInt8s{2, 2, 2, 2, 2}),
				),
			},
			bets: &analyse.MinMaxInt64{Count: 1, Total: 100, Min: 100, Max: 100, Counts: map[int64]uint64{100: 1}},
			wins: &analyse.MinMaxInt64{Count: 1, Total: 2000000, Min: 2000000, Max: 2000000, Counts: map[int64]uint64{2000000: 1}},
			symbols: []*analyse.Symbol{nil,
				{ID: 1, Name: "A", TotalCount: 3, TotalReels: []uint64{1, 0, 1, 0, 1}, FirstCount: 3, FirstReels: []uint64{1, 0, 1, 0, 1}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 2, Name: "B", TotalCount: 3, TotalReels: []uint64{1, 0, 1, 0, 1}, FirstCount: 3, FirstReels: []uint64{1, 0, 1, 0, 1}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 3, Name: "C", TotalCount: 0, TotalReels: []uint64{0, 0, 0, 0, 0}, FirstCount: 0, FirstReels: []uint64{0, 0, 0, 0, 0}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 4, Name: "D", TotalCount: 2, TotalReels: []uint64{0, 1, 0, 1, 0}, FirstCount: 2, FirstReels: []uint64{0, 1, 0, 1, 0}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 5, Name: "E", TotalCount: 2, TotalReels: []uint64{0, 1, 0, 1, 0}, FirstCount: 2, FirstReels: []uint64{0, 1, 0, 1, 0}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 6, Name: "F", TotalCount: 5, TotalReels: []uint64{1, 1, 1, 1, 1}, FirstCount: 5, FirstReels: []uint64{1, 1, 1, 1, 1}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{0, 0, 0, 1}},
			},
			best: []results.Results{
				{
					results.AcquireResult(
						slots.AcquireSpinResultFromData(utils.Indexes{1, 2, 6, 4, 5, 6, 1, 2, 6, 4, 5, 6, 1, 2, 6}, nil, nil, nil, nil, 0, 0, 0),
						results.SpinData,
						slots.WinlinePayoutFromData(10000, 2, 6, 5, slots.PayLTR, 2, utils.UInt8s{2, 2, 2, 2, 2}),
					),
				},
			},
		},
		{
			name:      "1 round, 1 payout, almost best",
			playerID:  "jimmy",
			reelCount: 5,
			rowCount:  3,
			bet:       100,
			results: results.Results{
				results.AcquireResult(
					slots.AcquireSpinResultFromData(utils.Indexes{1, 2, 6, 4, 5, 6, 1, 2, 6, 4, 5, 6, 1, 2, 6}, nil, nil, nil, nil, 0, 0, 0),
					results.SpinData,
					slots.WinlinePayoutFromData(199.9, 0, 6, 5, slots.PayLTR, 2, utils.UInt8s{2, 2, 2, 2, 2}),
				),
			},
			bets: &analyse.MinMaxInt64{Count: 1, Total: 100, Min: 100, Max: 100, Counts: map[int64]uint64{100: 1}},
			wins: &analyse.MinMaxInt64{Count: 1, Total: 19990, Min: 19990, Max: 19990, Counts: map[int64]uint64{19990: 1}},
			symbols: []*analyse.Symbol{nil,
				{ID: 1, Name: "A", TotalCount: 3, TotalReels: []uint64{1, 0, 1, 0, 1}, FirstCount: 3, FirstReels: []uint64{1, 0, 1, 0, 1}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 2, Name: "B", TotalCount: 3, TotalReels: []uint64{1, 0, 1, 0, 1}, FirstCount: 3, FirstReels: []uint64{1, 0, 1, 0, 1}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 3, Name: "C", TotalCount: 0, TotalReels: []uint64{0, 0, 0, 0, 0}, FirstCount: 0, FirstReels: []uint64{0, 0, 0, 0, 0}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 4, Name: "D", TotalCount: 2, TotalReels: []uint64{0, 1, 0, 1, 0}, FirstCount: 2, FirstReels: []uint64{0, 1, 0, 1, 0}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 5, Name: "E", TotalCount: 2, TotalReels: []uint64{0, 1, 0, 1, 0}, FirstCount: 2, FirstReels: []uint64{0, 1, 0, 1, 0}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{}},
				{ID: 6, Name: "F", TotalCount: 5, TotalReels: []uint64{1, 1, 1, 1, 1}, FirstCount: 5, FirstReels: []uint64{1, 1, 1, 1, 1}, SecondReels: []uint64{0, 0, 0, 0, 0}, FreeReels: []uint64{0, 0, 0, 0, 0}, FreeSecondReels: []uint64{0, 0, 0, 0, 0}, Payouts: []uint64{0, 0, 0, 1}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := AcquireRounds(0, tc.playerID, 100000000, tc.reelCount, tc.rowCount, false, 100000, set1, nil, paylines, nil)
			require.NotNil(t, r)
			defer r.Release()

			assert.Equal(t, tc.playerID, r.PlayerID)
			assert.Equal(t, uint64(0), r.RoundCount)
			assert.NotNil(t, r.AllRounds)
			assert.NotNil(t, r.Symbols)
			assert.NotNil(t, r.FirstPayouts)
			assert.NotNil(t, r.Best)
			assert.Equal(t, tc.reelCount, r.reelCount)
			assert.Equal(t, tc.rowCount, r.rowCount)
			assert.Equal(t, tc.rowCount, r.FirstPayouts.rowCount)

			for _, res := range tc.results {
				r.Analyse(tc.bet, tc.bet, results.Results{res})
			}

			assert.Equal(t, uint64(len(tc.results)), r.RoundCount)
			if tc.bets != nil {
				if !tc.bets.Equals(r.AllRounds.Bets) {
					assert.EqualValues(t, tc.bets, r.AllRounds.Bets)
				}
			}
			if tc.wins != nil {
				if !tc.wins.Equals(r.AllRounds.Wins) {
					assert.EqualValues(t, tc.wins, r.AllRounds.Wins)
				}
			}
			if tc.symbols != nil {
				bad := len(tc.symbols) != len(r.Symbols)
				if !bad {
					for ix := range tc.symbols {
						if tc.symbols[ix] != nil {
							if !tc.symbols[ix].Equals(r.Symbols[ix]) {
								bad = true
							}
						} else {
							if r.Symbols[ix] != nil {
								bad = true
							}
						}
					}
				}
				if bad {
					assert.EqualValues(t, tc.symbols, r.Symbols)
				}
			}
			if tc.best != nil {
				assert.EqualValues(t, len(tc.best), len(r.Best))
			}

			n := r.Clone().(*Rounds)
			require.NotNil(t, n)
			defer n.Release()

			assert.Equal(t, uint64(len(tc.results)), n.RoundCount)
			if tc.bets != nil {
				if !tc.bets.Equals(n.AllRounds.Bets) {
					assert.EqualValues(t, tc.bets, n.AllRounds.Bets)
				}
			}
			if tc.wins != nil {
				if !tc.wins.Equals(n.AllRounds.Wins) {
					assert.EqualValues(t, tc.wins, n.AllRounds.Wins)
				}
			}
			if tc.symbols != nil {
				bad := len(tc.symbols) != len(n.Symbols)
				if !bad {
					for ix := range tc.symbols {
						if tc.symbols[ix] != nil {
							if !tc.symbols[ix].Equals(n.Symbols[ix]) {
								bad = true
							}
						} else {
							if n.Symbols[ix] != nil {
								bad = true
							}
						}
					}
				}
				if bad {
					assert.EqualValues(t, tc.symbols, n.Symbols)
				}
			}
			if tc.best != nil {
				assert.EqualValues(t, len(tc.best), len(n.Best))
			}

			n.ResetData()
			assert.Zero(t, n.RoundCount)
			assert.Equal(t, r.reelCount, n.reelCount)
			assert.Equal(t, r.rowCount, n.rowCount)
		})
	}
}

func TestRounds_BestSorted(t *testing.T) {
	testCases := []struct {
		name    string
		results []results.Results
		want    []results.Results
	}{
		{
			name:    "single",
			results: []results.Results{{r1}},
			want:    []results.Results{{r1}},
		},
		{
			name:    "few ordered",
			results: []results.Results{{r3}, {r2}, {r1}},
			want:    []results.Results{{r3}, {r2}, {r1}},
		},
		{
			name:    "few unordered (1)",
			results: []results.Results{{r1}, {r3}, {r2}},
			want:    []results.Results{{r3}, {r2}, {r1}},
		},
		{
			name:    "few unordered (2)",
			results: []results.Results{{r2}, {r1}, {r3}},
			want:    []results.Results{{r3}, {r2}, {r1}},
		},
		{
			name:    "many ordered",
			results: []results.Results{{r6}, {r5}, {r4}, {r3}, {r2}, {r1}},
			want:    []results.Results{{r6}, {r5}, {r4}, {r3}, {r2}, {r1}},
		},
		{
			name:    "many unordered (1)",
			results: []results.Results{{r1}, {r6}, {r3}, {r5}, {r2}, {r4}},
			want:    []results.Results{{r6}, {r5}, {r4}, {r3}, {r2}, {r1}},
		},
		{
			name:    "many unordered (2)",
			results: []results.Results{{r4}, {r2}, {r5}, {r1}, {r3}, {r6}},
			want:    []results.Results{{r6}, {r5}, {r4}, {r3}, {r2}, {r1}},
		},
		{
			name:    "exact ordered",
			results: []results.Results{{r7}, {r7}, {r7}, {r7}, {r6}, {r6}, {r6}, {r5}, {r5}, {r5}, {r5}, {r4}, {r4}, {r4}, {r4}, {r3}, {r3}, {r3}, {r2}, {r2}, {r2}, {r2}, {r1}, {r1}, {r1}},
			want:    []results.Results{{r7}, {r7}, {r7}, {r7}, {r6}, {r6}, {r6}, {r5}, {r5}, {r5}, {r5}, {r4}, {r4}, {r4}, {r4}, {r3}, {r3}, {r3}, {r2}, {r2}, {r2}, {r2}, {r1}, {r1}, {r1}},
		},
		{
			name:    "exact unordered (1)",
			results: []results.Results{{r1}, {r1}, {r1}, {r2}, {r2}, {r2}, {r2}, {r3}, {r3}, {r3}, {r4}, {r4}, {r4}, {r4}, {r5}, {r5}, {r5}, {r5}, {r6}, {r6}, {r6}, {r7}, {r7}, {r7}, {r7}},
			want:    []results.Results{{r7}, {r7}, {r7}, {r7}, {r6}, {r6}, {r6}, {r5}, {r5}, {r5}, {r5}, {r4}, {r4}, {r4}, {r4}, {r3}, {r3}, {r3}, {r2}, {r2}, {r2}, {r2}, {r1}, {r1}, {r1}},
		},
		{
			name:    "exact unordered (2)",
			results: []results.Results{{r1}, {r3}, {r6}, {r7}, {r3}, {r4}, {r2}, {r6}, {r1}, {r2}, {r7}, {r6}, {r7}, {r2}, {r4}, {r1}, {r7}, {r3}, {r4}, {r3}, {r2}, {r7}, {r6}, {r4}, {r1}},
			want:    []results.Results{{r7}, {r7}, {r7}, {r7}, {r7}, {r6}, {r6}, {r6}, {r6}, {r4}, {r4}, {r4}, {r4}, {r3}, {r3}, {r3}, {r3}, {r2}, {r2}, {r2}, {r2}, {r1}, {r1}, {r1}, {r1}},
		},
		{
			name:    "overflow ordered",
			results: []results.Results{{r7}, {r7}, {r7}, {r7}, {r6}, {r6}, {r6}, {r6}, {r5}, {r5}, {r5}, {r5}, {r4}, {r4}, {r4}, {r4}, {r3}, {r3}, {r3}, {r3}, {r2}, {r2}, {r2}, {r2}, {r1}, {r1}, {r1}, {r1}},
			want:    []results.Results{{r7}, {r7}, {r7}, {r7}, {r6}, {r6}, {r6}, {r6}, {r5}, {r5}, {r5}, {r5}, {r4}, {r4}, {r4}, {r4}, {r3}, {r3}, {r3}, {r3}, {r2}, {r2}, {r2}, {r2}, {r1}},
		},
		{
			name:    "overflow unordered (1)",
			results: []results.Results{{r1}, {r1}, {r1}, {r1}, {r2}, {r2}, {r2}, {r2}, {r3}, {r3}, {r3}, {r3}, {r4}, {r4}, {r4}, {r4}, {r5}, {r5}, {r5}, {r5}, {r6}, {r6}, {r6}, {r6}, {r7}, {r7}, {r7}, {r7}},
			want:    []results.Results{{r7}, {r7}, {r7}, {r7}, {r6}, {r6}, {r6}, {r6}, {r5}, {r5}, {r5}, {r5}, {r4}, {r4}, {r4}, {r4}, {r3}, {r3}, {r3}, {r3}, {r2}, {r2}, {r2}, {r2}, {r1}},
		},
		{
			name:    "overflow unordered (2)",
			results: []results.Results{{r1}, {r3}, {r6}, {r7}, {r3}, {r4}, {r2}, {r6}, {r1}, {r2}, {r7}, {r5}, {r6}, {r7}, {r2}, {r4}, {r1}, {r6}, {r7}, {r3}, {r4}, {r3}, {r2}, {r7}, {r6}, {r3}, {r4}, {r1}},
			want:    []results.Results{{r7}, {r7}, {r7}, {r7}, {r7}, {r6}, {r6}, {r6}, {r6}, {r6}, {r5}, {r4}, {r4}, {r4}, {r4}, {r3}, {r3}, {r3}, {r3}, {r3}, {r2}, {r2}, {r2}, {r2}, {r1}},
		},
	}

	r := AcquireRounds(0, "x", 1000000, 5, 3, false, 10000, set1, nil, paylines, nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			best := make([]results.Results, 0, 8)
			for ix := range tc.results {
				best = r.addBestX(best, tc.results[ix], 0)
			}

			require.Equal(t, len(tc.want), len(best))

			for ix := range tc.want {
				gt1 := results.GrandTotal(tc.want[ix])
				gt2 := results.GrandTotal(best[ix])
				assert.Equal(t, gt1, gt2)
			}
		})
	}
}

var (
	s1   = slots.NewSymbol(1, slots.WithName("A"))
	s2   = slots.NewSymbol(2, slots.WithName("B"))
	s3   = slots.NewSymbol(3, slots.WithName("C"))
	s4   = slots.NewSymbol(4, slots.WithName("D"))
	s5   = slots.NewSymbol(5, slots.WithName("E"))
	s6   = slots.NewSymbol(6, slots.WithName("F"))
	set1 = slots.NewSymbolSet(s1, s2, s3, s4, s5, s6)

	l1       = slots.NewPayline(1, 3, 1, 1, 1, 1, 1)
	l2       = slots.NewPayline(2, 3, 0, 0, 0, 0, 0)
	l3       = slots.NewPayline(3, 3, 2, 2, 2, 2, 2)
	paylines = slots.Paylines{l1, l2, l3}

	r1 = results.AcquireResult(
		slots.AcquireSpinResultFromData(utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}, nil, nil, nil, nil, 0, 0, 0),
		results.SpinData,
		slots.WinlinePayoutFromData(50, 0, 1, 5, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
		slots.WinlinePayoutFromData(75, 0, 2, 5, slots.PayLTR, 2, utils.UInt8s{0, 0, 0, 0, 0}),
		slots.WinlinePayoutFromData(100, 0, 3, 5, slots.PayLTR, 3, utils.UInt8s{2, 2, 2, 2, 2}),
	)
	r2 = results.AcquireResult(
		slots.AcquireSpinResultFromData(utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}, nil, nil, nil, nil, 0, 0, 0),
		results.SpinData,
		slots.WinlinePayoutFromData(60, 0, 1, 5, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
		slots.WinlinePayoutFromData(85, 0, 2, 5, slots.PayLTR, 2, utils.UInt8s{0, 0, 0, 0, 0}),
		slots.WinlinePayoutFromData(110, 0, 3, 5, slots.PayLTR, 3, utils.UInt8s{2, 2, 2, 2, 2}),
	)
	r3 = results.AcquireResult(
		slots.AcquireSpinResultFromData(utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}, nil, nil, nil, nil, 0, 0, 0),
		results.SpinData,
		slots.WinlinePayoutFromData(70, 0, 1, 5, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
		slots.WinlinePayoutFromData(95, 0, 2, 5, slots.PayLTR, 2, utils.UInt8s{0, 0, 0, 0, 0}),
		slots.WinlinePayoutFromData(120, 0, 3, 5, slots.PayLTR, 3, utils.UInt8s{2, 2, 2, 2, 2}),
	)
	r4 = results.AcquireResult(
		slots.AcquireSpinResultFromData(utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}, nil, nil, nil, nil, 0, 0, 0),
		results.SpinData,
		slots.WinlinePayoutFromData(80, 0, 1, 5, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
		slots.WinlinePayoutFromData(105, 0, 2, 5, slots.PayLTR, 2, utils.UInt8s{0, 0, 0, 0, 0}),
		slots.WinlinePayoutFromData(130, 0, 3, 5, slots.PayLTR, 3, utils.UInt8s{2, 2, 2, 2, 2}),
	)
	r5 = results.AcquireResult(
		slots.AcquireSpinResultFromData(utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}, nil, nil, nil, nil, 0, 0, 0),
		results.SpinData,
		slots.WinlinePayoutFromData(90, 0, 1, 5, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
		slots.WinlinePayoutFromData(115, 0, 2, 5, slots.PayLTR, 2, utils.UInt8s{0, 0, 0, 0, 0}),
		slots.WinlinePayoutFromData(140, 0, 3, 5, slots.PayLTR, 3, utils.UInt8s{2, 2, 2, 2, 2}),
	)
	r6 = results.AcquireResult(
		slots.AcquireSpinResultFromData(utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}, nil, nil, nil, nil, 0, 0, 0),
		results.SpinData,
		slots.WinlinePayoutFromData(100, 0, 1, 5, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
		slots.WinlinePayoutFromData(125, 0, 2, 5, slots.PayLTR, 2, utils.UInt8s{0, 0, 0, 0, 0}),
		slots.WinlinePayoutFromData(150, 0, 3, 5, slots.PayLTR, 3, utils.UInt8s{2, 2, 2, 2, 2}),
	)
	r7 = results.AcquireResult(
		slots.AcquireSpinResultFromData(utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}, nil, nil, nil, nil, 0, 0, 0),
		results.SpinData,
		slots.WinlinePayoutFromData(110, 0, 1, 5, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
		slots.WinlinePayoutFromData(135, 0, 2, 5, slots.PayLTR, 2, utils.UInt8s{0, 0, 0, 0, 0}),
		slots.WinlinePayoutFromData(160, 0, 3, 5, slots.PayLTR, 3, utils.UInt8s{2, 2, 2, 2, 2}),
	)
)
