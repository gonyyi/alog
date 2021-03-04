package log

import (
	"github.com/gonyyi/alog"
	"os"
)

var al = alog.New(os.Stderr)

func Do(fn alog.DoFn) {
	al = al.Do(fn)
}

func NewTag(name string) alog.Tag {
	return al.NewTag(name)
}

func Trace(t alog.Tag) *alog.Entry {
	return al.Trace(t)
}

func Debug(t alog.Tag) *alog.Entry {
	return al.Debug(t)
}

func Info(t alog.Tag) *alog.Entry {
	return al.Info(t)
}

func Warn(t alog.Tag) *alog.Entry {
	return al.Warn(t)
}

func Error(t alog.Tag) *alog.Entry {
	return al.Error(t)
}

func Fatal(t alog.Tag) *alog.Entry {
	return al.Fatal(t)
}
