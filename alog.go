package alog

import (
	"io"
	"time"
)

// New will return a Alog logger pointer with default values.
// This function will take an io.Writer and convert it to AlWriter.
// A user's custom AlWriter will let the user steer more control.
func New(w io.Writer) *Logger {
	l := Logger{
		out: w,
		buf: newAbuffer(512),
	}
	if w == nil {
		l.out = discard{}
	}
	l.fmt.init()
	l.Control.Level = Linfo
	l.Control.Tag = 0

	return &l
}

// logger is a main struct for Alog.
type Logger struct {
	out io.Writer
	buf abuffer
	fmt formatd
	// Formatter Formatter
	Control control
	Format  Format
}

func (l *Logger) NewTag(name string) Tag {
	return l.Control.Tags.MustGetTag(name)
}

// Do will run (series of) function(s) and is used for
// quick macro like settings for the logger.
func (l *Logger) Do(fns ...func(*Logger)) {
	for _, f := range fns {
		f(l)
	}
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

func (l *Logger) Trace(tag Tag) *entry {
	return l.logs(tag, Ltrace)
}
func (l *Logger) Debug(tag Tag) *entry {
	return l.logs(tag, Ldebug)
}
func (l *Logger) Info(tag Tag) *entry {
	return l.logs(tag, Linfo)
}
func (l *Logger) Warn(tag Tag) *entry {
	return l.logs(tag, Lwarn)
}
func (l *Logger) Error(tag Tag) *entry {
	return l.logs(tag, Lerror)
}
func (l *Logger) Fatal(tag Tag) *entry {
	return l.logs(tag, Lfatal)
}

func (l *Logger) logs(tag Tag, level Level) *entry {
	if (l.Control.CheckFn(level, tag) || l.Control.Check(level, tag)) && l.out != nil {
		buf := l.buf.Get(l)
		buf.tag = tag
		buf.level = level
		buf.buf = buf.buf[:0]
		buf.kvs = buf.kvs[:0]
		return buf
	}
	return nil
}

func (e *entry) Bool(key string, val bool) *entry {
	if e != nil {
		e.kvs = append(e.kvs, KeyVal{
			t: kvBool,
			k: key,
			b: val,
		})
	}
	return e
}
func (e *entry) Float(key string, val float64) *entry {
	if e != nil {
		e.kvs = append(e.kvs, KeyVal{
			t: kvFloat64,
			k: key,
			f64: val,
		})
	}
	return e
}
func (e *entry) Str(key string, val string) *entry {
	if e != nil {
		e.kvs = append(e.kvs, KeyVal{
			t: kvString,
			k: key,
			s: val,
		})
	}
	return e
}

func (e *entry) Int(key string, val int) *entry {
	if e != nil {
		e.kvs = append(e.kvs, KeyVal{
			t: kvInt,
			k: key,
			i: int64(val),
		})
	}
	return e
}

func (e *entry) Write(s string) {
	if e != nil {
		defer e.logger.buf.Put(e)
		e.buf = e.logger.fmt.addBegin(e.buf)

		// INTERFACE: AppendTime()
		if e.logger.Format&fUseTime != 0 {
			t := time.Now()
			if (FtimeUnix|FtimeUnixMs)&e.logger.Format != 0 {
				e.buf = e.logger.fmt.addKeyUnsafe(e.buf, "ts")
				if FtimeUnixMs&e.logger.Format != 0 {
					e.buf = e.logger.fmt.addTimeUnix(e.buf, t.UnixNano()/1e6)
				} else {
					e.buf = e.logger.fmt.addTimeUnix(e.buf, t.Unix())
				}
			} else {
				if FtimeUTC&e.logger.Format != 0 {
					t = t.UTC()
				}
				if Fdate&e.logger.Format != 0 {
					e.buf = e.logger.fmt.addKeyUnsafe(e.buf, "date")
					y, m, d := t.Date()
					e.buf = e.logger.fmt.addTimeDate(e.buf, y, int(m), d)
				}
				if FdateDay&e.logger.Format != 0 {
					e.buf = e.logger.fmt.addKeyUnsafe(e.buf, "day")
					e.buf = e.logger.fmt.addTimeDay(e.buf, int(t.Weekday()))
				}
				if (Ftime|FtimeMs)&e.logger.Format != 0 {
					e.buf = e.logger.fmt.addKeyUnsafe(e.buf, "time")
					h, m, s := t.Clock()
					if FtimeMs&e.logger.Format != 0 {
						e.buf = e.logger.fmt.addTimeMs(e.buf, h, m, s, t.Nanosecond())
					} else {
						e.buf = e.logger.fmt.addTime(e.buf, h, m, s)
					}
				}
			}
		}

		// INTERFACE: LEVEL
		if e.logger.Format&Flevel != 0 {
			e.buf = e.logger.fmt.addKeyUnsafe(e.buf, "level")
			e.buf = e.logger.fmt.addLevel(e.buf, e.level)
		}

		// INTERFACE: TAG
		if e.logger.Format&Ftag != 0 {
			e.buf = e.logger.fmt.addKeyUnsafe(e.buf, "tag")
			e.buf = e.logger.fmt.addTag(e.buf, &e.logger.Control.Tags, e.tag)
		}

		// INTERFACE: MSG
		e.buf = e.logger.fmt.addKeyUnsafe(e.buf, "message")
		if ok, _ := e.logger.fmt.isSimpleStr(s); ok {
			e.buf = e.logger.fmt.addValStringUnsafe(e.buf, s)
		} else {
			e.buf = e.logger.fmt.addValString(e.buf, s)
		}

		// INTERFACE: ADD kvs
		for i := 0; i < len(e.kvs); i++ {
			// Set name
			e.buf = e.logger.fmt.addKeyUnsafe(e.buf, e.kvs[i].k)
			switch e.kvs[i].t {
			case kvInt:
				e.buf = e.logger.fmt.addValInt(e.buf, e.kvs[i].i)
			case kvString:
				if ok, _ := e.logger.fmt.isSimpleStr(e.kvs[i].s); ok {
					e.buf = e.logger.fmt.addValStringUnsafe(e.buf, e.kvs[i].s)
				} else {
					e.buf = e.logger.fmt.addValString(e.buf, e.kvs[i].s)
				}
			case kvBool:
				e.buf = e.logger.fmt.addValBool(e.buf, e.kvs[i].b)
			case kvFloat64:
				e.buf = e.logger.fmt.addValFloat(e.buf, e.kvs[i].f64)
			default:
				e.buf = e.logger.fmt.addValStringUnsafe(e.buf, "err.unexpected")
			}
		}

		// INTERFACE: FINALIZE
		e.buf = e.logger.fmt.addEnd(e.buf)
		if e.logger.out!=nil {
			e.logger.out.Write(e.buf)
		}
	}
}
