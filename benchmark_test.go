package alog_test

import (
	"github.com/gonyyi/alog"
	"os"
	"testing"
)

func BenchmarkLogger_NewWriter(b *testing.B) {
	al := alog.New(nil).SetNewTags("backend", "frontend", "user", "req")
	al.UpdateFormat(alog.FdateDay, true)
	//al.UpdateFormat(alog.Fjson, true)
	// al.UpdateFormat(alog.Ftime|alog.Fdate, false)
	USER := al.GetTag("user")
	REQ := al.GetTag("req")
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

	al := alog.New(nil).SetNewTags("backend", "frontend", "user", "req").SetFormatter(alog.NewFormatterJSON())
	al.UpdateFormat(alog.Ftime|alog.Fdate, false)
	USER := al.GetTag("user")
	REQ := al.GetTag("req")

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
