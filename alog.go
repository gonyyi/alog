package alog

import (
	"io"
	"sync"
)

func New(w io.Writer) *Logger {
	l := Logger{
		out:     toAlWriter(w),
		fmt:     FmtText, // FmtJSON,
		fmtFlag: Fdefault,
	}
	l.ctl.ctlLevel = Linfo
	l.buf.Init(nil, nil, 512, 2048) // TODO: clean up this
	return &l
}

type Logger struct {
	mu      sync.Mutex // check if i really don't need this IF i am using sync.Pool
	out     AlWriter
	buf     bufSyncPool
	fmt     Formatter
	ctl     control
	fmtFlag Format
}

func (l *Logger) Do(fns ...func(*Logger)) {
	for _, f := range fns {
		f(l)
	}
}
func (l *Logger) Output() AlWriter {
	return l.out
}
func (l *Logger) SetOutput(w io.Writer) *Logger {
	l.out = toAlWriter(w)
	return l
}
func (l *Logger) SetFormatter(fmt Formatter) *Logger {
	if fmt != nil {
		l.fmt = fmt
	}
	return l
}
func (l *Logger) Format() Format {
	return l.fmtFlag
}
func (l *Logger) SetFormat(f Format) *Logger {
	l.fmtFlag = f
	return l
}
func (l *Logger) GetTag(name string) (tag Tag, ok bool) {
	return l.ctl.Tags.GetTag(name)
}
func (l *Logger) MustGetTag(name string) (tag Tag) {
	return l.ctl.Tags.MustGetTag(name)
}
func (l *Logger) SetControlFn(fn func(Level, Tag) bool) *Logger {
	l.ctl.CtlFn(fn)
	return l
}
func (l *Logger) SetControl(lv Level, tag Tag) *Logger {
	l.ctl.CtlTag(lv, tag)
	return l
}
func (l *Logger) SetHook(h HookFn) *Logger {
	l.ctl.hook = h
	return l
}
func (l *Logger) Log(level Level, tag Tag, msg string, a ...interface{}) {
	if l.ctl.Check(level, tag) {
		buf := l.buf.Get()
		defer l.buf.Reset(buf)

		buf.Head = l.fmt.Start(buf.Head, nil)

		if l.fmtFlag&fUseTime != 0 {
			buf.Head = l.fmt.AppendTime(buf.Head, l.fmtFlag)
		}
		if l.fmtFlag&Flevel != 0 {
			buf.Head = l.fmt.AppendLevel(buf.Head, level)
		}
		if l.fmtFlag&Ftag != 0 {
			buf.Head = l.fmt.AppendTag(buf.Head, &l.ctl.Tags, tag)
		}
		buf.Body = l.fmt.AppendMsg(buf.Body, msg)

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
						buf.Body = l.fmt.AppendKVString(buf.Body, key, next.(string))
					case nil:
						buf.Body = l.fmt.AppendKVString(buf.Body, key, `nil`)
					case error:
						buf.Body = l.fmt.AppendKVString(buf.Body, key, next.(error).Error())
					case bool:
						buf.Body = l.fmt.AppendKVBool(buf.Body, key, next.(bool))
					case int:
						buf.Body = l.fmt.AppendKVInt(buf.Body, key, next.(int))
					case int64:
						buf.Body = l.fmt.AppendKVInt(buf.Body, key, int(next.(int64)))
					case uint:
						buf.Body = l.fmt.AppendKVInt(buf.Body, key, int(next.(uint)))
					case float32:
						buf.Body = l.fmt.AppendKVFloat(buf.Body, key, float64(next.(float32)))
					case float64:
						buf.Body = l.fmt.AppendKVFloat(buf.Body, key, next.(float64))
					default:
						buf.Body = l.fmt.AppendKVString(buf.Body, key, `unsupp??`)
					}
				} else {
					buf.Body = l.fmt.AppendKVString(buf.Body, key, `null`)
				}
			}
		}

		if l.ctl.hook != nil {
			l.ctl.hook(level, tag, buf.Body)
		}

		l.out.WriteTag(level, tag, buf.Head, l.fmt.Final(buf.Body, nil))
	}
}
func (l *Logger) Iferr(err error, tag Tag, msg string) bool {
	if err != nil {
		l.Log(Lerror, tag, msg, "error", err)
		return true
	}
	return false
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
