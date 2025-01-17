package time

import (
	"time"
)

// Floats is the generic type for floating point values.
type Floats interface {
	~float64 | ~float32
}

// FV represents a time-series floating point value.
type FV[V Floats] struct {
	unixMilli int64
	value     V
}

// NewFV instantiates a new time-series floating point value.
func NewFV[V Floats](t time.Time, v V) FV[V] {
	return FV[V]{unixMilli: t.UnixMilli(), value: v}
}

// Time returns the time value of a time-series floating point value.
func (fv FV[V]) Time() time.Time {
	return time.UnixMilli(fv.unixMilli)
}

// Value returns the value of a time-series floating point value.
func (fv FV[V]) Value() V {
	return fv.value
}

// FVs is a slice of time-series float point values.
type FVs[V Floats] []FV[V]

// MinMaxTime determines the minimum & maximum time in the time-series.
func (v FVs[V]) MinMaxTime() (time.Time, time.Time) {
	if len(v) == 0 {
		return time.Time{}, time.Time{}
	}

	min, max := v[0].unixMilli, v[0].unixMilli
	for ix := 1; ix < len(v); ix++ {
		if v[ix].unixMilli < min {
			min = v[ix].unixMilli
		}
		if v[ix].unixMilli > max {
			max = v[ix].unixMilli
		}
	}
	return time.UnixMilli(min), time.UnixMilli(max)
}

// MinMaxValue determines the minimum & maximum value in the time-series.
func (v FVs[V]) MinMaxValue() (V, V) {
	if len(v) == 0 {
		return 0, 0
	}

	min, max := v[0].value, v[0].value
	for ix := 1; ix < len(v); ix++ {
		if v[ix].value < min {
			min = v[ix].value
		}
		if v[ix].value > max {
			max = v[ix].value
		}
	}
	return min, max
}
