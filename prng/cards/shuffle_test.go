package cards

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"
)

func TestFisherYatesShuffle(t *testing.T) {
	t.Run("FisherYates", func(t *testing.T) {
		cards := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52}
		r := rng.NewRNG()
		defer r.ReturnToPool()

		FisherYatesShuffle(cards, r)

		counts := make(map[int]int)
		var moved int
		for ix, v := range cards {
			assert.Greater(t, v, 0)
			assert.LessOrEqual(t, v, 52)

			if ix+1 != v {
				moved++
			}

			counts[v] = counts[v] + 1
		}

		require.GreaterOrEqual(t, moved, 40)
		require.Equal(t, 52, len(counts))

		for k, v := range counts {
			assert.Equal(t, v, 1, k)
		}
	})
}

func BenchmarkFisherYatesShuffle(b *testing.B) {
	cards := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52}
	r := rng.NewRNG()
	defer r.ReturnToPool()

	for i := 0; i < b.N; i++ {
		FisherYatesShuffle(cards, r)
	}
}
