package utils

import (
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexes_Contains(t *testing.T) {
	testCases := []struct {
		name    string
		indexes Indexes
		index   Index
		want    bool
	}{
		{"empty", Indexes{}, 1, false},
		{"single - no", Indexes{2}, 1, false},
		{"single - yes", Indexes{1}, 1, true},
		{"multi - no", Indexes{2, 3, 4, 5, 6, 7}, 1, false},
		{"multi - yes", Indexes{1, 2, 3, 4, 5, 6, 7}, 6, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.indexes.Contains(tc.index)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestIndexes_MarshalJSON(t *testing.T) {
	testCases := []struct {
		name    string
		indexes Indexes
		j       string
	}{
		{"empty", nil, "null"},
		{"single", Indexes{4}, "[4]"},
		{"three", Indexes{4, 5, 7}, "[4,5,7]"},
		{"five", Indexes{4, 6, 8, 10, 12}, "[4,6,8,10,12]"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			j, err := json.Marshal(tc.indexes)
			require.NoError(t, err)
			require.NotNil(t, j)
			assert.Equal(t, tc.j, string(j))
		})
	}
}

func TestSortIndexes(t *testing.T) {
	testCases := []struct {
		name    string
		indexes Indexes
		want    Indexes
	}{
		{"1", Indexes{1}, Indexes{1}},
		{"2 sorted", Indexes{1, 2}, Indexes{1, 2}},
		{"2 unsorted", Indexes{2, 1}, Indexes{1, 2}},
		{"3 sorted", Indexes{1, 2, 7}, Indexes{1, 2, 7}},
		{"3 unsorted", Indexes{2, 7, 1}, Indexes{1, 2, 7}},
		{"5 sorted", Indexes{1, 2, 7, 9, 11}, Indexes{1, 2, 7, 9, 11}},
		{"5 unsorted", Indexes{9, 11, 2, 7, 1}, Indexes{1, 2, 7, 9, 11}},
		{"11 sorted", Indexes{1, 2, 2, 3, 3, 5, 5, 5, 7, 8, 9}, Indexes{1, 2, 2, 3, 3, 5, 5, 5, 7, 8, 9}},
		{"11 unsorted", Indexes{2, 1, 7, 2, 4, 6, 4, 1, 8, 9, 4}, Indexes{1, 1, 2, 2, 4, 4, 4, 6, 7, 8, 9}},
		{"21 sorted", Indexes{1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 4, 4, 5, 6, 6, 7, 7, 7, 8, 9, 9}, Indexes{1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 4, 4, 5, 6, 6, 7, 7, 7, 8, 9, 9}},
		{"21 unsorted", Indexes{1, 2, 6, 2, 5, 3, 1, 4, 3, 7, 9, 2, 4, 7, 9, 3, 1, 6, 8, 7, 2}, Indexes{1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 4, 4, 5, 6, 6, 7, 7, 7, 8, 9, 9}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := SortIndexes(tc.indexes)
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestCopyIndexes(t *testing.T) {
	testCases := []struct {
		name   string
		input  Indexes
		output Indexes
		want   Indexes
	}{
		{
			name: "both nil",
		},
		{
			name:   "input nil",
			output: Indexes{0, 1, 2},
			want:   Indexes{},
		},
		{
			name:   "input empty",
			input:  Indexes{},
			output: Indexes{0, 1, 2},
			want:   Indexes{},
		},
		{
			name:  "output nil",
			input: Indexes{0, 1, 2},
			want:  Indexes{0, 1, 2},
		},
		{
			name:   "output empty",
			input:  Indexes{0, 1, 2},
			output: Indexes{},
			want:   Indexes{0, 1, 2},
		},
		{
			name:   "both filled",
			input:  Indexes{4, 5, 6},
			output: Indexes{7, 8, 9},
			want:   Indexes{4, 5, 6},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := CopyIndexes(tc.input, tc.output)
			if tc.want == nil {
				require.Nil(t, got)
			} else {
				require.NotNil(t, got)
				assert.EqualValues(t, tc.want, got)
			}
		})
	}
}
