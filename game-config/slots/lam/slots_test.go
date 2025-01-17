package lam

import (
	"math/rand"
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	analysis "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/analysis/game/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/rng"
	rng2 "git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"
)

func init() {
	rng.AcquireRNG = func() interfaces.Generator { return rng2.NewRNG() }
}

func TestNew(t *testing.T) {
	testCases := []struct {
		name string
		rtp  int
		fail bool
		min  float64
		max  float64
	}{
		{name: "RTP 90", rtp: 90, fail: true},
		{name: "RTP 91", rtp: 91, fail: true},
		{name: "RTP 92", rtp: 92, min: 0.0, max: 200.0},
		{name: "RTP 93", rtp: 93, fail: true},
		{name: "RTP 94", rtp: 94, min: 0.0, max: 200.0},
		{name: "RTP 95", rtp: 95, fail: true},
		{name: "RTP 96", rtp: 96, min: 0.0, max: 200.0},
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

				r := analysis.AcquireRounds(tg.LAMnr, "p1", 10000000, reels, rows, false, g.MaxPayout(), AllSymbols(), AllActions(tc.rtp), Paylines(), Flags())
				require.NotNil(t, r)
				defer r.Release()

				var res results.Results

				for ix := 0; ix < 100000; ix++ {
					r2 := g.Round(0)
					for iy := range r2 {
						res = append(res, r2[iy].Clone().(*results.Result))
					}

					if g.NeedPlayerChoice() {
						// player choice required!
						if rand.Int31n(10000) < 5000 {
							r2 = g.RoundResume(map[string]string{"wing": "south"})
						} else {
							r2 = g.RoundResume(map[string]string{"wing": "north"})
						}

						for iy := range r2 {
							res = append(res, r2[iy].Clone().(*results.Result))
						}
					}

					r.Analyse(100, 100, res)
					res = results.ReleaseResults(res)
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
		assert.Equal(t, symbols, a)
	})
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := New(96)
		g.Release()
	}
}

func BenchmarkRound92(b *testing.B) {
	g := New(92)
	defer g.Release()

	for i := 0; i < b.N; i++ {
		g.Round(0)
	}
}

func BenchmarkRound96(b *testing.B) {
	g := New(96)
	defer g.Release()

	for i := 0; i < b.N; i++ {
		g.Round(0)
	}
}

func BenchmarkRound92logged(b *testing.B) {
	g := NewLogged(92)
	defer g.Release()

	for i := 0; i < b.N; i++ {
		g.Round(0)
	}
}

func BenchmarkRound96logged(b *testing.B) {
	g := NewLogged(96)
	defer g.Release()

	for i := 0; i < b.N; i++ {
		g.Round(0)
	}
}
