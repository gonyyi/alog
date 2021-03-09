package alog

import (
	"io"
)

// Level const
const (
	TraceLevel Level = iota + 1 // TraceLevel shows trace level, the most detailed debugging level.
	DebugLevel                  // DebugLevel shows debug level or higher
	InfoLevel                   // InfoLevel shows information level or higher
	WarnLevel                   // WarnLevel is for a normal but a significant condition
	ErrorLevel                  // ErrorLevel shows error level or higher
	FatalLevel                  // FatalLevel shows fatal level or higher. This does not exit the process
)

// Flag const
const (
	WithLevel      Flag = 1 << iota // WithLevel show level in the log messsage.
	WithTag                         // WithTag will show tags
	WithDate                        // WithDate will show both CCYY and MMDD
	WithDay                         // WithDay will show 0-6 for JSON or (Sun-Mon)
	WithTime                        // WithTime will show HHMMSS
	WithTimeMs                      // WithTimeMs will show time + millisecond --> JSON: HHMMSS000, Text: HHMMSS,000
	WithUTC                         // WithUTC will show UTC time formats
	WithUnixTime                    // WithUnixTime will show unix time
	WithUnixTimeMs                  // WithUnixTimeMs will show unix time with millisecond

	// UseDefault holds default output format when no option is given.
	WithDefault = WithTime | WithDate | WithLevel | WithTag
	// fHasTime is precalculated time for internal functions. Not that if WithUTC is used by it self,
	// without any below, it won't print any time.
	fHasTime = WithDate | WithDay | WithTime | WithTimeMs | WithUnixTime | WithUnixTimeMs
)

// KeyValue const
const (
	KvInt     kvType = iota + 1 // KvInt indiciates int64 type KeyValue
	KvFloat64                   // KvFloat64 indicates float64 type KeyValue
	KvString                    // KvString indicates string type KeyValue
	KvBool                      // KvBool indicates bool type KeyValue
	KvError                     // KvError indicates error type KeyValue
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
		Flag:    WithDefault,
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
	if l.Control.bucket != nil {
		return l.Control.bucket.MustGetTag(name)
	}
	return 0
}

// Ext will run functions that will act as a
// quick macro like settings for the logger.
// See <https://github.com/gonyyi/alog/ext>
// for examples.
func (l Logger) Ext(fn LoggerFn) Logger {
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
		l.orFmtr.Init(l.w, &l.Flag, l.Control.bucket)
	}
	return l
}

// getEntry gets Entry from the Entry pool. This is the very first point
// where it evaluate if the tag/level is loggable.
func (l *Logger) getEntry(level Level, tags ...Tag) *Entry {
	// If a control function exists, BUT returns false,
	// otherwise, use result from level/tag check.
	var tag Tag = 0
	if len(tags) > 0 {
		for i := 0; i < len(tags); i++ {
			tag = tag | tags[i]
		}
	}
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
func (l *Logger) Trace(tags ...Tag) *Entry {
	return l.getEntry(TraceLevel, tags...)
}

// Debug takes a tag (0 for no tag) and returns an Entry point.
func (l *Logger) Debug(tags ...Tag) *Entry {
	return l.getEntry(DebugLevel, tags...)
}

// Info takes a tag (0 for no tag) and returns an Entry point.
func (l *Logger) Info(tags ...Tag) *Entry {
	return l.getEntry(InfoLevel, tags...)
}

// Warn takes a tag (0 for no tag) and returns an Entry point.
func (l *Logger) Warn(tags ...Tag) *Entry {
	return l.getEntry(WarnLevel, tags...)
}

// Error takes a tag (0 for no tag) and returns an Entry point.
func (l *Logger) Error(tags ...Tag) *Entry {
	return l.getEntry(ErrorLevel, tags...)
}

// Fatal takes a tag (0 for no tag) and returns an Entry point.
func (l *Logger) Fatal(tags ...Tag) *Entry {
	return l.getEntry(FatalLevel, tags...)
}
