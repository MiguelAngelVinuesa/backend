package zjson

import (
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEncoder(t *testing.T) {
	t.Run("new encoder", func(t *testing.T) {
		e := AcquireEncoder(512)
		require.NotNil(t, e)
		defer e.Release()

		e.Object(d1)
		j2 := e.Bytes()
		if !reflect.DeepEqual(j2, j1) {
			log.Fatalf("encoding failed:\n%s\n%s", string(j2), string(j1))
		}
	})
}

func TestEncoder_IntBoolField(t *testing.T) {
	t.Run("int bool field", func(t *testing.T) {
		e := AcquireEncoder(32)
		require.NotNil(t, e)
		defer e.Release()

		e.IntBoolField("bool", true)
		assert.Equal(t, `"bool":1`, string(e.Bytes()))

		e.Reset()
		e.IntBoolField("boolVal", false)
		assert.Equal(t, `"boolVal":0`, string(e.Bytes()))
	})
}

func TestEncoder_TimestampField(t *testing.T) {
	t.Run("timestamp field", func(t *testing.T) {
		e := AcquireEncoder(32)
		require.NotNil(t, e)
		defer e.Release()

		e.TimestampField("epoch", time.Unix(0, 0).UTC())
		assert.Equal(t, `"epoch":"19700101000000Z"`, string(e.Bytes()))

		e.Reset()
		e.TimestampField("1999", time.Date(1999, 12, 31, 23, 59, 59, 999000000, time.UTC))
		assert.Equal(t, `"1999":"19991231235959.999Z"`, string(e.Bytes()))

		e.Reset()
		e.TimestampField("2000", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
		assert.Equal(t, `"2000":"20000101000000Z"`, string(e.Bytes()))
	})
}

func TestEncoder_StringMapOpt(t *testing.T) {
	testCases := []struct {
		name string
		f    string
		m    map[string]string
		j    string
		j2   string
	}{
		{
			name: "nil",
			f:    "a",
			j:    ``,
		},
		{
			name: "empty",
			f:    "b",
			m:    map[string]string{},
			j:    ``,
		},
		{
			name: "one",
			f:    "c",
			m:    map[string]string{"x": "y"},
			j:    `"c":{"x":"y"}`,
		},
		{
			name: "two",
			f:    "c",
			m:    map[string]string{"x": "y", "a": "11"},
			j:    `"c":{"x":"y","a":"11"}`,
			j2:   `"c":{"a":"11","x":"y"}`,
		},
		{
			name: "two escape",
			f:    "c",
			m:    map[string]string{"x": "y", "a": "1 of \n15"},
			j:    `"c":{"x":"y","a":"1 of \n15"}`,
			j2:   `"c":{"a":"1 of \n15","x":"y"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := AcquireEncoder(512)
			require.NotNil(t, e)
			defer e.Release()

			e.StringMapFieldOpt(tc.f, tc.m)
			j := string(e.Bytes())
			if j != tc.j && j != tc.j2 {
				assert.Equal(t, tc.j, j)
				assert.Equal(t, tc.j2, j)
			}
		})
	}
}

func BenchmarkEncode(b *testing.B) {
	e := AcquireEncoder(512)
	defer e.Release()

	for i := 0; i < b.N; i++ {
		e.Reset()
		e.Object(d1)
	}
}
