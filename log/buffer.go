package log

import "sync"

// refer to golang native logger implementation. here just used the same buffer
// pool implementation.

// bufPool is a global buffer pool for log formatting.
var bufPool = sync.Pool{
	New: func() interface{} {
		return new([]byte)
	},
}

// getBuf returns a buffer from the pool.
func getBuf() *[]byte {
	b := bufPool.Get().(*[]byte)
	*b = (*b)[:0]
	return b
}

// putBuf returns a buffer to the pool.
func putBuf(b *[]byte) {
	if cap(*b) > 64<<10 {
		b = nil
	}
	bufPool.Put(b)
}
