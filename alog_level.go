package alog

type Leveller interface {
	String() string
	name_terminal() string
}

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

func (l *Level) name_terminal() string {
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
