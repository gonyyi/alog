package alog

import (
	"strconv"
	"time"
)

func NewFmtrText() Fmtr {
	return &FmtrText{}
}

// FmtrText is a alog formatter that will print the result in text Format.
type FmtrText struct {
	FormatFlag   *Format
	CustomHeader func([]byte) []byte
}

func (f FmtrText) LogCustomHeader(dst []byte) []byte {
	if f.CustomHeader != nil {
		return f.CustomHeader(dst)
	}
	return dst
}

func (f FmtrText) Begin(dst []byte, prefix []byte) []byte {
	dst = dst[:0] // reset first
	if prefix != nil {
		dst = append(dst, prefix...) // prefix not to be escaped
	}
	return dst
}

func (f FmtrText) End(dst []byte) []byte {
	return append(dst, '\n')
}
func (f FmtrText) Space(dst []byte) []byte {
	return append(dst, ' ')
}

// Log specific type
func (f FmtrText) LogLevel(dst []byte, lv Level) []byte {
	return append(dst, lv.ShortName()...)
}

func (f FmtrText) LogTag(dst []byte, tag Tag, alogTagStr [64]string, alogTagIssued int) []byte {
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
func (f FmtrText) LogMsg(dst []byte, s string, suffix byte) []byte {
	dst = f.escape(dst, s, false)
	if suffix != 0 {
		dst = append(dst, suffix)
	}
	return dst
}

func (f FmtrText) LogMsgb(dst []byte, b []byte, suffix byte) []byte {
	dst = f.escapeb(dst, b, false)
	if suffix != 0 {
		dst = append(dst, suffix)
	}
	return dst
}

func (f FmtrText) LogTime(dst []byte, t time.Time) []byte {
	// "t":  time shows up to millisecond: 3_04_05_000 = h:3, m:4, s:5, ms: 000
	h, m, s := t.Clock()
	dst = f.itoa(dst, h*10000+m*100+s, 6, '.')
	return f.itoa(dst, t.Nanosecond()/1e6, 3, 0)
}
func (f FmtrText) LogTimeDate(dst []byte, t time.Time) []byte {
	y, m, d := t.Date()
	return f.itoa(dst, y*10000+int(m)*100+d, 8, 0)
}
func (f FmtrText) LogTimeDay(dst []byte, t time.Time) []byte {
	return append(dst, t.Weekday().String()[0:3]...)
}

func (f FmtrText) LogTimeUnix(dst []byte, t time.Time) []byte {
	// "ts": unix second
	return f.itoa(dst, int(t.Unix()), 8, 0)
}
func (f FmtrText) LogTimeUnixMs(dst []byte, t time.Time) []byte {
	// "ts": unix second
	return f.itoa(dst, int(t.UnixNano()/1e6), 8, 0)
}

// Special type
func (f FmtrText) Nil(dst []byte, k string) []byte {
	dst = f.escape(dst, k, false)
	return append(dst, `=null`...)
}

func (f FmtrText) Error(dst []byte, k string, v error) []byte {
	dst = f.addKey(dst, k)
	if v != nil {
		return f.String(dst, k, v.Error())
	} else {
		return f.Nil(dst, k)
	}
}
func (f FmtrText) Errors(dst []byte, k string, v []error) []byte {
	dst = f.addKey(dst, k)

	idxv := len(v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range v {
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
func (f FmtrText) Bool(dst []byte, k string, v bool) []byte {
	dst = f.addKey(dst, k)
	if v {
		return append(dst, "true"...)
	}
	return append(dst, "false"...)
}
func (f FmtrText) String(dst []byte, k string, v string) []byte {
	dst = f.addKey(dst, k)
	return f.escape(dst, v, true)
}
func (f FmtrText) Int(dst []byte, k string, v int) []byte {
	dst = f.addKey(dst, k)
	return f.itoa(dst, v, 1, 0)
}
func (f FmtrText) Int8(dst []byte, k string, v int8) []byte {
	dst = f.addKey(dst, k)
	return f.itoa(dst, int(v), 1, 0)
}
func (f FmtrText) Int16(dst []byte, k string, v int16) []byte {
	dst = f.addKey(dst, k)
	return f.itoa(dst, int(v), 1, 0)
}
func (f FmtrText) Int32(dst []byte, k string, v int32) []byte {
	dst = f.addKey(dst, k)
	return f.itoa(dst, int(v), 1, 0)
}
func (f FmtrText) Int64(dst []byte, k string, v int64) []byte {
	dst = f.addKey(dst, k)
	return strconv.AppendInt(dst, v, 10)
}
func (f FmtrText) Uint(dst []byte, k string, v uint) []byte {
	dst = f.addKey(dst, k)
	return strconv.AppendUint(dst, uint64(v), 10)
}
func (f FmtrText) Uint8(dst []byte, k string, v uint8) []byte {
	dst = f.addKey(dst, k)
	return f.itoa(dst, int(v), 1, 0)
}
func (f FmtrText) Uint16(dst []byte, k string, v uint16) []byte {
	dst = f.addKey(dst, k)
	return f.itoa(dst, int(v), 1, 0)
}
func (f FmtrText) Uint32(dst []byte, k string, v uint32) []byte {
	dst = f.addKey(dst, k)
	return strconv.AppendUint(dst, uint64(v), 10)
}
func (f FmtrText) Uint64(dst []byte, k string, v uint64) []byte {
	dst = f.addKey(dst, k)
	return strconv.AppendUint(dst, v, 10)
}
func (f FmtrText) Float32(dst []byte, k string, v float32) []byte {
	dst = f.addKey(dst, k)
	return f.ftoa(dst, float64(v), 2)
}
func (f FmtrText) Float64(dst []byte, k string, v float64) []byte {
	dst = f.addKey(dst, k)
	return f.ftoa(dst, v, 2)
}

// Slice of basic data type
func (f FmtrText) Bools(dst []byte, k string, v []bool) []byte {
	dst = f.addKey(dst, k)
	idxv := len(v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range v {
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
func (f FmtrText) Strings(dst []byte, k string, v []string) []byte {
	dst = f.addKey(dst, k)
	idxv := len(v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v := range v {
		dst = f.escape(dst, v, true)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}
func (f FmtrText) Ints(dst []byte, k string, v []int) []byte {
	dst = f.addKey(dst, k)
	idxv := len(v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range v {
		dst = f.itoa(dst, v2, 1, 0)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}
func (f FmtrText) Int32s(dst []byte, k string, v []int32) []byte {
	dst = f.addKey(dst, k)
	idxv := len(v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range v {
		dst = f.itoa(dst, int(v2), 1, 0)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}
func (f FmtrText) Int64s(dst []byte, k string, v []int64) []byte {
	dst = f.addKey(dst, k)
	idxv := len(v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range v {
		dst = strconv.AppendInt(dst, v2, 10)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}
func (f FmtrText) Uints(dst []byte, k string, v []uint) []byte {
	dst = f.addKey(dst, k)
	idxv := len(v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range v {
		dst = strconv.AppendUint(dst, uint64(v2), 10)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}
func (f FmtrText) Uint8s(dst []byte, k string, v []uint8) []byte {
	dst = f.addKey(dst, k)
	idxv := len(v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range v {
		dst = f.itoa(dst, int(v2), 1, 0)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}
func (f FmtrText) Uint32s(dst []byte, k string, v []uint32) []byte {
	dst = f.addKey(dst, k)
	idxv := len(v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range v {
		dst = f.itoa(dst, int(v2), 1, 0)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}
func (f FmtrText) Uint64s(dst []byte, k string, v []uint64) []byte {
	dst = f.addKey(dst, k)
	idxv := len(v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range v {
		dst = strconv.AppendUint(dst, v2, 10)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}
func (f FmtrText) Float32s(dst []byte, k string, v []float32) []byte {
	dst = f.addKey(dst, k)
	idxv := len(v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range v {
		// dst = strconv.AppendUint(dst, v2, 10)
		dst = f.ftoa(dst, float64(v2), 2)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}
func (f FmtrText) Float64s(dst []byte, k string, v []float64) []byte {
	dst = f.addKey(dst, k)
	idxv := len(v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range v {
		// dst = strconv.AppendUint(dst, v2, 10)
		dst = f.ftoa(dst, v2, 2)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}

func (f FmtrText) addKey(dst []byte, s string) []byte {
	dst = f.escape(dst, s, false)
	return append(dst, '=')
}

func (f FmtrText) escape(dst []byte, s string, addQuote bool) []byte {
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
func (f FmtrText) escapeb(dst []byte, b []byte, addQuote bool) []byte {
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

// itoa converts int to []byte
// if minLength == 0, it will print without padding 0
// due to limit on int type, 19 digit max; 18 digit is safe.
func (f FmtrText) itoa(dst []byte, i int, minLength int, suffix byte) []byte {
	var b [22]byte
	var positiveNum = true
	if i < 0 {
		positiveNum = false
		i = -i // change the sign to positive
	}
	bIdx := len(b) - 1
	if suffix != 0 {
		b[bIdx] = suffix
		bIdx--
	}

	for i >= 10 || minLength > 1 {
		minLength--
		q := i / 10
		b[bIdx] = byte('0' + i - q*10)
		bIdx--
		i = q
	}

	b[bIdx] = byte('0' + i)
	if positiveNum == false {
		bIdx--
		b[bIdx] = '-'
	}
	return append(dst, b[bIdx:]...)
}

// ftoa takes float64 and converts and add to dst byte slice pointer.
// this is used to reduce memory allocation.
func (f FmtrText) ftoa(dst []byte, f64 float64, decPlace int) []byte {
	if int(f64) == 0 && f64 < 0 {
		dst = append(dst, '-')
	}
	dst = f.itoa(dst, int(f64), 0, 0) // add full number first

	if decPlace > 0 {
		// if decPlace == 3, multiplier will be 1000
		// get nth power
		var multiplier = 1
		for i := decPlace + 1; i > 0; i-- {
			multiplier = multiplier * 10
		}
		dst = append(dst, '.')
		tmp := int((f64 - float64(int(f64))) * float64(multiplier))
		if tmp%10 > 4 {
			tmp = tmp + 10
		}
		tmp = tmp / 10
		if f64 > 0 { // 2nd num shouldn't include decimala
			dst = f.itoa(dst, tmp, decPlace, 0)
		} else {
			dst = f.itoa(dst, -tmp, decPlace, 0)
		}
	}
	return dst
}
