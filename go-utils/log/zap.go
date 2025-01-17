package log

import (
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitZAP instantiates a ZAP backed logger.
// output defines the output stream; supported values: STDOUT, STDERR; default STDOUT.
// level defines the initial logging level; supported values: DEBUG, INFO, WARNING, ERROR; default INFO.
// format defines the output format; supported values: TEXT, JSON; default JSON.
// dev sets up the logger to include stack traces.
func InitZAP(output, level, format string, dev bool) Logger {
	cfg := &zap.Config{
		Level:            zap.NewAtomicLevelAt(encodeLevel(level)),
		Development:      dev,
		Encoding:         encodeEncoding(format),
		OutputPaths:      []string{encodeOutput(output)},
		ErrorOutputPaths: []string{encodeOutput(output)},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "m",
			LevelKey:       "l",
			TimeKey:        "t",
			NameKey:        "n",
			CallerKey:      "c",
			FunctionKey:    "",
			StacktraceKey:  traceKey(dev),
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     timeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
	}

	l, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic("failed to initialize logger")
	}

	return &z{cfg: cfg, l: l.Sugar()}
}

// SetLevel dynamically changes the logging level.
func (z *z) SetLevel(lvl Level) {
	switch lvl {
	case InfoLevel:
		z.cfg.Level.SetLevel(zap.InfoLevel)
	case DebugLevel:
		z.cfg.Level.SetLevel(zap.DebugLevel)
	case WarnLevel:
		z.cfg.Level.SetLevel(zap.WarnLevel)
	case ErrorLevel:
		z.cfg.Level.SetLevel(zap.ErrorLevel)
	default:
		z.cfg.Level.SetLevel(zap.PanicLevel)
	}
}

// Enabled indicates if the given level matches the minimum level on the logger instance.
func (z *z) Enabled(lvl Level) bool {
	switch lvl {
	case InfoLevel:
		return z.cfg.Level.Enabled(zap.InfoLevel)
	case DebugLevel:
		return z.cfg.Level.Enabled(zap.DebugLevel)
	case WarnLevel:
		return z.cfg.Level.Enabled(zap.WarnLevel)
	case ErrorLevel:
		return z.cfg.Level.Enabled(zap.ErrorLevel)
	default:
		return z.cfg.Level.Enabled(zap.PanicLevel)
	}
}

// Debug implements the Logger interface.
func (z *z) Debug(msg string, keysAndValues ...interface{}) {
	z.l.Debugw(msg, keysAndValues...)
}

// Info implements the Logger interface.
func (z *z) Info(msg string, keysAndValues ...interface{}) {
	z.l.Infow(msg, keysAndValues...)
}

// Warn implements the Logger interface.
func (z *z) Warn(msg string, keysAndValues ...interface{}) {
	z.l.Warnw(msg, keysAndValues...)
}

// Error implements the Logger interface.
func (z *z) Error(msg string, keysAndValues ...interface{}) {
	z.l.Errorw(msg, keysAndValues...)
}

// Panic implements the Logger interface.
func (z *z) Panic(msg string, keysAndValues ...interface{}) {
	z.l.Panicw(msg, keysAndValues...)
}

type z struct {
	cfg *zap.Config
	l   *zap.SugaredLogger
}

func traceKey(dev bool) string {
	if dev {
		return "s"
	}
	return ""
}

func encodeEncoding(format string) string {
	if strings.EqualFold(format, "TEXT") {
		return "console"
	}
	return "json"
}

func encodeOutput(output string) string {
	if strings.EqualFold(output, "STDERR") {
		return "stderr"
	}
	return "stdout"
}

func encodeLevel(level string) zapcore.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return zap.DebugLevel
	case "WARN", "WARNING":
		return zap.WarnLevel
	case "ERR", "ERROR":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("20060102T150405.999999"))
}
