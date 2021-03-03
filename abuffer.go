package alog

import "sync"

// entry is the main entry format used in Alog.
// This can be used as standalone or within the sync.Pool
type entry struct {
	buf    []byte
	level  Level
	tag    Tag
	logger *Logger
	kvs    []KeyVal
}

func newAbuffer(size int) abuffer {
	return abuffer{
		pool: sync.Pool{
			New: func() interface{} {
				return &entry{
					buf: make([]byte, 512),
					kvs: make([]KeyVal, 10),
				}
			},
		},
	}
}

// buf is an a Buffer implementation of sync.Pool.
type abuffer struct {
	pool   sync.Pool
}

func (p *abuffer) Get(logger *Logger) *entry {
	b := p.pool.Get().(*entry)
	b.logger = logger
	return b
}

func (p *abuffer) Put(b *entry) {
	b.buf = b.buf[:512]
	b.kvs = b.kvs[:10]
	p.pool.Put(b)
}
