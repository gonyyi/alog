package alog

import "sync"

// buffer is the main buffer format used in Alog.
// This can be used as standalone or within the sync.Pool
type buffer struct {
	HeadCap int
	BodyCap int
	Head    []byte
	Body    []byte
}

func (pi *buffer) Init(headCap, bodyCap int) {
	if headCap < 1 {
		headCap = 256
	}
	if bodyCap < 1 {
		bodyCap = 1024
	}
	pi.HeadCap, pi.BodyCap = headCap, bodyCap
	pi.Head = make([]byte, headCap)
	pi.Body = make([]byte, bodyCap)
}

// Reset will resize the buffers. In case a large size data came in,
// this will reset the size of buffers.
func (pi *buffer) Reset() {
	pi.Head = pi.Head[:pi.HeadCap]
	pi.Body = pi.Body[:pi.BodyCap]
}

// Future use
//	type Buffer interface {
//		Init(headCap, bodyCap int)
//		Get() *buffer
//		Reset(b *buffer)
//	}

// bufSyncPool is an a Buffer implementation of sync.Pool.
type bufSyncPool struct {
	pool sync.Pool
}

func (p *bufSyncPool) Init(headCap, bodyCap int) {
	p.pool = sync.Pool{
		New: func() interface{} {
			b := buffer{}
			b.Init(headCap, bodyCap)
			return &b
		},
	}
}

func (p *bufSyncPool) Get() *buffer {
	b := p.pool.Get().(*buffer)
	b.Body = b.Body[:0]
	b.Head = b.Head[:0]
	return b
}

func (p *bufSyncPool) Reset(b *buffer) {
	b.Reset()
	p.pool.Put(b)
}
