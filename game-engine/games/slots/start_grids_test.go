package slots

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewStartGridsCCB(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	ss := make([]*slots.Symbol, 13)
	for ix := range ss {
		ss[ix] = slots.NewSymbol(utils.Index(ix+1), slots.WithWeights(20, 20, 20, 20, 20))
	}
	symbols := slots.NewSymbolSet(ss...)
	g := slots.NewSlots(slots.Grid(5, 3), slots.WithSymbols(symbols))

	a1 := slots.NewBestSymbolStickyAction()

	t.Run("CCB", func(t *testing.T) {
		f, err := os.Open("../../bin/ccb/startgrids.bin")
		require.NoError(t, err)

		sg := NewStartGrids(f, true)
		require.NotNil(t, sg)
		defer sg.Release()

		spin := slots.AcquireSpin(g, prng)
		defer spin.Release()
		spin.SetSpinner(sg)

		counts := make(map[int]int, 100000)

		for ix := 0; ix < 100000; ix++ {
			// first spin
			spin.ResetSpin()
			spin.Spin()

			lastGrid := sg.LastGrid()
			assert.GreaterOrEqual(t, lastGrid, 0)
			assert.Less(t, lastGrid, 500000)

			counts[lastGrid] = counts[lastGrid] + 1

			a1.FindBestSymbol(spin, nil)

			// second spin
			spin.Spin()
			assert.Equal(t, lastGrid, sg.LastGrid())
		}

		for ix, c := range counts {
			assert.NotZero(t, c, ix)
			assert.LessOrEqual(t, c, 6, ix)
		}
	})
}
