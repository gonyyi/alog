package alog_test

import (
	"github.com/gonyyi/alog"
	"os"
	"testing"
)

func TestAlog_New(t *testing.T) {
	al := alog.New(os.Stderr)
	al.UpdateFormat(alog.FtimeUnix, true)
	USER := al.GetTag("user")
	REQ := al.GetTag("req")

	al.SetTriggerFn(func(level alog.Level, tag alog.Tag, bytes []byte) {
		if tag&USER != 0 {
			println("-----\nUSER: " + string(bytes) + "\n-----\n")
		}
	})

	al.Info(USER|REQ, "test", "name", "gon", "age", 17, "married", false)
	al.Info(USER, "test", "name", "gon", "age", 17, "married", false)
	al.Info(0, "test", "name", "gon", "age", 17, "married", false)
	al.UpdateFormat(alog.Fjson, true)
	al.Info(USER|REQ, "test", "name", "gon", "age", 17, "married", false)
	al.Info(USER, "test", "name", "gon", "age", 17, "married", false)
	al.Info(0, "test", "name", "gon", "age", 17, "married", false)

	//al.Info(USER, "test", "name", "gon", "age", 17, "married", false)
	//al.Info(REQ, "test", "name", "gon", "ages", []int{17, 18, 20}, "married", []bool{true, false, true})
	//
	//al.UpdateFormat(alog.Fjson, true)
	//al.Info(REQ, "test", "name", "gon", "age", 17, "married", false)
	//al.Info(REQ, "test", "name", "gon", "ages", []int{17, 18, 20}, "married", []bool{true, false, true})

	// JSON
	// {"d":20210125,"t":102100324,"level":"info","wTag":["user","req"],"msg":"test","name":"gon","age":17,"married":false}
	// {"d":20210125,"t":102100324,"level":"info","wTag":["user"],"msg":"test","name":"gon","age":17,"married":false}
	// {"d":20210125,"t":102100324,"level":"info","wTag":["req"],"msg":"test","name":"gon","age":17,"married":false}
	// {"d":20210125,"t":102100324,"level":"info","wTag":[],"msg":"test","name":"gon","age":17,"married":false}

	// TEXT
	// 20210125 102128.161 info wTag=[user,req] test; name="gon" age=17 married=false
	// 20210125 102128.161 info wTag=[user] test; name="gon" age=17 married=false
	// 20210125 102128.161 info wTag=[req] test; name="gon" age=17 married=false
	// 20210125 102128.161 info wTag=[] test; name="gon" age=17 married=false
}
