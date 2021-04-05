package alog

// This formatdString.go is from Zerolog's (internal/json/string.go)

import "unicode/utf8"

const hex = "0123456789abcdef"

var noEscapeTable = [256]bool{}

func init() {
	for i := 0; i <= 0x7e; i++ {
		noEscapeTable[i] = i >= 0x20 && i != '\\' && i != '"'
	}
}

func appendString(dst []byte, s string, addQuote bool) []byte {
	// Start with a double quote.
	if addQuote {
		dst = append(dst, '"')
		for i := 0; i < len(s); i++ {
			if !noEscapeTable[s[i]] {
				// We encountered a character that needs to be encoded. Switch
				// to complex version of the algorithm.
				dst = appendStringComplex(dst, s, i)
				return append(dst, '"')
			}
		}
		dst = append(dst, s...)
		return append(dst, '"')
	}


	for i := 0; i < len(s); i++ {
		if !noEscapeTable[s[i]] {
			return appendStringComplex(dst, s, i)
		}
	}
	return append(dst, s...)
}

func appendStringComplex(dst []byte, s string, i int) []byte {
	start := 0
	for i < len(s) {
		b := s[i]
		if b >= utf8.RuneSelf {
			r, size := utf8.DecodeRuneInString(s[i:])
			if r == utf8.RuneError && size == 1 {
				// In case of error, first append previous simple characters to
				// the byte slice if any and append a remplacement character code
				// in place of the invalid sequence.
				if start < i {
					dst = append(dst, s[start:i]...)
				}
				dst = append(dst, `\ufffd`...)
				i += size
				start = i
				continue
			}
			i += size
			continue
		}
		if noEscapeTable[b] {
			i++
			continue
		}
		// We encountered a character that needs to be encoded.
		// Let's append the previous simple characters to the byte slice
		// and switch our operation to read and encode the remainder
		// characters byte-by-byte.
		if start < i {
			dst = append(dst, s[start:i]...)
		}
		switch b {
		case '"', '\\':
			dst = append(dst, '\\', b)
		case '\b':
			dst = append(dst, '\\', 'b')
		case '\f':
			dst = append(dst, '\\', 'f')
		case '\n':
			dst = append(dst, '\\', 'n')
		case '\r':
			dst = append(dst, '\\', 'r')
		case '\t':
			dst = append(dst, '\\', 't')
		default:
			dst = append(dst, '\\', 'u', '0', '0', hex[b>>4], hex[b&0xF])
		}
		i++
		start = i
	}
	if start < len(s) {
		dst = append(dst, s[start:]...)
	}
	return dst
}
