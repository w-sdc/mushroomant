package log

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

// TraceID represents a trace ID.
type TraceID [32]byte

// TraceValue represents a trace value.
type TraceValue struct {
	ID   TraceID
	Name string
}

// TraceScope represents a trace chain. which includes a sequence of trace
// values. the last trace value is the current trace value, and the previous
// trace values are the parent traces.
type TraceScope []TraceValue

// type of trace key.
type traceKeyT string

// key to store trace value in context.
const traceKey traceKeyT = "mra-log-trace"

// traceCtx represents a context with trace.
type traceCtx struct {
	context.Context
	TraceValue
}

// unique trace prefix
var traceInit []byte

// trace sequence
var traceSeq uint64

func init() {
	executable, _ := os.Executable()
	// processname(pid)@user.hostname-timestamp
	initstr := fmt.Sprintf("%s(%d)@%s.%s-%d", executable, os.Getpid(),
		os.Getenv("USER"), os.Getenv("HOSTNAME"), time.Now().UnixNano())
	hash := sha256.New()
	hash.Write([]byte(initstr))
	traceInit = hash.Sum(nil)[:16]
	traceSeq = rand.Uint64() & 0x7fffffffffffffff
}

// nextTraceID generates a new trace ID.
func nextTraceID() TraceID {
	var id [32]byte
	s := atomic.AddUint64(&traceSeq, 1)
	copy(id[8:], traceInit)
	id[24] = byte(s >> 56)
	id[25] = byte(s >> 48)
	id[26] = byte(s >> 40)
	id[27] = byte(s >> 32)
	id[28] = byte(s >> 24)
	id[29] = byte(s >> 16)
	id[30] = byte(s >> 8)
	id[31] = byte(s)
	hash := md5.New()
	hash.Write(id[8:])
	copy(id[:8], hash.Sum(nil))
	return TraceID(id)
}

// WithTrace assigns a trace ID to the context. when applying this context to
// the logger, the trace ID will be logged in the log message.
// The trace ID will be generated automatically.
// User can specify a name for the trace, which will be logged in the log
// message as well.
// If provided context have nested trace, they will be logged as a parent-child
// relationship.
func WithTrace(ctx context.Context, name string) context.Context {
	return &traceCtx{
		Context:    ctx,
		TraceValue: TraceValue{ID: nextTraceID(), Name: name},
	}
}

// GetTrace returns the trace value from the context.
func GetTrace(ctx context.Context) TraceScope {
	if ctx == nil {
		return nil
	}
	rst := ctx.Value(traceKey)
	if rst == nil {
		return nil
	}
	return TraceScope(rst.([]TraceValue))
}

// Value reimplements the Value method of context.Context.
// If the key is TraceKey, it will return a slice of TraceValue which contains
// whole trace chain.
func (t *traceCtx) Value(key interface{}) interface{} {
	if key == traceKey {
		parent := t.Context.Value(key)
		if parent == nil {
			return []TraceValue{t.TraceValue}
		}
		return append(parent.([]TraceValue), t.TraceValue)
	}
	return t.Context.Value(key)
}

// String returns the string representation of the trace ID.
func (t TraceID) String() string {
	return fmt.Sprintf("%016x-%032x-%016x", t[:8], t[8:24], t[24:])
}

// Sum returns the first 8 bytes summary of the trace ID.
func (t TraceID) Sum() string {
	return fmt.Sprintf("%016x", t[:8])
}

// Hex returns the hex string of the trace ID.
func (t TraceID) Hex() string {
	return fmt.Sprintf("%064x", t[:])
}

// String returns the string representation of the trace chain.
func (t TraceScope) String() string {
	if len(t) == 0 {
		return ""
	}
	sb := strings.Builder{}
	for i, v := range t {
		if i > 0 {
			sb.WriteString("/")
		}
		sb.WriteString(v.ID.Sum())
	}
	if t[len(t)-1].Name != "" {
		sb.WriteString("#")
		sb.WriteString(t[len(t)-1].Name)
	}
	return sb.String()
}

// Format returns the formatted string representation of the trace chain.
func (t TraceScope) Format(prefix, suffix, sep string) string {
	if len(t) == 0 {
		return ""
	}
	sb := strings.Builder{}
	for i, v := range t {
		if i > 0 && sep != "" {
			sb.WriteString(sep)
		}
		if prefix != "" {
			sb.WriteString(prefix)
		}
		sb.WriteString(v.ID.String())
		sb.WriteString("#")
		sb.WriteString(v.Name)
		if suffix != "" {
			sb.WriteString(suffix)
		}
	}
	return sb.String()
}
