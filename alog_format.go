package alog

import (
	"strconv"
	"time"
)

// header will add date/time/prefix/Level.
func (l *Logger) header(buf *[]byte, lvl Level, tag Tag) {
	if l.flag&Fjson != 0 {
		// l.sbufc = 0
		// JSON format: `{ d: 20201012, t:151223, ms:12345, type: "info", tag: [], msg: "my message", err: "additional error", add: "additional data" }`
		if l.flag&(FdateYYMMDD|FdateYYYYMMDD|FdateMMDD|Ftime|FtimeMs|FtimeUTC) != 0 {
			l.time = time.Now()
			if l.flag&FtimeUTC != 0 {
				l.time = l.time.UTC()
			}
			// YYYYMMDD
			year, month, day := l.time.Date()
			// *buf = append(*buf, `{"d":`...)
			// itoa(buf, year*10000+int(month)*100+day, 0, ',')
			l.buf = append(l.buf, `{"d":`...)
			l.buf = strconv.AppendInt(l.buf, int64(year*10000+int(month)*100+day), 10)

			// HHMMSS
			hour, min, sec := l.time.Clock()
			// *buf = append(*buf, `"t":`...)
			// itoa(buf, hour*10000+min*100+sec, 0, ',')
			l.buf = append(l.buf, `,"t":`...)
			l.buf = strconv.AppendInt(l.buf, int64(hour*10000+min*100+sec), 10)

			// MS
			if l.flag&FtimeMs != 0 {
				// *buf = append(*buf, `"ns":`...)
				// itoa(buf, l.time.Nanosecond()/1e3, 6, ',')
				l.buf = append(l.buf, `,"ns":`...)
				l.buf = strconv.AppendInt(l.buf, int64(l.time.Nanosecond()/1e3), 10)
			}
		} else {
			// This can be hardcoded as it will be the first line
			// *buf = append(*buf, '{')
			l.buf = append(l.buf, '{')
		}

		// Add log lvl if lvl is to shown and valid range (0-6) where 0 will not show lvl prefix.
		if lvl < 7 {
			l.buf = append(l.buf, `,"lv":"`...)
			l.buf = append(l.buf, l.levelStringForJson[lvl]...)
			l.buf = append(l.buf, '"')
		}

		if tag != 0 {
			l.buf = append(l.buf, `,"tag":[`...)
			firstItem := true
			for i := 0; i < l.logTagIssued; i++ {
				if tag&(1<<i) != 0 {
					if firstItem {
						firstItem = false
						l.buf = append(l.buf, '"')
					} else {
						l.buf = append(l.buf, `,"`...)
					}
					l.buf = append(l.buf, l.logTagString[i]...)
					l.buf = append(l.buf, `"`...)

				}
			}
			l.buf = append(l.buf, ']')
		}
		return
	}

	if l.flag&(FdateYYMMDD|FdateYYYYMMDD|FdateMMDD|Ftime|FtimeMs|FtimeUTC) != 0 {
		l.time = time.Now()
		if l.flag&FtimeUTC != 0 {
			l.time = l.time.UTC()
		}
		if l.flag&(FdateYYMMDD|FdateYYYYMMDD|FdateMMDD) != 0 {
			year, month, day := l.time.Date()
			// if both YYMMDD and YYYYMMDD is given, YYYYMMDD will be used
			if l.flag&FdateYYYYMMDD != 0 {
				itoa(buf, year, 4, '/')
			} else if l.flag&FdateYYMMDD != 0 {
				itoa(buf, year%100, 2, '/')
			}
			// MMDD will be always added ass it's a common denominator of
			// FdateYYMMDD|FdateYYYYMMDD|FdateMMDD
			itoa(buf, int(month), 2, '/')
			itoa(buf, day, 2, ' ')
		}
		if l.flag&(Ftime|FtimeMs|FtimeUTC) != 0 {
			hour, min, sec := l.time.Clock()
			itoa(buf, hour, 2, ':')
			itoa(buf, min, 2, ':')
			if l.flag&FtimeMs != 0 {
				itoa(buf, sec, 2, '.')
				itoa(buf, l.time.Nanosecond()/1e3, 6, ' ')
			} else {
				itoa(buf, sec, 2, ' ')
			}
		}
	}
	// Add prefix
	if l.flag&Fprefix != 0 {
		*buf = append(*buf, l.prefix...)
	}

	// Add log lvl if lvl is to shown and valid range (0-6) where 0 will not show lvl prefix.
	if l.flag&Flevel != 0 && lvl < 7 {
		*buf = append(*buf, l.levelString[lvl]...)
	}

	if tag != 0 {
		// *buf = append(*buf, `"tag":[`...)
		firstItem := true
		for i := 0; i < l.logTagIssued; i++ {
			if tag&(1<<i) != 0 {
				if firstItem {
					firstItem = false
				} else {
					*buf = append(*buf, `/`...)
				}
				*buf = append(*buf, l.logTagString[i]...)
			}
		}
		*buf = append(*buf, "; "...)
	}
}

// finalize will add newline to the end of log if missing,
// also write it to writer, and clear the buffer.
func (l *Logger) finalize() (n int, err error) {
	if l.flag&Fjson != 0 {
		for i, v := range l.buf {
			if v == newline {
				// l.buf[i] = byte(' ')
				l.buf = append(l.buf[0:i], append(newlineRepl, l.buf[i+1:]...)...)
			}
		}
		// l.sbufc = 0
		// l.sbufc += copy(l.sbuf, "\"}\n")
		l.buf = append(l.buf, "}\n"...)
		// l.buf = append(l.buf, newline)
	} else {
		l.buf = append(l.buf, newline)
	}

	// If bufUseBuffer is false or current size is bigger than the buffer size,
	// print the buffer and reset it.
	n, err = l.out.Write(l.buf)
	l.buf = l.buf[:0]
	return n, err
}
