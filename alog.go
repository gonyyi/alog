package alog

import (
	"io"
)

// New will return a Alog logger pointer with default values.
// This function will take an io.Writer and convert it to AlWriter.
// A user'Vstr custom AlWriter will let the user steer more control.
func New(w io.Writer) Logger {
	if w == nil {
		w = io.Discard
	}
	return Logger{
		w:       w,
		pool:    newEntryPool(),
		Control: newControl(),
		Flag:    UseDefault,
	}
}

// Logger is a main struct for Alog.
// This struct is 80 bytes.
type Logger struct {
	w       io.Writer
	pool    *entryPool
	orFmtr  Formatter
	Control control // 32 bytes
	Flag    Flag
}

// NewTag will create a new tag
// Using value receiver as this won't be used many times anyway
func (l Logger) NewTag(name string) Tag {
	return l.Control.Bucket().MustGetTag(name)
}

// Do will run functions that will act as a
// quick macro like settings for the logger.
// See <https://github.com/gonyyi/alog/ext>
// for examples.
func (l Logger) Do(fn DoFn) Logger {
	if fn != nil {
		return fn(l)
	}
	return l
}

// Close will close io.Writer if applicable
func (l Logger) Close() error {
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
// See: <https://github.com/gonyyi/alog/ext> for examples.
func (l Logger) SetFormatter(f Formatter) Logger {
	l.orFmtr = f
	if l.orFmtr != nil {
		l.orFmtr.Init(l.w, l.Flag, *l.Control.bucket)
	}
	return l
}

// getEntry gets Entry from the Entry pool. This is the very first point
// where it evaluate if the tag/level is loggable.
func (l *Logger) getEntry(tag Tag, level Level) *Entry {
	// If a control function exists, BUT returns false,
	// otherwise, use result from level/tag check.
	if l.Control.Fn != nil {
		if l.Control.Fn(level, tag) == false {
			return nil
		}
	} else if l.Control.Check(level, tag) == false {
		return nil
	}

	e := l.pool.Get(entryInfo{
		flag:    l.Flag,
		tbucket: l.Control.bucket,
		pool:    l.pool,
		orFmtr:  l.orFmtr,
		w:       l.w,
	})
	e.tag = tag
	e.level = level

	e.buf = e.buf[:0]
	e.kvs = e.kvs[:0]
	return e

}

// Trace takes a tag (0 for no tag) and returns an Entry point.
func (l *Logger) Trace(tag Tag) *Entry {
	return l.getEntry(tag, TraceLevel)
}

// Debug takes a tag (0 for no tag) and returns an Entry point.
func (l *Logger) Debug(tag Tag) *Entry {
	return l.getEntry(tag, DebugLevel)
}

// Info takes a tag (0 for no tag) and returns an Entry point.
func (l *Logger) Info(tag Tag) *Entry {
	return l.getEntry(tag, InfoLevel)
}

// Warn takes a tag (0 for no tag) and returns an Entry point.
func (l *Logger) Warn(tag Tag) *Entry {
	return l.getEntry(tag, WarnLevel)
}

// Error takes a tag (0 for no tag) and returns an Entry point.
func (l *Logger) Error(tag Tag) *Entry {
	return l.getEntry(tag, ErrorLevel)
}

// Fatal takes a tag (0 for no tag) and returns an Entry point.
func (l *Logger) Fatal(tag Tag) *Entry {
	return l.getEntry(tag, FatalLevel)
}
