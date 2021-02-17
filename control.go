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

// Format a bit-formatFlag formatFlag options that is used for variety of configuration.
type Format uint32

func (f Format) Reset() Format {
	return Format(uint32(0))
}
func (f Format) On(item Format) Format {
	return f | item
}
func (f Format) Off(item Format) Format {
	return f &^ item
}

const (
	// Fprefix will show prefix when printing log message
	Fprefix Format = 1 << iota
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
	fUseTime = Fdate | FdateDay | Ftime | FtimeUnix | FtimeUnixMs
)

type control struct {
	Tags     TagBucket
	hook     HookFn
	ctlFn    func(Level, Tag) bool
	ctlLevel Level
	ctlTag   Tag
}

func (c *control) SetHook(h HookFn) {
	c.hook = h
}

// HookFn is a type for a function designed to run when certain condition meets
type HookFn func(lvl Level, tag Tag, p []byte)

// SetFilter will define what level or tags to show.
// Integer 0 can be used, and when it's used, it will not Filter anything.
func (c *control) CtlTag(lv Level, tags Tag) {
	c.ctlLevel = lv
	c.ctlTag = tags
}

func (c *control) CtlFn(fn func(Level, Tag) bool) {
	// didn't check for nil, because if it's nil, it will simple remove current one.
	c.ctlFn = fn
}

// check will check if Level and Tag given is good to be printed.
func (c *control) Check(lvl Level, tag Tag) bool {
	switch {
	case c.ctlFn != nil: // FilterFn has the highest order if Set.
		return c.ctlFn(lvl, tag)
	case c.ctlLevel < lvl: // if level is higher than set level, print
		return true
	case c.ctlTag&tag != 0: // even if level is not high, if tag matches, print
		return true
	default:
		return false
	}
}

// Tag is a bit-formatFlag used to show only necessary part of process to show
// in the log. For instance, if there's an web service, there can be different
// Tag such as UI, HTTP request, HTTP response, etc. By alConf a Tag
// for each log using `Print` or `Printf`, a user can only print certain
// Tag of log messages for better debugging.
type Tag uint64

func (tag Tag) Has(t Tag) bool {
	if tag&t != 0 {
		return true
	}
	return false
}

type TagBucket struct {
	count int        // count stores number of Tag issued.
	names [64]string // names stores Tag names.
}

// GetTag returns a tag if found
func (t *TagBucket) GetTag(name string) (tag Tag, ok bool) {
	for i := 0; i < t.count; i++ {
		if t.names[i] == name {
			return 1 << i, true
		}
	}
	return 0, false
}

// MustGetTag returns a tag if found. If not, create a new tag.
func (t *TagBucket) MustGetTag(name string) Tag {
	if tag, ok := t.GetTag(name); ok {
		return tag
	}

	// create a new tag if tag not found
	t.names[t.count] = name
	tag := t.count // this is the value to be printed.

	t.count += 1
	return 1 << tag
}

func (t *TagBucket) AppendSelectedTags(dst []byte, delimiter byte, quote bool, tag Tag) []byte {
	if tag == 0 {
		return dst
	}
	cntDst := len(dst)
	for i := 0; i < t.count; i++ {
		if tag&(1<<i) != 0 {
			if quote { // redundant; as speed matter rather than the binary size
				dst = append(dst, '"')
				dst = append(dst, t.names[i]...)
				dst = append(dst, '"', delimiter)
			} else {
				dst = append(dst, t.names[i]...)
				dst = append(dst, delimiter)
			}
		}
	}
	if cntDst < len(dst) {
		return dst[:len(dst)-1] // last delimiter to be omitted
	}
	return dst
}
