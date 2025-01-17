package time

import (
	"time"
)

// Integers is the generic type for integer values.
type Integers interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// IV represents a time-series integer value.
type IV[V Integers] struct {
	unixMilli int64
	value     V
}

// NewIV instantiates a new time-series integer value.
func NewIV[V Integers](t time.Time, v V) IV[V] {
	return IV[V]{unixMilli: t.UnixMilli(), value: v}
}

// Time returns the time value of a time-series integer value.
func (iv IV[V]) Time() time.Time {
	return time.UnixMilli(iv.unixMilli)
}

// Value returns the value of a time-series integer value.
func (iv IV[V]) Value() V {
	return iv.value
}

// IVs is a slice of time-series integer values.
type IVs[V Integers] []IV[V]

// MinMaxTime determines the minimum & maximum time in the time-series.
func (v IVs[V]) MinMaxTime() (time.Time, time.Time) {
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
func (v IVs[V]) MinMaxValue() (V, V) {
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
