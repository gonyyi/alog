package alog

import (
	"io"
	"time"
)

// entryInfo is 56 bytes
type entryInfo struct {
	flag    Flag
	tbucket *TagBucket
	pool    *entryPool
	orFmtr  Formatter
	w       io.Writer
}

// entry is a log entry will be used with a entryPool to
// reuse the resource.
// entry is 120 bytes, use pointer.
type entry struct {
	buf   []byte
	level Level
	tag   Tag
	kvs   []KeyValue
	info  entryInfo
}

// Writes will finalize the log message, format it, and
// write it to writer. Besides *logger.getEntry(), this is
// the only other method which isn't inline-able.
func (e *entry) Write(s string) {
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
		defer e.info.pool.Put(e)

		// if custom formatter exists, use it instead of default formatter.
		// for default formatter (formatd), it's a concrete function for speed.
		// rather than using from the interface.
		if e.info.orFmtr != nil {
			// CUSTOM FORMATTER
			e.buf = e.info.orFmtr.Begin(e.buf)
			e.buf = e.info.orFmtr.AddTime(e.buf)
			e.buf = e.info.orFmtr.AddLevel(e.buf, e.level)
			e.buf = e.info.orFmtr.AddTag(e.buf, e.tag)
			e.buf = e.info.orFmtr.AddMsg(e.buf, s)
			e.buf = e.info.orFmtr.AddKVs(e.buf, e.kvs)
			e.buf = e.info.orFmtr.End(e.buf)
			e.info.orFmtr.Write(e.buf)
		} else {
			// BUILT-IN FORMATTER
			// using dFmt (of formatd)
			e.buf = dFmt.addBegin(e.buf)

			// APPEND TIME
			if e.info.flag&fUseTime != 0 {
				t := time.Now()
				if (FtimeUnix|FtimeUnixMs)&e.info.flag != 0 {
					e.buf = dFmt.addKeyUnsafe(e.buf, "ts")
					if FtimeUnixMs&e.info.flag != 0 {
						e.buf = dFmt.addTimeUnix(e.buf, t.UnixNano()/1e6)
					} else {
						e.buf = dFmt.addTimeUnix(e.buf, t.Unix())
					}
				} else {
					if FUTC&e.info.flag != 0 {
						t = t.UTC()
					}
					if Fdate&e.info.flag != 0 {
						e.buf = dFmt.addKeyUnsafe(e.buf, "date")
						y, m, d := t.Date()
						e.buf = dFmt.addTimeDate(e.buf, y, int(m), d)
					}
					if FdateDay&e.info.flag != 0 {
						e.buf = dFmt.addKeyUnsafe(e.buf, "day")
						e.buf = dFmt.addTimeDay(e.buf, int(t.Weekday()))
					}
					if (Ftime|FtimeMs)&e.info.flag != 0 {
						e.buf = dFmt.addKeyUnsafe(e.buf, "time")
						h, m, s := t.Clock()
						if FtimeMs&e.info.flag != 0 {
							e.buf = dFmt.addTimeMs(e.buf, h, m, s, t.Nanosecond())
						} else {
							e.buf = dFmt.addTime(e.buf, h, m, s)
						}
					}
				}
			}

			// APPEND LEVEL
			if e.info.flag&Flevel != 0 {
				e.buf = dFmt.addKeyUnsafe(e.buf, "level")
				e.buf = dFmt.addLevel(e.buf, e.level)
			}

			// APPEND TAG
			if e.info.flag&Ftag != 0 {
				e.buf = dFmt.addKeyUnsafe(e.buf, "tag")
				e.buf = dFmt.addTag(e.buf, e.info.tbucket, e.tag)
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
			if e.info.w != nil {
				//e.logger.w.Write(e.buf)
				e.info.w.Write(e.buf)
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
