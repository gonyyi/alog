package alog_test

import (
	"github.com/gonyyi/alog"
	"os"
	"testing"
)

func TestAlog_New(t *testing.T) {
	al := alog.New(os.Stderr).
		SetFormatItem(alog.Flevel, true).
		SetNewTags("backend", "frontend", "user", "req").
		SetFormatter(alog.NewFmtrText())

	USER := al.MustGetTag("user")
	REQ := al.MustGetTag("req")

	al.Info(USER|REQ, "test", "name", "gon", "age", 17, "married", false)
	al.Info(USER, "test", "name", "gon", "age", 17, "married", false)
	al.Info(REQ, "test", "name", "gon", "age", 17, "married", false)
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
	al := alog.New(nil).SetNewTags("backend", "frontend", "user", "req").SetFormatter(alog.NewFmtrJSON())
	al.SetFormatItem(alog.FdateDay, true)
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
			// todo: subWriter will take 1 allocation when converting to string.. Create separate one for byte array.
		}
	})
	al.SetOutput(os.Stderr)
	sw.Write(txt)
}

func BenchmarkLogger_New(b *testing.B) {
	al := alog.New(nil).SetNewTags("backend", "frontend", "user", "req").SetFormatter(alog.NewFmtrJSON()).SetFormatItem(alog.Ftime|alog.Fdate, false)
	USER := al.MustGetTag("user")
	REQ := al.MustGetTag("req")

	b.Run("alog", func(c *testing.B) {
		c.ReportAllocs()
		for i := 0; i < c.N; i++ {
			al.Info(USER|REQ, "message")
			// al.Info(0, "test", "val", "okay")
			// al.Info(0, "", "val\t\r", "ok\tay", "message", "te\tst")
			// al.Info(0, "test", "name", "gon", "age", 17, "married", false)
		}
	})

	al.SetOutput(os.Stderr)
	al.Info(USER|REQ, "message", "name", "gon", "age", 17, "lat", 123.45, "lon", 456.789)
	// al.Info(USER|REQ, "test", "name", "gon", "age", 17, "married", false)

	// zl = zl.Output(os.Stderr)
	// zl.Info().Str("val\t\r", "ok\tay").Msg("te\tst")

}
