package alog

import "time"

type Formatter interface {
	AppendPrefix(dst []byte, prefix []byte) []byte
	AppendTime(dst []byte, format Format) []byte
	AppendLevel(dst []byte, level Level) []byte
	AppendTag(dst []byte, tb *TagBucket, tag Tag) []byte
	AppendMsg(dst []byte, s string) []byte
	AppendMsgBytes(dst []byte, p []byte) []byte
	AppendKVInt(dst []byte, key string, val int) []byte
	AppendKVString(dst []byte, key string, val string) []byte
	AppendKVFloat(dst []byte, key string, val float64) []byte
	AppendKVBool(dst []byte, key string, val bool) []byte
	AppendKVError(dst []byte, key string, val error) []byte
	AppendSuffix(dst []byte, suffix []byte) []byte
	TrimSuffix(dst []byte, c byte) []byte
}

type defaultFormatter struct{}

func (f *defaultFormatter) AppendPrefix(dst []byte, prefix []byte) []byte {
	dst = append(dst, '{')
	if prefix != nil {
		return append(dst, prefix...)
	}
	return dst
}
func (f *defaultFormatter) AppendTime(dst []byte, format Format) []byte {
	dst = conv.EscKey(dst, "ts", true, ':')
	t := time.Now()
	if FtimeUnixMs&format != 0 {
		dst = append(dst, `"ts":`...) // faster without addKey
		return conv.Intf(dst, int(t.UnixNano())/1e6, 0, ',')
	} else if FtimeUnix&format != 0 {
		dst = append(dst, `"ts":`...) // faster without addKey
		return conv.Intf(dst, int(t.Unix()), 0, ',')
	} else {
		if FtimeUTC&format != 0 {
			t = t.UTC()
		}
		if Fdate&format != 0 {
			dst = append(dst, `"d":`...) // faster without addKey
			y, m, d := t.Date()
			dst = conv.Intf(dst, y*10000+int(m)*100+d, 4, ',')
		}
		if FdateDay&format != 0 {
			// "wd": 0 being sunday, 6 being saturday
			dst = append(dst, `"wd":`...) // faster without addKey
			dst = conv.Intf(dst, int(t.Weekday()), 1, ',')
		}
		if Ftime&format != 0 {
			dst = append(dst, `"t":`...) // faster without addKey
			h, m, s := t.Clock()
			dst = conv.Intf(dst, h*10000+m*100+s, 1, '.')
			dst = conv.Intf(dst, t.Nanosecond()/1e6, 3, ',')
		}
	}

	return dst
}
func (f *defaultFormatter) AppendLevel(dst []byte, level Level) []byte {
	dst = conv.EscKey(dst, "lv", true, ':')
	return conv.EscKey(dst, level.String(), true, ',')
}
func (f *defaultFormatter) AppendTag(dst []byte, tb *TagBucket, tag Tag) []byte {
	dst = conv.EscKey(dst, "tag", true, ':')
	dst = tb.AppendSelectedTags(append(dst, '['), ',', true, tag)
	return append(dst, ']', ',')
}
func (f *defaultFormatter) AppendMsg(dst []byte, s string) []byte {
	dst = conv.EscKey(dst, "msg", true, ':')
	return conv.EscString(dst, s, true, ',')
}
func (f *defaultFormatter) AppendMsgBytes(dst []byte, p []byte) []byte {
	dst = conv.EscKey(dst, "msg", true, ':')
	return conv.EscStringBytes(dst, p, true, ',')
}

func (f *defaultFormatter) AppendKVInt(dst []byte, key string, val int) []byte {
	dst = conv.EscKey(dst, key, true, ':')
	return conv.Int(dst, val, false, ',')
}

func (f *defaultFormatter) AppendKVString(dst []byte, key string, val string) []byte {
	dst = conv.EscKey(dst, key, true, ':')
	return conv.EscString(dst, val, true, ',')
}

func (f *defaultFormatter) AppendKVFloat(dst []byte, key string, val float64) []byte {
	dst = conv.EscKey(dst, key, true, ':')
	return conv.Float(dst, val, false, ',')
}

func (f *defaultFormatter) AppendKVBool(dst []byte, key string, val bool) []byte {
	dst = conv.EscKey(dst, key, true, ':')
	if val == true {
		return append(dst, `true,`...)
	}
	return append(dst, `false,`...)
}

func (f *defaultFormatter) AppendKVError(dst []byte, key string, val error) []byte {
	dst = conv.EscKey(dst, key, true, ':')
	if val != nil {
		return conv.EscString(dst, val.Error(), true, ',')
	}
	return append(dst, `null,`...)
}

func (f *defaultFormatter) AppendSuffix(dst []byte, suffix []byte) []byte {
	dst[len(dst)-1] = '}'
	if suffix != nil {
		return append(dst, suffix...)
	}
	return dst
}

func (f *defaultFormatter) TrimSuffix(dst []byte, c byte) []byte {
	if len(dst) > 0 && dst[len(dst)-1] == c {
		return dst[:len(dst)-1]
	}
	return dst
}
