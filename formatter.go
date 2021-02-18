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
	FoutJSON           // FoutJSON will print to a JSON
	FoutText           // FoutText will print to a TEXT

	// Fdefault will show month/day with time, and Level of logging.
	Fdefault = Fdate | Ftime | Flevel | Ftag | FoutText
	// fUseTime is precalculated time for internal functions.
	fUseTime = Fdate | FdateDay | Ftime | FtimeMs | FtimeUnix | FtimeUnixMs
)

// Formatter interface allows Alog to have different format of output.
// Default formatter in Alog is set to formatText, but also has formatJSON built in.
type Formatter interface {
	Start(dst []byte, prefix []byte) []byte                   // Start to be used at starting of any log message.
	AppendTime(dst []byte, format Format) []byte              // AppendTime will take format flag and add formatted time.
	AppendLevel(dst []byte, level Level) []byte               // AppendLevel will add level string.
	AppendTag(dst []byte, tb *TagBucket, tag Tag) []byte      // AppendTag will add a tag to the log.
	AppendMsg(dst []byte, s string) []byte                    // AppendMsg will add a main message. (For JSON, "msg")
	AppendMsgBytes(dst []byte, p []byte) []byte               // AppendMsgBytes is same as AppendMsg but take a byte slice.
	AppendKVInt(dst []byte, key string, val int) []byte       // AppendKVInt will add key/value for integer
	AppendKVString(dst []byte, key string, val string) []byte // AppendKVString will add key/value for string
	AppendKVFloat(dst []byte, key string, val float64) []byte // AppendKVFloat will add key/value for float64
	AppendKVBool(dst []byte, key string, val bool) []byte     // AppendKVBool will add key/value for boolean
	AppendKVError(dst []byte, key string, val error) []byte   // AppendKVError will add error value with a key
	Final(dst []byte, suffix []byte) []byte                   // Final to be used at the end of each log message
}

// formatJSON is a formatter struct for JSON
type formatJSON struct{}

func (f *formatJSON) Start(dst []byte, prefix []byte) []byte {
	if prefix != nil {
		dst = append(dst, prefix...)
	}
	return append(dst, '{')
}
func (f *formatJSON) AppendTime(dst []byte, format Format) []byte {
	t := time.Now()
	if FtimeUnixMs&format != 0 {
		dst = append(dst, `"ts":`...) // faster without addKey
		return Defaults.Converter.Intf(dst, int(t.UnixNano())/1e6, 0, ',')
	} else if FtimeUnix&format != 0 {
		dst = append(dst, `"ts":`...) // faster without addKey
		return Defaults.Converter.Intf(dst, int(t.Unix()), 0, ',')
	} else {
		if FtimeUTC&format != 0 {
			t = t.UTC()
		}
		if Fdate&format != 0 {
			dst = append(dst, `"d":`...) // faster without addKey
			y, m, d := t.Date()
			dst = Defaults.Converter.Intf(dst, y*10000+int(m)*100+d, 4, ',')
		}
		if FdateDay&format != 0 {
			// "wd": 0 being sunday, 6 being saturday
			dst = append(dst, `"wd":`...) // faster without addKey
			dst = Defaults.Converter.Intf(dst, int(t.Weekday()), 1, ',')
		}
		if (Ftime|FtimeMs)&format != 0 {
			dst = append(dst, `"t":`...) // faster without addKey
			h, m, s := t.Clock()
			if FtimeMs&format != 0 {
				dst = Defaults.Converter.Intf(dst, h*10000+m*100+s, 1, 0)
				dst = Defaults.Converter.Intf(dst, t.Nanosecond()/1e6, 0, ',')
			} else {
				dst = Defaults.Converter.Intf(dst, h*10000+m*100+s, 1, ',')
			}
		}
	}

	return dst
}
func (f *formatJSON) AppendLevel(dst []byte, level Level) []byte {
	dst = append(dst, `"lv":`...)
	return Defaults.Converter.EscKey(dst, level.String(), true, ',')
}
func (f *formatJSON) AppendTag(dst []byte, tb *TagBucket, tag Tag) []byte {
	if tag == 0 {
		return append(dst, `"tag":[],`...)
	}
	dst = append(dst, `"tag":[`...)
	dst = tb.AppendSelectedTags(dst, ',', true, tag)
	return append(dst, ']', ',')
}
func (f *formatJSON) AppendMsg(dst []byte, s string) []byte {
	if len(s) == 0 {
		return append(dst, `"msg":null,`...)
	}
	dst = append(dst, `"msg":`...)
	return Defaults.Converter.EscString(dst, s, true, ',')
}

func (f *formatJSON) AppendMsgBytes(dst []byte, p []byte) []byte {
	dst = append(dst, `"msg":`...)
	return Defaults.Converter.EscStringBytes(dst, p, true, ',')
}

func (f *formatJSON) AppendKVInt(dst []byte, key string, val int) []byte {
	dst = Defaults.Converter.EscKey(dst, key, true, ':')
	return Defaults.Converter.Int(dst, val, false, ',')
}

func (f *formatJSON) AppendKVString(dst []byte, key string, val string) []byte {
	dst = Defaults.Converter.EscKey(dst, key, true, ':')
	return Defaults.Converter.EscString(dst, val, true, ',')
}

func (f *formatJSON) AppendKVFloat(dst []byte, key string, val float64) []byte {
	dst = Defaults.Converter.EscKey(dst, key, true, ':')
	return Defaults.Converter.Float(dst, val, false, ',')
}

func (f *formatJSON) AppendKVBool(dst []byte, key string, val bool) []byte {
	dst = Defaults.Converter.EscKey(dst, key, true, ':')
	if val == true {
		return append(dst, `true,`...)
	}
	return append(dst, `false,`...)
}

func (f *formatJSON) AppendKVError(dst []byte, key string, val error) []byte {
	dst = Defaults.Converter.EscKey(dst, key, true, ':')
	if val != nil {
		return Defaults.Converter.EscString(dst, val.Error(), true, ',')
	}
	return append(dst, `null,`...)
}

func (f *formatJSON) AppendSuffix(dst []byte, suffix []byte) []byte {
	if suffix != nil {
		return append(dst, suffix...)
	}
	return dst
}

func (f *formatJSON) Final(dst, suffix []byte) []byte {
	if len(dst) > 0 { // only do this if dst exists,
		dst[len(dst)-1] = '}'
	}
	if suffix != nil {
		return append(dst, suffix...)
	}
	return dst
}

// formatText is a text formatter
type formatText struct{}

func (f *formatText) Start(dst []byte, prefix []byte) []byte {
	if prefix != nil {
		return append(dst, prefix...)
	}
	return dst
}

func (f *formatText) AppendTime(dst []byte, format Format) []byte {
	t := time.Now()
	if FtimeUnixMs&format != 0 {
		return Defaults.Converter.Intf(dst, int(t.UnixNano())/1e6, 0, ' ')
	} else if FtimeUnix&format != 0 {
		return Defaults.Converter.Intf(dst, int(t.Unix()), 0, ' ')
	} else {
		if FtimeUTC&format != 0 {
			t = t.UTC()
		}
		if Fdate&format != 0 {
			y, m, d := t.Date()
			dst = Defaults.Converter.Intf(dst, y, 4, '/')
			dst = Defaults.Converter.Intf(dst, int(m), 2, '/')
			dst = Defaults.Converter.Intf(dst, d, 2, ' ')
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
			dst = Defaults.Converter.Intf(dst, h, 2, ':')
			dst = Defaults.Converter.Intf(dst, m, 2, ':')
			if FtimeMs&format != 0 {
				dst = Defaults.Converter.Intf(dst, s, 2, ',')
				dst = Defaults.Converter.Intf(dst, t.Nanosecond()/1e6, 3, ' ')
			} else {
				dst = Defaults.Converter.Intf(dst, s, 2, ' ')
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

func (f *formatText) AppendKVInt(dst []byte, key string, val int) []byte {
	dst = Defaults.Converter.EscKey(dst, key, false, '=')
	return Defaults.Converter.Int(dst, val, false, ',')
}

func (f *formatText) AppendKVString(dst []byte, key string, val string) []byte {
	dst = Defaults.Converter.EscKey(dst, key, false, '=')
	return Defaults.Converter.EscString(dst, val, true, ',')
}

func (f *formatText) AppendKVFloat(dst []byte, key string, val float64) []byte {
	dst = Defaults.Converter.EscKey(dst, key, false, '=')
	return Defaults.Converter.Float(dst, val, false, ',')
}

func (f *formatText) AppendKVBool(dst []byte, key string, val bool) []byte {
	dst = Defaults.Converter.EscKey(dst, key, false, '=')
	if val == true {
		return append(dst, `true,`...)
	}
	return append(dst, `false,`...)
}

func (f *formatText) AppendKVError(dst []byte, key string, val error) []byte {
	dst = Defaults.Converter.EscKey(dst, key, false, '=')
	if val != nil {
		return Defaults.Converter.EscString(dst, val.Error(), true, ',')
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
