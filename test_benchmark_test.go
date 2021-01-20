// (c) 2020 Gon Y Yi. <https://gonyyi.com>
// Version 0.2.0, 12/31/2020

package alog_test

import (
	"github.com/gonyyi/alog"
	"testing"
)

var (
	tStr   = "hello this is gon"
	tInt   = 123
	tInt64 = 123
	tFloat = 1.234
	tBool  = true
)

func BenchmarkLogger_Info(b *testing.B) {
	al := alog.New(nil)
	{
		b.Run("fmt=json,arg=", func(b2 *testing.B) {
			b2.ReportAllocs()
			al.SetFormat(alog.Flevel | alog.Ftag | alog.Fjson)
			for i := 0; i < b2.N; i++ {
				al.Info(0, "test")
			}
		})

		b.Run("fmt=json,arg=int", func(b2 *testing.B) {
			b2.ReportAllocs()
			al.SetFormat(alog.Flevel | alog.Ftag | alog.Fjson)
			for i := 0; i < b2.N; i++ {
				al.Info(0, "test", "a1", 123)
			}
		})

		b.Run("fmt=json,arg=str", func(b2 *testing.B) {
			b2.ReportAllocs()
			al.SetFormat(alog.Flevel | alog.Ftag | alog.Fjson)
			for i := 0; i < b2.N; i++ {
				al.Info(0, "test", "a1", "str1")
			}
		})
		b.Run("fmt=json,arg=str+int", func(b2 *testing.B) {
			b2.ReportAllocs()
			al.SetFormat(alog.Flevel | alog.Ftag | alog.Fjson)
			for i := 0; i < b2.N; i++ {
				al.Info(0, "test", "a1", "str1", "b1", 123)
			}
		})
	}

	{
		b.Run("fmt=text,arg=", func(b2 *testing.B) {
			b2.ReportAllocs()
			al.SetFormat(alog.Flevel | alog.Ftag)
			for i := 0; i < b2.N; i++ {
				al.Info(0, "test")
			}
		})

		b.Run("fmt=text,arg=int", func(b2 *testing.B) {
			b2.ReportAllocs()
			al.SetFormat(alog.Flevel | alog.Ftag)
			for i := 0; i < b2.N; i++ {
				al.Info(0, "test", "a2", 123)
			}
		})
		b.Run("fmt=text,arg=str", func(b2 *testing.B) {
			b2.ReportAllocs()
			al.SetFormat(alog.Flevel | alog.Ftag)
			for i := 0; i < b2.N; i++ {
				al.Info(0, "test", "a1", "str1")
			}
		})
		b.Run("fmt=text,arg=str+int", func(b2 *testing.B) {
			b2.ReportAllocs()
			al.SetFormat(alog.Flevel | alog.Ftag)
			for i := 0; i < b2.N; i++ {
				al.Info(0, "test", "a1", "str1", "a2", 123)
			}
		})
	}
}

// func BenchmarkLogger_Info(b *testing.B) {
// 	b.Run("fmt=json, error=str", func(b2 *testing.B) {
// 		l := alog.New(nil).SetFormat(alog.Fjson)
// 		b2.ReportAllocs()
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Error(testLogMsg)
// 			}
// 		})
// 	})
// 	b.Run("fmt=json+utc, error=str", func(b2 *testing.B) {
// 		l := alog.New(nil).SetFormat(alog.Fjson | alog.FtimeUTC)
// 		b2.ReportAllocs()
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Error(testLogMsg)
// 			}
// 		})
// 	})
// 	b.Run("fmt=json+time, error=str", func(b2 *testing.B) {
// 		l := alog.New(nil).SetFormat(alog.Fjson | alog.Ftime)
// 		b2.ReportAllocs()
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Error(testLogMsg)
// 			}
// 		})
// 	})
// 	b.Run("fmt=json+time, errorf=str", func(b2 *testing.B) {
// 		l := alog.New(nil).SetFormat(alog.Fjson | alog.Ftime)
// 		b2.ReportAllocs()
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Errorf("test: %s", "abc")
// 			}
// 		})
// 	})
// 	b.Run("fmt=json+time, errorf=int", func(b2 *testing.B) {
// 		l := alog.New(nil).SetFormat(alog.Fjson | alog.Ftime)
// 		b2.ReportAllocs()
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Errorf("test: %d", 1234)
// 			}
// 		})
// 	})
// }
