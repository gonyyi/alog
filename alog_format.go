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

func (l *Logger) header_json(lvl Level, tag Tag) {
	isFirstPrinted := false
	// ----------------------
	// Handling the opening
	// ----------------------
	if l.flag&(FtimeUnix|FtimeUnixNano) != 0 {
		isFirstPrinted = true
		l.time = time.Now()
		l.buf = append(l.buf, `{"ts":`...)
		if l.flag&FtimeUnix != 0 {
			l.buf = strconv.AppendInt(l.buf, l.time.Unix(), 10)
		} else {
			l.buf = strconv.AppendInt(l.buf, l.time.UnixNano(), 10)
		}
	} else if l.flag&(Fyear|Fdate|Ftime|FtimeMs|FtimeUTC) != 0 {
		isFirstPrinted = true
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

	// ----------------------
	// After opening
	// ----------------------
	// Add log lvl if lvl is to shown and valid range (0-6) where 0 will not show lvl prefix.
	if lvl < 7 {
		if !isFirstPrinted {
			l.buf = append(l.buf, `"lv":"`...)
		} else {
			l.buf = append(l.buf, `,"lv":"`...)
		}
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

func (l *Logger) header_text(lvl Level, tag Tag) {
	// Add prefix
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
	} else if l.flag&(Fyear|Fdate|Ftime|FtimeMs|FtimeUTC) != 0 {
		l.time = time.Now()
		if l.flag&FtimeUTC != 0 {
			l.time = l.time.UTC()
		}
		if l.flag&(Fyear|Fdate) != 0 {
			year, month, day := l.time.Date()
			// if both YYMMDD and YYYYMMDD is given, YYYYMMDD will be used
			if l.flag&Fyear != 0 {
				l.buf = itoa(l.buf, year, 4, '/')
			}
			// MMDD will be always added ass it's a common denominator of
			// FdateYYMMDD|Fyear|Fdate
			l.buf = itoa(l.buf, int(month), 2, '/')
			l.buf = itoa(l.buf, day, 2, ' ')
		}
		if l.flag&(Ftime|FtimeMs|FtimeUTC) != 0 {
			hour, min, sec := l.time.Clock()
			l.buf = itoa(l.buf, hour, 2, ':')
			l.buf = itoa(l.buf, min, 2, ':')
			if l.flag&FtimeMs != 0 {
				l.buf = itoa(l.buf, sec, 2, '.')
				l.buf = itoa(l.buf, l.time.Nanosecond()/1e3, 6, ' ')
			} else {
				l.buf = itoa(l.buf, sec, 2, ' ')
			}
		}
	}

	// Add log lvl if lvl is to shown and valid range (0-6) where 0 will not show lvl prefix.
	if l.flag&Flevel != 0 && lvl < 7 {
		l.buf = append(l.buf, l.levelString[lvl]...)
	}

	if tag != 0 {
		l.buf = append(l.buf, `tag=`...)
		firstItem := true
		for i := 0; i < l.logTagIssued; i++ {
			if tag&(1<<i) != 0 {
				if firstItem {
					firstItem = false
				} else {
					l.buf = append(l.buf, `/`...)
				}
				l.buf = append(l.buf, l.logTagString[i]...)
			}
		}
		l.buf = append(l.buf, ", "...)
	}
}

// // header will add date/time/prefix/Level.
// // Printing priority:
// // 	Fjson > FtimeUnix | FtimeUnixNano
// func (l *Logger) header(buf *[]byte, lvl Level, tag Tag) {
// 	// ======================
// 	// IF JSON
// 	// ======================
// 	if l.flag&Fjson != 0 {
// 		isFirstPrinted := false
// 		// ----------------------
// 		// Handling the opening
// 		// ----------------------
// 		if l.flag&(FtimeUnix|FtimeUnixNano) != 0 {
// 			isFirstPrinted = true
// 			l.time = time.Now()
// 			l.buf = append(l.buf, `{"ts":`...)
// 			if l.flag&FtimeUnix != 0 {
// 				l.buf = strconv.AppendInt(l.buf, l.time.Unix(), 10)
// 			} else {
// 				l.buf = strconv.AppendInt(l.buf, l.time.UnixNano(), 10)
// 			}
// 		} else if l.flag&(Fyear|Fdate|Ftime|FtimeMs|FtimeUTC) != 0 {
// 			isFirstPrinted = true
// 			l.time = time.Now()
// 			if l.flag&FtimeUTC != 0 {
// 				l.time = l.time.UTC()
// 			}
// 			// YYYYMMDD
// 			year, month, day := l.time.Date()
// 			// *buf = append(*buf, `{"d":`...)
// 			// itoa(buf, year*10000+int(month)*100+day, 0, ',')
// 			l.buf = append(l.buf, `{"d":`...)
// 			l.buf = strconv.AppendInt(l.buf, int64(year*10000+int(month)*100+day), 10)
//
// 			// HHMMSS
// 			hour, min, sec := l.time.Clock()
// 			// *buf = append(*buf, `"t":`...)
// 			// itoa(buf, hour*10000+min*100+sec, 0, ',')
// 			l.buf = append(l.buf, `,"t":`...)
// 			l.buf = strconv.AppendInt(l.buf, int64(hour*10000+min*100+sec), 10)
//
// 			// MS
// 			if l.flag&FtimeMs != 0 {
// 				// *buf = append(*buf, `"ns":`...)
// 				// itoa(buf, l.time.Nanosecond()/1e3, 6, ',')
// 				l.buf = append(l.buf, `,"ns":`...)
// 				l.buf = strconv.AppendInt(l.buf, int64(l.time.Nanosecond()/1e3), 10)
// 			}
// 		} else {
// 			// This can be hardcoded as it will be the first line
// 			// *buf = append(*buf, '{')
// 			l.buf = append(l.buf, '{')
// 		}
//
// 		// ----------------------
// 		// After opening
// 		// ----------------------
// 		// Add log lvl if lvl is to shown and valid range (0-6) where 0 will not show lvl prefix.
// 		if lvl < 7 {
// 			if !isFirstPrinted {
// 				l.buf = append(l.buf, `"lv":"`...)
// 			} else {
// 				l.buf = append(l.buf, `,"lv":"`...)
// 			}
// 			l.buf = append(l.buf, l.levelStringForJson[lvl]...)
// 			l.buf = append(l.buf, '"')
// 		}
//
// 		if tag != 0 {
// 			l.buf = append(l.buf, `,"tag":[`...)
// 			firstItem := true
// 			for i := 0; i < l.logTagIssued; i++ {
// 				if tag&(1<<i) != 0 {
// 					if firstItem {
// 						firstItem = false
// 						l.buf = append(l.buf, '"')
// 					} else {
// 						l.buf = append(l.buf, `,"`...)
// 					}
// 					l.buf = append(l.buf, l.logTagString[i]...)
// 					l.buf = append(l.buf, `"`...)
// 				}
// 			}
// 			l.buf = append(l.buf, ']')
// 		}
// 		return
// 	}
//
// 	// ======================
// 	// NOT JSON
// 	// ======================
// 	if l.flag&(FtimeUnix|FtimeUnixNano) != 0 {
// 		l.time = time.Now()
// 		if l.flag&FtimeUnix != 0 {
// 			l.buf = strconv.AppendInt(l.buf, l.time.Unix(), 10)
// 		} else {
// 			l.buf = strconv.AppendInt(l.buf, l.time.UnixNano(), 10)
// 		}
// 		l.buf = append(l.buf, ' ')
// 	} else if l.flag&(Fyear|Fdate|Ftime|FtimeMs|FtimeUTC) != 0 {
// 		l.time = time.Now()
// 		if l.flag&FtimeUTC != 0 {
// 			l.time = l.time.UTC()
// 		}
// 		if l.flag&(Fyear|Fdate) != 0 {
// 			year, month, day := l.time.Date()
// 			// if both YYMMDD and YYYYMMDD is given, YYYYMMDD will be used
// 			if l.flag&Fyear != 0 {
// 				itoa(buf, year, 4, '/')
// 			}
// 			// MMDD will be always added ass it's a common denominator of
// 			// FdateYYMMDD|Fyear|Fdate
// 			itoa(buf, int(month), 2, '/')
// 			itoa(buf, day, 2, ' ')
// 		}
// 		if l.flag&(Ftime|FtimeMs|FtimeUTC) != 0 {
// 			hour, min, sec := l.time.Clock()
// 			itoa(buf, hour, 2, ':')
// 			itoa(buf, min, 2, ':')
// 			if l.flag&FtimeMs != 0 {
// 				itoa(buf, sec, 2, '.')
// 				itoa(buf, l.time.Nanosecond()/1e3, 6, ' ')
// 			} else {
// 				itoa(buf, sec, 2, ' ')
// 			}
// 		}
// 	}
// 	// Add prefix
// 	if l.flag&Fprefix != 0 {
// 		*buf = append(*buf, l.prefix...)
// 	}
//
// 	// Add log lvl if lvl is to shown and valid range (0-6) where 0 will not show lvl prefix.
// 	if l.flag&Flevel != 0 && lvl < 7 {
// 		*buf = append(*buf, l.levelString[lvl]...)
// 	}
//
// 	if tag != 0 {
// 		// *buf = append(*buf, `"tag":[`...)
// 		firstItem := true
// 		for i := 0; i < l.logTagIssued; i++ {
// 			if tag&(1<<i) != 0 {
// 				if firstItem {
// 					firstItem = false
// 				} else {
// 					*buf = append(*buf, `/`...)
// 				}
// 				*buf = append(*buf, l.logTagString[i]...)
// 			}
// 		}
// 		*buf = append(*buf, "; "...)
// 	}
// }

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
