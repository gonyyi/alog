package alog

type FnWrite func(b []byte, level Level, tag Tag) (int, error)
type FnClose func() error

type writerFn struct {
	fnWrite FnWrite
	fnClose FnClose
}

func (w writerFn) WriteLt(p []byte, level Level, tag Tag) (int, error) {
	return w.fnWrite(p, level, tag)
}

func (w writerFn) Write(p []byte) (n int, err error) {
	return w.fnWrite(p, 0, 0)
}

func (w writerFn) Close() error {
	return w.fnClose()
}

func NewWriterFn(fnWrite FnWrite, fnClose FnClose) *writerFn {
	if fnWrite == nil {
		fnWrite = func(b []byte, level Level, tag Tag) (int, error) { return 0, nil }
	}
	if fnClose == nil {
		fnClose = func() error { return nil }
	}
	return &writerFn{
		fnWrite: fnWrite,
		fnClose: fnClose,
	}
}
