package slots

import (
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMinMaxInt64(t *testing.T) {
	testCases := []struct {
		name      string
		increases []int64
		want      *MinMaxInt64
	}{
		{
			name: "empty",
			want: &MinMaxInt64{Counts: map[int64]uint64{}, first: true},
		},
		{
			name:      "0",
			increases: []int64{0},
			want:      &MinMaxInt64{Count: 1, Keys: []int64{0}, Counts: map[int64]uint64{0: 1}},
		},
		{
			name:      "0,0,0,0,0,0",
			increases: []int64{0, 0, 0, 0, 0, 0},
			want:      &MinMaxInt64{Count: 6, Keys: []int64{0}, Counts: map[int64]uint64{0: 6}},
		},
		{
			name:      "1",
			increases: []int64{1},
			want:      &MinMaxInt64{Count: 1, Total: 1, Min: 1, Max: 1, Keys: []int64{1}, Counts: map[int64]uint64{1: 1}},
		},
		{
			name:      "111",
			increases: []int64{111},
			want:      &MinMaxInt64{Count: 1, Total: 111, Min: 111, Max: 111, Keys: []int64{111}, Counts: map[int64]uint64{111: 1}},
		},
		{
			name:      "1,2,3,4,5",
			increases: []int64{1, 2, 3, 4, 5},
			want:      &MinMaxInt64{Count: 5, Total: 15, Min: 1, Max: 5, Keys: []int64{1, 2, 3, 4, 5}, Counts: map[int64]uint64{1: 1, 2: 1, 3: 1, 4: 1, 5: 1}},
		},
		{
			name:      "81,22,53,14,65",
			increases: []int64{81, 22, 53, 14, 65},
			want:      &MinMaxInt64{Count: 5, Total: 235, Min: 14, Max: 81, Keys: []int64{14, 22, 53, 65, 81}, Counts: map[int64]uint64{14: 1, 22: 1, 53: 1, 65: 1, 81: 1}},
		},
		{
			name:      "1,2,1,4,3,2,2,1,3",
			increases: []int64{1, 2, 1, 4, 3, 2, 2, 1, 3},
			want:      &MinMaxInt64{Count: 9, Total: 19, Min: 1, Max: 4, Keys: []int64{1, 2, 3, 4}, Counts: map[int64]uint64{1: 3, 2: 3, 3: 2, 4: 1}},
		},
		{
			name:      "1,2,0,0,1,4,3,0,2,2,1,0,3",
			increases: []int64{1, 2, 0, 0, 1, 4, 3, 0, 2, 2, 1, 0, 3},
			want:      &MinMaxInt64{Count: 13, Total: 19, Min: 1, Max: 4, Keys: []int64{0, 1, 2, 3, 4}, Counts: map[int64]uint64{0: 4, 1: 3, 2: 3, 3: 2, 4: 1}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := AcquireMinMaxInt64()
			require.NotNil(t, m)
			defer m.Release()

			for _, i := range tc.increases {
				m.Increase(i)
			}

			if !tc.want.Equals(m) {
				assert.EqualValues(t, tc.want, m)
			}

			m.ResetData()
			assert.Zero(t, m.Total)
			assert.Equal(t, int64(0), m.Min)
			assert.Equal(t, int64(0), m.Max)
			assert.Zero(t, len(m.Counts))
		})
	}
}

func TestMinMaxInt64_Merge(t *testing.T) {
	testCases := []struct {
		name  string
		in    *MinMaxInt64
		other *MinMaxInt64
		want  *MinMaxInt64
	}{
		{
			name:  "merge empties",
			in:    AcquireMinMaxInt64(),
			other: AcquireMinMaxInt64(),
			want:  AcquireMinMaxInt64(),
		},
		{
			name:  "merge with empty",
			in:    &MinMaxInt64{Total: 10, Min: 1, Max: 2, Keys: []int64{1, 2}, Counts: map[int64]uint64{1: 4, 2: 3}},
			other: AcquireMinMaxInt64(),
			want:  &MinMaxInt64{Total: 10, Min: 1, Max: 2, Keys: []int64{1, 2}, Counts: map[int64]uint64{1: 4, 2: 3}},
		},
		{
			name:  "merge to empty",
			in:    AcquireMinMaxInt64(),
			other: &MinMaxInt64{Total: 10, Min: 1, Max: 2, Keys: []int64{1, 2}, Counts: map[int64]uint64{1: 4, 2: 3}},
			want:  &MinMaxInt64{Total: 10, Min: 1, Max: 2, Keys: []int64{1, 2}, Counts: map[int64]uint64{1: 4, 2: 3}},
		},
		{
			name:  "merge non-empties",
			in:    &MinMaxInt64{Total: 10, Min: 1, Max: 2, Keys: []int64{1, 2}, Counts: map[int64]uint64{1: 4, 2: 3}},
			other: &MinMaxInt64{Total: 20, Min: 1, Max: 5, Keys: []int64{1, 2, 5}, Counts: map[int64]uint64{1: 6, 2: 2, 5: 2}},
			want:  &MinMaxInt64{Total: 30, Min: 1, Max: 5, Keys: []int64{1, 2, 5}, Counts: map[int64]uint64{1: 10, 2: 5, 5: 2}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.in
			got.Merge(tc.other)
			if !tc.want.Equals(got) {
				assert.EqualValues(t, tc.want, got)
			}
		})
	}
}

func TestNewMinMaxUInt64(t *testing.T) {
	testCases := []struct {
		name      string
		increases []uint64
		want      *MinMaxUInt64
	}{
		{
			name: "empty",
			want: &MinMaxUInt64{Counts: map[uint64]uint64{}, first: true},
		},
		{
			name:      "0",
			increases: []uint64{0},
			want:      &MinMaxUInt64{Count: 1, Keys: []uint64{0}, Counts: map[uint64]uint64{0: 1}},
		},
		{
			name:      "0,0,0,0,0,0",
			increases: []uint64{0, 0, 0, 0, 0, 0},
			want:      &MinMaxUInt64{Count: 6, Keys: []uint64{0}, Counts: map[uint64]uint64{0: 6}},
		},
		{
			name:      "1",
			increases: []uint64{1},
			want:      &MinMaxUInt64{Count: 1, Total: 1, Min: 1, Max: 1, Keys: []uint64{1}, Counts: map[uint64]uint64{1: 1}},
		},
		{
			name:      "111",
			increases: []uint64{111},
			want:      &MinMaxUInt64{Count: 1, Total: 111, Min: 111, Max: 111, Keys: []uint64{111}, Counts: map[uint64]uint64{111: 1}},
		},
		{
			name:      "1,2,3,4,5",
			increases: []uint64{1, 2, 3, 4, 5},
			want:      &MinMaxUInt64{Count: 5, Total: 15, Min: 1, Max: 5, Keys: []uint64{1, 2, 3, 4, 5}, Counts: map[uint64]uint64{1: 1, 2: 1, 3: 1, 4: 1, 5: 1}},
		},
		{
			name:      "81,22,53,14,65",
			increases: []uint64{81, 22, 53, 14, 65},
			want:      &MinMaxUInt64{Count: 5, Total: 235, Min: 14, Max: 81, Keys: []uint64{14, 22, 53, 65, 81}, Counts: map[uint64]uint64{14: 1, 22: 1, 53: 1, 65: 1, 81: 1}},
		},
		{
			name:      "1,2,1,4,3,2,2,1,3",
			increases: []uint64{1, 2, 1, 4, 3, 2, 2, 1, 3},
			want:      &MinMaxUInt64{Count: 9, Total: 19, Min: 1, Max: 4, Keys: []uint64{1, 2, 3, 4}, Counts: map[uint64]uint64{1: 3, 2: 3, 3: 2, 4: 1}},
		},
		{
			name:      "1,2,0,0,1,4,3,0,2,2,1,0,3",
			increases: []uint64{1, 2, 0, 0, 1, 4, 3, 0, 2, 2, 1, 0, 3},
			want:      &MinMaxUInt64{Count: 13, Total: 19, Min: 1, Max: 4, Keys: []uint64{0, 1, 2, 3, 4}, Counts: map[uint64]uint64{0: 4, 1: 3, 2: 3, 3: 2, 4: 1}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := AcquireMinMaxUInt64()
			require.NotNil(t, m)
			defer m.Release()

			for _, i := range tc.increases {
				m.Increase(i)
			}

			if !tc.want.Equals(m) {
				assert.EqualValues(t, tc.want, m)
			}

			m.ResetData()
			assert.Zero(t, m.Total)
			assert.Equal(t, uint64(0), m.Min)
			assert.Zero(t, m.Max)
			assert.Zero(t, len(m.Counts))
		})
	}
}

func TestMinMaxUInt64_Merge(t *testing.T) {
	testCases := []struct {
		name  string
		in    *MinMaxUInt64
		other *MinMaxUInt64
		want  *MinMaxUInt64
	}{
		{
			name:  "merge empties",
			in:    AcquireMinMaxUInt64(),
			other: AcquireMinMaxUInt64(),
			want:  AcquireMinMaxUInt64(),
		},
		{
			name:  "merge with empty",
			in:    &MinMaxUInt64{Count: 7, Total: 10, Min: 1, Max: 2, Keys: []uint64{1, 2}, Counts: map[uint64]uint64{1: 4, 2: 3}},
			other: AcquireMinMaxUInt64(),
			want:  &MinMaxUInt64{Count: 7, Total: 10, Min: 1, Max: 2, Keys: []uint64{1, 2}, Counts: map[uint64]uint64{1: 4, 2: 3}},
		},
		{
			name:  "merge to empty",
			in:    AcquireMinMaxUInt64(),
			other: &MinMaxUInt64{Count: 7, Total: 10, Min: 1, Max: 2, Keys: []uint64{1, 2}, Counts: map[uint64]uint64{1: 4, 2: 3}},
			want:  &MinMaxUInt64{Count: 7, Total: 10, Min: 1, Max: 2, Keys: []uint64{1, 2}, Counts: map[uint64]uint64{1: 4, 2: 3}},
		},
		{
			name:  "merge non-empties",
			in:    &MinMaxUInt64{Count: 7, Total: 10, Min: 1, Max: 2, Keys: []uint64{1, 2}, Counts: map[uint64]uint64{1: 4, 2: 3}},
			other: &MinMaxUInt64{Count: 10, Total: 20, Min: 1, Max: 5, Keys: []uint64{1, 2, 5}, Counts: map[uint64]uint64{1: 6, 2: 2, 5: 2}},
			want:  &MinMaxUInt64{Count: 17, Total: 30, Min: 1, Max: 5, Keys: []uint64{1, 2, 5}, Counts: map[uint64]uint64{1: 10, 2: 5, 5: 2}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.in
			got.Merge(tc.other)
			if !tc.want.Equals(got) {
				assert.EqualValues(t, tc.want, got)
			}
		})
	}
}

func TestNewMinMaxFloat64(t *testing.T) {
	testCases := []struct {
		name      string
		increases []float64
		want      *MinMaxFloat64
		j         string
	}{
		{
			name: "empty",
			want: &MinMaxFloat64{Counts: map[int64]uint64{}, first: true, Decimals: 1, Factor: 10},
			j:    `{"count":0,"total":0}`,
		},
		{
			name:      "0",
			increases: []float64{0},
			want:      &MinMaxFloat64{Count: 1, Keys: []int64{0}, Counts: map[int64]uint64{0: 1}, Decimals: 1, Factor: 10},
			j:         `{"count":1,"total":0,"counts":{"0.00":1}}`,
		},
		{
			name:      "0,0,0,0,0,0",
			increases: []float64{0, 0, 0, 0, 0, 0},
			want:      &MinMaxFloat64{Count: 6, Keys: []int64{0}, Counts: map[int64]uint64{0: 6}, Decimals: 1, Factor: 10},
			j:         `{"count":6,"total":0,"counts":{"0.00":6}}`,
		},
		{
			name:      "0.5",
			increases: []float64{0.5},
			want:      &MinMaxFloat64{Count: 1, Total: 0.5, Min: 0.5, Max: 0.5, Keys: []int64{5}, Counts: map[int64]uint64{5: 1}, Decimals: 1, Factor: 10},
			j:         `{"count":1,"total":0.5,"min":0.5,"max":0.5,"counts":{"0.50":1}}`,
		},
		{
			name:      "111",
			increases: []float64{111},
			want:      &MinMaxFloat64{Count: 1, Total: 111, Min: 111, Max: 111, Keys: []int64{1110}, Counts: map[int64]uint64{1110: 1}, Decimals: 1, Factor: 10},
			j:         `{"count":1,"total":111,"min":111,"max":111,"counts":{"111.00":1}}`,
		},
		{
			name:      "1,2,3,4,5",
			increases: []float64{1, 2, 3, 4, 5},
			want:      &MinMaxFloat64{Count: 5, Total: 15, Min: 1, Max: 5, Keys: []int64{10, 20, 30, 40, 50}, Counts: map[int64]uint64{10: 1, 20: 1, 30: 1, 40: 1, 50: 1}, Decimals: 1, Factor: 10},
			j:         `{"count":5,"total":15,"min":1,"max":5,"counts":{"1.00":1,"2.00":1,"3.00":1,"4.00":1,"5.00":1}}`,
		},
		{
			name:      "8.1,2.2,5.3,1.4,6.5",
			increases: []float64{8.1, 2.2, 5.3, 1.4, 6.5},
			want:      &MinMaxFloat64{Count: 5, Total: 23.5, Min: 1.4, Max: 8.1, Keys: []int64{14, 22, 53, 65, 81}, Counts: map[int64]uint64{14: 1, 22: 1, 53: 1, 65: 1, 81: 1}, Decimals: 1, Factor: 10},
			j:         `{"count":5,"total":23.5,"min":1.4,"max":8.1,"counts":{"1.40":1,"2.20":1,"5.30":1,"6.50":1,"8.10":1}}`,
		},
		{
			name:      "1,2.5,1,4,3,2.5,2.5,1,3",
			increases: []float64{1, 2.5, 1, 4, 3, 2.5, 2.5, 1, 3},
			want:      &MinMaxFloat64{Count: 9, Total: 20.5, Min: 1, Max: 4, Keys: []int64{10, 25, 30, 40}, Counts: map[int64]uint64{10: 3, 25: 3, 30: 2, 40: 1}, Decimals: 1, Factor: 10},
			j:         `{"count":9,"total":20.5,"min":1,"max":4,"counts":{"1.00":3,"2.50":3,"3.00":2,"4.00":1}}`,
		},
		{
			name:      "1,2,0,0,1,4,3,0,2,2,1,0,3",
			increases: []float64{1, 2, 0, 0, 1, 4, 3, 0, 2, 2, 1, 0, 3},
			want:      &MinMaxFloat64{Count: 13, Total: 19, Min: 1, Max: 4, Keys: []int64{0, 10, 20, 30, 40}, Counts: map[int64]uint64{0: 4, 10: 3, 20: 3, 30: 2, 40: 1}, Decimals: 1, Factor: 10},
			j:         `{"count":13,"total":19,"min":1,"max":4,"counts":{"0.00":4,"1.00":3,"2.00":3,"3.00":2,"4.00":1}}`,
		},
		{
			name:      "8.12,2.27,5.311,1.498,6.52788",
			increases: []float64{8.12, 2.27, 5.311, 1.498, 6.52788},
			want:      &MinMaxFloat64{Count: 5, Total: 23.726879999999998, Min: 1.498, Max: 8.12, Keys: []int64{15, 23, 53, 65, 81}, Counts: map[int64]uint64{15: 1, 23: 1, 53: 1, 65: 1, 81: 1}, Decimals: 1, Factor: 10},
			j:         `{"count":5,"total":23.726879999999998,"min":1.498,"max":8.12,"counts":{"1.50":1,"2.30":1,"5.30":1,"6.50":1,"8.10":1}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := AcquireMinMaxFloat64(1)
			require.NotNil(t, m)
			defer m.Release()

			assert.Equal(t, 1, m.Decimals)
			assert.Equal(t, 10.0, m.Factor)

			for _, i := range tc.increases {
				m.Increase(i)
			}

			if !tc.want.Equals(m) {
				assert.EqualValues(t, tc.want, m)
			}

			j, err := json.Marshal(m)
			require.NoError(t, err)
			require.NotNil(t, j)
			assert.Equal(t, tc.j, string(j))

			m.ResetData()
			assert.Zero(t, m.Total)
			assert.Zero(t, m.Min)
			assert.Zero(t, m.Max)
			assert.Zero(t, len(m.Counts))
			assert.True(t, m.first)
			assert.Equal(t, 1, m.Decimals)
			assert.Equal(t, 10.0, m.Factor)
		})
	}
}

func TestMinMaxFloat64_Merge(t *testing.T) {
	testCases := []struct {
		name  string
		in    *MinMaxFloat64
		other *MinMaxFloat64
		want  *MinMaxFloat64
	}{
		{
			name:  "merge empties",
			in:    AcquireMinMaxFloat64(1),
			other: AcquireMinMaxFloat64(1),
			want:  AcquireMinMaxFloat64(1),
		},
		{
			name:  "merge with empty",
			in:    &MinMaxFloat64{Total: 10, Min: 1, Max: 2, Keys: []int64{10, 20}, Counts: map[int64]uint64{10: 4, 20: 3}, Decimals: 1, Factor: 10.0},
			other: AcquireMinMaxFloat64(1),
			want:  &MinMaxFloat64{Total: 10, Min: 1, Max: 2, Keys: []int64{10, 20}, Counts: map[int64]uint64{10: 4, 20: 3}, Decimals: 1, Factor: 10.0},
		},
		{
			name:  "merge to empty",
			in:    AcquireMinMaxFloat64(1),
			other: &MinMaxFloat64{Total: 10, Min: 1, Max: 2, Keys: []int64{10, 20}, Counts: map[int64]uint64{10: 4, 20: 3}, Decimals: 1, Factor: 10.0},
			want:  &MinMaxFloat64{Total: 10, Min: 1, Max: 2, Keys: []int64{10, 20}, Counts: map[int64]uint64{10: 4, 20: 3}, Decimals: 1, Factor: 10.0},
		},
		{
			name:  "merge non-empties",
			in:    &MinMaxFloat64{Total: 10, Min: 1, Max: 2, Keys: []int64{10, 20}, Counts: map[int64]uint64{10: 4, 20: 3}, Decimals: 1, Factor: 10.0},
			other: &MinMaxFloat64{Total: 20, Min: 1, Max: 5, Keys: []int64{10, 20, 50}, Counts: map[int64]uint64{10: 6, 20: 2, 50: 2}, Decimals: 1, Factor: 10.0},
			want:  &MinMaxFloat64{Total: 30, Min: 1, Max: 5, Keys: []int64{10, 20, 50}, Counts: map[int64]uint64{10: 10, 20: 5, 50: 2}, Decimals: 1, Factor: 10.0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.in
			got.Merge(tc.other)
			if !tc.want.Equals(got) {
				assert.EqualValues(t, tc.want, got)
			}
		})
	}
}

func BenchmarkMinMaxFloat64_Merge(b *testing.B) {
	m1 := AcquireMinMaxFloat64(2)
	defer m1.Release()

	m1.Increase(12.00)
	m1.Increase(12.00)
	m1.Increase(4.00)
	m1.Increase(4.00)
	m1.Increase(4.00)
	m1.Increase(7.50)
	m1.Increase(12.75)
	m1.Increase(7.50)

	m2 := AcquireMinMaxFloat64(2)
	defer m2.Release()

	for i := 0; i < b.N; i++ {
		m2.Merge(m1)
	}
}
