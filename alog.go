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

// Logger is a main struct for Alog.
type Logger struct {
	out io.Writer
	buf abuffer
	fmt formatd
	// Formatter Formatter
	Control control
	Format  Format
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

// Log is a main method of logging. This takes level, tag, message, as well as optional
// data. The optional data has to be in pairs of name and value. For the speed, Alog only
// supports basic types: int, int64, uint, string, bool, float32, float64.
func (l *Logger) Log(level Level, tag Tag, msg string, kvs ...KeyVal) {
	if (l.Control.CheckFn(level, tag) || l.Control.Check(level, tag)) && l.out != nil {
		buf := l.buf.Get()
		defer l.buf.Put(buf)

		buf.Buf = l.fmt.addBegin(buf.Buf)

		// INTERFACE: AppendTime()
		if l.Format&fUseTime != 0 {
			t := time.Now()
			if FtimeUnix&l.Format != 0 {
				buf.Buf = l.fmt.addKeyUnsafe(buf.Buf, "ts")
				buf.Buf = l.fmt.addTimeUnix(buf.Buf, t.Unix())
			} else {
				if FtimeUTC&l.Format != 0 {
					t = t.UTC()
				}
				if Fdate&l.Format != 0 {
					buf.Buf = l.fmt.addKeyUnsafe(buf.Buf, "date")
					y, m, d := t.Date()
					buf.Buf = l.fmt.addTimeDate(buf.Buf, y, int(m), d)
				}
				if FdateDay&l.Format != 0 {
					buf.Buf = l.fmt.addKeyUnsafe(buf.Buf, "day")
					buf.Buf = l.fmt.addTimeDay(buf.Buf, int(t.Weekday()))
				}
				if (Ftime|FtimeMs)&l.Format != 0 {
					buf.Buf = l.fmt.addKeyUnsafe(buf.Buf, "time")
					h, m, s := t.Clock()
					if FtimeMs&l.Format != 0 {
						buf.Buf = l.fmt.addTimeMs(buf.Buf, h, m, s, t.Nanosecond()/1e6)
					} else {
						buf.Buf = l.fmt.addTime(buf.Buf, h, m, s)
					}
				}
			}
		}

		// INTERFACE: LEVEL
		if l.Format&Flevel != 0 {
			buf.Buf = l.fmt.addKeyUnsafe(buf.Buf, "level")
			buf.Buf = l.fmt.addLevel(buf.Buf, level)
		}

		// INTERFACE: TAG
		if l.Format&Ftag != 0 {
			buf.Buf = l.fmt.addKeyUnsafe(buf.Buf, "tag")
			buf.Buf = l.fmt.addTag(buf.Buf, &l.Control.Tags, tag)
		}

		// INTERFACE: MSG
		buf.Buf = l.fmt.addKeyUnsafe(buf.Buf, "message")
		if l.fmt.isSimpleStr(msg) {
			buf.Buf = l.fmt.addValStringUnsafe(buf.Buf, msg)
		} else {
			buf.Buf = l.fmt.addValString(buf.Buf, msg)
		}

		// INTERFACE: SEPARATOR

		// INTERFACE: KV
		if kvs != nil {
			for i := 0; i < len(kvs); i++ {
				buf.Buf = l.fmt.addKeyUnsafe(buf.Buf, kvs[i].k)
				switch kvs[i].t {
				case kvString:
					if l.fmt.isSimpleStr(kvs[i].s) {
						buf.Buf = l.fmt.addValStringUnsafe(buf.Buf, kvs[i].s)
					} else {
						buf.Buf = l.fmt.addValString(buf.Buf, kvs[i].s)
					}
				case kvInt:
					buf.Buf = l.fmt.addValInt(buf.Buf, kvs[i].i)
				case kvBool:
					buf.Buf = l.fmt.addValBool(buf.Buf, kvs[i].b)
				case kvFloat64:
					buf.Buf = l.fmt.addValFloat(buf.Buf, kvs[i].f64)
				}
			}
		}

		buf.Buf = l.fmt.addEnd(buf.Buf)
		l.out.Write(buf.Buf)
		return
	}
	return
}
