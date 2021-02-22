package alog_test

import (
	"bytes"
	"errors"
	"github.com/gonyyi/alog"
	"os"
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	var out, outHook bytes.Buffer
	{
		alog.Conf.BufferHead = 0
		alog.Conf.BufferBody = 0
		alog.Conf.SetBuffer(nil)
		alog.Conf.SetConverter(nil)
		alog.Conf.Converter()
		alog.Conf.SetFormatter(nil)
	}
	al := alog.New(nil)
	al.Info(0, "no tag")

	IO := al.MustGetTag("IO")
	SYS := al.MustGetTag("SYS")
	IO1, _ := al.GetTag("IO1")
	IO2, _ := al.GetTag("IO2")
	al.Do(func(l *alog.Logger) {
		al.SetOutput(&out)
		al.SetFormat(
			alog.Fdefault.On(alog.Fprefix|alog.Fsuffix).Off(alog.Ftime) |
				alog.FdateDay | alog.FtimeUTC | alog.FtimeMs | alog.Fprefix | alog.Fsuffix)
		al.SetFormatter(nil)
		al.SetBufferSize(512, 1024)
	})

	al.SetFormatter(alog.Conf.Formatter())

	// These two should print as it's logging level is below default (info)
	al.Trace(0, "alog trace", "traceStr", "hello")
	al.Debug(IO, "alog debug", "debugInt", 123)

	// By setting logging level to debug, only debug will print below
	al.SetControl(alog.Ldebug, 0)
	al.Trace(0, "alog trace", "traceStr", "hello")
	al.Debug(IO, "alog debug", "debugInt", 123)

	// By setting level to fatal, and tag to IO,
	// every logs with a tag IO and logs with fatal will be printed.
	al.SetControl(alog.Lfatal, IO)
	al.Debug(0, "alog debug", "debugInt", 123)
	al.Debug(IO, "alog debug", "debugInt", 123) // print
	al.Fatal(0, "alog debug", "debugInt", 123)  // print

	al.SetControl(alog.Ldebug, IO)
	al.Trace(0, "alog trace", "traceStr", "hello")
	al.Debug(IO, "alog debug", "debugInt", 123)
	al.SetControlFn(func(level alog.Level, tag alog.Tag) bool {
		if tag != 0 {
			return true
		}
		return false
	})
	al.Info(SYS, "alog info", "infoBool", true, "infoBool2", false)
	al.Warn(IO|SYS, "alog warn", "warnError", errors.New("theError"))

	al.SetFormat(alog.FtimeUnix)
	al.Error(IO1, "alog error", "errFloat", 3.1415)
	al.SetFormat(alog.FtimeUnixMs)
	al.Fatal(IO2, "alog fatal", "fatalNothing", nil, "string", "string")

	al.SetAffix(
		[]byte("prefix"), []byte("suffix"))
	al.Info(SYS, "alog info", "infoBool", true, 123, false, "i64", int64(123), "u", uint(123), "f32", float32(3.14), "f64", float64(3.14))
	al.Warn(SYS, "test", "string", "string", "nil", nil, "str", struct{ a string }{}, "itself")
	al.SetHook(func(lvl alog.Level, tag alog.Tag, p []byte) {
		if lvl == alog.Lwarn && tag == IO {
			outHook.WriteString("found it")
		}
	})
	al.Warn(IO, "yay")

	sw := al.NewSubWriter(alog.Ldebug, SYS)
	sw.Write([]byte("swByte" + "a\t\b\n\r\f\\\""))
	al.SetAffix(nil, nil)
	sw.Write([]byte("swByte" + "a\t\b\n\r\f\\\""))
	sw.Trace("swTrace")
	sw.Debug("swDebug")
	sw.Info("swInfo")
	sw.Warn("swWarn")
	sw.Error("swErr")
	sw.Fatal("swFatal")
	sw.Write(nil)

	{
		al.SetControl(alog.Linfo, 0)
		sw := al.NewSubWriter(alog.Ltrace, SYS)
		sw.Write([]byte("123"))
	}
	{
		al.Iferr(nil, 0, "test")
		al.Iferr(errors.New("my err"), 0, "test")
	}
	{
		_, err := os.Create("./test.txt")
		al.Iferr(err, 0, "shit1")
		t, err := os.Open("./test.txt")
		al.Iferr(err, 0, "shit2")
		al.SetOutput(t)
		al.Info(0, "bufio writer")
		al.Close()
	}
	{
		t := func(level alog.Level) {
			level.String()
			level.ShortName()
		}
		t(alog.Ltrace)
		t(alog.Ldebug)
		t(alog.Linfo)
		t(alog.Lwarn)
		t(alog.Lerror)
		t(alog.Lfatal)
		t(99)
	}

	{
		(IO | SYS).Has(IO)
		(IO | SYS).Has(IO1)
		al.MustGetTag("abc")
		al.GetTag("abc")
		al.MustGetTag("abc")
		for i := 0; i < 63; i++ {
			al.MustGetTag("test:" + strconv.Itoa(i))
		}
	}

	{
		tb := alog.TagBucket{}
		var out []byte
		tb.AppendSelectedTags(out, 0, true, 0)
		tb.AppendSelectedTags(out, 0, false, 0)
		tb.AppendSelectedTags(out, 0, true, 99)
		tb.AppendSelectedTags(out, 0, false, 99)
		tb.AppendSelectedTags(nil, 0, true, 99)
		tb.AppendSelectedTags(nil, 0, false, 99)

		abc := tb.MustGetTag("abc")
		tb.AppendSelectedTags(out, 0, false, abc)
		tb.AppendSelectedTags(out, 0, false, 0)
		tb.AppendSelectedTags(out, 0, false, 99)
		tb.AppendSelectedTags(out, ' ', false, 99)
		tb.AppendSelectedTags(out, ' ', true, 99)
		al.Warn(0, "g", "ab", "a\t\b\n\r\f\\\"")
	}

	{
		var out []byte
		c := alog.Conf.Converter()
		c.EscKeyBytes(nil, []byte("test"), false, 0)
		c.EscKeyBytes(nil, []byte("test"), true, 0)
		c.EscKeyBytes(nil, []byte("test"), false, ' ')
		c.EscKeyBytes(nil, []byte("test"), true, ' ')
		c.EscKeyBytes(nil, nil, true, ' ')
		c.EscKeyBytes(out, nil, true, ' ')
		c.EscKey(nil, "", false, 0)
		c.EscKey(out, "", true, 0)
		c.EscKey(nil, "", false, ' ')
		c.EscKey(out, "", true, ' ')
		c.EscString(nil, " ", false, 0)
		c.EscString(nil, " ", true, 0)
		c.EscString(nil, " ", true, ' ')
		c.EscString(nil, " ", false, ' ')

		c.EscStringBytes(nil, []byte(" "), false, 0)
		c.EscStringBytes(nil, []byte(" "), true, 0)
		c.EscStringBytes(nil, []byte(" "), true, ' ')
		c.EscStringBytes(nil, []byte(" "), false, ' ')

		c.Int(nil, 0, true, ' ')
		c.Int(nil, 0, false, ' ')
		c.Int(nil, 0, true, 0)
		c.Int(nil, 0, false, 0)

		c.Intf(nil, 0, 0, ' ')
		c.Intf(nil, 0, -1, ' ')
		c.Intf(nil, 0, 3, ' ')
		c.Intf(nil, 0, 0, 0)
		c.Intf(nil, 0, -1, 0)
		c.Intf(nil, 0, 3, 0)
		c.Intf(nil, -3, 0, 0)
		c.Intf(nil, -3, -1, 0)
		c.Intf(nil, -3, 3, 0)

		c.Float(nil, 1.0, true, 0)
		c.Float(nil, 1.0, true, ' ')
		c.Float(nil, -1.0, true, 0)
		c.Float(nil, -1.0, true, ' ')
		c.Float(nil, -0.01, true, 0)
		c.Float(nil, -0.01, true, ' ')
		c.Float(nil, -0.15, true, 0)
		c.Float(nil, -0.14, true, ' ')

		c.Floatf(nil, 0.1555, 3, 0)
		c.Floatf(nil, 0.1444, 3, ' ')
		c.Floatf(nil, 0.1555, 0, 0)
		c.Floatf(nil, 0.1444, 0, ' ')

		c.Bool(nil, true, true, ' ')
		c.Bool(nil, true, false, ' ')
		c.Bool(nil, false, true, 0)
		c.Bool(nil, false, false, 0)

		e := errors.New("he")
		c.Error(nil, nil, true, 0)
		c.Error(nil, nil, false, 0)
		c.Error(nil, nil, true, ' ')
		c.Error(nil, nil, false, ' ')
		c.Error(nil, e, true, 0)
		c.Error(nil, e, false, 0)
		c.Error(nil, e, true, ' ')
		c.Error(nil, e, false, ' ')
	}
	al.Close()

	{
		fw := fakeWriter{}
		al.SetOutput(&fw)
		al.Info(0, "fake writer test")
		al.Close()
	}
	println(out.String())
}

type fakeWriter struct{}

func (w *fakeWriter) Write(p []byte) (int, error) {
	return 0, nil
}
