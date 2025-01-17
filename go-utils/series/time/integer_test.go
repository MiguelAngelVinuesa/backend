package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValues_MinMaxTime(t *testing.T) {
	now := time.Now().Round(time.Second)

	testCases := []struct {
		name    string
		values  IVs[int]
		wantMin time.Time
		wantMax time.Time
	}{
		{
			name:   "empty",
			values: IVs[int]{},
		},
		{
			name:    "single",
			values:  IVs[int]{NewIV(now, 1)},
			wantMin: now,
			wantMax: now,
		},
		{
			name:    "few",
			values:  IVs[int]{NewIV(now, 1), NewIV(now.Add(time.Hour), 2), NewIV(now.Add(-time.Hour), 3)},
			wantMin: now.Add(-time.Hour),
			wantMax: now.Add(time.Hour),
		},
		{
			name: "many",
			values: IVs[int]{
				NewIV(now, 1),
				NewIV(now.Add(time.Second), 2),
				NewIV(now.Add(-time.Second), 3),
				NewIV(now.Add(time.Minute), 4),
				NewIV(now.Add(-time.Minute), 5),
				NewIV(now.Add(time.Hour), 6),
				NewIV(now.Add(-time.Hour), 7),
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

func TestValues_MinMaxInt(t *testing.T) {
	now := time.Now().Round(time.Second)

	testCases := []struct {
		name    string
		values  IVs[int]
		wantMin int
		wantMax int
	}{
		{
			name:   "empty",
			values: IVs[int]{},
		},
		{
			name:    "single",
			values:  IVs[int]{NewIV(now, 1)},
			wantMin: 1,
			wantMax: 1,
		},
		{
			name:    "few",
			values:  IVs[int]{NewIV(now, 1), NewIV(now.Add(time.Hour), 2), NewIV(now.Add(-time.Hour), 3)},
			wantMin: 1,
			wantMax: 3,
		},
		{
			name: "many",
			values: IVs[int]{
				NewIV(now, 5),
				NewIV(now.Add(time.Second), 4),
				NewIV(now.Add(-time.Second), 3),
				NewIV(now.Add(time.Minute), 1),
				NewIV(now.Add(-time.Minute), 7),
				NewIV(now.Add(time.Hour), 6),
				NewIV(now.Add(-time.Hour), 2),
			},
			wantMin: 1,
			wantMax: 7,
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

func TestValues_MinMaxUint8(t *testing.T) {
	now := time.Now().Round(time.Second)

	testCases := []struct {
		name    string
		values  IVs[uint8]
		wantMin uint8
		wantMax uint8
	}{
		{
			name:   "empty",
			values: IVs[uint8]{},
		},
		{
			name:    "single",
			values:  IVs[uint8]{NewIV[uint8](now, 1)},
			wantMin: 1,
			wantMax: 1,
		},
		{
			name:    "few",
			values:  IVs[uint8]{NewIV[uint8](now, 1), NewIV[uint8](now.Add(time.Hour), 2), NewIV[uint8](now.Add(-time.Hour), 3)},
			wantMin: 1,
			wantMax: 3,
		},
		{
			name: "many",
			values: IVs[uint8]{
				NewIV[uint8](now, 5),
				NewIV[uint8](now.Add(time.Second), 4),
				NewIV[uint8](now.Add(-time.Second), 3),
				NewIV[uint8](now.Add(time.Minute), 1),
				NewIV[uint8](now.Add(-time.Minute), 7),
				NewIV[uint8](now.Add(time.Hour), 6),
				NewIV[uint8](now.Add(-time.Hour), 2),
			},
			wantMin: 1,
			wantMax: 7,
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
