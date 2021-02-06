package alog

import (
	"strconv"
	"time"
)

// fmtText is a alog formatter that will print the result in text Format.
type fmtText struct {
	FormatFlag   *Format
	CustomHeader func([]byte) []byte
}

func (f fmtText) LogCustomHeader(dst []byte) []byte {
	if f.CustomHeader != nil {
		return f.CustomHeader(dst)
	}
	return dst
}

func (f fmtText) Begin(dst []byte, prefix []byte) []byte {
	dst = dst[:0] // reset first
	if prefix != nil {
		dst = append(dst, prefix...) // prefix not to be escaped
	}
	return dst
}

func (f fmtText) End(dst []byte) []byte {
	return append(dst, '\n')
}
func (f fmtText) Space(dst []byte) []byte {
	return append(dst, ' ')
}

// Log specific type
func (f fmtText) LogLevel(dst []byte, lv Level) []byte {
	return append(dst, lv.ShortName()...)
}

func (f fmtText) LogTag(dst []byte, tag Tag, alogTagStr *[64]string, alogTagIssued int) []byte {
	dst = append(dst, '[')

	firstItem := true

	for i := 0; i < alogTagIssued; i++ {
		if tag&(1<<i) != 0 {
			if firstItem {
				firstItem = false
			} else {
				dst = append(dst, ',')
			}
			dst = append(dst, alogTagStr[i]...)
		}
	}

	return append(dst, ']')
}
func (f fmtText) LogMsg(dst []byte, s string, suffix byte) []byte {
	dst = f.escape(dst, s, false)
	if suffix != 0 {
		dst = append(dst, suffix)
	}
	return dst
}

func (f fmtText) LogMsgb(dst []byte, b []byte, suffix byte) []byte {
	dst = f.escapeb(dst, b, false)
	if suffix != 0 {
		dst = append(dst, suffix)
	}
	return dst
}

func (f fmtText) LogTime(dst []byte, t time.Time) []byte {
	// "t":  time shows up to millisecond: 3_04_05_000 = h:3, m:4, s:5, ms: 000
	h, m, s := t.Clock()
	dst = itoa(dst, h*10000+m*100+s, 6, '.')
	return itoa(dst, t.Nanosecond()/1e6, 3, 0)
}
func (f fmtText) LogTimeDate(dst []byte, t time.Time) []byte {
	y, m, d := t.Date()
	return itoa(dst, y*10000+int(m)*100+d, 8, 0)
}
func (f fmtText) LogTimeDay(dst []byte, t time.Time) []byte {
	return append(dst, t.Weekday().String()[0:3]...)
}

func (f fmtText) LogTimeUnix(dst []byte, t time.Time) []byte {
	// "ts": unix second
	return itoa(dst, int(t.Unix()), 8, 0)
}
func (f fmtText) LogTimeUnixMs(dst []byte, t time.Time) []byte {
	// "ts": unix second
	return itoa(dst, int(t.UnixNano()/1e6), 8, 0)
}

// Special type
func (f fmtText) Nil(dst []byte, k string) []byte {
	dst = f.escape(dst, k, false)
	return append(dst, `=null`...)
}

func (f fmtText) Error(dst []byte, k string, v error) []byte {
	dst = f.addKey(dst, k)
	if v != nil {
		return f.String(dst, k, v.Error())
	} else {
		return f.Nil(dst, k)
	}
}
func (f fmtText) Errors(dst []byte, k string, v *[]error) []byte {
	dst = f.addKey(dst, k)

	idxv := len(*v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range *v {
		if v2 != nil {
			dst = f.escape(dst, v2.Error(), true)
		} else {
			// dst = append(dst, '"', '"')
			dst = append(dst, "null"...) // todo: check if this is acceptable (null in string array)
		}
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}

// Basic data type
// byte and rune are skipped
func (f fmtText) Bool(dst []byte, k string, v bool) []byte {
	dst = f.addKey(dst, k)
	if v {
		return append(dst, "true"...)
	}
	return append(dst, "false"...)
}
func (f fmtText) String(dst []byte, k string, v string) []byte {
	dst = f.addKey(dst, k)
	return f.escape(dst, v, true)
}
func (f fmtText) Int(dst []byte, k string, v int) []byte {
	dst = f.addKey(dst, k)
	return itoa(dst, v, 1, 0)
}
func (f fmtText) Int8(dst []byte, k string, v int8) []byte {
	dst = f.addKey(dst, k)
	return itoa(dst, int(v), 1, 0)
}
func (f fmtText) Int16(dst []byte, k string, v int16) []byte {
	dst = f.addKey(dst, k)
	return itoa(dst, int(v), 1, 0)
}
func (f fmtText) Int32(dst []byte, k string, v int32) []byte {
	dst = f.addKey(dst, k)
	return itoa(dst, int(v), 1, 0)
}
func (f fmtText) Int64(dst []byte, k string, v int64) []byte {
	dst = f.addKey(dst, k)
	return strconv.AppendInt(dst, v, 10)
}
func (f fmtText) Uint(dst []byte, k string, v uint) []byte {
	dst = f.addKey(dst, k)
	return strconv.AppendUint(dst, uint64(v), 10)
}
func (f fmtText) Uint8(dst []byte, k string, v uint8) []byte {
	dst = f.addKey(dst, k)
	return itoa(dst, int(v), 1, 0)
}
func (f fmtText) Uint16(dst []byte, k string, v uint16) []byte {
	dst = f.addKey(dst, k)
	return itoa(dst, int(v), 1, 0)
}
func (f fmtText) Uint32(dst []byte, k string, v uint32) []byte {
	dst = f.addKey(dst, k)
	return strconv.AppendUint(dst, uint64(v), 10)
}
func (f fmtText) Uint64(dst []byte, k string, v uint64) []byte {
	dst = f.addKey(dst, k)
	return strconv.AppendUint(dst, v, 10)
}
func (f fmtText) Float32(dst []byte, k string, v float32) []byte {
	dst = f.addKey(dst, k)
	return ftoa(dst, float64(v), 2)
}
func (f fmtText) Float64(dst []byte, k string, v float64) []byte {
	dst = f.addKey(dst, k)
	return ftoa(dst, v, 2)
}

// Slice of basic data type
func (f fmtText) Bools(dst []byte, k string, v *[]bool) []byte {
	dst = f.addKey(dst, k)
	idxv := len(*v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range *v {
		if v2 {
			dst = append(dst, "true"...)
		} else {
			dst = append(dst, "false"...)
		}
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}
func (f fmtText) Strings(dst []byte, k string, v *[]string) []byte {
	dst = f.addKey(dst, k)
	idxv := len(*v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v := range *v {
		dst = f.escape(dst, v, true)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}
func (f fmtText) Ints(dst []byte, k string, v *[]int) []byte {
	dst = f.addKey(dst, k)
	idxv := len(*v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range *v {
		dst = itoa(dst, v2, 1, 0)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}
func (f fmtText) Int32s(dst []byte, k string, v *[]int32) []byte {
	dst = f.addKey(dst, k)
	idxv := len(*v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range *v {
		dst = itoa(dst, int(v2), 1, 0)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}
func (f fmtText) Int64s(dst []byte, k string, v *[]int64) []byte {
	dst = f.addKey(dst, k)
	idxv := len(*v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range *v {
		dst = strconv.AppendInt(dst, v2, 10)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}
func (f fmtText) Uints(dst []byte, k string, v *[]uint) []byte {
	dst = f.addKey(dst, k)
	idxv := len(*v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range *v {
		dst = strconv.AppendUint(dst, uint64(v2), 10)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}
func (f fmtText) Uint8s(dst []byte, k string, v *[]uint8) []byte {
	dst = f.addKey(dst, k)
	idxv := len(*v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range *v {
		dst = itoa(dst, int(v2), 1, 0)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}
func (f fmtText) Uint32s(dst []byte, k string, v *[]uint32) []byte {
	dst = f.addKey(dst, k)
	idxv := len(*v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range *v {
		dst = itoa(dst, int(v2), 1, 0)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}
func (f fmtText) Uint64s(dst []byte, k string, v *[]uint64) []byte {
	dst = f.addKey(dst, k)
	idxv := len(*v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range *v {
		dst = strconv.AppendUint(dst, v2, 10)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}
func (f fmtText) Float32s(dst []byte, k string, v *[]float32) []byte {
	dst = f.addKey(dst, k)
	idxv := len(*v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range *v {
		// dst = strconv.AppendUint(dst, v2, 10)
		dst = ftoa(dst, float64(v2), 2)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}
func (f fmtText) Float64s(dst []byte, k string, v *[]float64) []byte {
	dst = f.addKey(dst, k)
	idxv := len(*v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range *v {
		// dst = strconv.AppendUint(dst, v2, 10)
		dst = ftoa(dst, v2, 2)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}

func (f fmtText) addKey(dst []byte, s string) []byte {
	dst = f.escape(dst, s, false)
	return append(dst, '=')
}

func (f fmtText) escape(dst []byte, s string, addQuote bool) []byte {
	if addQuote {
		dst = append(dst, '"')
	}
	for j := 0; j < len(s); j++ {
		switch s[j] {
		case '\n':
			dst = append(dst, ';')
		case '\r':
			// ignore
		default:
			dst = append(dst, s[j])
		}
	}
	if addQuote {
		dst = append(dst, '"')
	}
	return dst
}
func (f fmtText) escapeb(dst []byte, b []byte, addQuote bool) []byte {
	if addQuote {
		dst = append(dst, '"')
	}
	for j := 0; j < len(b); j++ {
		switch b[j] {
		case '\n':
			dst = append(dst, ';')
		case '\r':
			// ignore
		default:
			dst = append(dst, b[j])
		}
	}
	if addQuote {
		dst = append(dst, '"')
	}
	return dst
}
