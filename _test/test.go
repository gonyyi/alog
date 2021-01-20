package main

import (
	"github.com/gonyyi/alog"
	"os"
)

func main() {
	l := alog.New(os.Stdout).SetNewTags("t1", "t2", "t3")
	t1, t2 := l.MustGetTag("t1"), l.MustGetTag("t2")

	//l.SetFilter(0, t1)
	l.SetFormat(alog.Flevel | alog.Ftag | alog.FtimeUnix)
	l.Error(t1, "hello error 1")
	l.Error(t2, "hello error 2")
	l.Error(t1|t2, "hello error 1")

	l.SetFormat(alog.Fjson | alog.FtimeUnix)
	l.Error(0, "hello error 1")
	l.Error(t1, "hello error 1", "a", "b")
	l.Error(t2, "hello error 1", "a", "b", "c")
	l.Error(t1|t2, "hello error 1", "a", "b", "c", "d")

}
