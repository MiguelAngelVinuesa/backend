package sharedlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewShuffler(t *testing.T) {
	testCases := []struct {
		name  string
		input []int
	}{
		{"5 element sorted", []int{1, 2, 3, 4, 5}},
		{"7 element sorted", []int{1, 2, 3, 4, 5, 6, 7}},
		{"13 element sorted", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}},
		{"5 element unsorted", []int{3, 5, 2, 1, 4}},
		{"7 element unsorted", []int{5, 7, 3, 6, 2, 4, 1}},
		{"11 element unsorted", []int{6, 11, 10, 5, 6, 9, 2, 3, 1, 4, 7}},
		{"13 element unsorted", []int{12, 7, 3, 9, 5, 2, 4, 13, 6, 11, 8, 10, 1}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := AcquireShuffler()
			require.NotNil(t, s)
			defer s.Release()

			got := make([]int, len(tc.input))
			copy(got, tc.input)
			s.Shuffle(got)
			assert.NotEqual(t, tc.input, got)
		})
	}
}
