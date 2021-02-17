package alog

import "sync"

// This is the main buffer
type buffer struct {
	HeadCap int
	BodyCap int
	Prefix  []byte
	Head    []byte
	Body    []byte
	Suffix  []byte
}

func (pi *buffer) Init(prefix, suffix []byte, headCap, bodyCap int) {
	if headCap < 1 {
		headCap = 256
	}
	if bodyCap < 1 {
		bodyCap = 1024
	}
	pi.HeadCap, pi.BodyCap = headCap, bodyCap
	pi.Prefix = prefix
	pi.Suffix = suffix
	pi.Head = make([]byte, headCap)
	pi.Body = make([]byte, bodyCap)
}

func (pi *buffer) Reset() {
	pi.Head = pi.Head[:0]
	pi.Body = pi.Body[:0]
}

type bufSyncPool struct {
	pool sync.Pool
}

func (p *bufSyncPool) Init(prefix, suffix []byte, headCap, bodyCap int) {
	p.pool = sync.Pool{
		New: func() interface{} {
			b := buffer{}
			b.Init(prefix, suffix, headCap, bodyCap)
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
	//b.Head = b.Head[:b.HeadCap]
	//b.Body = b.Body[:b.BodyCap]
	p.pool.Put(b)
}
