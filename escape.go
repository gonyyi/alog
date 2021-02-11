package alog

func escapeBytes(dst []byte, b []byte, addQuote bool) []byte {
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

func escapeString(dst []byte, s string, addQuote bool, suffix byte) []byte {
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

var keyable = newEscKeys()

func newEscKeys() *escKeys {
	k := &escKeys{}
	k.Init()
	return k
}

type escKeys [128]bool

func (k *escKeys) Init() {
	for i := 0; i < 128; i++ {
		if i == 45 || i == 46 || (47 < i && i < 58) || (64 < i && i < 91) || i == 95 || (96 < i && i < 123) {
			k[i] = true
		}
	}
}

func (k escKeys) String(dst []byte, s string) []byte {
	for i := 0; i < len(s); i++ {
		if k[s[i]] {
			dst = append(dst, s[i])
		}
	}
	return dst
}

func (k escKeys) Bytes(dst []byte, b []byte) []byte {
	for i := 0; i < len(b); i++ {
		if k[b[i]] {
			dst = append(dst, b[i])
		}
	}
	return dst
}
