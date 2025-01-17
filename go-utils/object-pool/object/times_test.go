package object

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
)

func TestResetTimes(t *testing.T) {
	testCases := []struct {
		name  string
		ints  Times
		min   int
		max   int
		clear bool
	}{
		{"nil", nil, 4, 8, true},
		{"empty", Times{}, 4, 8, true},
		{"single - no clear", Times{1}, 4, 8, false},
		{"single - clear", Times{1}, 4, 8, true},
		{"few - no clear", Times{2, 3, 1}, 4, 8, false},
		{"few - clear", Times{3, 1, 2}, 4, 8, true},
		{"many - no clear", Times{3, 4, 7, 1, 2, 5, 9, 8, 0, 6}, 4, 8, false},
		{"many - clear", Times{3, 4, 7, 1, 2, 5, 9, 8, 0, 6}, 4, 8, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n := ResetTimes(tc.ints, tc.min, tc.max, tc.clear)
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

func TestTimes_Replace(t *testing.T) {
	testCases := []struct {
		name  string
		times Times
		input Times
		want  Times
	}{
		{"nil - nil", nil, nil, Times{}},
		{"empty - nil", Times{}, nil, Times{}},
		{"nil - empty", nil, Times{}, Times{}},
		{"empty - empty", Times{}, Times{}, Times{}},
		{"one - nil", Times{1}, nil, Times{}},
		{"one - empty", Times{1}, Times{}, Times{}},
		{"one - one", Times{2}, Times{1}, Times{1}},
		{"one - few", Times{2}, Times{3, 5, 2, 7}, Times{3, 5, 2, 7}},
		{"one - many", Times{3}, Times{1, 5, 9, 8, 3, 4, 2, 6, 7}, Times{1, 5, 9, 8, 3, 4, 2, 6, 7}},
		{"few - nil", Times{3, 1, 2}, nil, Times{}},
		{"few - empty", Times{1, 3, 2}, Times{}, Times{}},
		{"few - one", Times{2, 3, 1}, Times{7}, Times{7}},
		{"few - few", Times{9, 8, 7}, Times{1, 0, 8, 4}, Times{1, 0, 8, 4}},
		{"few - many", Times{7, 8, 9}, Times{1, 5, 9, 8, 3, 4, 2, 6, 7}, Times{1, 5, 9, 8, 3, 4, 2, 6, 7}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.times.Replace(tc.input)
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

func TestTimes_IsEmpty(t *testing.T) {
	testCases := []struct {
		name  string
		times Times
		want  bool
	}{
		{"nil", nil, true},
		{"empty", Times{}, true},
		{"not empty", Times{1, 2, 3}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.times.IsEmpty()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestTimes_EncodeDecode(t *testing.T) {
	testCases := []struct {
		name  string
		times Times
		want  string
	}{
		{"nil", nil, "[]"},
		{"empty", Times{}, "[]"},
		{"one", Times{1}, "[1]"},
		{"few", Times{3, 6, 9, 12}, "[3,6,9,12]"},
		{"many", Times{0, 7, 99, 99, 3, 1, 9911, 55, 7}, "[0,7,99,99,3,1,9911,55,7]"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			enc := zjson.AcquireEncoder(512)
			defer enc.Release()

			tc.times.Encode(enc)
			j := enc.Bytes()
			require.NotNil(t, j)
			assert.Equal(t, tc.want, string(j))

			n := ResetTimes(nil, 8, 16, false)
			require.NotNil(t, n)
			assert.Equal(t, 8, cap(n))
			assert.Zero(t, len(n))

			dec := zjson.AcquireDecoder(j)
			defer dec.Release()

			err := n.Decode(dec)
			require.NoError(t, err)
			if tc.times != nil {
				assert.EqualValues(t, tc.times, n)
			}
		})
	}
}

func TestNewProducerTimes(t *testing.T) {
	now := time.Now().UTC().Round(time.Millisecond)

	testCases := []struct {
		name    string
		minSize int
		maxSize int
		clear   bool
		items   []time.Time
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
			name:    "small, one negative item, no clear",
			minSize: 4,
			maxSize: 16,
			items:   []time.Time{time.Date(1960, 1, 1, 0, 0, 0, 0, time.UTC)},
			minCap:  4,
			maxCap:  4,
		},
		{
			name:    "small, one negative item, clear",
			minSize: 4,
			maxSize: 16,
			clear:   true,
			items:   []time.Time{time.Date(1960, 1, 1, 0, 0, 0, 0, time.UTC)},
			minCap:  4,
			maxCap:  4,
		},
		{
			name:    "small, few empty items, no clear",
			minSize: 4,
			maxSize: 16,
			items:   []time.Time{{}, {}, {}},
			minCap:  4,
			maxCap:  4,
		},
		{
			name:    "small, few empty items, clear",
			minSize: 4,
			maxSize: 16,
			clear:   true,
			items:   []time.Time{{}, {}, {}},
			minCap:  4,
			maxCap:  4,
		},
		{
			name:    "small, few negative items, no clear",
			minSize: 4,
			maxSize: 16,
			items: []time.Time{
				time.Date(1940, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(1960, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			minCap: 4,
			maxCap: 4,
		},
		{
			name:    "small, few negative items, clear",
			minSize: 4,
			maxSize: 16,
			clear:   true,
			items: []time.Time{
				time.Date(1940, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(1960, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			minCap: 4,
			maxCap: 4,
		},
		{
			name:    "small, few items, no clear",
			minSize: 4,
			maxSize: 16,
			items:   []time.Time{now, now.Add(time.Second), now.Add(2 * time.Second)},
			minCap:  4,
			maxCap:  4,
		},
		{
			name:    "small, few items, clear",
			minSize: 4,
			maxSize: 16,
			clear:   true,
			items:   []time.Time{now, now.Add(time.Second), now.Add(2 * time.Second)},
			minCap:  4,
			maxCap:  4,
		},
		{
			name:    "small, many items, no clear",
			minSize: 4,
			maxSize: 16,
			items: []time.Time{now, now.Add(time.Second), now.Add(2 * time.Second), now.Add(3 * time.Second), now.Add(4 * time.Second),
				now.Add(5 * time.Second), now.Add(6 * time.Second), now.Add(7 * time.Second), now.Add(8 * time.Second), now.Add(9 * time.Second),
				now.Add(10 * time.Second)},
			minCap: 12,
			maxCap: 16,
		},
		{
			name:    "small, many items, clear",
			minSize: 4,
			maxSize: 16,
			clear:   true,
			items: []time.Time{now, now.Add(time.Second), now.Add(2 * time.Second), now.Add(3 * time.Second), now.Add(4 * time.Second),
				now.Add(5 * time.Second), now.Add(6 * time.Second), now.Add(7 * time.Second), now.Add(8 * time.Second), now.Add(9 * time.Second),
				now.Add(10 * time.Second)},
			minCap: 12,
			maxCap: 16,
		},
		{
			name:    "small, too many items, no clear",
			minSize: 4,
			maxSize: 16,
			items: []time.Time{now, now.Add(time.Second), now.Add(2 * time.Second), now.Add(3 * time.Second), now.Add(4 * time.Second),
				now.Add(5 * time.Second), now.Add(6 * time.Second), now.Add(7 * time.Second), now.Add(8 * time.Second), now.Add(9 * time.Second),
				now.Add(10 * time.Second), now.Add(11 * time.Second), now.Add(12 * time.Second), now.Add(13 * time.Second), now.Add(14 * time.Second),
				now.Add(15 * time.Second), now.Add(16 * time.Second), now.Add(17 * time.Second), now.Add(18 * time.Second), now.Add(19 * time.Second)},
			minCap: 20,
			maxCap: 32,
		},
		{
			name:    "small, too many items, clear",
			minSize: 4,
			maxSize: 16,
			clear:   true,
			items: []time.Time{now, now.Add(time.Second), now.Add(2 * time.Second), now.Add(3 * time.Second), now.Add(4 * time.Second),
				now.Add(5 * time.Second), now.Add(6 * time.Second), now.Add(7 * time.Second), now.Add(8 * time.Second), now.Add(9 * time.Second),
				now.Add(10 * time.Second), now.Add(11 * time.Second), now.Add(12 * time.Second), now.Add(13 * time.Second), now.Add(14 * time.Second),
				now.Add(15 * time.Second), now.Add(16 * time.Second), now.Add(17 * time.Second), now.Add(18 * time.Second), now.Add(19 * time.Second)},
			minCap: 20,
			maxCap: 32,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewTimesProducer(tc.minSize, tc.maxSize, tc.clear)
			require.NotNil(t, p)

			o1 := p.Acquire().(*TimesManager)
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

			o2 := n2.(*TimesManager)
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
						assert.Equal(t, tc.items[ix], o1.Value(ix))
						o1.Items[ix] = 0
					}
				}
			}

			o2.Release()
			o1.Release()
		})
	}
}

func TestTimesManager_Clone(t *testing.T) {
	now := time.Now().UTC().Round(time.Millisecond)

	testCases := []struct {
		name   string
		values []time.Time
	}{
		{"empty", []time.Time{}},
		{"1 value", []time.Time{now}},
		{"few values", []time.Time{now, now, now, now}},
		{"many values", []time.Time{now, now, now, now, now, now, now, now, now, now, now}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewTimesProducer(8, 64, true)
			require.NotNil(t, p)

			n1 := p.Acquire()
			require.NotNil(t, n1)

			o1, ok1 := n1.(*TimesManager)
			require.True(t, ok1)
			require.NotNil(t, o1)
			assert.True(t, o1.IsEmpty())

			o1.Append(tc.values...)

			require.Equal(t, len(tc.values), len(o1.Items))
			for ix := range tc.values {
				assert.Equal(t, tc.values[ix], o1.Value(ix))
			}
			assert.Equal(t, len(tc.values) == 0, o1.IsEmpty())

			n2 := o1.Clone()
			require.NotNil(t, n2)

			o2, ok2 := n2.(*TimesManager)
			require.True(t, ok2)
			require.NotNil(t, o2)

			o1.Release()

			require.Equal(t, len(tc.values), len(o2.Items))
			for ix := range tc.values {
				assert.Equal(t, tc.values[ix], o2.Value(ix))
			}
			assert.Equal(t, len(tc.values) == 0, o2.IsEmpty())

			o2.Release()
		})
	}
}
