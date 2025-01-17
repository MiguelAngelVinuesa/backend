package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
)

func TestResetFloats(t *testing.T) {
	testCases := []struct {
		name   string
		floats Floats
		min    int
		max    int
		clear  bool
	}{
		{"nil", nil, 4, 8, true},
		{"empty", Floats{}, 4, 8, true},
		{"single - no clear", Floats{1}, 4, 8, false},
		{"single - clear", Floats{1}, 4, 8, true},
		{"few - no clear", Floats{2, 3, 1}, 4, 8, false},
		{"few - clear", Floats{3, 1, 2}, 4, 8, true},
		{"many - no clear", Floats{3, 4, 7, 1, 2, 5, 9, 8, 0, 6}, 4, 8, false},
		{"many - clear", Floats{3, 4, 7, 1, 2, 5, 9, 8, 0, 6}, 4, 8, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n := ResetFloats(tc.floats, tc.min, tc.max, tc.clear)
			require.NotNil(t, n)
			assert.GreaterOrEqual(t, cap(n), tc.min)
			assert.LessOrEqual(t, cap(n), tc.max)

			if tc.clear {
				l := cap(n)
				n = n[:l]
				for ix := range n {
					assert.Zero(t, n[ix])
				}
			}
		})
	}
}

func TestFloats_Replace(t *testing.T) {
	testCases := []struct {
		name   string
		floats Floats
		input  Floats
		want   Floats
	}{
		{"nil - nil", nil, nil, Floats{}},
		{"empty - nil", Floats{}, nil, Floats{}},
		{"nil - empty", nil, Floats{}, Floats{}},
		{"empty - empty", Floats{}, Floats{}, Floats{}},
		{"one - nil", Floats{1}, nil, Floats{}},
		{"one - empty", Floats{1}, Floats{}, Floats{}},
		{"one - one", Floats{2}, Floats{1}, Floats{1}},
		{"one - few", Floats{2}, Floats{3, 5, 2, 7}, Floats{3, 5, 2, 7}},
		{"one - many", Floats{3}, Floats{1, 5, 9, 8, 3, 4, 2, 6, 7}, Floats{1, 5, 9, 8, 3, 4, 2, 6, 7}},
		{"few - nil", Floats{3, 1, 2}, nil, Floats{}},
		{"few - empty", Floats{1, 3, 2}, Floats{}, Floats{}},
		{"few - one", Floats{2, 3, 1}, Floats{7}, Floats{7}},
		{"few - few", Floats{9, 8, 7}, Floats{1, 0, 8, 4}, Floats{1, 0, 8, 4}},
		{"few - many", Floats{7, 8, 9}, Floats{1, 5, 9, 8, 3, 4, 2, 6, 7}, Floats{1, 5, 9, 8, 3, 4, 2, 6, 7}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.floats.Replace(tc.input)
			if tc.input == nil {
				require.Nil(t, got)
			} else {
				require.NotNil(t, got)
				assert.Equal(t, tc.input.IsEmpty(), got.IsEmpty())
				assert.EqualValues(t, tc.want, got)
			}
		})
	}
}

func TestFloats_IsEmpty(t *testing.T) {
	testCases := []struct {
		name   string
		floats Floats
		want   bool
	}{
		{"nil", nil, true},
		{"empty", Floats{}, true},
		{"not empty", Floats{1, 2, 3}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.floats.IsEmpty()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFloats_EncodeDecode(t *testing.T) {
	testCases := []struct {
		name   string
		floats Floats
		want   string
	}{
		{"nil", nil, "[]"},
		{"empty", Floats{}, "[]"},
		{"one", Floats{1.23}, "[1.23]"},
		{"few", Floats{3.4, 6.66, 9.0, 12345.23}, "[3.4,6.66,9,12345.23]"},
		{"many", Floats{0, 7.01, 99.10, 99.12, 9999.11, 5.5, 7}, "[0,7.01,99.1,99.12,9999.11,5.5,7]"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			enc := zjson.AcquireEncoder(512)
			defer enc.Release()

			tc.floats.Encode(enc)
			j := enc.Bytes()
			require.NotNil(t, j)
			assert.Equal(t, tc.want, string(j))

			n := ResetFloats(nil, 8, 16, false)
			require.NotNil(t, n)
			assert.Equal(t, 8, cap(n))
			assert.Zero(t, len(n))

			dec := zjson.AcquireDecoder(j)
			defer dec.Release()

			err := n.Decode(dec)
			require.NoError(t, err)
			if tc.floats != nil {
				assert.EqualValues(t, tc.floats, n)
			}
		})
	}
}

func TestNewProducerFloats(t *testing.T) {
	testCases := []struct {
		name    string
		minSize int
		maxSize int
		clear   bool
		items   []float64
		minCap  int
		maxCap  int
	}{
		{
			name:    "small, no items, no clear",
			minSize: 4,
			maxSize: 16,
			minCap:  4,
			maxCap:  4,
		},
		{
			name:    "small, no items, clear",
			minSize: 4,
			maxSize: 16,
			clear:   true,
			minCap:  4,
			maxCap:  4,
		},
		{
			name:    "small, few items, no clear",
			minSize: 4,
			maxSize: 16,
			items:   []float64{1.1, -2.8, 3.4},
			minCap:  4,
			maxCap:  4,
		},
		{
			name:    "small, few items, clear",
			minSize: 4,
			maxSize: 16,
			clear:   true,
			items:   []float64{1.1, -2.8, 3.4},
			minCap:  4,
			maxCap:  4,
		},
		{
			name:    "small, many items, no clear",
			minSize: 4,
			maxSize: 16,
			items:   []float64{1.5, -2.2, 3.9, -4.2, 5.1, -6.6, 7.6, -8.3, 9.8, -10.2, 11.5},
			minCap:  12,
			maxCap:  16,
		},
		{
			name:    "small, many items, clear",
			minSize: 4,
			maxSize: 16,
			clear:   true,
			items:   []float64{1.5, -2.2, 3.9, -4.2, 5.1, -6.6, 7.6, -8.3, 9.8, -10.2, 11.5},
			minCap:  12,
			maxCap:  16,
		},
		{
			name:    "small, too many items, no clear",
			minSize: 4,
			maxSize: 16,
			items:   []float64{1.6, -2.2, 3.7, -4.1, 5.9, -6.5, 7.3, -8.4, 9.7, -10.3, 11.9, -12.4, 13.1, -14.6, 15.5, -16.3, 17.8, -18.4, 19.7},
			minCap:  20,
			maxCap:  32,
		},
		{
			name:    "small, too many items, clear",
			minSize: 4,
			maxSize: 16,
			clear:   true,
			items:   []float64{1.6, -2.2, 3.7, -4.1, 5.9, -6.5, 7.3, -8.4, 9.7, -10.3, 11.9, -12.4, 13.1, -14.6, 15.5, -16.3, 17.8, -18.4, 19.7},
			minCap:  20,
			maxCap:  32,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewFloatsProducer(tc.minSize, tc.maxSize, tc.clear, 'f', 2)
			require.NotNil(t, p)

			o1 := p.Acquire().(*FloatsManager)
			require.NotNil(t, o1)
			assert.Equal(t, tc.clear, o1.fullClear)
			assert.Equal(t, tc.minSize, o1.minSize)
			assert.Equal(t, tc.maxSize, o1.maxSize)
			assert.NotNil(t, o1.Items)
			assert.Zero(t, len(o1.Items))
			assert.True(t, o1.IsEmpty())

			o1.Append(tc.items...)

			assert.Equal(t, len(tc.items) == 0, o1.IsEmpty())
			assert.GreaterOrEqual(t, cap(o1.Items), tc.minCap)
			assert.LessOrEqual(t, cap(o1.Items), tc.maxCap)

			enc := zjson.AcquireEncoder(512)
			defer enc.Release()

			enc.Array(o1)
			j := enc.Bytes()
			assert.NotZero(t, len(j))

			n2, err := p.AcquireFromJSON(j)
			require.NoError(t, err)
			require.NotNil(t, n2)

			o2 := n2.(*FloatsManager)
			require.NotNil(t, o2)
			assert.EqualValues(t, o1.Items, o2.Items)

			o1.reset()
			assert.Zero(t, len(o1.Items))
			assert.True(t, o1.IsEmpty())

			if cap(o1.Items) >= len(tc.items) {
				o1.Items = o1.Items[:len(tc.items)]
				if tc.clear {
					for ix := range o1.Items {
						assert.Zero(t, o1.Items[ix])
					}
				} else {
					for ix := range o1.Items {
						assert.Equal(t, tc.items[ix], o1.Items[ix])
						o1.Items[ix] = 0
					}
				}
			}

			o2.Release()
			o1.Release()
		})
	}
}

func TestFloatsManager_Clone(t *testing.T) {
	testCases := []struct {
		name   string
		values []float64
	}{
		{"empty", []float64{}},
		{"1 value", []float64{66.5}},
		{"few values", []float64{16.0, -77.77, 200.1, 0.0}},
		{"many values", []float64{12.12, -9.0, 101.101, -45.4004, 678.0, 345.12, -1.0, 0.0, 34.12, -99.109, 15.6}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewFloatsProducer(8, 128, true, 'f', 2)
			require.NotNil(t, p)

			n1 := p.Acquire()
			require.NotNil(t, n1)

			o1, ok1 := n1.(*FloatsManager)
			require.True(t, ok1)
			require.NotNil(t, o1)
			assert.True(t, o1.IsEmpty())

			o1.Append(tc.values...)

			assert.EqualValues(t, tc.values, o1.Items)
			assert.Equal(t, len(tc.values) == 0, o1.IsEmpty())

			n2 := o1.Clone()
			require.NotNil(t, n2)

			o2, ok2 := n2.(*FloatsManager)
			require.True(t, ok2)
			require.NotNil(t, o2)

			o1.Release()

			assert.EqualValues(t, tc.values, o2.Items)
			assert.Equal(t, len(tc.values) == 0, o2.IsEmpty())

			o2.Release()
		})
	}
}
