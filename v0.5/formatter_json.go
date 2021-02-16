package alog

import "time"

type formatterJSON struct {
	conv FormatterConverter
	esc  FormatterEscaper
	t    time.Time
}

func (f *formatterJSON) Init() {
	f.conv = &formatterConvBasic{}
	f.esc = &formatterEscBasic{}
	f.esc.Init()
}

func (f formatterJSON) AppendPrefix(dst []byte, prefix []byte) []byte {
	if prefix != nil {
		return append(dst, prefix...)
	}
	return append(dst, '{')
}

func (f formatterJSON) AppendTime(dst []byte, format Format) []byte {
	if fUseTime&format != 0 {
		f.t = time.Now()
		if FtimeUnixMs&format != 0 {
			dst = append(dst, `"ts":`...) // faster without addKey
			return f.conv.Intf(dst, int(f.t.UnixNano())/1e6, 0, ',')
		} else if FtimeUnix&format != 0 {
			dst = append(dst, `"ts":`...) // faster without addKey
			return f.conv.Intf(dst, int(f.t.Unix()), 0, ',')
		} else {
			if FtimeUTC&format != 0 {
				f.t = f.t.UTC()
			}
			if Fdate&format != 0 {
				dst = append(dst, `"d":`...) // faster without addKey
				y, m, d := f.t.Date()
				dst = f.conv.Intf(dst, y*10000+int(m)*100+d, 4, ',')
			}
			if FdateDay&format != 0 {
				// "wd": 0 being sunday, 6 being saturday
				dst = append(dst, `"wd":`...) // faster without addKey
				dst = f.conv.Intf(dst, int(f.t.Weekday()), 1, ',')
			}
			if Ftime&format != 0 {
				dst = append(dst, `"t":`...) // faster without addKey
				h, m, s := f.t.Clock()
				dst = f.conv.Intf(dst, h*10000+m*100+s, 1, '.')
				dst = f.conv.Intf(dst, f.t.Nanosecond()/1e6, 3, ',')
			}
		}
	}
	return dst
}

func (f formatterJSON) AppendTag(dst []byte, tb *TagBucket, tag Tag) []byte {
	dst = append(dst, `"tag":[`...)
	if tag != 0 {
		dst = tb.AppendSelectedTags(dst, ',', true, tag)
	}
	dst = append(dst, ']', ',')
	return dst
}

func (f formatterJSON) AppendMsg(dst []byte, s string) []byte {
	if len(s) != 0 {
		dst = append(dst, `"msg":`...)
		return f.esc.Val(dst, s, true, ',')
	}
	return dst
}

func (f formatterJSON) AppendMsgBytes(dst []byte, p []byte) []byte {
	if p != nil {
		dst = append(dst, `"msg":`...)
		return f.esc.ValBytes(dst, p, true, ',')
	}
	return dst
}

func (f formatterJSON) AppendAdd(dst []byte, a ...interface{}) []byte {
	lenA := len(a)
	idxA := lenA - 1
	for i := 0; i < lenA; i += 2 { // 0, 2, 4..
		if key, ok := a[i].(string); !ok {
			key = "badKey??"
		} else {
			dst = f.esc.Key(dst, key, true, ':') // key --> "key":
		}

		if i < idxA {
			next := a[i+1]
			switch next.(type) {
			case string:
				dst = f.esc.Val(dst, next.(string), true, ',')
			case nil:
				dst = append(dst, `null,`...)
			case error:
				dst = f.conv.Error(dst, next.(error), true, ',')
			case bool:
				dst = f.conv.Bool(dst, next.(bool), false, ',')
			case int:
				dst = f.conv.Int(dst, next.(int), false, ',')
			case int32:
				dst = f.conv.Int(dst, int(next.(int32)), false, ',')
			case int64:
				dst = f.conv.Int(dst, int(next.(int64)), false, ',')
			case uint:
				dst = f.conv.Int(dst, int(next.(uint)), false, ',')
			case uint32:
				dst = f.conv.Int(dst, int(next.(uint32)), false, ',')
			case uint64:
				dst = f.conv.Int(dst, int(next.(uint64)), false, ',')
			case float32:
				dst = f.conv.Float(dst, float64(next.(float32)), false, ',')
			case float64:
				dst = f.conv.Float(dst, next.(float64), false, ',')
			case *[]string:
				dst = append(dst, '[')
				for _, v := range *next.(*[]string) {
					dst = f.esc.Val(dst, v, true, ',')
				}
				dst = f.TrimLast(dst, ',')
				dst = append(dst, ']', ',')
			case *[]int:
				dst = append(dst, '[')
				for _, v := range *next.(*[]int) {
					dst = f.conv.Int(dst, v, false, ',')
				}
				dst = f.TrimLast(dst, ',')
				dst = append(dst, ']', ',')
			case *[]float64:
				dst = append(dst, '[')
				for _, v := range *next.(*[]float64) {
					dst = f.conv.Float(dst, v, false, ',')
				}
				dst = f.TrimLast(dst, ',')
				dst = append(dst, ']', ',')
			case *[]bool:
				dst = append(dst, '[')
				for _, v := range *next.(*[]bool) {
					dst = f.conv.Bool(dst, v, false, ',')
				}
				dst = f.TrimLast(dst, ',')
				dst = append(dst, ']', ',')
			case *[]error:
				dst = append(dst, '[')
				for _, v := range *next.(*[]error) {
					dst = f.conv.Error(dst, v, false, ',')
				}
				dst = f.TrimLast(dst, ',')
				dst = append(dst, ']', ',')
			default:
				dst = append(dst, `null,`...)
			}
		} else {
			dst = append(dst, `null,`...)
		}
	}
	return dst
}

func (f formatterJSON) AppendSuffix(dst []byte, suffix []byte) []byte {
	dst = f.TrimLast(dst, ',')
	if suffix != nil {
		return append(dst, suffix...)
	}
	return append(dst, '}', '\n')
}

func (f formatterJSON) TrimLast(dst []byte, b byte) []byte {
	if len(dst) > 0 && dst[len(dst)-1] == b {
		return dst[:len(dst)-1]
	}
	return dst
}
