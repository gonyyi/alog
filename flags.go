package alog

// format a bit-flag flag options that is used for variety of configuration.
type format uint32

const (
	// Fprefix will show prefix when printing log message
	Fprefix format = 1 << iota
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
	Fall = format(^uint32(0))
	// Fnone for all options off
	Fnone = format(uint32(0))
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

func (l *Level) ShortName() string {
	switch *l {
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
