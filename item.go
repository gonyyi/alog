package alog

import "time"

type alogItem struct {
	bufCap int
	buf    []byte // this is a buffer that will be created multiple and used by multiple goroutines by sync.Pool
}

func newItem(size int) *alogItem {
	return &alogItem{
		bufCap: size,
		buf:    make([]byte, size),
	}
}

func (i *alogItem) reset() {
	i.buf = i.buf[:0]
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
func (l *Logger) Warn(tag Tag, msg string, a ...interface{}) {
	l.Log(Lwarn, tag, msg, a...)
}
func (l *Logger) Error(tag Tag, msg string, a ...interface{}) {
	l.Log(Lerror, tag, msg, a...)
}
func (l *Logger) Fatal(tag Tag, msg string, a ...interface{}) {
	l.Log(Lfatal, tag, msg, a...)
}
