package conv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolFromAny(t *testing.T) {
	testCases := []struct {
		name  string
		in    any
		want1 bool
		want2 bool
	}{
		{"empty", nil, false, true},
		{"false", false, false, false},
		{"true", true, true, true},
		{"0", 0, false, false},
		{"1", 1, true, true},
		{"0.0", 0.0, false, false},
		{"1.0", 1.0, true, true},
		{"empty string", "", false, true},
		{"string xyz", "xyz", false, true},
		{"string false", "false", false, false},
		{"string true", "true", true, true},
		{"string 0", "0", false, false},
		{"string 1", "1", true, true},
		{"string 0.0", "0.0", false, true},
		{"string 1.0", "1.0", false, true},
		{"string f", "f", false, false},
		{"string t", "t", true, true},
		{"string no", "no", false, false},
		{"string yes", "yes", true, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := BoolFromAny(tc.in)
			assert.Equal(t, tc.want1, got)

			got = BoolFromAny(tc.in, true)
			assert.Equal(t, tc.want2, got)
		})
	}
}
