package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	sw1  = NewSymbol(1, WithWeights(80, 70, 60, 50, 40))
	sw2  = NewSymbol(2, WithWeights(70, 60, 50, 40, 30))
	sw3  = NewSymbol(3, WithWeights(60, 50, 40, 30, 20))
	sw4  = NewSymbol(4, WithWeights(50, 40, 30, 20, 10))
	wfw1 = NewSymbol(5, WildFor(0, 1), WithWeights(5, 4, 3, 4, 5))
	wfw2 = NewSymbol(6, WildFor(2, 3), WithWeights(5, 4, 3, 4, 5))
	ww1  = NewSymbol(7, WithKind(Wild), WithWeights(5, 4, 3, 4, 5))
	hw1  = NewSymbol(8, WithKind(Hero), WithWeights(5, 4, 3, 4, 5))
	scw1 = NewSymbol(9, WithKind(Scatter), WithWeights(5, 4, 3, 4, 5))

	set1 = NewSymbolSet(sw1, sw2, sw3, sw4)
	set2 = NewSymbolSet(sw1, sw2, sw3, sw4, wfw1, wfw2)
	set3 = NewSymbolSet(sw1, sw2, sw3, sw4, wfw1, wfw2, hw1, ww1, scw1)
)

func TestNewReel(t *testing.T) {
	testCases := []struct {
		name    string
		index   uint8
		rows    uint8
		symbols *SymbolSet
	}{
		{
			name:    "simple",
			index:   0,
			rows:    3,
			symbols: set1,
		},
		{
			name:    "with splits",
			index:   1,
			rows:    3,
			symbols: set2,
		},
		{
			name:    "with all",
			index:   2,
			rows:    5,
			symbols: set3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prng := rng.NewRNG()
			defer prng.ReturnToPool()

			r := AcquireReel(tc.index, tc.rows, tc.symbols, tc.rows-1)
			require.NotNil(t, r)
			defer r.Release()

			assert.Equal(t, tc.index, r.index)
			assert.EqualValues(t, tc.rows, r.rows)
			assert.Equal(t, tc.symbols, r.symbols)

			work := make([]utils.Index, tc.rows)
			for ix := 0; ix < 1000; ix++ {
				r.Spin(prng, work)
				require.NotNil(t, work)
				for iy, n := range work {
					assert.NotZero(t, n)
					work[iy] = 0
				}
			}
		})
	}
}
