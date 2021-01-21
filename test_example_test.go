// (c) 2020 Gon Y Yi. <https://gonyyi.com>
// Version 0.1.3, 12/29/2020

package alog_test

import (
	"github.com/gonyyi/alog"
	"os"
)

func ExampleNew() {
	// Default alog will record date (YYYYMMDD) and time.
	// So disable date and time, and only show level for output comparison
	l := alog.New(os.Stdout).SetNewTags("t1", "t2", "t3")
	t1, t2 := l.MustGetTag("t1"), l.MustGetTag("t2")

	f := func() {
		l.Error(t1, "hello error 1")
		l.Error(t1|t2, "hello error 1/2")
		l.Error(t1|t2, "") // this shouldn't be printed
		l.Error(0, "starting new test")
		l.Error(0, "", "name", "gon") // this should be printed
		l.Error(0, "arg", "name", "gon")
		l.Error(0, "arg", "age", 17)
		l.Error(0, "arg", "weight", 180.1)
		l.Error(0, "arg", "name", "gon", "age", 17, "weight", 180.1)
		l.Error(0, "bad arg", "name", "gon", "age", 17, "weight") // should add null

		l.SetFormatItem(alog.Flevel, false)
		l.SetFormatItem(alog.Ftag, false)
		l.Error(t1|t2, "t1+t2")
		l.Error(t1|t2, "") // this shouldn't be printed
		l.Error(0, "blah")
		l.Error(0, "", "emptyMsg", true) // this should be printed
		l.Error(0, "blah", "name", "gon")

	}
	l.SetFormat(0)
	l.SetFormatItem(alog.Fdate|alog.Ftag, true).SetFormatItem(alog.Fjson, true)
	f()

	println()
	l.SetFormatItem(alog.FtimeMs|alog.Ftag, true).SetFormatItem(alog.Fjson|alog.Fdate, false)
	f()

	// Output:
	// {"lv":"error","tag":["t1"],"msg":"hello error 1"}
	// {"lv":"error","tag":["t1","t2"],"msg":"hello error 1/2"}
	// {"lv":"error","msg":"starting new test"}
	// {"lv":"error","msg":"arg","name":"gon"}
	// {"lv":"error","msg":"arg","age":17}
	// {"lv":"error","msg":"arg","weight":180.1}
	// {"lv":"error","msg":"arg","name":"gon","age":17,"weight":180.1}
	// {"lv":"error","msg":"bad arg","name":"gon","age":17,"weight":null}
	// ERR tag=[t1] msg="hello error 1"
	// ERR tag=[t1,t2] msg="hello error 1/2"
	// ERR msg="starting new test"
	// ERR msg="arg" name="gon"
	// ERR msg="arg" age=17
	// ERR msg="arg" weight=180.1
	// ERR msg="arg" name="gon" age=17 weight=180.1
	// ERR msg="bad arg" name="gon" age=17 weight
	//
}
