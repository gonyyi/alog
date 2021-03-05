package alog_test

import (
	"bytes"
	"github.com/gonyyi/alog"
	"github.com/gonyyi/alog/ext"
	"io"
	"testing"
)

type fakeWriter2 struct{}

func (fakeWriter2) Write([]byte) (int, error) {
	return 0, nil
}

type fakeWriter struct{}

func (fakeWriter) Write([]byte) (int, error) {
	return 0, nil
}
func (fakeWriter) Close() error {
	out.WriteString("closed\n")
	return nil
}

var out bytes.Buffer
var tag1, tag2 alog.Tag
var log = alog.New(&out).Do(func(l alog.Logger) alog.Logger {
	l.Flag = alog.UseTag | alog.UseLevel
	tag1 = l.NewTag("TAG1")
	tag2 = l.NewTag("TAG2")
	return l
})

func check(t *testing.T, exp string) {
	if exp != "" {
		exp += "\n"
	}

	if act := out.String(); act != exp {
		//t.Errorf("Expected: <%s>, Actual: <%s>", strconv.Quote(exp), strconv.Quote(act))
		t.Errorf("%s\nExpected: <%s>\nActual  : <%s>", t.Name(), exp, act)
	}
	out.Reset()
	log.Control.Fn = nil
	log.Control.Level = alog.InfoLevel
	log.Control.Tags = 0
	log = log.SetOutput(&out).SetFormatter(nil)
}

func TestLogger_Close(t *testing.T) {
	{
		log = log.SetOutput(&fakeWriter{})
		log.Close()
		check(t, `closed`)
	}

	{
		log = log.SetFormatter(ext.NewFormatterTerminal())
		log.Close()
		check(t, ``)
	}

	{
		fw := &fakeWriter2{}
		var tmp io.Writer
		tmp = fw

		log = log.SetOutput(tmp)
		log.Info(0).Write("test123")
		_, _ = log.Output().(*fakeWriter2)
		check(t, ``)
	}

	{
		log = log.SetOutput(nil)
		log.Info(0).Write("test")
		log.Close()
		check(t, ``)
	}

}

func TestLogger_Trace(t *testing.T) {
	log.Trace(0).Write("test1")
	check(t, ``)
	log.Control.Level = alog.TraceLevel
	log.Trace(0).Write("test2")
	check(t, `{"level":"trace","tag":[],"message":"test2"}`)
}

func TestLogger_Debug(t *testing.T) {
	log.Debug(0).Write("test")
	check(t, ``)
	log.Control.Level = alog.DebugLevel
	log.Debug(0).Write("test")
	check(t, `{"level":"debug","tag":[],"message":"test"}`)
}

func TestLogger_Info(t *testing.T) {
	log.Info(0).Write("test")
	check(t, `{"level":"info","tag":[],"message":"test"}`)
	log.Control.Level = alog.ErrorLevel
	log.Info(0).Write("test")
	check(t, ``)
}

func TestLogger_Warn(t *testing.T) {
	log.Warn(0).Write("test")
	check(t, `{"level":"warn","tag":[],"message":"test"}`)
	log.Control.Level = alog.ErrorLevel
	log.Warn(0).Write("test")
	check(t, ``)
}

func TestLogger_Error(t *testing.T) {
	log.Error(0).Write("test")
	check(t, `{"level":"error","tag":[],"message":"test"}`)
	log.Control.Level = alog.FatalLevel
	log.Error(0).Write("test")
	check(t, ``)
}
func TestLogger_Fatal(t *testing.T) {
	log.Fatal(0).Write("test")
	check(t, `{"level":"fatal","tag":[],"message":"test"}`)
	log.Control.Level = alog.FatalLevel
	log.Fatal(0).Write("test")
	check(t, `{"level":"fatal","tag":[],"message":"test"}`)
}

func TestLogger_SetFormatter(t *testing.T) {
	log = log.Do(nil).Do(ext.DoFmt.TXT())
	log.Info(0).Str("test", "ok").Write("done")
	check(t, `INF [] done // test="ok"`)

	tmp := log.SetFormatter(nil)
	tmp.Info(0).Str("test", "ok").Write("done")
	check(t, `{"level":"info","tag":[],"message":"done","test":"ok"}`)
}
func TestNew(t *testing.T) {
	log = alog.New(nil)
	log.Flag = alog.UseLevel | alog.UseTag
	log.Fatal(0).Write("error!")
	check(t, ``)

	log = alog.New(&out)
	log.Flag = alog.UseLevel | alog.UseTag
	log.Fatal(0).Write("error!")
	check(t, `{"level":"fatal","tag":[],"message":"error!"}`)
}
func TestLogger_getEntry(t *testing.T) {
	newFakeControlFn := func(retVal bool) alog.ControlFn {
		return func(level alog.Level, tag alog.Tag) bool {
			return retVal
		}
	}

	// ControlFn= YES --> false
	// Control 	= YES --> true
	{
		log.Control.Fn = newFakeControlFn(false) // always return false
		log.Fatal(0).Write("test")               // this shouldn't print
		check(t, "")
	}

	// ControlFn= YES --> false
	// Control 	= YES --> false
	{
		log.Control.Fn = newFakeControlFn(false) // always return false
		log.Trace(0).Write("test")               // this shouldn't print
		check(t, "")
	}

	// ControlFn= YES --> true
	// Control 	= YES --> true
	{
		log.Control.Fn = newFakeControlFn(true)
		log.Fatal(0).Write("test")
		check(t, `{"level":"fatal","tag":[],"message":"test"}`)
	}

	// ControlFn= YES --> true
	// Control 	= YES --> false
	{
		log.Control.Fn = newFakeControlFn(true)
		log.Trace(0).Write("test")
		check(t, `{"level":"trace","tag":[],"message":"test"}`)
	}

	// ControlFn= NO
	// Control 	= true
	{
		log.Control.Fn = nil
		log.Info(0).Write("test")
		check(t, `{"level":"info","tag":[],"message":"test"}`)
	}

	// ControlFn= NO
	// Control 	= false
	{
		log.Control.Fn = nil
		log.Trace(0).Write("test")
		check(t, ``)
	}
}
