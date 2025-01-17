package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMultiplier(t *testing.T) {
	testCases := []struct {
		name string
		in   []float64
		want float64
	}{
		{name: "none", want: 1.0},
		{name: "single invalid (1)", in: []float64{0.0}, want: 1.0},
		{name: "single invalid (2)", in: []float64{1.0}, want: 1.0},
		{name: "few invalid", in: []float64{0.0, 0.0, 1.0}, want: 1.0},
		{name: "single good (1)", in: []float64{1.50}, want: 1.5},
		{name: "single good (2)", in: []float64{5.0}, want: 5.0},
		{name: "few good", in: []float64{2.0, 3.0, 1.5}, want: 9.0},
		{name: "many good", in: []float64{0.5, 6.0, 3.0, 1.5, 2.0, 3.25, 2}, want: 175.5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := NewMultiplier(tc.in...)
			assert.Equal(t, tc.want, got)
		})
	}
}
