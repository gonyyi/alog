package alog_test

import (
	"errors"
	"github.com/gonyyi/alog"
	"github.com/gonyyi/alog/ext"
	"os"
	"testing"
)

var err1 = errors.New("error msg my")

func TestRandom(t *testing.T) {
	al := alog.New(os.Stderr)
	al.Info(0)
}

func TestNew(t *testing.T) {
	al := alog.New(os.Stderr)

	tOS := al.NewTag("OS")
	tSYS := al.NewTag("SYS")

	test := func() {
		al.Info(tOS|tSYS).Err("err1", nil).Err("err2", err1).Str("ok", "yes okay").Write("")
	}

	test()
	al.SetFormatter(ext.NewFormatterTerminal())
	test()
}
