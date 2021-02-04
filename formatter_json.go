package alog

import (
	"strconv"
	"time"
)

func NewFormatterJSON() Formatter {
	return &FormatterJSON{}
}

// FormatterJSON is a alog formatter that will print the result in JSON format.
type FormatterJSON struct {
	CustomHeader func([]byte) []byte
}

func (f FormatterJSON) LogCustomHeader(dst []byte) []byte {
	if f.CustomHeader != nil {
		return f.CustomHeader(dst)
	}
	return dst
}

func (f FormatterJSON) Begin(dst []byte, prefix []byte) []byte {
	dst = dst[:0] // reset first
	if prefix != nil {
		dst = append(dst, prefix...) // prefix not to be escaped
	}
	return append(dst, '{')
}

func (f FormatterJSON) End(dst []byte) []byte {
	return append(dst, '}', '\n')
}

func (f FormatterJSON) Space(dst []byte) []byte {
	return append(dst, ',')
}

// Log specific type
func (f FormatterJSON) LogLevel(dst []byte, lv Level) []byte {
	return f.safeString(dst, "level", lv.String())
}

func (f FormatterJSON) LogTag(dst []byte, tag Tag, alogTagStr [64]string, alogTagIssued int) []byte {
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
func (f FormatterJSON) LogMsg(dst []byte, s string, suffix byte) []byte { // suffix to be used for text version only
	// For JSON suffix won't be applied.
	return f.String(dst, "msg", s) // faster without addKey
}
func (f FormatterJSON) LogMsgb(dst []byte, b []byte, suffix byte) []byte {
	dst = f.addKey(dst, "msg")
	return f.escapeb(dst, b, true)
}

func (f FormatterJSON) LogTime(dst []byte, t time.Time) []byte {
	// "t":  time shows up to millisecond: 3_04_05_000 = h:3, m:4, s:5, ms: 000
	dst = append(dst, `"t":`...) // faster without addKey
	h, m, s := t.Clock()
	dst = f.itoa(dst, h*10000+m*100+s, 1, 0)
	return f.itoa(dst, t.Nanosecond()/1e6, 3, 0)
}
func (f FormatterJSON) LogTimeDate(dst []byte, t time.Time) []byte {
	dst = append(dst, `"d":`...) // faster without addKey
	y, m, d := t.Date()
	return f.itoa(dst, y*10000+int(m)*100+d, 4, 0)
}
func (f FormatterJSON) LogTimeDay(dst []byte, t time.Time) []byte {
	// "wd": 0 being sunday, 6 being saturday
	dst = append(dst, `"wd":`...) // faster without addKey
	dst = f.itoa(dst, int(t.Weekday()), 1, 0)
	return dst
}
func (f FormatterJSON) LogTimeUnix(dst []byte, t time.Time) []byte {
	// "ts": unix second
	dst = append(dst, `"ts":`...) // faster without addKey
	return f.itoa(dst, int(t.Unix()), 8, 0)
}
func (f FormatterJSON) LogTimeUnixMs(dst []byte, t time.Time) []byte {
	// "ts": unix second
	dst = append(dst, `"ts":`...) // faster without addKey
	return f.itoa(dst, int(t.UnixNano())/1e6, 8, 0)
}

// Special type
func (f FormatterJSON) Nil(dst []byte, k string) []byte {
	dst = f.escape(dst, k, true) // faster without addKey
	return append(dst, `:null`...)
}
func (f FormatterJSON) Error(dst []byte, k string, v error) []byte {
	dst = f.addKey(dst, k)
	if v != nil {
		return f.String(dst, k, v.Error())
	} else {
		return f.Nil(dst, k)
	}
}

func (f FormatterJSON) Errors(dst []byte, k string, v []error) []byte {
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
			dst = append(dst, '"', '"')
			// dst = append(dst, "null"...) // todo: check if this is acceptable (null in string array)
		}
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}

// Basic data type
// byte and rune are skipped
func (f FormatterJSON) Bool(dst []byte, k string, v bool) []byte {
	dst = f.addKey(dst, k)
	if v {
		return append(dst, "true"...)
	}
	return append(dst, "false"...)
}

func (f FormatterJSON) String(dst []byte, k string, v string) []byte {
	dst = f.addKey(dst, k)
	return f.escape(dst, v, true)
}

func (f FormatterJSON) Int(dst []byte, k string, v int) []byte {
	dst = f.addKey(dst, k)
	return f.itoa(dst, v, 1, 0)
}
func (f FormatterJSON) Int8(dst []byte, k string, v int8) []byte {
	dst = f.addKey(dst, k)
	return f.itoa(dst, int(v), 1, 0)
}
func (f FormatterJSON) Int16(dst []byte, k string, v int16) []byte {
	dst = f.addKey(dst, k)
	return f.itoa(dst, int(v), 1, 0)
}
func (f FormatterJSON) Int32(dst []byte, k string, v int32) []byte {
	dst = f.addKey(dst, k)
	return f.itoa(dst, int(v), 1, 0)
}
func (f FormatterJSON) Int64(dst []byte, k string, v int64) []byte {
	dst = f.addKey(dst, k)
	return strconv.AppendInt(dst, v, 10)
}
func (f FormatterJSON) Uint(dst []byte, k string, v uint) []byte {
	dst = f.addKey(dst, k)
	return f.itoa(dst, int(v), 1, 0)
}
func (f FormatterJSON) Uint8(dst []byte, k string, v uint8) []byte {
	dst = f.addKey(dst, k)
	return f.itoa(dst, int(v), 1, 0)
}
func (f FormatterJSON) Uint16(dst []byte, k string, v uint16) []byte {
	dst = f.addKey(dst, k)
	return f.itoa(dst, int(v), 1, 0)
}
func (f FormatterJSON) Uint32(dst []byte, k string, v uint32) []byte {
	dst = f.addKey(dst, k)
	return strconv.AppendUint(dst, uint64(v), 10)
}
func (f FormatterJSON) Uint64(dst []byte, k string, v uint64) []byte {
	dst = f.addKey(dst, k)
	return strconv.AppendUint(dst, v, 10)
}
func (f FormatterJSON) Float32(dst []byte, k string, v float32) []byte {
	dst = f.addKey(dst, k)
	return f.ftoa(dst, float64(v), 2)
}
func (f FormatterJSON) Float64(dst []byte, k string, v float64) []byte {
	dst = f.addKey(dst, k)
	return f.ftoa(dst, v, 2)
}

// Slice of basic data type
func (f FormatterJSON) Bools(dst []byte, k string, v []bool) []byte {
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
func (f FormatterJSON) Strings(dst []byte, k string, v []string) []byte {
	dst = f.addKey(dst, k)
	idxv := len(v) - 1
	if idxv == -1 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	for i, v2 := range v {
		dst = f.escape(dst, v2, true)
		if i != idxv { // if not last item
			dst = append(dst, ',')
		}
	}
	return append(dst, ']')
}

func (f FormatterJSON) Ints(dst []byte, k string, v []int) []byte {
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
func (f FormatterJSON) Int32s(dst []byte, k string, v []int32) []byte {
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
func (f FormatterJSON) Int64s(dst []byte, k string, v []int64) []byte {
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
func (f FormatterJSON) Uints(dst []byte, k string, v []uint) []byte {
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
func (f FormatterJSON) Uint8s(dst []byte, k string, v []uint8) []byte {
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
func (f FormatterJSON) Uint32s(dst []byte, k string, v []uint32) []byte {
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
func (f FormatterJSON) Uint64s(dst []byte, k string, v []uint64) []byte {
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
func (f FormatterJSON) Float32s(dst []byte, k string, v []float32) []byte {
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
func (f FormatterJSON) Float64s(dst []byte, k string, v []float64) []byte {
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

func (f FormatterJSON) addKey(dst []byte, s string) []byte {
	dst = f.escape(dst, s, true)
	return append(dst, ':')
}

func (f FormatterJSON) safeString(dst []byte, k string, v string) []byte {
	dst = append(dst, '"')
	dst = append(dst, k...)
	dst = append(dst, `":"`...)
	dst = append(dst, v...)
	return append(dst, '"')
}

func (f FormatterJSON) escape(dst []byte, s string, addQuote bool) []byte {
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
func (f FormatterJSON) escapeb(dst []byte, b []byte, addQuote bool) []byte {
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

// itoa converts int to []byte
// if minLength == 0, it will print without padding 0
// due to limit on int type, 19 digit max; 18 digit is safe.
func (f FormatterJSON) itoa(dst []byte, i int, minLength int, suffix byte) []byte {
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
func (f FormatterJSON) ftoa(dst []byte, f64 float64, decPlace int) []byte {
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
