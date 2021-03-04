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
	UseLevel      Flag = 1 << iota // UseLevel show level in the log messsage.
	UseTag                         // UseTag will show tags
	UseDate                        // UseDate will show both CCYY and MMDD
	UseDay                         // UseDay will show 0-6 for JSON or (Sun-Mon)
	UseTime                        // UseTime will show HHMMSS
	UseTimeMs                      // UseTimeMs will show time + millisecond --> JSON: HHMMSS000, Text: HHMMSS,000
	UseUnixTime                    // UseUnixTime will show unix time
	UseUnixTimeMs                  // UseUnixTimeMs will show unix time with millisecond
	UseUTC                         // UseUTC will show UTC time formats

	UseDefault = UseTime | UseDate | UseLevel | UseTag
	// fUseTime is precalculated time for internal functions. Not that if UseUTC is used by it self,
	// without any below, it won't print any time.
	fUseTime = UseDate | UseDay | UseTime | UseTimeMs | UseUnixTime | UseUnixTimeMs
)

const (
	// TraceLevel shows trace level, thee most detailed debugging level.
	// This will show everything.
	TraceLevel Level = iota + 1
	DebugLevel       // DebugLevel shows debug level or higher
	InfoLevel        // InfoLevel shows information level or higher
	WarnLevel        // WarnLevel is for a normal but a significant condition
	ErrorLevel       // ErrorLevel shows error level or higher
	FatalLevel       // FatalLevel shows fatal level or higher. This does not exit the process
)

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

// ControlFn is used to trigger whether log or not in control.
// Once ControlFn is set, level/tag conditions will be ignored.
type ControlFn func(Level, Tag) bool
