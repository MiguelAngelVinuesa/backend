package utils

import (
	"math"
)

// Index is the type for an index of a weighted object.
type Index uint16

const (
	// MaxIndex can be used to represent an invalid index.
	MaxIndex = Index(math.MaxUint16)
	// NullIndex can be used to represent the zero (null) index.
	NullIndex = Index(0)
)

// Indexes is a slice of indexes.
type Indexes []Index

// Contains returns true if the slice contains the given index.
func (i Indexes) Contains(index Index) bool {
	for ix := range i {
		if i[ix] == index {
			return true
		}
	}
	return false
}

// PurgeIndexes returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
func PurgeIndexes(input Indexes, capacity int) Indexes {
	if cap(input) < capacity {
		return make(Indexes, 0, capacity)
	}
	return input[:0]
}

// CopyIndexes creates a deep copy of the input slice into the output slice and returns the output slice.
// If the input slice is nil, the output slice will be truncated if it wasn't nil itself.
// The output slice is reused if it has enough capacity, otherwise a new slice will be created on the heap.
func CopyIndexes(input, output Indexes) Indexes {
	if l := len(input); l > 0 {
		output = PurgeIndexes(output, l)[:l]
		copy(output, input)
		return output
	}
	if len(output) > 0 {
		return output[:0]
	}
	return output
}

// SortIndexes sorts the slice of indexes in ascending order.
// We use our own quick-sort as the standard sort package uses reflection which results in heap allocations.
func SortIndexes(input Indexes) Indexes {
	max := len(input)
	if max < 2 {
		return input
	}

	left, right, split := 0, max-1, max/2
	input[right], input[split] = input[split], input[right]

	for i := range input {
		if input[i] < input[right] {
			input[left], input[i] = input[i], input[left]
			left++
		}
	}

	input[left], input[right] = input[right], input[left]

	SortIndexes(input[:left])
	SortIndexes(input[left+1:])

	return input
}

var (
	ClearIndexes = make(Indexes, 100)
)
