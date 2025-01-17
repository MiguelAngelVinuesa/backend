package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDuration(t *testing.T) {
	testCases := []struct {
		name   string
		inputs []time.Duration
		count  uint64
		total  time.Duration
		min    time.Duration
		max    time.Duration
		avg    time.Duration
		j      string
	}{
		{
			name: "no inputs",
			j:    `{"count":0,"total":0,"min":0,"max":0,"avg":0}`,
		},
		{
			name:   "1 input",
			inputs: []time.Duration{time.Second},
			count:  1,
			total:  time.Second,
			min:    time.Second,
			max:    time.Second,
			avg:    time.Second,
			j:      `{"count":1,"total":1000000000,"min":1000000000,"max":1000000000,"avg":1000000000}`,
		},
		{
			name:   "few inputs",
			inputs: []time.Duration{time.Second, 2 * time.Second, 3 * time.Second, 3 * time.Second, time.Second},
			count:  5,
			total:  10 * time.Second,
			min:    time.Second,
			max:    3 * time.Second,
			avg:    2 * time.Second,
			j:      `{"count":5,"total":10000000000,"min":1000000000,"max":3000000000,"avg":2000000000}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := NewDuration()
			require.NotNil(t, d)
			assert.Zero(t, d.Count())
			assert.Zero(t, d.Total())
			assert.Zero(t, d.Min())
			assert.Zero(t, d.Max())
			assert.Zero(t, d.Avg())

			for _, in := range tc.inputs {
				d.Add(in)
			}

			assert.Equal(t, tc.count, d.Count())
			assert.Equal(t, tc.total, d.Total())
			assert.Equal(t, tc.min, d.Min())
			assert.Equal(t, tc.max, d.Max())
			assert.Equal(t, tc.avg, d.Avg())

			j, err := d.MarshalJSON()
			require.NoError(t, err)
			require.NotNil(t, j)
			assert.Equal(t, tc.j, string(j))

			d.Reset()
			assert.Zero(t, d.Count())
			assert.Zero(t, d.Total())
			assert.Zero(t, d.Min())
			assert.Zero(t, d.Max())
			assert.Zero(t, d.Avg())
		})
	}
}
