package alog

import "strconv"

func (l *Logger) Log(lvl Level, tag Tag, msg string, a ...interface{}) (n int, err error) {
	if !l.check(lvl, tag) || (len(msg) == 0 && len(a) == 0) {
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

func (l *Logger) log_json(lvl Level, tag Tag, msg string, a ...interface{}) {
	isBufSpace := l.isBufSpace()
	if len(msg) != 0 {
		if !isBufSpace {
			l.addComma()
		}
		l.buf = append(l.buf, `"msg":`...)
		l.addEscapeJStr(msg)
		isBufSpace = false
	}

	if len(a) == 0 {
		return
	}

	if !isBufSpace {
		l.buf = append(l.buf, ',')
	}

	isOdd := true // even items to be type cased

	lenA := len(a) - 1

	for i, v := range a {
		if isOdd {
			if _, ok := v.(string); ok {
				l.addEscapeJStr(v.(string))
			} else {
				return
			}
		} else {
			l.addColon()
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
	isMsg := false
	isBufSpace := l.isBufSpace()
	if len(msg) != 0 {
		isMsg = true
		if !isBufSpace {
			l.buf = append(l.buf, ' ')
		}
		l.addEscapeBasic(msg)
	}

	// If not args pairs, then return
	lenA := len(a)
	if lenA == 0 {
		return // todo: if nothing has been added to the buffer, do not finalize, just clear the buf?
	}

	// l.addColon()
	if isMsg {
		l.buf = append(l.buf, ' ', '-')
	}

	// args pair(s) exist
	isOdd := true
	// firstRec := true
	for _, v := range a {
		if isOdd { // todo: get index of `range a`, and %2 instead of using isOdd variable.
			if _, ok := v.(string); ok {
				if isMsg || !isBufSpace {
					l.buf = append(l.buf, ' ')
				}
				// }
				// l.buf = strconv.AppendQuote(l.buf, v.(string))
				// l.addEscapeJStr(v.(string))
				l.addEscapeBasic(v.(string))
				// l.buf = append(l.buf, v.(string)...)
			} else {
				// todo: add a msg like unknown type key
				return
			}
		} else {
			l.buf = append(l.buf, "="...)

			switch v.(type) {
			// Frequently used first
			// Todo: add error type
			case string:
				l.addEscapeJStr(v.(string))
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
				l.buf = append(l.buf, `"?unsupp?"`...) // todo: move this to global
			}
		}
		isOdd = !isOdd // toggle
	}
	if lenA%2 == 1 {
		l.buf = append(l.buf, "=null"...)
	}
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
