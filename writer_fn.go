package alog

type FnWrite func(b []byte, level Level, tag Tag) (int, error)
type FnClose func() error

type fnWriter struct {
	fnWrite FnWrite
	fnClose FnClose
}

func (w fnWriter) WriteLt(p []byte, level Level, tag Tag) (int, error) {
	return w.fnWrite(p, level, tag)
}

func (w fnWriter) Write(p []byte) (n int, err error) {
	return w.fnWrite(p, 0, 0)
}

func (w fnWriter) Close() error {
	return w.fnClose()
}

func NewFnWriter(fnWrite FnWrite, fnClose FnClose) *fnWriter {
	if fnWrite == nil {
		fnWrite = func(b []byte, level Level, tag Tag) (int, error) { return 0, nil }
	}
	if fnClose == nil {
		fnClose = func() error { return nil }
	}
	return &fnWriter{
		fnWrite: fnWrite,
		fnClose: fnClose,
	}
}
