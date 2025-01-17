package cards

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"
)

// FisherYatesShuffle shuffles the slice into a random order.
// See https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
func FisherYatesShuffle(c []int, r *rng.RNG) {
	// Our ChaCha based prng has a 256-bit cycle, so it can shuffle upto 57 cards and produce every possible permutation.
	for a := len(c) - 1; a > 0; a-- {
		b := r.IntN(a + 1)
		c[a], c[b] = c[b], c[a]
	}
}
