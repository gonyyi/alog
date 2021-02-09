package alog

import (
	"io"
	"sync"
	"time"
)

// Logger is a main struct for logger.
// Available methods are:
//    Simple: Print(), Trace(), Debug(), Info(), Warn(), Error(), Fatal()
//    Format: Printf(), Tracef(), Debugf(), Infof(), Warnf(), Errorf(), Fatalf()
//    Other:  NewPrint(), NewWriter()
type Logger struct {
	time       time.Time
	formatFlag Format

	pool   sync.Pool
	fmtr   Fmtr
	hookFn HookFn

	lvtag tagger // lvtag to replace 5 items below

	out      io.Writer
	mu       sync.Mutex
	useMutex bool

	// buf    []byte // buf is a main buffer; reset per each log entry
	prefix []byte // prefix will be stored as a byte slice.
}

// New function creates new logger.
// This takes an output writer for its argument. If nil is given, it will discard logs.
func New(output io.Writer) *Logger {
	// If output is given as `nil`, it will use io.discard as a default.
	if output == nil {
		output = discard
	}

	// Creating new logger object and returns pointer to logger.
	// Default value will be set here. If a user uses *alog.SetFormat to provoke
	// unless specifically set certain value, it will not be overwritten.
	// eg. If a user called SetFormat with other config formatFlag except the Level, then the log
	//     Level will not be changed. Therefore default Level should be defined here.
	l := &Logger{
		pool: sync.Pool{
			New: func() interface{} {
				return newPoolbuf(256, 512)
			},
		},
		fmtr:       &fmtText{},
		out:        output,
		prefix:     []byte(""), // prefix will be saved as a byte slice to prevent need to be converted later.
		formatFlag: Fdefault,   // default formatFlag is given
		// buf:      make([]byte, 1024),
	}
	l.lvtag.filter.lvl = Linfo // default logging Level to INFO

	return l
}

// Do can run a function(s) that were created by a user
// An example would be set Level prefix with ANSI color
// or series of frequentyly used settings.
// planned for v0.2.1 release.
func (l *Logger) Do(fn ...func(*Logger)) *Logger {
	for _, f := range fn {
		f(l)
	}
	return l
}

// SetMutex will set mutex when used for writing to writer.
func (l *Logger) SetMutex(on bool) *Logger {
	l.useMutex = on
	return l
}

// GetTag takes a wTag name and returns a wTag if found.
func (l *Logger) GetTag(name string) Tag {
	return l.lvtag.mustGetTag(name)
}

// Output returns the output destination for the logger.
func (l *Logger) Output() io.Writer {
	return l.out
}

// SetOutput can redefined the output after logger has been created.
// If output is nil, the logger will set it to ioutil.Discard instead.
func (l *Logger) SetOutput(output io.Writer) *Logger {
	l.mu.Lock()
	if output == nil {
		l.out = discard
	} else {
		l.out = output
	}
	l.mu.Unlock()
	return l
}

// Format will return current format flag.
// This can be modified and set again using SetFormat method.
func (l *Logger) Format() Format {
	return l.formatFlag
}

// SetFormat will reconfigure the logger after it has been created.
// This will first copy formatFlag into *alog.formatFlag, and sets few that
// need additional parsing.
func (l *Logger) SetFormat(flag Format) *Logger {
	l.mu.Lock()

	// Previous had a JSON formatFlag on, but anymore, then use text formatter
	// If previous didn't have JSON formatFlag, but now do, then use json formatter
	if l.formatFlag&Fjson != 0 && flag&Fjson == 0 {
		l.SetFormatter(&fmtText{})
	} else if l.formatFlag&Fjson == 0 && flag&Fjson != 0 {
		l.SetFormatter(&fmtJSON{})
	}

	l.formatFlag = flag

	l.mu.Unlock()
	return l
}

// SetFormatter will set logger formatter. Without this setting,
// the logger's default will be text formatter (fmtText).
// If switching between default JSON and TEXT formatter, it should
// be done by ModFormat() or SetFormat(). This SetFormatter is
// mainly for a customized formatter.
func (l *Logger) SetFormatter(fmtr Fmtr) *Logger {
	if fmtr != nil {
		l.fmtr = fmtr
		return l
	}
	l.fmtr = fmtText{}
	return l
}

// SetPrefix can redefine prefix after the logger has been created.
func (l *Logger) SetPrefix(s string) *Logger {
	l.mu.Lock()
	l.prefix = []byte(s)
	l.mu.Unlock()
	return l
}

// SetFilter will define what level or tags to show.
// Integer 0 can be used, and when it's used, it will not filter anything.
func (l *Logger) SetFilter(lv Level, tags Tag) *Logger {
	l.lvtag.filter.fn = nil
	l.lvtag.filter.lvl = lv
	l.lvtag.filter.tag = tags
	return l
}

// SetFilterFn can control more precisely by taking a FilterFn.
func (l *Logger) SetFilterFn(fn FilterFn) *Logger {
	l.lvtag.filter.fn = fn
	return l
}

// SetHookFn will create a hookFn that works addition to filter.
// Example would be log everything but for HTTP request tags,
// also write it to a file.
func (l *Logger) SetHookFn(fn HookFn) {
	l.hookFn = fn
}

// Log is a main logging method.
func (l *Logger) Log(lvl Level, tag Tag, msg string, a ...interface{}) (n int, err error) {
	return l.log(lvl, tag, msg, nil, a...)
}

// LogIf will check and log error if exist (not nil)
// For instance, when running multiple lines of error check
// This can save error checking.
// added as v0.1.6c3, 12/30/2020
func (l *Logger) LogIf(e error, lvl Level, tag Tag, msg string) {
	if e != nil {
		l.log(lvl, tag, msg, nil, "err", e.Error())
	}
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
	s.bufHeader = s.bufHeader[:0]
	s.bufMain = s.bufMain[:0]

	if l.formatFlag&Fprefix != 0 {
		s.bufHeader = l.fmtr.Begin(s.bufHeader, l.prefix)
	} else {
		s.bufHeader = l.fmtr.Begin(s.bufHeader, nil)
	}

	if l.formatFlag&(FtimeUnix|FtimeUnixMs) != 0 {
		l.time = time.Now()

		firstItem = false
		if l.formatFlag&FtimeUnixMs != 0 {
			s.bufHeader = l.fmtr.LogTimeUnixMs(s.bufHeader, l.time)
		} else {
			s.bufHeader = l.fmtr.LogTimeUnix(s.bufHeader, l.time)
		}

	} else if l.formatFlag&(Fdate|FdateDay|Ftime|FtimeUTC) != 0 {
		// at least one item will be printed here, so just check once.
		l.time = time.Now()
		if l.formatFlag&FtimeUTC != 0 {
			l.time = l.time.UTC()
		}

		if l.formatFlag&Fdate != 0 {
			firstItem = false
			s.bufHeader = l.fmtr.LogTimeDate(s.bufHeader, l.time)
		}
		if l.formatFlag&FdateDay != 0 {
			if !firstItem {
				s.bufHeader = l.fmtr.Space(s.bufHeader)
			}
			firstItem = false
			s.bufHeader = l.fmtr.LogTimeDay(s.bufHeader, l.time)
		}
		if l.formatFlag&Ftime != 0 {
			if !firstItem {
				s.bufHeader = l.fmtr.Space(s.bufHeader)
			}
			s.bufHeader = l.fmtr.LogTime(s.bufHeader, l.time)
		}
		firstItem = false
	}

	if l.formatFlag&Flevel != 0 {
		if !firstItem {
			s.bufHeader = l.fmtr.Space(s.bufHeader)
		}
		s.bufHeader = l.fmtr.LogLevel(s.bufHeader, lvl)
		firstItem = false
	}

	if l.formatFlag&Ftag != 0 {
		if !firstItem {
			s.bufHeader = l.fmtr.Space(s.bufHeader)
		}
		s.bufHeader = l.fmtr.LogTag(s.bufHeader, tag, &l.lvtag.tagNames, l.lvtag.numTagIssued)
		firstItem = false
	}

	// print msg
	if lenMsg > 0 {
		if !firstItem {
			s.bufHeader = l.fmtr.Space(s.bufHeader) // add space to very previous step
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
			case *[]string:
				s.bufMain = l.fmtr.Strings(s.bufMain, key, next.(*[]string))
			case *[]error:
				s.bufMain = l.fmtr.Errors(s.bufMain, key, next.(*[]error))
			case *[]bool:
				s.bufMain = l.fmtr.Bools(s.bufMain, key, next.(*[]bool))
			case *[]float32:
				s.bufMain = l.fmtr.Float32s(s.bufMain, key, next.(*[]float32))
			case *[]float64:
				s.bufMain = l.fmtr.Float64s(s.bufMain, key, next.(*[]float64))
			case *[]int:
				s.bufMain = l.fmtr.Ints(s.bufMain, key, next.(*[]int))
			case *[]int32:
				s.bufMain = l.fmtr.Int32s(s.bufMain, key, next.(*[]int32))
			case *[]int64:
				s.bufMain = l.fmtr.Int64s(s.bufMain, key, next.(*[]int64))
			case *[]uint:
				s.bufMain = l.fmtr.Uints(s.bufMain, key, next.(*[]uint))
			case *[]uint8:
				s.bufMain = l.fmtr.Uint8s(s.bufMain, key, next.(*[]uint8))
			case *[]uint32:
				s.bufMain = l.fmtr.Uint32s(s.bufMain, key, next.(*[]uint32))
			case *[]uint64:
				s.bufMain = l.fmtr.Uint64s(s.bufMain, key, next.(*[]uint64))
			default:
				s.bufMain = l.fmtr.String(s.bufMain, key, "?unsupp?")
			}
		} else {
			s.bufMain = l.fmtr.Nil(s.bufMain, key)
		}
	}

	// any custom func using bufMain should be run here.
	if l.hookFn != nil {
		l.hookFn(lvl, tag, s.bufMain)
	}

	// Finalize
	s.bufMain = l.fmtr.End(s.bufMain)

	// Use mutex only when necessary
	if l.useMutex {
		l.mu.Lock()
		l.out.Write(append(s.bufHeader, s.bufMain...))
		l.mu.Unlock()
	} else {
		l.out.Write(append(s.bufHeader, s.bufMain...))
	}

	s.reset() // reset buffer to prevent potentially large one left in the pool
	l.pool.Put(s)

	return 0, nil
}

// Close will call .Close() method if supported
func (l *Logger) Close() error {
	if c, ok := l.out.(io.Closer); ok && c != nil {
		return c.Close()
	}
	return nil
}

// NewWriter takes a Level and a Tag and create an Alog writer (SubWriter)
// that is compatible with io.Writer interface. This can be used as a
// logger hookFn.
func (l *Logger) NewWriter(lvl Level, tag Tag) *SubWriter {
	return &SubWriter{
		l:      l,
		wLevel: lvl,
		wTag:   tag,
	}
}

func (l *Logger) Trace(tag Tag, msg string, a ...interface{}) {
	l.log(Ltrace, tag, msg, nil, a...)
}
func (l *Logger) Debug(tag Tag, msg string, a ...interface{}) {
	l.log(Ldebug, tag, msg, nil, a...)
}
func (l *Logger) Info(tag Tag, msg string, a ...interface{}) {
	l.log(Linfo, tag, msg, nil, a...)
}
func (l *Logger) Warn(tag Tag, msg string, a ...interface{}) {
	l.log(Lwarn, tag, msg, nil, a...)
}
func (l *Logger) Error(tag Tag, msg string, a ...interface{}) {
	l.log(Lerror, tag, msg, nil, a...)
}
func (l *Logger) Fatal(tag Tag, msg string, a ...interface{}) {
	l.log(Lfatal, tag, msg, nil, a...)
}
