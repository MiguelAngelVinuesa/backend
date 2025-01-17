package log

import (
	"strings"
)

// Level indicates the log level.
type Level uint8

const (
	// DebugLevel logs messages at debug level and higher.
	DebugLevel Level = iota + 1
	// InfoLevel logs messages at info level and higher.
	// InfoLevel logs messages at info level and higher.
	InfoLevel
	// WarnLevel logs messages at warning level and higher.
	WarnLevel
	// ErrorLevel logs messages at error level and higher.
	ErrorLevel
	// PanicLevel logs messages at panic level and higher.
	PanicLevel
	TraceLevel = DebugLevel // added for compatibility with Logrus
	FatalLevel = PanicLevel // added for compatibility with Logrus
)

// String implements the stringer interface.
func (l Level) String() string {
	switch l {
	case InfoLevel:
		return "INFO"
	case DebugLevel:
		return "DEBUG"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case PanicLevel:
		return "PANIC"
	default:
		return "INFO"
	}
}

// LevelFromName determines the appropriate log level from the input string.
func LevelFromName(name string) Level {
	switch strings.ToUpper(name) {
	case "INFO":
		return InfoLevel
	case "DEBUG":
		return DebugLevel
	case "WARN", "WARNING":
		return WarnLevel
	case "ERR", "ERROR":
		return ErrorLevel
	case "PANIC":
		return PanicLevel
	default:
		return InfoLevel
	}
}
