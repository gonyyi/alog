package alog

import "sync"

// bufItem is the main bufItem format used in Alog.
// This can be used as standalone or within the sync.Pool
type bufItem struct {
	Buf []byte
}

func newAbuffer(size int) abuffer {
	return abuffer{
		size: size,
		pool: sync.Pool{
			New: func() interface{} {
				return &bufItem{
					Buf: make([]byte, size),
				}
			},
		},
	}
}

// buf is an a Buffer implementation of sync.Pool.
type abuffer struct {
	size int
	pool sync.Pool
}

func (p *abuffer) Get() *bufItem {
	b := p.pool.Get().(*bufItem)
	b.Buf = b.Buf[:0]
	return b
}

func (p *abuffer) Put(b *bufItem) {
	b.Buf = b.Buf[:p.size]
	p.pool.Put(b)
}
