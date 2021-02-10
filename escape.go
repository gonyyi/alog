package alog

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

func (k keyable) String(dst []byte, s string) []byte {
	for i := 0; i < len(s); i++ {
		if k[s[i]] {
			dst = append(dst, s[i])
		}
	}
	return dst
}

func (k keyable) Bytes(dst []byte, b []byte) []byte {
	for i := 0; i < len(b); i++ {
		if k[b[i]] {
			dst = append(dst, b[i])
		}
	}
	return dst
}
