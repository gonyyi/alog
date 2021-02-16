package alog

type bufBasic struct {
	buf BufferBuf
}

func (s *bufBasic) Init(headCap, mainCap int) {
	s.buf = BufferBuf{}
	s.buf.Init(headCap, mainCap)
}

func (s *bufBasic) Get() *BufferBuf {
	s.buf.Reset()
	return &s.buf
}

func (s *bufBasic) Reset(*BufferBuf) {
}
