package alog

import (
	"os"
	"testing"
	"time"
)

func Test_formatd(t *testing.T) {
	f := formatd{}
	f.init()

	var out []byte
	out = f.addBegin(out)
	out = f.addKey(out, "age")
	out = f.addValInt(out, 39)
	out = f.addKey(out, "name")
	out = f.addValString(out, "Gon Yi")
	out = f.addKey(out, "isMarried")
	out = f.addValBool(out, false)
	out = f.addKey(out, "height")
	out = f.addValFloat(out, 5.8)
	out = f.addEnd(out)

	println(string(out))
}
func Benchmark_new(b *testing.B) {
	b.ReportAllocs()
	al := New(nil)
	al.Format = Flevel
	TEST := al.Control.Tags.MustGetTag("TEST")
	IO := al.Control.Tags.MustGetTag("IO")
	_, _, _ = al, TEST, IO
	for i := 0; i < b.N; i++ {
		al.Log(Linfo, 0, "hello", Kstr("name", `Gon Yi		ha`))
	}

	al.SetOutput(os.Stderr)
	al.Log(Linfo, 0, "hello", Kstr("name", "Gon Yi"))
}
func Benchmark_formatd(b *testing.B) {

	f := formatd{}
	f.init()
	var out []byte
	b.ReportAllocs()

	al := New(nil)
	TEST := al.Control.Tags.MustGetTag("TEST")
	IO := al.Control.Tags.MustGetTag("IO")
	_, _, _ = al, TEST, IO

	for i := 0; i < b.N; i++ {
		out = out[:0]
		out = f.addBegin(out)
		//out = f.addKeyUnsafe(out, "age")
		//out = f.addValInt(out, 39)
		//out = f.addKeyUnsafe(out, "tag")
		//out = f.addTag(out, &al.Control.Tags, IO|TEST)
		out = f.addKeyUnsafe(out, "ts")
		//time.Now()
		t := time.Now()
		//t := time.Now().UTC()
		_ = t
		//out = f.addTimeDay(out, &t)

		//out = f.addTimeUnix(out, &t)

		//y, m, d := t.Date()
		//out = f.addTimeDate(out, y, int(m), d)

		//h,m,s := t.Clock()
		//out = f.addTime(out,h,m,s)

		//h,m,s := t.Clock()
		//out = f.addTimeNano(out, h,m,s,t.Nanosecond())

		//out = f.addKeyUnsafe(out, "level")
		//out = f.addLevel(out, Linfo)
		//out = f.addValInt2(out, 39, 0, ',')
		//out = f.addKey(out, "name")
		//out = f.addKeyUnsafe(out, "name")
		//out = f.addValStringUnsafe(out, "Gon Yi")
		//out = f.addValString(out, "Gon Yi")
		//out = f.addKeyUnsafe(out, "isMarried")
		//out = f.addKey(out, "isMarried")
		//out = f.addValBool(out, false)
		//out = f.addKey(out, "height")
		//out = f.addValFloat(out, 5.8)
		out = f.addEnd(out)
	}

	println(string(out))
}
