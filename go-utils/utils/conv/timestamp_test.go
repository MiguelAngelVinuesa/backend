package conv

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimestampFromAny(t *testing.T) {
	n := time.Now().Round(time.Millisecond)

	testCases := []struct {
		name  string
		in    any
		dflts []time.Time
		want  time.Time
	}{
		{
			name: "nil",
			want: time.Time{},
		},
		{
			name: "unix milli",
			in:   n.UnixMilli(),
			want: n,
		},
		{
			name: "time",
			in:   n,
			want: n,
		},
		{
			name: "*time",
			in:   &n,
			want: n,
		},
		{
			name: "empty string",
			in:   "",
			want: time.Time{},
		},
		{
			name: "bad string",
			in:   "abcdef",
			want: time.Time{},
		},
		{
			name: "RFC3339 nano",
			in:   n.Format(time.RFC3339Nano),
			want: n,
		},
		{
			name: "RFC3339",
			in:   n.Round(time.Second).Format(time.RFC3339),
			want: n.Round(time.Second),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := TimestampFromAny(tc.in, tc.dflts...)
			assert.Equal(t, tc.want, got)
		})
	}
}
