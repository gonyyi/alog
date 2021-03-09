package alog

import "io"

// devNull is a type for discard
type Discard struct{}

// Write discards everything
func (Discard) Write([]byte) (int, error) {
	return 0, nil
}

// WriteLt is for AlWriter compatible; discards everything
func (Discard) WriteLt([]byte, Level, Tag) (int, error) {
	return 0, nil
}
func (Discard) Close() error {
	return nil
}

// AlWriter is a writer with Level and Tag function
type AlWriter interface {
	WriteLt([]byte, Level, Tag) (int, error)
	Close() error
}

func writerToAlWriter(w io.Writer) AlWriter {
	if alw, ok := w.(AlWriter); ok {
		return alw
	}
	return alWriterAdapter{w: w}
}

type alWriterAdapter struct {
	w io.Writer
}

func (w alWriterAdapter) WriteLt(p []byte, lvl Level, tag Tag) (int, error) {
	return w.w.Write(p)
}
func (w alWriterAdapter) Close() error {
	if c, ok := w.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
