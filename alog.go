package alog

import (
	"io"
)

// New will return a Alog logger pointer with default values.
// This function will take an io.Writer and convert it to AlWriter.
// A user'Vstr custom AlWriter will let the user steer more control.
func New(w io.Writer) *Logger {
	l := Logger{
		outd: w,
		pool: newEntryPool(),
	}
	if w == nil {
		l.outd = io.Discard
	}
	//l.fmat.init()
	l.Control.Level = Linfo
	l.Control.Tag = 0
	l.Control.TagBucket = &TagBucket{}
	l.Flag = Fdefault

	return &l
}

// logger is a main struct for Alog.
type Logger struct {
	outd    io.Writer
	pool    entryPool
	orFmtr  Formatter
	Control control
	Flag    Flag
}

func (l *Logger) NewTag(name string) Tag {
	return l.Control.TagBucket.MustGetTag(name)
}

// Do will run (series of) function(Vstr) and is used for
// quick macro like settings for the logger.
func (l *Logger) Do(fn func(*Logger)) *Logger {
	fn(l)
	return l
}

// Close will close io.Writer if applicable
func (l *Logger) Close() error {
	if l.orFmtr != nil {
		return l.orFmtr.Close()
	}
	if c, ok := l.outd.(io.Closer); ok && c != nil {
		return c.Close()
	}
	return nil
}

// SetOutput will set the output writer to be used
// in the logger. If nil is given, it will discard the output.
func (l *Logger) SetOutput(w io.Writer) *Logger {
	l.outd = w
	if w == nil {
		l.outd = io.Discard
	}
	return l
}

// Output will return currently used default writer.
func (l *Logger) Output() io.Writer {
	return l.outd
}

// SetFormatter will take an object with Formatter interface
// For Alog, nil can be used to disable the override.
func (l *Logger) SetFormatter(f Formatter) *Logger {
	l.orFmtr = f
	if l.orFmtr != nil {
		l.orFmtr.Init(l.outd, l.Flag, *l.Control.TagBucket)
	}
	return l
}

// getEntry gets entry from the entry pool. This is the very first point
// where it evaluate if the tag/level is loggable.
func (l *Logger) getEntry(tag Tag, level Level) *entry {
	if (l.Control.CheckFn(level, tag) || l.Control.Check(level, tag)) && l.outd != nil {
		buf := l.pool.Get(l)
		buf.tag = tag
		buf.level = level
		buf.buf = buf.buf[:0]
		buf.kvs = buf.kvs[:0]
		return buf
	}
	return nil
}

// Trace takes a tag (0 for no tag) and returns an entry point.
func (l *Logger) Trace(tag Tag) *entry {
	return l.getEntry(tag, Ltrace)
}

// Debug takes a tag (0 for no tag) and returns an entry point.
func (l *Logger) Debug(tag Tag) *entry {
	return l.getEntry(tag, Ldebug)
}

// Info takes a tag (0 for no tag) and returns an entry point.
func (l *Logger) Info(tag Tag) *entry {
	return l.getEntry(tag, Linfo)
}

// Warn takes a tag (0 for no tag) and returns an entry point.
func (l *Logger) Warn(tag Tag) *entry {
	return l.getEntry(tag, Lwarn)
}

// Error takes a tag (0 for no tag) and returns an entry point.
func (l *Logger) Error(tag Tag) *entry {
	return l.getEntry(tag, Lerror)
}

// Fatal takes a tag (0 for no tag) and returns an entry point.
func (l *Logger) Fatal(tag Tag) *entry {
	return l.getEntry(tag, Lfatal)
}
