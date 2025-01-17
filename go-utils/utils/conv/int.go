package conv

import (
	"math"
	"strconv"

	"github.com/goccy/go-json"
)

// IntFromAny returns the integer value from the input if possible, otherwise the default value.
// If the default value is not supplied, zero will be used as the default value.
func IntFromAny(in any, dflt ...int) int {
	var out int
	if len(dflt) > 0 {
		out = dflt[0]
	}

	switch t := in.(type) {
	case int:
		return t
	case int8:
		return int(t)
	case int16:
		return int(t)
	case int32:
		return int(t)
	case int64:
		return int(t)

	case json.Number:
		f, _ := t.Float64()
		return int(f)

	case float64:
		return int(t)
	case float32:
		return int(t)

	case uint:
		return int(t)
	case uint8:
		return int(t)
	case uint16:
		return int(t)
	case uint32:
		return int(t)
	case uint64:
		if t < math.MaxInt {
			return int(t)
		}

	case string:
		if i, err := strconv.Atoi(t); err == nil {
			return i
		}
		if f, err := strconv.ParseFloat(t, 64); err == nil {
			return int(f)
		}

	case bool:
		if t {
			return 1
		}
		return 0
	}

	return out
}

// IntsFromAny returns a slice of integers from the input if possible.
// If not possible to convert the input, the default slice will be returned if it is supplied.
// In all other cases a nil slice will be returned.
func IntsFromAny(in any, dflt ...int) []int {
	out := dflt

	switch t := in.(type) {
	case []int:
		out = t
	case int:
		return []int{t}

	case []string:
		out = make([]int, len(t))
		for ix := range t {
			out[ix] = IntFromAny(t[ix])
		}
	case []any:
		out = make([]int, len(t))
		for ix := range t {
			out[ix] = IntFromAny(t[ix])
		}
	case []int8:
		out = make([]int, len(t))
		for ix := range t {
			out[ix] = int(t[ix])
		}
	case []int16:
		out = make([]int, len(t))
		for ix := range t {
			out[ix] = int(t[ix])
		}
	case []int64:
		out = make([]int, len(t))
		for ix := range t {
			out[ix] = int(t[ix])
		}
	case []uint:
		out = make([]int, len(t))
		for ix := range t {
			out[ix] = int(t[ix])
		}
	case []uint8:
		out = make([]int, len(t))
		for ix := range t {
			out[ix] = int(t[ix])
		}
	case []uint16:
		out = make([]int, len(t))
		for ix := range t {
			out[ix] = int(t[ix])
		}
	case []uint64:
		out = make([]int, len(t))
		for ix := range t {
			out[ix] = int(t[ix])
		}
	}

	return out
}
