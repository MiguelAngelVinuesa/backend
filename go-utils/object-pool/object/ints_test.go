package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
)

func TestResetInts(t *testing.T) {
	testCases := []struct {
		name  string
		ints  Ints
		min   int
		max   int
		clear bool
	}{
		{"nil", nil, 4, 8, true},
		{"empty", Ints{}, 4, 8, true},
		{"single - no clear", Ints{1}, 4, 8, false},
		{"single - clear", Ints{1}, 4, 8, true},
		{"few - no clear", Ints{2, 3, 1}, 4, 8, false},
		{"few - clear", Ints{3, 1, 2}, 4, 8, true},
		{"many - no clear", Ints{3, 4, 7, 1, 2, 5, 9, 8, 0, 6}, 4, 8, false},
		{"many - clear", Ints{3, 4, 7, 1, 2, 5, 9, 8, 0, 6}, 4, 8, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n := ResetInts(tc.ints, tc.min, tc.max, tc.clear)
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

func TestInts_Replace(t *testing.T) {
	testCases := []struct {
		name  string
		ints  Ints
		input Ints
		want  Ints
	}{
		{"nil - nil", nil, nil, Ints{}},
		{"empty - nil", Ints{}, nil, Ints{}},
		{"nil - empty", nil, Ints{}, Ints{}},
		{"empty - empty", Ints{}, Ints{}, Ints{}},
		{"one - nil", Ints{1}, nil, Ints{}},
		{"one - empty", Ints{1}, Ints{}, Ints{}},
		{"one - one", Ints{2}, Ints{1}, Ints{1}},
		{"one - few", Ints{2}, Ints{3, 5, 2, 7}, Ints{3, 5, 2, 7}},
		{"one - many", Ints{3}, Ints{1, 5, 9, 8, 3, 4, 2, 6, 7}, Ints{1, 5, 9, 8, 3, 4, 2, 6, 7}},
		{"few - nil", Ints{3, 1, 2}, nil, Ints{}},
		{"few - empty", Ints{1, 3, 2}, Ints{}, Ints{}},
		{"few - one", Ints{2, 3, 1}, Ints{7}, Ints{7}},
		{"few - few", Ints{9, 8, 7}, Ints{1, 0, 8, 4}, Ints{1, 0, 8, 4}},
		{"few - many", Ints{7, 8, 9}, Ints{1, 5, 9, 8, 3, 4, 2, 6, 7}, Ints{1, 5, 9, 8, 3, 4, 2, 6, 7}},
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

func TestInts_IsEmpty(t *testing.T) {
	testCases := []struct {
		name string
		ints Ints
		want bool
	}{
		{"nil", nil, true},
		{"empty", Ints{}, true},
		{"not empty", Ints{1, 2, 3}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.ints.IsEmpty()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestInts_EncodeDecode(t *testing.T) {
	testCases := []struct {
		name string
		ints Ints
		want string
	}{
		{"nil", nil, "[]"},
		{"empty", Ints{}, "[]"},
		{"one", Ints{1}, "[1]"},
		{"few", Ints{3, 6, 9, 12}, "[3,6,9,12]"},
		{"many", Ints{0, 7, 99, 99, 3, 1, 9911, 55, 7}, "[0,7,99,99,3,1,9911,55,7]"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			enc := zjson.AcquireEncoder(512)
			defer enc.Release()

			tc.ints.Encode(enc)
			j := enc.Bytes()
			require.NotNil(t, j)
			assert.Equal(t, tc.want, string(j))

			n := ResetInts(nil, 8, 16, false)
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

func TestNewIntsProducer(t *testing.T) {
	testCases := []struct {
		name    string
		minSize int
		maxSize int
		clear   bool
		items   []int
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
			items:   []int{1, 2, 3},
			minCap:  4,
			maxCap:  4,
		},
		{
			name:    "small, few items, clear",
			minSize: 4,
			maxSize: 16,
			clear:   true,
			items:   []int{1, 2, 3},
			minCap:  4,
			maxCap:  4,
		},
		{
			name:    "small, many items, no clear",
			minSize: 4,
			maxSize: 16,
			items:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			minCap:  12,
			maxCap:  16,
		},
		{
			name:    "small, many items, clear",
			minSize: 4,
			maxSize: 16,
			clear:   true,
			items:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			minCap:  12,
			maxCap:  16,
		},
		{
			name:    "small, too many items, no clear",
			minSize: 4,
			maxSize: 16,
			items:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
			minCap:  20,
			maxCap:  32,
		},
		{
			name:    "small, too many items, clear",
			minSize: 4,
			maxSize: 16,
			clear:   true,
			items:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
			minCap:  20,
			maxCap:  32,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewIntsProducer(tc.minSize, tc.maxSize, tc.clear)
			require.NotNil(t, p)

			o1, ok := p.Acquire().(*IntsManager)
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

			o2 := n2.(*IntsManager)
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

func TestIntsManager_Clone(t *testing.T) {
	testCases := []struct {
		name   string
		values []int
	}{
		{"empty", []int{}},
		{"1 value", []int{66}},
		{"few values", []int{16, -202775678, 60000000, 0}},
		{"many values", []int{12, -9, 1010101010101, -45, -12378, 545, 1, 0, 34, -99, -3421598765}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewIntsProducer(8, 128, true)
			require.NotNil(t, p)

			n1 := p.Acquire()
			require.NotNil(t, n1)

			o1, ok1 := n1.(*IntsManager)
			require.True(t, ok1)
			require.NotNil(t, o1)
			assert.True(t, o1.IsEmpty())

			o1.Append(tc.values...)

			assert.EqualValues(t, tc.values, o1.Items)
			assert.Equal(t, len(tc.values) == 0, o1.IsEmpty())

			n2 := o1.Clone()
			require.NotNil(t, n2)

			o2, ok2 := n2.(*IntsManager)
			require.True(t, ok2)
			require.NotNil(t, o2)

			o1.Release()

			assert.EqualValues(t, tc.values, o2.Items)
			assert.Equal(t, len(tc.values) == 0, o2.IsEmpty())

			o2.Release()
		})
	}
}
