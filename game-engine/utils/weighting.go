package utils

import (
	"bytes"
	"fmt"
	"strconv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
)

// WeightedGenerator is the interface for generating weighted random numbers.
type WeightedGenerator interface {
	AddWeight(index Index, weight float64) WeightedGenerator
	AddWeights(indexes []Index, weights []float64) WeightedGenerator
	Options() Indexes
	RandomIndex(prng interfaces.Generator) Index
	FillRandom(prng interfaces.Generator, count int, out Indexes)
	fmt.Stringer
	pool.Objecter
}

// AcquireWeighting instantiates a weighting from the memory pool.
func AcquireWeighting() WeightedGenerator {
	return weightingPool.Acquire().(*WeightingNoDedup)
}

// AcquireWeightingDedup3 instantiates a weighting from the memory pool which deduplicates across 3 rows.
func AcquireWeightingDedup3() WeightedGenerator {
	return weightingDedup3Pool.Acquire().(*WeightingDedup3)
}

// AcquireWeightingDedup4 instantiates a weighting from the memory pool which deduplicates across 4 rows.
func AcquireWeightingDedup4() WeightedGenerator {
	return weightingDedup4Pool.Acquire().(*WeightingDedup4)
}

// AcquireWeightingDedup5 instantiates a weighting from the memory pool which deduplicates across 5 rows.
func AcquireWeightingDedup5() WeightedGenerator {
	return weightingDedup5Pool.Acquire().(*WeightingDedup5)
}

// AcquireWeightingUnique3 instantiates a weighting from the memory pool which will always produce 3 unique results.
func AcquireWeightingUnique3() WeightedGenerator {
	return weightingUnique3Pool.Acquire().(*WeightingUnique3)
}

// AddWeight adds an item index and the corresponding weight.
// The given weight will be rounded to 3 decimals.
// If the weight is zero or less the index is not added.
func (w *WeightingNoDedup) AddWeight(index Index, weight float64) WeightedGenerator {
	if weight <= 0 {
		return w
	}
	total := int(weight * 1000)
	if l := len(w.weights); l > 0 {
		total += w.weights[l-1]
	}
	w.indexes = append(w.indexes, index)
	w.weights = append(w.weights, total)
	w.total = total
	w.size = len(w.weights)
	return w
}

// AddWeights adds item indexes and the corresponding weights.
// The given weights will be rounded to 3 decimals.
// If a weight is zero or less the index is not added.
func (w *WeightingNoDedup) AddWeights(indexes []Index, weights []float64) WeightedGenerator {
	for ix := range indexes {
		w.AddWeight(indexes[ix], weights[ix])
	}
	return w
}

// Options returns the possible item indexes.
func (w *WeightingNoDedup) Options() Indexes {
	return w.indexes
}

// RandomIndex calculates a random item index using the given PRNG if needed.
// The function panics if the Weighting hasn't been initialized properly.
// Even if the no-repeat feature is set, it cannot be used here,
// so the caller is always responsible for de-duplicating repetitive outputs.
func (w *WeightingNoDedup) RandomIndex(prng interfaces.Generator) Index {
	return w.nextPRNG(prng)
}

// FillRandom fills the given slice with count random item indexes and returns the slice.
// The function panics if count >= len(in) or if the Weighting hasn't been initialized properly.
// This is useful for filling reels. It honours the no-repeat feature if set.
func (w *WeightingNoDedup) FillRandom(prng interfaces.Generator, count int, out Indexes) {
	for ix := 0; ix < count; ix++ {
		out[ix] = w.nextPRNG(prng)
	}
}

// nextPRNG retrieves the next random index using the given PRNG.
func (w *WeightingNoDedup) nextPRNG(prng interfaces.Generator) Index {
	switch w.size {
	case 0:
		return 0
	case 1:
		return w.indexes[0]
	default:
		return w.search(prng.IntN(w.total))
	}
}

// search uses the given weight to find the corresponding index.
func (w *WeightingNoDedup) search(m int) Index {
	var ix int

	last := len(w.weights) - 1
	if last >= 16 {
		// binary search if 16 or more symbols.
		first := 0
		for first <= last {
			ix = (first + last) / 2
			if m < w.weights[ix] {
				if ix == 0 || m >= w.weights[ix-1] {
					return w.indexes[ix]
				}
				last = ix - 1
			} else {
				first = ix + 1
			}
		}
		panic("Weighting: binary search failed")
	}

	// sequential search if less than 16 symbols.
	for m >= w.weights[ix] {
		ix++
	}
	return w.indexes[ix]
}

// WeightingNoDedup is used to randomly select one or more items by index based on their individual weights.
type WeightingNoDedup struct {
	total   int
	size    int
	indexes Indexes
	weights []int
	pool.Object
}

var weightingPool = pool.NewProducer(func() (pool.Objecter, func()) {
	w := &WeightingNoDedup{
		indexes: make(Indexes, 0, 16),
		weights: make([]int, 0, 16),
	}
	return w, w.reset
})

// reset clears the weighting.
func (w *WeightingNoDedup) reset() {
	if w != nil {
		w.ResetWeights()
	}
}

// ResetWeights resets the weighting to initial state.
func (w *WeightingNoDedup) ResetWeights() {
	w.total = 0
	w.size = 0
	w.indexes = w.indexes[:0]
	w.weights = w.weights[:0]
}

// String implements the Stringer interface.
func (w *WeightingNoDedup) String() string {
	b := bytes.Buffer{}
	b.WriteByte('[')
	for ix := range w.indexes {
		if ix > 1 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.Itoa(int(w.indexes[ix])))
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(w.weights[ix]))
	}
	b.WriteByte(']')
	return b.String()
}

// WeightingDedup3 is used to randomly select one or more items by index based on their individual weights.
// WeightingDedup3 deduplicates across 3 rows during FillRandom.
type WeightingDedup3 struct {
	WeightingNoDedup
}

var weightingDedup3Pool = pool.NewProducer(func() (pool.Objecter, func()) {
	w := &WeightingDedup3{}
	w.indexes = make(Indexes, 0, 16)
	w.weights = make([]int, 0, 16)
	return w, w.reset
})

// AddWeight is overridden to make sure we return the underlying struct.
func (w *WeightingDedup3) AddWeight(index Index, weight float64) WeightedGenerator {
	w.WeightingNoDedup.AddWeight(index, weight)
	return w
}

// AddWeights is overridden to make sure we return the underlying struct.
func (w *WeightingDedup3) AddWeights(indexes []Index, weights []float64) WeightedGenerator {
	w.WeightingNoDedup.AddWeights(indexes, weights)
	return w
}

// FillRandom fills the given slice with count random item indexes and returns the slice.
// The function panics if count >= len(in) or if the Weighting hasn't been initialized properly.
// This is useful for filling reels. It deduplicates across 3 rows.
func (w *WeightingDedup3) FillRandom(prng interfaces.Generator, count int, out Indexes) {
	for ix := 0; ix < count; ix++ {
		var n Index

		repeat := true
		for repeat {
			n = w.nextPRNG(prng)
			switch ix {
			case 0:
				repeat = false
			case 1:
				repeat = n == out[ix-1]
			default:
				repeat = n == out[ix-1] || n == out[ix-2]
			}
		}

		out[ix] = n
	}
}

// WeightingDedup4 is used to randomly select one or more items by index based on their individual weights.
// WeightingDedup4 deduplicates across 4 rows during FillRandom.
type WeightingDedup4 struct {
	WeightingNoDedup
}

var weightingDedup4Pool = pool.NewProducer(func() (pool.Objecter, func()) {
	w := &WeightingDedup4{}
	w.indexes = make(Indexes, 0, 16)
	w.weights = make([]int, 0, 16)
	return w, w.reset
})

// AddWeight is overridden to make sure we return the underlying struct.
func (w *WeightingDedup4) AddWeight(index Index, weight float64) WeightedGenerator {
	w.WeightingNoDedup.AddWeight(index, weight)
	return w
}

// AddWeights is overridden to make sure we return the underlying struct.
func (w *WeightingDedup4) AddWeights(indexes []Index, weights []float64) WeightedGenerator {
	w.WeightingNoDedup.AddWeights(indexes, weights)
	return w
}

// FillRandom fills the given slice with count random item indexes and returns the slice.
// The function panics if count >= len(in) or if the Weighting hasn't been initialized properly.
// This is useful for filling reels. It deduplicates across 4 rows.
func (w *WeightingDedup4) FillRandom(prng interfaces.Generator, count int, out Indexes) {
	for ix := 0; ix < count; ix++ {
		var n Index

		repeat := true
		for repeat {
			n = w.nextPRNG(prng)
			switch ix {
			case 0:
				repeat = false
			case 1:
				repeat = n == out[ix-1]
			case 2:
				repeat = n == out[ix-1] || n == out[ix-2]
			default:
				repeat = n == out[ix-1] || n == out[ix-2] || n == out[ix-3]
			}
		}

		out[ix] = n
	}
}

// WeightingDedup5 is used to randomly select one or more items by index based on their individual weights.
// WeightingDedup5 deduplicates across 5 rows during FillRandom.
type WeightingDedup5 struct {
	WeightingNoDedup
}

var weightingDedup5Pool = pool.NewProducer(func() (pool.Objecter, func()) {
	w := &WeightingDedup5{}
	w.indexes = make(Indexes, 0, 16)
	w.weights = make([]int, 0, 16)
	return w, w.reset
})

// AddWeight is overridden to make sure we return the underlying struct.
func (w *WeightingDedup5) AddWeight(index Index, weight float64) WeightedGenerator {
	w.WeightingNoDedup.AddWeight(index, weight)
	return w
}

// AddWeights is overridden to make sure we return the underlying struct.
func (w *WeightingDedup5) AddWeights(indexes []Index, weights []float64) WeightedGenerator {
	w.WeightingNoDedup.AddWeights(indexes, weights)
	return w
}

// FillRandom fills the given slice with count random item indexes and returns the slice.
// The function panics if count >= len(in) or if the Weighting hasn't been initialized properly.
// This is useful for filling reels. It deduplicates across 5 rows.
func (w *WeightingDedup5) FillRandom(prng interfaces.Generator, count int, out Indexes) {
	for ix := 0; ix < count; ix++ {
		var n Index

		repeat := true
		for repeat {
			n = w.nextPRNG(prng)
			switch ix {
			case 0:
				repeat = false
			case 1:
				repeat = n == out[ix-1]
			case 2:
				repeat = n == out[ix-1] || n == out[ix-2]
			case 3:
				repeat = n == out[ix-1] || n == out[ix-2] || n == out[ix-3]
			default:
				repeat = n == out[ix-1] || n == out[ix-2] || n == out[ix-3] || n == out[ix-4]
			}
		}

		out[ix] = n
	}
}

// WeightingUnique3 is used to randomly select one or more items by index based on their individual weights.
// WeightingUnique3 always produces exactly 3 unique results.
type WeightingUnique3 struct {
	WeightingNoDedup
}

var weightingUnique3Pool = pool.NewProducer(func() (pool.Objecter, func()) {
	w := &WeightingUnique3{}
	w.indexes = make(Indexes, 0, 16)
	w.weights = make([]int, 0, 16)
	return w, w.reset
})

// AddWeight is overridden to make sure we return the underlying struct.
func (w *WeightingUnique3) AddWeight(index Index, weight float64) WeightedGenerator {
	w.WeightingNoDedup.AddWeight(index, weight)
	if len(w.weights) > 3 {
		panic("WeightingUnique3: too many weights")
	}
	return w
}

// AddWeights is overridden to make sure we return the underlying struct.
func (w *WeightingUnique3) AddWeights(indexes []Index, weights []float64) WeightedGenerator {
	w.WeightingNoDedup.AddWeights(indexes, weights)
	if len(w.weights) > 3 {
		panic("WeightingUnique3: too many weights")
	}
	return w
}

// FillRandom fills the given slice with 3 random item indexes and returns the slice.
// The function panics if count != 3 or if the Weighting hasn't been initialized properly.
func (w *WeightingUnique3) FillRandom(prng interfaces.Generator, count int, out Indexes) {
	if count != 3 || len(w.weights) != 3 {
		panic("WeightingUnique3: invalid count")
	}

	m := w.nextPRNG(prng)
	out[0] = m

	if prng.IntN(10000) < 5000 {
		iy := 1
		for ix := range w.indexes {
			n := w.indexes[ix]
			if n != m {
				out[iy] = n
				iy++
			}
		}
	} else {
		iy := 2
		for ix := range w.indexes {
			n := w.indexes[ix]
			if n != m {
				out[iy] = n
				iy--
			}
		}
	}
}
