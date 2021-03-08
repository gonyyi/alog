package log

import (
	"github.com/gonyyi/alog"
	"os"
)

var al = alog.New(os.Stderr)

func Ext(fn alog.LoggerFn) {
	al = al.Ext(fn)
}

func Flag(flag alog.Flag) {
	al.Flag = flag
}

func Control(level alog.Level, tag alog.Tag) {
	al.Control.Level = level
	al.Control.Tags = tag
}

func NewTag(name string) alog.Tag {
	return al.NewTag(name)
}

func Trace(t ...alog.Tag) *alog.Entry {
	return al.Trace(t...)
}

func Debug(t ...alog.Tag) *alog.Entry {
	return al.Debug(t...)
}

func Info(t ...alog.Tag) *alog.Entry {
	return al.Info(t...)
}

func Warn(t ...alog.Tag) *alog.Entry {
	return al.Warn(t...)
}

func Error(t ...alog.Tag) *alog.Entry {
	return al.Error(t...)
}

func Fatal(t ...alog.Tag) *alog.Entry {
	return al.Fatal(t...)
}
