package metrics

import (
	"time"

	"github.com/goccy/go-json"
)

// Durations contains a list of event duration metrics
type Durations struct {
	list []*Duration
}

// NewDurations instantiates a list of event duration metrics.
func NewDurations(max DurationType) *Durations {
	d := &Durations{list: make([]*Duration, max+1)}
	for ix := DurationType(0); ix <= max; ix++ {
		d.list[ix] = NewDuration()
	}
	return d
}

// Add adds an event duration of the given type to the list.
func (d *Durations) Add(t DurationType, in time.Duration) {
	d.list[t].Add(in)
}

// GetType gets the event duration metrics of the given type.
func (d *Durations) GetType(t DurationType) *Duration {
	return d.list[t]
}

// GetAll returns all event duration metrics as a slice.
func (d *Durations) GetAll() []*Duration {
	return d.list
}

// Reset resets the list.
func (d *Durations) Reset() {
	for ix := range d.list {
		d.list[ix].Reset()
	}
}

// Clone creates a deep copy of the metrics.
func (d *Durations) Clone() *Durations {
	l := len(d.list) + 1
	n := &Durations{list: make([]*Duration, l)}
	for ix, o := range d.list {
		n.list[ix] = &Duration{count: o.count, total: o.total, min: o.min, max: o.max}
	}
	return n
}

// MarshalJSON implements the JSON Marshaller interface.
func (d *Durations) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.list)
}
