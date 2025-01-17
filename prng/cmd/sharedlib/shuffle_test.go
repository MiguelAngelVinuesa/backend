package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFisherYatesShuffle(t *testing.T) {
	t.Run("FisherYates", func(t *testing.T) {
		r := NewRNG()
		require.NotNil(t, r)
		defer FreeRNG(r)

		cards := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		FisherYatesShuffle(&cards[0], len(cards), r)

		counts := make(map[int]int)
		var moved int
		for ix, v := range cards {
			assert.Greater(t, v, 0)
			assert.LessOrEqual(t, v, 10)

			if ix+1 != v {
				moved++
			}

			counts[v] = counts[v] + 1
		}

		require.GreaterOrEqual(t, moved, 5)
		require.Equal(t, 10, len(counts))

		for k, v := range counts {
			assert.Equal(t, v, 1, k)
		}
	})
}
