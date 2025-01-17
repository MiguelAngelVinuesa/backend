package utils

import (
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
)

func TestStringFromAny(t *testing.T) {
	testCases := []struct {
		name string
		in   any
		want string
	}{
		{name: "nil", want: "<nil>"},
		{name: "empty string", in: "", want: ""},
		{name: "int 0", in: 0, want: "0"},
		{name: "int 1", in: 1, want: "1"},
		{name: "int 123", in: 123, want: "123"},
		{name: "int64 0", in: int64(0), want: "0"},
		{name: "int64 1", in: int64(1), want: "1"},
		{name: "int64 123", in: int64(123), want: "123"},
		{name: "float64 0.0", in: 0.0, want: "0"},
		{name: "float64 1.5", in: 1.5, want: "1.5"},
		{name: "float64 123.123", in: 123.123, want: "123.123"},
		{name: "json.Number 0", in: json.Number("0"), want: "0"},
		{name: "json.Number 1", in: json.Number("1"), want: "1"},
		{name: "json.Number 123", in: json.Number("123"), want: "123"},
		{name: "json.Number 123.123", in: json.Number("123.123"), want: "123.123"},
		{name: "true", in: true, want: "1"},
		{name: "false", in: false, want: "0"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StringFromAny(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
}
