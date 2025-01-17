package conv

import (
	"strconv"

	"github.com/goccy/go-json"
)

// UintFromAny returns the unsigned integer value from the input if possible, otherwise the default value.
// If the default value is not supplied, zero will be used as the default value.
func UintFromAny(in any, dflt ...uint) uint {
	var out uint
	if len(dflt) > 0 {
		out = dflt[0]
	}

	switch t := in.(type) {
	case uint:
		return t
	case uint8:
		return uint(t)
	case uint16:
		return uint(t)
	case uint32:
		return uint(t)
	case uint64:
		return uint(t)

	case json.Number:
		if f, _ := t.Float64(); f >= 0.0 {
			return uint(f)
		}

	case float64:
		if t >= 0.0 {
			return uint(t)
		}
	case float32:
		if t >= 0.0 {
			return uint(t)
		}

	case int:
		if t >= 0.0 {
			return uint(t)
		}
	case int8:
		if t >= 0.0 {
			return uint(t)
		}
	case int16:
		if t >= 0.0 {
			return uint(t)
		}
	case int32:
		if t >= 0.0 {
			return uint(t)
		}
	case int64:
		if t >= 0.0 {
			return uint(t)
		}

	case string:
		if i, err := strconv.ParseUint(t, 10, 64); err == nil {
			return uint(i)
		}
		if f, err := strconv.ParseFloat(t, 64); err == nil && f >= 0.0 {
			return uint(f)
		}

	case bool:
		if t {
			return 1
		}
		return 0
	}

	return out
}
