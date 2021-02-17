package alog

var conv = func() *convert {
	c := convert{}
	c.Init()
	return &c
}()

// For future need:
//	type Converter interface {
//		Init()
//		EscKey(dst []byte, s string, addQuote bool, suffix byte) []byte
//		EscKeyBytes(dst []byte, b []byte, addQuote bool, suffix byte) []byte
//		EscString(dst []byte, s string, addQuote bool, suffix byte) []byte
//		EscStringBytes(dst []byte, p []byte, addQuote bool, suffix byte) []byte
//		Int(dst []byte, i int, quote bool, suffix byte) []byte
//		Intf(dst []byte, i int, minLength int, suffix byte) []byte
//		Float(dst []byte, f float64, quote bool, suffix byte) []byte
//		Floatf(dst []byte, f float64, decPlace int, suffix byte) []byte
//		Bool(dst []byte, b bool, quote bool, suffix byte) []byte
//		Error(dst []byte, err error, quote bool, suffix byte) []byte
//	}

type convert [128]bool

func (c *convert) Init() {
	// Initialize key values
	for i := 0; i < 128; i++ {
		if i == 45 || i == 46 || (47 < i && i < 58) || (64 < i && i < 91) || i == 95 || (96 < i && i < 123) {
			c[i] = true
		}
	}
}

func (c convert) EscKey(dst []byte, s string, addQuote bool, suffix byte) []byte {
	if addQuote {
		dst = append(dst, '"')
	}

	// Only allowed
	for i := 0; i < len(s); i++ {
		if c[s[i]] {
			dst = append(dst, s[i])
		}
	}

	// Add all
	//dst = append(dst, s...)

	// Smarty way, but doesn't seem fast
	//cur := 0
	//for i := 0; i < len(s); i++ {
	//	if !c[s[i]] {
	//		dst = append(dst, s[cur:i+1]...)
	//		cur = i + 1
	//	}
	//}
	//dst = append(dst, s[cur:]...)

	if addQuote {
		if suffix != 0 {
			return append(dst, '"', suffix)
		}
		return append(dst, '"')
	} else if suffix != 0 {
		return append(dst, suffix)
	}
	return dst
}

func (c convert) EscKeyBytes(dst []byte, b []byte, addQuote bool, suffix byte) []byte {
	if addQuote {
		dst = append(dst, '"')
	}
	for i := 0; i < len(b); i++ {
		if c[b[i]] {
			dst = append(dst, b[i])
		}
	}
	if addQuote {
		if suffix != 0 {
			return append(dst, '"', suffix)
		}
		return append(dst, '"')
	} else if suffix != 0 {
		return append(dst, suffix)
	}
	return dst
}

func (c convert) EscString(dst []byte, s string, addQuote bool, suffix byte) []byte {
	if addQuote {
		dst = append(dst, '"')
	}

	// For for-loop, using len() is the faster than using for-loop with a range.
	cur := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"', '\\':
			dst = append(dst, s[cur:i]...)
			dst = append(dst, '\\', s[i])
			cur = i + 1
		case '\n':
			dst = append(dst, s[cur:i]...)
			dst = append(dst, '\\', 'n')
			cur = i + 1
		case '\t':
			dst = append(dst, s[cur:i]...)
			dst = append(dst, '\\', 't')
			cur = i + 1
		case '\r':
			dst = append(dst, s[cur:i]...)
			dst = append(dst, '\\', 'r')
			cur = i + 1
		case '\b':
			dst = append(dst, s[cur:i]...)
			dst = append(dst, '\\', 'b')
			cur = i + 1
		case '\f':
			dst = append(dst, s[cur:i]...)
			dst = append(dst, '\\', 'f')
			cur = i + 1
		}
	}
	dst = append(dst, s[cur:]...)
	if addQuote {
		if suffix != 0 {
			// quote and suffix
			return append(dst, '"', suffix)
		}
		// quote, but no suffix
		return append(dst, '"')
	} else if suffix != 0 {
		// no quote, but suffix
		return append(dst, suffix)
	}
	return dst
}

func (c convert) EscStringBytes(dst []byte, p []byte, addQuote bool, suffix byte) []byte {
	if addQuote {
		dst = append(dst, '"')
	}

	// For for-loop, using len() is the faster than using for-loop with a range.
	cur := 0
	for i := 0; i < len(p); i++ {
		switch p[i] {
		case '"', '\\':
			dst = append(dst, p[cur:i]...)
			dst = append(dst, '\\', p[i])
			cur = i + 1
		case '\n':
			dst = append(dst, p[cur:i]...)
			dst = append(dst, '\\', 'n')
			cur = i + 1
		case '\t':
			dst = append(dst, p[cur:i]...)
			dst = append(dst, '\\', 't')
			cur = i + 1
		case '\r':
			dst = append(dst, p[cur:i]...)
			dst = append(dst, '\\', 'r')
			cur = i + 1
		case '\b':
			dst = append(dst, p[cur:i]...)
			dst = append(dst, '\\', 'b')
			cur = i + 1
		case '\f':
			dst = append(dst, p[cur:i]...)
			dst = append(dst, '\\', 'f')
			cur = i + 1
		}
	}
	dst = append(dst, p[cur:]...)
	if addQuote {
		if suffix != 0 {
			// quote and suffix
			return append(dst, '"', suffix)
		}
		// quote, but no suffix
		return append(dst, '"')
	} else if suffix != 0 {
		// no quote, but suffix
		return append(dst, suffix)
	}
	return dst
}

func (c convert) Int(dst []byte, i int, quote bool, suffix byte) []byte {
	if quote {
		dst = append(dst, '"')
		dst = c.Intf(dst, i, 0, 0)
		if suffix != 0 {
			return append(dst, '"', suffix)
		} else {
			return append(dst, '"')
		}
	} else {
		return c.Intf(dst, i, 0, suffix)
	}
}

func (c convert) Intf(dst []byte, i int, minLength int, suffix byte) []byte {
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

func (c convert) Float(dst []byte, f float64, quote bool, suffix byte) []byte {
	if quote {
		dst = append(dst, '"')
		dst = c.Floatf(dst, f, 2, 0)
		if suffix != 0 {
			return append(dst, '"', suffix)
		} else {
			return append(dst, '"')
		}
	}
	return c.Floatf(dst, f, 2, suffix)
}

func (c convert) Floatf(dst []byte, f float64, decPlace int, suffix byte) []byte {
	if int(f) == 0 && f < 0 {
		dst = append(dst, '-')
	}

	if decPlace > 0 {
		dst = c.Intf(dst, int(f), 0, 0) // add full number first

		// if decPlace == 3, multiplier will be 10000
		// get nth power
		// if decPlace == 2, then 3:10, 2:10*10, 1:10*10*10 = 1000
		var multiplier = 1
		for i := decPlace + 1; i > 0; i-- {
			multiplier = multiplier * 10
		}
		dst = append(dst, '.')
		// 3.145 --> 3.145 - 3 = 0.145 * 100 = 14.5
		// (fmt - float32(int(fmt))) * float32(multiplier)
		// (3.145 - float32(3)) * float32(1000) = 145
		tmp := int((f - float64(int(f))) * float64(multiplier))
		if tmp%10 > 4 {
			tmp = tmp + 10
		}
		tmp = tmp / 10

		if f > 0 { // 2nd num shouldn't include decimala
			return c.Intf(dst, tmp, decPlace, suffix)
		} else {
			return c.Intf(dst, -tmp, decPlace, suffix)
		}
	}

	return c.Intf(dst, int(f), 0, suffix)
}

func (c convert) Bool(dst []byte, b bool, quote bool, suffix byte) []byte {
	if quote {
		if b {
			dst = append(dst, `"true"`...)
		} else {
			dst = append(dst, `"false"`...)
		}
	} else {
		if b {
			dst = append(dst, "true"...)
		} else {
			dst = append(dst, "false"...)
		}
	}
	if suffix != 0 {
		return append(dst, suffix)
	}
	return dst
}

func (c convert) Error(dst []byte, err error, quote bool, suffix byte) []byte {
	if err != nil {
		if quote {
			dst = append(dst, '"')
			dst = append(dst, err.Error()...)
			if suffix != 0 {
				dst = append(dst, '"', suffix)
			} else {
				dst = append(dst, '"')
			}
		} else {
			dst = append(dst, err.Error()...)
			if suffix != 0 {
				return append(dst, suffix)
			}
		}
	} else {
		if quote {
			dst = append(dst, '"')
			dst = append(dst, "null"...)
			if suffix != 0 {
				return append(dst, '"', suffix)
			} else {
				return append(dst, '"')
			}
		} else {
			dst = append(dst, "null"...)
			if suffix != 0 {
				return append(dst, suffix)
			}
		}
	}
	return dst
}
