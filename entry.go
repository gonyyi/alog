package alog

import (
	"time"
)

// KvType hold diff
type KvType uint8

// KeyValue holds Key and value info
type KeyValue struct {
	Key   string
	Vtype KvType
	Vint  int64
	Vf64  float64
	Vstr  string
	Vbool bool
	Verr  error
}

// entryInfo is 56 bytes
type entryInfo struct {
	flag    Flag
	tbucket *TagBucket
	w       Writer
	orFmtr  Formatter
	//w       io.Writer
}

// Entry is a log Entry will be used with a entryPool to
// reuse the resource.
// Entry is 120 bytes, use pointer.
type Entry struct {
	buf   []byte
	level Level
	tag   Tag
	kvs   []KeyValue
	info  entryInfo
}

// Writes will finalize the log message, format it, and
// write it to writer. Besides *logger.getEntry(), this is
// the only other method which isn't inline-able.
func (e *Entry) Write(msg string) {
	// When log message was created from *Logger.getEntry(),
	// it examines logability (should log or not). Once it's not eligible,
	// it will return nil.
	// eg. log.Trace(0).Str("name", "gon").Int("age", 39).Write()
	//   1. Trace(0) will call getEntry(), if it determines the log shouldn't be
	//      logged, it will return nil.
	//   2. Next method with Str() receives nil for the pointer and will ignore,
	//      and return nil to next.
	//   3. Int method will receive nil, and just pass nil to next.
	//   4. Write method finally receives it, if Entry pointer is nil, it won't
	//      do anything as it's not eligible to log.

	// if Entry is not nil (=loggable),
	if e != nil {
		// since pointer receiver *Entry is obtained from the pool,
		// make sure this will be put back to memory.
		defer pool.Put(e)

		// if custom formatter exists, use it instead of default formatter.
		// for default formatter (formatd), it's a concrete function for speed.
		// rather than using from the interface.
		if e.info.orFmtr != nil {
			// CUSTOM FORMATTER
			e.buf = e.info.orFmtr.Begin(e.buf)
			e.buf = e.info.orFmtr.AddTime(e.buf)
			e.buf = e.info.orFmtr.AddLevel(e.buf, e.level)
			e.buf = e.info.orFmtr.AddTag(e.buf, e.tag)
			e.buf = e.info.orFmtr.AddMsg(e.buf, msg)
			e.buf = e.info.orFmtr.AddKVs(e.buf, e.kvs)
			e.buf = e.info.orFmtr.End(e.buf)
			//e.info.orFmtr.Write(e.buf)
			e.info.orFmtr.Write(e.buf, e.level, e.tag)
		} else {
			// BUILT-IN FORMATTER
			// using dFmt (of formatd)
			e.buf = dFmt.addBegin(e.buf)

			// APPEND TIME
			if e.info.flag&fHasTime != 0 {
				t := time.Now()
				if (WithUnixTime|WithUnixTimeMs)&e.info.flag != 0 {
					e.buf = dFmt.addKeyUnsafe(e.buf, "ts")
					if WithUnixTimeMs&e.info.flag != 0 {
						e.buf = dFmt.addTimeUnix(e.buf, t.UnixNano()/1e6)
					} else {
						e.buf = dFmt.addTimeUnix(e.buf, t.Unix())
					}
				} else {
					if WithUTC&e.info.flag != 0 {
						t = t.UTC()
					}
					if WithDate&e.info.flag != 0 {
						e.buf = dFmt.addKeyUnsafe(e.buf, "date")
						y, m, d := t.Date()
						e.buf = dFmt.addTimeDate(e.buf, y, int(m), d)
					}
					if WithDay&e.info.flag != 0 {
						e.buf = dFmt.addKeyUnsafe(e.buf, "day")
						e.buf = dFmt.addTimeDay(e.buf, int(t.Weekday()))
					}
					if (WithTime|WithTimeMs)&e.info.flag != 0 {
						e.buf = dFmt.addKeyUnsafe(e.buf, "time")
						h, m, s := t.Clock()
						if WithTimeMs&e.info.flag != 0 {
							e.buf = dFmt.addTimeMs(e.buf, h, m, s, t.Nanosecond())
						} else {
							e.buf = dFmt.addTime(e.buf, h, m, s)
						}
					}
				}
			}

			// APPEND LEVEL
			if e.info.flag&WithLevel != 0 {
				e.buf = dFmt.addKeyUnsafe(e.buf, "level")
				e.buf = dFmt.addLevel(e.buf, e.level)
			}

			// APPEND TAG
			if e.info.flag&WithTag != 0 {
				e.buf = dFmt.addKeyUnsafe(e.buf, "tag")
				e.buf = dFmt.addTag(e.buf, e.info.tbucket, e.tag)
			}

			// APPEND MSG
			if msg != "" {
				e.buf = dFmt.addKeyUnsafe(e.buf, "message")
				if ok, _ := dFmt.isSimpleStr(msg); ok {
					e.buf = dFmt.addValStringUnsafe(e.buf, msg)
				} else {
					e.buf = dFmt.addValString(e.buf, msg)
				}
			}

			// APPEND KEY VALUES
			for i := 0; i < len(e.kvs); i++ {
				// Set name
				if ok, _ := dFmt.isSimpleStr(e.kvs[i].Key); ok {
					e.buf = dFmt.addKeyUnsafe(e.buf, e.kvs[i].Key)
				} else {
					e.buf = dFmt.addKey(e.buf, e.kvs[i].Key)
				}

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
				//e.info.w.Write(e.buf)
				e.info.w.WriteLt(e.buf, e.level, e.tag)
			}
		}
	}
}

// Bool adds KeyValue of boolean into kvs slice.
func (e *Entry) Bool(key string, val bool) *Entry {
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
// To minimize the size of Entry, alog only supports float64.
func (e *Entry) Float(key string, val float64) *Entry {
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
func (e *Entry) Str(key string, val string) *Entry {
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
func (e *Entry) Int(key string, val int) *Entry {
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
func (e *Entry) Int64(key string, val int64) *Entry {
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
func (e *Entry) Err(key string, val error) *Entry {
	if e != nil {
		e.kvs = append(e.kvs, KeyValue{
			Vtype: KvError,
			Key:   key,
			Verr:  val,
		})
	}
	return e
}

// Ext will take EntryFn and add an entry to it.
func (e *Entry) Ext(fn EntryFn) *Entry {
	if e == nil || fn == nil {
		return e
	}
	return fn(e)
}
