package metrics

import (
	"math"

	"github.com/goccy/go-json"
)

// Size contains metrics for object sizes.
type Size struct {
	count uint64
	total uint64
	min   uint64
	max   uint64
}

// NewSize instantiates a new object size metrics.
func NewSize() *Size {
	return &Size{}
}

// Add adds an object size to the metrics.
func (s *Size) Add(in uint64) *Size {
	s.count++
	s.total += in
	if in < s.min || s.min == 0 {
		s.min = in
	}
	if in > s.max {
		s.max = in
	}
	return s
}

// Count returns the number of reported objects.
func (s *Size) Count() uint64 {
	return s.count
}

// Total returns the total size of all reported objects.
func (s *Size) Total() uint64 {
	return s.total
}

// Min returns the minimum object size.
func (s *Size) Min() uint64 {
	return s.min
}

// Max returns the maximum object size.
func (s *Size) Max() uint64 {
	return s.max
}

// Avg returns the average object size.
func (s *Size) Avg() uint64 {
	if s.count == 0 {
		return 0
	}
	return uint64(math.Round(float64(s.total) / float64(s.count)))
}

// Reset rests the object size metrics.
func (s *Size) Reset() {
	s.count = 0
	s.total = 0
	s.min = 0
	s.max = 0
}

// MarshalJSON implements the JSON Marshaller interface.
func (s *Size) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonSize{Count: s.count, Total: s.total, Min: s.min, Max: s.max, Avg: s.Avg()})
}

type jsonSize struct {
	Count uint64 `json:"count"`
	Total uint64 `json:"total"`
	Min   uint64 `json:"min"`
	Max   uint64 `json:"max"`
	Avg   uint64 `json:"avg"`
}
