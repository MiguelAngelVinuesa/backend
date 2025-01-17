package utils

import (
	"fmt"
	"strings"
)

// UInt8s is a convenience type for a slice of uint8.
type UInt8s []uint8

// PurgeUInt8s returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
func PurgeUInt8s(input UInt8s, capacity int) UInt8s {
	if cap(input) < capacity {
		return make(UInt8s, 0, capacity)
	}
	return input[:0]
}

// CopyUInt8s makes a deep copy of the input slice into the output slice and returns the output slice.
// If the input slice is nil, the output slice will be truncated if it wasn't nil itself.
// The output slice is reused if it has enough capacity, otherwise a new slice will be created on the heap.
func CopyUInt8s(input, output UInt8s) UInt8s {
	if l := len(input); l > 0 {
		output = PurgeUInt8s(output, l)[:l]
		copy(output, input)
		return output
	}
	if len(output) > 0 {
		return output[:0]
	}
	return output
}

// CopyPurgeUInt8s makes a deep copy of the input slice into the output slice and returns the output slice.
// The output slice is reused if it has enough capacity, otherwise a new slice will be created on the heap.
func CopyPurgeUInt8s(input, output UInt8s, capacity int) UInt8s {
	l := len(input)
	if l > capacity {
		capacity = l
	}
	output = PurgeUInt8s(output, capacity)[:l]
	copy(output, input)
	return output
}

// Contains returns true if the slice contains the given number.
func (u UInt8s) Contains(n uint8) bool {
	for _, i := range u {
		if i == n {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json marshaller interface.
func (u UInt8s) MarshalJSON() ([]byte, error) {
	if u == nil {
		return []byte("null"), nil
	}
	return []byte(strings.ReplaceAll(fmt.Sprintf("%v", u), " ", ",")), nil
}

// PurgeUInt16s returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
func PurgeUInt16s(input []uint16, capacity int) []uint16 {
	if cap(input) < capacity {
		return make([]uint16, 0, capacity)
	}
	return input[:0]
}

// PurgeUInt64s returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
func PurgeUInt64s(input []uint64, capacity int) []uint64 {
	if cap(input) < capacity {
		return make([]uint64, 0, capacity)
	}
	return input[:0]
}

// PurgeInt64s returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
func PurgeInt64s(input []int64, capacity int) []int64 {
	if cap(input) < capacity {
		return make([]int64, 0, capacity)
	}
	return input[:0]
}

// PurgeInts returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
func PurgeInts(input []int, capacity int) []int {
	if cap(input) < capacity {
		return make([]int, 0, capacity)
	}
	return input[:0]
}

// CopyInts makes a deep copy of the input slice into the output slice and returns the output slice.
// If the input slice is nil, the output slice will be truncated if it wasn't nil itself.
// The output slice is reused if it has enough capacity, otherwise a new slice will be created on the heap.
func CopyInts(input, output []int) []int {
	if l := len(input); l > 0 {
		output = PurgeInts(output, l)[:l]
		copy(output, input)
		return output
	}
	if len(output) > 0 {
		return output[:0]
	}
	return output
}

// PurgeUInt32s returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
func PurgeUInt32s(input []uint32, capacity int) []uint32 {
	if cap(input) < capacity {
		return make([]uint32, 0, capacity)
	}
	return input[:0]
}

// PurgeFloats returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
func PurgeFloats(input []float64, capacity int) []float64 {
	if cap(input) < capacity {
		return make([]float64, 0, capacity)
	}
	return input[:0]
}

// PurgeBools returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
func PurgeBools(input []bool, capacity int) []bool {
	if cap(input) < capacity {
		return make([]bool, 0, capacity)
	}
	return input[:0]
}
