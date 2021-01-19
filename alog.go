// (c) 2020 Gon Y Yi. <https://gonyyi.com>

// version 0.4 candidate

package alog

import (
	"io"
	"sync"
	"time"
)

const (
	// newline constant provides byte of newline, so it can be usd right away.
	newline     = byte('\n')
	quotation   = byte('"')
	unsuppTypes = "{??}"
	// noTag will reset tag if already set
	noTag Tag = 0
)

// unsuppType is a slice of byte and will be used when unknown formats string is being used in
// any formatted prints such as `outputf`, `infof`, `debugf`, etc. This is pre-converted to
// a byte slice and reused to save process time.
var unsuppType = []byte("{??}")
var newlineRepl = []byte(`\n`)
var quotationRepl = []byte(`\"`)

// flags a bit-flag flag options that is used for variety of configuration.
type flags uint32

const (
	// Fall for all options on
	Fall = flags(^uint32(0))
	// Fnone for all options off
	Fnone = flags(uint32(0))

	// Fprefix will show prefix when printing log message
	Fprefix flags = 1 << iota
	// Fyear will show 4 digit year such as 2006
	Fyear
	// Fdate will show 01/02 date formats.
	Fdate
	// Ftime will show HH:MM:SS formats such as 05:02:03
	Ftime
	// FtimeMs will show millisecond in its time such as 10:12:13.1234
	FtimeMs
	// FtimeUnix will show unix time
	FtimeUnix
	// FtimeUnixNano will show unix time
	FtimeUnixNano
	// FtimeUTC will show UTC time formats
	FtimeUTC
	// Flevel show Level in the log messsage.
	Flevel
	// Ftag will show tags
	Ftag
	// Fjson will print to a JSON
	Fjson
	// Fdefault will show month/day with time, and Level of logging.
	Fdefault = Fyear | Fdate | Ftime | Flevel | Ftag
)

// LevelPrefix is a bit-flag used for different Level of log activity:
// - Ltrace: detailed debugging Level
// - Ldebug: general debugging Level
// - Linfo: information Level
// - Lwarn: warning
// - Lerror: error
// - Lfatal: fatal, the process will can be terminated
type Level uint8

func (l *Level) String() string {
	switch *l {
	case Ltrace:
		return "trace"
	case Ldebug:
		return "debug"
	case Linfo:
		return "info"
	case Lwarn:
		return "warn"
	case Lerror:
		return "error"
	case Lfatal:
		return "fatal"
	default:
		return ""
	}
}

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
	logTagIssued int
	logTagString [64]string

	out io.Writer
	mu  sync.Mutex

	// There are two buffers used. Both `buf` and `bufFormat` are being used regardless
	// of bufUseBuffer setting, however, if bufUseBuffer is false, this buffer will be
	// flushed for each log.
	buf    []byte // buf is a main buffer; reset per each log entry
	bufFmt []byte // bufFmt is a buffer for formatting
	// sbufc int

	// bufFormat is a buffer strictly used only for formatting - such as printing
	// date, time, prefix etc; and will be copied to `buf` (main buffer)
	// bufFormat []byte // for formatting only

	// Prefix will be stored as a byte slice.
	prefix []byte

	// levelString is an array of byte slice that stores what prefix per each log logLevel
	// will be used. Eg. "[DEBUG]", etc.
	levelString        [7][]byte
	levelStringForJson [7][]byte
}

// New function creates new logger. This takes an output writer for its argument (v0.2.0 change)
// All methods with suffix "Set" returns `*Logger`, therefore can be used together with `*Logger.New`.
// This is a v0.2.0 change that broke the backward compatibility, however, most of the time, people
// don't set logger prefix, also uses basic default setting. Therefore it's bit cumbersome to require
// two (prefix, flag), often, unused parameters. Duct-taping can be done following way:
// 1. alog.New(nil).SetOutput(os.Stderr) // initially set discard for output but overridden to os.Stderr
// 2. alog.New(os.Stderr).SetPrefix("TestLog: ").SetLogLevel(alog.Linfo) // set prefix and Level
// 3. alog.New(os.Stderr).SetPrefix("TestLog: ").SetFormat(alog.Fdefault|alog.FtimeUTC)
func New(output io.Writer) *Logger {
	// If output is given as `nil`, it will use io.discard as a default.
	if output == nil {
		output = discard
	}

	// Creating new logger object and returns pointer to logger.
	// Default value will be set here. If a user uses *alog.SetFormat to provoke
	// unless specifically set certain value, it will not be overwritten.
	// eg. If a user called SetFormat with other config flag except the Level, then the log
	//     Level will not be changed. Therefore default Level should be defined here.
	l := &Logger{
		out:      output,
		prefix:   []byte(""), // prefix will be saved as a byte slice to prevent need to be converted later.
		logLevel: Linfo,      // default logging Level to INFO
		flag:     Fdefault,   // default flag is given
		buf:      make([]byte, 1024),
	}

	// Default prefixes for each Level. This can be changed by a user using *alog.setLevelPrefix()
	l.levelString[0] = []byte("")
	l.levelString[1] = []byte("[TRC] ")
	l.levelString[2] = []byte("[DBG] ")
	l.levelString[3] = []byte("[INF] ")
	l.levelString[4] = []byte("[WRN] ")
	l.levelString[5] = []byte("[ERR] ")
	l.levelString[6] = []byte("[FTL] ")

	// For JSON output, this is hardcoded
	l.levelStringForJson[0] = []byte("")
	l.levelStringForJson[1] = []byte("trace")
	l.levelStringForJson[2] = []byte("debug")
	l.levelStringForJson[3] = []byte("info")
	l.levelStringForJson[4] = []byte("warn")
	l.levelStringForJson[5] = []byte("error")
	l.levelStringForJson[6] = []byte("fatal")
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

// Close will call .Close() method if supported
func (l *Logger) Close() error {
	if c, ok := l.out.(io.Closer); ok && c != nil {
		return c.Close()
	}
	return nil
}

// MustGetTag will ignore "ok" from GetTag.
// If not found, it will return 0 (none) for the tag.
func (l *Logger) MustGetTag(name string) Tag {
	tag, _ := l.GetTag(name)
	return tag
}

// GetTag takes a tag name and returns a tag if found.
func (l *Logger) GetTag(name string) (tag Tag, ok bool) {
	for i := 0; i < l.logTagIssued; i++ {
		if l.logTagString[i] == name {
			return 1 << i, true
		}
	}
	return 0, false
}

// GetWriter returns the output destination for the logger.
func (l *Logger) GetWriter() io.Writer {
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

// SetPrefix can redefine prefix after the logger has been created.
func (l *Logger) SetPrefix(s string) *Logger {
	l.mu.Lock()
	l.prefix = []byte(s)
	l.mu.Unlock()
	return l
}

// SetFormat will reconfigure the logger after it has been created.
// This will first copy flag into *alog.flag, and sets few that
// need additional parsing.
func (l *Logger) SetFormat(flag flags) *Logger {
	l.mu.Lock()
	l.flag = flag
	l.mu.Unlock()
	return l
}

func (l *Logger) SetNewTags(names ...string) *Logger {
	for _, name := range names {
		if _, ok := l.GetTag(name); !ok {
			l.logTagString[l.logTagIssued] = name
			l.logTagIssued += 1
		}
	}
	return l
}

func (l *Logger) SetFilter(fn FilterFn, lv Level, tags Tag) *Logger {
	if fn != nil {
		l.logFn = fn
	} else {
		l.logLevel = lv
		l.logTag = tags
	}
	return l
}

//
// // Output prints a byte array log message.
// // Both Level and Tag has to match with what's in the config.
// // (However, if Tag is 0, then it will be printed regardless of logTag).
// // For Print, even if fatal Level is given, it will not exit.
// func (l *Logger) Output(lvl Level, tag Tag, b []byte) {
// 	// Check if given lvl/logTag are printable
// 	if l.check(lvl, tag) {
// 		l.mu.Lock()
// 		l.header(&l.buf, lvl, tag)
// 		if l.flag&Fjson != 0 {
// 			lastUpdate := 0
// 			escapeKey := false
// 			for i := 0; i < len(b); i++ {
// 				switch b[i] {
// 				case '\\':
// 					if escapeKey == true {
// 						l.buf = append(l.buf, `\`...)
// 					} else {
// 						escapeKey = true
// 					}
// 				case '\n':
// 					l.buf = append(l.buf, b[lastUpdate:i]...)
// 					l.buf = append(l.buf, `\n`...)
// 					lastUpdate = i + 1
// 					escapeKey = false
// 				case '"':
// 					if escapeKey == false {
// 						l.buf = append(l.buf, b[lastUpdate:i]...)
// 						l.buf = append(l.buf, `\"`...)
// 						lastUpdate = i + 1
// 					}
// 				default:
// 					escapeKey = false
// 				}
// 			}
// 			l.buf = append(l.buf, b[lastUpdate:]...)
// 		} else {
// 			l.buf = append(l.buf, b...)
// 		}
// 		l.finalize()
// 		l.mu.Unlock()
// 	}
// }

// LogIferr will check and log error if exist (not nil)
// For instance, when running multiple lines of error check
// This can save error checking.
// added as v0.1.6c3, 12/30/2020
func (l *Logger) LogIferr(e error, lvl Level, tag Tag, msg string) {
	if e != nil {
		l.Log(lvl, tag, msg, "err", e.Error())
	}
}

// Iferr will run function "do" if error is not nil.
func (l *Logger) Iferr(e error, do func()) {
	if e != nil {
		do()
	}
}

// NewWriter takes a Level and a Tag and create an Alog writer (SubWriter)
// that is compatible with io.Writer interface. This can be used as a
// logger hook.
func (l *Logger) NewWriter(lvl Level, tag Tag) *SubWriter {
	return &SubWriter{
		l:   l,
		lvl: lvl,
		tag: tag,
	}
}

// Tag is a bit-flag used to show only necessary part of process to show
// in the log. For instance, if there's an web service, there can be different
// Tag such as UI, HTTP request, HTTP response, etc. By setting a Tag
// for each log using `Print` or `Printf`, a user can only print certain
// Tag of log messages for better debugging.
type Tag uint64

// SubWriter is a writer with predefined Level and Tag.
type SubWriter struct {
	l   *Logger
	lvl Level
	tag Tag
}

// Write is to be used as io.Writer interface
func (w *SubWriter) Write(b []byte) (n int, err error) {
	if w.l.check(w.lvl, w.tag) {
		w.l.mu.Lock()
		w.l.header(&w.l.buf, w.lvl, w.tag)
		w.l.buf = append(w.l.buf, b...) // todo: check if this works with JSON
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

// FilterFn is a function type to be used with SetFilter.
type FilterFn func(Level, Tag) bool

// itoa converts int to []byte
// if minLength == 0, it will print without padding 0
// due to limit on int type, 19 digit max; 18 digit is safe.
// Keeping this because of minLength and suffix...
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
