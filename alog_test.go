package alog_test

import (
	"errors"
	"github.com/gonyyi/alog"
	"os"
	"testing"
)

func TestAlogArray(t *testing.T) {
	al := alog.New(os.Stdout).SetFormat(alog.Fjson)
	boolArr := []bool{true, false, true}
	strArr := []string{"okay", "not okay"}
	intArr := []int{1, 2, 3, 4}
	f64Arr := []float64{1, 2, 3, 4}
	errArr := []error{nil, errors.New("err1"), errors.New("err2")}

	al.Info(0, "test", "val", &errArr)
	al.Info(0, "test", "val", &boolArr)
	al.Info(0, "test", "val", &intArr)
	al.Info(0, "test", "val", &f64Arr)
	al.Info(0, "test", "val", &strArr)
	al.Info(0, "test", "val", &strArr)
}

func TestAlog_New(t *testing.T) {
	al := alog.New(os.Stderr)
	USER := al.GetTag("user")
	REQ := al.GetTag("req")

	// al.SetHookFn(func(level alog.Level, tag alog.Tag, bytes []byte) {
	// 	if !tag.Has(REQ) {
	// 		println("-----\nREQ: " + string(bytes) + "\n-----\n")
	// 	}
	// })

	al.SetFormat(al.Format().Off(alog.FtimeUnixMs))

	al.Info(USER|REQ, "test", "name", "gon", "age", 17, "married", false, "req", true)
	al.Info(USER, "test", "name", "gon", "age", 17, "married", false)
	al.Info(0, "test", "name", "gon", "age", 17, "married", false)

	al.SetFormat(al.Format().On(alog.Fjson).Off(alog.FtimeUnixMs))

	al.Info(USER|REQ, "test", "name", "gon", "age", 17, "married", false, "req", true)
	al.Info(USER, "test", "name", "gon", "age", 17, "married", false)
	al.Info(0, "test", "name", "gon", "age", 17, "married", false)

	// al.Info(USER, "test", "name", "gon", "age", 17, "married", false)
	// al.Info(REQ, "test", "name", "gon", "ages", []int{17, 18, 20}, "married", []bool{true, false, true})
	//
	// al.ModFormat(alog.Fjson, true)
	// al.Info(REQ, "test", "name", "gon", "age", 17, "married", false)
	// al.Info(REQ, "test", "name", "gon", "ages", []int{17, 18, 20}, "married", []bool{true, false, true})

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
