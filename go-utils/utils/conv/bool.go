package conv

import (
	"strings"

	"github.com/goccy/go-json"
)

// BoolFromAny returns the boolean value from the input if possible, otherwise the default value.
// If the default value is not supplied, false will be used as the default value.
func BoolFromAny(in any, dflt ...bool) bool {
	var out bool
	if len(dflt) > 0 {
		out = dflt[0]
	}

	switch t := in.(type) {
	case bool:
		return t

	case string:
		t = "|" + t + "|"
		if out {
			if strings.Contains(falseStrings, t) {
				return false
			}
		} else {
			if !out && strings.Contains(trueStrings, t) {
				return true
			}
		}

	case json.Number:
		f, _ := t.Float64()
		return f != 0.0

	case int:
		return t != 0
	case int8:
		return t != 0
	case int16:
		return t != 0
	case int32:
		return t != 0
	case int64:
		return t != 0

	case uint:
		return t != 0
	case uint8:
		return t != 0
	case uint16:
		return t != 0
	case uint32:
		return t != 0
	case uint64:
		return t != 0

	case float64:
		return t != 0.0
	case float32:
		return t != 0.0
	}

	return out
}

const (
	falseStrings = "|false|f|.f.|0|no|n|.n."
	trueStrings  = "|true|t|.t.|1|yes|y|.y."
)
