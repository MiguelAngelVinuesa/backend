package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSize(t *testing.T) {
	testCases := []struct {
		name   string
		inputs []uint64
		count  uint64
		total  uint64
		min    uint64
		max    uint64
		avg    uint64
		j      string
	}{
		{
			name: "no inputs",
			j:    `{"count":0,"total":0,"min":0,"max":0,"avg":0}`,
		},
		{
			name:   "1 input",
			inputs: []uint64{100},
			count:  1,
			total:  100,
			min:    100,
			max:    100,
			avg:    100,
			j:      `{"count":1,"total":100,"min":100,"max":100,"avg":100}`,
		},
		{
			name:   "few inputs",
			inputs: []uint64{100, 200, 100, 300, 200, 400, 100},
			count:  7,
			total:  1400,
			min:    100,
			max:    400,
			avg:    200,
			j:      `{"count":7,"total":1400,"min":100,"max":400,"avg":200}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSize()
			require.NotNil(t, s)
			assert.Zero(t, s.Count())
			assert.Zero(t, s.Total())
			assert.Zero(t, s.Min())
			assert.Zero(t, s.Max())
			assert.Zero(t, s.Avg())

			for _, in := range tc.inputs {
				s.Add(in)
			}

			assert.Equal(t, tc.count, s.Count())
			assert.Equal(t, tc.total, s.Total())
			assert.Equal(t, tc.min, s.Min())
			assert.Equal(t, tc.max, s.Max())
			assert.Equal(t, tc.avg, s.Avg())

			j, err := s.MarshalJSON()
			require.NoError(t, err)
			require.NotNil(t, j)
			assert.Equal(t, tc.j, string(j))

			s.Reset()
			assert.Zero(t, s.Count())
			assert.Zero(t, s.Total())
			assert.Zero(t, s.Min())
			assert.Zero(t, s.Max())
			assert.Zero(t, s.Avg())
		})
	}
}
