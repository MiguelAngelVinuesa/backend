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

func TestNewPayouts(t *testing.T) {
	p1 := slots.NewPayline(1, 3, 1, 1, 1, 1, 1)
	p2 := slots.NewPayline(2, 3, 0, 0, 0, 0, 0)
	p3 := slots.NewPayline(3, 3, 2, 2, 2, 2, 2)

	testCases := []struct {
		name       string
		rowCount   int
		maxSymbol  utils.Index
		paylines   []*slots.Payline
		results    results.Results
		payouts    analyse.Paylines
		allPayouts map[int]*analyse.Payline
		wilds      *analyse.ScatterPayout
		scatters   *analyse.ScatterPayout
		bonuses    *analyse.ScatterPayout
	}{
		{
			name:      "paylines, no results",
			rowCount:  3,
			maxSymbol: 3,
			paylines:  slots.Paylines{p1, p2, p3},
		},
		{
			name:      "paylines, 1 winline",
			rowCount:  3,
			maxSymbol: 12,
			paylines:  slots.Paylines{p1, p2, p3},
			results: results.Results{
				results.AcquireResult(nil, 0,
					slots.WinlinePayoutFromData(1.0, 0, 4, 3, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
				),
			},
			payouts: analyse.Paylines{nil,
				analyse.NewPayline(1, 12, utils.UInt8s{1, 1, 1, 1, 1}).Increase(4, 3, 1.0),
				analyse.NewPayline(2, 12, utils.UInt8s{0, 0, 0, 0, 0}),
				analyse.NewPayline(3, 12, utils.UInt8s{2, 2, 2, 2, 2}),
			},
		},
		{
			name:      "paylines, 5 winlines",
			rowCount:  3,
			maxSymbol: 12,
			paylines:  slots.Paylines{p1, p2, p3},
			results: results.Results{
				results.AcquireResult(nil, 0,
					slots.WinlinePayoutFromData(1.0, 0, 2, 4, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
					slots.WinlinePayoutFromData(5.0, 0, 4, 5, slots.PayLTR, 2, utils.UInt8s{0, 0, 0, 0, 0}),
					slots.WinlinePayoutFromData(1.5, 0, 6, 2, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
					slots.WinlinePayoutFromData(0.5, 0, 1, 3, slots.PayLTR, 3, utils.UInt8s{2, 2, 2, 2, 2}),
					slots.WinlinePayoutFromData(10.0, 0, 8, 4, slots.PayLTR, 2, utils.UInt8s{0, 0, 0, 0, 0}),
				),
			},
			payouts: analyse.Paylines{nil,
				analyse.NewPayline(1, 12, utils.UInt8s{1, 1, 1, 1, 1}).Increase(2, 4, 1.0).Increase(6, 2, 1.5),
				analyse.NewPayline(2, 12, utils.UInt8s{0, 0, 0, 0, 0}).Increase(4, 5, 5.0).Increase(8, 4, 10.0),
				analyse.NewPayline(3, 12, utils.UInt8s{2, 2, 2, 2, 2}).Increase(1, 3, 0.5),
			},
		},
		{
			name:      "paylines, 1 wilds payout",
			rowCount:  3,
			maxSymbol: 11,
			paylines:  slots.Paylines{p1, p2, p3},
			results: results.Results{
				results.AcquireResult(nil, 0,
					slots.WildSymbolPayoutWithMap(5.5, 1, 11, 4, utils.UInt8s{1, 3, 4, 7}),
				),
			},
			wilds: analyse.NewScatterPayout(11).Increase(11, 4, 5.5),
		},
		{
			name:      "paylines, 3 wilds payouts",
			rowCount:  3,
			maxSymbol: 11,
			paylines:  slots.Paylines{p1, p2, p3},
			results: results.Results{
				results.AcquireResult(nil, 0,
					slots.WildSymbolPayoutWithMap(2, 1, 10, 3, utils.UInt8s{2, 5, 8}),
					slots.WildSymbolPayoutWithMap(5.5, 1, 11, 4, utils.UInt8s{1, 3, 8, 9}),
					slots.WildSymbolPayoutWithMap(200, 1, 11, 5, utils.UInt8s{0, 6, 7, 12, 13}),
				),
			},
			wilds: analyse.NewScatterPayout(11).
				Increase(10, 3, 2).
				Increase(11, 4, 5.5).
				Increase(11, 5, 200),
		},
		{
			name:      "paylines, 1 scatters payout",
			rowCount:  3,
			maxSymbol: 11,
			paylines:  slots.Paylines{p1, p2, p3},
			results: results.Results{
				results.AcquireResult(nil, 0,
					slots.ScatterSymbolPayoutWithMap(5.5, 1, 11, 4, utils.UInt8s{1, 5, 9, 13}),
				),
			},
			scatters: analyse.NewScatterPayout(11).Increase(11, 4, 5.5),
		},
		{
			name:      "paylines, 3 scatters payouts",
			rowCount:  3,
			maxSymbol: 11,
			paylines:  slots.Paylines{p1, p2, p3},
			results: results.Results{
				results.AcquireResult(nil, 0,
					slots.ScatterSymbolPayoutWithMap(2, 1, 10, 3, utils.UInt8s{7, 9, 13}),
					slots.ScatterSymbolPayoutWithMap(5.5, 1, 11, 4, utils.UInt8s{3, 4, 8, 10}),
					slots.ScatterSymbolPayoutWithMap(200, 1, 11, 5, utils.UInt8s{0, 2, 6, 12, 14}),
				),
			},
			scatters: analyse.NewScatterPayout(11).
				Increase(10, 3, 2).
				Increase(11, 4, 5.5).
				Increase(11, 5, 200),
		},
		{
			name:      "paylines, 1 bonus symbol payout",
			rowCount:  3,
			maxSymbol: 11,
			paylines:  slots.Paylines{p1, p2, p3},
			results: results.Results{
				results.AcquireResult(nil, 0,
					slots.BonusSymbolPayout(5.5, 1, 7, 4),
				),
			},
			bonuses: analyse.NewScatterPayout(11).Increase(7, 4, 5.5),
		},
		{
			name:      "paylines, 3 bonus symbol payouts",
			rowCount:  3,
			maxSymbol: 11,
			paylines:  slots.Paylines{p1, p2, p3},
			results: results.Results{
				results.AcquireResult(nil, 0,
					slots.BonusSymbolPayout(2, 1, 5, 3),
					slots.BonusSymbolPayout(5.5, 1, 7, 4),
					slots.BonusSymbolPayout(200, 1, 2, 5),
				),
			},
			bonuses: analyse.NewScatterPayout(11).
				Increase(5, 3, 2).
				Increase(7, 4, 5.5).
				Increase(2, 5, 200),
		},
		{
			name:      "paylines, mixed payouts",
			rowCount:  3,
			maxSymbol: 11,
			paylines:  slots.Paylines{p1, p2, p3},
			results: results.Results{
				results.AcquireResult(nil, 0,
					slots.WinlinePayoutFromData(5.0, 0, 4, 5, slots.PayLTR, 2, utils.UInt8s{0, 0, 0, 0, 0}),
					slots.BonusSymbolPayout(2, 1, 5, 3),
					slots.WinlinePayoutFromData(1.5, 0, 6, 2, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
					slots.WinlinePayoutFromData(10.0, 0, 8, 4, slots.PayLTR, 2, utils.UInt8s{0, 0, 0, 0, 0}),
					slots.WildSymbolPayoutWithMap(5.5, 1, 11, 4, utils.UInt8s{1, 3, 6, 9}),
					slots.WinlinePayoutFromData(1.0, 0, 2, 4, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
					slots.ScatterSymbolPayoutWithMap(200, 1, 11, 5, utils.UInt8s{0, 3, 5, 11, 12}),
					slots.WinlinePayoutFromData(0.5, 0, 1, 3, slots.PayLTR, 3, utils.UInt8s{2, 2, 2, 2, 2}),
				),
			},
			payouts: analyse.Paylines{nil,
				analyse.NewPayline(1, 11, utils.UInt8s{1, 1, 1, 1, 1}).Increase(2, 4, 1.0).Increase(6, 2, 1.5),
				analyse.NewPayline(2, 11, utils.UInt8s{0, 0, 0, 0, 0}).Increase(4, 5, 5.0).Increase(8, 4, 10.0),
				analyse.NewPayline(3, 11, utils.UInt8s{2, 2, 2, 2, 2}).Increase(1, 3, 0.5),
			},
			wilds:    analyse.NewScatterPayout(11).Increase(11, 4, 5.5),
			scatters: analyse.NewScatterPayout(11).Increase(11, 5, 200),
			bonuses:  analyse.NewScatterPayout(11).Increase(5, 3, 2),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPayouts(tc.rowCount, tc.maxSymbol, tc.paylines, false)
			require.NotNil(t, p)
			defer p.Release()

			assert.Equal(t, tc.rowCount, p.rowCount)
			assert.Equal(t, tc.maxSymbol, p.maxSymbol)

			if len(tc.paylines) > 0 {
				assert.NotNil(t, p.Paylines)
			}

			assert.NotNil(t, p.AllPaylines)
			assert.NotNil(t, p.WildPayouts)
			assert.NotNil(t, p.ScatterPayouts)
			assert.NotNil(t, p.BonusPayouts)

			for ix := range tc.results {
				p.analyse(tc.results[ix])
			}

			if tc.payouts != nil {
				require.Equal(t, len(tc.payouts), len(p.Paylines))
				for ix := range tc.payouts {
					x1, x2 := tc.payouts[ix], p.Paylines[ix]
					if (x1 != nil || x2 != nil) && !x1.Equals(x2) {
						assert.EqualValues(t, x1, x2)
					}
				}
			}
			if tc.allPayouts != nil {
				require.Equal(t, len(tc.allPayouts), len(p.AllPaylines))
				for ix := range tc.allPayouts {
					x1, x2 := tc.allPayouts[ix], p.AllPaylines[ix]
					if (x1 != nil || x2 != nil) && !x1.Equals(x2) {
						assert.EqualValues(t, x1, x2)
					}
				}
			}
			if tc.wilds != nil {
				if !tc.wilds.Equals(p.WildPayouts) {
					assert.EqualValues(t, tc.wilds, p.WildPayouts)
				}
			}
			if tc.scatters != nil {
				if !tc.scatters.Equals(p.ScatterPayouts) {
					assert.EqualValues(t, tc.scatters, p.ScatterPayouts)
				}
			}
			if tc.bonuses != nil {
				if !tc.bonuses.Equals(p.BonusPayouts) {
					assert.EqualValues(t, tc.bonuses, p.BonusPayouts)
				}
			}

			n := p.Clone().(*Payouts)
			require.NotNil(t, n)
			defer n.Release()

			if tc.payouts != nil {
				require.Equal(t, len(tc.payouts), len(n.Paylines))
				for ix := range tc.payouts {
					x1, x2 := tc.payouts[ix], n.Paylines[ix]
					if (x1 != nil || x2 != nil) && !x1.Equals(x2) {
						assert.EqualValues(t, x1, x2)
					}
				}
			}
			if tc.allPayouts != nil {
				require.Equal(t, len(tc.allPayouts), len(n.AllPaylines))
				for ix := range tc.allPayouts {
					x1, x2 := tc.allPayouts[ix], n.AllPaylines[ix]
					if (x1 != nil || x2 != nil) && !x1.Equals(x2) {
						assert.EqualValues(t, x1, x2)
					}
				}
			}
			if tc.wilds != nil {
				if !tc.wilds.Equals(n.WildPayouts) {
					assert.EqualValues(t, tc.wilds, n.WildPayouts)
				}
			}
			if tc.scatters != nil {
				if !tc.scatters.Equals(n.ScatterPayouts) {
					assert.EqualValues(t, tc.scatters, n.ScatterPayouts)
				}
			}
			if tc.bonuses != nil {
				if !tc.bonuses.Equals(n.BonusPayouts) {
					assert.EqualValues(t, tc.bonuses, n.BonusPayouts)
				}
			}

			n.ResetData()

			for _, x := range n.Paylines {
				if x != nil {
					assert.Zero(t, x.Count)
				}
			}

			assert.Zero(t, len(n.AllPaylines))
			assert.Zero(t, n.rowCount)
			assert.Zero(t, n.maxSymbol)
		})
	}
}
