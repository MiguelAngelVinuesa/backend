package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewPaylineSet(t *testing.T) {
	testCases := []struct {
		name      string
		direction PayDirection
		highest   bool
		paylines  Paylines
	}{
		{
			name:      "single, normal, LTR",
			direction: PayLTR,
			paylines:  Paylines{pl01},
		},
		{
			name:      "single, normal, RTL",
			direction: PayRTL,
			paylines:  Paylines{pl02},
		},
		{
			name:      "single, normal, both",
			direction: PayBoth,
			paylines:  Paylines{pl02},
		},
		{
			name:      "few, normal, LTR",
			direction: PayRTL,
			paylines:  Paylines{pl01, pl02, pl03},
		},
		{
			name:      "few, normal, RTL",
			direction: PayRTL,
			paylines:  Paylines{pl01, pl02, pl03},
		},
		{
			name:      "few, normal, both",
			direction: PayBoth,
			paylines:  Paylines{pl01, pl02, pl03},
		},
		{
			name:      "ten, normal, LTR",
			direction: PayRTL,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
		},
		{
			name:      "ten, normal, RTL",
			direction: PayRTL,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
		},
		{
			name:      "ten, normal, both",
			direction: PayBoth,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
		},
		{
			name:      "all, normal, LTR",
			direction: PayRTL,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10, pl11, pl12, pl13, pl14, pl15, pl16, pl17, pl18, pl19, pl20, pl21},
		},
		{
			name:      "all, normal, RTL",
			direction: PayRTL,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10, pl11, pl12, pl13, pl14, pl15, pl16, pl17, pl18, pl19, pl20, pl21},
		},
		{
			name:      "all, normal, both",
			direction: PayBoth,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10, pl11, pl12, pl13, pl14, pl15, pl16, pl17, pl18, pl19, pl20, pl21},
		},
		{
			name:      "single, highest, LTR",
			direction: PayLTR,
			highest:   true,
			paylines:  Paylines{pl01},
		},
		{
			name:      "single, highest, RTL",
			direction: PayRTL,
			highest:   true,
			paylines:  Paylines{pl02},
		},
		{
			name:      "single, highest, both",
			direction: PayBoth,
			highest:   true,
			paylines:  Paylines{pl02},
		},
		{
			name:      "few, highest, LTR",
			direction: PayRTL,
			highest:   true,
			paylines:  Paylines{pl01, pl02, pl03},
		},
		{
			name:      "few, highest, RTL",
			direction: PayRTL,
			highest:   true,
			paylines:  Paylines{pl01, pl02, pl03},
		},
		{
			name:      "few, highest, both",
			direction: PayBoth,
			highest:   true,
			paylines:  Paylines{pl01, pl02, pl03},
		},
		{
			name:      "ten, highest, LTR",
			direction: PayRTL,
			highest:   true,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
		},
		{
			name:      "ten, highest, RTL",
			direction: PayRTL,
			highest:   true,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
		},
		{
			name:      "ten, highest, both",
			direction: PayBoth,
			highest:   true,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
		},
		{
			name:      "all, highest, LTR",
			direction: PayRTL,
			highest:   true,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10, pl11, pl12, pl13, pl14, pl15, pl16, pl17, pl18, pl19, pl20, pl21},
		},
		{
			name:      "all, highest, RTL",
			direction: PayRTL,
			highest:   true,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10, pl11, pl12, pl13, pl14, pl15, pl16, pl17, pl18, pl19, pl20, pl21},
		},
		{
			name:      "all, highest, both",
			direction: PayBoth,
			highest:   true,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10, pl11, pl12, pl13, pl14, pl15, pl16, pl17, pl18, pl19, pl20, pl21},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewPaylineSet(tc.direction, tc.highest, tc.paylines...)
			require.NotNil(t, s)

			assert.Equal(t, tc.direction, s.Directions())
			assert.Equal(t, tc.highest, s.HighestPayout())
			assert.EqualValues(t, tc.paylines, s.Paylines())
		})
	}
}

func TestPaylineSet_GetPayouts_Single(t *testing.T) {
	testCases := []struct {
		name        string
		payline     *Payline
		indexes     utils.Indexes
		multipliers []uint16
		direction   PayDirection
		want        bool
		symbol      utils.Index
		wantDir     PayDirection
		count       uint8
		factor      float64
		multiplier  float64
	}{
		{
			name:       "no win LTR",
			indexes:    utils.Indexes{0, 1, 2, 3, 4, 5, 4, 3, 2, 1, 0, 1, 2, 3, 4},
			direction:  PayLTR,
			multiplier: 1.0,
		},
		{
			name:       "no win RTL",
			indexes:    utils.Indexes{0, 1, 2, 3, 4, 5, 4, 3, 2, 1, 0, 1, 2, 3, 4},
			direction:  PayRTL,
			multiplier: 1.0,
		},
		{
			name:       "no win both",
			indexes:    utils.Indexes{0, 1, 2, 3, 4, 5, 4, 3, 2, 1, 0, 1, 2, 3, 4},
			direction:  PayBoth,
			multiplier: 1.0,
		},
		{
			name:       "almost win LTR",
			indexes:    utils.Indexes{0, 1, 2, 3, 1, 5, 4, 0, 2, 1, 3, 1, 2, 3, 4},
			direction:  PayLTR,
			multiplier: 1.0,
		},
		{
			name:       "almost win RTL",
			indexes:    utils.Indexes{0, 1, 2, 3, 1, 5, 4, 0, 2, 1, 3, 1, 2, 3, 4},
			direction:  PayRTL,
			multiplier: 1.0,
		},
		{
			name:       "almost win both",
			indexes:    utils.Indexes{0, 1, 2, 3, 1, 5, 4, 0, 2, 1, 3, 1, 2, 3, 4},
			direction:  PayBoth,
			multiplier: 1.0,
		},
		{
			name:       "small win LTR",
			indexes:    utils.Indexes{0, 1, 2, 3, 1, 5, 4, 1, 2, 1, 3, 1, 2, 3, 4},
			want:       true,
			direction:  PayLTR,
			symbol:     1,
			wantDir:    PayLTR,
			count:      3,
			factor:     0.5,
			multiplier: 1.0,
		},
		{
			name:       "small win RTL",
			indexes:    utils.Indexes{0, 1, 2, 3, 4, 5, 4, 3, 2, 1, 3, 1, 2, 3, 4},
			want:       true,
			direction:  PayRTL,
			symbol:     3,
			wantDir:    PayRTL,
			count:      3,
			factor:     1,
			multiplier: 1.0,
		},
		{
			name:       "small win both",
			indexes:    utils.Indexes{0, 1, 2, 3, 1, 5, 4, 3, 2, 1, 3, 1, 2, 3, 4},
			want:       true,
			direction:  PayBoth,
			symbol:     3,
			wantDir:    PayRTL,
			count:      3,
			factor:     1,
			multiplier: 1.0,
		},
		{
			name:       "big win LTR",
			indexes:    utils.Indexes{0, 1, 2, 3, 1, 5, 4, 1, 2, 1, 1, 1, 2, 1, 4},
			want:       true,
			direction:  PayLTR,
			symbol:     1,
			wantDir:    PayLTR,
			count:      5,
			factor:     4,
			multiplier: 1.0,
		},
		{
			name:       "big win RTL",
			indexes:    utils.Indexes{0, 1, 2, 3, 3, 5, 4, 3, 2, 1, 3, 1, 2, 3, 4},
			want:       true,
			direction:  PayRTL,
			symbol:     3,
			wantDir:    PayRTL,
			count:      4,
			factor:     2.5,
			multiplier: 1.0,
		},
		{
			name:       "big win both",
			indexes:    utils.Indexes{0, 1, 2, 3, 3, 5, 4, 3, 2, 1, 3, 1, 2, 3, 4},
			want:       true,
			direction:  PayBoth,
			symbol:     3,
			wantDir:    PayRTL,
			count:      4,
			factor:     2.5,
			multiplier: 1.0,
		},
		{
			name:       "2 scatter LTR",
			indexes:    utils.Indexes{0, 11, 2, 3, 11, 5, 4, 0, 2, 1, 3, 1, 2, 3, 4},
			want:       true,
			direction:  PayLTR,
			symbol:     11,
			wantDir:    PayLTR,
			count:      2,
			factor:     1,
			multiplier: 1.0,
		},
		{
			name:       "2 scatter RTL",
			indexes:    utils.Indexes{0, 1, 2, 3, 1, 5, 4, 0, 2, 1, 11, 1, 2, 11, 4},
			want:       true,
			direction:  PayRTL,
			symbol:     11,
			wantDir:    PayRTL,
			count:      2,
			factor:     1,
			multiplier: 1.0,
		},
		{
			name:       "2 scatter both",
			indexes:    utils.Indexes{0, 11, 2, 3, 11, 5, 4, 0, 2, 1, 11, 1, 2, 11, 4},
			want:       true,
			direction:  PayBoth,
			symbol:     11,
			wantDir:    PayLTR,
			count:      2,
			factor:     1,
			multiplier: 1.0,
		},
		{
			name:       "with 1 splits LTR",
			indexes:    utils.Indexes{0, 7, 2, 3, 1, 4, 4, 1, 2, 1, 3, 1, 2, 3, 4},
			want:       true,
			direction:  PayLTR,
			symbol:     1,
			wantDir:    PayLTR,
			count:      3,
			factor:     0.5,
			multiplier: 1.0,
		},
		{
			name:       "with 1 splits RTL",
			indexes:    utils.Indexes{0, 1, 2, 3, 4, 5, 4, 3, 2, 1, 8, 1, 2, 3, 4},
			want:       true,
			direction:  PayRTL,
			symbol:     3,
			wantDir:    PayRTL,
			count:      3,
			factor:     1,
			multiplier: 1.0,
		},
		{
			name:       "with 1 splits both",
			indexes:    utils.Indexes{0, 1, 2, 3, 1, 7, 4, 7, 2, 1, 3, 1, 2, 3, 4},
			want:       true,
			direction:  PayBoth,
			symbol:     1,
			wantDir:    PayLTR,
			count:      3,
			factor:     0.5,
			multiplier: 1.0,
		},
		{
			name:       "with 2 splits LTR",
			indexes:    utils.Indexes{0, 7, 2, 3, 7, 7, 4, 1, 2, 1, 3, 1, 2, 3, 4},
			want:       true,
			direction:  PayLTR,
			symbol:     1,
			wantDir:    PayLTR,
			count:      3,
			factor:     0.5,
			multiplier: 1.0,
		},
		{
			name:       "with 2 splits RTL",
			indexes:    utils.Indexes{0, 7, 2, 3, 7, 7, 4, 3, 2, 1, 8, 1, 2, 8, 4},
			want:       true,
			direction:  PayRTL,
			symbol:     3,
			wantDir:    PayRTL,
			count:      3,
			factor:     1,
			multiplier: 1.0,
		},
		{
			name:       "with 2 splits both",
			indexes:    utils.Indexes{0, 7, 2, 3, 7, 7, 4, 1, 2, 1, 8, 1, 2, 8, 4},
			want:       true,
			direction:  PayBoth,
			symbol:     1,
			wantDir:    PayLTR,
			count:      3,
			factor:     0.5,
			multiplier: 1.0,
		},
		{
			name:       "with 4 splits LTR",
			indexes:    utils.Indexes{0, 8, 2, 3, 8, 5, 4, 3, 2, 1, 8, 1, 2, 8, 4},
			want:       true,
			direction:  PayLTR,
			symbol:     3,
			wantDir:    PayLTR,
			count:      5,
			factor:     5,
			multiplier: 1.0,
		},
		{
			name:       "with 4 splits RTL",
			indexes:    utils.Indexes{0, 8, 2, 3, 8, 5, 4, 8, 2, 1, 8, 1, 2, 3, 4},
			want:       true,
			direction:  PayRTL,
			symbol:     3,
			wantDir:    PayRTL,
			count:      5,
			factor:     5,
			multiplier: 1.0,
		},
		{
			name:       "with 4 splits both",
			indexes:    utils.Indexes{0, 7, 2, 3, 7, 7, 4, 7, 2, 1, 7, 1, 2, 1, 4},
			want:       true,
			direction:  PayBoth,
			symbol:     1,
			wantDir:    PayLTR,
			count:      5,
			factor:     4,
			multiplier: 1.0,
		},
		{
			name:       "with 1 wild+split LTR",
			indexes:    utils.Indexes{0, 9, 2, 3, 1, 7, 4, 7, 2, 1, 3, 1, 2, 3, 4},
			want:       true,
			direction:  PayLTR,
			symbol:     1,
			wantDir:    PayLTR,
			count:      3,
			factor:     0.5,
			multiplier: 1.0,
		},
		{
			name:       "with 1 wild+split RTL",
			indexes:    utils.Indexes{0, 1, 2, 3, 4, 7, 4, 3, 2, 1, 12, 1, 2, 8, 4},
			want:       true,
			direction:  PayRTL,
			symbol:     3,
			wantDir:    PayRTL,
			count:      3,
			factor:     1,
			multiplier: 1.0,
		},
		{
			name:       "with 1 wild+split both",
			indexes:    utils.Indexes{0, 7, 2, 3, 1, 7, 4, 9, 2, 1, 8, 1, 2, 3, 4},
			want:       true,
			direction:  PayBoth,
			symbol:     3,
			wantDir:    PayRTL,
			count:      3,
			factor:     1,
			multiplier: 1.0,
		},
		{
			name:       "with 2 wilds+splits LTR",
			indexes:    utils.Indexes{0, 12, 2, 3, 7, 7, 4, 1, 2, 1, 7, 1, 2, 9, 4},
			want:       true,
			direction:  PayLTR,
			symbol:     1,
			wantDir:    PayLTR,
			count:      5,
			factor:     4,
			multiplier: 1.0,
		},
		{
			name:       "with 2 wilds+splits RTL",
			indexes:    utils.Indexes{0, 3, 2, 3, 9, 5, 4, 8, 2, 1, 12, 1, 2, 8, 4},
			want:       true,
			direction:  PayRTL,
			symbol:     3,
			wantDir:    PayRTL,
			count:      5,
			factor:     5,
			multiplier: 1.0,
		},
		{
			name:       "with 2 wilds+splits both",
			indexes:    utils.Indexes{0, 12, 2, 3, 8, 5, 4, 9, 2, 1, 8, 1, 2, 3, 4},
			want:       true,
			direction:  PayBoth,
			symbol:     3,
			wantDir:    PayLTR,
			count:      5,
			factor:     5,
			multiplier: 1.0,
		},
		{
			name:       "3-4-5-4-3 silly payline",
			payline:    NewPayline(1, 5, 3, 3, 3, 3, 3),
			indexes:    utils.Indexes{1, 2, 3, 0, 0, 5, 4, 3, 2, 0, 1, 2, 3, 4, 5, 5, 4, 3, 2, 0, 1, 2, 3, 0, 0},
			want:       false,
			direction:  PayLTR,
			multiplier: 1.0,
		},
		{
			name:        "small win LTR + multiplier",
			indexes:     utils.Indexes{0, 1, 2, 3, 12, 7, 4, 1, 2, 1, 3, 1, 2, 3, 4},
			multipliers: []uint16{0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want:        true,
			direction:   PayLTR,
			symbol:      1,
			wantDir:     PayLTR,
			count:       3,
			factor:      0.5,
			multiplier:  5.0,
		},
		{
			name:        "small win LTR + multipliers",
			indexes:     utils.Indexes{0, 1, 2, 3, 12, 7, 4, 1, 2, 1, 3, 1, 2, 3, 4},
			multipliers: []uint16{0, 2, 0, 0, 5, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0},
			want:        true,
			direction:   PayLTR,
			symbol:      1,
			wantDir:     PayLTR,
			count:       3,
			factor:      0.5,
			multiplier:  30.0,
		},
		{
			name:       "highest payout 2 wilds",
			indexes:    utils.Indexes{0, 12, 2, 3, 12, 5, 4, 1, 6, 1, 3, 1, 2, 3, 4},
			want:       true,
			direction:  PayLTR,
			symbol:     12,
			wantDir:    PayLTR,
			count:      2,
			factor:     2,
			multiplier: 1.0,
		},
		{
			name:       "not highest payout 2 wilds",
			indexes:    utils.Indexes{0, 12, 2, 3, 12, 5, 4, 1, 6, 2, 1, 2, 2, 3, 4},
			want:       true,
			direction:  PayLTR,
			symbol:     1,
			wantDir:    PayLTR,
			count:      4,
			factor:     2,
			multiplier: 1.0,
		},
		{
			name:       "highest payout 3 wilds",
			indexes:    utils.Indexes{0, 12, 2, 3, 12, 5, 4, 12, 6, 2, 1, 2, 4, 5, 4},
			want:       true,
			direction:  PayLTR,
			symbol:     12,
			wantDir:    PayLTR,
			count:      3,
			factor:     4,
			multiplier: 1.0,
		},
		{
			name:       "not highest payout 3 wilds",
			indexes:    utils.Indexes{0, 12, 2, 3, 12, 5, 4, 12, 6, 2, 1, 2, 4, 1, 4},
			want:       true,
			direction:  PayLTR,
			symbol:     1,
			wantDir:    PayLTR,
			count:      5,
			factor:     4,
			multiplier: 1.0,
		},
		{
			name:       "highest payout 4 wilds",
			indexes:    utils.Indexes{0, 12, 2, 3, 12, 5, 4, 12, 6, 2, 12, 2, 4, 1, 4},
			want:       true,
			direction:  PayLTR,
			symbol:     12,
			wantDir:    PayLTR,
			count:      4,
			factor:     6,
			multiplier: 1.0,
		},
		{
			name:       "not highest payout 4 wilds",
			indexes:    utils.Indexes{0, 12, 2, 3, 12, 5, 4, 12, 6, 2, 12, 2, 4, 6, 4},
			want:       true,
			direction:  PayLTR,
			symbol:     6,
			wantDir:    PayLTR,
			count:      5,
			factor:     7,
			multiplier: 1.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pl := tc.payline
			if pl == nil {
				pl = pl01
			}

			s := NewPaylineSet(tc.direction, true, pl)
			require.NotNil(t, s)

			spin := &Spin{
				slots: &Slots{
					reelCount:     5,
					rowCount:      3,
					symbols:       setF1,
					directions:    tc.direction,
					highestPayout: true,
				},
				symbols:     setF1,
				indexes:     tc.indexes,
				multipliers: tc.multipliers,
				reelCount:   5,
				rowCount:    3,
			}

			result := results.AcquireResult(nil, 0)
			defer result.Release()

			got := s.GetPayouts(spin, result)

			if !tc.want {
				require.Empty(t, result.Payouts)
				assert.False(t, got)
			} else {
				require.NotEmpty(t, result.Payouts)
				assert.True(t, got)

				p, ok := result.Payouts[0].(*SpinPayout)
				require.True(t, ok)
				require.NotNil(t, p)

				assert.Equal(t, tc.symbol, p.symbol)
				assert.Equal(t, tc.wantDir, p.direction)
				assert.Equal(t, tc.count, p.count)
				assert.Equal(t, tc.factor, p.factor)
				assert.Equal(t, tc.multiplier, p.multiplier)
			}
		})
	}
}

func TestPaylineSet_GetPayouts_Multi(t *testing.T) {
	testCases := []struct {
		name        string
		direction   PayDirection
		paylines    Paylines
		indexes     utils.Indexes
		multipliers []uint16
		spinMult    float64
		want        int
		total       float64
	}{
		{
			name:      "10 lines LTR, no multipliers, no payouts",
			direction: PayLTR,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:      "10 lines RTL, no multipliers, no payouts",
			direction: PayRTL,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:      "10 lines Both, no multipliers, no payouts",
			direction: PayBoth,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:        "10 lines LTR, multipliers, no payouts",
			direction:   PayLTR,
			paylines:    Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:     utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			multipliers: []uint16{2, 2, 2, 3, 4, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:        "10 lines RTL, multipliers, no payouts",
			paylines:    Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			direction:   PayRTL,
			indexes:     utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			multipliers: []uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 4, 5, 2, 2, 2},
		},
		{
			name:        "10 lines Both, multipliers, no payouts",
			indexes:     utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			paylines:    Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			multipliers: []uint16{2, 2, 2, 3, 4, 5, 0, 0, 0, 3, 4, 5, 2, 2, 2},
			direction:   PayBoth,
		},
		{
			name:      "10 lines LTR, no multipliers, almost",
			direction: PayLTR,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:   utils.Indexes{1, 2, 3, 1, 2, 3, 4, 5, 6, 1, 2, 3, 1, 2, 3},
		},
		{
			name:      "10 lines RTL, no multipliers, almost",
			direction: PayRTL,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:   utils.Indexes{1, 2, 3, 1, 2, 3, 4, 5, 6, 1, 2, 3, 1, 2, 3},
		},
		{
			name:      "10 lines Both, no multipliers, almost",
			direction: PayBoth,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:   utils.Indexes{1, 2, 3, 1, 2, 3, 4, 5, 6, 1, 2, 3, 1, 2, 3},
		},
		{
			name:        "10 lines LTR, multipliers, almost",
			direction:   PayLTR,
			paylines:    Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:     utils.Indexes{1, 2, 3, 1, 2, 3, 4, 5, 6, 1, 2, 3, 1, 2, 3},
			multipliers: []uint16{2, 2, 2, 3, 4, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:        "10 lines RTL, multipliers, almost",
			direction:   PayRTL,
			paylines:    Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:     utils.Indexes{1, 2, 3, 1, 2, 3, 4, 5, 6, 1, 2, 3, 1, 2, 3},
			multipliers: []uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 4, 5, 2, 2, 2},
		},
		{
			name:        "10 lines Both, multipliers, almost",
			direction:   PayBoth,
			paylines:    Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:     utils.Indexes{1, 2, 3, 1, 2, 3, 4, 5, 6, 1, 2, 3, 1, 2, 3},
			multipliers: []uint16{2, 2, 2, 3, 4, 5, 0, 0, 0, 3, 4, 5, 2, 2, 2},
		},
		{
			name:      "10 lines LTR, no multipliers, 1 small",
			direction: PayLTR,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:   utils.Indexes{1, 2, 3, 1, 2, 3, 1, 5, 6, 4, 2, 3, 4, 2, 3},
			want:      1,
			total:     0.5,
		},
		{
			name:      "10 lines RTL, no multipliers, 1 small",
			direction: PayRTL,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:   utils.Indexes{1, 2, 3, 1, 2, 6, 4, 5, 3, 1, 2, 3, 1, 2, 3},
			want:      1,
			total:     1,
		},
		{
			name:      "10 lines Both, no multipliers, 1 small LTR",
			direction: PayBoth,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:   utils.Indexes{1, 2, 3, 1, 2, 3, 4, 2, 6, 1, 5, 3, 1, 2, 3},
			want:      1,
			total:     0.5,
		},
		{
			name:      "10 lines Both, no multipliers, 1 small RTL",
			direction: PayBoth,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:   utils.Indexes{1, 2, 3, 1, 5, 3, 4, 2, 6, 1, 2, 3, 1, 2, 3},
			spinMult:  5,
			want:      1,
			total:     2.5,
		},
		{
			name:      "10 lines Both, no multipliers, 1 highest RTL",
			direction: PayBoth,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:   utils.Indexes{1, 2, 3, 4, 2, 6, 1, 9, 3, 4, 5, 6, 1, 5, 3},
			spinMult:  3.7,
			want:      1,
			total:     3.7,
		},
		{
			name:        "10 lines LTR, multipliers, 1 small",
			direction:   PayLTR,
			paylines:    Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:     utils.Indexes{1, 2, 3, 1, 2, 3, 1, 5, 6, 4, 2, 3, 4, 2, 3},
			multipliers: []uint16{2, 2, 2, 3, 4, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want:        1,
			total:       3,
		},
		{
			name:        "10 lines RTL, multipliers, 1 small",
			direction:   PayRTL,
			paylines:    Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:     utils.Indexes{1, 2, 3, 1, 2, 6, 4, 5, 3, 1, 2, 3, 1, 2, 3},
			multipliers: []uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 4, 5, 2, 2, 2},
			spinMult:    2.5,
			want:        1,
			total:       25,
		},
		{
			name:        "10 lines Both, multipliers, 1 small LTR",
			direction:   PayBoth,
			paylines:    Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:     utils.Indexes{1, 2, 3, 1, 2, 3, 4, 2, 6, 1, 5, 3, 1, 2, 3},
			multipliers: []uint16{2, 2, 2, 3, 4, 5, 0, 0, 0, 3, 4, 5, 2, 2, 2},
			want:        1,
			total:       4,
		},
		{
			name:        "10 lines Both, multipliers, 1 small RTL",
			direction:   PayBoth,
			paylines:    Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:     utils.Indexes{1, 2, 3, 1, 5, 3, 4, 2, 6, 1, 2, 3, 1, 2, 3},
			multipliers: []uint16{2, 2, 2, 3, 4, 5, 0, 0, 0, 3, 4, 5, 2, 2, 2},
			want:        1,
			total:       4,
		},
		{
			name:        "10 lines Both, multipliers, 1 highest RTL",
			direction:   PayBoth,
			paylines:    Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:     utils.Indexes{1, 2, 3, 4, 2, 6, 1, 9, 3, 4, 5, 6, 1, 5, 3},
			multipliers: []uint16{2, 2, 2, 3, 4, 5, 0, 0, 0, 3, 4, 5, 2, 2, 2},
			want:        1,
			total:       8,
		},
		{
			name:      "10 lines LTR, no multipliers, 3",
			direction: PayLTR,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:   utils.Indexes{1, 2, 3, 1, 2, 3, 1, 3, 3, 3, 2, 6, 3, 2, 6},
			want:      3,
			total:     6.5,
		},
		{
			name:      "10 lines RTL, no multipliers, 3",
			direction: PayRTL,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:   utils.Indexes{4, 2, 1, 4, 2, 1, 5, 1, 3, 1, 2, 3, 1, 2, 3},
			want:      3,
			total:     5.5,
		},
		{
			name:      "10 lines Both, no multipliers, 3 LTR",
			direction: PayBoth,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:   utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 5, 3, 1, 2, 6},
			want:      3,
			total:     7,
		},
		{
			name:      "10 lines Both, no multipliers, 3 RTL",
			direction: PayBoth,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 1, 2, 3, 1, 2, 3},
			want:      3,
			total:     2,
		},
		{
			name:      "10 lines Both, no multipliers, 3 highest mixed",
			direction: PayBoth,
			paylines:  Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:   utils.Indexes{1, 2, 3, 4, 2, 3, 3, 2, 3, 4, 3, 6, 1, 5, 3},
			want:      3,
			total:     2.5,
		},
		{
			name:        "10 lines LTR, multipliers, 3",
			direction:   PayLTR,
			paylines:    Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:     utils.Indexes{1, 2, 3, 1, 2, 3, 1, 3, 3, 3, 2, 6, 3, 2, 6},
			multipliers: []uint16{2, 2, 2, 0, 0, 0, 3, 4, 5, 0, 0, 0, 2, 2, 2},
			want:        3,
			total:       93,
		},
		{
			name:        "10 lines RTL, multipliers, 3",
			direction:   PayRTL,
			paylines:    Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:     utils.Indexes{4, 2, 1, 4, 2, 1, 5, 1, 3, 1, 2, 3, 1, 2, 3},
			multipliers: []uint16{2, 2, 2, 0, 0, 0, 3, 4, 5, 0, 0, 0, 2, 2, 2},
			want:        3,
			total:       78,
		},
		{
			name:        "10 lines Both, multipliers, 3 LTR",
			direction:   PayBoth,
			paylines:    Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:     utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 5, 3, 1, 2, 6},
			multipliers: []uint16{2, 2, 2, 0, 0, 0, 3, 4, 5, 0, 0, 0, 2, 2, 2},
			want:        3,
			total:       77,
		},
		{
			name:        "10 lines Both, multipliers, 3 RTL",
			direction:   PayBoth,
			paylines:    Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:     utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 1, 2, 3, 1, 2, 3},
			multipliers: []uint16{2, 2, 2, 0, 0, 0, 3, 4, 5, 0, 0, 0, 2, 2, 2},
			want:        3,
			total:       17,
		},
		{
			name:        "10 lines Both, multipliers, 3 highest mixed",
			direction:   PayBoth,
			paylines:    Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10},
			indexes:     utils.Indexes{1, 2, 3, 4, 2, 3, 3, 2, 3, 4, 3, 6, 1, 5, 3},
			multipliers: []uint16{2, 2, 2, 0, 0, 0, 3, 4, 5, 0, 0, 0, 2, 2, 2},
			want:        3,
			total:       20,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewPaylineSet(tc.direction, true, tc.paylines...)
			require.NotNil(t, s)

			spin := &Spin{
				slots: &Slots{
					reelCount:  5,
					rowCount:   3,
					symbols:    setF1,
					directions: tc.direction,
				},
				symbols:     setF1,
				indexes:     tc.indexes,
				multipliers: tc.multipliers,
				multiplier:  tc.spinMult,
				reelCount:   5,
				rowCount:    3,
			}

			result := results.AcquireResult(nil, 0)
			defer result.Release()

			got := s.GetPayouts(spin, result)
			if tc.want > 0 {
				assert.True(t, got)
			} else {
				assert.False(t, got)
			}

			assert.Equal(t, tc.want, len(result.Payouts))
			assert.Equal(t, tc.total, result.Total)
		})
	}
}

var (
	pl01 = NewPayline(1, 3, 1, 1, 1, 1, 1)
	pl02 = NewPayline(2, 3, 0, 0, 0, 0, 0)
	pl03 = NewPayline(3, 3, 2, 2, 2, 2, 2)
	pl04 = NewPayline(4, 3, 0, 1, 1, 1, 2)
	pl05 = NewPayline(5, 3, 2, 1, 1, 1, 0)
	pl06 = NewPayline(6, 3, 0, 1, 2, 1, 0)
	pl07 = NewPayline(7, 3, 2, 1, 0, 1, 2)
	pl08 = NewPayline(8, 3, 0, 0, 1, 2, 2)
	pl09 = NewPayline(9, 3, 2, 2, 1, 0, 0)
	pl10 = NewPayline(10, 3, 0, 0, 1, 0, 0)
	pl11 = NewPayline(11, 3, 2, 2, 1, 2, 2)
	pl12 = NewPayline(12, 3, 0, 1, 0, 1, 0)
	pl13 = NewPayline(13, 3, 2, 1, 2, 1, 2)
	pl14 = NewPayline(14, 3, 0, 2, 0, 2, 0)
	pl15 = NewPayline(15, 3, 2, 0, 2, 0, 2)
	pl16 = NewPayline(16, 3, 1, 0, 0, 0, 1)
	pl17 = NewPayline(17, 3, 1, 2, 2, 2, 1)
	pl18 = NewPayline(18, 3, 1, 1, 0, 1, 1)
	pl19 = NewPayline(19, 3, 1, 1, 2, 1, 1)
	pl20 = NewPayline(20, 3, 1, 0, 1, 2, 1)
	pl21 = NewPayline(21, 3, 1, 2, 1, 0, 1)

	lines10 = Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10}
	lines21 = Paylines{pl01, pl02, pl03, pl04, pl05, pl06, pl07, pl08, pl09, pl10, pl11, pl12, pl13, pl14, pl15, pl16, pl17, pl18, pl19, pl20, pl21}
)

func BenchmarkPaylineSet_GetPayouts_10linesLTR_none(b *testing.B) {
	indexes := utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3}

	s := NewPaylineSet(PayLTR, true, lines10...)

	spin := &Spin{
		slots: &Slots{
			reelCount:  5,
			rowCount:   3,
			symbols:    setF1,
			directions: PayLTR,
		},
		symbols:   setF1,
		indexes:   indexes,
		reelCount: 5,
		rowCount:  3,
	}

	result := results.AcquireResult(nil, 0)
	defer result.Release()

	for i := 0; i < b.N; i++ {
		s.GetPayouts(spin, result)
		result.Payouts = results.ReleasePayouts(result.Payouts)
	}
}

func BenchmarkPaylineSet_GetPayouts_10linesRTL_none(b *testing.B) {
	indexes := utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3}

	s := NewPaylineSet(PayRTL, true, lines10...)

	spin := &Spin{
		slots: &Slots{
			reelCount:  5,
			rowCount:   3,
			symbols:    setF1,
			directions: PayRTL,
		},
		symbols:   setF1,
		indexes:   indexes,
		reelCount: 5,
		rowCount:  3,
	}

	result := results.AcquireResult(nil, 0)
	defer result.Release()

	for i := 0; i < b.N; i++ {
		s.GetPayouts(spin, result)
		result.Payouts = results.ReleasePayouts(result.Payouts)
	}
}

func BenchmarkPaylineSet_GetPayouts_10linesBoth_none(b *testing.B) {
	indexes := utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3}

	s := NewPaylineSet(PayBoth, true, lines10...)

	spin := &Spin{
		slots: &Slots{
			reelCount:  5,
			rowCount:   3,
			symbols:    setF1,
			directions: PayBoth,
		},
		symbols:   setF1,
		indexes:   indexes,
		reelCount: 5,
		rowCount:  3,
	}

	result := results.AcquireResult(nil, 0)
	defer result.Release()

	for i := 0; i < b.N; i++ {
		s.GetPayouts(spin, result)
		result.Payouts = results.ReleasePayouts(result.Payouts)
	}
}

func BenchmarkPaylineSet_GetPayouts_21linesLTR_none(b *testing.B) {
	indexes := utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3}

	s := NewPaylineSet(PayLTR, true, lines21...)

	spin := &Spin{
		slots: &Slots{
			reelCount:  5,
			rowCount:   3,
			symbols:    setF1,
			directions: PayLTR,
		},
		symbols:   setF1,
		indexes:   indexes,
		reelCount: 5,
		rowCount:  3,
	}

	result := results.AcquireResult(nil, 0)
	defer result.Release()

	for i := 0; i < b.N; i++ {
		s.GetPayouts(spin, result)
		result.Payouts = results.ReleasePayouts(result.Payouts)
	}
}

func BenchmarkPaylineSet_GetPayouts_21linesRTL_none(b *testing.B) {
	indexes := utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3}

	s := NewPaylineSet(PayRTL, true, lines21...)

	spin := &Spin{
		slots: &Slots{
			reelCount:  5,
			rowCount:   3,
			symbols:    setF1,
			directions: PayRTL,
		},
		symbols:   setF1,
		indexes:   indexes,
		reelCount: 5,
		rowCount:  3,
	}

	result := results.AcquireResult(nil, 0)
	defer result.Release()

	for i := 0; i < b.N; i++ {
		s.GetPayouts(spin, result)
		result.Payouts = results.ReleasePayouts(result.Payouts)
	}
}

func BenchmarkPaylineSet_GetPayouts_21linesBoth_none(b *testing.B) {
	indexes := utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3}

	s := NewPaylineSet(PayBoth, true, lines21...)

	spin := &Spin{
		slots: &Slots{
			reelCount:  5,
			rowCount:   3,
			symbols:    setF1,
			directions: PayBoth,
		},
		symbols:   setF1,
		indexes:   indexes,
		reelCount: 5,
		rowCount:  3,
	}

	result := results.AcquireResult(nil, 0)
	defer result.Release()

	for i := 0; i < b.N; i++ {
		s.GetPayouts(spin, result)
		result.Payouts = results.ReleasePayouts(result.Payouts)
	}
}
