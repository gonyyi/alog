package ext

import (
	"github.com/gonyyi/alog"
	"strconv"
	"time"
)

func NewFormatterTerminal() *fmtTxt {
	return &fmtTxt{}
}

type fmtTxt struct{}

func (f fmtTxt) Begin(dst []byte) []byte {
	return dst
}
func (f fmtTxt) AddTime(dst []byte, format alog.Format) []byte {
	return append(append(dst, time.Now().String()...), ' ')
}
func (f fmtTxt) AddLevel(dst []byte, level alog.Level) []byte {
	return append(append(dst, level.NameShort()...), ' ')
}
func (f fmtTxt) AddTag(dst []byte, tag alog.Tag, bucket *alog.TagBucket) []byte {
	return append(bucket.AppendTag(append(dst, '['), tag), ']', ' ')
}
func (f fmtTxt) AddMsg(dst []byte, s string) []byte {
	return append(append(dst, s...), ` // `...)
}
func (f fmtTxt) AddKvs(dst []byte, kvs []alog.KeyValue) []byte {
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
func (f fmtTxt) End(dst []byte) []byte {
	if len(dst) > 1 {
		dst[len(dst)-2] = '\n'
		return dst[:len(dst)-1]
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
