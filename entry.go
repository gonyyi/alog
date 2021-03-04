package alog

import (
	"sync"
	"time"
)

const (
	entry_buf_size = 512
	entry_kv_size  = 10
)

func newEntryPoolItem() interface{} {
	return &entry{
		buf: make([]byte, entry_buf_size),
		kvs: make([]KeyValue, entry_kv_size),
	}
}

func newEntryPool() entryPool {
	return entryPool{
		pool: sync.Pool{
			New: newEntryPoolItem,
		},
	}
}

// pool is an a Buffer implementation of sync.Pool.
type entryPool struct {
	pool sync.Pool
}

func (p *entryPool) Get(logger *Logger) *entry {
	b := p.pool.Get().(*entry)
	b.logger = logger
	//b.format = logger.Format
	//b.out = logger.out
	return b
}

func (p *entryPool) Put(b *entry) {
	b.buf = b.buf[:entry_buf_size]
	b.kvs = b.kvs[:entry_kv_size]
	p.pool.Put(b)
}

// entry is the main entry format used in Alog.
// This can be used as standalone or within the sync.Pool
type entry struct {
	buf    []byte
	level  Level
	tag    Tag
	logger *Logger
	//format Format
	//out    io.Writer
	kvs []KeyValue
}

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

func (e *entry) Write(s string) {
	if e != nil {
		defer e.logger.pool.Put(e)

		if e.logger.CusFmat != nil {
			e.buf = e.logger.CusFmat.Begin(e.buf)
			//e.buf = e.logger.CusFmat.AddTime(e.buf, e.logger.Format)
			e.buf = e.logger.CusFmat.AddLevel(e.buf, e.level)
			e.buf = e.logger.CusFmat.AddTag(e.buf, e.tag, e.logger.Control.Tags)
			e.buf = e.logger.CusFmat.AddMsg(e.buf, s)
			e.buf = e.logger.CusFmat.AddKvs(e.buf, e.kvs)
			e.buf = e.logger.CusFmat.End(e.buf)
		} else {
			// BUILTIN PRINTER
			e.buf = dFmt.addBegin(e.buf)
			// INTERFACE: AppendTime()
			if e.logger.Format&fUseTime != 0 {
				t := time.Now()
				if (FtimeUnix|FtimeUnixMs)&e.logger.Format != 0 {
					e.buf = dFmt.addKeyUnsafe(e.buf, "ts")
					if FtimeUnixMs&e.logger.Format != 0 {
						e.buf = dFmt.addTimeUnix(e.buf, t.UnixNano()/1e6)
					} else {
						e.buf = dFmt.addTimeUnix(e.buf, t.Unix())
					}
				} else {
					if FUTC&e.logger.Format != 0 {
						t = t.UTC()
					}
					if Fdate&e.logger.Format != 0 {
						e.buf = dFmt.addKeyUnsafe(e.buf, "date")
						y, m, d := t.Date()
						e.buf = dFmt.addTimeDate(e.buf, y, int(m), d)
					}
					if FdateDay&e.logger.Format != 0 {
						e.buf = dFmt.addKeyUnsafe(e.buf, "day")
						e.buf = dFmt.addTimeDay(e.buf, int(t.Weekday()))
					}
					if (Ftime|FtimeMs)&e.logger.Format != 0 {
						e.buf = dFmt.addKeyUnsafe(e.buf, "time")
						h, m, s := t.Clock()
						if FtimeMs&e.logger.Format != 0 {
							e.buf = dFmt.addTimeMs(e.buf, h, m, s, t.Nanosecond())
						} else {
							e.buf = dFmt.addTime(e.buf, h, m, s)
						}
					}
				}
			}

			// INTERFACE: LEVEL
			if e.logger.Format&Flevel != 0 {
				e.buf = dFmt.addKeyUnsafe(e.buf, "level")
				e.buf = dFmt.addLevel(e.buf, e.level)
			}

			// INTERFACE: TAG
			if e.logger.Format&Ftag != 0 {
				e.buf = dFmt.addKeyUnsafe(e.buf, "tag")
				e.buf = dFmt.addTag(e.buf, e.logger.Control.Tags, e.tag)
			}

			// INTERFACE: MSG
			e.buf = dFmt.addKeyUnsafe(e.buf, "message")
			if ok, _ := dFmt.isSimpleStr(s); ok {
				e.buf = dFmt.addValStringUnsafe(e.buf, s)
			} else {
				e.buf = dFmt.addValString(e.buf, s)
			}

			// INTERFACE: ADD kvs
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
						e.buf = dFmt.addValString(e.buf, e.kvs[i].Verr.Error())
					} else {
						e.buf = append(e.buf, "null,"...)
					}
				default:
					e.buf = append(e.buf, "null,"...)
				}
			}
			// INTERFACE: FINALIZE
			e.buf = dFmt.addEnd(e.buf)
		}

		if e.logger.out != nil {
			e.logger.out.Write(e.buf)
		}
	}
}
