package alog

var dFmtChars [256]bool
var dFmt formatd

func init() {
	for i := 0; i <= 0x7e; i++ {
		dFmtChars[i] = i >= 0x20 && i != '\\' && i != '"' // all printable will be true
	}
}

// Flag a bit-formatFlag formatFlag options that is used for variety of configuration.
type Flag uint32

// level is a flag for logging level
type Level uint8

// Name will print level'Vstr full name
func (l Level) Name() string {
	switch l {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return ""
	}
}

// NameShort will print level'Vstr abbreviated name
func (l Level) NameShort() string {
	switch l {
	case TraceLevel:
		return "TRC"
	case DebugLevel:
		return "DBG"
	case InfoLevel:
		return "INF"
	case WarnLevel:
		return "WRN"
	case ErrorLevel:
		return "ERR"
	case FatalLevel:
		return "FTL"
	default:
		return ""
	}
}

// DoFn will be used to manipulate multiple functionality at once.
type DoFn func(Logger) Logger

// EntryFn will
type EntryFn func(*Entry) *Entry

// ControlFn is used to trigger whether log or not in control.
// Once ControlFn is set, level/tag conditions will be ignored.
type ControlFn func(Level, Tag) bool
