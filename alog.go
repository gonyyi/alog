package alog

import (
	"io"
)

// New will return a Alog logger pointer with default values.
// This function will take an io.Writer and convert it to AlWriter.
// A user'Vstr custom AlWriter will let the user steer more control.
func New(w io.Writer) Logger {
	l := Logger{
		w:    w,
		pool: newEntryPool(),
	}
	if w == nil {
		l.w = io.Discard
	}
	//l.fmat.init()
	l.Control.Level = Linfo
	l.Control.Tag = 0
	l.Control.TagBucket = &TagBucket{}
	l.Flag = Fdefault

	return l
}

// logger is a main struct for Alog.
type Logger struct {
	w       io.Writer
	pool    *entryPool
	orFmtr  Formatter
	Control control
	Flag    Flag
}

// NewTag will create a new tag
// Using value receiver as this won't be used many times anyway
func (l Logger) NewTag(name string) Tag {
	return l.Control.TagBucket.MustGetTag(name)
}

// Do will run (series of) function(Vstr) and is used for
// quick macro like settings for the logger.
func (l Logger) Do(fn DoFn) Logger {
	return fn(l)
}

// Close will close io.Writer if applicable
func (l *Logger) Close() error {
	if l.orFmtr != nil {
		return l.orFmtr.Close()
	}
	if c, ok := l.w.(io.Closer); ok && c != nil {
		return c.Close()
	}
	return nil
}

// SetOutput will set the output writer to be used
// in the logger. If nil is given, it will discard the output.
func (l Logger) SetOutput(w io.Writer) Logger {
	l.w = w
	if w == nil {
		l.w = io.Discard
	}
	return l
}

// Output will return currently used default writer.
func (l Logger) Output() io.Writer {
	return l.w
}

// SetFormatter will take an object with Formatter interface
// For Alog, nil can be used to disable the override.
func (l Logger) SetFormatter(f Formatter) Logger {
	l.orFmtr = f
	if l.orFmtr != nil {
		l.orFmtr.Init(l.w, l.Flag, *l.Control.TagBucket)
	}
	return l
}

// getEntry gets entry from the entry pool. This is the very first point
// where it evaluate if the tag/level is loggable.
func (l *Logger) getEntry(tag Tag, level Level) *entry {
	if l.Control.CheckFn(level, tag) || l.Control.Check(level, tag) {
		e := l.pool.Get(entryInfo{
			flag:    l.Flag,
			tbucket: l.Control.TagBucket,
			pool:    l.pool,
			orFmtr:  l.orFmtr,
			w:       l.w,
		})
		e.tag = tag
		e.level = level
		// buf and kvs are reset when *entryPool.Put()
		// e.buf = e.buf[:0]
		// e.kvs = e.kvs[:0]
		return e
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
