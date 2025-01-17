package conv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringFromAny(t *testing.T) {
	testCases := []struct {
		name string
		in   any
		dflt []string
		want string
	}{
		{name: "nil, no default", want: ""},
		{name: "nil, default", dflt: []string{"haha"}, want: "haha"},
		{name: "nil, defaults", dflt: []string{"haha", "hoho", "hihi"}, want: "haha"},
		{name: "empty, no default", in: "", want: ""},
		{name: "empty, default", in: "", dflt: []string{"haha"}, want: ""},
		{name: "string, no default", in: "yoyo", want: "yoyo"},
		{name: "string, default", in: "bobo", dflt: []string{"haha", "hoho", "hihi"}, want: "bobo"},
		{name: "int, no default", in: -123, want: "-123"},
		{name: "int8, no default", in: int8(124), want: "124"},
		{name: "int16, no default", in: int16(-1234), want: "-1234"},
		{name: "int32, no default", in: int32(12345), want: "12345"},
		{name: "int64, no default", in: int64(-123456), want: "-123456"},
		{name: "uint, no default", in: uint(123), want: "123"},
		{name: "uint8, no default", in: uint8(124), want: "124"},
		{name: "uint16, no default", in: uint16(1234), want: "1234"},
		{name: "uint32, no default", in: uint32(12345), want: "12345"},
		{name: "uint64, no default", in: uint64(123456), want: "123456"},
		{name: "float64, no default", in: 123.124, want: "123.124"},
		{name: "true, no default", in: true, want: "1"},
		{name: "false, no default", in: false, want: "0"},
		{name: "struct, no default", in: struct{ x int }{x: 123}, want: ""},
		{name: "struct, default", in: struct{ x int }{x: 123}, dflt: []string{"nono"}, want: "nono"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StringFromAny(tc.in, tc.dflt...)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStringsFromAny(t *testing.T) {
	testCases := []struct {
		name string
		in   any
		dflt []string
		want []string
	}{
		{name: "nil, no default", want: nil},
		{name: "nil, default", dflt: []string{"haha"}, want: []string{"haha"}},
		{name: "nil, defaults", dflt: []string{"haha", "hoho", "hihi"}, want: []string{"haha", "hoho", "hihi"}},
		{name: "empty, no default", in: "", want: []string{""}},
		{name: "empty, default", in: "", dflt: []string{"haha"}, want: []string{""}},
		{name: "string, no default", in: "yoyo", want: []string{"yoyo"}},
		{name: "string, default", in: "bobo", dflt: []string{"haha", "hoho", "hihi"}, want: []string{"bobo"}},
		{name: "strings, no default", in: []string{"yoyo", "lala", "rara"}, want: []string{"yoyo", "lala", "rara"}},
		{name: "ints, no default", in: []int{123, -124, 1256}, want: []string{"123", "-124", "1256"}},
		{name: "struct, no default", in: struct{ x int }{x: 123}, want: nil},
		{name: "struct, default", in: struct{ x int }{x: 123}, dflt: []string{"nono", "dodo"}, want: []string{"nono", "dodo"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StringsFromAny(tc.in, tc.dflt...)
			assert.Equal(t, tc.want, got)
		})
	}
}
