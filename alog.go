// (c) 2020 Gon Y Yi. <https://gonyyi.com>

// version 0.4 candidate

package alog

import (
	"io"
	"os"
	"sync"
	"time"
)

const (
	// newline constant provides byte of newline, so it can be usd right away.
	newline = byte('\n')
	// noTag will reset tag if already set
	noTag Tag = 0
)

// unsuppType is a slice of byte and will be used when unknown formats string is being used in
// any formatted prints such as `outputf`, `infof`, `debugf`, etc. This is pre-converted to
// a byte slice and reused to save process time.
var unsuppType = []byte("{??}")

// flags a bit-flag flag options that is used for variety of configuration.
type flags uint32

const (
	// Fall for all options on
	Fall = flags(^uint32(0))
	// Fnone for all options off
	Fnone = flags(uint32(0))

	// Ftime will show HH:MM:SS formats such as 05:02:03
	Ftime flags = 1 << iota
	// FtimeMs will show millisecond in its time such as 10:12:13.1234
	FtimeMs
	// FtimeUTC will show UTC time formats
	FtimeUTC
	// FdateMMDD will show 01/02 date formats.
	FdateMMDD
	// FdateYYMMDD will show 06/01/02 date formats.
	FdateYYMMDD
	// FdateYYYYMMDD will show 2006/01/02 date formats.
	FdateYYYYMMDD
	// Fprefix will show prefix when printing log message
	Fprefix
	// Flevel show Level in the log messsage.
	Flevel
	// Fnewline will enable newlines within the log (v0.1.4)
	Fnewline
	// Fdefault will show month/day with time, and Level of logging.
	Fdefault = FdateYYYYMMDD | Ftime | Flevel
)

// LevelPrefix is a bit-flag used for different Level of log activity:
// - Ltrace: detailed debugging Level
// - Ldebug: general debugging Level
// - Linfo: information Level
// - Lwarn: warning
// - Lerror: error
// - Lfatal: fatal, the process will can be terminated
type Level uint8

const (
	// Ltrace shows trace Level, thee most detailed debugging Level.
	// This will show everything.
	Ltrace Level = iota + 1
	// Ldebug shows debug Level or higher
	Ldebug
	// Linfo shows information Level or higher
	Linfo
	// Lwarn shows warning Level or higher
	Lwarn
	// Lerror shows error Level or higher
	Lerror
	// Lfatal shows fatal Level or higher
	// If Fatal() or Fatalf() is called, it will exit the process with
	// os.Exit(1)
	Lfatal
)

// Logger is a main struct for logger.
// Available methods are:
//    Simple: Print(), Trace(), Debug(), Info(), Warn(), Error(), Fatal()
//    Format: Printf(), Tracef(), Debugf(), Infof(), Warnf(), Errorf(), Fatalf()
//    Other:  NewPrint(), NewWriter()
type Logger struct {
	time time.Time
	flag flags

	// logFn is a customizable function space and supercedes builtin Level and Tag filters if set.
	logFn        func(Level, Tag) bool
	logLevel     Level
	logTag       Tag
	logTagIssued Tag // stores last Tag issued. Whenever NewTag is called, this value will be doubled.

	out io.Writer
	mu  sync.Mutex

	// There are two buffers used. Both `buf` and `bufFormat` are being used regardless
	// of bufUseBuffer setting, however, if bufUseBuffer is false, this buffer will be
	// flushed for each log.
	buf []byte // main buffer; reset per each log entry

	// bufFormat is a buffer strictly used only for formatting - such as printing
	// date, time, prefix etc; and will be copied to `buf` (main buffer)
	// bufFormat []byte // for formatting only

	// Prefix will be stored as a byte slice.
	prefix []byte

	// levelString is an array of byte slice that stores what prefix per each log logLevel
	// will be used. Eg. "[DEBUG]", etc.
	levelString [7][]byte
}

// New function creates new logger. This takes an output writer for its argument (v0.2.0 change)
// All methods with suffix "Set" returns `*Logger`, therefore can be used together with `*Logger.New`.
// This is a v0.2.0 change that broke the backward compatibility, however, most of the time, people
// don't set logger prefix, also uses basic default setting. Therefore it's bit cumbersome to require
// two (prefix, flag), often, unused parameters. Duct-taping can be done following way:
// 1. alog.New(nil).SetOutput(os.Stderr) // initially set discard for output but overridden to os.Stderr
// 2. alog.New(os.Stderr).SetPrefix("TestLog: ").SetLogLevel(alog.Linfo) // set prefix and Level
// 3. alog.New(os.Stderr).SetPrefix("TestLog: ").SetFlag(alog.Fdefault|alog.FtimeUTC)
func New(output io.Writer) *Logger {
	// If output is given as `nil`, it will use io.discard as a default.
	if output == nil {
		output = discard
	}

	// Creating new logger object and returns pointer to logger.
	// Default value will be set here. If a user uses *alog.SetFlag to provoke
	// unless specifically set certain value, it will not be overwritten.
	// eg. If a user called SetFlag with other config flag except the Level, then the log
	//     Level will not be changed. Therefore default Level should be defined here.
	l := &Logger{
		out:      output,
		prefix:   []byte(""), // prefix will be saved as a byte slice to prevent need to be converted later.
		logLevel: Linfo,      // default logging Level to INFO
		flag:     Fdefault,   // default flag is given
	}

	// Default prefixes for each Level. This can be changed by a user using *alog.SetLevelPrefix()
	l.SetLevelPrefix("[TRC] ", "[DBG] ", "[INF] ", "[WRN] ", "[ERR] ", "[FTL] ")
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

// SetPrefix can redefine prefix after the logger has been created.
func (l *Logger) SetPrefix(s string) *Logger {
	l.mu.Lock()
	l.prefix = []byte(s)
	l.mu.Unlock()
	return l
}

// SetFlag will reconfigure the logger after it has been created.
// This will first copy flag into *alog.flag, and sets few that
// need additional parsing.
func (l *Logger) SetFlag(flag flags) *Logger {
	l.mu.Lock()
	l.flag = flag
	l.mu.Unlock()
	return l
}

// SetLevelPrefix will set different prefixes for each levelled log messages.
// eg. "[DEBUG]".
func (l *Logger) SetLevelPrefix(trace, debug, info, warn, error, fatal string) *Logger {
	l.levelString[0] = []byte("")
	l.levelString[1] = []byte(trace)
	l.levelString[2] = []byte(debug)
	l.levelString[3] = []byte(info)
	l.levelString[4] = []byte(warn)
	l.levelString[5] = []byte(error)
	l.levelString[6] = []byte(fatal)
	return l
}

// SetLogLevel will set the minimum logging Level. If this is set to INFO, anything below
// info, such as TRACE/DEBUG, will be not printed.
func (l *Logger) SetLogLevel(minLogLevel Level) *Logger {
	l.logLevel = minLogLevel
	return l
}

// SetLogTag will take a bit-flag Tag and sets what categories will be printed.
func (l *Logger) SetLogTag(tags Tag) *Logger {
	l.logTag = tags
	return l
}

// SetLogFn will set a filter based on a func `func(Level, Tag)bool`.
// Once this set, it will supercedes SetLogLevel and SetLogTag.
// If `nil` is given as a param, it will be deactivated.
func (l *Logger) SetLogFn(f func(Level, Tag) bool) *Logger {
	l.logFn = f
	return l
}

// UseTag will issue tags to Tag(s) pointers.
// If a logger is created with dots such as `alog.New(out).SetPrefix("nptest ")...`
// This can be used. This maybe useful when there are many tags to be initialized.
// Usage:
//    var TEST1, TEST2, TEST3 alog.Tag
//    l := alog.New(out).UseTag(&TEST1, &TEST2, &TEST3).SetLogTag(TEST1)
func (l *Logger) UseTag(tags ...*Tag) *Logger {
	for _, t := range tags {
		*t = l.NewTag()
	}
	return l
}

// header will add date/time/prefix/Level.
func (l *Logger) header(buf *[]byte, lvl Level) {
	if l.flag&(FdateYYMMDD|FdateYYYYMMDD|FdateMMDD|Ftime|FtimeMs) != 0 {
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
		if l.flag&(Ftime|FtimeMs) != 0 {
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
}

// finalize will add newline to the end of log if missing,
// also write it to writer, and clear the buffer.
func (l *Logger) finalize() (n int, err error) {
	// This can be moved to output().
	if l.flag&Fnewline != 0 {
		// If the log message doesn't end with newline, add a newline.
		if curBufSize := len(l.buf); curBufSize > 1 && l.buf[curBufSize-1] != newline {
			l.buf = append(l.buf, newline)
		}
	} else {
		// Remove all newlines
		for i, v := range l.buf {
			if v == newline {
				l.buf[i] = byte(' ')
			}
		}
		// Append a newline at the end
		l.buf = append(l.buf, newline)
	}

	// If bufUseBuffer is false or current size is bigger than the buffer size,
	// print the buffer and reset it.
	n, err = l.out.Write(l.buf)
	l.buf = l.buf[:0]
	return n, err
}

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

// outputf creates formatted string
func (l *Logger) outputf(lvl Level, format string, a ...interface{}) {
	// Format the header and add it to buffer
	l.header(&l.buf, lvl)
	// Parse formatted string and add it to buffer
	formats(&l.buf, format, a...)
	// Check newline at the end, if missing add it.
	// Then, print log, reset the buffer.
	_, _ = l.finalize()
}

// Output prints a byte array log message.
// Both Level and Tag has to match with what's in the config.
// (However, if Tag is 0, then it will be printed regardless of logTag).
// For Print, even if fatal Level is given, it will not exit.
func (l *Logger) Output(lvl Level, tag Tag, b []byte) {
	// Check if given lvl/logTag are printable
	if l.check(lvl, tag) {
		l.mu.Lock()
		l.header(&l.buf, lvl)
		l.buf = append(l.buf, b...)
		_, _ = l.finalize()
		l.mu.Unlock()
	}
}

// Print prints a single string log message.
// Both Level and Tag has to match with what's in the config.
// (However, if Tag is 0, then it will be printed regardless of logTag).
// For Print, even if fatal Level is given, it will not exit.
func (l *Logger) Print(lvl Level, tag Tag, s string) {
	// Check if given lvl/logTag are printable
	if l.check(lvl, tag) {
		l.mu.Lock()
		l.header(&l.buf, lvl)
		l.buf = append(l.buf, s...)
		_, _ = l.finalize()
		l.mu.Unlock()
	}
}

// Printf prints formatted logs if Level and Tag is met.
// For Printf, even if fatal Level is given, it will not exit.
// If Tag is 0, it will print regardless of Tag being filtered/set.
func (l *Logger) Printf(lvl Level, tag Tag, format string, a ...interface{}) {
	// Both lvl and Tag has to match
	// If logTag is 0, then all logTag.
	if l.check(lvl, tag) {
		l.mu.Lock()
		l.outputf(lvl, format, a...)
		l.mu.Unlock()
	}
}

// Trace will take a single string and print log without Tag
func (l *Logger) Trace(s string) {
	l.Print(Ltrace, noTag, s)
}

// Tracef will formats and print log without Tag
func (l *Logger) Tracef(format string, a ...interface{}) {
	if len(a) == 0 {
		l.Print(Ltrace, noTag, format)
		return
	}
	if l.check(Ltrace, noTag) {
		l.mu.Lock()
		l.outputf(Ltrace, format, a...)
		l.mu.Unlock()
	}
}

// Debug will take a single string and print log without a Tag
func (l *Logger) Debug(s string) {
	l.Print(Ldebug, noTag, s)
}

// Debugf will formats and print log without a Tag
func (l *Logger) Debugf(format string, a ...interface{}) {
	if len(a) == 0 {
		l.Print(Ldebug, noTag, format)
		return
	}
	if l.check(Ldebug, noTag) {
		l.mu.Lock()
		l.outputf(Ldebug, format, a...)
		l.mu.Unlock()
	}
}

// Info will take a single string and print log without a Tag
func (l *Logger) Info(s string) {
	l.Print(Linfo, noTag, s)
}

// Infof will formats and print log without a Tag
func (l *Logger) Infof(format string, a ...interface{}) {
	if len(a) == 0 {
		l.Print(Linfo, noTag, format)
		return
	}
	if l.check(Linfo, noTag) {
		l.mu.Lock()
		l.outputf(Linfo, format, a...)
		l.mu.Unlock()
	}
}

// Warn will take a single string and print log without a Tag
func (l *Logger) Warn(s string) {
	l.Print(Lwarn, noTag, s)
}

// Warnf will formats and print log without a Tag
func (l *Logger) Warnf(format string, a ...interface{}) {
	if len(a) == 0 {
		l.Print(Lwarn, noTag, format)
		return
	}
	if l.check(Lwarn, noTag) {
		l.mu.Lock()
		l.outputf(Lwarn, format, a...)
		l.mu.Unlock()
	}
}

// Error will take a single string and print log without a Tag
func (l *Logger) Error(s string) {
	l.Print(Lerror, noTag, s)
}

// Errorf will formats and print log without a Tag
func (l *Logger) Errorf(format string, a ...interface{}) {
	if len(a) == 0 {
		l.Print(Lerror, noTag, format)
		return
	}
	if l.check(Lerror, noTag) {
		l.mu.Lock()
		l.outputf(Lerror, format, a...)
		l.mu.Unlock()
	}
}

// Fatal will take a single string and print log without a Tag
// and this will terminate process with exit code 1
// updated with Close() as v0.1.6c4, 12/30/2020
func (l *Logger) Fatal(s string) {
	l.Print(Lfatal, noTag, s)
	_ = l.Close()
	os.Exit(1)
}

// Fatalf will formats and print log without a Tag
// and this will terminate process with exit code 1
// updated with Close() as v0.1.6c4, 12/30/2020
func (l *Logger) Fatalf(format string, a ...interface{}) {
	if len(a) == 0 {
		l.Print(Lfatal, noTag, format)
	} else if l.check(Lfatal, noTag) {
		l.mu.Lock()
		l.outputf(Lfatal, format, a...)
		l.mu.Unlock()
	}
	_ = l.Close()
	os.Exit(1)
}

// IfError will check and log error if exist (not nil)
// For instance, when running multiple lines of error check
// This can save error checking.
// added as v0.1.6c3, 12/30/2020
func (l *Logger) IfError(e error) {
	if e != nil {
		l.Print(Lerror, noTag, e.Error())
	}
}

// IfFatal will check and log error if exist (not nil)
// For instance, when running multiple lines of error check
// This can save error checking.
// Unlike IfError, IfFatal will exit the program
// added as v0.1.6c3, 12/30/2020
// updated with Close() as v0.1.6c4, 12/30/2020
func (l *Logger) IfFatal(e error) {
	if e != nil {
		l.Print(Lfatal, noTag, e.Error())
		_ = l.Close()
		os.Exit(1)
	}
}

// Writer returns the output destination for the logger.
func (l *Logger) Writer() io.Writer {
	return l.out
}

// Close will call .Close() method if supported
func (l *Logger) Close() error {
	if c, ok := l.out.(io.Closer); ok && c != nil {
		return c.Close()
	}
	return nil
}

// NewPrint takes Level and Tag and create a print function.
// This is to make such as custom `*Logger.Debug()` that has a Tag
// predefined.
// Note: for outputf, due to memory allocation occurred it is not included.
func (l *Logger) NewPrint(lvl Level, tag Tag, prefix string) func(string) {
	return func(s string) {
		if l.check(lvl, tag) {
			l.mu.Lock()
			l.header(&l.buf, lvl)
			l.buf = append(l.buf, append([]byte(prefix), s...)...)
			_, _ = l.finalize()
			l.mu.Unlock()
		}
	}
}

// NewWriter takes a Level and a Tag and create an Alog writer (SubWriter)
// that is compatible with io.Writer interface. This can be used as a
// logger hook.
func (l *Logger) NewWriter(lvl Level, tag Tag, prefix string) *SubWriter {
	return &SubWriter{
		l:      l,
		lvl:    lvl,
		tag:    tag,
		prefix: []byte(prefix),
	}
}

// NewTag will generate new Tag to be used for user.
// This is nothing but creating a big-flag, but easier for the user
// who aren't familiar with a bit-flag.
func (l *Logger) NewTag() Tag {
	if l.logTagIssued == 0 {
		l.logTagIssued = 1
		return 1
	}
	l.logTagIssued *= 2
	return l.logTagIssued
}

// Tag is a bit-flag used to show only necessary part of process to show
// in the log. For instance, if there's an web service, there can be different
// Tag such as UI, HTTP request, HTTP response, etc. By setting a Tag
// for each log using `Print` or `Printf`, a user can only print certain
// Tag of log messages for better debugging.
type Tag uint64

// SubWriter is a writer with predefined Level and Tag.
type SubWriter struct {
	l      *Logger
	lvl    Level
	tag    Tag
	prefix []byte
}

// Write is to be used as io.Writer interface
func (w *SubWriter) Write(b []byte) (n int, err error) {
	if w.l.check(w.lvl, w.tag) {
		w.l.mu.Lock()
		w.l.header(&w.l.buf, w.lvl)
		w.l.buf = append(w.l.buf, w.prefix...)
		w.l.buf = append(w.l.buf, b...)
		n, err := w.l.finalize()
		w.l.mu.Unlock()
		return n, err
	}
	return 0, nil
}

// devNull is a type for discard
type devNull int

// discard is defined here to get rid of needs to import of ioutil package.
var discard io.Writer = devNull(0)

// Write discards everything
func (devNull) Write([]byte) (int, error) {
	return 0, nil
}

// itoa converts int to []byte
// if minLength == 0, it will print without padding 0
// due to limit on int type, 19 digit max; 18 digit is safe.
func itoa(dst *[]byte, i int, minLength int, suffix byte) {
	var b [22]byte
	var positiveNum = true
	if i < 0 {
		positiveNum = false
		i = -i // change the sign to positive
	}
	bIdx := len(b) - 1
	if suffix != 0 {
		b[bIdx] = suffix
		bIdx--
	}

	for i >= 10 || minLength > 1 {
		minLength--
		q := i / 10
		b[bIdx] = byte('0' + i - q*10)
		bIdx--
		i = q
	}

	b[bIdx] = byte('0' + i)
	if positiveNum == false {
		bIdx--
		b[bIdx] = '-'
	}
	*dst = append(*dst, b[bIdx:]...)
}

// ftoa takes float64 and converts and add to dst byte slice pointer.
// this is used to reduce memory allocation.
func ftoa(dst *[]byte, f float64, decPlace int) {
	if int(f) == 0 && f < 0 {
		*dst = append(*dst, '-')
	}
	itoa(dst, int(f), 0, 0) // add full number first

	if decPlace > 0 {
		// if decPlace == 3, multiplier will be 1000
		// get nth power
		var multiplier = 1
		for i := decPlace; i > 0; i-- {
			multiplier = multiplier * 10
		}
		*dst = append(*dst, '.')
		tmp := int((f - float64(int(f))) * float64(multiplier))
		if f > 0 { // 2nd num shouldn't include decimala
			itoa(dst, tmp, decPlace, 0)
		} else {
			itoa(dst, -tmp, decPlace, 0)
		}
	}
}

// formats method is a replacement for fmt.Sprintf(). This is to save memory allocation.
// This utilizes bufFormat and each run, it will reset it and reuse it.
func formats(dst *[]byte, s string, a ...interface{}) {
	flagKeyword := false

	var aIdx = 0
	var aLen = len(a)

	for _, c := range s {
		if flagKeyword == false {
			if c == '%' {
				flagKeyword = true
			} else {
				*dst = append(*dst, byte(c))
			}
		} else {
			// flagKeyword == true
			if c == '%' {
				*dst = append(*dst, '%')
				flagKeyword = false
				continue
			}
			if aIdx >= aLen {
				flagKeyword = false
				continue
			}
			switch c {
			case 'd':
				if v, ok := a[aIdx].(int); ok {
					itoa(dst, v, 0, 0)
				} else {
					*dst = append(*dst, unsuppType...)
				}
				aIdx++
			case 's':
				if v, ok := a[aIdx].(string); ok {
					*dst = append(*dst, v...)
				} else {
					*dst = append(*dst, unsuppType...)
				}
				aIdx++
			case 'f':
				switch a[aIdx].(type) {
				case float64:
					if v, ok := a[aIdx].(float64); ok {
						ftoa(dst, v, 2)
					} else {
						*dst = append(*dst, unsuppType...)
					}
				case float32:
					if v, ok := a[aIdx].(float32); ok {
						ftoa(dst, float64(v), 2)
					} else {
						*dst = append(*dst, unsuppType...)
					}
				}
				aIdx++
			case 't':
				if v, ok := a[aIdx].(bool); ok {
					if v {
						*dst = append(*dst, "true"...)
					} else {
						*dst = append(*dst, "false"...)
					}
				} else {
					*dst = append(*dst, unsuppType...)
				}
				aIdx++
			}
			flagKeyword = false
		}
	}
}

// DoColor is an example of Do function creation.
// This function returns do-function for alog, and is an example for `*Logger.Do` application.
// Usage: `alog.New(os.Stderr).Do(alog.DoColor())`
func DoColor() func(*Logger) {
	trc := "[TRC]"
	dbg := "[DBG]"
	inf := "[INF]"
	wrn := "[WRN]"
	err := "[ERR]"
	ftl := "[FTL]"

	return func(l *Logger) {
		l.SetLevelPrefix(
			"\u001B[0;35m"+trc+"\u001B[0m ",
			"\u001B[0;36m"+dbg+"\u001B[0m ",
			"\u001B[0;34m"+inf+"\u001B[0m ",
			"\u001B[1;33m"+wrn+"\u001B[0m ",
			"\u001B[1;31m"+err+"\u001B[0m ",
			"\u001B[1;41;30m"+ftl+"\u001B[0m ",
		)
		// IF output is set to os.Stderr OR os.Stdout, it can be done by checking output.
		// if l.Writer() != nil && (l.Writer() == os.Stderr || l.Writer() == os.Stdout) {
		// 	l.SetLevelPrefix(
		// 		"[\u001B[0;35mTRC\u001B[0m] ",
		//      ...
		// 	)
		// } else {
		// 	l.SetLevelPrefix(
		// 		Trace,
		//      ...
		// 	)
		// }
	}
}
