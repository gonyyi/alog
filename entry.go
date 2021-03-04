package alog

import (
	"time"
)

// entry is a log entry will be used with a entryPool to
// reuse the resource.
type entry struct {
	buf    []byte
	level  Level
	tag    Tag
	logger *Logger
	kvs    []KeyValue
}

// Writes will finalize the log message, format it, and
// write it to writer. Besides *logger.getEntry(), this is
// the only other method which isn't inline-able.
func (e *entry) Writes(s string) {
	// When log message was created from *Logger.getEntry(),
	// it examines logability (should log or not). Once it's not eligible,
	// it will return nil.
	// eg. log.Trace(0).Str("name", "gon").Int("age", 39).Write()
	//   1. Trace(0) will call getEntry(), if it determines the log shouldn't be
	//      logged, it will return nil.
	//   2. Next method with Str() receives nil for the pointer and will ignore,
	//      and return nil to next.
	//   3. Int method will receive nil, and just pass nil to next.
	//   4. Write method finally receives it, if entry pointer is nil, it won't
	//      do anything as it's not eligible to log.

	// if entry is not nil (=loggable),
	if e != nil {
		// since pointer receiver *entry is obtained from the pool,
		// make sure this will be put back to memory.
		defer e.logger.pool.Put(e)

		// if custom formatter exists, use it instead of default formatter.
		// for default formatter (formatd), it's a concrete function for speed.
		// rather than using from the interface.
		if e.logger.orFmtr != nil {
			// CUSTOM FORMATTER
			e.buf = e.logger.orFmtr.Begin(e.buf)
			e.buf = e.logger.orFmtr.AddTime(e.buf)
			e.buf = e.logger.orFmtr.AddLevel(e.buf, e.level)
			e.buf = e.logger.orFmtr.AddTag(e.buf, e.tag)
			e.buf = e.logger.orFmtr.AddMsg(e.buf, s)
			e.buf = e.logger.orFmtr.AddKVs(e.buf, e.kvs)
			e.buf = e.logger.orFmtr.End(e.buf)
			e.logger.orFmtr.Write(e.buf)
		} else {
			// BUILT-IN FORMATTER
			// using dFmt (of formatd)
			e.buf = dFmt.addBegin(e.buf)

			// APPEND TIME
			if e.logger.Flag&fUseTime != 0 {
				t := time.Now()
				if (FtimeUnix|FtimeUnixMs)&e.logger.Flag != 0 {
					e.buf = dFmt.addKeyUnsafe(e.buf, "ts")
					if FtimeUnixMs&e.logger.Flag != 0 {
						e.buf = dFmt.addTimeUnix(e.buf, t.UnixNano()/1e6)
					} else {
						e.buf = dFmt.addTimeUnix(e.buf, t.Unix())
					}
				} else {
					if FUTC&e.logger.Flag != 0 {
						t = t.UTC()
					}
					if Fdate&e.logger.Flag != 0 {
						e.buf = dFmt.addKeyUnsafe(e.buf, "date")
						y, m, d := t.Date()
						e.buf = dFmt.addTimeDate(e.buf, y, int(m), d)
					}
					if FdateDay&e.logger.Flag != 0 {
						e.buf = dFmt.addKeyUnsafe(e.buf, "day")
						e.buf = dFmt.addTimeDay(e.buf, int(t.Weekday()))
					}
					if (Ftime|FtimeMs)&e.logger.Flag != 0 {
						e.buf = dFmt.addKeyUnsafe(e.buf, "time")
						h, m, s := t.Clock()
						if FtimeMs&e.logger.Flag != 0 {
							e.buf = dFmt.addTimeMs(e.buf, h, m, s, t.Nanosecond())
						} else {
							e.buf = dFmt.addTime(e.buf, h, m, s)
						}
					}
				}
			}

			// APPEND LEVEL
			if e.logger.Flag&Flevel != 0 {
				e.buf = dFmt.addKeyUnsafe(e.buf, "level")
				e.buf = dFmt.addLevel(e.buf, e.level)
			}

			// APPEND TAG
			if e.logger.Flag&Ftag != 0 {
				e.buf = dFmt.addKeyUnsafe(e.buf, "tag")
				e.buf = dFmt.addTag(e.buf, e.logger.Control.TagBucket, e.tag)
			}

			// APPEND MSG
			if s != "" {
				e.buf = dFmt.addKeyUnsafe(e.buf, "message")
				if ok, _ := dFmt.isSimpleStr(s); ok {
					e.buf = dFmt.addValStringUnsafe(e.buf, s)
				} else {
					e.buf = dFmt.addValString(e.buf, s)
				}
			}

			// APPEND KEY VALUES
			for i := 0; i < len(e.kvs); i++ {
				// Set name
				e.buf = dFmt.addKeyUnsafe(e.buf, e.kvs[i].Key)
				switch e.kvs[i].Vtype {
				case KvInt:
					e.buf = dFmt.addValInt(e.buf, e.kvs[i].Vint)
				case KvString:
					if ok, _ := dFmt.isSimpleStr(e.kvs[i].Vstr); ok {
						e.buf = dFmt.addValStringUnsafe(e.buf, e.kvs[i].Vstr)
					} else {
						e.buf = dFmt.addValString(e.buf, e.kvs[i].Vstr)
					}
				case KvBool:
					e.buf = dFmt.addValBool(e.buf, e.kvs[i].Vbool)
				case KvFloat64:
					e.buf = dFmt.addValFloat(e.buf, e.kvs[i].Vf64)
				case KvError:
					if e.kvs[i].Verr != nil {
						errStr := e.kvs[i].Verr.Error()
						if ok, _ := dFmt.isSimpleStr(errStr); ok {
							e.buf = dFmt.addValStringUnsafe(e.buf, errStr)
						} else {
							e.buf = dFmt.addValString(e.buf, errStr)
						}
					} else {
						e.buf = append(e.buf, `null,`...)
					}
				default:
					e.buf = append(e.buf, `null,`...)
				}
			}

			// APPEND FINAL
			e.buf = dFmt.addEnd(e.buf)

			// Write to output
			if e.logger.outd != nil {
				e.logger.outd.Write(e.buf)
			}
		}
	}
}

// Bool adds KeyValue of boolean into kvs slice.
func (e *entry) Bool(key string, val bool) *entry {
	if e != nil {
		e.kvs = append(e.kvs, KeyValue{
			Vtype: KvBool,
			Key:   key,
			Vbool: val,
		})
	}
	return e
}

// Float adds KeyValue of float64 into kvs slice.
// To minimize the size of entry, alog only supports float64.
func (e *entry) Float(key string, val float64) *entry {
	if e != nil {
		e.kvs = append(e.kvs, KeyValue{
			Vtype: KvFloat64,
			Key:   key,
			Vf64:  val,
		})
	}
	return e
}

// Str adds KeyValue item of a string into kvs slice.
func (e *entry) Str(key string, val string) *entry {
	if e != nil {
		e.kvs = append(e.kvs, KeyValue{
			Vtype: KvString,
			Key:   key,
			Vstr:  val,
		})
	}
	return e
}

// Int adds KeyValue item for integer. This will convert int to int64.
// As both int and int64 are widely used, Alog has both Int and Int64
// but they share same kind (Vint).
func (e *entry) Int(key string, val int) *entry {
	if e != nil {
		e.kvs = append(e.kvs, KeyValue{
			Vtype: KvInt,
			Key:   key,
			Vint:  int64(val),
		})
	}
	return e
}

// Int64 adds KeyValue item for 64 bit integer.
// As both int and int64 are widely used, Alog has both Int and Int64
// but they share same kind (Vint).
func (e *entry) Int64(key string, val int64) *entry {
	if e != nil {
		e.kvs = append(e.kvs, KeyValue{
			Vtype: KvInt,
			Key:   key,
			Vint:  val,
		})
	}
	return e
}

// Err adds KeyValue for error item.
func (e *entry) Err(key string, val error) *entry {
	if e != nil {
		e.kvs = append(e.kvs, KeyValue{
			Vtype: KvError,
			Key:   key,
			Verr:  val,
		})
	}
	return e
}

// Write will call *Logger.Writes() without param.
// As Writes without parameter can be used frequently,
// Alog has both Write() and Writes(string).
func (e *entry) Write() {
	e.Writes("")
}
