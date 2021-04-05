package ext

import (
	"github.com/gonyyi/alog"
	"io"
)

func NewFilterWriter(w io.Writer, defaultLevel alog.Level, defaultTag alog.Tag) filterWriter {
	return filterWriter{
		w:        w,
		defLevel: defaultLevel,
		defTag:   defaultTag,
	}
}

type filterWriter struct {
	w        io.Writer
	defLevel alog.Level
	defTag   alog.Tag
}

func (w filterWriter) Write(p []byte) (int, error) {
	return w.WriteLt(p, w.defLevel, w.defTag)
}

func (w filterWriter) WriteLt(p []byte, level alog.Level, tag alog.Tag) (int, error) {
	//if level >= w.defLevel {
	return w.w.Write(p)
	//}
	//return 0, nil
}
