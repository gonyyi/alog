package alog

import (
	"sync"
)

// Default value for buffer is 512 bytes, and 10 KV items.
const (
	entry_buf_size = 512
	entry_kv_size  = 10
)

// newEntryPoolItem is a function that returns new Entry
// This was with newEntryPool() but separated w to make it
// inline-able.
func newEntryPoolItem() interface{} {
	return &Entry{
		buf: make([]byte, entry_buf_size),
		kvs: make([]KeyValue, entry_kv_size),
	}
}

// newEntryPool will create Entry pool.
func newEntryPool() *entryPool {
	return &entryPool{
		pool: sync.Pool{
			New: newEntryPoolItem,
		},
	}
}

// entryPool is an a Buffer implementation of sync.Pool.
type entryPool struct {
	pool sync.Pool
}

// Disable this for inlining
// Get will obtain Entry (pointer) from the pool
//func (p *entryPool) Get(f Flag, tb *bucket, pool *entryPool, w io.Writer, orfmtr Formatter) *Entry {
//	b := p.pool.Get().(*Entry)
//	b.flag, b.tbucket, b.pool, b.orFmtr, b.w = f, tb, pool, orfmtr, w
//	return b
//}

// Get will obtain Entry (pointer) from the pool
func (p *entryPool) Get(info entryInfo) *Entry {
	b := p.pool.Get().(*Entry) // cost: 69
	b.info = info
	return b
}

// Put will put Entry back to the pool
func (p *entryPool) Put(b *Entry) {
	// When buffer became too big, do not put it back.
	if cap(b.buf) > 64<<10 {
		return
	}
	b.buf = b.buf[:0]
	b.kvs = b.kvs[:0]
	p.pool.Put(b)
}
