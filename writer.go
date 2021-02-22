package alog

import (
	"io"
)

// AlWriter is a combination of a formatter and a writer for Alog.
// It will be in charge of its own buf and formatting.
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
	if w == nil {
		w = discard
	}
	// If not, just create a compatible one.
	return &alWriter{w: w}
}

// alWriter is a default writer
type alWriter struct {
	w io.Writer
}

// WriteLvt will take level and tag. This is to be used as a conditional writer.
// Eg. When certain tag and/or level, this can write it to a different place.
func (alw alWriter) WriteTag(lv Level, tag Tag, head, body []byte) (int, error) {
	body = append(body, '\n')
	return alw.w.Write(append(head, body...))
}
func (alw alWriter) Write(p []byte) (n int, err error) {
	return alw.w.Write(p)
}

func (alw alWriter) Close() error {
	if c, ok := alw.w.(io.Closer); ok && c != nil {
		return c.Close()
	}
	return nil
}

// SubWriter is a writer with predefined Level and Tag.
type SubWriter struct {
	l      *Logger
	dLevel Level // default level for the SubWriter
	dTag   Tag   // default tag for the SubWriter
}

// Write is to be used as io.Writer interface
func (w SubWriter) Write(p []byte) (n int, err error) { return w.l.logb(w.dLevel, w.dTag, p) }
func (w SubWriter) Trace(s string)                    { w.l.Log(Ltrace, w.dTag, s) }
func (w SubWriter) Debug(s string)                    { w.l.Log(Ldebug, w.dTag, s) }
func (w SubWriter) Info(s string)                     { w.l.Log(Linfo, w.dTag, s) }
func (w SubWriter) Warn(s string)                     { w.l.Log(Lwarn, w.dTag, s) }
func (w SubWriter) Error(s string)                    { w.l.Log(Lerror, w.dTag, s) }
func (w SubWriter) Fatal(s string)                    { w.l.Log(Lfatal, w.dTag, s) }

// discard will be used instead of ioutil.Discard
const discard = devnull(true)

// devNull is a type for discard
type devnull bool

// Write discards everything
func (devnull) Write([]byte) (int, error) {
	return 0, nil
}
