package alog

import (
	"sync"
)

// Default value for buffer is 512 bytes, and 10 KV items.
const (
	entry_buf_size = 512
	entry_kv_size  = 10
)

// newEntryPoolItem is a function that returns new entry
// This was with newEntryPool() but separated w to make it
// inline-able.
func newEntryPoolItem() interface{} {
	return &entry{
		buf: make([]byte, entry_buf_size),
		kvs: make([]KeyValue, entry_kv_size),
	}
}

// newEntryPool will create entry pool.
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
// Get will obtain entry (pointer) from the pool
//func (p *entryPool) Get(f Flag, tb *TagBucket, pool *entryPool, w io.Writer, orfmtr Formatter) *entry {
//	b := p.pool.Get().(*entry)
//	b.flag, b.tbucket, b.pool, b.orFmtr, b.w = f, tb, pool, orfmtr, w
//	return b
//}

// Get will obtain entry (pointer) from the pool
func (p *entryPool) Get(info entryInfo) *entry {
	b := p.pool.Get().(*entry) // cost: 69
	b.info = info
	return b
}

// Put will put entry back to the pool
func (p *entryPool) Put(b *entry) {
	b.buf = b.buf[:0]
	b.kvs = b.kvs[:0]
	p.pool.Put(b)
}