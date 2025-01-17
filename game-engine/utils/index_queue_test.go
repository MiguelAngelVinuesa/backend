package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewIndexQueue(t *testing.T) {
	testCases := []struct {
		name     string
		capacity int
	}{
		{name: "small", capacity: 10},
		{name: "big", capacity: 1000},
		{name: "medium", capacity: 200},
		{name: "very small", capacity: 5},
		{name: "very big", capacity: 10000},
		{name: "default", capacity: DefaultIndexQueueCap},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := NewIndexQueue(tc.capacity)
			require.NotNil(t, q)
			defer q.Release()

			require.NotNil(t, q.buf)
			assert.Equal(t, tc.capacity, q.bufMax)
			assert.Zero(t, q.putPtr)
			assert.Zero(t, q.getPtr)
			assert.Zero(t, q.getMax)
			assert.Zero(t, q.Available())

			i, ok := q.Get()
			assert.Zero(t, i)
			assert.False(t, ok)
		})
	}
}

func TestIndexQueue_PutGet(t *testing.T) {
	testCases := []struct {
		name     string
		capacity int
		indexes  Indexes
	}{
		{
			name:     "small, no wrap",
			capacity: 5,
			indexes:  Indexes{2, 1, 4, 3},
		},
		{
			name:     "small, wrap once",
			capacity: 5,
			indexes:  Indexes{4, 5, 6, 7, 1, 2, 3},
		},
		{
			name:     "small, wrap 3 times",
			capacity: 5,
			indexes:  Indexes{8, 9, 10, 17, 18, 19, 11, 12, 2, 3, 1, 4, 6, 7, 13, 14, 5, 15, 16},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := NewIndexQueue(tc.capacity)
			require.NotNil(t, q)
			defer q.Release()

			got := make(Indexes, 0, len(tc.indexes))
			for ix := range tc.indexes {
				q.Put(tc.indexes[ix])
				assert.NotZero(t, q.Available())

				i, ok := q.Get()
				require.True(t, ok)
				assert.Zero(t, q.Available())

				got = append(got, i)
			}

			assert.EqualValues(t, tc.indexes, got)
		})
	}
}

func TestIndexQueue_Put4Get4(t *testing.T) {
	testCases := []struct {
		name     string
		capacity int
		indexes  Indexes
	}{
		{
			name:     "small, no wrap",
			capacity: 5,
			indexes:  Indexes{2, 1, 4, 3},
		},
		{
			name:     "small, wrap once",
			capacity: 5,
			indexes:  Indexes{4, 5, 6, 7, 1, 2, 3},
		},
		{
			name:     "small, wrap 3 times",
			capacity: 5,
			indexes:  Indexes{8, 9, 10, 17, 18, 19, 11, 12, 2, 3, 1, 4, 6, 7, 13, 14, 5, 15, 16},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := NewIndexQueue(tc.capacity)
			require.NotNil(t, q)
			defer q.Release()

			l := len(tc.indexes)
			got := make(Indexes, 0, l)
			for ix := 0; ix < l; ix++ {
				q.Put(tc.indexes[ix])
				assert.NotZero(t, q.Available())

				if ix%4 == 3 || ix == l-1 {
					for iy := 0; iy <= ix%4; iy++ {
						i, ok := q.Get()
						require.True(t, ok)
						got = append(got, i)
					}

					i, ok := q.Get()
					assert.Zero(t, i)
					assert.False(t, ok)
				}
			}

			assert.EqualValues(t, tc.indexes, got)
		})
	}
}

func TestIndexQueue_Reset(t *testing.T) {
	t.Run("reset", func(t *testing.T) {
		q := NewIndexQueue(10)
		require.NotNil(t, q)
		defer q.Release()

		for ix := 0; ix < 5; ix++ {
			q.Put(1)
		}
		require.NotZero(t, q.Available())

		q2 := q.Reset()
		require.NotNil(t, q2)
		assert.Equal(t, q, q2)
		assert.Zero(t, q.Available())
	})
}
