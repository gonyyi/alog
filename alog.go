package alog

import (
	"io"
)

// New will return a Alog logger pointer with default values.
// This function will take an io.Writer and convert it to AlWriter.
// A user's custom AlWriter will let the user steer more control.
// (eg. per tag, or per level)
// Default formatter: 		FmtText
// Default format flag:		date, time, level, tag
// Default logging level: 	Info
// Default buf size: 	512, 2048
func New(w io.Writer) *Logger {
	l := Logger{
		out:     newAlWriter(w),
		buf:     Conf.Buffer(),
		fmtr:    Conf.Formatter(),
		fmtFlag: Conf.FormatFlag,
	}

	l.ctl.CtlTag(Conf.ControlLevel, 0)
	l.buf.Init(Conf.BufferHead, Conf.BufferBody)
	return &l
}

// Logger is a main struct for Alog.
type Logger struct {
	out       AlWriter
	buf       Buffer
	fmt       formatJSON
	fmtr      Formatter
	fmtPrefix []byte
	fmtSuffix []byte
	ctl       control
	fmtFlag   Format
}

// Do will run (series of) function(s) and is used for
// quick macro like settings for the logger.
func (l *Logger) Do(fns ...func(*Logger)) {
	for _, f := range fns {
		f(l)
	}
}

// Close will close io.Writer if applicable
func (l *Logger) Close() error {
	return l.out.Close()
}

// SetOutput will set the output writer to be used
// in the logger. If nil is given, it will discard the output.
func (l *Logger) SetOutput(w io.Writer) *Logger {
	l.out = newAlWriter(w)
	return l
}

// SetFormatter will take Formatter compatible objects.
// This will only reset the formatter if it's not nil.
// In case nil is given as its argument, it will ignore.
func (l *Logger) SetFormatter(fmt Formatter) *Logger {
	if fmt != nil {
		l.fmtr = fmt
	} else {
		l.fmtr = Conf.Formatter()
	}
	return l
}

// SetFormat will set the format flag.
func (l *Logger) SetFormat(f Format) *Logger {
	l.fmtFlag = f
	return l
}

// SetBufferSize will initialize head and body buffers.
func (l *Logger) SetBufferSize(head, body int) *Logger {
	l.buf.Init(head, body)
	return l
}

// SetAffix will set both prefix and suffix. If only is not to be set,
// use nil. Eg. SetAffix(nil, []byte("--end"))
func (l *Logger) SetAffix(prefix, suffix []byte) *Logger {
	l.fmtPrefix, l.fmtSuffix = prefix, suffix
	return l
}

// GetTag will take a name of tag and return it. If not found, it will
// return false for the output ok.
func (l *Logger) GetTag(name string) (tag Tag, ok bool) {
	return l.ctl.Tags.GetTag(name)
}

// MustGetTag will return a tag. If a required tag is not exists,
// it will create one.
func (l *Logger) MustGetTag(name string) (tag Tag) {
	return l.ctl.Tags.MustGetTag(name)
}

// SetControlFn will set a ControlFn that determines what to log.
// By using this instead of SetControl, a user can control precisely.
func (l *Logger) SetControlFn(fn ControlFn) *Logger {
	l.ctl.CtlFn(fn)
	return l
}

// SetControl will set logging level and tag.
// Note that this is an OR condition: if level has met the minimum logging level OR
// tag is met, the logger will log. For any precise control, use SetControlFn.
func (l *Logger) SetControl(lv Level, tag Tag) *Logger {
	l.ctl.CtlTag(lv, tag)
	return l
}

// SetHook will run HookFn if set. This can be used to special custom situation.
// As HookFn will run AFTER right before formatter's method Final is being called,
// its argument p []byte will have already formatted body.
func (l *Logger) SetHook(h HookFn) *Logger {
	l.ctl.SetHook(h)
	return l
}

// Log is a main method of logging. This takes level, tag, message, as well as optional
// data. The optional data has to be in pairs of name and value. For the speed, Alog only
// supports basic types: int, int64, uint, string, bool, float32, float64.
func (l *Logger) Log(level Level, tag Tag, msg string, a ...interface{}) (n int, err error) {
	// Below recover may not needed but worst possible case..
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case error:
				err = r.(error)
				println("alog recovered: panic(" + err.Error() + ")")
			case string:
				println("alog recovered: panic(" + r.(string) + ")")
			}
		}
	}()

	if l.ctl.Check(level, tag) && l.fmtr != nil && l.out != nil {
		buf := l.buf.Get()
		defer l.buf.Reset(buf)

		buf.Head = l.logHead(level, tag, buf.Head)
		buf.Body = l.fmtr.AppendMsg(buf.Body, msg)

		if a != nil {
			buf.Body = l.fmtr.AppendSeparator(buf.Body)
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
		buf.Body = l.logTail(level, tag, buf.Body)
		return l.out.WriteTag(level, tag, buf.Head, buf.Body)
	}
	return 0, nil
}

// logb will log a simple msg but takes byte slice instead of string
func (l *Logger) logb(level Level, tag Tag, msgb []byte) (n int, err error) {
	if l.ctl.Check(level, tag) && l.fmtr != nil && l.out != nil {
		buf := l.buf.Get()
		defer l.buf.Reset(buf)
		buf.Head = l.logHead(level, tag, buf.Head)
		buf.Body = l.fmtr.AppendMsgBytes(buf.Body, msgb)
		buf.Body = l.logTail(level, tag, buf.Body)
		return l.out.WriteTag(level, tag, buf.Head, buf.Body)
	}
	return 0, nil
}

// logHead creates head part of the log message.
func (l *Logger) logHead(level Level, tag Tag, dst []byte) []byte {
	if l.fmtFlag&Fprefix != 0 {
		dst = l.fmtr.Start(dst, l.fmtPrefix)
	} else {
		dst = l.fmtr.Start(dst, nil)
	}
	// If any time components are in the format flag, then append the time
	if l.fmtFlag&fUseTime != 0 {
		dst = l.fmtr.AppendTime(dst, l.fmtFlag)
	}
	if l.fmtFlag&Flevel != 0 {
		dst = l.fmtr.AppendLevel(dst, level)
	}
	if l.fmtFlag&Ftag != 0 {
		dst = l.fmtr.AppendTag(dst, &l.ctl.Tags, tag)
	}
	return dst
}

// logTail creates tail part of the log message
func (l *Logger) logTail(level Level, tag Tag, dst []byte) []byte {
	// Run control hook func if any.
	if l.ctl.hook != nil {
		l.ctl.hook(level, tag, dst)
	}
	// Check prefix flag, if exist run.
	if l.fmtFlag&Fsuffix != 0 {
		return l.fmtr.Final(dst, l.fmtSuffix)
	} else {
		return l.fmtr.Final(dst, nil)
	}
}

// Iferr method will log an error when argument err is not nil.
// This also returns true/false if error is or not nil.
func (l *Logger) Iferr(err error, tag Tag, msg string) bool {
	if err != nil {
		l.Log(Lerror, tag, msg, "error", err)
		return true
	}
	return false
}

// Trace records a msg with a trace level with optional additional variables
func (l *Logger) Trace(tag Tag, msg string, a ...interface{}) {
	l.Log(Ltrace, tag, msg, a...)
}

// Debug records a msg with a debug level with optional additional variables
func (l *Logger) Debug(tag Tag, msg string, a ...interface{}) {
	l.Log(Ldebug, tag, msg, a...)
}

// Info records a msg with an info level with optional additional variables
// And info level is default log level of Alog.
func (l *Logger) Info(tag Tag, msg string, a ...interface{}) {
	l.Log(Linfo, tag, msg, a...)
}

// Warn records a msg with a warning level with optional additional variables
func (l *Logger) Warn(tag Tag, msg string, a ...interface{}) {
	l.Log(Lwarn, tag, msg, a...)
}

// Error records a msg with an error level with optional additional variables
func (l *Logger) Error(tag Tag, msg string, a ...interface{}) {
	l.Log(Lerror, tag, msg, a...)
}

// Fatal records a msg with a fatal level with optional additional variables.
// Unlike other logger, Alog will NOT terminal the program with a Fatal method.
// A user need to handle what to do.
func (l *Logger) Fatal(tag Tag, msg string, a ...interface{}) {
	l.Log(Lfatal, tag, msg, a...)
}

// NewSubWriter will create a SubWriter that can be used by other
// library to write a log message using Alog. SubWriter meets io.Writer
// interface format.
func (l *Logger) NewSubWriter(dLevel Level, dTag Tag) SubWriter {
	return SubWriter{
		l:      l,
		dLevel: dLevel,
		dTag:   dTag,
	}
}
