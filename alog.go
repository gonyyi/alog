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
	time time.Time
	flag format

	pool sync.Pool
	fmtr Formatter

	lvtag tagger // lvtag to replace 5 items below

	// filterFn is a customizable function space and supercedes builtin Level and Tag filters if set.
	logFn        func(Level, Tag) bool
	logLevel     Level      // filterLvl stores current logging level
	logTag       Tag        // filterTag stores current logging Tag (bitflag)
	logTagIssued int        // numTagIssued stores number of Tag issued.
	logTagString [64]string // tagNames stores Tag names.

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
		fmtr:     &FormatterText{},
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

// SetNewTags will initialize the wTag.
// Although GetTag returns a tag, SetNewTags will initialize those
// even though those aren't used later on.
func (l *Logger) SetNewTags(names ...string) *Logger {
	l.lvtag.newTags(names...)
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

// SetPrefix can redefine prefix after the logger has been created.
func (l *Logger) SetPrefix(s string) *Logger {
	l.mu.Lock()
	l.prefix = []byte(s)
	l.mu.Unlock()
	return l
}

// SetFormatter will set logger formatter. Without this setting,
// the logger's default will be text formatter (FormatterText).
// If switching between default JSON and TEXT formatter, it should
// be done by UpdateFormat() or SetFormat(). This SetFormatter is
// mainly for a customized formatter.
func (l *Logger) SetFormatter(fmtr Formatter) *Logger {
	if fmtr != nil {
		l.fmtr = fmtr
		return l
	}
	l.fmtr = FormatterText{}
	return l
}

// SetFormat will reconfigure the logger after it has been created.
// This will first copy flag into *alog.flag, and sets few that
// need additional parsing.
func (l *Logger) SetFormat(flag format) *Logger {
	l.mu.Lock()

	// Previous had a JSON flag on, but anymore, then use text formatter
	// If previous didn't have JSON flag, but now do, then use json formatter
	if l.flag&Fjson != 0 && flag&Fjson == 0 {
		l.SetFormatter(&FormatterText{})
	} else if l.flag&Fjson == 0 && flag&Fjson != 0 {
		l.SetFormatter(&FormatterJSON{})
	}

	l.flag = flag

	l.mu.Unlock()
	return l
}

// UpdateFormat allow to adjust item(s) on/off without
// impacting what is already set. This is helpful when override
// certain flag(s).
func (l *Logger) UpdateFormat(item format, on bool) {
	l.mu.Lock()

	// Previous had a JSON flag on, but anymore, then use text formatter
	// If previous didn't have JSON flag, but now do, then use json formatter
	if item&Fjson != 0 {
		if on && l.flag&Fjson == 0 {
			l.SetFormatter(&FormatterJSON{})
		} else if on == false && l.flag&Fjson != 0 {
			l.SetFormatter(&FormatterText{})
		}
	}

	if on {
		l.flag = l.flag | item
	} else {
		l.flag = l.flag &^ item
	}

	l.mu.Unlock()
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

// IfErr will check and log error if exist (not nil)
// For instance, when running multiple lines of error check
// This can save error checking.
// added as v0.1.6c3, 12/30/2020
func (l *Logger) IfErr(e error, lvl Level, tag Tag, msg string) {
	if e != nil {
		l.Log(lvl, tag, msg, "err", e.Error())
	}
}

// NewWriter takes a Level and a Tag and create an Alog writer (SubWriter)
// that is compatible with io.Writer interface. This can be used as a
// logger hook.
func (l *Logger) NewWriter(lvl Level, tag Tag) *SubWriter {
	return &SubWriter{
		l:      l,
		wLevel: lvl,
		wTag:   tag,
	}
}
