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
		kvs: make([]KeyVal, entry_kv_size),
	}
}

func newEntryPool() entryPool {
	return entryPool{
		pool: sync.Pool{
			New: newEntryPoolItem,
		},
	}
}

// buf is an a Buffer implementation of sync.Pool.
type entryPool struct {
	pool sync.Pool
}

func (p *entryPool) Get(logger *Logger) *entry {
	b := p.pool.Get().(*entry)
	b.logger = logger
	return b
}

func (p *entryPool) Put(b *entry) {
	b.buf = b.buf[:512]
	b.kvs = b.kvs[:10]
	p.pool.Put(b)
}

// entry is the main entry format used in Alog.
// This can be used as standalone or within the sync.Pool
type entry struct {
	buf    []byte
	level  Level
	tag    Tag
	logger *Logger
	kvs    []KeyVal
}

func (e *entry) Bool(key string, val bool) *entry {
	if e != nil {
		e.kvs = append(e.kvs, KeyVal{
			vType: kvBool,
			key:   key,
			vBool: val,
		})
	}
	return e
}
func (e *entry) Float(key string, val float64) *entry {
	if e != nil {
		e.kvs = append(e.kvs, KeyVal{
			vType: kvFloat64,
			key:   key,
			vF64:  val,
		})
	}
	return e
}
func (e *entry) Str(key string, val string) *entry {
	if e != nil {
		e.kvs = append(e.kvs, KeyVal{
			vType: kvString,
			key:   key,
			vStr:  val,
		})
	}
	return e
}
func (e *entry) Int(key string, val int) *entry {
	if e != nil {
		e.kvs = append(e.kvs, KeyVal{
			vType: kvInt,
			key:   key,
			vInt:  int64(val),
		})
	}
	return e
}
func (e *entry) Err(key string, val error) *entry {
	if e != nil {
		e.kvs = append(e.kvs, KeyVal{
			vType: kvError,
			key:   key,
			vErr:  val,
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
			e.buf = e.logger.fmt.addKeyUnsafe(e.buf, e.kvs[i].key)
			switch e.kvs[i].vType {
			case kvInt:
				e.buf = e.logger.fmt.addValInt(e.buf, e.kvs[i].vInt)
			case kvString:
				if ok, _ := e.logger.fmt.isSimpleStr(e.kvs[i].vStr); ok {
					e.buf = e.logger.fmt.addValStringUnsafe(e.buf, e.kvs[i].vStr)
				} else {
					e.buf = e.logger.fmt.addValString(e.buf, e.kvs[i].vStr)
				}
			case kvBool:
				e.buf = e.logger.fmt.addValBool(e.buf, e.kvs[i].vBool)
			case kvFloat64:
				e.buf = e.logger.fmt.addValFloat(e.buf, e.kvs[i].vF64)
			case kvError:
				if e.kvs[i].vErr != nil {
					e.buf = e.logger.fmt.addValString(e.buf, e.kvs[i].vErr.Error())
				} else {
					e.buf = append(e.buf, "null,"...)
				}
			default:
				e.buf = append(e.buf, "null,"...)
			}
		}

		// INTERFACE: FINALIZE
		e.buf = e.logger.fmt.addEnd(e.buf)
		if e.logger.out != nil {
			e.logger.out.Write(e.buf)
		}
	}
}
