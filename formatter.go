package alog

import "time"

// Format a bit-formatFlag formatFlag options that is used for variety of configuration.
type Format uint32

// On will turn flags on for the given item(s)
func (f Format) On(item Format) Format {
	return f | item
}

// Off will turn flags off for the given item(s)
func (f Format) Off(item Format) Format {
	return f &^ item
}

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

	// Fdefault will show month/day with time, and Level of logging.
	Fdefault = Fdate | Ftime | Flevel | Ftag
	// fUseTime is precalculated time for internal functions.
	fUseTime = Fdate | FdateDay | Ftime | FtimeMs | FtimeUnix | FtimeUnixMs
)

// Formatter interface allows Alog to have different format of output.
// Default formatter in Alog is set to formatText, but also has formatJSON built in.
type Formatter interface {
	Init()
	Start(dst []byte, prefix []byte) []byte                   // Start to be used at starting of any log message.
	AppendTime(dst []byte, format Format) []byte              // AppendTime will take format flag and add formatted time.
	AppendLevel(dst []byte, level Level) []byte               // AppendLevel will add level string.
	AppendTag(dst []byte, tb *TagBucket, tag Tag) []byte      // AppendTag will add a tag to the log.
	AppendMsg(dst []byte, s string) []byte                    // AppendMsg will add a main message. (For JSON, "msg")
	AppendMsgBytes(dst []byte, p []byte) []byte               // AppendMsgBytes is same as AppendMsg but take a byte slice.
	AppendSeparator(dst []byte) []byte                        // AppendSeparator will add separate key if any
	AppendKVInt(dst []byte, key string, val int) []byte       // AppendKVInt will add key/value for integer
	AppendKVString(dst []byte, key string, val string) []byte // AppendKVString will add key/value for string
	AppendKVFloat(dst []byte, key string, val float64) []byte // AppendKVFloat will add key/value for float64
	AppendKVBool(dst []byte, key string, val bool) []byte     // AppendKVBool will add key/value for boolean
	AppendKVError(dst []byte, key string, val error) []byte   // AppendKVError will add error value with a key
	Final(dst []byte, suffix []byte) []byte                   // Final to be used at the end of each log message
}

// formatText is a text formatter
type formatText struct {
	conv Converter
}

func (f *formatText) Init() {
	if f.conv == nil {
		f.conv = Defaults.Converter()
	}
}

func (f *formatText) Start(dst []byte, prefix []byte) []byte {
	if prefix != nil {
		return append(dst, prefix...)
	}
	return dst
}

func (f *formatText) AppendTime(dst []byte, format Format) []byte {
	t := time.Now()
	if FtimeUnixMs&format != 0 {
		return Defaults.converter.Intf(dst, int(t.UnixNano())/1e6, 0, ' ')
	} else if FtimeUnix&format != 0 {
		return Defaults.converter.Intf(dst, int(t.Unix()), 0, ' ')
	} else {
		if FtimeUTC&format != 0 {
			t = t.UTC()
		}
		if Fdate&format != 0 {
			y, m, d := t.Date()
			dst = Defaults.converter.Intf(dst, y, 4, '/')
			dst = Defaults.converter.Intf(dst, int(m), 2, '/')
			dst = Defaults.converter.Intf(dst, d, 2, ' ')
		}
		if FdateDay&format != 0 {
			// "wd": 0 being sunday, 6 being saturday
			switch t.Weekday() {
			case 0:
				dst = append(dst, `Sun`...)
			case 1:
				dst = append(dst, `Mon`...)
			case 2:
				dst = append(dst, `Tue`...)
			case 3:
				dst = append(dst, `Wed`...)
			case 4:
				dst = append(dst, `Thu`...)
			case 5:
				dst = append(dst, `Fri`...)
			case 6:
				dst = append(dst, `Sat`...)
			}
			dst = append(dst, ' ')
		}
		if (Ftime|FtimeMs)&format != 0 {
			h, m, s := t.Clock()
			dst = Defaults.converter.Intf(dst, h, 2, ':')
			dst = Defaults.converter.Intf(dst, m, 2, ':')
			if FtimeMs&format != 0 {
				dst = Defaults.converter.Intf(dst, s, 2, ',')
				dst = Defaults.converter.Intf(dst, t.Nanosecond()/1e6, 3, ' ')
			} else {
				dst = Defaults.converter.Intf(dst, s, 2, ' ')
			}
		}
	}
	return dst
}
func (f *formatText) AppendLevel(dst []byte, level Level) []byte {
	dst = append(dst, level.ShortName()...)
	return append(dst, ' ')
	// return conv.EscKey(dst, level.ShortName(), false, ' ')
}
func (f *formatText) AppendTag(dst []byte, tb *TagBucket, tag Tag) []byte {
	if tag == 0 {
		return append(dst, `[] `...)
	}
	dst = tb.AppendSelectedTags(append(dst, '['), ',', false, tag)
	return append(dst, ']', ' ')
}
func (f *formatText) AppendMsg(dst []byte, s string) []byte {
	if len(s) == 0 {
		return append(dst, `null `...)
	}
	for i := 0; i < len(s); i++ {
		if s[i] != '\n' {
			dst = append(dst, s[i])
		} else {
			dst = append(dst, ';')
		}
	}

	return append(dst, ' ') // return conv.EscString(dst, s, false, ' ')
}

func (f *formatText) AppendMsgBytes(dst []byte, p []byte) []byte {
	for i := 0; i < len(p); i++ {
		if p[i] != '\n' {
			dst = append(dst, p[i])
		} else {
			dst = append(dst, ' ')
		}
	}

	return append(dst, ' ') // return conv.EscStringBytes(dst, p, false, ' ')
}
func (f *formatText) AppendSeparator(dst []byte) []byte {
	return append(dst, `// `...)
}
func (f *formatText) AppendKVInt(dst []byte, key string, val int) []byte {
	dst = Defaults.converter.EscKey(dst, key, false, '=')
	return Defaults.converter.Int(dst, val, false, ',')
}

func (f *formatText) AppendKVString(dst []byte, key string, val string) []byte {
	dst = Defaults.converter.EscKey(dst, key, false, '=')
	return Defaults.converter.EscString(dst, val, true, ',')
}

func (f *formatText) AppendKVFloat(dst []byte, key string, val float64) []byte {
	dst = Defaults.converter.EscKey(dst, key, false, '=')
	return Defaults.converter.Float(dst, val, false, ',')
}

func (f *formatText) AppendKVBool(dst []byte, key string, val bool) []byte {
	dst = Defaults.converter.EscKey(dst, key, false, '=')
	if val == true {
		return append(dst, `true,`...)
	}
	return append(dst, `false,`...)
}

func (f *formatText) AppendKVError(dst []byte, key string, val error) []byte {
	dst = Defaults.converter.EscKey(dst, key, false, '=')
	if val != nil {
		return Defaults.converter.EscString(dst, val.Error(), true, ',')
	}
	return append(dst, `nil,`...)
}

func (f *formatText) AppendSuffix(dst []byte, suffix []byte) []byte {
	if suffix != nil {
		return append(dst, suffix...)
	}
	return dst
}

func (f *formatText) Final(dst, suffix []byte) []byte {
	if len(dst) > 0 { // only do this if dst exists,
		dst = dst[:len(dst)-1] // trim last space
	}
	if suffix != nil {
		return append(dst, suffix...)
	}
	return dst
}
