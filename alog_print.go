package alog

import "strconv"

func (l *Logger) Trace(tag Tag, msg string, a ...interface{}) {
	l.Log(Ltrace, tag, msg, a...)
}
func (l *Logger) Debug(tag Tag, msg string, a ...interface{}) {
	l.Log(Ldebug, tag, msg, a...)
}
func (l *Logger) Info(tag Tag, msg string, a ...interface{}) {
	l.Log(Linfo, tag, msg, a...)
}
func (l *Logger) Warn(tag Tag, msg string, a ...interface{}) {
	l.Log(Lwarn, tag, msg, a...)
}
func (l *Logger) Error(tag Tag, msg string, a ...interface{}) {
	l.Log(Lerror, tag, msg, a...)
}
func (l *Logger) Fatal(tag Tag, msg string, a ...interface{}) {
	l.Log(Lfatal, tag, msg, a...)
}

func (l *Logger) Log(lvl Level, tag Tag, msg string, a ...interface{}) (n int, err error) {
	if !l.check(lvl, tag) {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.header(&l.buf, lvl, tag)
	// l.buf = append(l.buf, msg...)
	if l.flag&Fjson != 0 {
		if msg != "" {
			l.buf = append(l.buf, `,"msg":"`...)
			l.buf = append(l.buf, msg...)
			l.buf = append(l.buf, '"')
		}
	} else {
		l.buf = append(l.buf, msg...)
	}

	if len(a) == 0 {
		return l.finalize()
	}

	if l.flag&Fjson != 0 {
		l.buf = append(l.buf, ","...)
	} else {
		l.buf = append(l.buf, ", "...)
	}

	l.bufFmt = l.bufFmt[:0]
	isOdd := true
	for _, v := range a {
		if isOdd {
			if _, ok := v.(string); ok {
				if l.flag&Fjson != 0 {
					l.bufFmt = append(l.bufFmt, '"')
					l.bufFmt = append(l.bufFmt, v.(string)...)
					l.bufFmt = append(l.bufFmt, '"')
				} else {
					l.bufFmt = append(l.bufFmt, v.(string)...)
				}
			} else {
				return
			}
		} else {
			if l.flag&Fjson != 0 {
				l.bufFmt = append(l.bufFmt, ":"...)
			} else {
				l.bufFmt = append(l.bufFmt, "="...)
			}
			switch v.(type) {
			// Frequently used first
			// Todo: add error type
			case string:
				l.bufFmt = strconv.AppendQuote(l.bufFmt, v.(string))
			case int:
				l.bufFmt = strconv.AppendInt(l.bufFmt, int64(v.(int)), 10)
			case int64:
				l.bufFmt = strconv.AppendInt(l.bufFmt, v.(int64), 10)
			case bool:
				l.bufFmt = strconv.AppendBool(l.bufFmt, v.(bool))
			case uint:
				l.bufFmt = strconv.AppendInt(l.bufFmt, int64(v.(uint)), 10)
			case uint64:
				l.bufFmt = strconv.AppendUint(l.bufFmt, v.(uint64), 10)
			case float32:
				l.bufFmt = strconv.AppendFloat(l.bufFmt, float64(v.(float32)), 'f', -1, 32)
			case float64:
				l.bufFmt = strconv.AppendFloat(l.bufFmt, v.(float64), 'f', -1, 64)

			default:
				if l.flag&Fjson != 0 {
					l.bufFmt = strconv.AppendQuote(l.bufFmt, unsuppTypes)
				} else {
					l.bufFmt = append(l.bufFmt, unsuppType...)
				}
				// l.bufFmt = append(l.bufFmt, '"')
				// l.bufFmt = append(l.bufFmt, v.(string)...)
				// l.bufFmt = append(l.bufFmt, '"')
			}
			if l.flag&Fjson != 0 {
				l.bufFmt = append(l.bufFmt, ","...)
			} else {
				l.bufFmt = append(l.bufFmt, ", "...)
			}
		}
		isOdd = !isOdd // toggle
	}
	if l.flag&Fjson != 0 {
		if len(l.bufFmt) != 0 && l.bufFmt[len(l.bufFmt)-1] == ',' {
			l.buf = append(l.buf, l.bufFmt[:len(l.bufFmt)-1]...)
		} else {
			l.buf = append(l.buf, l.bufFmt...)
		}
	} else {
		if len(l.bufFmt) > 1 && l.bufFmt[len(l.bufFmt)-2] == ',' {
			l.buf = append(l.buf, l.bufFmt[:len(l.bufFmt)-2]...)
		} else {
			l.buf = append(l.buf, l.bufFmt...)
		}
	}

	return l.finalize()
}
