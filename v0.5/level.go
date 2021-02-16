package alog

const (
	// Ltrace shows trace Level, thee most detailed debugging Level.
	// This will show everything.
	Ltrace Level = iota + 1
	// Ldebug shows debug Level or higher
	Ldebug
	// Linfo shows information Level or higher
	Linfo
	// Lnotice is when requires special handling.
	Lnotice
	// Lwarn is for a normal but a significant condition
	Lwarn
	// Lerror shows error Level or higher
	Lerror
	// Lfatal shows fatal Level or higher. This does not exit the process
	Lfatal
)

// Level is a flag for logging level
type Level uint8

// String will print level's full name
func (l *Level) String() string {
	switch *l {
	case Ltrace:
		return "trace"
	case Ldebug:
		return "debug"
	case Linfo:
		return "info"
	case Lnotice:
		return "notice"
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

// ShortName will print level's abbreviated name
func (l *Level) ShortName() string {
	switch *l {
	case Ltrace:
		return "TRC"
	case Ldebug:
		return "DBG"
	case Linfo:
		return "INF"
	case Lnotice:
		return "NTC"
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
