package alog

import "sync"

type bufSyncPool struct {
	pool sync.Pool
}

func (p *bufSyncPool) Init(headCap, mainCap int) {
	p.pool = sync.Pool{
		New: func() interface{} {
			b := BufferBuf{
				HeadCap: headCap,
				MainCap: mainCap,
			}
			b.Init(headCap, mainCap)
			return &b
		},
	}
}

func (p *bufSyncPool) Get() *BufferBuf {
	b := p.pool.Get().(*BufferBuf)
	b.Main = b.Main[:0]
	b.Head = b.Head[:0]
	return b
}

func (p *bufSyncPool) Reset(b *BufferBuf) {
	b.Head = b.Head[:0]
	b.Main = b.Main[:0]
	p.pool.Put(b)
}
