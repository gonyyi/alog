package alog

import "time"

// check will check if Level and Tag given is good to be printed.
// Eg. if setting is Level INFO, Tag USER, then
//     any log Level below INFO shouldn't be printed.
//     Also, any Tag other than USER shouldn't be printed either.
func (l *Logger) check(lvl Level, tag Tag) bool {
	switch {
	case l.logFn != nil: // filterFn has the highest order if set.
		return l.logFn(lvl, tag)
	case l.logLevel > lvl: // if defaultLevel is below defaultLevel limit, the do not print
		return false
	case l.logTag != 0 && l.logTag&tag == 0: // if filterTag is set but Tag is not matching, then do not print
		return false
	default:
		return true
	}
}

// Log is a main method for log.
func (l *Logger) Log(lvl Level, tag Tag, msg string, a ...interface{}) (n int, err error) {
	return l.log(lvl, tag, msg, nil, a...)
}

// logb was created for SubWriter to reduce converting string to byte.
func (l *Logger) logb(lvl Level, tag Tag, msg []byte) (n int, err error) {
	return l.log(lvl, tag, "", msg)
}

// log method will take both stringed msg and []byte msgb assume only one will be used.
func (l *Logger) log(lvl Level, tag Tag, msg string, msgb []byte, a ...interface{}) (n int, err error) {
	lenA, lenMsg, lenMsgb := len(a), len(msg), len(msgb)

	if !l.check(lvl, tag) || (lenMsg == 0 && lenMsgb == 0 && lenA == 0) {
		return
	}

	firstItem := true

	s := l.pool.Get().(*alogItem)

	if l.flag&Fprefix != 0 {
		s.buf = l.fmtr.Begin(s.buf, l.prefix)
	} else {
		s.buf = l.fmtr.Begin(s.buf, nil)
	}

	if l.flag&(FtimeUnix|FtimeUnixMs) != 0 {
		l.time = time.Now()

		if l.flag&FtimeUnixMs != 0 {
			if !firstItem {
				s.buf = l.fmtr.Space(s.buf)
			}
			s.buf = l.fmtr.LogTimeUnixMs(s.buf, l.time)
		} else {
			if !firstItem {
				s.buf = l.fmtr.Space(s.buf)
			}
			s.buf = l.fmtr.LogTimeUnix(s.buf, l.time)
		}

		firstItem = false
	} else if l.flag&(Fdate|FdateDay|Ftime|FtimeUTC) != 0 {
		// at least one item will be printed here, so just check once.
		l.time = time.Now()
		if l.flag&FtimeUTC != 0 {
			l.time = l.time.UTC()
		}

		if l.flag&Fdate != 0 {
			if !firstItem {
				s.buf = l.fmtr.Space(s.buf)
			}
			firstItem = false
			s.buf = l.fmtr.LogTimeDate(s.buf, l.time)
		}
		if l.flag&FdateDay != 0 {
			if !firstItem {
				s.buf = l.fmtr.Space(s.buf)
			}
			firstItem = false
			s.buf = l.fmtr.LogTimeDay(s.buf, l.time)
		}

		if l.flag&Ftime != 0 {
			if !firstItem {
				s.buf = l.fmtr.Space(s.buf)
			}
			s.buf = l.fmtr.LogTime(s.buf, l.time)
		}

		firstItem = false
	}

	if l.flag&Flevel != 0 {
		if !firstItem {
			s.buf = l.fmtr.Space(s.buf)
		}
		s.buf = l.fmtr.LogLevel(s.buf, lvl)
		firstItem = false
	}

	if l.flag&Ftag != 0 {
		if !firstItem {
			s.buf = l.fmtr.Space(s.buf)
		}
		s.buf = l.fmtr.LogTag(s.buf, tag, l.lvtag.tagNames, l.lvtag.numTagIssued)
		firstItem = false
	}

	// print msg
	if lenMsg > 0 {
		if !firstItem {
			s.buf = l.fmtr.Space(s.buf)
		}
		s.buf = l.fmtr.LogMsg(s.buf, msg, ';') // suffix is only for text one.
		firstItem = false
	} else if lenMsgb > 0 {
		if !firstItem {
			s.buf = l.fmtr.Space(s.buf)
		}
		s.buf = l.fmtr.LogMsgb(s.buf, msgb, ';') // suffix is only for text one.
		firstItem = false
	}

	idxA := lenA - 1
	for i := 0; i < lenA; i += 2 { // 0, 2, 4..
		key, ok := a[i].(string)
		if !ok {
			key = "?badKey?"
		}
		if !firstItem {
			s.buf = l.fmtr.Space(s.buf)
		}
		firstItem = false
		if i < idxA {
			next := a[i+1]
			switch next.(type) {
			case string:
				s.buf = l.fmtr.String(s.buf, key, next.(string))
			case nil:
				s.buf = l.fmtr.Nil(s.buf, key)
			case error:
				s.buf = l.fmtr.Error(s.buf, key, next.(error))
			case bool:
				s.buf = l.fmtr.Bool(s.buf, key, next.(bool))
			case int:
				s.buf = l.fmtr.Int(s.buf, key, next.(int))
			case int8:
				s.buf = l.fmtr.Int8(s.buf, key, next.(int8))
			case int16:
				s.buf = l.fmtr.Int16(s.buf, key, next.(int16))
			case int32:
				s.buf = l.fmtr.Int32(s.buf, key, next.(int32))
			case int64:
				s.buf = l.fmtr.Int64(s.buf, key, next.(int64))
			case uint:
				s.buf = l.fmtr.Uint(s.buf, key, next.(uint))
			case uint8:
				s.buf = l.fmtr.Uint8(s.buf, key, next.(uint8))
			case uint16:
				s.buf = l.fmtr.Uint16(s.buf, key, next.(uint16))
			case uint32:
				s.buf = l.fmtr.Uint32(s.buf, key, next.(uint32))
			case uint64:
				s.buf = l.fmtr.Uint64(s.buf, key, next.(uint64))
			case float32:
				s.buf = l.fmtr.Float32(s.buf, key, next.(float32))
			case float64:
				s.buf = l.fmtr.Float64(s.buf, key, next.(float64))
			case []string:
				s.buf = l.fmtr.Strings(s.buf, key, next.([]string))
			case []error:
				s.buf = l.fmtr.Errors(s.buf, key, next.([]error))
			case []bool:
				s.buf = l.fmtr.Bools(s.buf, key, next.([]bool))
			case []float32:
				s.buf = l.fmtr.Float32s(s.buf, key, next.([]float32))
			case []float64:
				s.buf = l.fmtr.Float64s(s.buf, key, next.([]float64))
			case []int:
				s.buf = l.fmtr.Ints(s.buf, key, next.([]int))
			case []int32:
				s.buf = l.fmtr.Int32s(s.buf, key, next.([]int32))
			case []int64:
				s.buf = l.fmtr.Int64s(s.buf, key, next.([]int64))
			case []uint:
				s.buf = l.fmtr.Uints(s.buf, key, next.([]uint))
			case []uint8:
				s.buf = l.fmtr.Uint8s(s.buf, key, next.([]uint8))
			case []uint32:
				s.buf = l.fmtr.Uint32s(s.buf, key, next.([]uint32))
			case []uint64:
				s.buf = l.fmtr.Uint64s(s.buf, key, next.([]uint64))
			default:
				s.buf = l.fmtr.String(s.buf, key, "?unsupp?")
			}
		} else {
			s.buf = l.fmtr.Nil(s.buf, key)
		}
	}

	s.buf = l.fmtr.End(s.buf)

	l.mu.Lock()
	l.out.Write(s.buf)
	l.mu.Unlock()
	s.reset() // reset buffer to prevent potentially large one left in the pool
	l.pool.Put(s)

	return 0, nil
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
