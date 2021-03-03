package alog

import (
	"strconv"
)

type Formatter interface {
}

type formatd [256]bool

func (f *formatd) init() {
	for i := 0; i <= 0x7e; i++ {
		f[i] = i >= 0x20 && i != '\\' && i != '"' // all printable will be true
	}
}
func (f *formatd) isSimpleStr(s string) (bool, int) {
	for i := 0; i < len(s); i++ {
		if f[s[i]] == false {
			return false, i
		}
	}
	return true, 0
}
func (formatd) addBegin(dst []byte) []byte {
	return append(dst, '{')
}
func (formatd) addEnd(dst []byte) []byte {
	dst[len(dst)-1] = '}'
	return append(dst, '\n')
}
func (formatd) addKey(dst []byte, s string) []byte {
	return append(strconv.AppendQuote(dst, s), ':')
}
func (formatd) addKeyUnsafe(dst []byte, s string) []byte {
	return append(append(append(dst, '"'), s...), '"', ':')
}
func (formatd) addValString(dst []byte, s string) []byte {
	return append(strconv.AppendQuote(dst, s), ',')
}
func (formatd) addValStringUnsafe(dst []byte, s string) []byte {
	return append(append(append(dst, '"'), s...), '"', ',')
}
func (formatd) addValBool(dst []byte, b bool) []byte {
	if b {
		return append(dst, `true,`...)
	}
	return append(dst, `false,`...)
}
func (formatd) addValInt(dst []byte, i int64) []byte {
	return append(strconv.AppendInt(dst, i, 10), ',')
}
func (formatd) addValFloat(dst []byte, f float64) []byte {
	return append(strconv.AppendFloat(dst, f, 'f', -1, 64), ',')
}
func (formatd) addTag(dst []byte, bucket *TagBucket, tag Tag) []byte {
	dst = bucket.AppendTagForJSON(append(dst, '['), tag)
	return append(dst, ']', ',')
}
func (formatd) addLevel(dst []byte, level Level) []byte {
	return append(append(append(dst, '"'), level.Name()...), '"', ',')
}
func (formatd) addTimeUnix(dst []byte, ts int64) []byte {
	return append(strconv.AppendInt(dst, ts, 10), ',')
}
func (formatd) addTimeDate(dst []byte, y, m, d int) []byte {
	return append(strconv.AppendInt(dst, int64(y*10000+int(m)*100+d), 10), ',')
}
func (formatd) addTime(dst []byte, h, m, s int) []byte {
	return append(strconv.AppendInt(dst, int64(h*10000+m*100+s), 10), ',')
}
func (formatd) addTimeMs(dst []byte, h, m, s, ns int) []byte {
	return append(strconv.AppendInt(dst, int64(h*10000000+m*100000+s*1000+ns/1e6), 10), ',')
}
func (formatd) addTimeDay(dst []byte, weekday int) []byte {
	return append(strconv.AppendInt(dst, int64(weekday), 10), ',')
}
