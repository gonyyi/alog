package alog

import (
	"strconv"
	"time"
)

// check will check if Level and Tag given is good to be printed.
// If
// Eg. if setting is Level INFO, Tag USER, then
//     any log Level below INFO shouldn't be printed.
//     Also, any Tag other than USER shouldn't be printed either.
func (l *Logger) check(lvl Level, tag Tag) bool {
	switch {
	case l.logFn != nil: // logFn has the highest order if set.
		return l.logFn(lvl, tag)
	case l.logLevel > lvl: // if lvl is below lvl limit, the do not print
		return false
	case l.logTag != noTag && l.logTag&tag == noTag: // if logTag is set but Tag is not matching, then do not print
		return false
	default:
		return true
	}
}

func (l *Logger) header(lvl Level, tag Tag) {
	if l.flag&Fjson != 0 {
		l.header_json(lvl, tag)
	} else {
		l.header_text(lvl, tag)
	}
}

func (l *Logger) header_json(lvl Level, tag Tag) {
	isFirstPrinted := false
	// ----------------------
	// Handling the opening
	// ----------------------
	l.buf = l.buf[:0]
	l.addCurlyOpen()

	if l.flag&(FtimeUnix|FtimeUnixNano) != 0 {
		isFirstPrinted = true
		l.time = time.Now()
		l.buf = append(l.buf, `"ts":`...)
		if l.flag&FtimeUnix != 0 {
			l.buf = strconv.AppendInt(l.buf, l.time.Unix(), 10)
		} else {
			l.buf = strconv.AppendInt(l.buf, l.time.UnixNano(), 10)
		}
	} else if l.flag&(Fdate|Fdatemmdd|Ftime|FtimeMs|FtimeUTC) != 0 {
		isFirstPrinted = true
		l.time = time.Now()
		if l.flag&FtimeUTC != 0 {
			l.time = l.time.UTC()
		}
		// YYYYMMDD
		if l.flag&(Fdate|Fdatemmdd) != 0 {
			year, month, day := l.time.Date()
			l.buf = append(l.buf, `"d":`...)
			l.buf = strconv.AppendInt(l.buf, int64(year*10000+int(month)*100+day), 10)
		}

		// HHMMSS
		if l.flag&(Ftime|FtimeMs|FtimeUTC) != 0 {
			hour, min, sec := l.time.Clock()
			if !l.isBufSpace() {
				l.addComma()
			}
			l.buf = append(l.buf, `"t":`...)
			l.buf = strconv.AppendInt(l.buf, int64(hour*10000+min*100+sec), 10)

			// MS
			if l.flag&FtimeMs != 0 {
				l.buf = append(l.buf, `,"ms":`...)
				l.buf = strconv.AppendInt(l.buf, int64(l.time.Nanosecond()/1e6), 10)
			}
		}
	} else {
		// This can be hardcoded as it will be the first line
		l.buf = append(l.buf, '{')
	}

	// ----------------------
	// After opening
	// ----------------------
	// Add log lvl if lvl is to shown and valid range (0-6) where 0 will not show lvl prefix.
	if l.flag&Flevel != 0 && lvl < 7 {
		if !isFirstPrinted {
			isFirstPrinted = true
			l.buf = append(l.buf, `"lv":"`...)
		} else {
			l.buf = append(l.buf, `,"lv":"`...)
		}
		l.buf = append(l.buf, lvl.String()...)
		l.buf = append(l.buf, '"')
	}

	if l.flag&Ftag != 0 {
		if !isFirstPrinted {
			isFirstPrinted = true
			l.buf = append(l.buf, `"tag":[`...)
		} else {
			l.buf = append(l.buf, `,"tag":[`...)
		}
		firstItem := true
		for i := 0; i < l.logTagIssued; i++ {
			if tag&(1<<i) != 0 {
				if firstItem {
					firstItem = false
					l.buf = append(l.buf, '"')
				} else {
					l.buf = append(l.buf, ',', '"')
				}
				l.buf = append(l.buf, l.logTagString[i]...)
				l.buf = append(l.buf, '"')
			}
		}
		l.buf = append(l.buf, ']')
	}
	return
}

func (l *Logger) header_text(lvl Level, tag Tag) {
	// Add prefix
	l.buf = l.buf[:0]
	if l.flag&Fprefix != 0 {
		l.buf = append(l.buf, l.prefix...)
	}

	if l.flag&(FtimeUnix|FtimeUnixNano) != 0 {
		l.time = time.Now()
		if l.flag&FtimeUnix != 0 {
			l.buf = strconv.AppendInt(l.buf, l.time.Unix(), 10)
		} else {
			l.buf = strconv.AppendInt(l.buf, l.time.UnixNano(), 10)
		}
		l.buf = append(l.buf, ' ')
	} else if l.flag&(Fdate|Fdatemmdd|Ftime|FtimeMs|FtimeUTC) != 0 {
		l.time = time.Now()
		if l.flag&FtimeUTC != 0 {
			l.time = l.time.UTC()
		}
		if l.flag&(Fdate|Fdatemmdd) != 0 {
			year, month, day := l.time.Date()
			// if both YYMMDD and YYYYMMDD is given, YYYYMMDD will be used
			if l.flag&Fdate != 0 {
				l.buf = itoa(l.buf, year, 4, '/')
			}
			// MMDD will be always added ass it's a common denominator of
			// FdateYYMMDD|Fdate|Fdatemmdd
			l.buf = itoa(l.buf, int(month), 2, '/')
			l.buf = itoa(l.buf, day, 2, ' ')
		}
		if l.flag&(Ftime|FtimeMs|FtimeUTC) != 0 {
			hour, min, sec := l.time.Clock()
			l.buf = itoa(l.buf, hour, 2, ':')
			l.buf = itoa(l.buf, min, 2, ':')
			if l.flag&FtimeMs != 0 {
				l.buf = itoa(l.buf, sec, 2, '.')
				l.buf = itoa(l.buf, l.time.Nanosecond()/1e6, 3, 0)
			} else {
				l.buf = itoa(l.buf, sec, 2, 0)
			}
		}
	}

	// Add log lvl if lvl is to shown and valid range (0-6) where 0 will not show lvl prefix.
	if l.flag&Flevel != 0 && lvl < 7 {
		if !l.isBufSpace() {
			l.buf = append(l.buf, ' ')
		}
		l.buf = append(l.buf, lvl.name_terminal()...)
	}

	if l.flag&Ftag != 0 {
		if !l.isBufSpace() {
			l.buf = append(l.buf, ' ')
		}
		l.buf = append(l.buf, '[')

		if tag != 0 {
			firstItem := true
			for i := 0; i < l.logTagIssued; i++ {
				if tag&(1<<i) != 0 {
					if firstItem {
						firstItem = false
					} else {
						l.buf = append(l.buf, ',')
					}
					l.buf = append(l.buf, l.logTagString[i]...)
				}
			}
		}
		l.buf = append(l.buf, ']')
	}

}

func (l *Logger) isBufSpace() bool {
	// If a buffer is fresh (nothing in it), or last item has a space, or json just started..
	if lb := len(l.buf); lb == 0 || l.buf[lb-1] == ' ' || l.buf[lb-1] == '{' {
		return true
	}
	return false
}

// finalize will add newline to the end of log if missing,
// also write it to writer, and clear the buffer.
func (l *Logger) finalize() (n int, err error) {
	if l.flag&Fjson != 0 {
		l.addCurlyClose()
	}
	l.addNewline()
	n, err = l.out.Write(l.buf)
	return n, err
}
func (l *Logger) addEscapeBasic(s string) {
	for j := 0; j < len(s); j++ {
		switch s[j] {
		case '\n':
			l.buf = append(l.buf, '\\', 'n')
		case '\r':
			l.buf = append(l.buf, '\\', 'r')
		case '\b':
			l.buf = append(l.buf, '\\', 'b')
		case '\f':
			l.buf = append(l.buf, '\\', 'f')
		default:
			l.buf = append(l.buf, s[j])
		}
	}
}
func (l *Logger) addEscapeJStr(s string) {
	l.buf = append(l.buf, '"')
	for j := 0; j < len(s); j++ {
		switch s[j] {
		case '\\':
			l.buf = append(l.buf, '\\')
		case '"':
			l.buf = append(l.buf, '\\', '"')
		case '\n':
			l.buf = append(l.buf, '\\', 'n')
		case '\t':
			l.buf = append(l.buf, '\\', 't')
		case '\r':
			l.buf = append(l.buf, '\\', 'r')
		case '\b':
			l.buf = append(l.buf, '\\', 'b')
		case '\f':
			l.buf = append(l.buf, '\\', 'f')
		default:
			l.buf = append(l.buf, s[j])
		}
	}
	l.buf = append(l.buf, '"')
}
func (l *Logger) addSquareOpen() {
	l.buf = append(l.buf, '[')
}
func (l *Logger) addSquareClose() {
	l.buf = append(l.buf, ']')
}
func (l *Logger) addCurlyOpen() {
	l.buf = append(l.buf, '{')
}
func (l *Logger) addCurlyClose() {
	l.buf = append(l.buf, '}')
}
func (l *Logger) addColon() {
	l.buf = append(l.buf, ':')
}
func (l *Logger) addDblQuote() {
	l.buf = append(l.buf, '"')
}
func (l *Logger) addComma() {
	l.buf = append(l.buf, ',')
}
func (l *Logger) addNewline() {
	l.buf = append(l.buf, newline)
}
