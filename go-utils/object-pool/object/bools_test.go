package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
)

func TestResetBools(t *testing.T) {
	testCases := []struct {
		name  string
		bools Bools
		min   int
		max   int
		clear bool
	}{
		{"nil", nil, 4, 8, true},
		{"empty", Bools{}, 4, 8, true},
		{"single - no clear", Bools{true}, 4, 8, false},
		{"single - clear", Bools{true}, 4, 8, true},
		{"few - no clear", Bools{true, false, true}, 4, 8, false},
		{"few - clear", Bools{true, false, true}, 4, 8, true},
		{"many - no clear", Bools{true, false, true, true, false, true, true, false, false, true}, 4, 8, false},
		{"many - clear", Bools{true, false, true, true, false, true, true, false, false, true}, 4, 8, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n := ResetBools(tc.bools, tc.min, tc.max, tc.clear)
			require.NotNil(t, n)
			assert.GreaterOrEqual(t, cap(n), tc.min)
			assert.LessOrEqual(t, cap(n), tc.max)

			if tc.clear {
				l := cap(n)
				n = n[:l]
				for ix := range n {
					assert.False(t, n[ix])
				}
			}
		})
	}
}

func TestBools_Replace(t *testing.T) {
	testCases := []struct {
		name  string
		bools Bools
		input Bools
		want  Bools
	}{
		{"nil - nil", nil, nil, Bools{}},
		{"empty - nil", Bools{}, nil, Bools{}},
		{"nil - empty", nil, Bools{}, Bools{}},
		{"empty - empty", Bools{}, Bools{}, Bools{}},
		{"one - nil", Bools{true}, nil, Bools{}},
		{"one - empty", Bools{true}, Bools{}, Bools{}},
		{"one - one", Bools{true}, Bools{false}, Bools{false}},
		{"one - few", Bools{true}, Bools{true, false, true, true}, Bools{true, false, true, true}},
		{"one - many", Bools{true}, Bools{true, false, true, true, false, true, true, false, true}, Bools{true, false, true, true, false, true, true, false, true}},
		{"few - nil", Bools{false, false, true}, nil, Bools{}},
		{"few - empty", Bools{false, false, true}, Bools{}, Bools{}},
		{"few - one", Bools{false, false, true}, Bools{true}, Bools{true}},
		{"few - few", Bools{false, false, true}, Bools{true, false, true, true}, Bools{true, false, true, true}},
		{"few - many", Bools{false, false, true}, Bools{true, false, true, true, false, true, true, false, true}, Bools{true, false, true, true, false, true, true, false, true}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.bools.Replace(tc.input)
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

func TestBools_IsEmpty(t *testing.T) {
	testCases := []struct {
		name  string
		bools Bools
		want  bool
	}{
		{"nil", nil, true},
		{"empty", Bools{}, true},
		{"not empty", Bools{true, false, true}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.bools.IsEmpty()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestBools_EncodeDecode(t *testing.T) {
	testCases := []struct {
		name  string
		bools Bools
		want  string
	}{
		{"nil", nil, "[]"},
		{"empty", Bools{}, "[]"},
		{"one", Bools{true}, "[1]"},
		{"few", Bools{true, false, true, false}, "[1,0,1,0]"},
		{"many", Bools{true, true, false, true, false, false, true, false}, "[1,1,0,1,0,0,1,0]"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			enc := zjson.AcquireEncoder(512)
			defer enc.Release()

			tc.bools.Encode(enc)
			j := enc.Bytes()
			require.NotNil(t, j)
			assert.Equal(t, tc.want, string(j))

			n := ResetBools(nil, 8, 16, false)
			require.NotNil(t, n)
			assert.Equal(t, 8, cap(n))
			assert.Zero(t, len(n))

			dec := zjson.AcquireDecoder(j)
			defer dec.Release()

			err := n.Decode(dec)
			require.NoError(t, err)
			if tc.bools != nil {
				assert.EqualValues(t, tc.bools, n)
			}
		})
	}
}

func TestNewProducerBools(t *testing.T) {
	testCases := []struct {
		name    string
		minSize int
		maxSize int
		clear   bool
		items   []bool
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
			items:   []bool{true, false, true},
			minCap:  4,
			maxCap:  4,
		},
		{
			name:    "small, few items, clear",
			minSize: 4,
			maxSize: 16,
			clear:   true,
			items:   []bool{true, false, true},
			minCap:  4,
			maxCap:  4,
		},
		{
			name:    "small, many items, no clear",
			minSize: 4,
			maxSize: 16,
			items:   []bool{true, true, false, true, false, false, true, false, true, true, true},
			minCap:  12,
			maxCap:  16,
		},
		{
			name:    "small, many items, clear",
			minSize: 4,
			maxSize: 16,
			clear:   true,
			items:   []bool{true, true, false, true, false, false, true, false, true, true, true},
			minCap:  12,
			maxCap:  16,
		},
		{
			name:    "small, too many items, no clear",
			minSize: 4,
			maxSize: 16,
			items:   []bool{true, true, false, true, false, false, true, false, true, true, true, true, false, false, true, false, true, true, true},
			minCap:  20,
			maxCap:  32,
		},
		{
			name:    "small, too many items, clear",
			minSize: 4,
			maxSize: 16,
			clear:   true,
			items:   []bool{true, true, false, true, false, false, true, false, true, true, true, true, false, false, true, false, true, true, true},
			minCap:  20,
			maxCap:  32,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewBoolsProducer(tc.minSize, tc.maxSize, tc.clear)
			require.NotNil(t, p)

			o1 := p.Acquire().(*BoolsManager)
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
			if len(tc.items) > 0 {
				assert.EqualValues(t, tc.items, o1.Items)
			}

			enc := zjson.AcquireEncoder(512)
			defer enc.Release()

			enc.Array(o1)
			j := enc.Bytes()
			assert.NotZero(t, len(j))

			n2, err := p.AcquireFromJSON(j)
			require.NoError(t, err)
			require.NotNil(t, n2)

			o2 := n2.(*BoolsManager)
			require.NotNil(t, o2)
			assert.EqualValues(t, o1.Items, o2.Items)
			if len(tc.items) > 0 {
				assert.EqualValues(t, tc.items, o2.Items)
			}

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
					}
				}
			}

			o2.Release()
			o1.Release()
		})
	}
}

func TestBoolsManager_Clone(t *testing.T) {
	testCases := []struct {
		name   string
		values []bool
	}{
		{"empty", []bool{}},
		{"1 value", []bool{true}},
		{"few values", []bool{true, false, false, true}},
		{"many values", []bool{false, true, true, false, false, true, true, false, true, false, true}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewBoolsProducer(8, 64, true)
			require.NotNil(t, p)

			n1 := p.Acquire()
			require.NotNil(t, n1)

			o1, ok1 := n1.(*BoolsManager)
			require.True(t, ok1)
			require.NotNil(t, o1)
			assert.True(t, o1.IsEmpty())

			o1.Append(tc.values...)

			assert.EqualValues(t, tc.values, o1.Items)
			assert.Equal(t, len(tc.values) == 0, o1.IsEmpty())

			n2 := o1.Clone()
			require.NotNil(t, n2)

			o2, ok2 := n2.(*BoolsManager)
			require.True(t, ok2)
			require.NotNil(t, o2)

			o1.Release()

			assert.EqualValues(t, tc.values, o2.Items)
			assert.Equal(t, len(tc.values) == 0, o2.IsEmpty())

			o2.Release()
		})
	}
}
