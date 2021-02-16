package alog

type FormatterConverter interface {
	Int(dst []byte, i int, quote bool, suffix byte) []byte
	Intf(dst []byte, i int, minLength int, suffix byte) []byte
	Float(dst []byte, f64 float64, quote bool, suffix byte) []byte
	Floatf(dst []byte, f64 float64, decPlace int, suffix byte) []byte
	Bool(dst []byte, b bool, quote bool, suffix byte) []byte
	Error(dst []byte, err error, quote bool, suffix byte) []byte
}

type formatterConvBasic struct{}

func (c *formatterConvBasic) Int(dst []byte, i int, quote bool, suffix byte) []byte {
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

func (c *formatterConvBasic) Intf(dst []byte, i int, minLength int, suffix byte) []byte {
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

func (c *formatterConvBasic) Float(dst []byte, f float64, quote bool, suffix byte) []byte {
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

func (c *formatterConvBasic) Floatf(dst []byte, f float64, decPlace int, suffix byte) []byte {
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
		// (f - float32(int(f))) * float32(multiplier)
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

func (c *formatterConvBasic) Bool(dst []byte, b bool, quote bool, suffix byte) []byte {
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

func (c *formatterConvBasic) Error(dst []byte, err error, quote bool, suffix byte) []byte {
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
