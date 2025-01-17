package sharedlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRNG(t *testing.T) {
	t.Run("new RNG", func(t *testing.T) {
		r := AcquireRNG()
		require.NotNil(t, r)
		defer r.ReturnToPool()
	})
}

func TestRNG_UInt32(t *testing.T) {
	t.Run("RNG uint32", func(t *testing.T) {
		r := AcquireRNG()
		require.NotNil(t, r)
		defer r.ReturnToPool()

		max := 100000
		counts := make(map[uint32]int, max)
		for ix := 0; ix < max; ix++ {
			n := r.Uint32()
			counts[n] = counts[n] + 1
		}

		for k, v := range counts {
			assert.LessOrEqual(t, v, 2, k)
		}
	})
}

func TestRNG_UInt64(t *testing.T) {
	t.Run("RNG uint64", func(t *testing.T) {
		r := AcquireRNG()
		require.NotNil(t, r)
		defer r.ReturnToPool()

		max := 100000
		counts := make(map[uint64]int, max)
		for ix := 0; ix < max; ix++ {
			n := r.Uint64()
			counts[n] = counts[n] + 1
		}

		for k, v := range counts {
			assert.LessOrEqual(t, v, 2, k)
		}
	})
}

func TestRNG_IntN(t *testing.T) {
	t.Run("RNG intN", func(t *testing.T) {
		r := AcquireRNG()
		require.NotNil(t, r)
		defer r.ReturnToPool()

		max := 100000
		counts := make(map[int]int, 2000)
		for ix := 0; ix < max; ix++ {
			n := r.IntN(123)
			counts[n] = counts[n] + 1
		}

		for k, v := range counts {
			assert.GreaterOrEqual(t, k, 0)
			assert.Less(t, k, 123)
			assert.LessOrEqual(t, v, 1500, k)
		}
	})
}

func TestRNG_IntsN(t *testing.T) {
	t.Run("RNG intsN", func(t *testing.T) {
		r := AcquireRNG()
		require.NotNil(t, r)
		defer r.ReturnToPool()

		out := make([]int, 100)
		max := 1000
		counts := make(map[int]int, 2000)
		for ix := 0; ix < max; ix++ {
			r.IntsN(123, out)
			for iy := range out {
				n := out[iy]
				counts[n] = counts[n] + 1
			}
		}

		for k, v := range counts {
			assert.GreaterOrEqual(t, k, 0)
			assert.Less(t, k, 123)
			assert.LessOrEqual(t, v, 1500, k)
		}
	})
}

func BenchmarkNewRNG(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := AcquireRNG()
		r.ReturnToPool()
	}
}

func BenchmarkRNG_UInt32(b *testing.B) {
	r := AcquireRNG()
	defer r.ReturnToPool()

	for i := 0; i < b.N; i++ {
		r.Uint32()
	}
}

func BenchmarkRNG_UInt64(b *testing.B) {
	r := AcquireRNG()
	defer r.ReturnToPool()

	for i := 0; i < b.N; i++ {
		r.Uint64()
	}
}

func BenchmarkRNG_IntN_Small(b *testing.B) {
	r := AcquireRNG()
	defer r.ReturnToPool()

	n := 9
	for i := 0; i < b.N; i++ {
		r.IntN(n)
	}
}

func BenchmarkRNG_IntN_Medium(b *testing.B) {
	r := AcquireRNG()
	defer r.ReturnToPool()

	n := 1000000
	for i := 0; i < b.N; i++ {
		r.IntN(n)
	}
}

func BenchmarkRNG_IntN_Large(b *testing.B) {
	r := AcquireRNG()
	defer r.ReturnToPool()

	n := 1<<62 + 1
	for i := 0; i < b.N; i++ {
		r.IntN(n)
	}
}

func BenchmarkRNG_IntsN_Small_15(b *testing.B) {
	r := AcquireRNG()
	defer r.ReturnToPool()

	n := 9
	out := make([]int, 15)
	for i := 0; i < b.N; i++ {
		r.IntsN(n, out)
	}
}

func BenchmarkRNG_IntsN_Medium_15(b *testing.B) {
	r := AcquireRNG()
	defer r.ReturnToPool()

	n := 1000000
	out := make([]int, 15)
	for i := 0; i < b.N; i++ {
		r.IntsN(n, out)
	}
}

func BenchmarkRNG_IntsN_Large_15(b *testing.B) {
	r := AcquireRNG()
	defer r.ReturnToPool()

	n := 1<<62 + 1
	out := make([]int, 15)
	for i := 0; i < b.N; i++ {
		r.IntsN(n, out)
	}
}

func BenchmarkRNG_IntsN_Small_100(b *testing.B) {
	r := AcquireRNG()
	defer r.ReturnToPool()

	n := 9
	out := make([]int, 100)
	for i := 0; i < b.N; i++ {
		r.IntsN(n, out)
	}
}

func BenchmarkRNG_IntsN_Medium_100(b *testing.B) {
	r := AcquireRNG()
	defer r.ReturnToPool()

	n := 1000000
	out := make([]int, 100)
	for i := 0; i < b.N; i++ {
		r.IntsN(n, out)
	}
}

func BenchmarkRNG_IntsN_Large_100(b *testing.B) {
	r := AcquireRNG()
	defer r.ReturnToPool()

	n := 1<<62 + 1
	out := make([]int, 100)
	for i := 0; i < b.N; i++ {
		r.IntsN(n, out)
	}
}
