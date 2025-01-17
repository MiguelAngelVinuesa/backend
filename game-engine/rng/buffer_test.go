package rng

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBuffer(t *testing.T) {
	testCases := []struct {
		name    string
		keys    []int
		counts  []int
		max     int
		buffers int
		reloads int
	}{
		{
			name:    "simple",
			keys:    []int{10000},
			counts:  []int{1000},
			max:     5,
			buffers: 1,
			reloads: 32,
		},
		{
			name:    "several",
			keys:    []int{10000, 5100, 7000, 9300},
			counts:  []int{100, 200, 200, 200},
			max:     5,
			buffers: 4,
			reloads: 25,
		},
		{
			name:    "many",
			keys:    []int{10000, 5100, 7000, 9300, 2560, 8100, 4490, 12341, 7777, 20000},
			counts:  []int{100, 200, 200, 200, 300, 200, 300, 50, 75, 150},
			max:     10,
			buffers: 10,
			reloads: 62,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prng := AcquireRNG()

			buf := AcquireBuffer(prng, false)
			require.NotNil(t, buf)
			defer buf.Release()

			assert.Equal(t, prng, buf.prng)
			assert.NotNil(t, buf.buffers)
			assert.Nil(t, buf.inputs)
			assert.Nil(t, buf.outputs)

			counts := make(map[int]int, 100)

			for ix := range tc.keys {
				key := tc.keys[ix]
				count := tc.counts[ix]

				for iy := 0; iy < count; iy++ {
					n := buf.IntN(key)
					assert.GreaterOrEqual(t, n, 0)
					assert.Less(t, n, key)
					counts[n] = counts[n] + 1
				}
			}

			assert.Equal(t, tc.buffers, buf.bufferCount())
			assert.Equal(t, tc.reloads, buf.reloadCount())

			l1, l2 := buf.Log()
			assert.Nil(t, l1)
			assert.Nil(t, l2)

			for ix := range counts {
				assert.NotZero(t, counts[ix])
				assert.LessOrEqual(t, counts[ix], tc.max)
			}
		})
	}
}

func TestNewBufferWithLog(t *testing.T) {
	t.Run("new buffer with log", func(t *testing.T) {
		prng := AcquireRNG()

		buf := AcquireBuffer(prng, true)
		require.NotNil(t, buf)
		defer buf.Release()

		n := 100000
		max := 10000
		counts := make(map[int]int, 100)
		for ix := 0; ix < max; ix++ {
			i := buf.IntN(n)
			counts[i] = counts[i] + 1
		}

		l1, l2 := buf.Log()
		require.NotNil(t, l1)
		require.NotNil(t, l2)
		assert.Equal(t, max, len(l1))
		assert.Equal(t, max, len(l2))

		for ix := range l1 {
			assert.Equal(t, n, l1[ix])
			assert.NotZero(t, counts[l2[ix]])
		}
	})
}

func TestBuffer_WithCache(t *testing.T) {
	testCases := []struct {
		name   string
		cache  []int
		input  []int
		output []int
		fail   bool
	}{
		{
			name: "empty",
		},
		{
			name:   "short input",
			cache:  []int{10000},
			input:  []int{10000},
			output: []int{10000},
			fail:   true,
		},
		{
			name:   "fail on first",
			cache:  []int{10000, 9000},
			input:  []int{10000},
			output: []int{9001},
			fail:   true,
		},
		{
			name:   "fail on middle",
			cache:  []int{10000, 9000, 11000, 8000, 12000, 7000, 13000, 6000, 14000, 5000},
			input:  []int{10000, 11000, 12000, 13000, 14000},
			output: []int{9000, 8000, 7001, 6000, 5000},
			fail:   true,
		},
		{
			name:   "fail on last",
			cache:  []int{10000, 9000, 11000, 8000, 12000, 7000, 13000, 6000, 14000, 5000},
			input:  []int{10000, 11000, 12000, 13000, 14000},
			output: []int{9000, 8000, 7000, 6000, 5001},
			fail:   true,
		},
		{
			name:   "fail cache too short",
			cache:  []int{10000, 9000, 11000, 8000, 12000, 7000, 13000, 6000, 14000, 5000},
			input:  []int{10000, 11000, 12000, 13000, 14000, 15000},
			output: []int{9000, 8000, 7000, 6000, 5000, 4000},
			fail:   true,
		},
		{
			name:   "good",
			cache:  []int{10000, 9000, 11000, 8000, 12000, 7000, 13000, 6000, 14000, 5000},
			input:  []int{10000, 11000, 12000, 13000, 14000},
			output: []int{9000, 8000, 7000, 6000, 5000},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prng := AcquireRNG()

			buf := AcquireBuffer(prng, false).WithCache(tc.cache)
			require.NotNil(t, buf)

			var fail bool
			for ix := range tc.input {
				got := buf.IntN(tc.input[ix])
				if got != tc.output[ix] {
					fail = true
				}
			}

			assert.Equal(t, tc.fail, fail)
		})
	}
}
