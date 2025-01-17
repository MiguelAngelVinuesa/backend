package metrics

import (
	"fmt"
	"math"
	"time"

	"github.com/goccy/go-json"
)

// Duration contains metrics for duration events.
type Duration struct {
	count uint64
	total time.Duration
	min   time.Duration
	max   time.Duration
}

// NewDuration instantiates a new duration event metrics.
func NewDuration() *Duration {
	return &Duration{}
}

// Add adds a duration event to the metrics.
func (d *Duration) Add(in time.Duration) *Duration {
	d.count++
	d.total += in
	if in < d.min || d.min == 0 {
		d.min = in
	}
	if in > d.max {
		d.max = in
	}
	return d
}

// Count returns the number of reported events.
func (d *Duration) Count() uint64 {
	return d.count
}

// Total returns the total duration of all reported events.
func (d *Duration) Total() time.Duration {
	return d.total
}

// Min returns the minimum event duration.
func (d *Duration) Min() time.Duration {
	return d.min
}

// Max returns the maximum event duration.
func (d *Duration) Max() time.Duration {
	return d.max
}

// Avg returns the average event duration.
func (d *Duration) Avg() time.Duration {
	if d.count == 0 {
		return 0
	}
	return time.Duration(math.Round(float64(d.total) / float64(d.count)))
}

// Reset rests the event duration metrics.
func (d *Duration) Reset() {
	d.count = 0
	d.total = 0
	d.min = 0
	d.max = 0
}

// Merge merges the given metrics with this one.
func (d *Duration) Merge(other *Duration) {
	d.count += other.count
	d.total += other.total
	if d.min == 0 || (other.min > 0 && other.min < d.min) {
		d.min = other.min
	}
	if other.max > d.max {
		d.max = other.max
	}
}

// String implements the Stringer interface.
func (d *Duration) String() string {
	return fmt.Sprintf("count:%d total:%s min:%s max:%s avg:%s", d.count, d.total, d.min, d.max, d.Avg())
}

// MarshalJSON implements the JSON Marshaller interface.
func (d *Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonDuration{Count: d.count, Total: d.total, Min: d.min, Max: d.max, Avg: d.Avg()})
}

type jsonDuration struct {
	Count uint64        `json:"count"`
	Total time.Duration `json:"total"`
	Min   time.Duration `json:"min"`
	Max   time.Duration `json:"max"`
	Avg   time.Duration `json:"avg"`
}
