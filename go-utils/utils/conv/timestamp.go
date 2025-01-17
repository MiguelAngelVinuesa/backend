package conv

import (
	"time"
)

func TimestampFromAny(in any, dflt ...time.Time) time.Time {
	switch t := in.(type) {
	case time.Time:
		return t

	case *time.Time:
		if t != nil {
			return *t
		}

	case string:
		for ix := range formats {
			if s, err := time.Parse(formats[ix], t); err == nil {
				return s
			}
		}

	case int64:
		return time.UnixMilli(t)

	case int:
		return time.UnixMilli(int64(t))
	}

	if len(dflt) > 0 {
		return dflt[0]
	}
	return time.Time{}
}

var formats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"20060102",
	"2006-01-02",
	"20060102T150405",
	"2006-01-02T15:04:05",
}
