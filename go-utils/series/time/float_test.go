package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFloats_MinMaxTime(t *testing.T) {
	now := time.Now().Round(time.Second)

	testCases := []struct {
		name    string
		values  FVs[float64]
		wantMin time.Time
		wantMax time.Time
	}{
		{
			name:   "empty",
			values: FVs[float64]{},
		},
		{
			name:    "single",
			values:  FVs[float64]{NewFV[float64](now, 1)},
			wantMin: now,
			wantMax: now,
		},
		{
			name:    "few",
			values:  FVs[float64]{NewFV[float64](now, 1), NewFV[float64](now.Add(time.Hour), 2), NewFV[float64](now.Add(-time.Hour), 3)},
			wantMin: now.Add(-time.Hour),
			wantMax: now.Add(time.Hour),
		},
		{
			name: "many",
			values: FVs[float64]{
				NewFV[float64](now, 1),
				NewFV[float64](now.Add(time.Second), 2),
				NewFV[float64](now.Add(-time.Second), 3),
				NewFV[float64](now.Add(time.Minute), 4),
				NewFV[float64](now.Add(-time.Minute), 5),
				NewFV[float64](now.Add(time.Hour), 6),
				NewFV[float64](now.Add(-time.Hour), 7),
			},
			wantMin: now.Add(-time.Hour),
			wantMax: now.Add(time.Hour),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			min, max := tc.values.MinMaxTime()
			assert.EqualValues(t, tc.wantMin, min)
			assert.EqualValues(t, tc.wantMax, max)
		})
	}
}

func TestFloats_MinMaxValue(t *testing.T) {
	now := time.Now().Round(time.Second)

	testCases := []struct {
		name    string
		values  FVs[float64]
		wantMin float64
		wantMax float64
	}{
		{
			name:   "empty",
			values: FVs[float64]{},
		},
		{
			name:    "single",
			values:  FVs[float64]{NewFV[float64](now, 1.1)},
			wantMin: 1.1,
			wantMax: 1.1,
		},
		{
			name:    "few",
			values:  FVs[float64]{NewFV[float64](now, 1.1), NewFV[float64](now.Add(time.Hour), 2.2), NewFV[float64](now.Add(-time.Hour), 3.3)},
			wantMin: 1.1,
			wantMax: 3.3,
		},
		{
			name: "many",
			values: FVs[float64]{
				NewFV[float64](now, 5.5),
				NewFV[float64](now.Add(time.Second), 4.4),
				NewFV[float64](now.Add(-time.Second), 3.3),
				NewFV[float64](now.Add(time.Minute), 1.1),
				NewFV[float64](now.Add(-time.Minute), 7.7),
				NewFV[float64](now.Add(time.Hour), 6.6),
				NewFV[float64](now.Add(-time.Hour), 2.2),
			},
			wantMin: 1.1,
			wantMax: 7.7,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			min, max := tc.values.MinMaxValue()
			assert.EqualValues(t, tc.wantMin, min)
			assert.EqualValues(t, tc.wantMax, max)
		})
	}
}
