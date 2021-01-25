package alog

import (
	"io"
	"sync"
	"time"
)

// flags a bit-flag flag options that is used for variety of configuration.
type flags uint32

const (
	// Fprefix will show prefix when printing log message
	Fprefix flags = 1 << iota
	// Fdate will show both CCYY and MMDD
	Fdate
	// FdateDay will show 0-6 for JSON or (Sun-Mon)
	FdateDay
	// Ftime will show HHMMSS.000; for json, it will be HHMMSS000
	Ftime
	// FtimeUnix will show unix time
	FtimeUnix
	// FtimeUnixNano will show unix time
	FtimeUnixMs
	// FtimeUTC will show UTC time formats
	FtimeUTC
	// Flevel show Level in the log messsage.
	Flevel
	// Ftag will show tags
	Ftag
	// Fjson will print to a JSON
	Fjson

	// Fdefault will show month/day with time, and Level of logging.
	Fdefault = Fdate | Ftime | Flevel | Ftag
	// Fall for all options on
	Fall = flags(^uint32(0))
	// Fnone for all options off
	Fnone = flags(uint32(0))
)

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
	// Lfatal shows fatal Level or higher. This does not exit the process
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

	pool sync.Pool
	fmtr AlogFmtr

	// logFn is a customizable function space and supercedes builtin Level and Tag filters if set.
	logFn        func(Level, Tag) bool
	logLevel     Level      // logLevel stores current logging level
	logTag       Tag        // logTag stores current logging Tag (bitflag)
	logTagIssued int        // logTagIssued stores number of Tag issued.
	logTagString [64]string // logTagString stores Tag names.

	out io.Writer
	mu  sync.Mutex

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
	// eg. If a user called SetFormat with other config flag except the Level, then the log
	//     Level will not be changed. Therefore default Level should be defined here.
	l := &Logger{
		pool: sync.Pool{
			New: func() interface{} {
				return newItem(512)
			},
		},
		fmtr:     TextFmtr{},
		out:      output,
		prefix:   []byte(""), // prefix will be saved as a byte slice to prevent need to be converted later.
		logLevel: Linfo,      // default logging Level to INFO
		flag:     Fdefault,   // default flag is given
		// buf:      make([]byte, 1024),
	}
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

// SetFormatter will set logger formatter
func (l *Logger) SetFormatter(fmtr AlogFmtr) *Logger {
	if fmtr != nil {
		l.fmtr = fmtr
		return l
	}
	l.fmtr = TextFmtr{}
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

// SetFormat will reconfigure the logger after it has been created.
// This will first copy flag into *alog.flag, and sets few that
// need additional parsing.
func (l *Logger) SetFormat(flag flags) *Logger {
	l.mu.Lock()
	l.flag = flag
	l.mu.Unlock()
	return l
}

// SetFormatItem allow to adjust a single item on/off without
// impacting what is already set.
func (l *Logger) SetFormatItem(item flags, on bool) *Logger {
	l.mu.Lock()
	if on {
		l.flag = l.flag | item
	} else {
		l.flag = l.flag &^ item
	}
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

// SetFilter will define what level or tags to show.
// Integer 0 can be used, and when it's used, it will not filter anything.
func (l *Logger) SetFilter(lv Level, tags Tag) *Logger {
	l.logFn = nil
	l.logLevel = lv
	l.logTag = tags

	return l
}

// SetFilterFn can control more precisely by taking a FilterFn.
func (l *Logger) SetFilterFn(fn FilterFn) *Logger {
	l.logFn = fn
	return l
}

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
		if do != nil {
			do()
		} else {
			l.Error(0, "", "err", e)
		}
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
	case l.logTag != 0 && l.logTag&tag == 0: // if logTag is set but Tag is not matching, then do not print
		return false
	default:
		return true
	}
}
