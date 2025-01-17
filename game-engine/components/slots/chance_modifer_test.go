package slots

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPowerFunc(t *testing.T) {
	testCases := []struct {
		a    float64
		b    float64
		c    float64
		x    float64
		want float64
	}{
		{a: 0.2, b: 1.1, c: 0, x: 0, want: 0.2},
		{a: 0.2, b: 1.1, c: 0, x: 1, want: 0.22},
		{a: 0.2, b: 1.1, c: 0, x: 2, want: 0.242},
		{a: 0.2, b: 1.1, c: 0, x: 3, want: 0.2662},
		{a: 0.2, b: 1.1, c: 0, x: 4, want: 0.29282},
		{a: 0.2, b: 1.1, c: 0.5, x: 0, want: 0.7},
		{a: 0.2, b: 1.1, c: 0.5, x: 1, want: 0.72},
		{a: 0.2, b: 1.1, c: 0.5, x: 2, want: 0.742},
		{a: 0.2, b: 1.1, c: 0.5, x: 3, want: 0.7662},
		{a: 0.2, b: 1.1, c: 0.5, x: 4, want: 0.79282},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("%.2f * %.2f ^ %.2f + %.2f", tc.a, tc.b, tc.x, tc.c)
		t.Run(name, func(t *testing.T) {
			f := NewPowerFunc(tc.b, tc.c, func(*Spin) float64 { return tc.x })
			require.NotNil(t, f)

			got := f.Exec(tc.a, nil)
			assert.Less(t, math.Abs(tc.want-got), 0.0000001)
		})
	}
}

func TestInvalidFunc(t *testing.T) {
	t.Run("invalid func", func(t *testing.T) {
		f := &chanceModifier{}
		got := f.Exec(0.2, nil)
		assert.Equal(t, 0.2, got)
	})
}
