package utils

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWeighting(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		for ix := 0; ix < 10; ix++ {
			w := AcquireWeighting().(*WeightingNoDedup)
			require.NotNil(t, w)

			assert.NotNil(t, w.indexes)
			assert.NotNil(t, w.weights)
			assert.Zero(t, len(w.indexes))
			assert.Zero(t, len(w.weights))
			assert.Zero(t, w.total)

			w.Release()
		}
	})
}

func TestWeighting_AddWeight(t *testing.T) {
	testCases := []struct {
		name    string
		indexes Indexes
		weights []float64
		length  int
		total   int
	}{
		{"3 consecutive", []Index{1, 2, 3}, []float64{24.5, 9.5, 4.5}, 3, 38500},
		{"3 non-consecutive", []Index{1, 4, 7}, []float64{24.5, 9.5, 14.5}, 3, 48500},
		{"5 consecutive", []Index{0, 1, 2, 3, 4}, []float64{55.5, 24.5, 12, 9.5, 4.5}, 5, 106000},
		{"5 non-consecutive", []Index{0, 1, 3, 5, 6}, []float64{55.5, 24.5, 12, 9.5, 4.5}, 5, 106000},
		{"5 non-consecutive with zero weight", []Index{0, 1, 3, 5, 6}, []float64{55.5, 24.5, 0, 9.5, 4.5}, 4, 94000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, len(tc.indexes), len(tc.weights))

			w := AcquireWeighting().(*WeightingNoDedup)
			require.NotNil(t, w)
			defer w.Release()

			for ix, index := range tc.indexes {
				w2 := w.AddWeight(index, tc.weights[ix])
				require.Equal(t, w, w2)
			}

			assert.Equal(t, tc.length, len(w.indexes))
			assert.Equal(t, tc.length, len(w.weights))
			assert.Equal(t, tc.total, w.total)
		})
	}
}

func TestWeighting_AddWeights(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	testCases := []struct {
		name    string
		indexes Indexes
		weights []float64
		length  int
		total   int
	}{
		{"3 consecutive", []Index{1, 2, 3}, []float64{24.5, 9.5, 4.5}, 3, 38500},
		{"3 non-consecutive", []Index{1, 4, 7}, []float64{24.5, 9.5, 14.5}, 3, 48500},
		{"5 consecutive", []Index{0, 1, 2, 3, 4}, []float64{55.5, 24.5, 12, 9.5, 4.5}, 5, 106000},
		{"5 non-consecutive", []Index{0, 1, 3, 5, 6}, []float64{55.5, 24.5, 12, 9.5, 4.5}, 5, 106000},
		{"5 non-consecutive with zero weight", []Index{0, 1, 3, 5, 6}, []float64{55.5, 24.5, 0, 9.5, 4.5}, 4, 94000},
		{"scripted rounds scenarios (1)", []Index{1, 2, 3, 4, 9}, []float64{0.25, 0.25, 0.25, 0.25, 999}, 5, 1000000},
		{"scripted rounds scenarios (2)", []Index{1, 2, 3, 4, 9}, []float64{2, 2, 2, 2, 9992}, 5, 10000000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, len(tc.indexes), len(tc.weights))

			w := AcquireWeighting().(*WeightingNoDedup)
			require.NotNil(t, w)
			defer w.Release()

			w2 := w.AddWeights(tc.indexes, tc.weights)
			require.Equal(t, w, w2)

			assert.Equal(t, tc.length, len(w.indexes))
			assert.Equal(t, tc.length, len(w.weights))
			assert.Equal(t, tc.total, w.total)

			max := w.total / 25
			counts := make(map[Index]int, tc.length)
			for ix := 0; ix < max; ix++ {
				id := w2.RandomIndex(prng)
				counts[id] = counts[id] + 1
			}

			for ix, id := range tc.indexes {
				if tc.weights[ix] > 0 {
					assert.NotZero(t, counts[id])
				}
			}
		})
	}
}

func TestWeighting_RandomIndex(t *testing.T) {
	t.Run("random index", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		w := AcquireWeighting().(*WeightingNoDedup)
		require.NotNil(t, w)
		defer w.Release()

		data := map[Index]float64{
			0:  75,
			1:  50.25,
			2:  33.333,
			4:  25.5,
			5:  20,
			7:  10.75,
			8:  8,
			10: 3.666,
		}

		for k, v := range data {
			w.AddWeight(k, v)
		}

		assert.Equal(t, len(data), len(w.indexes))
		assert.Equal(t, len(data), len(w.weights))

		count := 1000000
		counts := make(map[Index]int, 8)
		for ix := 0; ix < count; ix++ {
			n := w.RandomIndex(prng)
			counts[n] = counts[n] + 1
		}

		for k, v := range counts {
			avg := float64(count) * math.Round(data[k]*1000) / float64(w.total)
			min := int(math.Round(avg * 0.7))
			max := int(math.Round(avg * 1.3))
			assert.GreaterOrEqual(t, v, min, k)
			assert.LessOrEqual(t, v, max, k)
		}
	})
}

func TestWeighting_RandomIndex2(t *testing.T) {
	t.Run("random index 2", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		w := AcquireWeighting().(*WeightingNoDedup)
		require.NotNil(t, w)
		defer w.Release()

		indexes := Indexes{1, 2}
		weights := []float64{1, 9999}
		w.AddWeights(indexes, weights)

		assert.Equal(t, len(indexes), len(w.indexes))
		assert.Equal(t, len(indexes), len(w.weights))

		count := 1000000
		counts := make(map[Index]int, 8)
		for ix := 0; ix < count; ix++ {
			n := w.RandomIndex(prng)
			counts[n] = counts[n] + 1
		}

		require.Equal(t, len(indexes), len(counts))

		for k, v := range counts {
			avg := float64(count) * math.Round(weights[k-1]*1000) / float64(w.total)
			min := int(math.Round(avg * 0.7))
			max := int(math.Round(avg * 1.3))
			assert.GreaterOrEqual(t, v, min, k)
			assert.LessOrEqual(t, v, max, k)
		}
	})
}

func TestWeighting_RandomIndexNoRepeat2(t *testing.T) {
	t.Run("random index", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		w := AcquireWeightingDedup3().(*WeightingDedup3)
		require.NotNil(t, w)
		defer w.Release()

		data := map[Index]float64{
			0:  75,
			1:  50.25,
			2:  33.333,
			4:  25.5,
			5:  20,
			7:  10.75,
			8:  8,
			10: 3.666,
		}

		for k, v := range data {
			w.AddWeight(k, v)
		}

		assert.Equal(t, len(data), len(w.indexes))
		assert.Equal(t, len(data), len(w.weights))

		count := 1000000
		counts := make(map[Index]int, 8)
		for ix := 0; ix < count; ix++ {
			n := w.RandomIndex(prng)
			counts[n] = counts[n] + 1
		}

		for k, v := range counts {
			avg := float64(count) * math.Round(data[k]*1000) / float64(w.total)
			min := int(math.Round(avg * 0.5))
			max := int(math.Round(avg * 2.2))
			assert.GreaterOrEqual(t, v, min, k)
			assert.LessOrEqual(t, v, max, k)
		}
	})
}

func TestWeighting_RandomIndexManyItems(t *testing.T) {
	t.Run("random index many items", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		w := AcquireWeighting().(*WeightingNoDedup)
		require.NotNil(t, w)
		defer w.Release()

		data := map[Index]float64{
			0:  75,
			1:  50.25,
			2:  33.333,
			4:  25.5,
			5:  20,
			7:  10.75,
			8:  8,
			10: 3.666,
			12: 33.333,
			14: 25.5,
			15: 20,
			17: 10.75,
			22: 33.333,
			24: 25.5,
			25: 20,
			27: 10.75,
		}

		for k, v := range data {
			w.AddWeight(k, v)
		}

		assert.Equal(t, len(data), len(w.indexes))
		assert.Equal(t, len(data), len(w.weights))

		count := 1000000
		counts := make(map[Index]int, 8)
		for ix := 0; ix < count; ix++ {
			n := w.RandomIndex(prng)
			counts[n] = counts[n] + 1
		}

		for k, v := range counts {
			avg := float64(count) * math.Round(data[k]*1000) / float64(w.total)
			min := int(math.Round(avg * 0.7))
			max := int(math.Round(avg * 1.3))
			assert.GreaterOrEqual(t, v, min, k)
			assert.LessOrEqual(t, v, max, k)
		}
	})
}

func TestWeighting_FillRandom(t *testing.T) {
	t.Run("fill random", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		w := AcquireWeighting().(*WeightingNoDedup)
		require.NotNil(t, w)
		defer w.Release()

		data := map[Index]float64{
			0:  75,
			1:  50.25,
			2:  33.333,
			4:  25.5,
			5:  20,
			7:  10.75,
			8:  8,
			10: 3.666,
		}

		for k, v := range data {
			w.AddWeight(k, v)
		}

		assert.Equal(t, len(data), len(w.indexes))
		assert.Equal(t, len(data), len(w.weights))

		count := 100000
		counts := make(map[Index]int, 8)
		tempLen := 13
		temp := make([]Index, tempLen)
		for ix := 0; ix < count; ix++ {
			w.FillRandom(prng, len(temp), temp)
			for _, n := range temp {
				counts[n] = counts[n] + 1
			}
		}

		for k, v := range counts {
			avg := float64(count) * float64(tempLen) * math.Round(data[k]*1000) / float64(w.total)
			min := int(math.Round(avg * 0.9))
			max := int(math.Round(avg * 1.1))
			assert.GreaterOrEqual(t, v, min, k)
			assert.LessOrEqual(t, v, max, k)
		}
	})
}

func TestWeighting_FillRandomDedup3(t *testing.T) {
	t.Run("fill random dedup3", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		w := AcquireWeightingDedup3().(*WeightingDedup3)
		require.NotNil(t, w)
		defer w.Release()

		data := map[Index]float64{
			1:  50.25,
			2:  33.333,
			4:  25.5,
			5:  20,
			7:  10.75,
			8:  8,
			10: 3.666,
		}

		for k, v := range data {
			w.AddWeight(k, v)
		}

		assert.Equal(t, len(data), len(w.indexes))
		assert.Equal(t, len(data), len(w.weights))

		temp := make([]Index, 3)
		for ix := 0; ix < 1000; ix++ {
			w.FillRandom(prng, len(temp), temp)
			assert.NotZero(t, temp[0])
			assert.NotZero(t, temp[1])
			assert.NotZero(t, temp[2])
			assert.NotEqual(t, temp[0], temp[1])
			assert.NotEqual(t, temp[0], temp[2])
			assert.NotEqual(t, temp[1], temp[2])
		}

	})
}

func TestWeighting_FillRandomDedup4(t *testing.T) {
	t.Run("fill random dedup4", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		w := AcquireWeightingDedup4().(*WeightingDedup4)
		require.NotNil(t, w)
		defer w.Release()

		data := map[Index]float64{
			1:  50.25,
			2:  33.333,
			4:  25.5,
			5:  20,
			7:  10.75,
			8:  8,
			10: 3.666,
		}

		for k, v := range data {
			w.AddWeight(k, v)
		}

		assert.Equal(t, len(data), len(w.indexes))
		assert.Equal(t, len(data), len(w.weights))

		temp := make([]Index, 4)
		for ix := 0; ix < 1000; ix++ {
			w.FillRandom(prng, len(temp), temp)
			assert.NotZero(t, temp[0])
			assert.NotZero(t, temp[1])
			assert.NotZero(t, temp[2])
			assert.NotZero(t, temp[3])
			assert.NotEqual(t, temp[0], temp[1])
			assert.NotEqual(t, temp[0], temp[2])
			assert.NotEqual(t, temp[0], temp[3])
			assert.NotEqual(t, temp[1], temp[2])
			assert.NotEqual(t, temp[1], temp[3])
			assert.NotEqual(t, temp[2], temp[3])
		}

	})
}

func TestWeighting_FillRandomDedup5(t *testing.T) {
	t.Run("fill random dedup5", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		w := AcquireWeightingDedup5().(*WeightingDedup5)
		require.NotNil(t, w)
		defer w.Release()

		data := map[Index]float64{
			1:  50.25,
			2:  33.333,
			4:  25.5,
			5:  20,
			7:  10.75,
			8:  8,
			10: 3.666,
		}

		for k, v := range data {
			w.AddWeight(k, v)
		}

		assert.Equal(t, len(data), len(w.indexes))
		assert.Equal(t, len(data), len(w.weights))

		temp := make([]Index, 5)
		for ix := 0; ix < 1000; ix++ {
			w.FillRandom(prng, len(temp), temp)
			assert.NotZero(t, temp[0])
			assert.NotZero(t, temp[1])
			assert.NotZero(t, temp[2])
			assert.NotZero(t, temp[3])
			assert.NotZero(t, temp[4])
			assert.NotEqual(t, temp[0], temp[1])
			assert.NotEqual(t, temp[0], temp[2])
			assert.NotEqual(t, temp[0], temp[3])
			assert.NotEqual(t, temp[0], temp[4])
			assert.NotEqual(t, temp[1], temp[2])
			assert.NotEqual(t, temp[1], temp[3])
			assert.NotEqual(t, temp[1], temp[4])
			assert.NotEqual(t, temp[2], temp[3])
			assert.NotEqual(t, temp[2], temp[4])
			assert.NotEqual(t, temp[3], temp[4])
		}

	})
}

func TestWeighting_FillRandomUnique3(t *testing.T) {
	t.Run("fill random unique3", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		w := AcquireWeightingUnique3().(*WeightingUnique3)
		require.NotNil(t, w)
		defer w.Release()

		w.AddWeights(Indexes{4, 5, 6}, []float64{65, 25, 10})

		assert.Equal(t, 3, len(w.indexes))
		assert.Equal(t, 3, len(w.weights))

		temp := make([]Index, 3)
		for ix := 0; ix < 1000; ix++ {
			w.FillRandom(prng, len(temp), temp)
			assert.NotZero(t, temp[0])
			assert.NotZero(t, temp[1])
			assert.NotZero(t, temp[2])
			assert.NotEqual(t, temp[0], temp[1])
			assert.NotEqual(t, temp[0], temp[2])
			assert.NotEqual(t, temp[1], temp[2])
		}

	})
}

func TestWeights1_1_1_0_0_0(t *testing.T) {
	t.Run("weighting 1,1,1,0,0,0", func(t *testing.T) {
		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		output := Indexes{0, 1, 2, 3, 4, 5}
		weights := []float64{25000, 2400, 680, 10, 1, 90000}

		w := AcquireWeighting().AddWeights(output, weights)
		require.NotNil(t, w)

		high := 100000000
		counts := make(map[Index]int, len(output))

		for ix := 0; ix < high; ix++ {
			q := w.RandomIndex(prng)
			counts[q] = counts[q] + 1
		}

		fmt.Printf("total count: %d\n", high)
		fmt.Printf("outputs    : %v\n", output)
		fmt.Printf("weights    : %v\n", weights)

		for ix := range output {
			fmt.Printf("output %.2d  : %d\n", output[ix], counts[output[ix]])
		}
	})
}

func TestWeights1_1_1_0_0_0_golang(t *testing.T) {
	t.Run("weighting 1,1,1,0,0,0 golang", func(t *testing.T) {
		prng := &gorng{r: rand.New(rand.NewSource(time.Now().UnixMicro()))}

		output := Indexes{0, 1, 2, 3, 4, 5}
		weights := []float64{90000, 25000, 2400, 680, 10, 1}

		w := AcquireWeighting().AddWeights(output, weights)
		require.NotNil(t, w)

		high := 100000000
		counts := make(map[Index]int, len(output))

		for ix := 0; ix < high; ix++ {
			q := w.RandomIndex(prng)
			counts[q] = counts[q] + 1
		}

		fmt.Printf("total count: %d\n", high)
		fmt.Printf("outputs    : %v\n", output)
		fmt.Printf("weights    : %v\n", weights)

		for ix := range output {
			fmt.Printf("output %.2d  : %d\n", output[ix], counts[output[ix]])
		}
	})
}

type gorng struct {
	r *rand.Rand
}

func (r *gorng) ReturnToPool()  {}
func (r *gorng) Uint32() uint32 { return r.r.Uint32() }
func (r *gorng) Uint64() uint64 { return r.r.Uint64() }
func (r *gorng) IntN(n int) int { return r.r.Intn(n) }

func (r *gorng) IntsN(n int, out []int) {
	for ix := range out {
		out[ix] = r.r.Intn(n)
	}
	return
}

func BenchmarkWeighting_RandomGet7(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	w := AcquireWeighting().(*WeightingNoDedup)
	defer w.Release()

	data := map[Index]float64{
		0:  75,
		1:  50.25,
		2:  33.333,
		4:  25.5,
		5:  20,
		7:  10.75,
		10: 3.666,
	}

	for k, v := range data {
		w.AddWeight(k, v)
	}

	for i := 0; i < b.N; i++ {
		w.nextPRNG(prng)
	}
}

func BenchmarkWeighting_RandomGet8(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	w := AcquireWeighting().(*WeightingNoDedup)
	defer w.Release()

	data := map[Index]float64{
		0:  50,
		1:  48,
		2:  46,
		4:  42,
		5:  40,
		7:  36,
		8:  34,
		10: 30,
	}

	for k, v := range data {
		w.AddWeight(k, v)
	}

	for i := 0; i < b.N; i++ {
		w.nextPRNG(prng)
	}
}

func BenchmarkWeighting_RandomGet9(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	w := AcquireWeighting().(*WeightingNoDedup)
	defer w.Release()

	data := map[Index]float64{
		0:  50,
		1:  48,
		2:  46,
		4:  42,
		5:  40,
		7:  36,
		8:  34,
		9:  32,
		10: 30,
	}

	for k, v := range data {
		w.AddWeight(k, v)
	}

	for i := 0; i < b.N; i++ {
		w.nextPRNG(prng)
	}
}

func BenchmarkWeighting_RandomGet10(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	w := AcquireWeighting().(*WeightingNoDedup)
	defer w.Release()

	data := map[Index]float64{
		0:  50,
		1:  48,
		2:  46,
		4:  42,
		5:  40,
		6:  38,
		7:  36,
		8:  34,
		9:  32,
		10: 30,
	}

	for k, v := range data {
		w.AddWeight(k, v)
	}

	for i := 0; i < b.N; i++ {
		w.nextPRNG(prng)
	}
}

func BenchmarkWeighting_RandomGet11(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	w := AcquireWeighting().(*WeightingNoDedup)
	defer w.Release()

	data := map[Index]float64{
		0:  50,
		1:  48,
		2:  46,
		4:  42,
		5:  40,
		6:  38,
		7:  36,
		8:  34,
		9:  32,
		10: 30,
		12: 26,
	}

	for k, v := range data {
		w.AddWeight(k, v)
	}

	for i := 0; i < b.N; i++ {
		w.nextPRNG(prng)
	}
}

func BenchmarkWeighting_RandomGet12(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	w := AcquireWeighting().(*WeightingNoDedup)
	defer w.Release()

	data := map[Index]float64{
		0:  50,
		1:  48,
		2:  46,
		4:  42,
		5:  40,
		6:  38,
		7:  36,
		8:  34,
		9:  32,
		10: 30,
		11: 28,
		12: 26,
	}

	for k, v := range data {
		w.AddWeight(k, v)
	}

	for i := 0; i < b.N; i++ {
		w.nextPRNG(prng)
	}
}

func BenchmarkWeighting_RandomGet14(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	w := AcquireWeighting().(*WeightingNoDedup)
	defer w.Release()

	data := map[Index]float64{
		0:  50,
		1:  48,
		2:  46,
		3:  48,
		4:  42,
		5:  40,
		6:  38,
		7:  36,
		8:  34,
		9:  32,
		10: 30,
		11: 28,
		12: 26,
		13: 24,
	}

	for k, v := range data {
		w.AddWeight(k, v)
	}

	for i := 0; i < b.N; i++ {
		w.nextPRNG(prng)
	}
}

func BenchmarkWeighting_RandomGet15(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	w := AcquireWeighting().(*WeightingNoDedup)
	defer w.Release()

	data := map[Index]float64{
		0:  50,
		1:  48,
		2:  46,
		3:  48,
		4:  42,
		5:  40,
		6:  38,
		7:  36,
		8:  34,
		9:  32,
		10: 30,
		11: 28,
		12: 26,
		13: 24,
		14: 20,
	}

	for k, v := range data {
		w.AddWeight(k, v)
	}

	for i := 0; i < b.N; i++ {
		w.nextPRNG(prng)
	}
}

func BenchmarkWeighting_RandomGet16(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	w := AcquireWeighting().(*WeightingNoDedup)
	defer w.Release()

	data := map[Index]float64{
		0:  50,
		1:  48,
		2:  46,
		3:  48,
		4:  42,
		5:  40,
		6:  38,
		7:  36,
		8:  34,
		9:  32,
		10: 30,
		11: 28,
		12: 26,
		13: 24,
		14: 20,
		15: 15,
	}

	for k, v := range data {
		w.AddWeight(k, v)
	}

	for i := 0; i < b.N; i++ {
		w.nextPRNG(prng)
	}
}

func BenchmarkWeighting_RandomGet17(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	w := AcquireWeighting().(*WeightingNoDedup)
	defer w.Release()

	data := map[Index]float64{
		0:  50,
		1:  48,
		2:  46,
		3:  48,
		4:  42,
		5:  40,
		6:  38,
		7:  36,
		8:  34,
		9:  32,
		10: 30,
		11: 28,
		12: 26,
		13: 24,
		14: 20,
		15: 15,
		16: 10,
	}

	for k, v := range data {
		w.AddWeight(k, v)
	}

	for i := 0; i < b.N; i++ {
		w.nextPRNG(prng)
	}
}

func BenchmarkWeighting_RandomIndex(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	w := AcquireWeighting().(*WeightingNoDedup)
	defer w.Release()

	data := map[Index]float64{
		0:  75,
		1:  50.25,
		2:  33.333,
		4:  25.5,
		5:  20,
		7:  10.75,
		8:  8,
		10: 3.666,
	}

	for k, v := range data {
		w.AddWeight(k, v)
	}

	for i := 0; i < b.N; i++ {
		w.RandomIndex(prng)
	}
}

func BenchmarkWeighting_RandomIndexManyItems(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	w := AcquireWeighting().(*WeightingNoDedup)
	defer w.Release()

	data := map[Index]float64{
		0:  75,
		1:  50.25,
		2:  33.333,
		4:  25.5,
		5:  20,
		7:  10.75,
		8:  8,
		10: 3.666,
		12: 33.333,
		14: 25.5,
		15: 20,
		17: 10.75,
		22: 33.333,
		24: 25.5,
		25: 20,
		27: 10.75,
	}

	for k, v := range data {
		w.AddWeight(k, v)
	}

	for i := 0; i < b.N; i++ {
		w.RandomIndex(prng)
	}
}

func BenchmarkWeighting_FillRandom3Repeat(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	w := AcquireWeighting().(*WeightingNoDedup)
	defer w.Release()

	data := map[Index]float64{
		0:  75,
		1:  50.25,
		2:  33.333,
		4:  25.5,
		5:  20,
		7:  10.75,
		8:  8,
		10: 3.666,
		11: 3,
		12: 2.4,
	}

	for k, v := range data {
		w.AddWeight(k, v)
	}

	max := 3
	temp := make([]Index, max)

	for i := 0; i < b.N; i++ {
		w.FillRandom(prng, max, temp)
	}
}

func BenchmarkWeighting_FillRandom3NoRepeat2(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	w := AcquireWeightingDedup3().(*WeightingDedup3)
	defer w.Release()

	data := map[Index]float64{
		0:  75,
		1:  50.25,
		2:  33.333,
		4:  25.5,
		5:  20,
		7:  10.75,
		8:  8,
		10: 3.666,
		11: 3,
		12: 2.4,
	}

	for k, v := range data {
		w.AddWeight(k, v)
	}

	max := 3
	temp := make([]Index, max)

	for i := 0; i < b.N; i++ {
		w.FillRandom(prng, max, temp)
	}
}

func BenchmarkWeighting_FillRandom5Repeat(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	w := AcquireWeighting().(*WeightingNoDedup)
	defer w.Release()

	data := map[Index]float64{
		0:  75,
		1:  50.25,
		2:  33.333,
		3:  30,
		4:  25.5,
		5:  20,
		6:  15,
		7:  10.75,
		8:  8,
		10: 3.666,
		11: 3,
		12: 2.4,
	}

	for k, v := range data {
		w.AddWeight(k, v)
	}

	max := 5
	temp := make([]Index, max)

	for i := 0; i < b.N; i++ {
		w.FillRandom(prng, max, temp)
	}
}

func BenchmarkWeighting_FillRandom5NoRepeat4(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	w := AcquireWeightingDedup5().(*WeightingDedup5)
	defer w.Release()

	data := map[Index]float64{
		0:  75,
		1:  50.25,
		2:  33.333,
		3:  30,
		4:  25.5,
		5:  20,
		6:  15,
		7:  10.75,
		8:  8,
		10: 3.666,
		11: 3,
		12: 2.4,
	}

	for k, v := range data {
		w.AddWeight(k, v)
	}

	max := 5
	temp := make([]Index, max)

	for i := 0; i < b.N; i++ {
		w.FillRandom(prng, max, temp)
	}
}
