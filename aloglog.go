package alog

import "time"

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

	if !l.lvtag.filter.check(lvl, tag) || (lenMsg == 0 && lenMsgb == 0 && lenA == 0) {
		return
	}

	firstItem := true

	s := l.pool.Get().(*poolbuf)

	if l.flag&Fprefix != 0 {
		s.bufHeader = l.fmtr.Begin(s.bufHeader, l.prefix)
	} else {
		s.bufHeader = l.fmtr.Begin(s.bufHeader, nil)
	}

	if l.flag&(FtimeUnix|FtimeUnixMs) != 0 {
		l.time = time.Now()

		if l.flag&FtimeUnixMs != 0 {
			if !firstItem {
				s.bufHeader = l.fmtr.Space(s.bufHeader)
			}
			s.bufHeader = l.fmtr.LogTimeUnixMs(s.bufHeader, l.time)
		} else {
			if !firstItem {
				s.bufHeader = l.fmtr.Space(s.bufHeader)
			}
			s.bufHeader = l.fmtr.LogTimeUnix(s.bufHeader, l.time)
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
				s.bufHeader = l.fmtr.Space(s.bufHeader)
			}
			firstItem = false
			s.bufHeader = l.fmtr.LogTimeDate(s.bufHeader, l.time)
		}
		if l.flag&FdateDay != 0 {
			if !firstItem {
				s.bufHeader = l.fmtr.Space(s.bufHeader)
			}
			firstItem = false
			s.bufHeader = l.fmtr.LogTimeDay(s.bufHeader, l.time)
		}

		if l.flag&Ftime != 0 {
			if !firstItem {
				s.bufHeader = l.fmtr.Space(s.bufHeader)
			}
			s.bufHeader = l.fmtr.LogTime(s.bufHeader, l.time)
		}

		firstItem = false
	}

	if l.flag&Flevel != 0 {
		if !firstItem {
			s.bufHeader = l.fmtr.Space(s.bufHeader)
		}
		s.bufHeader = l.fmtr.LogLevel(s.bufHeader, lvl)
		firstItem = false
	}

	if l.flag&Ftag != 0 {
		if !firstItem {
			s.bufHeader = l.fmtr.Space(s.bufHeader)
		}
		s.bufHeader = l.fmtr.LogTag(s.bufHeader, tag, l.lvtag.tagNames, l.lvtag.numTagIssued)
		firstItem = false
	}

	// print msg
	if lenMsg > 0 {
		if !firstItem {
			s.bufMain = l.fmtr.Space(s.bufMain)
		}
		s.bufMain = l.fmtr.LogMsg(s.bufMain, msg, ';') // suffix is only for text one.
		firstItem = false
	} else if lenMsgb > 0 {
		if !firstItem {
			s.bufMain = l.fmtr.Space(s.bufMain)
		}
		s.bufMain = l.fmtr.LogMsgb(s.bufMain, msgb, ';') // suffix is only for text one.
		firstItem = false
	}

	idxA := lenA - 1
	for i := 0; i < lenA; i += 2 { // 0, 2, 4..
		key, ok := a[i].(string)
		if !ok {
			key = "?badKey?"
		}
		if !firstItem {
			s.bufMain = l.fmtr.Space(s.bufMain)
		}
		firstItem = false
		if i < idxA {
			next := a[i+1]
			switch next.(type) {
			case string:
				s.bufMain = l.fmtr.String(s.bufMain, key, next.(string))
			case nil:
				s.bufMain = l.fmtr.Nil(s.bufMain, key)
			case error:
				s.bufMain = l.fmtr.Error(s.bufMain, key, next.(error))
			case bool:
				s.bufMain = l.fmtr.Bool(s.bufMain, key, next.(bool))
			case int:
				s.bufMain = l.fmtr.Int(s.bufMain, key, next.(int))
			case int8:
				s.bufMain = l.fmtr.Int8(s.bufMain, key, next.(int8))
			case int16:
				s.bufMain = l.fmtr.Int16(s.bufMain, key, next.(int16))
			case int32:
				s.bufMain = l.fmtr.Int32(s.bufMain, key, next.(int32))
			case int64:
				s.bufMain = l.fmtr.Int64(s.bufMain, key, next.(int64))
			case uint:
				s.bufMain = l.fmtr.Uint(s.bufMain, key, next.(uint))
			case uint8:
				s.bufMain = l.fmtr.Uint8(s.bufMain, key, next.(uint8))
			case uint16:
				s.bufMain = l.fmtr.Uint16(s.bufMain, key, next.(uint16))
			case uint32:
				s.bufMain = l.fmtr.Uint32(s.bufMain, key, next.(uint32))
			case uint64:
				s.bufMain = l.fmtr.Uint64(s.bufMain, key, next.(uint64))
			case float32:
				s.bufMain = l.fmtr.Float32(s.bufMain, key, next.(float32))
			case float64:
				s.bufMain = l.fmtr.Float64(s.bufMain, key, next.(float64))
			case []string:
				s.bufMain = l.fmtr.Strings(s.bufMain, key, next.([]string))
			case []error:
				s.bufMain = l.fmtr.Errors(s.bufMain, key, next.([]error))
			case []bool:
				s.bufMain = l.fmtr.Bools(s.bufMain, key, next.([]bool))
			case []float32:
				s.bufMain = l.fmtr.Float32s(s.bufMain, key, next.([]float32))
			case []float64:
				s.bufMain = l.fmtr.Float64s(s.bufMain, key, next.([]float64))
			case []int:
				s.bufMain = l.fmtr.Ints(s.bufMain, key, next.([]int))
			case []int32:
				s.bufMain = l.fmtr.Int32s(s.bufMain, key, next.([]int32))
			case []int64:
				s.bufMain = l.fmtr.Int64s(s.bufMain, key, next.([]int64))
			case []uint:
				s.bufMain = l.fmtr.Uints(s.bufMain, key, next.([]uint))
			case []uint8:
				s.bufMain = l.fmtr.Uint8s(s.bufMain, key, next.([]uint8))
			case []uint32:
				s.bufMain = l.fmtr.Uint32s(s.bufMain, key, next.([]uint32))
			case []uint64:
				s.bufMain = l.fmtr.Uint64s(s.bufMain, key, next.([]uint64))
			default:
				s.bufMain = l.fmtr.String(s.bufMain, key, "?unsupp?")
			}
		} else {
			s.bufMain = l.fmtr.Nil(s.bufMain, key)
		}
	}

	// any custom func using bufMain should be run here.
	if l.trigfn != nil {
		l.trigfn(lvl, tag, s.bufMain)
	}

	//println("----\ns.bufHeader", string(s.bufHeader), "\n\n")
	//println("----\ns.bufMain", string(s.bufMain), "\n\n")

	// Finalize
	s.bufMain = l.fmtr.End(s.bufMain)
	l.mu.Lock()
	l.out.Write(append(s.bufHeader, s.bufMain...))
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
