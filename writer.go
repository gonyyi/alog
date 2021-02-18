package alog

import (
	"io"
)

// AlWriter is a combination of a formatter and a writer for Alog.
// It will be in charge of its own buf and formatting.
type AlWriter interface {
	io.Writer
	// WriteLvt will take level and tag. This is to be used as a conditional writer.
	// Eg. When certain tag and/or level, this can write it to a different io.Writer.
	WriteTag(lv Level, tag Tag, head, body []byte) (int, error)
}

// toAlWriter convers io.Writer to AlWriter.
// If it can't be converted, use alwBasic.
func toAlWriter(w io.Writer) AlWriter {
	if w == nil {
		w = discard
	}
	// If AlWriter is given, take it
	if alw, ok := w.(AlWriter); ok {
		return alw
	}
	// If not, just create a compatible one.
	return &alwBasic{w}
}

// alWriter is a default writer
type alwBasic struct {
	io.Writer
}

// WriteLvt will take level and tag. This is to be used as a conditional writer.
// Eg. When certain tag and/or level, this can write it to a different place.
func (alwp alwBasic) WriteTag(lv Level, tag Tag, head, body []byte) (int, error) {
	body = append(body, '\n')
	return alwp.Write(append(head, body...))
}

// SubWriter is a writer with predefined Level and Tag.
type SubWriter struct {
	l      *Logger
	dLevel Level // default level for the SubWriter
	dTag   Tag   // default tag for the SubWriter
}

// Write is to be used as io.Writer interface
func (w *SubWriter) Write(p []byte) (n int, err error) { return w.l.logb(w.dLevel, w.dTag, p) }
func (w *SubWriter) Trace(s string)                    { w.l.Log(Ltrace, w.dTag, s) }
func (w *SubWriter) Debug(s string)                    { w.l.Log(Ldebug, w.dTag, s) }
func (w *SubWriter) Info(s string)                     { w.l.Log(Linfo, w.dTag, s) }
func (w *SubWriter) Warn(s string)                     { w.l.Log(Lwarn, w.dTag, s) }
func (w *SubWriter) Error(s string)                    { w.l.Log(Lerror, w.dTag, s) }
func (w *SubWriter) Fatal(s string)                    { w.l.Log(Lfatal, w.dTag, s) }

// discard will be used instead of ioutil.Discard
const discard = devnull(true)

// devNull is a type for discard
type devnull bool

// Write discards everything
func (devnull) Write([]byte) (int, error) {
	return 0, nil
}
