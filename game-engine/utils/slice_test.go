package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPurgeUInt8s(t *testing.T) {
	testCases := []struct {
		name    string
		in      UInt8s
		cap     int
		wantCap int
	}{
		{"nil", nil, 5, 5},
		{"empty", UInt8s{}, 5, 5},
		{"short", UInt8s{1, 2, 3}, 5, 5},
		{"exact", UInt8s{1, 2, 3, 4, 5}, 5, 5},
		{"long", UInt8s{1, 2, 3, 4, 5, 6, 7}, 5, 7},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := PurgeUInt8s(tc.in, tc.cap)
			require.NotNil(t, out)
			assert.Equal(t, tc.wantCap, cap(out))
			assert.Zero(t, len(out))
		})
	}
}

func TestUInt8s_MarshalJSON(t *testing.T) {
	testCases := []struct {
		name  string
		input UInt8s
		json  string
	}{
		{"nil", nil, "null"},
		{"empty", UInt8s{}, "[]"},
		{"single", UInt8s{5}, "[5]"},
		{"multi", UInt8s{3, 2, 9, 6, 0, 1}, "[3,2,9,6,0,1]"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			json, err := tc.input.MarshalJSON()
			require.NoError(t, err)
			require.NotNil(t, json)
			assert.Equal(t, tc.json, string(json))
		})
	}
}

func TestPurgeInts(t *testing.T) {
	testCases := []struct {
		name    string
		in      []int
		cap     int
		wantCap int
	}{
		{"nil", nil, 5, 5},
		{"empty", []int{}, 5, 5},
		{"short", []int{1, 2, 3}, 5, 5},
		{"exact", []int{1, 2, 3, 4, 5}, 5, 5},
		{"long", []int{1, 2, 3, 4, 5, 6, 7}, 5, 7},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := PurgeInts(tc.in, tc.cap)
			require.NotNil(t, out)
			assert.Equal(t, tc.wantCap, cap(out))
			assert.Zero(t, len(out))
		})
	}
}

func TestPurgeFloats(t *testing.T) {
	testCases := []struct {
		name    string
		in      []float64
		cap     int
		wantCap int
	}{
		{"nil", nil, 5, 5},
		{"empty", []float64{}, 5, 5},
		{"short", []float64{1, 2, 3}, 5, 5},
		{"exact", []float64{1, 2, 3, 4, 5}, 5, 5},
		{"long", []float64{1, 2, 3, 4, 5, 6, 7}, 5, 7},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := PurgeFloats(tc.in, tc.cap)
			require.NotNil(t, out)
			assert.Equal(t, tc.wantCap, cap(out))
			assert.Zero(t, len(out))
		})
	}
}

func TestPurgeBools(t *testing.T) {
	testCases := []struct {
		name    string
		in      []bool
		cap     int
		wantCap int
	}{
		{"nil", nil, 5, 5},
		{"empty", []bool{}, 5, 5},
		{"short", []bool{false, true, false}, 5, 5},
		{"exact", []bool{true, true, false, true, false}, 5, 5},
		{"long", []bool{false, false, true, true, false, true, false}, 5, 7},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := PurgeBools(tc.in, tc.cap)
			require.NotNil(t, out)
			assert.Equal(t, tc.wantCap, cap(out))
			assert.Zero(t, len(out))
		})
	}
}
