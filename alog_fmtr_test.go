package alog

import (
	"testing"
	"time"
)

var testCusHead = []byte("ok to do ")

func cushead(dst []byte) []byte {
	return append(dst, testCusHead...)
}

func TestFormatterJSON(t *testing.T) {
	var l AlogFmtr

	var buf []byte

	tag0 := []string{}
	tag1 := []string{"name"}
	tag2 := []string{"name", "age"}
	_, _, _ = tag0, tag1, tag2 // prevent error
	ts := time.Now()

	f := func() {
		buf = l.Begin(buf, nil)
		// buf = cushead(buf)
		buf = l.LogTime(buf, ts)
		buf = l.Space(buf)
		buf = l.LogTimeDate(buf, ts)
		buf = l.Space(buf)
		buf = l.LogTimeDay(buf, ts)
		buf = l.Space(buf)
		buf = l.LogTimeUnix(buf, ts)
		buf = l.Space(buf)
		buf = l.LogTimeUnixMs(buf, ts)
		buf = l.Space(buf)
		// buf = l.LogTag(buf, tag0)
		buf = l.Space(buf)
		buf = l.LogMsg(buf, `hello "world"`, ';')
		buf = l.Space(buf)
		buf = l.String(buf, "name", "gon")
		buf = l.Space(buf)
		buf = l.Float32(buf, "weight32", 123.450)
		buf = l.Space(buf)
		buf = l.Float64(buf, "weight64", 123.450)
		buf = l.Space(buf)
		buf = l.Strings(buf, "names", []string{"gon", "gone", "yi", "young"})
		buf = l.Space(buf)
		buf = l.Nil(buf, "ageinfo")
		buf = l.Space(buf)
		buf = l.Float32s(buf, "f32", []float32{1.1, 1.2, 1.3, 1.4})
		buf = l.Space(buf)
		buf = l.Float64s(buf, "f64", []float64{1.1, 1.2, 1.3, 1.4})
		buf = l.End(buf)
	}

	l = &JSONFmtr{}
	f()
	println(string(buf))
	l = &TextFmtr{}
	f()
	println(string(buf))
}

func BenchmarkFormatterJSON(b *testing.B) {
	b.ReportAllocs()
	var l AlogFmtr
	// l = &TextFmtr{}
	l = &JSONFmtr{}

	var buf []byte
	ts := time.Now()
	tag0 := []string{}
	tag1 := []string{"name"}
	tag2 := []string{"name", "age"}
	f32 := []float32{1.1, 1.2, 1.3, 1.4}
	f64 := []float64{1.1, 1.2, 1.3, 1.4}
	_, _, _ = tag0, tag1, tag2 // prevent error

	for i := 0; i < b.N; i++ {
		ts = time.Now()
		buf = l.Begin(buf, nil)
		// buf = cushead(buf)
		// buf = l.Space(buf)
		buf = l.LogTime(buf, ts)
		buf = l.Space(buf)
		buf = l.LogTimeDate(buf, ts)
		buf = l.Space(buf)
		buf = l.LogTimeDay(buf, ts)
		buf = l.Space(buf)
		buf = l.LogTimeUnix(buf, ts)
		buf = l.Space(buf)
		buf = l.LogTimeUnixMs(buf, ts)
		buf = l.Space(buf)
		buf = l.Float32(buf, "weight1", 123.45)
		buf = l.Space(buf)
		buf = l.Float64(buf, "weight2", 123.45)
		buf = l.Space(buf)
		// buf = l.LogTag(buf, tag0)
		buf = l.Space(buf)
		buf = l.LogMsg(buf, `hello "world"`, ';')
		buf = l.Space(buf)
		buf = l.String(buf, "name", "gon")
		buf = l.Space(buf)
		buf = l.Strings(buf, "names", tag2)
		buf = l.Space(buf)
		buf = l.Float32s(buf, "f32", f32)
		buf = l.Space(buf)
		buf = l.Float64s(buf, "f64", f64)
		buf = l.End(buf)
	}

	print(string(buf))
}
