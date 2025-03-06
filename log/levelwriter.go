package log

import "sync/atomic"

// LevelWriter provides a way to write log data with different levels.
type LevelWriter interface {
	Write(level Level) RawWriter
	SetLevel(level Level)
	Clone() LevelWriter
}

// levelWriter is the default implementation of LevelWriter.
type levelWriter struct {
	debugWriter   RawWriter
	infoWriter    RawWriter
	warningWriter RawWriter
	errorWriter   RawWriter
	fatalWriter   RawWriter
	panicWriter   RawWriter
	level         Level
}

// DefaultLevelWriter is the default LevelWriter.
var DefaultLevelWriter = NewLevelWriter(
	StdOutWriter, StdOutWriter, StdOutWriter,
	StdErrWriter, StdErrWriter, StdErrWriter,
	LLInfo,
)

// NewLevelWriter creates a new LevelWriter with the given RawWriters.
func NewLevelWriter(
	dbg, ifo, warn, err, ftl, pnc RawWriter,
	defaultLevel Level,
) LevelWriter {
	return &levelWriter{
		debugWriter:   dbg,
		infoWriter:    ifo,
		warningWriter: warn,
		errorWriter:   err,
		fatalWriter:   ftl,
		panicWriter:   pnc,
		level:         defaultLevel,
	}
}

// Write writes the log data with the given level.
func (w *levelWriter) Write(level Level) RawWriter {
	var writer RawWriter
	switch level {
	case LLDebug:
		if w.level <= LLDebug {
			writer = w.debugWriter
		}
	case LLInfo:
		if w.level <= LLInfo {
			writer = w.infoWriter
		}
	case LLWarn:
		if w.level <= LLWarn {
			writer = w.warningWriter
		}
	case LLError:
		if w.level <= LLError {
			writer = w.errorWriter
		}
	case llFatal:
		writer = w.fatalWriter
	case llPanic:
		writer = w.panicWriter
	}
	return writer
}

func checkLevel(level Level) Level {
	if level < LLDebug {
		return LLDebug
	} else if level > LLError {
		return LLError
	}
	return level
}

// SetLevel sets the log level.
func (w *levelWriter) SetLevel(level Level) {
	atomic.StoreInt32((*int32)(&w.level), int32(checkLevel(level)))
}

// Clone clones the LevelWriter with the given level.
func (w *levelWriter) Clone() LevelWriter {
	return &levelWriter{
		debugWriter:   w.debugWriter,
		infoWriter:    w.infoWriter,
		warningWriter: w.warningWriter,
		errorWriter:   w.errorWriter,
		fatalWriter:   w.fatalWriter,
		panicWriter:   w.panicWriter,
		level:         w.level,
	}
}
