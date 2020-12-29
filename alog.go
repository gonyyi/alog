// (c) 2020 Gon Y Yi. <https://gonyyi.com/copyright.txt>
// Version 0.1.2, 12/29/2020

package alog

import (
	"io"
	"os"
	"sync"
	"time"
)

// newline constant provides byte of newline, so it can be usd right away.
const (
	newline             = byte('\n')
	noCategory Category = 0
)

// unsuppType is a slice of byte and will be used when unknown formats string is being used in
// any formatted prints such as `printf`, `infof`, `debugf`, etc. This is pre-converted to
// a byte slice and reused to save process time.
var unsuppType = []byte("?{unexp}")

// formatflag a bit-flag flag options that is used for variety of configuration.
type formatflag uint32

const (
	Ftime         formatflag = 1 << iota // Ftime will show HH:MM:SS formats such as 05:02:03
	FtimeMs                              // FtimeMs will show millisecond in its time such as 10:12:13.1234
	FtimeUTC                             // FtimeUTC will show UTC time formats
	FdateMMDD                            // FdateMMDD will show 01/02 date formats.
	FdateYYMMDD                          // FdateYYMMDD will show 06/01/02 date formats.
	FdateYYYYMMDD                        // FdateYYYYMMDD will show 2006/01/02 date formats.
	Fprefix                              // Fprefix will show prefix when printing log message
	Flevel                               // Flevel show level in the log messsage.

	Fdefault = FdateMMDD | Ftime | Flevel // Fdefault will show month/day with time, and level of logging.
)

// LevelPrefix is a bit-flag used for different level of log activity:
// - Ltrace: detailed debugging level
// - Ldebug: general debugging level
// - Linfo: information level
// - Lwarn: warning
// - Lerror: error
// - Lfatal: fatal, the process will can be terminated
type level uint8

const (
	Ltrace level = iota + 1 // Ltrace: detailed debugging level
	Ldebug                  // Ldebug: general debugging level
	Linfo                   // Linfo: information level
	Lwarn                   // Lwarn: warning
	Lerror                  // Lerror: error
	Lfatal                  // Lfatal: fatal, the process will can be terminated
)

type Logger struct {
	flag  formatflag
	level level
	cat   Category

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

	// levelString is an array of byte slice that stores what prefix per each log level
	// will be used. Eg. "[DEBUG]", etc.
	levelString [7][]byte
}

// New function creates new logger. This takes three arguments:
// 1. output: io.Writer
// 2. prefix: any prefix messages to be printed,
// 3. flag: for detail configuration.
func New(output io.Writer, prefix string, flag formatflag) *Logger {
	// If output is given as `nil`, it will use io.discard as a default.
	if output == nil {
		output = discard
	}
	// Creating new logger object and returns pointer to logger.
	// Default value will be set here. If a user uses *alog.SetFlag to provoke
	// unless specifically set certain value, it will not be overwritten.
	// eg. If a user called SetFlag with other config flag except the level, then the log
	//     level will not be changed. Therefore default level should be defined here.
	l := &Logger{
		out:    output,
		prefix: []byte(prefix), // prefix will be saved as a byte slice to prevent need to be converted later.
		level:  Linfo,          // default logging level to INFO
	}

	// Parse and set configuration flags.
	l.SetFlag(flag)

	// Default prefixes for each level. This can be changed by a user using *alog.SetLevelPrefix()
	l.SetLevelPrefix("[TRC] ", "[DBG] ", "[INF] ", "[WRN] ", "[ERR] ", "[FTL] ")

	return l
}

// SetOutput can redefined the output after logger has been created.
// If output is nil, the logger will set it to ioutil.Discard instead.
func (l *Logger) SetOutput(output io.Writer) {
	l.mu.Lock()
	if output == nil {
		l.out = discard
	} else {
		l.out = output
	}
	l.mu.Unlock()
}

// SetPrefix can redefine prefix after the logger has been created.
func (l *Logger) SetPrefix(s string) {
	l.mu.Lock()
	l.prefix = []byte(s)
	l.mu.Unlock()
}

// SetFlag will reconfigure the logger after it has been created.
// This will first copy flag into *alog.flag, and sets few that
// need additional parsing.
func (l *Logger) SetFlag(flag formatflag) {
	l.mu.Lock()
	l.flag = flag
	l.mu.Unlock()
}

// SetLevelPrefix will set different prefixes for each levelled log messages.
// eg. "[DEBUG]".
func (l *Logger) SetLevelPrefix(trace, debug, info, warn, error, fatal string) {
	l.levelString[0] = []byte("")
	l.levelString[1] = []byte(trace)
	l.levelString[2] = []byte(debug)
	l.levelString[3] = []byte(info)
	l.levelString[4] = []byte(warn)
	l.levelString[5] = []byte(error)
	l.levelString[6] = []byte(fatal)
}

// SetLevel will set the minimum logging level. If this is set to INFO, anything below
// info, such as TRACE/DEBUG, will be not printed.
func (l *Logger) SetLevel(lvl level) {
	l.level = lvl
}

// SetCategory will take a bit-flag Category and sets what categories will be printed.
func (l *Logger) SetCategory(category Category) {
	l.cat = category
}

// header will add date/time/prefix/level.
func (l *Logger) header(buf *[]byte, t time.Time, lvl level) {
	if l.flag&(FdateYYMMDD|FdateYYYYMMDD|FdateMMDD|Ftime|FtimeMs) != 0 {
		if l.flag&FtimeUTC != 0 {
			t = t.UTC()
		}
		if l.flag&(FdateYYMMDD|FdateYYYYMMDD|FdateMMDD) != 0 {
			year, month, day := t.Date()
			// if both YYMMDD and YYYYMMDD is given, YYYYMMDD will be used
			if l.flag&FdateYYYYMMDD != 0 {
				itoa(buf, year, 4)
				*buf = append(*buf, '/')
			} else if l.flag&FdateYYMMDD != 0 {
				itoa(buf, year%100, 2)
				*buf = append(*buf, '/')
			}
			// MMDD will be always added ass it's a common denominator of
			// FdateYYMMDD|FdateYYYYMMDD|FdateMMDD
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(Ftime|FtimeMs) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&FtimeMs != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
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
	// If the log message doesn't end with newline, add a newline.
	if curBufSize := len(l.buf); curBufSize > 1 && l.buf[curBufSize-1] != newline {
		l.buf = append(l.buf, newline)
	}

	// If bufUseBuffer is false or current size is bigger than the buffer size,
	// print the buffer and reset it.
	n, err = l.out.Write(l.buf)
	l.buf = l.buf[:0]
	return n, err
}

// Check method will check if level and category given is
// good to be printed.
// Eg. if setting is level INFO, category USER, then
//     any log level below INFO shouldn't be printed.
//     Also, any category other than USER shouldn't be printed either.
func (l *Logger) check(lvl level, cat Category) bool {
	switch {
	case l.level > lvl: // if lvl is below lvl limit, the do not print
		return false
	case l.cat != noCategory && l.cat&cat == noCategory: // if category is set but category is not matching, then do not print
		return false
	default:
		return true
	}
}

// printf creates formatted string
func (l *Logger) printf(lvl level, cat Category, format string, a ...interface{}) {
	t := time.Now()
	// Format the header and add it to buffer
	l.header(&l.buf, t, lvl)
	// Parse formatted string and add it to buffer
	formats(&l.buf, format, a...)
	// Check newline at the end, if missing add it.
	// Then, print log, reset the buffer.
	l.finalize()
}

// Print prints a single string log message.
// Both level and Category has to match with what's in the config.
// (However, if cat is 0, then it will be printed regardless of category).
// For Print, even if fatal level is given, it will not exit.
func (l *Logger) Print(lvl level, cat Category, s string) {
	// Check if given lvl/category are printable
	if l.check(lvl, cat) {
		t := time.Now()
		l.mu.Lock()
		l.header(&l.buf, t, lvl)
		l.buf = append(l.buf, s...)
		l.finalize()
		l.mu.Unlock()
	}
}

// Printf prints formatted logs if level and category is met.
// For Printf, even if fatal level is given, it will not exit.
// If category is 0, it will print regardless of category being filtered/set.
func (l *Logger) Printf(lvl level, cat Category, format string, a ...interface{}) {
	// Both lvl and Category has to match
	// If cat is 0, then all cat.
	if l.check(lvl, cat) {
		l.mu.Lock()
		l.printf(lvl, cat, format, a...)
		l.mu.Unlock()
	}
}

// Trace will take a single string and print log without category
func (l *Logger) Trace(s string) {
	l.Print(Ltrace, noCategory, s)
}

// Tracef will formats and print log without category
func (l *Logger) Tracef(format string, a ...interface{}) {
	if l.check(Ltrace, noCategory) {
		l.mu.Lock()
		l.printf(Ltrace, noCategory, format, a...)
		l.mu.Unlock()
	}
}

// Debug will take a single string and print log without category
func (l *Logger) Debug(s string) {
	l.Print(Ldebug, noCategory, s)
}

// Debugf will formats and print log without category
func (l *Logger) Debugf(format string, a ...interface{}) {
	if l.check(Ldebug, noCategory) {
		l.mu.Lock()
		l.printf(Ldebug, noCategory, format, a...)
		l.mu.Unlock()
	}
}

// Info will take a single string and print log without category
func (l *Logger) Info(s string) {
	l.Print(Linfo, noCategory, s)
}

// Infof will formats and print log without category
func (l *Logger) Infof(format string, a ...interface{}) {
	if l.check(Linfo, noCategory) {
		l.mu.Lock()
		l.printf(Linfo, noCategory, format, a...)
		l.mu.Unlock()
	}
}

// Warn will take a single string and print log without category
func (l *Logger) Warn(s string) {
	l.Print(Lwarn, noCategory, s)
}

// Warnf will formats and print log without category
func (l *Logger) Warnf(format string, a ...interface{}) {
	if l.check(Lwarn, noCategory) {
		l.mu.Lock()
		l.printf(Lwarn, noCategory, format, a...)
		l.mu.Unlock()
	}
}

// Error will take a single string and print log without category
func (l *Logger) Error(s string) {
	l.Print(Lerror, noCategory, s)
}

// Errorf will formats and print log without category
func (l *Logger) Errorf(format string, a ...interface{}) {
	if l.check(Lerror, noCategory) {
		l.mu.Lock()
		l.printf(Lerror, 0, format, a...)
		l.mu.Unlock()
	}
}

// Fatal will take a single string and print log without category
// and this will terminate process with exit code 1
func (l *Logger) Fatal(s string) {
	l.Print(Lfatal, noCategory, s)
	os.Exit(1)
}

// Fatalf will formats and print log without category
// and this will terminate process with exit code 1
func (l *Logger) Fatalf(format string, a ...interface{}) {
	if l.check(Lfatal, noCategory) {
		l.mu.Lock()
		l.printf(Lfatal, noCategory, format, a...)
		l.mu.Unlock()
	}
	os.Exit(1)
}

// Writer returns the output destination for the logger.
func (l *Logger) Writer() io.Writer {
	return l.out
}

// NewPrint takes level and category and create a print function.
// This is to make such as custom `*Logger.Debug()` that has category
// predefined. Added for v0.1.1.
// For printf, due to memory allocation occurred it is not included.
func (l *Logger) NewPrint(lvl level, cat Category) func(string) {
	return func(s string) {
		l.Print(lvl, cat, s)
	}
}

// NewWriter takes level and category and create an Alog writer (alogw)
// that is compatible with io.Writer interface. This can be used as a
// logger hook.
func (l *Logger) NewWriter(lvl level, cat Category) *alogw {
	return &alogw{
		l:   l,
		lvl: lvl,
		cat: cat,
	}
}

// Category is a bit-flag used to show only necessary part of process to show
// in the log. For instance, if there's an web service, there can be different
// category such as UI, HTTP request, HTTP response, etc. By setting category
// for each log using `Print` or `Printf`, a user can only print certain
// category of log messages for better debugging.
type Category uint64

// NewCategory will generate new categories to be used for user.
// This is nothing but creating a big-flag, but easier for the user
// who aren't familiar with a bit-flag.
func NewCategory() Category {
	return 0
}

// Add new category
func (c *Category) Add() Category {
	if *c == 0 {
		*c = 1
		return 1
	}
	*c *= 2
	return *c
}

// alogw is a writer with predefined level and category.
type alogw struct {
	l   *Logger
	lvl level
	cat Category
}

// Write is to be used as io.Writer interface
func (w *alogw) Write(b []byte) (n int, err error) {
	if w.l.check(w.lvl, w.cat) {
		t := time.Now()
		w.l.mu.Lock()
		w.l.header(&w.l.buf, t, w.lvl)
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

func (devNull) Write(p []byte) (int, error) {
	return 0, nil
}

// itoa converts int to []byte
// if minLength == 0, it will print without padding 0
// due to limit on int type, 19 digit max; 18 digit is safe.
func itoa(dst *[]byte, i int, minLength int) {
	var b [20]byte
	var positiveNum bool = true
	if i < 0 {
		positiveNum = false
		i = -i // change the sign to positive
	}
	bIdx := len(b) - 1

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
	itoa(dst, int(f), 0) // add full number first

	if decPlace > 0 {
		// if decPlace == 3, multiplier will be 1000
		// get nth power
		var multiplier int = 1
		for i := decPlace; i > 0; i-- {
			multiplier = multiplier * 10
		}
		*dst = append(*dst, '.')
		tmp := int((f - float64(int(f))) * float64(multiplier))
		if f > 0 { // 2nd num shouldn't include decimala
			itoa(dst, tmp, decPlace)
		} else {
			itoa(dst, -tmp, decPlace)
		}
	}
}

// formats method is a replacement for fmt.Sprintf(). This is to save memory allocation.
// This utilizes bufFormat and each run, it will reset it and reuse it.
func formats(dst *[]byte, s string, a ...interface{}) {
	flagKeyword := false

	var aIdx int = 0
	var aLen = len(a)

	// Reset bufFormat
	// *dst = []byte{}

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
					itoa(dst, v, 0)
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
						*dst = append(*dst, []byte("true")...)
					} else {
						*dst = append(*dst, []byte("false")...)
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
