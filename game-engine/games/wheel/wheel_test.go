package wheel

import (
	"math"
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestAcquireBonusWheel(t *testing.T) {
	testCases := []struct {
		name    string
		indexes utils.Indexes
		weights []float64
	}{
		{"one", utils.Indexes{11}, []float64{1}},
		{"two", utils.Indexes{1, 2}, []float64{50, 10}},
		{"three", utils.Indexes{1, 2, 3}, []float64{30, 20, 10}},
		{"six", utils.Indexes{1, 2, 3, 4, 5, 6}, []float64{100, 50, 30, 15, 5, 1}},
	}

	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := utils.AcquireWeighting().AddWeights(tc.indexes, tc.weights)
			require.NotNil(t, w)
			defer w.Release()

			bw := AcquireBonusWheel(prng, w)
			require.NotNil(t, bw)
			defer bw.Release()

			counts := make(map[utils.Index]int, 32)
			for ix := 0; ix < 10000; ix++ {
				_, data := bw.Run(nil)
				require.NotNil(t, data)

				wr, ok := data.(*BonusWheelResult)
				require.True(t, ok)
				require.NotNil(t, wr)

				counts[wr.result] = counts[wr.result] + 1

				data.Release()
			}

			n1 := math.MaxInt
			for ix := range tc.indexes {
				n2 := counts[tc.indexes[ix]]
				assert.NotZero(t, n2)
				assert.Less(t, n2, n1)
				n1 = n2
			}
		})
	}
}
