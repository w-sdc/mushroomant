package log

import (
	"errors"
	"io"
	"os"
	"sync"
	"time"
)

const (
	maxLogQueueSize = 50
)

var (
	ErrClosedWriter = errors.New("closed writer")
)

// logMeta represents the metadata of a log record.
type logMeta struct {
	ts     time.Time // timestamp
	lv     Level     // log level
	file   string    // source code file
	line   int       // source code line
	module string    // module name or logger name
}

// RawWriter represents a writer that writes raw log data.
type RawWriter interface {
	// WriteItem writes a single record to the writer.
	//
	// Upon the function's return, an internal EOR flag is triggered.
	// This flag is used to separate log records, replacing the traditional line
	// break, thereby allowing multiple lines within a single record.
	//
	// In such as database operations, each record can be stored as a separate
	// row.
	WriteItem(*logMeta, func(w io.Writer)) error
}

// simpWriter implements a very simple RawWriter.
// It writes the log data to the given io.Writer. It assume the IO never be
// close, also, it does not support any additional features like log rotation,
// limitaion, etc.
type simpWriter struct {
	mtx sync.Mutex
	w   io.Writer
}

// stdout and stderr writer is created by default.
var (
	StdOutWriter = NewSimpWriter(os.Stdout)
	StdErrWriter = NewSimpWriter(os.Stderr)
)

// NewSimpFileWriter creates a new file writer with simple implementation.
func NewSimpWriter(w io.Writer) RawWriter {
	return &simpWriter{
		w: w,
	}
}

func (w *simpWriter) WriteItem(m *logMeta, wproc func(w io.Writer)) error {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	wproc(w.w)
	w.w.Write([]byte{'\n'})
	return nil
}
