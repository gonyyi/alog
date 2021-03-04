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

const (
	// Fprefix will show prefix when printing log message
	Fprefix     Flag = 1 << iota
	Fsuffix          // Fsuffix will add suffix
	Fdate            // Fdate will show both CCYY and MMDD
	FdateDay         // FdateDay will show 0-6 for JSON or (Sun-Mon)
	Ftime            // Ftime will show HHMMSS
	FtimeMs          // FtimeMs will show time + millisecond --> JSON: HHMMSS000, Text: HHMMSS,000
	FtimeUnix        // FtimeUnix will show unix time
	FtimeUnixMs      // FtimeUnixNano will show unix time
	FUTC             // FUTC will show UTC time formats
	Flevel           // Flevel show level in the log messsage.
	Ftag             // Ftag will show tags

	Fdefault = Ftime | Fdate | Flevel | Ftag
	// fUseTime is precalculated time for internal functions.
	fUseTime = Fdate | FdateDay | Ftime | FtimeMs | FtimeUnix | FtimeUnixMs
)

const (
	// Ltrace shows trace level, thee most detailed debugging level.
	// This will show everything.
	Ltrace Level = iota + 1
	Ldebug       // Ldebug shows debug level or higher
	Linfo        // Linfo shows information level or higher
	Lwarn        // Lwarn is for a normal but a significant condition
	Lerror       // Lerror shows error level or higher
	Lfatal       // Lfatal shows fatal level or higher. This does not exit the process
)

// level is a flag for logging level
type Level uint8

// Name will print level'Vstr full name
func (l Level) Name() string {
	switch l {
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

// NameShort will print level'Vstr abbreviated name
func (l Level) NameShort() string {
	switch l {
	case Ltrace:
		return "TRC"
	case Ldebug:
		return "DBG"
	case Linfo:
		return "INF"
	case Lwarn:
		return "WRN"
	case Lerror:
		return "ERR"
	case Lfatal:
		return "FTL"
	default:
		return ""
	}
}

// DoFn will be used to manipulate multiple functionality at once.
type DoFn func(Logger) Logger

// ControlFn is used to trigger whether log or not in control.
// Once ControlFn is set, level/tag conditions will be ignored.
type ControlFn func(Level, Tag) bool
