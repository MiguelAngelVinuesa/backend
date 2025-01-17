package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSizes(t *testing.T) {
	const (
		t1 SizeType = iota + 1
		t2
		t3
	)

	t.Run("new sizes", func(t *testing.T) {
		d := NewSizes(t3)
		require.NotNil(t, d)

		d.Add(t1, 300)
		d.Add(t1, 100)

		d.Add(t2, 300)
		d.Add(t2, 200)
		d.Add(t2, 100)
		d.Add(t2, 200)

		d.Add(t3, 700)

		d0 := d.GetType(0)
		require.NotNil(t, d0)
		d1 := d.GetType(t1)
		require.NotNil(t, d1)
		d2 := d.GetType(t2)
		require.NotNil(t, d2)
		d3 := d.GetType(t3)
		require.NotNil(t, d3)

		assert.Equal(t, uint64(2), d1.Count())
		assert.Equal(t, uint64(200), d1.Avg())
		assert.Equal(t, uint64(4), d2.Count())
		assert.Equal(t, uint64(200), d2.Avg())
		assert.Equal(t, uint64(1), d3.Count())
		assert.Equal(t, uint64(700), d3.Avg())

		ds := d.GetAll()
		require.NotNil(t, ds)
		assert.Equal(t, 4, len(ds))

		d.Reset()
		assert.Zero(t, d1.Count())
		assert.Zero(t, d1.Avg())
		assert.Zero(t, d2.Count())
		assert.Zero(t, d2.Avg())
		assert.Zero(t, d3.Count())
		assert.Zero(t, d3.Avg())

		j, err := d.MarshalJSON()
		require.NoError(t, err)
		require.NotNil(t, j)
		assert.Equal(t, `[{"count":0,"total":0,"min":0,"max":0,"avg":0},{"count":0,"total":0,"min":0,"max":0,"avg":0},{"count":0,"total":0,"min":0,"max":0,"avg":0},{"count":0,"total":0,"min":0,"max":0,"avg":0}]`, string(j))
	})
}
