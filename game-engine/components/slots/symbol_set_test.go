package slots

import (
	"math"
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	s0  = NewSymbol(0)
	s1  = NewSymbol(1)
	s2  = NewSymbol(2)
	s3  = NewSymbol(3)
	wf1 = NewSymbol(4, WildFor(0, 1))
	w1  = NewSymbol(5, WithKind(Wild))
	h1  = NewSymbol(6, WithKind(Hero))
	sc1 = NewSymbol(7, WithKind(Scatter))
	st1 = NewSymbol(8, IsSticky())
	st2 = NewSymbol(9, IsSticky())
	st3 = NewSymbol(10, IsSticky())
	st4 = NewSymbol(11, IsSticky())
	st5 = NewSymbol(12, IsSticky())
	st6 = NewSymbol(13, IsSticky())
	w2  = NewSymbol(14, WithKind(Wild))
	h2  = NewSymbol(15, WithKind(Hero))
	sc2 = NewSymbol(16, WithKind(Scatter))
	b1  = NewSymbol(17, WithKind(Bomb), ClearPattern(-1, 1))
	b2  = NewSymbol(18, WithKind(Bomb), ClearPattern(-4, -3, -2, -1, 1, 2, 3, 4))
)

func TestNewSymbolSet(t *testing.T) {
	testCases := []struct {
		name      string
		symbols   Symbols
		splits    []utils.Index
		wilds     []utils.Index
		heroes    []utils.Index
		scatters  []utils.Index
		sticky    []utils.Index
		nonSticky []utils.Index
		withBomb  bool
		maxID     utils.Index
	}{
		{
			name:    "empty",
			symbols: Symbols{},
		},
		{
			name:    "single",
			symbols: Symbols{s0},
			maxID:   s0.ID(),
		},
		{
			name:    "multi",
			symbols: Symbols{s0, s1, s2, s3},
			maxID:   s3.ID(),
		},
		{
			name:    "multi with split",
			symbols: Symbols{s0, s1, wf1, s3},
			splits:  []utils.Index{wf1.ID()},
			maxID:   wf1.ID(),
		},
		{
			name:    "multi with wild",
			symbols: Symbols{s0, s1, s2, s3, w1, w2},
			wilds:   []utils.Index{w1.ID(), w2.ID()},
			maxID:   w2.ID(),
		},
		{
			name:    "multi with hero",
			symbols: Symbols{s0, s1, s2, s3, h1, h2},
			heroes:  []utils.Index{h1.ID(), h2.ID()},
			maxID:   h2.ID(),
		},
		{
			name:     "multi with scatter",
			symbols:  Symbols{s0, s1, s2, s3, sc1, sc2},
			scatters: []utils.Index{sc1.ID(), sc2.ID()},
			maxID:    sc2.ID(),
		},
		{
			name:     "multi with bomb",
			symbols:  Symbols{s0, s1, s2, s3, b1, b2},
			withBomb: true,
			maxID:    b2.ID(),
		},
		{
			name:      "one sticky",
			symbols:   Symbols{s0, s1, s2, s3, st1},
			nonSticky: []utils.Index{s0.ID(), s1.ID(), s2.ID(), s3.ID()},
			sticky:    []utils.Index{st1.ID()},
			maxID:     st1.ID(),
		},
		{
			name:      "50/50 sticky",
			symbols:   Symbols{s0, s1, s2, s3, st1, st2, st3, st4},
			nonSticky: []utils.Index{s0.ID(), s1.ID(), s2.ID(), s3.ID()},
			sticky:    []utils.Index{st1.ID(), st2.ID(), st3.ID(), st4.ID()},
			maxID:     st4.ID(),
		},
		{
			name:      "one non-sticky",
			symbols:   Symbols{s0, st1, st2, st3, st4, st5, st6},
			nonSticky: []utils.Index{s0.ID()},
			sticky:    []utils.Index{st1.ID(), st2.ID(), st3.ID(), st4.ID(), st5.ID(), st6.ID()},
			maxID:     st6.ID(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSymbolSet(tc.symbols...)
			require.NotNil(t, s)

			assert.EqualValues(t, tc.symbols, s.symbols)

			for _, index := range tc.splits {
				sym := s.GetSymbol(index)
				assert.True(t, sym.IsSplit())
			}
			for _, index := range tc.wilds {
				sym := s.GetSymbol(index)
				assert.True(t, sym.IsWild())
			}
			for _, index := range tc.heroes {
				sym := s.GetSymbol(index)
				assert.True(t, sym.IsHero())
			}
			for _, index := range tc.scatters {
				sym := s.GetSymbol(index)
				assert.True(t, sym.IsScatter())
			}
			for _, index := range tc.sticky {
				sym := s.GetSymbol(index)
				assert.True(t, sym.IsSticky())
			}
			for _, index := range tc.nonSticky {
				sym := s.GetSymbol(index)
				assert.False(t, sym.IsSticky())
			}

			assert.Equal(t, tc.withBomb, s.haveBombs)
			assert.Equal(t, tc.maxID, s.GetMaxSymbolID())

			assert.Nil(t, s.GetSymbol(99))
		})
	}
}

func TestSymbolSet_GetSymbol(t *testing.T) {
	t.Run("get symbol", func(t *testing.T) {
		s := NewSymbolSet(s1, s2, s3, s0, wf1, w1, h1, sc1)
		require.NotNil(t, s)

		i := s.GetSymbol(0)
		require.NotNil(t, i)
		assert.Equal(t, s0, i)

		i = s.GetSymbol(5)
		require.NotNil(t, i)
		assert.Equal(t, w1, i)

		i = s.GetSymbol(7)
		require.NotNil(t, i)
		assert.Equal(t, sc1, i)

		i = s.GetSymbol(9)
		require.Nil(t, i)
	})
}

func TestNewSymbolSetFail(t *testing.T) {
	t.Run("new symbol set fail", func(t *testing.T) {
		defer func() {
			if e := recover(); e == nil {
				t.Error("symbols.init should fail")
			}
		}()

		s := NewSymbolSet(s1, s2, s3, s1)
		require.NotNil(t, s)
	})
}

func TestSymbolSet_BonusWeights(t *testing.T) {
	t.Run("bonus weights", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		s := NewSymbolSet(s1, s2, s3, wf1, w1, h1, sc1)
		require.NotNil(t, s)

		w := utils.AcquireWeighting()
		require.NotNil(t, w)
		w.AddWeight(s1.id, 50)
		w.AddWeight(s2.id, 40)
		w.AddWeight(s3.id, 30)
		w.AddWeight(wf1.id, 20)
		w.AddWeight(h1.id, 10)

		s.SetBonusWeights(w)

		counts := make(map[utils.Index]int)
		for ix := 0; ix < 10000; ix++ {
			n := s.GetBonusSymbol(prng)
			counts[n] = counts[n] + 1
		}

		assert.NotZero(t, counts[s1.id])
		assert.NotZero(t, counts[s2.id])
		assert.NotZero(t, counts[s3.id])
		assert.NotZero(t, counts[wf1.id])
		assert.NotZero(t, counts[h1.id])
		assert.Zero(t, counts[0])
		assert.Zero(t, counts[w1.id])
		assert.Zero(t, counts[sc1.id])
	})
}

func TestSymbolSet_HighestPaying(t *testing.T) {
	t.Run("highest paying", func(t *testing.T) {
		s := NewSymbolSet(sf1, sf2, sf3, sf4, sf5, sf6, wff1, wwf1, hf1, scf1, scf2)
		require.NotNil(t, s)

		assert.Equal(t, []utils.Index{math.MaxUint16, math.MaxUint16, 14, 14, 14, 14}, s.bestWildSym)
		assert.Equal(t, []float64{0, 0, 1.5, 3, 5, 12}, s.bestWildPay)
	})
}
