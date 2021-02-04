package alog_test

import (
	"github.com/gonyyi/alog"
	"os"
	"testing"
)

func TestAlog_New(t *testing.T) {
	al := alog.New(os.Stderr).
		SetFormatItem(alog.Flevel, true).
		SetFormatItem(alog.Fjson, true).
		SetNewTags("backend", "frontend", "user", "req")

	USER := al.MustGetTag("user")
	REQ := al.MustGetTag("req")

	al.Info(USER|REQ, "test", "name", "gon", "age", 17, "married", false)
	al.SetFormatItem(alog.Fjson, false)
	al.Info(USER, "test", "name", "gon", "age", 17, "married", false)
	al.SetFormatItem(alog.Fjson, true)
	al.Info(REQ, "test", "name", "gon", "age", 17, "married", false)
	al.SetFormatItem(alog.Fjson, false)
	al.Info(0, "test", "name", "gon", "age", 17, "married", false)

	// JSON
	// {"d":20210125,"t":102100324,"level":"info","tag":["user","req"],"msg":"test","name":"gon","age":17,"married":false}
	// {"d":20210125,"t":102100324,"level":"info","tag":["user"],"msg":"test","name":"gon","age":17,"married":false}
	// {"d":20210125,"t":102100324,"level":"info","tag":["req"],"msg":"test","name":"gon","age":17,"married":false}
	// {"d":20210125,"t":102100324,"level":"info","tag":[],"msg":"test","name":"gon","age":17,"married":false}

	// TEXT
	// 20210125 102128.161 info tag=[user,req] test; name="gon" age=17 married=false
	// 20210125 102128.161 info tag=[user] test; name="gon" age=17 married=false
	// 20210125 102128.161 info tag=[req] test; name="gon" age=17 married=false
	// 20210125 102128.161 info tag=[] test; name="gon" age=17 married=false
}

func BenchmarkLogger_NewWriter(b *testing.B) {
	al := alog.New(nil).SetNewTags("backend", "frontend", "user", "req")
	al.SetFormatItem(alog.FdateDay, true)
	//al.SetFormatItem(alog.Fjson, true)
	// al.SetFormatItem(alog.Ftime|alog.Fdate, false)
	USER := al.MustGetTag("user")
	REQ := al.MustGetTag("req")
	sw := al.NewWriter(alog.Linfo, USER|REQ)
	txt := []byte("sub writer test")
	// al.SetFormatter(nil)

	b.Run("alog", func(c *testing.B) {
		c.ReportAllocs()
		for i := 0; i < c.N; i++ {
			// sw.WriteString("test") // this doesn't have allocation
			sw.Write(txt) // this has allocation
			sw.WriteString("sub writer tet2")
			// todo: subWriter will take 1 allocation when converting to string.. Create separate one for byte array.
		}
	})
	al.SetOutput(os.Stderr)
	sw.Write(txt)
}

func BenchmarkLogger_New(b *testing.B) {
	// BenchmarkLogger_New/msg-12         	 6399504	       170 ns/op	       0 B/op	       0 allocs/op
	// BenchmarkLogger_New/msg+s-12       	 6389883	       186 ns/op	       0 B/op	       0 allocs/op
	// BenchmarkLogger_New/s+s-12         	 6095625	       199 ns/op	       0 B/op	       0 allocs/op
	// BenchmarkLogger_New/msg+s+i+b-12   	 4642806	       262 ns/op	       0 B/op	       0 allocs/op

	al := alog.New(nil).SetNewTags("backend", "frontend", "user", "req").SetFormatter(alog.NewFormatterJSON()).SetFormatItem(alog.Ftime|alog.Fdate, false)
	USER := al.MustGetTag("user")
	REQ := al.MustGetTag("req")

	b.Run("msg", func(c *testing.B) {
		c.ReportAllocs()
		for i := 0; i < c.N; i++ {
			al.Info(USER|REQ, "message")
			// al.Info(0, "test", "val", "okay")
			// al.Info(0, "", "val\t\r", "ok\tay", "message", "te\tst")
			// al.Info(0, "test", "name", "gon", "age", 17, "married", false)
		}
	})

	b.Run("msg+s", func(c *testing.B) {
		c.ReportAllocs()
		for i := 0; i < c.N; i++ {
			// al.Info(USER|REQ, "message")
			al.Info(0, "test", "val", "okay")
			// al.Info(0, "", "val\t\r", "ok\tay", "message", "te\tst")
			// al.Info(0, "test", "name", "gon", "age", 17, "married", false)
		}
	})

	b.Run("s+s", func(c *testing.B) {
		c.ReportAllocs()
		for i := 0; i < c.N; i++ {
			// al.Info(USER|REQ, "message")
			// al.Info(0, "test", "val", "okay")
			al.Info(0, "", "val\t\r", "ok\tay", "message", "te\tst")
			// al.Info(0, "test", "name", "gon", "age", 17, "married", false)
		}
	})

	b.Run("msg+s+i+b", func(c *testing.B) {
		c.ReportAllocs()
		for i := 0; i < c.N; i++ {
			// al.Info(USER|REQ, "message")
			// al.Info(0, "test", "val", "okay")
			// al.Info(0, "", "val\t\r", "ok\tay", "message", "te\tst")
			al.Info(0, "test", "name", "gon", "age", 17, "married", false)
		}
	})
}
