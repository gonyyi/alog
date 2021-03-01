package alog

import (
	"io"
)

// devNull is a type for discard
type discard struct{}

// Write discards everything
func (discard) Write([]byte) (int, error) {
	return 0, nil
}

// AlWriter is a combination of a formatter and a writer for Alog.
// It will be in charge of its own bufItem and formatting.
type AlWriter interface {
	// Write for a Standard io.Writer method
	Write(p []byte) (int, error)
	// WriteLvt will take level and tag. This is to be used as a conditional writer.
	// Eg. When certain tag and/or level, this can write it to a different io.Writer.
	WriteTag(lv Level, tag Tag, head, body []byte) (int, error)
	Close() error
}

// newAlWriter convers io.Writer to AlWriter.
// If it can't be converted, use alWriter.
func newAlWriter(w io.Writer) AlWriter {
	// If nil is given, consider it a discard
	if w == nil {
		w = discard{}
	}
	// If AlWriter is given, use as is.
	if c, ok := w.(AlWriter); ok && c != nil {
		return c
	}
	// If not, create one.
	return &alWriter{w: w}
}

// alWriter is a default writer
type alWriter struct {
	w io.Writer
}

// WriteLvt will take level and tag. This is to be used as a conditional writer.
// Eg. When certain tag and/or level, this can write it to a different place.
func (alw *alWriter) WriteTag(lv Level, tag Tag, head, body []byte) (int, error) {
	body = append(body, '\n')
	return alw.w.Write(append(head, body...))
}
func (alw *alWriter) Write(p []byte) (n int, err error) {
	return alw.w.Write(p)
}

func (alw *alWriter) Close() error {
	if c, ok := alw.w.(io.Closer); ok && c != nil {
		return c.Close()
	}
	return nil
}
