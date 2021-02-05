package alog

type alogItem struct {
	bufCap int
	buf    []byte // this is a buffer that will be created multiple and used by multiple goroutines by sync.Pool
}

func newItem(size int) *alogItem {
	return &alogItem{
		bufCap: size,
		buf:    make([]byte, size),
	}
}

func (i *alogItem) reset() {
	i.buf = i.buf[:0]
}
