package cas

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	rng2 "git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	analysis "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/analysis/game/slots"
	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/rng"
)

func init() {
	rng.AcquireRNG = func() interfaces.Generator { return rng2.NewRNG() }
}

func TestNew(t *testing.T) {
	testCases := []struct {
		name    string
		rtp     int
		actions comp.SpinActions
		fail    bool
		min     float64
		max     float64
	}{
		{name: "RTP 90", rtp: 90, fail: true},
		{name: "RTP 91", rtp: 91, fail: true},
		{name: "RTP 92", rtp: 92, fail: true},
		{name: "RTP 93", rtp: 93, fail: true},
		{name: "RTP 94", rtp: 94, fail: true},
		{name: "RTP 95", rtp: 95, fail: true},
		{name: "RTP 96", rtp: 96, actions: actions96all, min: 0.001, max: 200.0},
		{name: "RTP 97", rtp: 97, fail: true},
		{name: "RTP 98", rtp: 98, fail: true},
		{name: "RTP 99", rtp: 99, fail: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := New(tc.rtp)
			if tc.fail {
				require.Nil(t, g)
			} else {
				require.NotNil(t, g)
				defer g.Release()

				r := analysis.AcquireRounds(tg.BOTnr, "p1", 10000000, reels, rows, false, g.MaxPayout(), symbols1, tc.actions, Paylines(), Flags())
				require.NotNil(t, r)
				defer r.Release()

				for ix := 0; ix < 100000; ix++ {
					results := g.Round(0)
					r.Analyse(100, 100, results)
				}

				assert.GreaterOrEqual(t, r.RTP(), tc.min)
				assert.LessOrEqual(t, r.RTP(), tc.max)
			}
		})
	}
}

func TestAllSymbols(t *testing.T) {
	t.Run("all symbols", func(t *testing.T) {
		a := AllSymbols()
		require.NotNil(t, a)
		assert.Equal(t, symbols1, a)
	})
}

func TestPaylines(t *testing.T) {
	t.Run("paylines", func(t *testing.T) {
		p := Paylines()
		require.NotNil(t, p)
		assert.Equal(t, paylines, p)
	})
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := New(96)
		g.Release()
	}
}

func BenchmarkNewAndRound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := New(96)
		g.Round(0)
		g.Release()
	}
}

func BenchmarkRound96(b *testing.B) {
	g := New(96)
	defer g.Release()

	for i := 0; i < b.N; i++ {
		g.Round(0)
	}
}

func BenchmarkAnalyze96(b *testing.B) {
	g := New(96)
	defer g.Release()

	for i := 0; i < b.N; i++ {
		a := analysis.AcquireRounds(tg.BOTnr, "x", 1000000, 5, 3, false, g.MaxPayout(), symbols1, actions96all, paylines, flags)
		for iy := 0; iy < 1000; iy++ {
			res := g.Round(0)
			a.Analyse(1000, 100, res)
		}
		a.Release()
	}
}

func BenchmarkRound96Logged(b *testing.B) {
	g := NewLogged(96)
	defer g.Release()

	for i := 0; i < b.N; i++ {
		g.Round(0)
	}
}
