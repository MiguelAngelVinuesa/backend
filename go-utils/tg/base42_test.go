package tg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToBase42(t *testing.T) {
	testCases := []struct {
		name string
		in   []byte
		want string
	}{
		{name: "empty", in: []byte{}, want: ""},
		{name: "1", in: []byte{1}, want: "4"},
		{name: "41", in: []byte{41}, want: "y"},
		{name: "42", in: []byte{42}, want: "34"},
		{name: "43", in: []byte{43}, want: "44"},
		{name: "42+41", in: []byte{83}, want: "y4"},
		{name: "42+42", in: []byte{84}, want: "36"},
		{name: "42+43", in: []byte{85}, want: "46"},
		{name: "42*42-1", in: []byte{227, 6}, want: "yy3"},
		{name: "42*42", in: []byte{228, 6}, want: "334"},
		{name: "42*42+1", in: []byte{229, 6}, want: "434"},
		{name: "medium", in: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, want: "DG3PY3ep3tA4D"},
		{name: "long", in: []byte("abcdefghijklmnopqrst"), want: "cFJrXJ9pJKAKaNKnfK74LHHLXaLjrL"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := ToBase42(tc.in)
			assert.Equal(t, tc.want, got)

			got2, err := FromBase42(got)
			require.NoError(t, err)
			assert.EqualValues(t, tc.in, got2)
		})
	}
}

func TestFromBase42_Error(t *testing.T) {
	t.Run("FromBase42 error", func(t *testing.T) {
		got, err := FromBase42("10")
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func BenchmarkToBase42(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToBase42([]byte("abcdefghijklmnopqrst"))
	}
}

func BenchmarkFromBase42(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FromBase42("cFJrXJ9pJKAKaNKnfK74LHHLXaLjrL")
	}
}
