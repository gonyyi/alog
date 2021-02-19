package alog

const (
	// Ltrace shows trace Level, thee most detailed debugging Level.
	// This will show everything.
	Ltrace  Level = iota + 1
	Ldebug        // Ldebug shows debug Level or higher
	Linfo         // Linfo shows information Level or higher
	Lnotice       // Lnotice is when requires special handling.
	Lwarn         // Lwarn is for a normal but a significant condition
	Lerror        // Lerror shows error Level or higher
	Lfatal        // Lfatal shows fatal Level or higher. This does not exit the process
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

type control struct {
	Tags     TagBucket
	hook     HookFn
	ctlFn    func(Level, Tag) bool
	ctlLevel Level
	ctlTag   Tag
}

// SetHook will add a HookFn to control.
// This can be used for when additional action is required by either log level OR tag.
// Since HookFn also takes buf body, it can record the certain type totally
// independently as a user customize it. OR can create a condition based on the buf body.
// Note that SetHook is called right BEFORE the Final() is called, which means, IF any log
// message didn't pass the control, it won't reach the hook function either. Alternate way
// of adding hook will be writing a custom AlWriter.
func (c *control) SetHook(h HookFn) {
	c.hook = h
}

// ControlFn enables for a user to set a precise control on what to log.
type ControlFn func(Level, Tag) bool

// HookFn is a type for a function designed to run when certain condition meets
type HookFn func(lvl Level, tag Tag, p []byte)

// CtlTag will define what level or tags to show.
// Integer 0 can be used, and when it's used, it will not Filter anything.
// Note that given level and tags will be used to check the logging eligibility,
// and is when either given level OR tags matches.
func (c *control) CtlTag(lv Level, tags Tag) {
	c.ctlLevel = lv
	c.ctlTag = tags
}

// CtlFn will set the control function.
func (c *control) CtlFn(fn func(Level, Tag) bool) {
	// didn't check for nil, because if it's nil, it will simple remove current one.
	c.ctlFn = fn
}

// check will check if Level and Tag given is good to be printed.
func (c *control) Check(lvl Level, tag Tag) bool {
	switch {
	case c.ctlFn != nil: // FilterFn has the highest order if Set.
		return c.ctlFn(lvl, tag)
	case c.ctlLevel <= lvl: // if a given level is equal to or higher than set level, print
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

// TagBucket can issue a tag and also holds the total number
// of tags issued AND also names given to each tag.
// Not that TagBucket is not using any mutex as it is designed
// to be set at the very beginning of the process.
// Also, the maximum number of tag can be issue is limited to 63.
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
	// If a tag is found, return it.
	if tag, ok := t.GetTag(name); ok {
		return tag
	}
	// If the tag is not found, issue a tag using most recently created.
	// When the maximum capacity of tag has met, return 0.
	if t.count >= 63 {
		return 0
	}
	// Create a new tag and return the tag.
	t.names[t.count] = name
	tag := t.count // this is the value to be printed.
	t.count += 1
	return 1 << tag
}

// AppendSelectedTags is to be used to append selected tags to the byte slice.
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
