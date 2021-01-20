package alog

import "strconv"

func (l *Logger) header(lvl Level, tag Tag) {
	if l.flag&Fjson != 0 {
		l.header_json(lvl, tag)
	} else {
		l.header_text(lvl, tag)
	}
}

func (l *Logger) log_json(lvl Level, tag Tag, msg string, a ...interface{}) {
	if len(msg) != 0 {
		if msg != "" {
			l.buf = append(l.buf, `,"msg":`...)
			l.buf = strconv.AppendQuote(l.buf, msg)
		}
	}

	if len(a) == 0 {
		return
	}

	l.buf = append(l.buf, ","...)

	//l.bufFmt = l.bufFmt[:0]
	isOdd := true // even items to be type cased

	lenA := len(a) - 1

	for i, v := range a {
		if isOdd {
			if _, ok := v.(string); ok {
				//l.bufFmt = append(l.bufFmt, '"')
				//l.bufFmt = append(l.bufFmt, v.(string)...)
				//l.bufFmt = append(l.bufFmt, '"')
				l.buf = append(l.buf, '"')
				l.buf = append(l.buf, v.(string)...)
				l.buf = append(l.buf, '"')
			} else {
				return
			}
		} else {
			//l.bufFmt = append(l.bufFmt, ":"...)
			l.buf = append(l.buf, ":"...)
			switch v.(type) {
			// Frequently used first
			// Todo: add error type
			case string:
				l.buf = strconv.AppendQuote(l.buf, v.(string))
			case int:
				l.buf = strconv.AppendInt(l.buf, int64(v.(int)), 10)
			case int64:
				l.buf = strconv.AppendInt(l.buf, v.(int64), 10)
			case bool:
				l.buf = strconv.AppendBool(l.buf, v.(bool))
			case uint:
				l.buf = strconv.AppendInt(l.buf, int64(v.(uint)), 10)
			case uint64:
				l.buf = strconv.AppendUint(l.buf, v.(uint64), 10)
			case float32:
				l.buf = strconv.AppendFloat(l.buf, float64(v.(float32)), 'f', -1, 32)
			case float64:
				l.buf = strconv.AppendFloat(l.buf, v.(float64), 'f', -1, 64)
			default:
				l.buf = strconv.AppendQuote(l.buf, unsuppTypes)
			}
			if lenA > i {
				l.buf = append(l.buf, ',')
			}
		}
		isOdd = !isOdd // toggle
	}
	// This is where incorrect pairs were received such as:
	// 		"age", 17, "name", "gon", "job"
	//		(where a value for job is missing)
	if lenA%2 == 0 { // lenA is index, so should be odd number if right pair
		l.buf = append(l.buf, ":null"...)
	}
}

func (l *Logger) log_text(lvl Level, tag Tag, msg string, a ...interface{}) {
	if len(msg) != 0 {
		{
			l.buf = append(l.buf, `msg=`...)
			l.buf = strconv.AppendQuote(l.buf, msg)
		}
	}

	if len(a) == 0 {
		return
	}

	l.buf = append(l.buf, ", "...)

	isOdd := true
	firstRec := true
	for _, v := range a {
		if isOdd {
			if _, ok := v.(string); ok {
				if firstRec {
					firstRec = false
				} else {
					l.buf = append(l.buf, ", "...)
				}

				l.buf = append(l.buf, v.(string)...)
			} else {
				return
			}
		} else {
			{
				l.buf = append(l.buf, "="...)
			}
			switch v.(type) {
			// Frequently used first
			// Todo: add error type
			case string:
				l.buf = strconv.AppendQuote(l.buf, v.(string))
			case int:
				l.buf = strconv.AppendInt(l.buf, int64(v.(int)), 10)
			case int64:
				l.buf = strconv.AppendInt(l.buf, v.(int64), 10)
			case bool:
				l.buf = strconv.AppendBool(l.buf, v.(bool))
			case uint:
				l.buf = strconv.AppendInt(l.buf, int64(v.(uint)), 10)
			case uint64:
				l.buf = strconv.AppendUint(l.buf, v.(uint64), 10)
			case float32:
				l.buf = strconv.AppendFloat(l.buf, float64(v.(float32)), 'f', -1, 32)
			case float64:
				l.buf = strconv.AppendFloat(l.buf, v.(float64), 'f', -1, 64)
			default:
				l.buf = append(l.buf, unsuppType...)
			}
		}
		isOdd = !isOdd // toggle
	}
}

func (l *Logger) Log(lvl Level, tag Tag, msg string, a ...interface{}) (n int, err error) {
	if !l.check(lvl, tag) {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.header(lvl, tag)
	if l.flag&Fjson != 0 {
		l.log_json(lvl, tag, msg, a...)
	} else {
		l.log_text(lvl, tag, msg, a...)
	}
	return l.finalize()
}

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
