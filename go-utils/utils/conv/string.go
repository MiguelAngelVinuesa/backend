package conv

import (
	"strconv"

	"github.com/goccy/go-json"
)

// StringFromAny returns the string value from the input if possible, otherwise the default value.
// If the default value is not supplied, an empty string will be used as the default value.
func StringFromAny(in any, dflt ...string) string {
	var out string
	if len(dflt) > 0 {
		out = dflt[0]
	}

	switch t := in.(type) {
	case string:
		return t

	case json.Number:
		f, _ := t.Float64()
		return strconv.FormatFloat(f, 'g', -1, 64)

	case int:
		return strconv.Itoa(t)
	case int8:
		return strconv.Itoa(int(t))
	case int16:
		return strconv.Itoa(int(t))
	case int32:
		return strconv.Itoa(int(t))
	case int64:
		return strconv.Itoa(int(t))

	case uint:
		return strconv.FormatUint(uint64(t), 10)
	case uint8:
		return strconv.FormatUint(uint64(t), 10)
	case uint16:
		return strconv.FormatUint(uint64(t), 10)
	case uint32:
		return strconv.FormatUint(uint64(t), 10)
	case uint64:
		return strconv.FormatUint(t, 10)

	case float64:
		return strconv.FormatFloat(t, 'g', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(t), 'g', -1, 64)

	case bool:
		if t {
			return "1"
		}
		return "0"
	}

	return out
}

// StringsFromAny returns a slice of strings from the input if possible.
// If not possible to convert the input, the default slice will be returned if it is supplied.
// In all other cases a nil slice will be returned.
func StringsFromAny(in any, dflt ...string) []string {
	out := dflt

	switch t := in.(type) {
	case []string:
		out = t
	case string:
		return []string{t}

	case []any:
		out = make([]string, len(t))
		for ix := range t {
			out[ix] = StringFromAny(t[ix])
		}
	case []int:
		out = make([]string, len(t))
		for ix := range t {
			out[ix] = strconv.Itoa(t[ix])
		}
	case []int8:
		out = make([]string, len(t))
		for ix := range t {
			out[ix] = strconv.Itoa(int(t[ix]))
		}
	case []int16:
		out = make([]string, len(t))
		for ix := range t {
			out[ix] = strconv.Itoa(int(t[ix]))
		}
	case []int64:
		out = make([]string, len(t))
		for ix := range t {
			out[ix] = strconv.FormatInt(t[ix], 10)
		}
	case []uint:
		out = make([]string, len(t))
		for ix := range t {
			out[ix] = strconv.FormatUint(uint64(t[ix]), 10)
		}
	case []uint8:
		out = make([]string, len(t))
		for ix := range t {
			out[ix] = strconv.FormatUint(uint64(t[ix]), 10)
		}
	case []uint16:
		out = make([]string, len(t))
		for ix := range t {
			out[ix] = strconv.FormatUint(uint64(t[ix]), 10)
		}
	case []uint64:
		out = make([]string, len(t))
		for ix := range t {
			out[ix] = strconv.FormatUint(t[ix], 10)
		}
	}

	return out
}
