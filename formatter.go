package alog

import "time"

type Formatter interface {
	Start(dst []byte, prefix []byte) []byte
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
	Final(dst []byte, suffix []byte) []byte
}

var FmtrJSON *formatterJSON = &formatterJSON{}
var FmtrText *formatterText = &formatterText{}

type formatterJSON struct{}

func (f *formatterJSON) Start(dst []byte, prefix []byte) []byte {
	dst = append(dst, '{')
	if prefix != nil {
		return append(dst, prefix...)
	}
	return dst
}
func (f *formatterJSON) AppendTime(dst []byte, format Format) []byte {
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
func (f *formatterJSON) AppendLevel(dst []byte, level Level) []byte {
	dst = append(dst, `"lv":`...)
	return conv.EscKey(dst, level.String(), true, ',')
}
func (f *formatterJSON) AppendTag(dst []byte, tb *TagBucket, tag Tag) []byte {
	if tag == 0 {
		return append(dst, `"tag":[],`...)
	}
	dst = append(dst, `"tag":[`...)
	dst = tb.AppendSelectedTags(dst, ',', true, tag)
	return append(dst, ']', ',')
}
func (f *formatterJSON) AppendMsg(dst []byte, s string) []byte {
	dst = append(dst, `"msg":`...)
	return conv.EscString(dst, s, true, ',')
}
func (f *formatterJSON) AppendMsgBytes(dst []byte, p []byte) []byte {
	dst = append(dst, `"msg":`...)
	return conv.EscStringBytes(dst, p, true, ',')
}

func (f *formatterJSON) AppendKVInt(dst []byte, key string, val int) []byte {
	dst = conv.EscKey(dst, key, true, ':')
	return conv.Int(dst, val, false, ',')
}

func (f *formatterJSON) AppendKVString(dst []byte, key string, val string) []byte {
	dst = conv.EscKey(dst, key, true, ':')
	return conv.EscString(dst, val, true, ',')
}

func (f *formatterJSON) AppendKVFloat(dst []byte, key string, val float64) []byte {
	dst = conv.EscKey(dst, key, true, ':')
	return conv.Float(dst, val, false, ',')
}

func (f *formatterJSON) AppendKVBool(dst []byte, key string, val bool) []byte {
	dst = conv.EscKey(dst, key, true, ':')
	if val == true {
		return append(dst, `true,`...)
	}
	return append(dst, `false,`...)
}

func (f *formatterJSON) AppendKVError(dst []byte, key string, val error) []byte {
	dst = conv.EscKey(dst, key, true, ':')
	if val != nil {
		return conv.EscString(dst, val.Error(), true, ',')
	}
	return append(dst, `null,`...)
}

func (f *formatterJSON) AppendSuffix(dst []byte, suffix []byte) []byte {
	if suffix != nil {
		return append(dst, suffix...)
	}
	return dst
}

func (f *formatterJSON) Final(dst, suffix []byte) []byte {
	if len(dst) > 0 { // only do this if dst exists,
		dst[len(dst)-1] = '}'
	}
	if suffix != nil {
		return append(dst, suffix...)
	}
	return dst
}

// formatterText is a text formatter
type formatterText struct{}

func (f *formatterText) Start(dst []byte, prefix []byte) []byte {
	if prefix != nil {
		return append(dst, prefix...)
	}
	return dst
}
func (f *formatterText) AppendTime(dst []byte, format Format) []byte {
	t := time.Now()
	if FtimeUnixMs&format != 0 {
		return conv.Intf(dst, int(t.UnixNano())/1e6, 0, ' ')
	} else if FtimeUnix&format != 0 {
		return conv.Intf(dst, int(t.Unix()), 0, ' ')
	} else {
		if FtimeUTC&format != 0 {
			t = t.UTC()
		}
		if Fdate&format != 0 {
			y, m, d := t.Date()
			dst = conv.Intf(dst, y, 4, '/')
			dst = conv.Intf(dst, int(m), 2, '/')
			dst = conv.Intf(dst, d, 2, ' ')
		}
		if FdateDay&format != 0 {
			// "wd": 0 being sunday, 6 being saturday
			dst = append(dst, t.Weekday().String()[:3]...)
			dst = append(dst, ' ')
		}
		if Ftime&format != 0 {
			h, m, s := t.Clock()
			dst = conv.Intf(dst, h, 2, ':')
			dst = conv.Intf(dst, m, 2, ':')
			dst = conv.Intf(dst, s, 2, ',')
			dst = conv.Intf(dst, t.Nanosecond()/1e6, 3, ' ')
		}
	}
	return dst
}
func (f *formatterText) AppendLevel(dst []byte, level Level) []byte {
	dst = append(dst, level.ShortName()...)
	return append(dst, ' ')
	//return conv.EscKey(dst, level.ShortName(), false, ' ')
}
func (f *formatterText) AppendTag(dst []byte, tb *TagBucket, tag Tag) []byte {
	if tag == 0 {
		return append(dst, `[] `...)
	}
	dst = tb.AppendSelectedTags(append(dst, '['), ',', false, tag)
	return append(dst, ']', ' ')
}
func (f *formatterText) AppendMsg(dst []byte, s string) []byte {
	return conv.EscString(dst, s, false, ' ')
}
func (f *formatterText) AppendMsgBytes(dst []byte, p []byte) []byte {
	return conv.EscStringBytes(dst, p, false, ' ')
}

func (f *formatterText) AppendKVInt(dst []byte, key string, val int) []byte {
	dst = conv.EscKey(dst, key, true, '=')
	return conv.Int(dst, val, false, ',')
}

func (f *formatterText) AppendKVString(dst []byte, key string, val string) []byte {
	dst = conv.EscKey(dst, key, true, '=')
	return conv.EscString(dst, val, true, ',')
}

func (f *formatterText) AppendKVFloat(dst []byte, key string, val float64) []byte {
	dst = conv.EscKey(dst, key, true, '=')
	return conv.Float(dst, val, false, ',')
}

func (f *formatterText) AppendKVBool(dst []byte, key string, val bool) []byte {
	dst = conv.EscKey(dst, key, true, '=')
	if val == true {
		return append(dst, `true,`...)
	}
	return append(dst, `false,`...)
}

func (f *formatterText) AppendKVError(dst []byte, key string, val error) []byte {
	dst = conv.EscKey(dst, key, true, '=')
	if val != nil {
		return conv.EscString(dst, val.Error(), true, ',')
	}
	return append(dst, `null,`...)
}

func (f *formatterText) AppendSuffix(dst []byte, suffix []byte) []byte {
	if suffix != nil {
		return append(dst, suffix...)
	}
	return dst
}

func (f *formatterText) Final(dst, suffix []byte) []byte {
	if len(dst) > 0 { // only do this if dst exists,
		dst = dst[:len(dst)-1] // trim last space
	}
	if suffix != nil {
		return append(dst, suffix...)
	}
	return dst
}
