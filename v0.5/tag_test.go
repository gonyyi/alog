package alog_test

import (
	alog "github.com/gonyyi/alog/v0.5"
	"testing"
)

func TestTagBucket_AppendSelectedTags(t *testing.T) {
	tb := alog.TagBucket{}
	tIo := tb.MustGetTag("io")
	tNet := tb.MustGetTag("net")
	tSys := tb.MustGetTag("sys")
	tIoNet := tIo | tNet
	if true != tIoNet.Has(tIo) || tIoNet.Has(tNet) != true || tIoNet.Has(tSys) != false {
		t.Error(tIoNet.Has(tIo), tIoNet.Has(tNet), tIoNet.Has(tSys))
	}
	{
		var out []byte
		out = tb.AppendSelectedTags(out[:0], ',', true, tIoNet)
		if string(out) != `"io","net"` {
			t.Error(string(out))
		}
		out = tb.AppendSelectedTags(out[:0], ',', false, tIoNet)
		if string(out) != `io,net` {
			t.Error(string(out))
		}
		out = tb.AppendSelectedTags(out[:0], ',', false, tSys)
		if string(out) != `sys` {
			t.Error(string(out))
		}
		out = tb.AppendSelectedTags(out[:0], ',', false, 0)
		if string(out) != `` {
			t.Error(string(out))
		}
		out = tb.AppendSelectedTags(out[:0], ',', false, 16)
		if string(out) != `` {
			t.Error(string(out))
		}
	}
}
