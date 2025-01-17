package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHash(t *testing.T) {
	t.Run("get hash", func(t *testing.T) {
		s := GetHash()
		require.NotNil(t, s)
	})
}

func TestNewRNG(t *testing.T) {
	t.Run("new RNG", func(t *testing.T) {
		r := NewRNG()
		require.NotNil(t, r)
		defer FreeRNG(r)
	})
}

func TestGetUInt32(t *testing.T) {
	t.Run("get uint32", func(t *testing.T) {
		r := NewRNG()
		require.NotNil(t, r)
		defer FreeRNG(r)

		counts := make(map[uint32]int)
		for ix := 0; ix < 100000; ix++ {
			n := GetUInt32(r)
			counts[n] = counts[n] + 1
		}

		for k, v := range counts {
			assert.GreaterOrEqual(t, k, uint32(0))
			assert.LessOrEqual(t, v, 2, k)
		}
	})
}

func TestGetUInt64(t *testing.T) {
	t.Run("get uint64", func(t *testing.T) {
		r := NewRNG()
		require.NotNil(t, r)
		defer FreeRNG(r)

		counts := make(map[uint64]int)
		for ix := 0; ix < 100000; ix++ {
			n := GetUInt64(r)
			counts[n] = counts[n] + 1
		}

		for k, v := range counts {
			assert.GreaterOrEqual(t, k, uint64(0))
			assert.LessOrEqual(t, v, 2, k)
		}
	})
}

func TestGetIntN(t *testing.T) {
	t.Run("get intN", func(t *testing.T) {
		r := NewRNG()
		require.NotNil(t, r)
		defer FreeRNG(r)

		counts := make(map[int]int)
		for ix := 0; ix < 100000; ix++ {
			n := GetIntN(r, 100)
			counts[n] = counts[n] + 1
		}

		for k, v := range counts {
			assert.GreaterOrEqual(t, k, 0)
			assert.Less(t, k, 100)
			assert.LessOrEqual(t, v, 2000, k)
		}
	})
}

func TestGetIntsN(t *testing.T) {
	t.Run("get intsN", func(t *testing.T) {
		r := NewRNG()
		require.NotNil(t, r)
		defer FreeRNG(r)

		out := make([]int, 5*3)
		counts := make(map[int]int)
		for ix := 0; ix < 1000; ix++ {
			GetIntsN(r, 100, &out[0], len(out))
			for iy := range out {
				n := out[iy]
				counts[n] = counts[n] + 1
			}
		}

		for k, v := range counts {
			assert.GreaterOrEqual(t, k, 0)
			assert.Less(t, k, 100)
			assert.LessOrEqual(t, v, 300, k)
		}
	})
}
