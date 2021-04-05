package alog

import "sync"

const(
	entry_buf_size = 1024
	entry_kv_size = 10
)

var pool = sync.Pool {
	New: func() interface{} {
		return &Entry{
			buf: make([]byte, entry_buf_size),
			kvs: make([]KeyValue, entry_kv_size),
		}
	},
}

func get(info entryInfo) *Entry {
	b := pool.Get().(*Entry)
	b.info = info
	return b
}

func put(b *Entry) {
	b.buf = b.buf[:0]
	b.kvs = b.kvs[:0]
	pool.Put(b)
}