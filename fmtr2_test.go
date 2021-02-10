package alog_test

import (
	"github.com/gonyyi/alog"
	"testing"
)

func TestFmtr2(t *testing.T) {
	var tgr alog.Tagger
	TEST := tgr.MustGetTag("test")
	_ = TEST
	tgr.Filter.Set(alog.Linfo, 0)

	var bufHead, bufBody []byte

	var f alog.Fmtr2
	f = &alog.Fmtr2JSON{}
	f.SetTagger(&tgr)

	f.Log(&bufHead, &bufBody, alog.Linfo, TEST, "hello")

	println("head:", string(bufHead))
	println("body:", string(bufBody))
}

func BenchmarkFmtr2(b *testing.B) {
	var tgr alog.Tagger
	TEST := tgr.MustGetTag("test")
	NAME := tgr.MustGetTag("name")
	URL := tgr.MustGetTag("url")
	_, _, _ = TEST, NAME, URL

	tgr.Filter.Set(alog.Linfo, 0)

	var bufHead, bufBody []byte

	var f alog.Fmtr2
	f = &alog.Fmtr2JSON{}
	f.SetFormat(alog.Ftag|alog.Flevel)
	f.SetTagger(&tgr)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		bufHead, bufBody = bufHead[:0], bufBody[:0]
		f.Log(&bufHead, &bufBody, alog.Linfo, TEST|NAME|URL, "hello")
	}
	//println("H:", string(bufHead))
	//println("B:", string(bufBody))
	println(string(append(bufHead,bufBody...)))
}
