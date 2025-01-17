package ber

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	analysis "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/analysis/game/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/rng"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"
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
		{name: "RTP 92", rtp: 92, min: 50.0, max: 200.0},
		{name: "RTP 93", rtp: 93, fail: true},
		{name: "RTP 94", rtp: 94, min: 50.0, max: 200.0},
		{name: "RTP 95", rtp: 95, fail: true},
		{name: "RTP 96", rtp: 96, min: 50.0, max: 200.0},
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

				r := analysis.AcquireRounds(tg.BERnr, "p1", 10000000, reels, rows, false, g.MaxPayout(), symbols, AllActions(tc.rtp), nil, flags)
				require.NotNil(t, r)
				defer r.Release()

				for ix := 0; ix < 100000; ix++ {
					results := g.Round(0)
					r.Analyse(100, 100, results)
				}

				assert.GreaterOrEqual(t, r.RTP(), tc.min)
				// assert.LessOrEqual(t, r.RTP(), tc.max)
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

func BenchmarkRound94(b *testing.B) {
	g := New(94)
	defer g.Release()

	for i := 0; i < b.N; i++ {
		g.Round(0)
	}
}

func BenchmarkRound92Logged(b *testing.B) {
	g := NewLogged(92)
	defer g.Release()

	for i := 0; i < b.N; i++ {
		g.Round(0)
	}
}

func BenchmarkRound94Logged(b *testing.B) {
	g := NewLogged(94)
	defer g.Release()

	for i := 0; i < b.N; i++ {
		g.Round(0)
	}
}
