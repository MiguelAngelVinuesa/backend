package log

// Logger is a simplified interface to wrap a backend logging package.
type Logger interface {
	SetLevel(lvl Level)
	Enabled(lvl Level) bool

	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
	Panic(msg string, keysAndValues ...any)

	// Logrus wrappers.
	WithField(field string, value any) Logger
	WithFields(fields Fields) Logger

	Trace(msg string, keysAndValues ...any) // alias for Debug().
	Print(msg string, keysAndValues ...any) // alias for Info().
	Fatal(msg string, keysAndValues ...any) // alias for Panic().

	Tracef(msg string, params ...any) // DEPRECATED: use one of the field mechanisms.
	Debugf(msg string, params ...any) // DEPRECATED: use one of the field mechanisms.
	Printf(msg string, params ...any) // DEPRECATED: use one of the field mechanisms.
	Infof(msg string, params ...any)  // DEPRECATED: use one of the field mechanisms.
	Warnf(msg string, params ...any)  // DEPRECATED: use one of the field mechanisms.
	Fatalf(msg string, params ...any) // DEPRECATED: use one of the field mechanisms.
	Errorf(msg string, params ...any) // DEPRECATED: use one of the field mechanisms.
	Panicf(msg string, params ...any) // DEPRECATED: use one of the field mechanisms.
}

// Fields is a convenience type for a map of key/value pairs.
type Fields map[string]any

// BasicLogger is a minimal interface for other modules within this library.
type BasicLogger interface {
	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
}
