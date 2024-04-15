package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// StandardLogger provides API compatibility with standard printf loggers
// eg. go-logging
type StandardLogger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
}

// NewLogger retrieves an event logger by name
func NewLogger(module string) *Logger {
	if len(module) == 0 {
		setuplog := getLogger("setup-logger")
		setuplog.Error("Missing name parameter")
		module = "undefined"
	}

	logger := getLogger(module)
	skipLogger := logger.Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar()

	return &Logger{
		module:        module,
		SugaredLogger: *logger,
		skipLogger:    *skipLogger,
	}
}

// Logger implements the EventLogger and wraps a go-logging EventLogger
type Logger struct {
	zap.SugaredLogger
	// used to fix the caller location when calling Warning and Warningf.
	skipLogger zap.SugaredLogger
	module     string
}

// Warning is for compatibility
// Deprecated: use Warn(args ...interface{}) instead
func (logger *Logger) Warning(args ...interface{}) {
	logger.skipLogger.Warn(args...)
}

// Warningf is for compatibility
// Deprecated: use Warnf(format string, args ...interface{}) instead
func (logger *Logger) Warningf(format string, args ...interface{}) {
	logger.skipLogger.Warnf(format, args...)
}

// FormatRFC3339 returns the given time in UTC with RFC3999Nano format.
func FormatRFC3339(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}

func WithStacktrace(l *Logger, level LogLevel) *Logger {
	copyLogger := *l
	copyLogger.SugaredLogger = *copyLogger.SugaredLogger.Desugar().
		WithOptions(zap.AddStacktrace(zapcore.Level(level))).Sugar()
	copyLogger.skipLogger = *copyLogger.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar()
	return &copyLogger
}

// WithSkip returns a new logger that skips the specified number of stack frames when reporting the
// line/file.
func WithSkip(l *Logger, skip int) *Logger {
	copyLogger := *l
	copyLogger.SugaredLogger = *copyLogger.SugaredLogger.Desugar().
		WithOptions(zap.AddCallerSkip(skip)).Sugar()
	copyLogger.skipLogger = *copyLogger.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar()
	return &copyLogger
}
