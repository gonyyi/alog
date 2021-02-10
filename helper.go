package alog

import "io"

// devNull is a type for discard
type devNull int

// discard is defined here to get rid of needs to import of ioutil package.
var discard io.Writer = devNull(0)

// Write discards everything
func (devNull) Write([]byte) (int, error) {
	return 0, nil
}

// HookFn is a type for a function designed to run when certain condition meets
type HookFn func(lvl Level, tag Tag, msg []byte)

// FilterFn is a function type to be used with SetFilter.
type FilterFn func(Level, Tag) bool

// itoa converts int to []byte
// if minLength == 0, it will print without padding 0
// due to limit on int type, 19 digit max; 18 digit is safe.
// Keeping this because of minLength and suffix...
func itoa(dst []byte, i int, minLength int, suffix byte) []byte {
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
func ftoa(dst []byte, f64 float64, decPlace int) []byte {
	if int(f64) == 0 && f64 < 0 {
		dst = append(dst, '-')
	}
	dst = itoa(dst, int(f64), 0, 0) // add full number first

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
			dst = itoa(dst, tmp, decPlace, 0)
		} else {
			dst = itoa(dst, -tmp, decPlace, 0)
		}
	}
	return dst
}

func escapebyte(dst []byte, b []byte, addQuote bool) []byte {
	if addQuote {
		dst = append(dst, '"')
	}

	// For for-loop, using len() is the faster than using for-loop with a range.
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

func escapestr(dst []byte, s string, addQuote bool) []byte {
	if addQuote {
		dst = append(dst, '"')
	}

	// For for-loop, using len() is the faster than using for-loop with a range.
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

type keyable [128]bool

func (k *keyable) Init() {
	for i := 0; i < 128; i++ {
		if i == 45 || i == 46 || (47 < i && i < 58) || (64 < i && i < 91) || i == 95 || (96 < i && i < 123) {
			k[i] = true
		}
	}
}
func (k keyable) FilterStr(dst []byte, s string) []byte {
	for i := 0; i < len(s); i++ {
		if k[s[i]] {
			dst = append(dst, s[i])
		}
	}
	return dst
}
func (k keyable) FilterBytes(dst []byte, b []byte) []byte {
	for i := 0; i < len(b); i++ {
		if k[b[i]] {
			dst = append(dst, b[i])
		}
	}
	return dst
}
