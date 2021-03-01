package alog

// Format a bit-formatFlag formatFlag options that is used for variety of configuration.
type Format uint32

const (
	// Fprefix will show prefix when printing log message
	Fprefix     Format = 1 << iota
	Fsuffix            // Fsuffix will add suffix
	Fdate              // Fdate will show both CCYY and MMDD
	FdateDay           // FdateDay will show 0-6 for JSON or (Sun-Mon)
	Ftime              // Ftime will show HHMMSS
	FtimeMs            // FtimeMs will show time + millisecond --> JSON: HHMMSS000, Text: HHMMSS,000
	FtimeUnix          // FtimeUnix will show unix time
	FtimeUnixMs        // FtimeUnixNano will show unix time
	FtimeUTC           // FtimeUTC will show UTC time formats
	Flevel             // Flevel show Level in the log messsage.
	Ftag               // Ftag will show tags

	Fdefault = Ftime | Flevel | Ftag
	// fUseTime is precalculated time for internal functions.
	fUseTime = Fdate | FdateDay | Ftime | FtimeMs | FtimeUnix | FtimeUnixMs
)

const (
	// Ltrace shows trace Level, thee most detailed debugging Level.
	// This will show everything.
	Ltrace Level = iota + 1
	Ldebug       // Ldebug shows debug Level or higher
	Linfo        // Linfo shows information Level or higher
	Lwarn        // Lwarn is for a normal but a significant condition
	Lerror       // Lerror shows error Level or higher
	Lfatal       // Lfatal shows fatal Level or higher. This does not exit the process
)

// Level is a flag for logging level
type Level uint8

// Name will print level's full name
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

// NameShort will print level's abbreviated name
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
