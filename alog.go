package alog

import (
	"io"
)

// New will return a Alog logger pointer with default values.
// This function will take an io.Writer and convert it to AlWriter.
// A user'Vstr custom AlWriter will let the user steer more control.
func New(w io.Writer) *Logger {
	l := Logger{
		out:  w,
		pool: newEntryPool(),
	}
	if w == nil {
		l.out = discard{}
	}
	//l.fmat.init()
	l.Control.Level = Linfo
	l.Control.Tag = 0
	l.Control.Tags = &TagBucket{}
	l.Format = Fdefault

	return &l
}

var dFmt formatd

func init() {
	dFmt.init()
}

// logger is a main struct for Alog.
type Logger struct {
	out  io.Writer
	pool entryPool
	//fmat formatd
	CusFmat Formatter
	Control control
	Format  Format
}

func (l *Logger) NewTag(name string) Tag {
	return l.Control.Tags.MustGetTag(name)
}

// Do will run (series of) function(Vstr) and is used for
// quick macro like settings for the logger.
func (l *Logger) Do(fn func(*Logger)) *Logger {
	fn(l)
	return l
}

// Close will close io.Writer if applicable
func (l *Logger) Close() error {
	if c, ok := l.out.(io.Closer); ok && c != nil {
		return c.Close()
	}
	return nil
}

// SetOutput will set the output writer to be used
// in the logger. If nil is given, it will discard the output.
func (l *Logger) SetOutput(w io.Writer) *Logger {
	l.out = newAlWriter(w)
	return l
}
func (l *Logger) getEntry(tag Tag, level Level) *entry {
	if (l.Control.CheckFn(level, tag) || l.Control.Check(level, tag)) && l.out != nil {
		buf := l.pool.Get(l)
		buf.tag = tag
		buf.level = level
		buf.buf = buf.buf[:0]
		buf.kvs = buf.kvs[:0]
		return buf
	}
	return nil
}

func (l *Logger) Trace(tag Tag) *entry {
	return l.getEntry(tag, Ltrace)
}
func (l *Logger) Debug(tag Tag) *entry {
	return l.getEntry(tag, Ldebug)
}
func (l *Logger) Info(tag Tag) *entry {
	return l.getEntry(tag, Linfo)
}
func (l *Logger) Warn(tag Tag) *entry {
	return l.getEntry(tag, Lwarn)
}
func (l *Logger) Error(tag Tag) *entry {
	return l.getEntry(tag, Lerror)
}
func (l *Logger) Fatal(tag Tag) *entry {
	return l.getEntry(tag, Lfatal)
}
