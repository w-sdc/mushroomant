package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"time"
)

// simpLogger provides a simple LevelLogger implementation.
type simpLogger struct {
	timefmt    string
	modulename string
	writer     LevelWriter
	traceinfo  string
}

// NewSimpleLogger returns a new simple LevelLogger.
func NewSimpleLogger(
	modulename string,
	timefmt string,
	writer LevelWriter,
) LevelLogger {
	tfmt := timefmt
	if tfmt == "" {
		tfmt = "01-02,2006 15:04:05.000"
	}
	return &simpLogger{
		timefmt:    tfmt,
		modulename: modulename,
		writer:     writer,
	}
}

// format the log message.
func (l *simpLogger) formatter(
	rww RawWriter, lv Level, cdp int, bfunc func([]byte), args ...interface{},
) {
	m := logMeta{
		ts:     time.Now(),
		lv:     lv,
		module: l.modulename,
	}
	_, file, line, ok := runtime.Caller(cdp)
	if ok {
		m.file = path.Base(file)
		m.line = line
	}

	levelname := levelName(lv)
	time := m.ts.Format(l.timefmt)

	buf := getBuf()
	defer putBuf(buf)
	*buf = fmt.Appendf(*buf, "%s % 7s %s (%s:%d)",
		time, levelname, l.modulename, m.file, m.line)
	if l.traceinfo != "" {
		*buf = fmt.Append(*buf, " [#T:", l.traceinfo, "] - ")
	} else {
		*buf = fmt.Append(*buf, " - ")
	}
	*buf = fmt.Append(*buf, args...)
	if bfunc != nil {
		bfunc(*buf)
	}
	rww.WriteItem(&m, func(w io.Writer) {
		w.Write(*buf)
	})
}

// direct log to writer.
func (l *simpLogger) tolog(level Level, args ...interface{}) {
	rww := l.writer.Write(level)
	if rww == nil {
		return
	}
	l.formatter(rww, level, 3, nil, args...)
}

// getPrinter returns log printer if the level is enabled.
func (l *simpLogger) getPrinter(level Level) func(args ...interface{}) {
	rww := l.writer.Write(level)
	if rww == nil {
		return nil
	}
	return func(args ...interface{}) {
		l.formatter(rww, level, 2, nil, args...)
	}
}

func (l *simpLogger) Debug(args ...interface{}) {
	l.tolog(LLDebug, args...)
}

func (l *simpLogger) Debugxt() func(args ...interface{}) {
	return l.getPrinter(LLDebug)
}

func (l *simpLogger) Info(args ...interface{}) {
	l.tolog(LLInfo, args...)
}

func (l *simpLogger) Infoxt() func(args ...interface{}) {
	return l.getPrinter(LLInfo)
}

func (l *simpLogger) Warn(args ...interface{}) {
	l.tolog(LLWarn, args...)
}

func (l *simpLogger) Warnxt() func(args ...interface{}) {
	return l.getPrinter(LLWarn)
}

func (l *simpLogger) Error(args ...interface{}) {
	l.tolog(LLError, args...)
}

func (l *simpLogger) Errorxt() func(args ...interface{}) {
	return l.getPrinter(LLError)
}

func (l *simpLogger) Fatal(args ...interface{}) {
	l.tolog(llFatal, args...)
	os.Exit(1)
}

func (l *simpLogger) Panic(args ...interface{}) {
	l.tolog(llPanic, args...)
	panic(fmt.Sprint(args...))
}

func (l *simpLogger) FromTrace(ctx context.Context) (LevelLogger, error) {
	t := GetTrace(ctx)
	if t == nil {
		return nil, fmt.Errorf("no trace info found")
	}
	nl := &simpLogger{
		timefmt:    l.timefmt,
		modulename: l.modulename,
		writer:     l.writer.Clone(),
		traceinfo:  t.String(),
	}
	return nl, nil
}

func (l *simpLogger) SetLevel(level Level) {
	l.writer.SetLevel(level)
}
