package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestAcquireAllPaylinesEvent(t *testing.T) {
	testCases := []struct {
		name       string
		symbol     utils.Index
		reels      int
		rows       int
		payouts    results.Payouts
		wantFactor float64
		wantCount  uint16
		wantMap    utils.UInt8s
		wantItems  int
	}{
		{
			name:   "single payout",
			symbol: 4,
			reels:  6,
			rows:   4,
			payouts: results.Payouts{
				&SpinPayout{kind: results.SlotWinline, count: 4, direction: PayLTR, symbol: 4, factor: 2.5, payRows: utils.UInt8s{0, 0, 1, 2, 0, 0}},
			},
			wantCount:  1,
			wantFactor: 2.5,
			wantMap:    utils.UInt8s{0, 4, 9, 14},
		},
		{
			name:   "2 payouts, same symbol",
			symbol: 4,
			reels:  6,
			rows:   4,
			payouts: results.Payouts{
				&SpinPayout{kind: results.SlotWinline, count: 4, direction: PayLTR, symbol: 4, factor: 2.5, payRows: utils.UInt8s{0, 0, 1, 2, 0, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 4, direction: PayLTR, symbol: 4, factor: 2.5, payRows: utils.UInt8s{0, 2, 1, 2, 0, 0}},
			},
			wantCount:  2,
			wantFactor: 5,
			wantMap:    utils.UInt8s{0, 4, 6, 9, 14},
		},
		{
			name:   "2 payouts, different symbol",
			symbol: 3,
			reels:  6,
			rows:   4,
			payouts: results.Payouts{
				&SpinPayout{kind: results.SlotWinline, count: 4, direction: PayLTR, symbol: 4, factor: 2.5, payRows: utils.UInt8s{0, 0, 1, 2, 0, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 3, direction: PayLTR, symbol: 3, factor: 1.5, payRows: utils.UInt8s{1, 2, 0, 0, 0, 0}},
			},
			wantCount:  1,
			wantFactor: 1.5,
			wantMap:    utils.UInt8s{1, 6, 8},
		},
		{
			name:   "9 payouts, same symbol",
			symbol: 6,
			reels:  6,
			rows:   4,
			payouts: results.Payouts{
				&SpinPayout{kind: results.SlotWinline, count: 5, direction: PayLTR, symbol: 6, factor: 15, payRows: utils.UInt8s{0, 0, 1, 3, 0, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 5, direction: PayLTR, symbol: 6, factor: 15, payRows: utils.UInt8s{0, 0, 1, 3, 0, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 5, direction: PayLTR, symbol: 6, factor: 15, payRows: utils.UInt8s{0, 0, 1, 3, 0, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 5, direction: PayLTR, symbol: 6, factor: 15, payRows: utils.UInt8s{0, 1, 1, 3, 1, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 5, direction: PayLTR, symbol: 6, factor: 15, payRows: utils.UInt8s{0, 1, 1, 3, 1, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 5, direction: PayLTR, symbol: 6, factor: 15, payRows: utils.UInt8s{0, 1, 1, 3, 1, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 5, direction: PayLTR, symbol: 6, factor: 15, payRows: utils.UInt8s{0, 2, 1, 3, 2, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 5, direction: PayLTR, symbol: 6, factor: 15, payRows: utils.UInt8s{0, 2, 1, 3, 2, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 5, direction: PayLTR, symbol: 6, factor: 15, payRows: utils.UInt8s{0, 2, 1, 3, 2, 0}},
			},
			wantCount:  9,
			wantFactor: 135,
			wantMap:    utils.UInt8s{0, 4, 5, 6, 9, 15, 16, 17, 18},
		},
		{
			name:   "9 payouts, different symbols, with multipliers",
			symbol: 3,
			reels:  6,
			rows:   4,
			payouts: results.Payouts{
				&SpinPayout{kind: results.SlotWinline, count: 3, direction: PayLTR, symbol: 6, factor: 2.5, payRows: utils.UInt8s{0, 0, 1, 0, 0, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 3, direction: PayLTR, symbol: 6, factor: 2.5, payRows: utils.UInt8s{0, 0, 1, 0, 0, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 3, direction: PayLTR, symbol: 6, factor: 2.5, multiplier: 2, payRows: utils.UInt8s{0, 0, 1, 0, 0, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 3, direction: PayLTR, symbol: 6, factor: 2.5, payRows: utils.UInt8s{0, 1, 1, 0, 0, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 3, direction: PayLTR, symbol: 6, factor: 2.5, multiplier: 3, payRows: utils.UInt8s{0, 1, 1, 0, 0, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 5, direction: PayLTR, symbol: 3, factor: 15, payRows: utils.UInt8s{1, 2, 0, 2, 1, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 5, direction: PayLTR, symbol: 3, factor: 15, payRows: utils.UInt8s{1, 2, 0, 2, 1, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 5, direction: PayLTR, symbol: 3, factor: 15, multiplier: 3, payRows: utils.UInt8s{1, 2, 3, 2, 2, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 5, direction: PayLTR, symbol: 3, factor: 15, multiplier: 5, payRows: utils.UInt8s{1, 2, 3, 2, 2, 0}},
			},
			wantCount:  10,
			wantFactor: 150,
			wantMap:    utils.UInt8s{1, 6, 8, 11, 14, 17, 18},
		},
		{
			name:   "9 payouts, same symbol, different counts, with multipliers",
			symbol: 5,
			reels:  6,
			rows:   4,
			payouts: results.Payouts{
				&SpinPayout{kind: results.SlotWinline, count: 3, direction: PayLTR, symbol: 5, factor: 2.5, payRows: utils.UInt8s{0, 0, 1, 0, 0, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 6, direction: PayLTR, symbol: 5, factor: 15, payRows: utils.UInt8s{0, 0, 1, 0, 0, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 3, direction: PayLTR, symbol: 5, factor: 2.5, payRows: utils.UInt8s{0, 0, 1, 0, 0, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 4, direction: PayLTR, symbol: 5, factor: 5, payRows: utils.UInt8s{0, 1, 1, 0, 0, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 6, direction: PayLTR, symbol: 5, factor: 15, payRows: utils.UInt8s{0, 1, 1, 0, 0, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 5, direction: PayLTR, symbol: 5, factor: 10, payRows: utils.UInt8s{1, 2, 0, 2, 1, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 5, direction: PayLTR, symbol: 5, factor: 10, payRows: utils.UInt8s{1, 2, 0, 2, 1, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 4, direction: PayLTR, symbol: 5, factor: 5, multiplier: 2, payRows: utils.UInt8s{1, 2, 3, 2, 2, 0}},
				&SpinPayout{kind: results.SlotWinline, count: 5, direction: PayLTR, symbol: 5, factor: 10, multiplier: 4, payRows: utils.UInt8s{1, 2, 3, 2, 2, 0}},
			},
			wantItems: 4,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			events := AcquireAllPaylinesEvents(tc.symbol, tc.reels, tc.rows, tc.payouts)

			if tc.wantItems != 0 {
				require.Equal(t, tc.wantItems, len(events))
			} else {
				require.Equal(t, 1, len(events))

				for _, e := range events {
					event, ok := e.(*AllPaylinesEvent)
					require.True(t, ok)
					require.NotNil(t, event)

					assert.Equal(t, tc.symbol, event.symbol)
					assert.Equal(t, tc.wantCount, event.paylines)
					assert.Equal(t, tc.wantFactor, event.factor)
					assert.EqualValues(t, tc.wantMap, event.payMap)

					defer event.Release()
				}
			}
		})
	}
}
