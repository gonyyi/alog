package alog

// this is a buffer that will be created multiple and used by multiple goroutines by sync.Pool
type poolbuf struct {
	capHeader int
	capMain   int
	bufHeader []byte // header only
	bufMain   []byte // main log data
}

func newPoolbuf(cheader, cmain int) *poolbuf {
	return &poolbuf{
		capHeader: cheader, // header only
		capMain:   cmain,   // main log data
		bufHeader: make([]byte, cheader),
		bufMain:   make([]byte, cmain),
	}
}

func (i *poolbuf) reset() {
	i.bufHeader = i.bufHeader[:i.capHeader] // this will shrink the size
	i.bufMain = i.bufMain[:i.capMain]       // this will shrink the size
}
