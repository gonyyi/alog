package alog

type Buffer interface {
	Init(headCap, mainCap int)
	Get() *BufferBuf
	Reset(*BufferBuf)
}

// This is the main buffer
type BufferBuf struct {
	HeadCap int
	MainCap int
	Head    []byte
	Main    []byte
}

func (pi *BufferBuf) Init(headCap, mainCap int) {
	if headCap < 0 {
		pi.HeadCap = 256
	}
	if mainCap < 0 {
		pi.MainCap = 1024
	}
	pi.Head = make([]byte, headCap)
	pi.Main = make([]byte, mainCap)
}

func (pi *BufferBuf) Reset() {
	pi.Head = pi.Head[:0]
	pi.Main = pi.Main[:0]
}
