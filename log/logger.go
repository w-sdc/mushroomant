package log

import (
	"context"
	"sync"
)

// Level represents the log level.
type Level int32

// Log levels.
const (
	LLDebug Level = iota
	LLInfo
	LLWarn
	LLError
	// Fatal and Panic are internal levels. can not change by SetLevel().
	llPanic
	llFatal
)

// LevelLogger represents a leveled logger.
type LevelLogger interface {
	Debug(args ...interface{})
	// The method with suffix "xt" returns log function if the log level is
	// enabled. It is useful to avoid unnecessary memory allocation when the
	// log level is disabled. will improve performance when the log message is
	// expensive to construct.
	Debugxt() func(args ...interface{})
	Info(args ...interface{})
	Infoxt() func(args ...interface{})
	Warn(args ...interface{})
	Warnxt() func(args ...interface{})
	Error(args ...interface{})
	Errorxt() func(args ...interface{})
	// Fatal logs a message at fatal level and exits the application os.Exit(1).
	Fatal(args ...interface{})
	// Panic logs a message at panic level and panics the application.
	Panic(args ...interface{})
	// Derivatively creates a new logger with context trace.
	// An TraceID should be assigned by function WithTrace() before calling this
	// method.
	// The TraceID will be logged in the log message.
	FromTrace(ctx context.Context) (LevelLogger, error)
	// SetLevel sets the log level.
	SetLevel(level Level)
}

// provide simple logger implementation for quick use.
var quickLogger = make(map[string]LevelLogger)

// mtxQuickLogger is the mutex for quickLogger map.
var mtxQuickLogger sync.Mutex

// GetQuickLogger returns a quick logger by name.
func GetQuickLogger(name string) LevelLogger {
	mtxQuickLogger.Lock()
	defer mtxQuickLogger.Unlock()
	if l, ok := quickLogger[name]; ok {
		return l
	}
	l := NewSimpleLogger(name, "", DefaultLevelWriter.Clone())
	quickLogger[name] = l
	return l
}

// convert level to string.
func levelName(level Level) string {
	switch level {
	case LLDebug:
		return "DEBUG"
	case LLInfo:
		return "INFO"
	case LLWarn:
		return "WARNING"
	case LLError:
		return "ERROR"
	case llPanic:
		return "PANIC!"
	case llFatal:
		return "FATAL!!"
	default:
		return "UNKNOWN"
	}
}
