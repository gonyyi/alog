package ext

import (
	"github.com/gonyyi/alog"
	"time"
)

var formatterText alog.Formatter

func FormatterText() alog.Formatter {
	if formatterText == nil {
		formatterText = &formatText{}
		formatterText.Init()
	}
	return formatterText
}

// formatText is a text formatter
type formatText struct {
	conv alog.Converter
}

func (f *formatText) Init() {
	if f.conv == nil {
		f.conv = alog.Defaults.Converter()
	}
}

func (f *formatText) Start(dst []byte, prefix []byte) []byte {
	if prefix != nil {
		return append(dst, prefix...)
	}
	return dst
}

func (f *formatText) AppendTime(dst []byte, format alog.Format) []byte {
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
			dst = f.conv.Intf(dst, y, 4, '/')
			dst = f.conv.Intf(dst, int(m), 2, '/')
			dst = f.conv.Intf(dst, d, 2, ' ')
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
			dst = f.conv.Intf(dst, h, 2, ':')
			dst = f.conv.Intf(dst, m, 2, ':')
			if alog.FtimeMs&format != 0 {
				dst = f.conv.Intf(dst, s, 2, ',')
				dst = f.conv.Intf(dst, t.Nanosecond()/1e6, 3, ' ')
			} else {
				dst = f.conv.Intf(dst, s, 2, ' ')
			}
		}
	}
	return dst
}
func (f *formatText) AppendLevel(dst []byte, level alog.Level) []byte {
	dst = append(dst, level.ShortName()...)
	return append(dst, ' ')
	// return conv.EscKey(dst, level.ShortName(), false, ' ')
}
func (f *formatText) AppendTag(dst []byte, tb *alog.TagBucket, tag alog.Tag) []byte {
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
	dst = f.conv.EscKey(dst, key, false, '=')
	return f.conv.Int(dst, val, false, ',')
}

func (f *formatText) AppendKVString(dst []byte, key string, val string) []byte {
	dst = f.conv.EscKey(dst, key, false, '=')
	return f.conv.EscString(dst, val, true, ',')
}

func (f *formatText) AppendKVFloat(dst []byte, key string, val float64) []byte {
	dst = f.conv.EscKey(dst, key, false, '=')
	return f.conv.Float(dst, val, false, ',')
}

func (f *formatText) AppendKVBool(dst []byte, key string, val bool) []byte {
	dst = f.conv.EscKey(dst, key, false, '=')
	if val == true {
		return append(dst, `true,`...)
	}
	return append(dst, `false,`...)
}

func (f *formatText) AppendKVError(dst []byte, key string, val error) []byte {
	dst = f.conv.EscKey(dst, key, false, '=')
	if val != nil {
		return f.conv.EscString(dst, val.Error(), true, ',')
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
