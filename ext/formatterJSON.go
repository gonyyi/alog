package ext

import (
	"github.com/gonyyi/alog"
	"time"
)

var formatterJson alog.Formatter

func FormatterJSON() alog.Formatter {
	if formatterJson == nil {
		formatterJson = &formatJSON{}
		formatterJson.Init()
	}
	return formatterJson
}

// formatJSON is a formatter struct for JSON
type formatJSON struct {
	conv alog.Converter
}

func (f *formatJSON) Init() {
	f.conv = alog.Conf.Converter()
}
func (f *formatJSON) Start(dst []byte, prefix []byte) []byte {
	if prefix != nil {
		dst = append(dst, prefix...)
	}
	return append(dst, '{')
}
func (f *formatJSON) AppendTime(dst []byte, format alog.Format) []byte {
	t := time.Now()
	if alog.FtimeUnixMs&format != 0 {
		dst = append(dst, `"ts":`...) // faster without addKey
		return f.conv.Intf(dst, int(t.UnixNano())/1e6, 0, ',')
	} else if alog.FtimeUnix&format != 0 {
		dst = append(dst, `"ts":`...) // faster without addKey
		return f.conv.Intf(dst, int(t.Unix()), 0, ',')
	} else {
		if alog.FtimeUTC&format != 0 {
			t = t.UTC()
		}
		if alog.Fdate&format != 0 {
			dst = append(dst, `"d":`...) // faster without addKey
			y, m, d := t.Date()
			dst = f.conv.Intf(dst, y*10000+int(m)*100+d, 4, ',')
		}
		if alog.FdateDay&format != 0 {
			// "wd": 0 being sunday, 6 being saturday
			dst = append(dst, `"wd":`...) // faster without addKey
			dst = f.conv.Intf(dst, int(t.Weekday()), 1, ',')
		}
		if (alog.Ftime|alog.FtimeMs)&format != 0 {
			dst = append(dst, `"t":`...) // faster without addKey
			h, m, s := t.Clock()
			if alog.FtimeMs&format != 0 {
				dst = f.conv.Intf(dst, h*10000+m*100+s, 1, 0)
				dst = f.conv.Intf(dst, t.Nanosecond()/1e6, 3, ',')
			} else {
				dst = f.conv.Intf(dst, h*10000+m*100+s, 1, ',')
			}
		}
	}

	return dst
}
func (f *formatJSON) AppendLevel(dst []byte, level alog.Level) []byte {
	dst = append(dst, `"lv":`...)
	return f.conv.EscKey(dst, level.String(), true, ',')
}
func (f *formatJSON) AppendTag(dst []byte, tb *alog.TagBucket, tag alog.Tag) []byte {
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
	return f.conv.EscString(dst, s, true, ',')
}
func (f *formatJSON) AppendSeparator(dst []byte) []byte {
	return dst
}
func (f *formatJSON) AppendMsgBytes(dst []byte, p []byte) []byte {
	dst = append(dst, `"msg":`...)
	return f.conv.EscStringBytes(dst, p, true, ',')
}

func (f *formatJSON) AppendKVInt(dst []byte, key string, val int) []byte {
	dst = f.conv.EscKey(dst, key, true, ':')
	return f.conv.Int(dst, val, false, ',')
}

func (f *formatJSON) AppendKVString(dst []byte, key string, val string) []byte {
	dst = f.conv.EscKey(dst, key, true, ':')
	return f.conv.EscString(dst, val, true, ',')
}

func (f *formatJSON) AppendKVFloat(dst []byte, key string, val float64) []byte {
	dst = f.conv.EscKey(dst, key, true, ':')
	return f.conv.Float(dst, val, false, ',')
}

func (f *formatJSON) AppendKVBool(dst []byte, key string, val bool) []byte {
	dst = f.conv.EscKey(dst, key, true, ':')
	if val == true {
		return append(dst, `true,`...)
	}
	return append(dst, `false,`...)
}

func (f *formatJSON) AppendKVError(dst []byte, key string, val error) []byte {
	dst = f.conv.EscKey(dst, key, true, ':')
	if val != nil {
		return f.conv.EscString(dst, val.Error(), true, ',')
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
