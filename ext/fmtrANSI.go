package ext

import (
	"github.com/gonyyi/alog"
	"time"
)

const (
	CLEAR  = "\033[0m"
	DIM    = "\033[0;90m"
	BOLD   = "\033[0;1m"
	ITALIC = "\033[0;1;3m"

	TRACE = "\033[0;37;100m"
	DEBUG = "\033[0;42;30m"
	INFO  = "\033[0;44;30m"
	WARN  = "\033[0;103;30m"
	ERROR = "\033[0;101;30m"
	FATAL = "\033[0;1;105;30m"
)

func NewFormatterANSI() *formatANSI {
	return &formatANSI{}
}

// formatANSI is a text formatter
type formatANSI struct{}

func (f *formatANSI) Start(dst []byte, prefix []byte) []byte {
	//dst = append(dst, DIM...)
	//dst = append(dst, "DIM"...)
	//dst = append(dst, CLEAR...)
	//dst = append(dst, " Normal"...)
	if prefix != nil {
		return append(dst, prefix...)
	}
	return dst
}

func (f *formatANSI) AppendTime(dst []byte, format alog.Format) []byte {
	t := time.Now()
	if alog.FtimeUnixMs&format != 0 {
		return alog.Defaults.Converter.Intf(dst, int(t.UnixNano())/1e6, 0, ' ')
	} else if alog.FtimeUnix&format != 0 {
		return alog.Defaults.Converter.Intf(dst, int(t.Unix()), 0, ' ')
	} else {
		if alog.FtimeUTC&format != 0 {
			t = t.UTC()
		}
		if alog.Fdate&format != 0 {
			y, m, d := t.Date()
			dst = append(dst, DIM...)
			dst = alog.Defaults.Converter.Intf(dst, y, 4, '/')
			dst = alog.Defaults.Converter.Intf(dst, int(m), 2, '/')
			dst = alog.Defaults.Converter.Intf(dst, d, 2, ' ')
			dst = append(dst, CLEAR...)
		}
		if alog.FdateDay&format != 0 {
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
		if (alog.Ftime|alog.FtimeMs)&format != 0 {
			h, m, s := t.Clock()
			dst = alog.Defaults.Converter.Intf(dst, h, 2, ':')
			dst = alog.Defaults.Converter.Intf(dst, m, 2, ':')
			if alog.FtimeMs&format != 0 {
				dst = alog.Defaults.Converter.Intf(dst, s, 2, 0)
				dst = append(dst, DIM...)
				dst = append(dst, ',')
				dst = alog.Defaults.Converter.Intf(dst, t.Nanosecond()/1e6, 3, ' ')
				dst = append(dst, CLEAR...)
			} else {
				dst = alog.Defaults.Converter.Intf(dst, s, 2, ' ')
			}
		}
	}
	return dst
}
func (f *formatANSI) AppendLevel(dst []byte, level alog.Level) []byte {
	switch level {
	case alog.Ltrace:
		dst = append(dst, TRACE...)
	case alog.Ldebug:
		dst = append(dst, DEBUG...)
	case alog.Linfo:
		dst = append(dst, INFO...)
	case alog.Lwarn:
		dst = append(dst, WARN...)
	case alog.Lerror:
		dst = append(dst, ERROR...)
	case alog.Lfatal:
		dst = append(dst, FATAL...)
	}
	dst = append(dst, ' ')
	dst = append(dst, level.ShortName()...)
	dst = append(dst, ' ')
	dst = append(dst, CLEAR...)
	return append(dst, ' ')
	// return conv.EscKey(dst, level.ShortName(), false, ' ')
}
func (f *formatANSI) AppendTag(dst []byte, tb *alog.TagBucket, tag alog.Tag) []byte {
	if tag == 0 {
		dst = append(dst, DIM...)
		dst = append(dst, `[] `...)
		return append(dst, CLEAR...)
	}
	dst = append(dst, DIM...)
	dst = append(dst, '[')
	dst = append(dst, CLEAR...)
	dst = tb.AppendSelectedTags(dst, ',', false, tag)
	dst = append(dst, DIM...)
	dst = append(dst, ']')
	dst = append(dst, CLEAR...)
	return append(dst, ' ')
}
func (f *formatANSI) AppendMsg(dst []byte, s string) []byte {
	if len(s) == 0 {
		dst = append(dst, ITALIC...)
		dst = append(dst, `null `...)
		return append(dst, CLEAR...)
	}
	dst = append(dst, ITALIC...)
	for i := 0; i < len(s); i++ {
		if s[i] != '\n' {
			dst = append(dst, s[i])
		} else {
			dst = append(dst, ';')
		}
	}
	dst = append(dst, CLEAR...)
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

func (f *formatANSI) AppendKVInt(dst []byte, key string, val int) []byte {
	dst = append(dst, DIM...)
	dst = alog.Defaults.Converter.EscKey(dst, key, false, '=')
	dst = append(dst, BOLD...)
	dst = alog.Defaults.Converter.Int(dst, val, false, 0)
	dst = append(dst, CLEAR...)
	return append(dst, ' ')
}

func (f *formatANSI) AppendKVString(dst []byte, key string, val string) []byte {
	dst = append(dst, DIM...)
	dst = alog.Defaults.Converter.EscKey(dst, key, false, '=')
	dst = append(dst, BOLD...)
	dst = alog.Defaults.Converter.EscString(dst, val, true, 0)
	dst = append(dst, CLEAR...)
	return append(dst, ' ')
}

func (f *formatANSI) AppendKVFloat(dst []byte, key string, val float64) []byte {
	dst = append(dst, DIM...)
	dst = alog.Defaults.Converter.EscKey(dst, key, false, '=')
	dst = append(dst, BOLD...)
	dst = alog.Defaults.Converter.Float(dst, val, false, 0)
	dst = append(dst, CLEAR...)
	return append(dst, ' ')
}

func (f *formatANSI) AppendKVBool(dst []byte, key string, val bool) []byte {
	dst = append(dst, DIM...)
	dst = alog.Defaults.Converter.EscKey(dst, key, false, '=')
	dst = append(dst, BOLD...)
	if val == true {
		dst = append(dst, `true`...)
		dst = append(dst, CLEAR...)
		return append(dst, ' ')
	}
	dst = append(dst, `false`...)
	dst = append(dst, CLEAR...)
	return append(dst, ' ')
}

func (f *formatANSI) AppendKVError(dst []byte, key string, val error) []byte {
	dst = append(dst, DIM...)
	dst = alog.Defaults.Converter.EscKey(dst, key, false, '=')
	dst = append(dst, BOLD...)
	if val != nil {
		dst = alog.Defaults.Converter.EscString(dst, val.Error(), true, 0)
		dst = append(dst, CLEAR...)
		return append(dst, ' ')
	}
	dst = append(dst, `nil`...)
	dst = append(dst, CLEAR...)
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
