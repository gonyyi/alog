package alog_test

import (
	"errors"
	"github.com/gonyyi/alog"
	"github.com/gonyyi/alog/ext"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	al := alog.New(os.Stderr)
	e1 := errors.New("error msg my")
	os := al.NewTag("OS")
	sys := al.NewTag("SYS")
	test := func() {
		al.Info(os|sys).Err("err1", nil).Err("err2", e1).Str("ok", "yes okay").Write("log starts")
	}

	test()
	al.CusFmat = ext.NewFormatterTerminal()
	test()
}
