package metrics

import (
	"github.com/goccy/go-json"
)

// Sizes contains a list of object size metrics
type Sizes struct {
	list []*Size
}

// NewSizes instantiates a list of object size metrics.
func NewSizes(max SizeType) *Sizes {
	d := &Sizes{list: make([]*Size, max+1)}
	for ix := SizeType(0); ix <= max; ix++ {
		d.list[ix] = NewSize()
	}
	return d
}

// Add adds an object size of the given type to the list.
func (s *Sizes) Add(t SizeType, in uint64) {
	s.list[t].Add(in)
}

// GetType gets the object size metrics of the given type.
func (s *Sizes) GetType(t SizeType) *Size {
	return s.list[t]
}

// GetAll returns all object size metrics as a slice.
func (s *Sizes) GetAll() []*Size {
	return s.list
}

// Reset resets the list.
func (s *Sizes) Reset() {
	for ix := range s.list {
		s.list[ix].Reset()
	}
}

// Clone creates a deep copy of the metrics.
func (s *Sizes) Clone() *Sizes {
	l := len(s.list) + 1
	n := &Sizes{list: make([]*Size, l)}
	for ix, o := range s.list {
		n.list[ix] = &Size{count: o.count, total: o.total, min: o.min, max: o.max}
	}
	return n
}

// MarshalJSON implements the JSON Marshaller interface.
func (s *Sizes) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.list)
}
