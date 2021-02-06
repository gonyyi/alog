package alog_test

import (
	"github.com/gonyyi/alog"
	"testing"
	"time"
)

func TestTextFmtr(t *testing.T) {
	var l alog.Fmtr
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
		buf = l.Strings(buf, "names", &[]string{"gon", "gone", "yi", "young"})
		buf = l.Space(buf)
		buf = l.Nil(buf, "ageinfo")
		buf = l.Space(buf)
		buf = l.Float32s(buf, "f32", &[]float32{1.1, 1.2, 1.3, 1.4})
		buf = l.Space(buf)
		buf = l.Float64s(buf, "f64", &[]float64{1.1, 1.2, 1.3, 1.4})
		buf = l.End(buf)
	}
	l = alog.Default.NewFmtText()

	f()
	println(string(buf))
}
