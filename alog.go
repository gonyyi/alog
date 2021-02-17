package alog

import (
	"io"
	"sync"
)

func New(w io.Writer) *Logger {
	l := Logger{
		out:   iowToAlw(w),
		fmtr:  &defaultFormatter{},
		fflag: Fdefault,
	}
	l.ctl.ctlLevel = Linfo
	l.buf.Init([]byte{'{'}, []byte{'}'}, 512, 2048)

	return &l
}

type Logger struct {
	mu    sync.Mutex // check if i really don't need this IF i am using sync.Pool
	out   AlWriter
	buf   bufSyncPool
	fmtr  Formatter
	ctl   control
	fflag Format
}

func (l *Logger) Output() io.Writer {
	return l.out
}
func (l *Logger) SetOutput(w io.Writer) *Logger {
	l.out = iowToAlw(w)
	return l
}
func (l *Logger) Format() Format {
	return l.fflag
}
func (l *Logger) SetFormat(f Format) *Logger {
	l.fflag = f
	return l
}
func (l *Logger) GetTag(name string) (tag Tag, ok bool) {
	return l.ctl.Tags.GetTag(name)
}
func (l *Logger) MustGetTag(name string) (tag Tag) {
	return l.ctl.Tags.MustGetTag(name)
}
func (l *Logger) Log(level Level, tag Tag, msg string, a ...interface{}) {
	if l.ctl.Check(level, tag) {
		buf := l.buf.Get()
		defer l.buf.Reset(buf)

		buf.Head = l.fmtr.AppendPrefix(buf.Head, nil)
		if l.fflag&fUseTime != 0 {
			buf.Head = l.fmtr.AppendTime(buf.Head, l.fflag)
		}
		if l.fflag&Flevel != 0 {
			buf.Head = l.fmtr.AppendLevel(buf.Head, level)
		}
		if l.fflag&Ftag != 0 {
			buf.Head = l.fmtr.AppendTag(buf.Head, &l.ctl.Tags, tag)
		}
		buf.Body = l.fmtr.AppendMsg(buf.Body, msg)

		if a != nil {
			lenA := len(a)
			idxA := lenA - 1
			for i := 0; i < lenA; i += 2 { // 0, 2, 4..
				key, ok := a[i].(string)
				if !ok {
					key = "badKey??"
				}
				if i < idxA {
					next := a[i+1]
					switch next.(type) {
					case string:
						buf.Body = l.fmtr.AppendKVString(buf.Body, key, next.(string))
					case nil:
						buf.Body = l.fmtr.AppendKVString(buf.Body, key, `nil`)
					case error:
						buf.Body = l.fmtr.AppendKVString(buf.Body, key, next.(error).Error())
					case bool:
						buf.Body = l.fmtr.AppendKVBool(buf.Body, key, next.(bool))
					case int:
						buf.Body = l.fmtr.AppendKVInt(buf.Body, key, next.(int))
					case int64:
						buf.Body = l.fmtr.AppendKVInt(buf.Body, key, int(next.(int64)))
					case uint:
						buf.Body = l.fmtr.AppendKVInt(buf.Body, key, int(next.(uint)))
					case float32:
						buf.Body = l.fmtr.AppendKVFloat(buf.Body, key, float64(next.(float32)))
					case float64:
						buf.Body = l.fmtr.AppendKVFloat(buf.Body, key, next.(float64))
					default:
						buf.Body = l.fmtr.AppendKVString(buf.Body, key, `unsupp??`)
					}
				} else {
					buf.Body = l.fmtr.AppendKVString(buf.Body, key, `null`)
				}
			}
		}

		// Replace extra comma
		buf.Body[len(buf.Body)-1] = '}'
		buf.Body = append(buf.Body, '\n')
		l.out.WriteTag(level, tag, buf.Head, buf.Body)
	}
}
func (l *Logger) Trace(tag Tag, msg string, a ...interface{}) {
	l.Log(Ltrace, tag, msg, a...)
}
func (l *Logger) Debug(tag Tag, msg string, a ...interface{}) {
	l.Log(Ldebug, tag, msg, a...)
}
func (l *Logger) Info(tag Tag, msg string, a ...interface{}) {
	l.Log(Linfo, tag, msg, a...)
}
func (l *Logger) Notice(tag Tag, msg string, a ...interface{}) {
	l.Log(Lnotice, tag, msg, a...)
}
func (l *Logger) Warn(tag Tag, msg string, a ...interface{}) {
	l.Log(Lwarn, tag, msg, a...)
}
func (l *Logger) Error(tag Tag, msg string, a ...interface{}) {
	l.Log(Lerror, tag, msg, a...)
}
func (l *Logger) Fatal(tag Tag, msg string, a ...interface{}) {
	l.Log(Lfatal, tag, msg, a...)
}
