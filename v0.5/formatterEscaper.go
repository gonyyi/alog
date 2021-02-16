package alog

type FormatterEscaper interface {
	Init()
	Key(dst []byte, s string, addQuote bool, suffix byte) []byte
	KeyBytes(dst []byte, p []byte, addQuote bool, suffix byte) []byte
	Val(dst []byte, s string, addQuote bool, suffix byte) []byte
	ValBytes(dst []byte, p []byte, addQuote bool, suffix byte) []byte
}

type formatterEscBasic [128]bool

func (e *formatterEscBasic) Init() {
	for i := 0; i < 128; i++ {
		if i == 45 || i == 46 || (47 < i && i < 58) || (64 < i && i < 91) || i == 95 || (96 < i && i < 123) {
			e[i] = true
		}
	}
}

func (e formatterEscBasic) Key(dst []byte, s string, addQuote bool, suffix byte) []byte {
	if addQuote {
		dst = append(dst, '"')
	}
	for i := 0; i < len(s); i++ {
		if e[s[i]] {
			dst = append(dst, s[i])
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

func (e formatterEscBasic) KeyBytes(dst []byte, b []byte, addQuote bool, suffix byte) []byte {
	if addQuote {
		dst = append(dst, '"')
	}
	for i := 0; i < len(b); i++ {
		if e[b[i]] {
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

func (e formatterEscBasic) Val(dst []byte, s string, addQuote bool, suffix byte) []byte {
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

func (e formatterEscBasic) ValBytes(dst []byte, p []byte, addQuote bool, suffix byte) []byte {
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
