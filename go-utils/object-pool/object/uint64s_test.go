package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
)

func TestResetUint64s(t *testing.T) {
	testCases := []struct {
		name  string
		ints  Uint64s
		min   int
		max   int
		clear bool
	}{
		{"nil", nil, 4, 8, true},
		{"empty", Uint64s{}, 4, 8, true},
		{"single - no clear", Uint64s{1}, 4, 8, false},
		{"single - clear", Uint64s{1}, 4, 8, true},
		{"few - no clear", Uint64s{2, 3, 1}, 4, 8, false},
		{"few - clear", Uint64s{3, 1, 2}, 4, 8, true},
		{"many - no clear", Uint64s{3, 4, 7, 1, 2, 5, 9, 8, 0, 6}, 4, 8, false},
		{"many - clear", Uint64s{3, 4, 7, 1, 2, 5, 9, 8, 0, 6}, 4, 8, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n := ResetUint64s(tc.ints, tc.min, tc.max, tc.clear)
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

func TestUint64s_Replace(t *testing.T) {
	testCases := []struct {
		name  string
		ints  Uint64s
		input Uint64s
		want  Uint64s
	}{
		{"nil - nil", nil, nil, Uint64s{}},
		{"empty - nil", Uint64s{}, nil, Uint64s{}},
		{"nil - empty", nil, Uint64s{}, Uint64s{}},
		{"empty - empty", Uint64s{}, Uint64s{}, Uint64s{}},
		{"one - nil", Uint64s{1}, nil, Uint64s{}},
		{"one - empty", Uint64s{1}, Uint64s{}, Uint64s{}},
		{"one - one", Uint64s{2}, Uint64s{1}, Uint64s{1}},
		{"one - few", Uint64s{2}, Uint64s{3, 5, 2, 7}, Uint64s{3, 5, 2, 7}},
		{"one - many", Uint64s{3}, Uint64s{1, 5, 9, 8, 3, 4, 2, 6, 7}, Uint64s{1, 5, 9, 8, 3, 4, 2, 6, 7}},
		{"few - nil", Uint64s{3, 1, 2}, nil, Uint64s{}},
		{"few - empty", Uint64s{1, 3, 2}, Uint64s{}, Uint64s{}},
		{"few - one", Uint64s{2, 3, 1}, Uint64s{7}, Uint64s{7}},
		{"few - few", Uint64s{9, 8, 7}, Uint64s{1, 0, 8, 4}, Uint64s{1, 0, 8, 4}},
		{"few - many", Uint64s{7, 8, 9}, Uint64s{1, 5, 9, 8, 3, 4, 2, 6, 7}, Uint64s{1, 5, 9, 8, 3, 4, 2, 6, 7}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.ints.Replace(tc.input)
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

func TestUint64s_IsEmpty(t *testing.T) {
	testCases := []struct {
		name string
		ints Uint64s
		want bool
	}{
		{"nil", nil, true},
		{"empty", Uint64s{}, true},
		{"not empty", Uint64s{1, 2, 3}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.ints.IsEmpty()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestUint64s_EncodeDecode(t *testing.T) {
	testCases := []struct {
		name string
		ints Uint64s
		want string
	}{
		{"nil", nil, "[]"},
		{"empty", Uint64s{}, "[]"},
		{"one", Uint64s{1}, "[1]"},
		{"few", Uint64s{3, 6, 9, 12}, "[3,6,9,12]"},
		{"many", Uint64s{0, 7, 99, 99, 3, 1, 11, 55, 7}, "[0,7,99,99,3,1,11,55,7]"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			enc := zjson.AcquireEncoder(512)
			defer enc.Release()

			tc.ints.Encode(enc)
			j := enc.Bytes()
			require.NotNil(t, j)
			assert.Equal(t, tc.want, string(j))

			n := ResetUint64s(nil, 8, 16, false)
			require.NotNil(t, n)
			assert.Equal(t, 8, cap(n))
			assert.Zero(t, len(n))

			dec := zjson.AcquireDecoder(j)
			defer dec.Release()

			err := n.Decode(dec)
			require.NoError(t, err)
			if tc.ints != nil {
				assert.EqualValues(t, tc.ints, n)
			}
		})
	}
}

func TestNewUint64sProducer(t *testing.T) {
	testCases := []struct {
		name    string
		minSize int
		maxSize int
		clear   bool
		items   []uint64
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
			items:   []uint64{1, 2, 3},
			minCap:  4,
			maxCap:  4,
		},
		{
			name:    "small, few items, clear",
			minSize: 4,
			maxSize: 16,
			clear:   true,
			items:   []uint64{1, 2, 3},
			minCap:  4,
			maxCap:  4,
		},
		{
			name:    "small, many items, no clear",
			minSize: 4,
			maxSize: 16,
			items:   []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			minCap:  12,
			maxCap:  16,
		},
		{
			name:    "small, many items, clear",
			minSize: 4,
			maxSize: 16,
			clear:   true,
			items:   []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			minCap:  12,
			maxCap:  16,
		},
		{
			name:    "small, too many items, no clear",
			minSize: 4,
			maxSize: 16,
			items:   []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
			minCap:  20,
			maxCap:  32,
		},
		{
			name:    "small, too many items, clear",
			minSize: 4,
			maxSize: 16,
			clear:   true,
			items:   []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
			minCap:  20,
			maxCap:  32,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewUint64sProducer(tc.minSize, tc.maxSize, tc.clear)
			require.NotNil(t, p)

			o1, ok := p.Acquire().(*Uint64sManager)
			require.True(t, ok)
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

			o2 := n2.(*Uint64sManager)
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

func TestUint64sManager_Clone(t *testing.T) {
	testCases := []struct {
		name   string
		values []uint64
	}{
		{"empty", []uint64{}},
		{"1 value", []uint64{66}},
		{"few values", []uint64{16, 20277, 60000, 0}},
		{"many values", []uint64{12, 9, 101, 45, 12378, 545, 1, 0, 34, 99, 34215}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewUint64sProducer(8, 128, true)
			require.NotNil(t, p)

			n1 := p.Acquire()
			require.NotNil(t, n1)

			o1, ok1 := n1.(*Uint64sManager)
			require.True(t, ok1)
			require.NotNil(t, o1)
			assert.True(t, o1.IsEmpty())

			o1.Append(tc.values...)

			assert.EqualValues(t, tc.values, o1.Items)
			assert.Equal(t, len(tc.values) == 0, o1.IsEmpty())

			n2 := o1.Clone()
			require.NotNil(t, n2)

			o2, ok2 := n2.(*Uint64sManager)
			require.True(t, ok2)
			require.NotNil(t, o2)

			o1.Release()

			assert.EqualValues(t, tc.values, o2.Items)
			assert.Equal(t, len(tc.values) == 0, o2.IsEmpty())

			o2.Release()
		})
	}
}
