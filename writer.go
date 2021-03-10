package alog

import "io"

// Writer is a Level and Tag writer
type Writer interface {
	WriteLt([]byte, Level, Tag) (int, error)
	Close() error
}

// iowToAlw converts io.Writer to alog.Writer interface compatible
func iowToAlw(w io.Writer) Writer {
	if alw, ok := w.(Writer); ok {
		return alw
	}
	return alwAdapter{w: w}
}

// alwAdapter is a struct to meet alog.Writer interface.
// This will be created by iowToAlw
type alwAdapter struct {
	w io.Writer
}

// WriteLt to meet requirements for alog.Writer interface.
// Since this is an adapter, it will write it to io.Writer
// regardless of level or tag.
func (w alwAdapter) WriteLt(p []byte, lvl Level, tag Tag) (int, error) {
	return w.w.Write(p)
}

// Close will close the io.Writer if it supports
func (w alwAdapter) Close() error {
	if c, ok := w.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

// devNull is a type for discard
type Discard struct{}

// Write discards everything
func (Discard) Write([]byte) (int, error) {
	return 0, nil
}

// WriteLt is for Writer compatible; discards everything
func (Discard) WriteLt([]byte, Level, Tag) (int, error) {
	return 0, nil
}

// Close to meet alog.Writer interface
func (Discard) Close() error {
	return nil
}
