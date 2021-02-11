package alog

var conv Converter = &Conv{}

type Converter interface {
	Intf(dst []byte, i int, minLength int, suffix byte) []byte
	Floatf(dst []byte, f64 float64, decPlace int) []byte
	Int(dst []byte, i int, quote bool, suffix byte) []byte
	Float(dst []byte, f float64, quote bool, suffix byte) []byte
	String(dst []byte, s string, quote bool, suffix byte) []byte
	Bytes(dst []byte, s []byte, quote bool, suffix byte) []byte
	Bool(dst []byte, b bool, quote bool, suffix byte) []byte
	Error(dst []byte, err error, quote bool, suffix byte) []byte
}

type Conv struct{}

// itoa converts int to []byte
// if minLength == 0, it will print without padding 0
// due to limit on int type, 19 digit max; 18 digit is safe.
// Keeping this because of minLength and suffix...
func (c *Conv) Intf(dst []byte, i int, minLength int, suffix byte) []byte {
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
func (c *Conv) Floatf(dst []byte, f64 float64, decPlace int) []byte {
	if int(f64) == 0 && f64 < 0 {
		dst = append(dst, '-')
	}
	dst = c.Intf(dst, int(f64), 0, 0) // add full number first

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
			dst = c.Intf(dst, tmp, decPlace, 0)
		} else {
			dst = c.Intf(dst, -tmp, decPlace, 0)
		}
	}
	return dst
}

func (c *Conv) Float(dst []byte, f float64, quote bool, suffix byte) []byte {
	if quote {
		dst = append(dst, '"')
		dst = c.Floatf(dst, f, 2)
		if suffix != 0 {
			return append(dst, '"', suffix)
		} else {
			return append(dst, '"')
		}
	} else {
		if suffix != 0 {
			dst = c.Floatf(dst, f, 2)
			return append(dst, suffix)
		} else {
			return c.Floatf(dst, f, 2)
		}
	}
}

func (c *Conv) Int(dst []byte, i int, quote bool, suffix byte) []byte {
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

func (c *Conv) Bytes(dst []byte, s []byte, quote bool, suffix byte) []byte {
	dst = append(dst, escapeBytes(dst, s, quote)...)
	if suffix != 0 {
		return append(dst, suffix)
	}
	return dst
}

func (c *Conv) String(dst []byte, s string, quote bool, suffix byte) []byte {
	return escapeString(dst, s, quote, suffix)
}

func (c *Conv) Bool(dst []byte, b bool, quote bool, suffix byte) []byte {
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

func (c *Conv) Error(dst []byte, err error, quote bool, suffix byte) []byte {
	if err != nil {
		if quote {
			dst = append(dst, '"')
			dst = append(dst, err.Error()...)
			dst = append(dst, '"')
		} else {
			dst = append(dst, err.Error()...)
		}
	} else {
		dst = append(dst, "null"...)
	}
	if suffix != 0 {
		return append(dst, suffix)
	}
	return dst
}
