package alog_test

import (
	"errors"
	"github.com/gonyyi/alog"
	"os"
	"testing"
)

func BenchmarkLogger_NewWriter(b *testing.B) {
	al := alog.New(nil)
	al.SetFormat(alog.Fdefault.Off(alog.Ftime | alog.Fdate))

	var USER, REQ, SW alog.Tag = 0, 0, 0
	_, _, _ = USER, REQ, SW

	USER = al.GetTag("user")
	REQ = al.GetTag("req")
	SW = al.GetTag("sw")

	// al.SetFilter(0, SW)

	sw := al.NewWriter(alog.Linfo, USER)
	txt := []byte("sub writer test")

	b.Run("subWriter-text", func(c *testing.B) {
		c.ReportAllocs()
		for i := 0; i < c.N; i++ {
			sw.Write(txt) // this has allocation
		}
	})

	al.SetFormat(alog.Fdefault.Off(alog.Ftime | alog.Fdate).On(alog.Fjson))
	b.Run("subWriter-JSON", func(c *testing.B) {
		c.ReportAllocs()
		for i := 0; i < c.N; i++ {
			sw.Write(txt) // this has allocation
		}
	})
}

func BenchmarkLogger_GetTag(b *testing.B) {
	al := alog.New(nil).SetFormat(alog.Fjson)

	b.Run("tag", func(c *testing.B) {
		c.ReportAllocs()
		test := al.GetTag("test")
		al.SetFormat(alog.Fjson | alog.Ftag)
		for i := 0; i < c.N; i++ {
			al.Info(test, "test")
		}
	})

	b.Run("notag", func(c *testing.B) {
		c.ReportAllocs()
		al.SetFormat(alog.Fjson | alog.Ftag)
		for i := 0; i < c.N; i++ {
			al.Info(0, "test")
		}
	})

	{
		al.SetOutput(os.Stdout)
		test := al.GetTag("test")
		al.SetFormat(alog.Fjson | alog.Ftag)
		al.Info(test, "test tag")
		al.SetFormat(alog.Fjson)
		al.Info(test, "test notag")
	}

}
func BenchmarkLogger_Array(b *testing.B) {
	al := alog.New(nil).SetFormat(alog.Flevel)
	boolArr := []bool{true, false, true}
	strArr := []string{"okay", "not okay"}
	intArr := []int{1, 2, 3}
	f64Arr := []float64{1, 2, 3}
	f32Arr := []float32{1, 2, 3}
	errArr := []error{nil, errors.New("err1"), errors.New("err2")}

	b.Run("msg+e[]", func(c *testing.B) {
		c.ReportAllocs()
		for i := 0; i < c.N; i++ {
			al.Info(0, "test", "val", &errArr)
		}
	})

	b.Run("msg+b[]", func(c *testing.B) {
		c.ReportAllocs()
		for i := 0; i < c.N; i++ {
			al.Info(0, "test", "val", &boolArr)
		}
	})

	b.Run("msg+i[]", func(c *testing.B) {
		c.ReportAllocs()
		for i := 0; i < c.N; i++ {
			al.Info(0, "test", "val", &intArr)
		}
	})
	b.Run("msg+f32[]", func(c *testing.B) {
		c.ReportAllocs()
		for i := 0; i < c.N; i++ {
			al.Info(0, "test", "val", &f32Arr)
		}
	})

	b.Run("msg+f64[]", func(c *testing.B) {
		c.ReportAllocs()
		for i := 0; i < c.N; i++ {
			al.Info(0, "test", "val", &f64Arr)
		}
	})

	b.Run("msg+s[]", func(c *testing.B) {
		c.ReportAllocs()
		for i := 0; i < c.N; i++ {
			al.Info(0, "test", "val", &strArr)
		}
	})

	b.Run("msg+s[*]", func(c *testing.B) {
		c.ReportAllocs()
		for i := 0; i < c.N; i++ {
			al.Info(0, "test", "val", &strArr)
		}
	})
}
func BenchmarkLogger_New(b *testing.B) {
	// BenchmarkLogger_New/msg-12         	 6399504	       170 ns/op	       0 B/op	       0 allocs/op
	// BenchmarkLogger_New/msg+s-12       	 6389883	       186 ns/op	       0 B/op	       0 allocs/op
	// BenchmarkLogger_New/s+s-12         	 6095625	       199 ns/op	       0 B/op	       0 allocs/op
	// BenchmarkLogger_New/msg+s+i+b-12   	 4642806	       262 ns/op	       0 B/op	       0 allocs/op

	al := alog.New(nil)
	al.SetFormat(alog.Fjson)

	USER := al.GetTag("user")
	REQ := al.GetTag("req")
	_, _ = USER, REQ

	strArr := []string{"okay", "not okay"}

	b.Run("msg", func(c *testing.B) {
		c.ReportAllocs()
		for i := 0; i < c.N; i++ {
			al.Info(USER, "message")
		}
	})

	b.Run("msg+s", func(c *testing.B) {
		c.ReportAllocs()
		for i := 0; i < c.N; i++ {
			al.Info(USER, "test", "val", "okay")
		}
	})

	b.Run("s+s", func(c *testing.B) {
		c.ReportAllocs()
		for i := 0; i < c.N; i++ {
			al.Info(USER, "", "val\t\r", "ok\tay", "message", "te\tst")
		}
	})

	b.Run("msg+s+i+b", func(c *testing.B) {
		c.ReportAllocs()
		for i := 0; i < c.N; i++ {
			al.Info(USER, "test", "name", "gon", "age", 17, "married", false)
		}
	})

	al.SetOutput(os.Stdout)
	al.Info(USER, "message")
	al.Info(USER, "test", "val", "okay")
	al.Info(USER, "test", "val", strArr)
	al.Info(USER, "test", "val", &strArr)
	al.Info(USER, "", "val\t\r", "ok\tay", "message", "te\tst")
	al.Info(USER, "test", "name", "gon", "age", 17, "married", false)
}
