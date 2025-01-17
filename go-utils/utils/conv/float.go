package conv

import (
	"strconv"

	"github.com/goccy/go-json"
)

// FloatFromAny returns the floating point value from the input if possible, otherwise the default value.
// If the default value is not supplied, zero will be used as the default value.
func FloatFromAny(in any, dflt ...float64) float64 {
	var out float64
	if len(dflt) > 0 {
		out = dflt[0]
	}

	switch t := in.(type) {
	case json.Number:
		f, _ := t.Float64()
		return f

	case float64:
		return t
	case float32:
		return float64(t)

	case int:
		return float64(t)
	case int8:
		return float64(t)
	case int16:
		return float64(t)
	case int32:
		return float64(t)
	case int64:
		return float64(t)

	case uint:
		return float64(t)
	case uint8:
		return float64(t)
	case uint16:
		return float64(t)
	case uint32:
		return float64(t)
	case uint64:
		return float64(t)

	case string:
		if f, err := strconv.ParseFloat(t, 64); err == nil {
			return f
		}
		if i, err := strconv.Atoi(t); err == nil {
			return float64(i)
		}

	case bool:
		if t {
			return 1
		}
		return 0
	}

	return out
}
