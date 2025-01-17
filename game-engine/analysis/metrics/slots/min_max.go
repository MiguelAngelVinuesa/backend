package slots

import (
	"fmt"
	"math"
	"reflect"

	"github.com/goccy/go-json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// AcquireMinMaxInt64 instantiates a new total/min/max/counts int64 metric from the memory pool.
func AcquireMinMaxInt64() *MinMaxInt64 {
	return minMaxIntPool.Acquire().(*MinMaxInt64)
}

// MinMaxInt64 contains the total count, min/max value and a map of value counts for an int64 metric.
type MinMaxInt64 struct {
	first  bool
	Count  uint64           `json:"count"`
	Total  int64            `json:"total"`
	Min    int64            `json:"min,omitempty"`
	Max    int64            `json:"max,omitempty"`
	Keys   []int64          `json:"-"`
	Counts map[int64]uint64 `json:"counts,omitempty"`
	pool.Object
}

var minMaxIntPool = pool.NewProducer(func() (pool.Objecter, func()) {
	m := &MinMaxInt64{
		first:  true,
		Keys:   make([]int64, 0, minMaxSize),
		Counts: make(map[int64]uint64, minMaxSize),
	}
	return m, m.reset
})

// reset cleart the metric.
func (m *MinMaxInt64) reset() {
	if m != nil {
		m.ResetData()
	}
}

// ResetData implements the Objecter interface.
func (m *MinMaxInt64) ResetData() {
	m.first = true
	m.Count = 0
	m.Total = 0
	m.Min = 0
	m.Max = 0
	m.Keys = m.Keys[:0]
	clear(m.Counts)
}

// Merge merges the given metrics.
func (m *MinMaxInt64) Merge(other *MinMaxInt64) {
	m.Count += other.Count
	m.Total += other.Total

	if m.first {
		m.Min = other.Min
		m.Max = other.Max
		m.first = other.first
	} else {
		if other.Min > 0 && other.Min < m.Min {
			m.Min = other.Min
		}
		if other.Max > m.Max {
			m.Max = other.Max
		}
	}

	for _, k := range other.Keys {
		if c, ok := m.Counts[k]; ok {
			m.Counts[k] = c + other.Counts[k]
		} else {
			m.Counts[k] = other.Counts[k]
			m.Keys = append(m.Keys, k)
		}
	}
}

// Increase updates a total/min/max/counts int64 metric with the given input.
func (m *MinMaxInt64) Increase(in int64) {
	m.Count++
	m.Total += in

	if m.first {
		m.Min = in
		m.Max = in
		m.first = false
	} else {
		if in > 0 && (m.Min == 0 || in < m.Min) {
			m.Min = in
		}
		if in > m.Max {
			m.Max = in
		}
	}

	if c, ok := m.Counts[in]; ok {
		m.Counts[in] = c + 1
	} else {
		m.Counts[in] = 1
		m.Keys = append(m.Keys, in)
	}
}

// Equals is used internally for unit tests!
func (m *MinMaxInt64) Equals(other *MinMaxInt64) bool {
	return m.first == other.first &&
		m.Count == other.Count &&
		m.Total == other.Total &&
		m.Min == other.Min &&
		m.Max == other.Max &&
		reflect.DeepEqual(m.Counts, other.Counts)
}

// AcquireMinMaxUInt64 instantiates a new total/min/max/counts uint64 metric from the memory pool.
func AcquireMinMaxUInt64() *MinMaxUInt64 {
	return minMaxUIntPool.Acquire().(*MinMaxUInt64)
}

// MinMaxUInt64 contains the total count, min/max value and a map of value counts for an uint64 metric.
type MinMaxUInt64 struct {
	first  bool
	Count  uint64            `json:"count"`
	Total  uint64            `json:"total"`
	Min    uint64            `json:"min,omitempty"`
	Max    uint64            `json:"max,omitempty"`
	Keys   []uint64          `json:"-"`
	Counts map[uint64]uint64 `json:"counts,omitempty"`
	pool.Object
}

var minMaxUIntPool = pool.NewProducer(func() (pool.Objecter, func()) {
	m := &MinMaxUInt64{
		first:  true,
		Keys:   make([]uint64, 0, minMaxSize),
		Counts: make(map[uint64]uint64, minMaxSize),
	}
	return m, m.reset
})

// reset clears the metric.
func (m *MinMaxUInt64) reset() {
	if m != nil {
		m.ResetData()
	}
}

// ResetData implements the Objecter interface.
func (m *MinMaxUInt64) ResetData() {
	m.first = true
	m.Count = 0
	m.Total = 0
	m.Min = 0
	m.Max = 0
	m.Keys = m.Keys[:0]
	clear(m.Counts)
}

// Merge merges the given metrics.
func (m *MinMaxUInt64) Merge(other *MinMaxUInt64) {
	m.Count += other.Count
	m.Total += other.Total

	if m.first {
		m.Min = other.Min
		m.Max = other.Max
		m.first = other.first
	} else {
		if other.Min > 0 && other.Min < m.Min {
			m.Min = other.Min
		}
		if other.Max > m.Max {
			m.Max = other.Max
		}
	}

	for _, k := range other.Keys {
		if c, ok := m.Counts[k]; ok {
			m.Counts[k] = c + other.Counts[k]
		} else {
			m.Counts[k] = other.Counts[k]
			m.Keys = append(m.Keys, k)
		}
	}
}

// First indicates if the uint64 metric is in its initial state.
func (m *MinMaxUInt64) First() bool {
	return m.first
}

// Increase updates a total/min/max/counts uint64 metric with the given input.
func (m *MinMaxUInt64) Increase(in uint64) {
	m.Count++
	m.Total += in

	if m.first {
		m.Min = in
		m.Max = in
		m.first = false
	} else {
		if in > 0 && (m.Min == 0 || in < m.Min) {
			m.Min = in
		}
		if in > m.Max {
			m.Max = in
		}
	}

	if c, ok := m.Counts[in]; ok {
		m.Counts[in] = c + 1
	} else {
		m.Counts[in] = 1
		m.Keys = append(m.Keys, in)
	}
}

// IncreaseOne increases the metric by one.
// It is a special use case for keeping track of a counter, like "#spins above start balance".
func (m *MinMaxUInt64) IncreaseOne() {
	if m.first {
		m.Count = 1
		m.first = false
	} else {
		delete(m.Counts, m.Total)
	}

	m.Total++
	m.Min++
	m.Max++

	r := m.Total
	if _, ok := m.Counts[r]; !ok {
		m.Counts[r] = 1
		m.Keys = append(m.Keys, r)
	}
}

// Equals is used internally for unit tests!
func (m *MinMaxUInt64) Equals(other *MinMaxUInt64) bool {
	return m.first == other.first &&
		m.Count == other.Count &&
		m.Total == other.Total &&
		m.Min == other.Min &&
		m.Max == other.Max &&
		reflect.DeepEqual(m.Counts, other.Counts)
}

// AcquireMinMaxFloat64 instantiates a new total/min/max/counts float64 metric from the memory pool.
func AcquireMinMaxFloat64(decimals int) *MinMaxFloat64 {
	m := minMaxFloatPool.Acquire().(*MinMaxFloat64)
	m.Decimals = decimals
	m.Factor = math.Pow10(decimals)
	return m
}

// MinMaxFloat64 contains the total count, min/max value and a map of value counts for a float64 metric.
type MinMaxFloat64 struct {
	first    bool
	Decimals int              `json:"-"`
	Count    uint64           `json:"count"`
	Total    float64          `json:"total"`
	Min      float64          `json:"min,omitempty"`
	Max      float64          `json:"max,omitempty"`
	Factor   float64          `json:"-"`
	Keys     []int64          `json:"-"`
	Counts   map[int64]uint64 `json:"counts,omitempty"`
	pool.Object
}

var minMaxFloatPool = pool.NewProducer(func() (pool.Objecter, func()) {
	m := &MinMaxFloat64{
		first:  true,
		Keys:   make([]int64, 0, minMaxSize),
		Counts: make(map[int64]uint64, minMaxSize),
	}
	return m, m.reset
})

func (m *MinMaxFloat64) reset() {
	if m != nil {
		m.ResetData()
		m.Decimals = 0
		m.Factor = 0
	}
}

// ResetData implements the Objecter interface.
func (m *MinMaxFloat64) ResetData() {
	m.first = true
	m.Count = 0
	m.Total = 0
	m.Min = 0.0
	m.Max = 0.0
	m.Keys = m.Keys[:0]
	clear(m.Counts)
}

// Merge merges the given metrics.
func (m *MinMaxFloat64) Merge(other *MinMaxFloat64) {
	m.Count += other.Count
	m.Total += other.Total

	if m.first {
		m.Min = other.Min
		m.Max = other.Max
		m.first = other.first
	} else {
		if other.Min > 0 && other.Min < m.Min {
			m.Min = other.Min
		}
		if other.Max > m.Max {
			m.Max = other.Max
		}
	}

	for _, k := range other.Keys {
		if c, ok := m.Counts[k]; ok {
			m.Counts[k] = c + other.Counts[k]
		} else {
			m.Counts[k] = other.Counts[k]
			m.Keys = append(m.Keys, k)
		}
	}
}

// Increase updates a total/min/max/counts float64 metric with the given input.
func (m *MinMaxFloat64) Increase(in float64) {
	m.Count++
	m.Total += in

	if m.first {
		m.Min = in
		m.Max = in
		m.first = false
	} else {
		if in > 0.0 && (m.Min == 0.0 || in < m.Min) {
			m.Min = in
		}
		if in > m.Max {
			m.Max = in
		}
	}

	r := int64(math.Round(m.Factor * in))
	if c, ok := m.Counts[r]; ok {
		m.Counts[r] = c + 1
	} else {
		m.Counts[r] = 1
		m.Keys = append(m.Keys, r)
	}
}

func (m *MinMaxFloat64) MarshalJSON() ([]byte, error) {
	if m == nil {
		return json.Marshal(nil)
	}

	j := &jsonMinMaxFloat64{
		Count:  m.Count,
		Total:  m.Total,
		Min:    m.Min,
		Max:    m.Max,
		Counts: make(map[string]uint64, len(m.Counts)),
	}

	for k, v := range m.Counts {
		j.Counts[fmt.Sprintf("%.2f", float64(k)/m.Factor)] = v
	}

	return json.Marshal(j)
}

type jsonMinMaxFloat64 struct {
	Count  uint64            `json:"count"`
	Total  float64           `json:"total"`
	Min    float64           `json:"min,omitempty"`
	Max    float64           `json:"max,omitempty"`
	Counts map[string]uint64 `json:"counts,omitempty"`
	pool.Object
}

// Equals is used internally for unit tests!
func (m *MinMaxFloat64) Equals(other *MinMaxFloat64) bool {
	if m.first != other.first || m.Count != other.Count ||
		!m.equalFloat(m.Total, other.Total) ||
		!m.equalFloat(m.Min, other.Min) ||
		!m.equalFloat(m.Max, other.Max) {
		return false
	}

	for k1, v1 := range m.Counts {
		var found bool
		for k2, v2 := range other.Counts {
			if k1 == k2 {
				found = true
				if v1 != v2 {
					return false
				}
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// used internally for unit tests!
func (m *MinMaxFloat64) equalFloat(x, y float64) bool {
	return int64(math.Round(x*m.Factor)) == int64(math.Round(y*m.Factor))
}

const minMaxSize = 32
