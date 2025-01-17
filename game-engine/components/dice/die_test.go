package dice

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"
)

func TestNewStandardDie(t *testing.T) {
	t.Run("new standard die", func(t *testing.T) {
		d := NewStandardDie()
		require.NotNil(t, d)
		defer d.Release()

		assert.Equal(t, 6, d.faces)
		assert.Equal(t, 1, d.first)
		assert.Zero(t, len(d.values))

		r := rng.NewRNG()
		defer r.ReturnToPool()

		counts := make(map[int]int)
		max := 100000
		for ix := 0; ix < max; ix++ {
			i := d.Roll(r)
			counts[i] = counts[i] + 1
		}

		avg := max / 6
		low := avg * 95 / 100
		high := avg * 105 / 100

		for i, count := range counts {
			assert.GreaterOrEqual(t, i, 1, i)
			assert.LessOrEqual(t, i, 6, i)
			assert.GreaterOrEqual(t, count, low, i)
			assert.LessOrEqual(t, count, high, i)
		}
	})
}

func TestNewFacedDie(t *testing.T) {
	testCases := []struct {
		name  string
		faces int
		first int
	}{
		{"3, 0-based", 3, 0},
		{"4, 0-based", 4, 0},
		{"6, 0-based", 6, 0},
		{"7, 0-based", 7, 0},
		{"8, 0-based", 8, 0},
		{"12, 0-based", 12, 0},
		{"4, 1-based", 4, 1},
		{"4, 10-based", 4, 10},
		{"8, 1-based", 8, 1},
		{"8, 10-based", 8, 10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := NewFacedDie(tc.faces, tc.first)
			require.NotNil(t, d)
			defer d.Release()

			assert.Equal(t, tc.faces, d.faces)
			assert.Equal(t, tc.first, d.first)
			assert.Zero(t, len(d.values))

			r := rng.NewRNG()
			defer r.ReturnToPool()

			counts := make(map[int]int)
			max := 100000
			for ix := 0; ix < max; ix++ {
				i := d.Roll(r)
				counts[i] = counts[i] + 1
			}

			avg := max / tc.faces
			low := avg * 95 / 100
			high := avg * 105 / 100

			for i, count := range counts {
				assert.GreaterOrEqual(t, i, tc.first, i)
				assert.Less(t, i, tc.first+tc.faces, i)
				assert.GreaterOrEqual(t, count, low, i)
				assert.LessOrEqual(t, count, high, i)
			}
		})
	}
}

func TestNewValuesDie(t *testing.T) {
	testCases := []struct {
		name   string
		values []int
	}{
		{"a", []int{1, 2, 3, 4, 5, 6}},
		{"b", []int{1, 3, 5, 7, 9, 11, 13, 15}},
		{"c", []int{1, 2, 4, 8, 16, 32}},
		{"d", []int{-1, 0, 1}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := NewValuesDie(tc.values...)
			require.NotNil(t, d)
			defer d.Release()

			assert.Equal(t, len(tc.values), d.faces)
			assert.EqualValues(t, tc.values, d.values)

			r := rng.NewRNG()
			defer r.ReturnToPool()

			counts := make(map[int]int)
			max := 100000
			for ix := 0; ix < max; ix++ {
				i := d.Roll(r)
				counts[i] = counts[i] + 1
			}

			avg := max / len(tc.values)
			low := avg * 95 / 100
			high := avg * 105 / 100

			for i, count := range counts {
				assert.GreaterOrEqual(t, count, low, i)
				assert.LessOrEqual(t, count, high, i)
			}
		})
	}
}
