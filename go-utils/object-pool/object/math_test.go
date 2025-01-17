package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeSize(t *testing.T) {
	testCases := []struct {
		name    string
		size    int
		minSIze int
		want    int
	}{
		{"zeroes", 0, 0, 4},
		{"0 / 2", 0, 2, 4},
		{"1 / 2", 1, 2, 4},
		{"2 / 2", 2, 2, 4},
		{"3 / 2", 3, 2, 4},
		{"4 / 2", 4, 2, 4},
		{"5 / 2", 5, 2, 8},
		{"8 / 2", 8, 2, 8},
		{"9 / 2", 9, 2, 12},
		{"0 / 8", 0, 8, 8},
		{"1 / 8", 1, 8, 8},
		{"2 / 8", 2, 8, 8},
		{"7 / 8", 7, 8, 8},
		{"8 / 8", 8, 8, 8},
		{"9 / 8", 9, 8, 16},
		{"0 / 16", 0, 16, 16},
		{"7 / 16", 7, 16, 16},
		{"8 / 16", 8, 16, 16},
		{"15 / 16", 15, 16, 16},
		{"16 / 16", 16, 16, 16},
		{"17 / 16", 17, 16, 32},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := NormalizeSize(tc.size, tc.minSIze)
			assert.Equal(t, tc.want, got)
		})
	}
}
