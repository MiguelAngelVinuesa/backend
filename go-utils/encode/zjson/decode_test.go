package zjson

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDecoder(t *testing.T) {
	t.Run("new decoder", func(t *testing.T) {
		dec := AcquireDecoder([]byte("haha"))
		require.NotNil(t, dec)
		defer dec.Release()
	})
}

func TestDecoder_Object(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  *Data1
		fail  bool
	}{
		{
			name:  "empty - good",
			input: "{}",
			want:  &Data1{},
		},
		{
			name:  "empty wsp - good",
			input: "  {  }  ",
			want:  &Data1{},
		},
		{
			name:  "missing end",
			input: "{",
			want:  &Data1{},
			fail:  true,
		},
		{
			name:  "missing end wsp",
			input: "  {  ",
			want:  &Data1{},
			fail:  true,
		},
		{
			name:  "missing start",
			input: "}",
			want:  &Data1{},
			fail:  true,
		},
		{
			name:  "missing start wsp",
			input: "  }  ",
			want:  &Data1{},
			fail:  true,
		},
		{
			name:  "b1 - good",
			input: `  { "b1"  : true  } `,
			want:  &Data1{B1: true},
		},
		{
			name:  "b1 - fail (1)",
			input: `  { "b1"  :  } `,
			fail:  true,
		},
		{
			name:  "b1 - fail (2)",
			input: `  { "b1"  true  } `,
			fail:  true,
		},
		{
			name:  "b1 - fail (3)",
			input: `  { "b1 : true  } `,
			fail:  true,
		},
		{
			name:  "b1 - fail (4)",
			input: `  { b1 : true  } `,
			fail:  true,
		},
		{
			name:  "b2 - good",
			input: `  { "b2"  : false  } `,
			want:  &Data1{},
		},
		{
			name:  "i1 - good (1)",
			input: `  { "i1"  : 12  } `,
			want:  &Data1{I1: 12},
		},
		{
			name:  "i2 - good (2)",
			input: `  { "i1"  : 255  } `,
			want:  &Data1{I1: 255},
		},
		{
			name:  "i1 - fail (1)",
			input: `  { "i1"  : -12  } `,
			fail:  true,
		},
		{
			name:  "i1 - fail (2)",
			input: `  { "i1"  : 256  } `,
			fail:  true,
		},
		{
			name:  "i1 - fail (3)",
			input: `  { "i1"  : haha  } `,
			fail:  true,
		},
		{
			name:  "i2 - good",
			input: `  { "i2"  : 1234  } `,
			want:  &Data1{I2: 1234},
		},
		{
			name:  "i2 - fail (1)",
			input: `  { "i2"  : -12  } `,
			fail:  true,
		},
		{
			name:  "i2 - fail (2)",
			input: `  { "i2"  : 65536  } `,
			fail:  true,
		},
		{
			name:  "i3 - good (1)",
			input: `  { "i3"  : 12345678  } `,
			want:  &Data1{I3: 12345678},
		},
		{
			name:  "i3 - good (2)",
			input: `  { "i3"  : -12345678  } `,
			want:  &Data1{I3: -12345678},
		},
		{
			name:  "i3 - good (3)",
			input: `  { "i3"  : -123456789  } `,
			want:  &Data1{I3: -123456789},
		},
		{
			name:  "i4 - good",
			input: `  { "i4"  : 1234567812345678  } `,
			want:  &Data1{I4: 1234567812345678},
		},
		{
			name:  "i4 - fail",
			input: `  { "i4"  : -1  } `,
			fail:  true,
		},
		{
			name:  "f1 - good (1)",
			input: `  { "f1"  : 1  } `,
			want:  &Data1{F1: 1},
		},
		{
			name:  "f1 - good (2)",
			input: `  { "f1"  : 1.25  } `,
			want:  &Data1{F1: 1.25},
		},
		{
			name:  "f1 - good (3)",
			input: `  { "f1"  : -987.654321  } `,
			want:  &Data1{F1: -987.654321},
		},
		{
			name:  "f1 - good (4)",
			input: `  { "f1"  : 10e12  } `,
			want:  &Data1{F1: 10e12},
		},
		{
			name:  "f1 - fail",
			input: `  { "f1"  : hoho  } `,
			fail:  true,
		},
		{
			name:  "s1 - good (1)",
			input: `  { "s1"  : ""  } `,
			want:  &Data1{},
		},
		{
			name:  "s1 - good (2)",
			input: `  { "s1"  : "hihi"  } `,
			want:  &Data1{S1: "hihi"},
		},
		{
			name:  "s2 - good",
			input: `  { "s2"  : "This is\n some text\t with\r whitespace  \r\n"  } `,
			want:  &Data1{S2: "This is\n some text\t with\r whitespace  \r\n"},
		},
		{
			name:  "d - good (1)",
			input: `  { "d"  :  [  ]  } `,
			want:  &Data1{},
		},
		{
			name:  "d - good (2)",
			input: `  { "d"  :  [ { } ]  } `,
			want:  &Data1{D: []*Data2{{}}},
		},
		{
			name:  "d - good (3)",
			input: `  { "d"  :  [ { "i2" : 500, "b2"  : true } ]  } `,
			want:  &Data1{D: []*Data2{{B2: true, I2: 500}}},
		},
		{
			name:  "d - good (4)",
			input: `  { "d"  :  [ { "i2" : 500, "b2"  : true } , { "i1" :100, "b1": true } ]  } `,
			want:  &Data1{D: []*Data2{{B2: true, I2: 500}, {B1: true, I1: 100}}},
		},
		{
			name:  "d - fail (1)",
			input: `  { "d"  :  [  } `,
			fail:  true,
		},
		{
			name:  "d - fail (2)",
			input: `  { "d"  :  ]  } `,
			fail:  true,
		},
		{
			name:  "d - fail (3)",
			input: `  { "d"  :  [  { ]  } `,
			fail:  true,
		},
		{
			name:  "d - fail (4)",
			input: `  { "d"  :  [ } ] } `,
			fail:  true,
		},
		{
			name:  "d - fail (5)",
			input: `  { "d"  :  [ , ] } `,
			fail:  true,
		},
		{
			name:  "combo - good",
			input: ` { "s2":  "non"  ,"d"  :  [ { "i2" : 500, "b2"  : true } , { "i1" :100, "b1": true } ]  , "b2" :true, "i3" : -54321, "f1"   : -5.67, "i1": 9  }  `,
			want:  &Data1{B2: true, I1: 9, I3: -54321, F1: -5.67, S2: "non", D: []*Data2{{B2: true, I2: 500}, {B1: true, I1: 100}}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dec := AcquireDecoder([]byte(tc.input))
			require.NotNil(t, dec)
			defer dec.Release()

			got := &Data1{}
			ok := dec.Object(got)

			if tc.fail {
				assert.False(t, ok)
				assert.Error(t, dec.Error())
			} else {
				assert.True(t, ok)
				assert.NoError(t, dec.Error())
				assert.EqualValues(t, tc.want, got)
			}
		})
	}
}

func TestReadString(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  string
		fail  bool
	}{
		{
			name:  "empty - good",
			input: `""`,
		},
		{
			name:  "empty wsp - good",
			input: `  ""  `,
		},
		{
			name:  "empty - missing start",
			input: ``,
			fail:  true,
		},
		{
			name:  "empty wsp - missing start",
			input: `  `,
			fail:  true,
		},
		{
			name:  "empty - missing end",
			input: `"`,
			fail:  true,
		},
		{
			name:  "empty wsp - missing end",
			input: `  "  `,
			fail:  true,
		},
		{
			name:  "short - good",
			input: `"Short"`,
			want:  "Short",
		},
		{
			name:  "short wsp - good",
			input: `  "Short"  `,
			want:  "Short",
		},
		{
			name:  "blanks - good",
			input: `" Short "`,
			want:  " Short ",
		},
		{
			name:  "blanks wsp - good",
			input: `  " Short "  `,
			want:  " Short ",
		},
		{
			name:  "escaped - good",
			input: `"\tShort\n sentence:\t\"hello world\""`,
			want:  "\tShort\n sentence:\t\"hello world\"",
		},
		{
			name:  "escaped wsp - good",
			input: `  "\tShort \n sentence:\t\"hello world\""  `,
			want:  "\tShort \n sentence:\t\"hello world\"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dec := AcquireDecoder([]byte(tc.input))
			require.NotNil(t, dec)
			defer dec.Release()

			start, end, escaped, ok := dec.readString()
			if tc.fail {
				assert.False(t, ok)
				assert.Error(t, dec.Error())
			} else {
				assert.True(t, ok)
				assert.NoError(t, dec.Error())

				b := dec.buf[start:end]
				if escaped {
					assert.Equal(t, tc.want, string(dec.Unescaped(b)))
				} else {
					assert.Equal(t, tc.want, string(b))
				}
			}
		})
	}
}

func TestDecode(t *testing.T) {
	dec := AcquireDecoder(j1)
	defer dec.Release()

	d2 := &Data1{}
	dec.Object(d2)
	assert.EqualValues(t, d1, d2)
}

func TestDecoder_StringMap(t *testing.T) {
	testCases := []struct {
		name string
		j    string
		want map[string]string
		fail bool
	}{
		{
			name: "nil",
			fail: true,
		},
		{
			name: "fail",
			j:    `not valid`,
			fail: true,
		},
		{
			name: "empty",
			j:    `{}`,
			want: map[string]string{},
		},
		{
			name: "one",
			j:    `{"x":"y"}`,
			want: map[string]string{"x": "y"},
		},
		{
			name: "few",
			j:    `{"x":"y","a":"11","cc":"one of \nmany"}`,
			want: map[string]string{"x": "y", "a": "11", "cc": "one of \nmany"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dec := AcquireDecoder([]byte(tc.j))
			require.NotNil(t, dec)
			defer dec.Release()

			got, ok := dec.StringMap(nil)
			if tc.fail {
				require.False(t, ok)
				require.Error(t, dec.Error())
			} else {
				require.True(t, ok)
				require.NoError(t, dec.Error())
				assert.EqualValues(t, tc.want, got)
			}
		})
	}
}

func BenchmarkDecodeSetup(b *testing.B) {
	dec := AcquireDecoder(nil)
	defer dec.Release()

	for i := 0; i < b.N; i++ {
		dec.Reset(j1)
	}
}

func BenchmarkDecodeEmpty(b *testing.B) {
	dec := AcquireDecoder(nil)
	defer dec.Release()

	d2 := &Data1{}
	s := []byte(`{}`)

	for i := 0; i < b.N; i++ {
		dec.Reset(s)
		dec.Object(d2)
	}
}

func BenchmarkDecodeBool(b *testing.B) {
	dec := AcquireDecoder(nil)
	defer dec.Release()

	d2 := &Data1{}
	s := []byte(`{"b1":true}`)

	for i := 0; i < b.N; i++ {
		dec.Reset(s)
		dec.Object(d2)
	}
}

func BenchmarkDecodeUint8_small(b *testing.B) {
	dec := AcquireDecoder(nil)
	defer dec.Release()

	d2 := &Data1{}
	s := []byte(`{"i1":1}`)

	for i := 0; i < b.N; i++ {
		dec.Reset(s)
		dec.Object(d2)
	}
}

func BenchmarkDecodeUint8_big(b *testing.B) {
	dec := AcquireDecoder(nil)
	defer dec.Release()

	d2 := &Data1{}
	s := []byte(`{"i1":123}`)

	for i := 0; i < b.N; i++ {
		dec.Reset(s)
		dec.Object(d2)
	}
}

func BenchmarkDecodeUint16_small(b *testing.B) {
	dec := AcquireDecoder(nil)
	defer dec.Release()

	d2 := &Data1{}
	s := []byte(`{"i2":1}`)

	for i := 0; i < b.N; i++ {
		dec.Reset(s)
		dec.Object(d2)
	}
}

func BenchmarkDecodeUint16_big(b *testing.B) {
	dec := AcquireDecoder(nil)
	defer dec.Release()

	d2 := &Data1{}
	s := []byte(`{"i2":12345}`)

	for i := 0; i < b.N; i++ {
		dec.Reset(s)
		dec.Object(d2)
	}
}

func BenchmarkDecodeString_short(b *testing.B) {
	dec := AcquireDecoder(nil)
	defer dec.Release()

	d2 := &Data1{}
	s := []byte(`{"s1":"haha"}`)

	for i := 0; i < b.N; i++ {
		dec.Reset(s)
		dec.Object(d2)
	}
}

func BenchmarkDecodeString_long(b *testing.B) {
	dec := AcquireDecoder(nil)
	defer dec.Release()

	d2 := &Data1{}
	s := []byte(`{"s1":"haha haha haha haha haha haha haha haha haha haha haha haha haha haha"}`)

	for i := 0; i < b.N; i++ {
		dec.Reset(s)
		dec.Object(d2)
	}
}

func BenchmarkDecode(b *testing.B) {
	dec := AcquireDecoder(nil)
	defer dec.Release()

	for i := 0; i < b.N; i++ {
		d2 := NewData1()
		dec.Reset(j1)
		dec.Object(d2)
		d2.ReturnToPool()
	}
}

// func BenchmarkTrueEqual(b *testing.B) {
// 	var b1 = []byte{'t', 'r', 'u', 'e'}
// 	var b2 = []byte{'t', 'r', 'u', 'e'}
//
// 	for i := 0; i < b.N; i++ {
// 		bytes.Equal(b1, b2)
// 	}
// }
//
// func BenchmarkTrueNotEqual(b *testing.B) {
// 	var b1 = []byte{'t', 'r', 'u', 'e'}
// 	var b2 = []byte{'f', 'a', 'l', 's', 'e'}
//
// 	for i := 0; i < b.N; i++ {
// 		bytes.Equal(b1, b2)
// 	}
// }
//
// func BenchmarkFalseEqual(b *testing.B) {
// 	var b1 = []byte{'f', 'a', 'l', 's', 'e'}
// 	var b2 = []byte{'f', 'a', 'l', 's', 'e'}
//
// 	for i := 0; i < b.N; i++ {
// 		bytes.Equal(b1, b2)
// 	}
// }
//
// func BenchmarkFalseNotEqual(b *testing.B) {
// 	var b1 = []byte{'f', 'a', 'l', 's', 'e'}
// 	var b2 = []byte{'t', 'r', 'u', 'e'}
//
// 	for i := 0; i < b.N; i++ {
// 		bytes.Equal(b1, b2)
// 	}
// }
//
// func BenchmarkTrueCompare(b *testing.B) {
// 	var b1 = []byte{'t', 'r', 'u', 'e'}
// 	var b2 = []byte{'t', 'r', 'u', 'e'}
//
// 	for i := 0; i < b.N; i++ {
// 		bytes.Compare(b1, b2)
// 	}
// }
//
// func BenchmarkTrueNotCompare(b *testing.B) {
// 	var b1 = []byte{'t', 'r', 'u', 'e'}
// 	var b2 = []byte{'f', 'a', 'l', 's', 'e'}
//
// 	for i := 0; i < b.N; i++ {
// 		bytes.Compare(b1, b2)
// 	}
// }
//
// func BenchmarkFalseCompare(b *testing.B) {
// 	var b1 = []byte{'f', 'a', 'l', 's', 'e'}
// 	var b2 = []byte{'f', 'a', 'l', 's', 'e'}
//
// 	for i := 0; i < b.N; i++ {
// 		bytes.Compare(b1, b2)
// 	}
// }
//
// func BenchmarkFalseNotCompare(b *testing.B) {
// 	var b1 = []byte{'f', 'a', 'l', 's', 'e'}
// 	var b2 = []byte{'t', 'r', 'u', 'e'}
//
// 	for i := 0; i < b.N; i++ {
// 		bytes.Compare(b1, b2)
// 	}
// }
