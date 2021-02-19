package alog

import "sync"

// buf is the main buf format used in Alog.
// This can be used as standalone or within the sync.Pool
type buf struct {
	HeadCap int
	BodyCap int
	Head    []byte
	Body    []byte
}

func (pi *buf) Init(headCap, bodyCap int) {
	if headCap < 1 {
		headCap = 512
	}
	if bodyCap < 1 {
		bodyCap = 2048
	}
	pi.HeadCap, pi.BodyCap = headCap, bodyCap
	pi.Head = make([]byte, headCap)
	pi.Body = make([]byte, bodyCap)
}

// Reset will resize the buffers. In case a large size data came in,
// this will reset the size of buffers.
func (pi *buf) Reset() {
	pi.Head = pi.Head[:pi.HeadCap]
	pi.Body = pi.Body[:pi.BodyCap]
}

type Buffer interface {
	Init(headCap, bodyCap int)
	Get() *buf
	Reset(b *buf)
}

// bufSyncPool is an a Buffer implementation of sync.Pool.
type bufSyncPool struct {
	pool sync.Pool
}

func (p *bufSyncPool) Init(headCap, bodyCap int) {
	p.pool = sync.Pool{
		New: func() interface{} {
			b := buf{}
			b.Init(headCap, bodyCap)
			return &b
		},
	}
}

func (p *bufSyncPool) Get() *buf {
	b := p.pool.Get().(*buf)
	b.Body = b.Body[:0]
	b.Head = b.Head[:0]
	return b
}

func (p *bufSyncPool) Reset(b *buf) {
	b.Reset()
	p.pool.Put(b)
}
