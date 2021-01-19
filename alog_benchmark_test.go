// // (c) 2020 Gon Y Yi. <https://gonyyi.com>
// // Version 0.2.0, 12/31/2020
//
package alog_test

//
// import (
// 	"fmt"
// 	"github.com/gonyyi/alog"
// 	"testing"
// )
//
// var (
// 	testLogMsg     = "Pairs for a logger with simple string without any formatting"
// 	testLogMsgFmt  = "Pairs for a logger with %s, number=%d, float=%f"
// 	testLogMsgFmt1 = "formats"
// 	testLogMsgFmt2 = 123
// 	testLogMsgFmt3 = 3.14
// )
//
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
// func BenchmarkLogger_Error(b *testing.B) {
// 	b.Run("skip by level, error=str", func(b2 *testing.B) {
// 		l := alog.New(nil).SetFormat(alog.Fnone).SetLogLevel(alog.Lfatal)
// 		b2.ReportAllocs()
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Error(testLogMsg)
// 			}
// 		})
// 	})
// 	b.Run("skip by tag, error=str", func(b2 *testing.B) {
// 		l := alog.New(nil).SetFormat(alog.Fnone)
// 		test1 := l.NewTag("test")
// 		test2 := l.NewTag("test")
// 		l.SetLogTag(test1)
//
// 		b2.ReportAllocs()
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Print(alog.Lerror, test2, testLogMsg)
// 			}
// 		})
// 	})
// 	b.Run("fmt=none, error=str", func(b2 *testing.B) {
// 		l := alog.New(nil).SetFormat(alog.Fnone)
// 		b2.ReportAllocs()
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Error(testLogMsg)
// 			}
// 		})
// 	})
// 	b.Run("fmt=default, error=str", func(b2 *testing.B) {
// 		l := alog.New(nil)
// 		b2.ReportAllocs()
//
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Error(testLogMsg)
// 			}
// 		})
// 	})
// 	b.Run("fmt=default+utc, error=str", func(b2 *testing.B) {
// 		l := alog.New(nil).SetFormat(alog.FtimeUTC)
// 		b2.ReportAllocs()
//
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Error(testLogMsg)
// 			}
// 		})
// 	})
// }
//
// func BenchmarkLogger_Errorf(b *testing.B) {
// 	b.Run("fmt=none, errorf=fmt", func(b2 *testing.B) {
// 		l := alog.New(nil).SetFormat(alog.Fnone)
// 		b2.ReportAllocs()
//
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Errorf(testLogMsg)
// 			}
// 		})
// 	})
// 	b.Run("fmt=default, errorf=fmt", func(b2 *testing.B) {
// 		l := alog.New(nil)
// 		b2.ReportAllocs()
//
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Errorf(testLogMsg)
// 			}
// 		})
// 	})
//
// 	b.Run("fmt=default+utc, errorf=fmt", func(b2 *testing.B) {
// 		l := alog.New(nil).SetFormat(alog.FtimeUTC)
// 		b2.ReportAllocs()
//
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Errorf(testLogMsg)
// 			}
// 		})
// 	})
//
// 	b.Run("fmt=none, errorf=int", func(b2 *testing.B) {
// 		l := alog.New(nil).SetFormat(alog.Fnone)
// 		b2.ReportAllocs()
//
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Errorf("test %d", 123)
// 			}
// 		})
// 	})
// 	b.Run("fmt=default, errorf=float", func(b2 *testing.B) {
// 		l := alog.New(nil)
// 		b2.ReportAllocs()
//
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Errorf("test %f", 1.234)
// 			}
// 		})
// 	})
// 	b.Run("fmt=default+utc, errorf=str", func(b2 *testing.B) {
// 		l := alog.New(nil).SetFormat(alog.FtimeUTC)
// 		b2.ReportAllocs()
//
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Errorf("test %s", "string")
// 			}
// 		})
// 	})
// 	b.Run("fmt=default+utc, errorf=bool", func(b2 *testing.B) {
// 		l := alog.New(nil).SetFormat(alog.FtimeUTC)
// 		b2.ReportAllocs()
//
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Errorf("test %t", true)
// 			}
// 		})
// 	})
// }
//
//
// func BenchmarkLogger_Print(b *testing.B) {
// 	b.Run("print", func(b2 *testing.B) {
// 		// S1: 1393 ns/op	       0 B/op	       0 allocs/op
// 		// S2:  770 ns/op	       0 B/op	       0 allocs/op
//
// 		l := alog.New(nil)
//
// 		b2.ReportAllocs()
// 		b2.ResetTimer()
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				l.Print(alog.Linfo, 0, "abc")
// 			}
// 		})
// 	})
// }
// func BenchmarkLogger_Printf(b *testing.B) {
// 	b.Run("printf with no arg", func(b2 *testing.B) {
// 		// S1: 1393 ns/op	       0 B/op	       0 allocs/op
// 		// S2:  770 ns/op	       0 B/op	       0 allocs/op
//
// 		l := alog.New(nil)
//
// 		b2.ReportAllocs()
// 		b2.ResetTimer()
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				l.Printf(alog.Linfo, 0, "abc")
// 			}
// 		})
// 	})
// }
//
// // Added as of 0.1.1 update, 12/29/2020
// func BenchmarkLogger_NewPrint(b *testing.B) {
// 	// S2:  293 ns/op	       0 B/op	       0 allocs/op
//
// 	l := alog.New(nil).SetFormat(alog.Fdefault).SetPrefix("test ")
// 	CAT1, CAT2 := l.NewTag("cat1"), l.NewTag("cat2")
// 	l.SetLogTag(CAT1)
// 	WarnCAT1 := l.NewPrint(alog.Lwarn, CAT1, "CAT1w ")
// 	WarnCAT2 := l.NewPrint(alog.Lwarn, CAT2, "CAT2w ")
// 	TraceCAT1 := l.NewPrint(alog.Ltrace, CAT1, "CAT1t ")
// 	TraceCAT2 := l.NewPrint(alog.Ltrace, CAT2, "CAT2t ")
//
// 	b.ReportAllocs()
// 	b.ResetTimer()
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			WarnCAT1("warn cat1 test")
// 			WarnCAT2("warn cat2 test")
// 			TraceCAT1("trace cat1 test")
// 			TraceCAT2("trace cat2 test")
// 		}
// 	})
// }
//
// // Added as of 0.1.2 update, 12/29/2020
// func BenchmarkLogger_NewWriter(b *testing.B) {
// 	// S2: 401 ns/op	       0 B/op	       0 allocs/op
// 	// S2: 438 ns/op	       0 B/op	       0 allocs/op
//
// 	l := alog.New(nil).SetFormat(alog.Fprefix | alog.Flevel).SetPrefix("nptest ")
// 	TEST1, TEST2, TEST3 := l.NewTag("test1"), l.NewTag("test2"), l.NewTag("test3")
//
// 	l.SetLogTag(TEST2) // only show TEST2
//
// 	wT1D := l.NewWriter(alog.Ldebug, TEST1, "T1D ")
// 	wT1I := l.NewWriter(alog.Linfo, TEST1, "T1I ")
// 	wT2D := l.NewWriter(alog.Ldebug, TEST2, "T2D ")
// 	wT2I := l.NewWriter(alog.Linfo, TEST2, "T2I ")
// 	wT3D := l.NewWriter(alog.Ldebug, TEST3, "T3D ")
// 	wT3I := l.NewWriter(alog.Linfo, TEST3, "T3I ")
//
// 	b.ReportAllocs()
// 	b.ResetTimer()
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			_, _ = fmt.Fprintf(wT1D, "test: %s fprintf", "T1D") // Not printed
// 			_, _ = fmt.Fprintf(wT1I, "test: %s fprintf", "T1I") // Not printed
// 			_, _ = fmt.Fprintf(wT2D, "test: %s fprintf", "T2D") // Not printed
// 			_, _ = fmt.Fprintf(wT2I, "test: %s fprintf", "T2I") // Printed
// 			_, _ = fmt.Fprintf(wT3D, "test: %s fprintf", "T3D") // Not printed
// 			_, _ = fmt.Fprintf(wT3I, "test: %s fprintf", "T3I") // Not printed
// 		}
// 	})
// }
