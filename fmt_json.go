package alog

import (
	"strconv"
	"time"
)

// fmtJSON is a alog formatter that will print the result in JSON Format.
type fmtJSON struct {
	FormatFlag   *Format
	CustomHeader func([]byte) []byte
}

func (f fmtJSON) LogCustomHeader(dst []byte) []byte {
	if f.CustomHeader != nil {
		return f.CustomHeader(dst)
	}
	return dst
}

func (f fmtJSON) Begin(dst []byte, prefix []byte) []byte {
	dst = dst[:0] // reset first
	if prefix != nil {
		dst = append(dst, prefix...) // prefix not to be escaped
	}
	return append(dst, '{')
}

func (f fmtJSON) End(dst []byte) []byte {
	return append(dst, '}', '\n')
}

func (f fmtJSON) Space(dst []byte) []byte {
	return append(dst, ',')
}

// Log specific type
func (f fmtJSON) LogLevel(dst []byte, lv Level) []byte {
	return f.safeString(dst, "level", lv.String())
}

func (f fmtJSON) LogTag(dst []byte, tag Tag, alogTagStr *[64]string, alogTagIssued int) []byte {
	dst = append(dst, `"tag":[`...)

	firstItem := true

	for i := 0; i < alogTagIssued; i++ {
		if tag&(1<<i) != 0 {
			if firstItem {
				firstItem = false
				dst = append(dst, '"')
			} else {
				dst = append(dst, ',', '"')
			}
			dst = append(dst, alogTagStr[i]...)
			dst = append(dst, '"')
		}
	}
	return append(dst, ']')
}
func (f fmtJSON) LogMsg(dst []byte, s string, suffix byte) []byte { // suffix to be used for text version only
	// For JSON suffix won't be applied.
	return f.String(dst, "msg", s) // faster without addKey
}
func (f fmtJSON) LogMsgb(dst []byte, b []byte, suffix byte) []byte {
	dst = f.addKey(dst, "msg")
	return f.escapeb(dst, b, true)
}

func (f fmtJSON) LogTime(dst []byte, t time.Time) []byte {
	// "t":  time shows up to millisecond: 3_04_05_000 = h:3, m:4, s:5, ms: 000
	dst = append(dst, `"t":`...) // faster without addKey
	h, m, s := t.Clock()
	dst = itoa(dst, h*10000+m*100+s, 1, 0)
	return itoa(dst, t.Nanosecond()/1e6, 3, 0)
}
func (f fmtJSON) LogTimeDate(dst []byte, t time.Time) []byte {
	dst = append(dst, `"d":`...) // faster without addKey
	y, m, d := t.Date()
	return itoa(dst, y*10000+int(m)*100+d, 4, 0)
}
func (f fmtJSON) LogTimeDay(dst []byte, t time.Time) []byte {
	// "wd": 0 being sunday, 6 being saturday
	dst = append(dst, `"wd":`...) // faster without addKey
	dst = itoa(dst, int(t.Weekday()), 1, 0)
	return dst
}
func (f fmtJSON) LogTimeUnix(dst []byte, t time.Time) []byte {
	// "ts": unix second
	dst = append(dst, `"ts":`...) // faster without addKey
	return itoa(dst, int(t.Unix()), 8, 0)
}
func (f fmtJSON) LogTimeUnixMs(dst []byte, t time.Time) []byte {
	// "ts": unix second
	dst = append(dst, `"ts":`...) // faster without addKey
	return itoa(dst, int(t.UnixNano())/1e6, 8, 0)
}

// Special type
func (f fmtJSON) Nil(dst []byte, k string) []byte {
	dst = f.escape(dst, k, true) // faster without addKey
	return append(dst, `:null`...)
}

func (f fmtJSON) Error(dst []byte, k string, v error) []byte {
	dst = f.addKey(dst, k)
	if v != nil {
		return f.String(dst, k, v.Error())
	} else {
		return f.Nil(dst, k)
	}
}

func (f fmtJSON) Errors(dst []byte, k string, v *[]error) []byte {
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
func (f fmtJSON) Bool(dst []byte, k string, v bool) []byte {
	dst = f.addKey(dst, k)
	if v {
		return append(dst, "true"...)
	}
	return append(dst, "false"...)
}

func (f fmtJSON) String(dst []byte, k string, v string) []byte {
	dst = f.addKey(dst, k)
	return f.escape(dst, v, true)
}

func (f fmtJSON) Int(dst []byte, k string, v int) []byte {
	dst = f.addKey(dst, k)
	return itoa(dst, v, 1, 0)
}
func (f fmtJSON) Int8(dst []byte, k string, v int8) []byte {
	dst = f.addKey(dst, k)
	return itoa(dst, int(v), 1, 0)
}
func (f fmtJSON) Int16(dst []byte, k string, v int16) []byte {
	dst = f.addKey(dst, k)
	return itoa(dst, int(v), 1, 0)
}
func (f fmtJSON) Int32(dst []byte, k string, v int32) []byte {
	dst = f.addKey(dst, k)
	return itoa(dst, int(v), 1, 0)
}
func (f fmtJSON) Int64(dst []byte, k string, v int64) []byte {
	dst = f.addKey(dst, k)
	return strconv.AppendInt(dst, v, 10)
}
func (f fmtJSON) Uint(dst []byte, k string, v uint) []byte {
	dst = f.addKey(dst, k)
	return itoa(dst, int(v), 1, 0)
}
func (f fmtJSON) Uint8(dst []byte, k string, v uint8) []byte {
	dst = f.addKey(dst, k)
	return itoa(dst, int(v), 1, 0)
}
func (f fmtJSON) Uint16(dst []byte, k string, v uint16) []byte {
	dst = f.addKey(dst, k)
	return itoa(dst, int(v), 1, 0)
}
func (f fmtJSON) Uint32(dst []byte, k string, v uint32) []byte {
	dst = f.addKey(dst, k)
	return strconv.AppendUint(dst, uint64(v), 10)
}
func (f fmtJSON) Uint64(dst []byte, k string, v uint64) []byte {
	dst = f.addKey(dst, k)
	return strconv.AppendUint(dst, v, 10)
}
func (f fmtJSON) Float32(dst []byte, k string, v float32) []byte {
	dst = f.addKey(dst, k)
	return ftoa(dst, float64(v), 2)
}
func (f fmtJSON) Float64(dst []byte, k string, v float64) []byte {
	dst = f.addKey(dst, k)
	return ftoa(dst, v, 2)
}

// Slice of basic data type
func (f fmtJSON) Bools(dst []byte, k string, v *[]bool) []byte {
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
func (f fmtJSON) Strings(dst []byte, k string, v *[]string) []byte {
	dst = f.addKey(dst, k)
	idxv := len(*v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range *v {
		dst = f.escape(dst, v2, true)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}
func (f fmtJSON) Ints(dst []byte, k string, v *[]int) []byte {
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
func (f fmtJSON) Int32s(dst []byte, k string, v *[]int32) []byte {
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
func (f fmtJSON) Int64s(dst []byte, k string, v *[]int64) []byte {
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
func (f fmtJSON) Uints(dst []byte, k string, v *[]uint) []byte {
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
func (f fmtJSON) Uint8s(dst []byte, k string, v *[]uint8) []byte {
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
func (f fmtJSON) Uint32s(dst []byte, k string, v *[]uint32) []byte {
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
func (f fmtJSON) Uint64s(dst []byte, k string, v *[]uint64) []byte {
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
func (f fmtJSON) Float32s(dst []byte, k string, v *[]float32) []byte {
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
func (f fmtJSON) Float64s(dst []byte, k string, v *[]float64) []byte {
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

func (f fmtJSON) addKey(dst []byte, s string) []byte {
	dst = f.escape(dst, s, true)
	return append(dst, ':')
}

func (f fmtJSON) safeString(dst []byte, k string, v string) []byte {
	dst = append(dst, '"')
	dst = append(dst, k...)
	dst = append(dst, `":"`...)
	dst = append(dst, v...)
	return append(dst, '"')
}

func (f fmtJSON) escape(dst []byte, s string, addQuote bool) []byte {
	if addQuote {
		dst = append(dst, '"')
	}
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"', '\\':
			dst = append(dst, '\\', s[i])
		case '\n':
			dst = append(dst, '\\', 'n')
		case '\t':
			dst = append(dst, '\\', 't')
		case '\r':
			dst = append(dst, '\\', 'r')
		case '\b':
			dst = append(dst, '\\', 'b')
		case '\f':
			dst = append(dst, '\\', 'f')
		default:
			dst = append(dst, s[i])
		}
	}
	if addQuote {
		dst = append(dst, '"')
	}
	return dst
}
func (f fmtJSON) escapeb(dst []byte, b []byte, addQuote bool) []byte {
	if addQuote {
		dst = append(dst, '"')
	}
	for i := 0; i < len(b); i++ {
		switch b[i] {
		case '"', '\\':
			dst = append(dst, '\\', b[i])
		case '\n':
			dst = append(dst, '\\', 'n')
		case '\t':
			dst = append(dst, '\\', 't')
		case '\r':
			dst = append(dst, '\\', 'r')
		case '\b':
			dst = append(dst, '\\', 'b')
		case '\f':
			dst = append(dst, '\\', 'f')
		default:
			dst = append(dst, b[i])
		}
	}
	if addQuote {
		dst = append(dst, '"')
	}
	return dst
}
