package utils

import (
	"fmt"
	"strconv"

	"github.com/goccy/go-json"
)

// StringFromAny converts the input to a string value.
func StringFromAny(in any) string {
	switch t := in.(type) {
	case string:
		return t
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case float64:
		return strconv.FormatFloat(t, 'g', -1, 64)
	case json.Number:
		return t.String()
	case bool:
		if t {
			return "1"
		}
		return "0"
	}

	if f, ok := in.(fmt.Stringer); ok {
		return f.String()
	}
	return fmt.Sprintf("%v", in)
}
