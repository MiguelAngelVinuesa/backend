package log

import (
	"fmt"
)

// WithField wraps the Logrus interface.
func (z *z) WithField(field string, value any) Logger {
	l := &logrus{logger: z, fields: make([]any, 0, 8)}
	l.fields = append(l.fields, field)
	l.fields = append(l.fields, value)
	return l
}

// WithFields wraps the Logrus interface.
func (z *z) WithFields(fields Fields) Logger {
	l := &logrus{logger: z, fields: make([]any, 0, 8)}
	for k, v := range fields {
		l.fields = append(l.fields, k)
		l.fields = append(l.fields, v)
	}
	return l
}

// Trace wraps the Logrus interface.
func (z *z) Trace(msg string, keysAndValues ...any) {
	z.Debug(msg, keysAndValues...)
}

// Print wraps the Logrus interface.
func (z *z) Print(msg string, keysAndValues ...any) {
	z.Info(msg, keysAndValues...)
}

// Fatal wraps the Logrus interface.
func (z *z) Fatal(msg string, keysAndValues ...any) {
	z.Panic(msg, keysAndValues...)
}

// Tracef wraps the Logrus interface.
func (z *z) Tracef(msg string, params ...any) {
	z.Debug(fmt.Sprintf(msg, params...))
}

// Debugf wraps the Logrus interface.
func (z *z) Debugf(msg string, params ...any) {
	z.Debug(fmt.Sprintf(msg, params...))
}

// Printf wraps the Logrus interface.
func (z *z) Printf(msg string, params ...any) {
	z.Info(fmt.Sprintf(msg, params...))
}

// Infof wraps the Logrus interface.
func (z *z) Infof(msg string, params ...any) {
	z.Info(fmt.Sprintf(msg, params...))
}

// Warnf wraps the Logrus interface.
func (z *z) Warnf(msg string, params ...any) {
	z.Warn(fmt.Sprintf(msg, params...))
}

// Errorf wraps the Logrus interface.
func (z *z) Errorf(msg string, params ...any) {
	z.Error(fmt.Sprintf(msg, params...))
}

// Fatalf wraps the Logrus interface.
func (z *z) Fatalf(msg string, params ...any) {
	z.Panic(fmt.Sprintf(msg, params...))
}

// Panicf wraps the Logrus interface.
func (z *z) Panicf(msg string, params ...any) {
	z.Panic(fmt.Sprintf(msg, params...))
}

// logrus is a wrapper to temporarily store key/value pairs.
type logrus struct {
	logger Logger
	fields []any
}

// SetLevel implements the Logger interface.
func (l *logrus) SetLevel(level Level) {
	l.logger.SetLevel(level)
}

// Enabled implements the Logger interface.
func (l *logrus) Enabled(level Level) bool {
	return l.logger.Enabled(level)
}

// WithField implements the Logger interface.
func (l *logrus) WithField(field string, value any) Logger {
	l.fields = append(l.fields, field)
	l.fields = append(l.fields, value)
	return l
}

// WithFields implements the Logger interface.
func (l *logrus) WithFields(fields Fields) Logger {
	for k, v := range fields {
		l.fields = append(l.fields, k)
		l.fields = append(l.fields, v)
	}
	return l
}

// Trace implements the Logger interface.
func (l *logrus) Trace(msg string, keysAndValues ...any) {
	l.logger.Debug(msg, append(l.fields, keysAndValues...))
}

// Tracef implements the Logger interface.
func (l *logrus) Tracef(msg string, params ...any) {
	l.logger.Debug(fmt.Sprintf(msg, params...), l.fields)
}

// Debug implements the Logger interface.
func (l *logrus) Debug(msg string, keysAndValues ...any) {
	l.logger.Debug(msg, append(l.fields, keysAndValues...))
}

// Debugf implements the Logger interface.
func (l *logrus) Debugf(msg string, params ...any) {
	l.logger.Debug(fmt.Sprintf(msg, params...), l.fields)
}

// Print implements the Logger interface.
func (l *logrus) Print(msg string, keysAndValues ...any) {
	l.logger.Info(msg, append(l.fields, keysAndValues...))
}

// Printf implements the Logger interface.
func (l *logrus) Printf(msg string, params ...any) {
	l.logger.Info(fmt.Sprintf(msg, params...), l.fields)
}

// Info implements the Logger interface.
func (l *logrus) Info(msg string, keysAndValues ...any) {
	l.logger.Info(msg, append(l.fields, keysAndValues...))
}

// Infof implements the Logger interface.
func (l *logrus) Infof(msg string, params ...any) {
	l.logger.Info(fmt.Sprintf(msg, params...), l.fields)
}

// Warn implements the Logger interface.
func (l *logrus) Warn(msg string, keysAndValues ...any) {
	l.logger.Warn(msg, append(l.fields, keysAndValues...))
}

// Warnf implements the Logger interface.
func (l *logrus) Warnf(msg string, params ...any) {
	l.logger.Warn(fmt.Sprintf(msg, params...), l.fields)
}

// Error implements the Logger interface.
func (l *logrus) Error(msg string, keysAndValues ...any) {
	l.logger.Error(msg, append(l.fields, keysAndValues...))
}

// Errorf implements the Logger interface.
func (l *logrus) Errorf(msg string, params ...any) {
	l.logger.Error(fmt.Sprintf(msg, params...), l.fields)
}

// Fatal implements the Logger interface.
func (l *logrus) Fatal(msg string, keysAndValues ...any) {
	l.logger.Panic(msg, append(l.fields, keysAndValues...))
}

// Fatalf implements the Logger interface.
func (l *logrus) Fatalf(msg string, params ...any) {
	l.logger.Panic(fmt.Sprintf(msg, params...), l.fields)
}

// Panic implements the Logger interface.
func (l *logrus) Panic(msg string, keysAndValues ...any) {
	l.logger.Panic(msg, append(l.fields, keysAndValues...))
}

// Panicf implements the Logger interface.
func (l *logrus) Panicf(msg string, params ...any) {
	l.logger.Panic(fmt.Sprintf(msg, params...), l.fields)
}
