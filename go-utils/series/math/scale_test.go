package math

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAutoScaleInt64(t *testing.T) {
	testCases := []struct {
		name     string
		min      int64
		max      int64
		steps    int64
		wantMin  int64
		wantMax  int64
		wantStep int64
	}{
		{"1-1000 / 15", 1, 1000, 15, 0, 1100, 100},
		{"999-1107 / 15", 999, 1107, 15, 990, 1110, 10},
		{"1501-2102 / 15", 1501, 2102, 15, 1500, 2150, 50},
		{"3170-5310 / 15", 3170, 5310, 15, 3000, 5400, 200},
		{"7099-10490 / 15", 7099, 10490, 15, 7000, 10500, 500},
		{"99-990490 / 15", 99, 990490, 15, 0, 1000000, 100000},
		{"799-442490 / 15", 799, 442490, 15, 0, 450000, 50000},
		{"7899-570490 / 15", 7899, 570490, 15, 0, 600000, 50000},
		{"107899-570490 / 15", 107899, 570490, 15, 100000, 600000, 50000},
		{"487899-570490 / 15", 487899, 570490, 15, 480000, 580000, 10000},
		{"567899-570490 / 15", 567899, 570490, 15, 567800, 570600, 200},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			min, max, step := AutoScaleInt64(tc.min, tc.max, tc.steps)
			assert.Equal(t, tc.wantMin, min)
			assert.Equal(t, tc.wantMax, max)
			assert.Equal(t, tc.wantStep, step)
		})
	}
}

func TestAutoScaleValuesInt64(t *testing.T) {
	testCases := []struct {
		name  string
		min   int64
		max   int64
		steps int64
		want  []int64
	}{
		{"1-1000 / 15", 1, 1000, 15, []int64{0, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000, 1100}},
		{"999-1107 / 12", 999, 1107, 12, []int64{990, 1000, 1010, 1020, 1030, 1040, 1050, 1060, 1070, 1080, 1090, 1100, 1110}},
		{"99-990490 / 12", 99, 990490, 12, []int64{0, 100000, 200000, 300000, 400000, 500000, 600000, 700000, 800000, 900000, 1000000}},
		{"107899-570490 / 10", 107899, 570490, 10, []int64{100000, 150000, 200000, 250000, 300000, 350000, 400000, 450000, 500000, 550000, 600000}},
		{"567899-570490 / 10", 567899, 570490, 10, []int64{567500, 568000, 568500, 569000, 569500, 570000, 570500}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			steps := AutoScaleValuesInt64(tc.min, tc.max, tc.steps)
			assert.EqualValues(t, tc.want, steps)
		})
	}
}

func TestAutoScaleTime(t *testing.T) {
	testCases := []struct {
		name     string
		min      time.Time
		max      time.Time
		steps    int64
		wantMin  time.Time
		wantMax  time.Time
		wantStep time.Duration
	}{
		{
			name:     "10:00-12:00 / 24",
			min:      time.Date(2022, 1, 1, 9, 2, 0, 0, time.UTC),
			max:      time.Date(2022, 1, 1, 17, 58, 15, 0, time.UTC),
			steps:    24,
			wantMin:  time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC),
			wantMax:  time.Date(2022, 1, 1, 18, 0, 0, 0, time.UTC),
			wantStep: 30 * time.Minute,
		},
		{
			name:     "10:00-12:00 / 24",
			min:      time.Date(2022, 1, 1, 10, 2, 0, 0, time.UTC),
			max:      time.Date(2022, 1, 1, 15, 58, 0, 0, time.UTC),
			steps:    24,
			wantMin:  time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
			wantMax:  time.Date(2022, 1, 1, 16, 0, 0, 0, time.UTC),
			wantStep: 15 * time.Minute,
		},
		{
			name:     "10:00-12:00 / 24",
			min:      time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
			max:      time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			steps:    24,
			wantMin:  time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
			wantMax:  time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			wantStep: 5 * time.Minute,
		},
		{
			name:     "10:00-11:00 / 24",
			min:      time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
			max:      time.Date(2022, 1, 1, 11, 0, 0, 0, time.UTC),
			steps:    24,
			wantMin:  time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
			wantMax:  time.Date(2022, 1, 1, 11, 0, 0, 0, time.UTC),
			wantStep: 2 * time.Minute,
		},
		{
			name:     "10:00-11:00 / 24",
			min:      time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
			max:      time.Date(2022, 1, 1, 10, 25, 0, 0, time.UTC),
			steps:    24,
			wantMin:  time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
			wantMax:  time.Date(2022, 1, 1, 10, 26, 0, 0, time.UTC),
			wantStep: 2 * time.Minute,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			min, max, step := AutoScaleUTC(tc.min, tc.max, tc.steps)
			assert.Equal(t, tc.wantMin, min)
			assert.Equal(t, tc.wantMax, max)
			assert.Equal(t, tc.wantStep, step)
		})
	}
}
