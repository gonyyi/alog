package ext

import (
	"github.com/gonyyi/alog"
	"time"
)

const (
	fmtAnsiCLEAR  = "\033[0m"
	fmtAnsiDIM    = "\033[0;90m"
	fmtAnsiBOLD   = "\033[0;1m"
	fmtAnsiITALIC = "\033[0;1;3m"

	fmtAnsiTRACE = "\033[0;37;100m"
	fmtAnsiDEBUG = "\033[0;42;30m"
	fmtAnsiINFO  = "\033[0;44;30m"
	fmtAnsiWARN  = "\033[0;103;30m"
	fmtAnsiERROR = "\033[0;101;30m"
	fmtAnsiFATAL = "\033[0;1;105;30m"
)

var formatterANSI alog.Formatter

func FormatterANSI() alog.Formatter {
	if formatterANSI == nil {
		formatterANSI = &formatANSI{}
		formatterANSI.Init()
	}
	return formatterANSI
}

// formatANSI is a text formatter
type formatANSI struct {
	conv alog.Converter
}

func (f *formatANSI) Init() {
	f.conv = alog.Defaults.Converter()
}
func (f *formatANSI) Start(dst []byte, prefix []byte) []byte {
	//dst = append(dst, fmtAnsiDIM...)
	//dst = append(dst, "fmtAnsiDIM"...)
	//dst = append(dst, fmtAnsiCLEAR...)
	//dst = append(dst, " Normal"...)
	if prefix != nil {
		return append(dst, prefix...)
	}
	return dst
}

func (f *formatANSI) AppendTime(dst []byte, format alog.Format) []byte {
	t := time.Now()
	if alog.FtimeUnixMs&format != 0 {
		return f.conv.Intf(dst, int(t.UnixNano())/1e6, 0, ' ')
	} else if alog.FtimeUnix&format != 0 {
		return f.conv.Intf(dst, int(t.Unix()), 0, ' ')
	} else {
		if alog.FtimeUTC&format != 0 {
			t = t.UTC()
		}
		if alog.Fdate&format != 0 {
			y, m, d := t.Date()
			dst = append(dst, fmtAnsiDIM...)
			dst = f.conv.Intf(dst, y, 4, '/')
			dst = f.conv.Intf(dst, int(m), 2, '/')
			dst = f.conv.Intf(dst, d, 2, ' ')
			dst = append(dst, fmtAnsiCLEAR...)
		}
		if alog.FdateDay&format != 0 {
			// "wd": 0 being sunday, 6 being saturday
			dst = append(dst, fmtAnsiDIM...)
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
			dst = append(dst, fmtAnsiCLEAR...)
			dst = append(dst, ' ')
		}
		if (alog.Ftime|alog.FtimeMs)&format != 0 {
			h, m, s := t.Clock()
			dst = f.conv.Intf(dst, h, 2, ':')
			dst = f.conv.Intf(dst, m, 2, ':')
			if alog.FtimeMs&format != 0 {
				dst = f.conv.Intf(dst, s, 2, 0)
				dst = append(dst, fmtAnsiDIM...)
				dst = append(dst, ',')
				dst = f.conv.Intf(dst, t.Nanosecond()/1e6, 3, ' ')
				dst = append(dst, fmtAnsiCLEAR...)
			} else {
				dst = f.conv.Intf(dst, s, 2, ' ')
			}
		}
	}
	return dst
}
func (f *formatANSI) AppendLevel(dst []byte, level alog.Level) []byte {
	switch level {
	case alog.Ltrace:
		dst = append(dst, fmtAnsiTRACE...)
	case alog.Ldebug:
		dst = append(dst, fmtAnsiDEBUG...)
	case alog.Linfo:
		dst = append(dst, fmtAnsiINFO...)
	case alog.Lwarn:
		dst = append(dst, fmtAnsiWARN...)
	case alog.Lerror:
		dst = append(dst, fmtAnsiERROR...)
	case alog.Lfatal:
		dst = append(dst, fmtAnsiFATAL...)
	}
	dst = append(dst, ' ')
	dst = append(dst, level.ShortName()...)
	dst = append(dst, ' ')
	dst = append(dst, fmtAnsiCLEAR...)
	return append(dst, ' ')
	// return conv.EscKey(dst, level.ShortName(), false, ' ')
}
func (f *formatANSI) AppendTag(dst []byte, tb *alog.TagBucket, tag alog.Tag) []byte {
	if tag == 0 {
		dst = append(dst, fmtAnsiDIM...)
		dst = append(dst, `[] `...)
		return append(dst, fmtAnsiCLEAR...)
	}
	dst = append(dst, fmtAnsiDIM...)
	dst = append(dst, '[')
	dst = append(dst, fmtAnsiCLEAR...)
	dst = tb.AppendSelectedTags(dst, ',', false, tag)
	dst = append(dst, fmtAnsiDIM...)
	dst = append(dst, ']')
	dst = append(dst, fmtAnsiCLEAR...)
	return append(dst, ' ')
}
func (f *formatANSI) AppendMsg(dst []byte, s string) []byte {
	if len(s) == 0 {
		dst = append(dst, fmtAnsiITALIC...)
		dst = append(dst, `null `...)
		return append(dst, fmtAnsiCLEAR...)
	}
	dst = append(dst, fmtAnsiITALIC...)
	for i := 0; i < len(s); i++ {
		if s[i] != '\n' {
			dst = append(dst, s[i])
		} else {
			dst = append(dst, ';')
		}
	}
	dst = append(dst, fmtAnsiCLEAR...)
	return append(dst, ' ')
}

func (f *formatANSI) AppendMsgBytes(dst []byte, p []byte) []byte {
	for i := 0; i < len(p); i++ {
		if p[i] != '\n' {
			dst = append(dst, p[i])
		} else {
			dst = append(dst, ' ')
		}
	}

	return append(dst, ' ') // return conv.EscStringBytes(dst, p, false, ' ')
}
func (f *formatANSI) AppendSeparator(dst []byte) []byte {
	return append(dst, `// `...)
}
func (f *formatANSI) AppendKVInt(dst []byte, key string, val int) []byte {
	dst = append(dst, fmtAnsiDIM...)
	dst = f.conv.EscKey(dst, key, false, '=')
	dst = append(dst, fmtAnsiBOLD...)
	dst = f.conv.Int(dst, val, false, 0)
	dst = append(dst, fmtAnsiCLEAR...)
	return append(dst, ' ')
}

func (f *formatANSI) AppendKVString(dst []byte, key string, val string) []byte {
	dst = append(dst, fmtAnsiDIM...)
	dst = f.conv.EscKey(dst, key, false, '=')
	dst = append(dst, fmtAnsiBOLD...)
	dst = f.conv.EscString(dst, val, true, 0)
	dst = append(dst, fmtAnsiCLEAR...)
	return append(dst, ' ')
}

func (f *formatANSI) AppendKVFloat(dst []byte, key string, val float64) []byte {
	dst = append(dst, fmtAnsiDIM...)
	dst = f.conv.EscKey(dst, key, false, '=')
	dst = append(dst, fmtAnsiBOLD...)
	dst = f.conv.Float(dst, val, false, 0)
	dst = append(dst, fmtAnsiCLEAR...)
	return append(dst, ' ')
}

func (f *formatANSI) AppendKVBool(dst []byte, key string, val bool) []byte {
	dst = append(dst, fmtAnsiDIM...)
	dst = f.conv.EscKey(dst, key, false, '=')
	dst = append(dst, fmtAnsiBOLD...)
	if val == true {
		dst = append(dst, `true`...)
		dst = append(dst, fmtAnsiCLEAR...)
		return append(dst, ' ')
	}
	dst = append(dst, `false`...)
	dst = append(dst, fmtAnsiCLEAR...)
	return append(dst, ' ')
}

func (f *formatANSI) AppendKVError(dst []byte, key string, val error) []byte {
	dst = append(dst, fmtAnsiDIM...)
	dst = f.conv.EscKey(dst, key, false, '=')
	dst = append(dst, fmtAnsiBOLD...)
	if val != nil {
		dst = f.conv.EscString(dst, val.Error(), true, 0)
		dst = append(dst, fmtAnsiCLEAR...)
		return append(dst, ' ')
	}
	dst = append(dst, `nil`...)
	dst = append(dst, fmtAnsiCLEAR...)
	return append(dst, ' ')
}

func (f *formatANSI) AppendSuffix(dst []byte, suffix []byte) []byte {
	if suffix != nil {
		return append(dst, suffix...)
	}
	return dst
}

func (f *formatANSI) Final(dst, suffix []byte) []byte {
	if len(dst) > 0 { // only do this if dst exists,
		dst = dst[:len(dst)-1] // trim last space
	}
	if suffix != nil {
		return append(dst, suffix...)
	}
	return dst
}
