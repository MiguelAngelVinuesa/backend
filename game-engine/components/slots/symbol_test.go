package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSymbol(t *testing.T) {
	testCases := []struct {
		name       string
		id         utils.Index
		opts       []SymbolOption
		symbol     string
		resource   string
		kind       SymbolKind
		weights    []float64
		payouts    []float64
		wildFor    utils.Indexes
		multiplier float64
		payable    utils.UInt8s
		minPayable uint8
		morphInto  utils.Index
		sticky     bool
		split      bool
		wild       bool
		hero       bool
		scatter    bool
		bomb       bool
		shooter    bool
		prize      bool
		varyMult   bool
	}{
		{
			name:       "basic",
			id:         0,
			kind:       Standard,
			multiplier: 1.0,
		},
		{
			name:       "with name",
			id:         1,
			opts:       []SymbolOption{WithName("haha")},
			kind:       Standard,
			symbol:     "haha",
			multiplier: 1.0,
		},
		{
			name:       "with resource",
			id:         2,
			opts:       []SymbolOption{WithResource("/haha.gif")},
			kind:       Standard,
			resource:   "/haha.gif",
			multiplier: 1.0,
		},
		{
			name:       "with kind wild",
			id:         3,
			opts:       []SymbolOption{WithKind(Wild)},
			kind:       Wild,
			wild:       true,
			multiplier: 1.0,
		},
		{
			name:       "with kind wild/scatter",
			id:         4,
			opts:       []SymbolOption{WithKind(WildScatter)},
			kind:       WildScatter,
			wild:       true,
			scatter:    true,
			multiplier: 1.0,
		},
		{
			name:       "with kind hero",
			id:         5,
			opts:       []SymbolOption{WithKind(Hero)},
			kind:       Hero,
			hero:       true,
			multiplier: 1.0,
		},
		{
			name:       "with kind hero/scatter",
			id:         6,
			opts:       []SymbolOption{WithKind(HeroScatter)},
			kind:       HeroScatter,
			hero:       true,
			scatter:    true,
			multiplier: 1.0,
		},
		{
			name:       "with kind scatter",
			id:         7,
			opts:       []SymbolOption{WithKind(Scatter)},
			kind:       Scatter,
			scatter:    true,
			multiplier: 1.0,
		},
		{
			name:       "with kind bomb",
			id:         8,
			opts:       []SymbolOption{WithKind(Bomb)},
			kind:       Bomb,
			bomb:       true,
			multiplier: 1.0,
		},
		{
			name:       "with kind shooter",
			id:         9,
			opts:       []SymbolOption{WithKind(Shooter)},
			kind:       Shooter,
			shooter:    true,
			multiplier: 1.0,
		},
		{
			name:       "with kind prize",
			id:         10,
			opts:       []SymbolOption{WithKind(Prize)},
			kind:       Prize,
			prize:      true,
			multiplier: 1.0,
		},
		{
			name:       "with kind wild/bomb",
			id:         11,
			opts:       []SymbolOption{WithKind(WildBomb)},
			kind:       WildBomb,
			wild:       true,
			bomb:       true,
			multiplier: 1.0,
		},
		{
			name:       "with weights",
			id:         12,
			opts:       []SymbolOption{WithWeights(25, 50, 25)},
			kind:       Standard,
			weights:    []float64{25, 50, 25},
			multiplier: 1.0,
		},
		{
			name:       "with payouts",
			id:         13,
			opts:       []SymbolOption{WithPayouts(0, 0, 0.5, 2, 7.5)},
			kind:       Standard,
			payouts:    []float64{0, 0, 0.5, 2, 7.5},
			payable:    utils.UInt8s{3, 4, 5, 6},
			minPayable: 3,
			multiplier: 1.0,
		},
		{
			name:       "with wild for",
			id:         14,
			opts:       []SymbolOption{WildFor(0, 1)},
			kind:       Split,
			wildFor:    utils.Indexes{0, 1},
			split:      true,
			multiplier: 1.0,
		},
		{
			name:       "sticky",
			id:         15,
			opts:       []SymbolOption{IsSticky()},
			kind:       Standard,
			sticky:     true,
			multiplier: 1.0,
		},
		{
			name:       "with multiplier",
			id:         16,
			opts:       []SymbolOption{WithMultiplier(4.0)},
			kind:       Standard,
			multiplier: 4.0,
		},
		{
			name:       "morph",
			id:         17,
			opts:       []SymbolOption{MorphInto(11)},
			kind:       Standard,
			morphInto:  11,
			multiplier: 1.0,
		},
		{
			name:       "vary multiplier",
			id:         18,
			opts:       []SymbolOption{VaryMultiplier()},
			kind:       Standard,
			varyMult:   true,
			multiplier: 1.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSymbol(tc.id, tc.opts...)
			require.NotNil(t, s)

			assert.Equal(t, tc.id, s.ID())
			assert.Equal(t, tc.symbol, s.Name())
			assert.Equal(t, tc.resource, s.Resource())
			assert.Equal(t, tc.kind, s.Kind())
			assert.Equal(t, tc.sticky, s.IsSticky())
			assert.Equal(t, tc.split, s.IsSplit())
			assert.Equal(t, tc.wild, s.IsWild())
			assert.Equal(t, tc.hero, s.IsHero())
			assert.Equal(t, tc.scatter, s.IsScatter())
			assert.Equal(t, tc.bomb, s.IsBomb())
			assert.Equal(t, tc.shooter, s.IsShooter())
			assert.Equal(t, tc.prize, s.IsPrize())
			assert.Equal(t, tc.varyMult, s.VaryMultiplier())
			assert.Equal(t, tc.multiplier, s.Multiplier())
			assert.Equal(t, tc.morphInto, s.morphInto)

			if len(tc.weights) > 0 {
				assert.EqualValues(t, tc.weights, s.Weights())
			}

			if len(tc.payouts) > 0 {
				assert.EqualValues(t, tc.payouts, s.Payouts())
				for ix := range tc.payouts {
					assert.Equal(t, tc.payouts[ix], s.Payout(uint8(ix)+1))
				}
				assert.Zero(t, s.Payout(0))
				assert.Equal(t, tc.payouts[len(tc.payouts)-1], s.Payout(9))
			}

			if len(tc.payable) > 0 {
				for _, c := range tc.payable {
					assert.True(t, s.Payable(c))
				}
				assert.False(t, s.Payable(0))
				assert.False(t, s.Payable(1))
			}

			if tc.minPayable > 0 {
				assert.Equal(t, tc.minPayable, s.MinPayable())
			}

			if len(tc.wildFor) > 0 {
				for _, f := range tc.wildFor {
					assert.True(t, s.WildFor(f))
				}
				assert.False(t, s.WildFor(99))
			}
		})
	}
}

func TestNewSymbolFail(t *testing.T) {
	t.Run("new symbol fail", func(t *testing.T) {
		defer func() {
			e := recover()
			require.NotNilf(t, e, "new symbold should fail")
		}()
		NewSymbol(MaxSymbolID + 1)
	})
}
