package alog_test

import (
	"bytes"
	"errors"
	"github.com/gonyyi/alog/v0.5"
	"io"
	"testing"
)

type testWriter struct {
	w0 io.Writer
	w1 io.Writer
}

func (w testWriter) Write(p []byte) (int, error) {
	return w.w0.Write(p)
}
func (w testWriter) WriteLvt(lv alog.Level, tag alog.Tag, p []byte) (int, error) {
	if lv == alog.Lerror && tag == 1 {
		return w.w0.Write(p)
	} else {
		return w.w1.Write(p)
	}
}

func TestDefaults_ToAlWriter(t *testing.T) {
	t.Run("usingStdWriter", func(t2 *testing.T) {
		var out bytes.Buffer
		alw := alog.Defaults.ToAlWriter(&out)
		alw.WriteLvt(0, 0, []byte("hello1"))
		if s := out.String(); s != `hello1` {
			t2.Error(s)
		}
	})

	t.Run("usingCustom", func(t2 *testing.T) {
		var out1, out2 bytes.Buffer
		var wt io.Writer = testWriter{w0: &out1, w1: &out2}
		alw := alog.Defaults.ToAlWriter(wt)
		alw.WriteLvt(alog.Lerror, 2, []byte("hello1")) // Print to out2
		if out1.String() != `` || out2.String() != `hello1` {
			t2.Error("out1:", out1.String())
			t2.Error("out2:", out2.String())
		}

		out1.Reset()
		out2.Reset()
		alw.WriteLvt(alog.Lerror, 1, []byte("hello2")) // Print to out1
		if out1.String() != `hello2` || out2.String() != `` {
			t2.Error("out1:", out1.String())
			t2.Error("out2:", out2.String())
		}
	})
}

func TestDefaults_FormatterJSON(t *testing.T) {
	fj := alog.Defaults.FormatterJSON()
	tb := alog.TagBucket{}
	tTest1 := tb.MustGetTag("test1")
	tTest2 := tb.MustGetTag("test2")

	var head, main []byte
	prefix, suffix := []byte{'{'}, []byte{'}'}
	head = fj.AppendPrefix(head, prefix)
	head = fj.AppendTime(head, alog.Ftime|alog.FtimeUTC|alog.FdateDay|alog.Fdate)
	head = fj.AppendTag(head, &tb, tTest1|tTest2)
	main = fj.AppendMsg(main, "Hello this is a test")
	main = fj.AppendAdd(main, "name", "gon yi", "age", 17, "test", 3.14)
	main = fj.TrimLast(main, ',')
	main = fj.AppendSuffix(main, suffix)
	println(string(head) + string(main))
}
func BenchmarkDefaults_FormatterJSON(b *testing.B) {
	fj := alog.Defaults.FormatterJSON()
	tb := alog.TagBucket{}
	tTest1 := tb.MustGetTag("test1")
	tTest2 := tb.MustGetTag("test2")
	_, _ = tTest1, tTest2

	var head, main []byte

	testStrArr := []string{"abcd1", "abcd2"}
	testIntArr := []int{1, 3, 5, 7, 9}
	testFloatArr := []float64{3.14, 3.14500}
	_ = testStrArr

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		head = fj.AppendPrefix(head[:0], nil)
		head = fj.AppendTime(head, alog.Ftime|alog.FtimeUTC|alog.FdateDay|alog.Fdate|alog.FtimeUnix)
		head = fj.AppendTag(head, &tb, tTest1|tTest2)
		main = fj.AppendMsg(main[:0], "Hello this is a test")
		main = fj.AppendAdd(main, "addr", &testStrArr)
		main = fj.AppendAdd(main, "num", &testIntArr)
		main = fj.AppendAdd(main, "float", &testFloatArr)
		main = fj.TrimLast(main, ',')
		main = fj.AppendSuffix(main, nil)
	}

	println(string(append(head, main...)))
}

func TestDefaults_BufferSingle(t *testing.T) {
	var buf alog.Buffer = alog.Defaults.BufferSingle()
	buf.Init(1024, 1024)
	{
		bb := buf.Get()
		bb.Head = append(bb.Head, "hello"...)
		if string(bb.Head) != "hello" {
			buf.Reset(bb)
			t.Fail()
		}
		buf.Reset(bb)
	}
	{
		bb := buf.Get()
		if string(bb.Head) != "" {
			t.Fail()
		}
	}
}

func TestDefaults_BufferSyncPool(t *testing.T) {
	var buf alog.Buffer = alog.Defaults.BufferSyncPool()
	buf.Init(1024, 1024)
	{
		bb := buf.Get()
		bb.Head = append(bb.Head, "hello"...)
		if string(bb.Head) != "hello" {
			buf.Reset(bb)
			t.Fail()
		}
	}
	{
		bb := buf.Get()
		if string(bb.Head) != "" {
			t.Fail()
		}
	}
}

func TestDefaults_ConverterSimple(t *testing.T) {
	var c alog.FormatterConverter = alog.Defaults.ConverterSimple()

	var buf []byte

	t.Run("Bool", func(t2 *testing.T) {
		buf = c.Bool(buf[:0], true, false, 0)
		if string(buf) != "true" {
			t2.Fail()
		}

		buf = c.Bool(buf[:0], false, false, 0)
		if string(buf) != "false" {
			t2.Fail()
		}

		buf = c.Bool(buf[:0], false, false, ',')
		if string(buf) != `false,` {
			t2.Fail()
		}

		buf = c.Bool(buf[:0], false, true, 0)
		if string(buf) != `"false"` {
			t2.Fail()
		}

		buf = c.Bool(buf[:0], false, true, ',')
		if string(buf) != `"false",` {
			t2.Fail()
		}
	})

	t.Run("Int", func(t2 *testing.T) {
		buf = c.Int(buf[:0], 123, false, 0)
		if string(buf) != `123` {
			t2.Fail()
		}

		buf = c.Int(buf[:0], 123, false, ',')
		if string(buf) != `123,` {
			t2.Fail()
		}
		buf = c.Int(buf[:0], 123, true, 0)
		if string(buf) != `"123"` {
			t2.Fail()
		}

		buf = c.Int(buf[:0], 123, true, ',')
		if string(buf) != `"123",` {
			t2.Fail()
		}
	})

	t.Run("Float", func(t2 *testing.T) {
		buf = c.Float(buf[:0], 123.456, false, 0)
		if string(buf) != `123.46` {
			t2.Fail()
		}

		buf = c.Float(buf[:0], 123.456, false, ',')
		if string(buf) != `123.46,` {
			t2.Fail()
		}
		buf = c.Float(buf[:0], 123.456, true, 0)
		if string(buf) != `"123.46"` {
			t2.Fail()
		}

		buf = c.Float(buf[:0], 123.456, true, ',')
		if string(buf) != `"123.46",` {
			t2.Fail()
		}
	})

	t.Run("Error", func(t2 *testing.T) {
		buf = c.Error(buf[:0], nil, false, 0)
		if string(buf) != "null" {
			t2.Error(string(buf))
			t2.Fail()
		}
		buf = c.Error(buf[:0], nil, false, ',')
		if string(buf) != `null,` {
			t2.Error(string(buf))
			t2.Fail()
		}
		buf = c.Error(buf[:0], nil, true, 0)
		if string(buf) != `"null"` {
			t2.Error(string(buf))
			t2.Fail()
		}
		buf = c.Error(buf[:0], nil, true, ',')
		if string(buf) != `"null",` {
			t2.Error(string(buf))
			t2.Fail()
		}

		buf = c.Error(buf[:0], errors.New("test error"), false, 0)
		if string(buf) != `test error` {
			t2.Error(string(buf))
			t2.Fail()
		}
		buf = c.Error(buf[:0], errors.New("test error"), false, ',')
		if string(buf) != `test error,` {
			t2.Error(string(buf))
			t2.Fail()
		}
		buf = c.Error(buf[:0], errors.New("test error"), true, 0)
		if string(buf) != `"test error"` {
			t2.Error(string(buf))
			t2.Fail()
		}
		buf = c.Error(buf[:0], errors.New("test error"), true, ',')
		if string(buf) != `"test error",` {
			t2.Error(string(buf))
			t2.Fail()
		}

	})

	t.Run("Intf", func(t2 *testing.T) {
		buf = c.Intf(buf[:0], 123, 0, 0)
		if string(buf) != `123` {
			t2.Error(string(buf))
			t2.Fail()
		}
		buf = c.Intf(buf[:0], 123, 5, 0)
		if string(buf) != `00123` {
			t2.Error(string(buf))
			t2.Fail()
		}
		buf = c.Intf(buf[:0], 123, 0, ',')
		if string(buf) != `123,` {
			t2.Error(string(buf))
			t2.Fail()
		}
		buf = c.Intf(buf[:0], 123, 5, ',')
		if string(buf) != `00123,` {
			t2.Error(string(buf))
			t2.Fail()
		}
	})

	t.Run("Floatf", func(t2 *testing.T) {
		buf = c.Floatf(buf[:0], 123.456, 0, 0)
		if string(buf) != `123` {
			t2.Error(string(buf))
			t2.Fail()
		}
		buf = c.Floatf(buf[:0], 123.446, 1, 0)
		if string(buf) != `123.4` {
			t2.Error(string(buf))
			t2.Fail()
		}
		buf = c.Floatf(buf[:0], 123.456, 1, 0)
		if string(buf) != `123.5` {
			t2.Error(string(buf))
			t2.Fail()
		}
		buf = c.Floatf(buf[:0], 123.456, 2, ',')
		if string(buf) != `123.46,` {
			t2.Error(string(buf))
			t2.Fail()
		}
	})
}
