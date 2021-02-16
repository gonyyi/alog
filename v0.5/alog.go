package alog

import (
	"io"
	"sync"
)

type Logger struct {
	mu  sync.Mutex
	out AlWriter
	buf Buffer
	f   Formatter
	// f      formatterJSON
	Ctrl   control
	format Format
}

func New(iow io.Writer, fns ...func(*Logger)) *Logger {
	l := Logger{
		out: iowToAlw(iow),
	}

	for _, f := range fns {
		f(&l)
	}
	if l.f == nil {
		fmtj := formatterJSON{}
		fmtj.Init()
		l.f = fmtj

	}
	if l.buf == nil {
		l.buf = &bufSyncPool{}
		l.buf.Init(256, 1024)
	}
	l.Ctrl.Filter.Set(Linfo, 0)
	return &l
}

// Output returns the output destination for the logger.
func (l *Logger) Output() io.Writer {
	return l.out
}

// SetOutput can redefined the output after logger has been created.
// If output is nil, the logger will Set it to ioutil.Discard instead.
func (l *Logger) SetOutput(output io.Writer) *Logger {
	l.mu.Lock()
	if output == nil {
		l.out = iowToAlw(discard)
	} else {
		l.out = iowToAlw(output)
	}
	l.mu.Unlock()
	return l
}

// SetHookFn will create a hookFn that works addition to Filter.
// Example would be log everything but for HTTP request tags,
// also write it to a file.
func (l *Logger) SetHookFn(fn HookFn) *Logger {
	l.Ctrl.SetHook(fn)
	return l
}

// GetTag takes a wTag name and returns a wTag if found.
func (l *Logger) GetTag(name string) Tag {
	return l.Ctrl.Tags.MustGetTag(name)
}

func (l *Logger) Log(lvl Level, tag Tag, msg string) {
	if l.Ctrl.Filter.Check(lvl, tag) {
		buf := l.buf.Get()
		buf.Head = l.f.AppendPrefix(buf.Head[:0], nil)

		if l.format&fUseTime != 0 {
			buf.Head = l.f.AppendTime(buf.Head, l.format)
		}
		if l.format&Ftag != 0 {
			buf.Head = l.f.AppendTag(buf.Head, &l.Ctrl.Tags, tag)
		}

		buf.Main = l.f.AppendMsg(buf.Main[:0], msg)
		buf.Main = l.f.AppendSuffix(buf.Main, nil)

		l.out.Write(append(buf.Head, buf.Main...))
		l.buf.Reset(buf)
	}
}

func (l *Logger) Loga(lvl Level, tag Tag, a ...interface{}) {
	if l.Ctrl.Filter.Check(lvl, tag) {
		buf := l.buf.Get()
		buf.Head = l.f.AppendPrefix(buf.Head[:0], nil)

		if l.format&fUseTime != 0 {
			buf.Head = l.f.AppendTime(buf.Head, l.format)
		}
		if l.format&Ftag != 0 {
			buf.Head = l.f.AppendTag(buf.Head, &l.Ctrl.Tags, tag)
		}
		buf.Main = l.f.AppendAdd(buf.Main[:0], a...)
		buf.Main = l.f.AppendSuffix(buf.Main, nil)

		l.out.Write(append(buf.Head, buf.Main...))
		l.buf.Reset(buf)
	}
}
