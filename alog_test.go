package alog_test

import (
	"errors"
	"github.com/gonyyi/alog"
	"github.com/gonyyi/alog/ext"
	"os"
	"testing"
	"unsafe"
)

var err1 = errors.New("error msg my")

func TestRandom(t *testing.T) {
	al := alog.New(os.Stderr)
	t1 := al.NewTag("t1")
	t2 := al.NewTag("t2")
	t3 := al.NewTag("t3")
	println(unsafe.Sizeof(al))

	println(t1, t2, t3)
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
