package ext

import (
	"github.com/gonyyi/alog"
	"io"
	"strconv"
	"time"
)

// NewFormatterTerminal will take io.Writer and returns
// a Logger.Formatter object.
func NewFormatterTerminal() *fmtTxt {
	return &fmtTxt{}
}

type fmtTxt struct {
	//out       io.Writer
	out       alog.Writer
	format    alog.Flag
	tagBucket *alog.TagBucket
}

//func (f *fmtTxt) Init(w io.Writer, formatFlag alog.Flag, tagBucket *alog.TagBucket) {
func (f *fmtTxt) Init(w alog.Writer, formatFlag alog.Flag, tagBucket *alog.TagBucket) {
	f.out = w
	if w == nil {
		f.out = alog.Discard{}
	}

	f.format = formatFlag
	f.tagBucket = tagBucket
}

func (f *fmtTxt) Write(dst []byte, level alog.Level, tag alog.Tag) (int, error) {
	//return f.out.Write(dst)
	return f.out.WriteLt(dst, level, tag)
}

func (f *fmtTxt) Close() error {
	if c, ok := f.out.(io.Closer); ok && c != nil {
		return c.Close()
	}
	return nil
}

func (fmtTxt) Begin(dst []byte) []byte {
	return dst
}

func (f *fmtTxt) AddTime(dst []byte) []byte {
	if (alog.WithUnixTime|alog.WithDate|alog.WithTime|alog.WithTimeMs)&f.format != 0 {
		switch {
		case alog.WithTimeMs&f.format != 0:
			return append(append(dst, time.Now().Format("2006/01/02 15:04:05.000")...), ' ')
		default:
			return append(append(dst, time.Now().Format("2006/01/02 15:04:05")...), ' ')
		}
	}
	return dst
}

func (fmtTxt) AddLevel(dst []byte, level alog.Level) []byte {
	return append(append(dst, level.NameShort()...), ' ')
}

func (f *fmtTxt) AddTag(dst []byte, tag alog.Tag) []byte {
	return append(f.tagBucket.AppendTag(append(dst, '['), tag), ']', ' ')
}

func (fmtTxt) AddMsg(dst []byte, s string) []byte {
	if s != "" {
		return append(dst, s...)
	}
	return dst
}

func (f *fmtTxt) AddKVs(dst []byte, kvs []alog.KeyValue) []byte {
	if len(kvs) > 0 {
		dst = append(dst, ` // `...)
	}

	for i := 0; i < len(kvs); i++ {
		dst = append(append(dst, kvs[i].Key...), '=')
		switch kvs[i].Vtype {
		case alog.KvString:
			dst = f.addValStringUnsafe(dst, kvs[i].Vstr)
		case alog.KvBool:
			dst = f.addValBool(dst, kvs[i].Vbool)
		case alog.KvError:
			if kvs[i].Verr == nil {
				dst = append(dst, "null, "...)
			} else {
				dst = f.addValString(dst, kvs[i].Verr.Error())
			}
		case alog.KvInt:
			dst = f.addValInt(dst, kvs[i].Vint)
		case alog.KvFloat64:
			dst = f.addValFloat(dst, kvs[i].Vf64)
		default:
			dst = append(dst, `null, `...)
		}
	}
	return dst
}

func (fmtTxt) End(dst []byte) []byte {
	if len(dst) > 1 {
		if dst[len(dst)-2] == ' ' || dst[len(dst)-2] == ',' {
			dst[len(dst)-2] = '\n'
			return dst[:len(dst)-1]
		}
		return append(dst, '\n')
	}
	return dst
}

func (fmtTxt) addKeyUnsafe(dst []byte, s string) []byte {
	return append(append(dst, s...), '=')
}

func (fmtTxt) addValString(dst []byte, s string) []byte {
	return append(strconv.AppendQuote(dst, s), ',', ' ')
}

func (fmtTxt) addValStringUnsafe(dst []byte, s string) []byte {
	return append(append(append(dst, '"'), s...), `", `...)
}

func (fmtTxt) addValBool(dst []byte, b bool) []byte {
	if b {
		return append(dst, `true, `...)
	}
	return append(dst, `false, `...)
}

func (fmtTxt) addValInt(dst []byte, i int64) []byte {
	return append(strconv.AppendInt(dst, i, 10), ',', ' ')
}

func (fmtTxt) addValFloat(dst []byte, f float64) []byte {
	return append(strconv.AppendFloat(dst, f, 'f', -1, 64), ',', ' ')
}
