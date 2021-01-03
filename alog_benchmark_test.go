// (c) 2020 Gon Y Yi. <https://gonyyi.com>
// Version 0.2.0, 12/31/2020

package alog_test

import (
	"fmt"
	"github.com/gonyyi/alog"
	"io/ioutil"
	"log"
	"testing"
)

/*
Based on MacBook Pro (15-inch, 2018)
- MacOS 10.15.7 Catalina
- 2.9 GHz 6-Core Intel Core i9
- 32 GB 2400 MHz DDR4
- Radeon Pro 560X 4 GB / Intel UHD Graphics 630 1536 MB

| Type    | Name                           | Test               | Count      | Speed       | Mem     | Alloc       |
|:--------|:-------------------------------|:-------------------|:-----------|:------------|:--------|:------------|
| Builtin | BenchmarkBuiltinLoggerBasic-12 |                    | 2969100    | 408 ns/op   | 80 B/op | 2 allocs/op |
| Builtin | BenchmarkBuiltinLoggerFmt-12   |                    | 2534346    | 477 ns/op   | 88 B/op | 3 allocs/op |
| Alog    | BenchmarkLogger_Info           | 1_eval_0_print-12  | 1000000000 | 0.420 ns/op | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_Info           | 5_eval_0_prints-12 | 651124946  | 1.89 ns/op  | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_Info           | 5_eval_3_prints-12 | 1517264    | 777 ns/op   | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_Infof          | 1_eval_0_print-12  | 1000000000 | 1.10 ns/op  | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_Infof          | 5_eval_0_prints-12 | 227060984  | 5.48 ns/op  | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_Infof          | 5_eval_3_prints-12 | 1078039    | 1112 ns/op  | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_Print          | 5_eval_1_prints-12 | 4073192    | 304 ns/op   | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_Print          | 5_eval_2_prints-12 | 2320356    | 478 ns/op   | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_Printf         | 5_eval_1_prints-12 | 2938634    | 419 ns/op   | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_Printf         | 5_eval_2_prints-12 | 1502662    | 805 ns/op   | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_NewPrint-12    |                    | 3516364    | 337 ns/op   | 0 B/op  | 0 allocs/op |
| Alog    | BenchmarkLogger_NewWriter-12   |                    | 2853715    | 409 ns/op   | 0 B/op  | 0 allocs/op |


*/

var (
	testLogMsg     = "Test for a logger with simple string without any formatting"
	testLogMsgFmt  = "Test for a logger with %s, number=%d, float=%f"
	testLogMsgFmt1 = "formats"
	testLogMsgFmt2 = 123
	testLogMsgFmt3 = 3.14
)

func BenchmarkBuiltinLoggerBasic(b *testing.B) {
	// S1: 570 ns/op	      80 B/op	       2 allocs/op
	// S2: 399 ns/op	      80 B/op	       2 allocs/op
	l := log.New(ioutil.Discard, "", log.LstdFlags)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Print(testLogMsg)
		}
	})
}
func BenchmarkBuiltinLoggerFmt(b *testing.B) {
	// S1: 775 ns/op	      88 B/op	       3 allocs/op
	// S2: 502 ns/op	      88 B/op	       3 allocs/op
	l := log.New(ioutil.Discard, "", log.LstdFlags)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Printf(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
		}
	})
}

func BenchmarkLogger_Info(b *testing.B) {
	b.Run("1 eval 0 print", func(b2 *testing.B) {
		// S1: 4.09 ns/op	       0 B/op	       0 allocs/op
		// S2: 1.12 ns/op	       0 B/op	       0 allocs/op
		// S2: 0.406 ns/op	       0 B/op	       0 allocs/op
		l := alog.New(nil).SetFlag(alog.Fdefault).SetPrefix("test ")
		l.SetLevel(alog.Lfatal)
		b2.ReportAllocs()
		b2.ResetTimer()

		b2.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// prints none
				l.Error(testLogMsg)
			}
		})
	})

	b.Run("5 eval 0 prints", func(b2 *testing.B) {
		// S1: 19.4 ns/op	       0 B/op	       0 allocs/op
		// S2: 5.30 ns/op	       0 B/op	       0 allocs/op
		// S2: 1.85 ns/op	       0 B/op	       0 allocs/op
		l := alog.New(nil).SetFlag(alog.Fdefault).SetPrefix("test ")
		l.SetLevel(alog.Lfatal)
		b2.ReportAllocs()
		b2.ResetTimer()
		b2.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// prints none
				l.Trace(testLogMsg)
				l.Debug(testLogMsg)
				l.Info(testLogMsg)
				l.Warn(testLogMsg)
				l.Error(testLogMsg)
			}
		})
	})

	b.Run("5 eval 3 prints", func(b2 *testing.B) {
		// S1: 2013 ns/op	       0 B/op	       0 allocs/op
		// S2: 1119 ns/op	       0 B/op	       0 allocs/op
		// S2:  762 ns/op	       0 B/op	       0 allocs/op
		l := alog.New(nil).SetFlag(alog.Fdefault).SetPrefix("test ")
		b2.ReportAllocs()
		b2.ResetTimer()
		b2.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Trace(testLogMsg)
				l.Debug(testLogMsg)
				l.Info(testLogMsg)  // print
				l.Warn(testLogMsg)  // print
				l.Error(testLogMsg) // print
			}
		})
	})
}

func BenchmarkLogger_Infof(b *testing.B) {
	b.Run("1 eval 0 print", func(b2 *testing.B) {
		// S1: 4.09 ns/op	       0 B/op	       0 allocs/op
		// S2: 1.12 ns/op	       0 B/op	       0 allocs/op
		// S2: 1.05 ns/op	       0 B/op	       0 allocs/op
		l := alog.New(nil).SetFlag(alog.Fdefault).SetPrefix("test ")
		l.SetLevel(alog.Lfatal)
		b2.ReportAllocs()
		b2.ResetTimer()

		b2.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// prints none
				l.Errorf(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
			}
		})
	})

	b.Run("5 eval 0 prints", func(b2 *testing.B) {
		// S1: 19.4 ns/op	       0 B/op	       0 allocs/op
		// S2: 5.30 ns/op	       0 B/op	       0 allocs/op
		// S2: 4.96 ns/op	       0 B/op	       0 allocs/op
		l := alog.New(nil).SetFlag(alog.Fdefault).SetPrefix("test ")
		l.SetLevel(alog.Lfatal)
		b2.ReportAllocs()
		b2.ResetTimer()
		b2.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// prints none
				l.Tracef(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
				l.Debugf(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
				l.Infof(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
				l.Warnf(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
				l.Errorf(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
			}
		})
	})

	b.Run("5 eval 3 prints", func(b2 *testing.B) {
		// S1: 2013 ns/op	       0 B/op	       0 allocs/op
		// S2: 1119 ns/op	       0 B/op	       0 allocs/op
		// S3: 1073 ns/op	       0 B/op	       0 allocs/op
		l := alog.New(nil).SetFlag(alog.Fdefault).SetPrefix("test ")
		b2.ReportAllocs()
		b2.ResetTimer()
		b2.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Tracef(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
				l.Debugf(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
				l.Infof(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)  // print
				l.Warnf(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)  // print
				l.Errorf(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3) // print
			}
		})
	})
}

func BenchmarkLogger_Print(b *testing.B) {
	b.Run("5 eval 1 prints", func(b2 *testing.B) {
		// S1: 382 ns/op	       0 B/op	       0 allocs/op
		// S2: 269 ns/op	       0 B/op	       0 allocs/op
		// S2: 300 ns/op	       0 B/op	       0 allocs/op

		var BACK, FRNT, CAT1, CAT2, CAT3 alog.Tag

		l := alog.New(nil).Do(
			func(l2 *alog.Logger) {
				BACK = l2.NewTag()
				FRNT = l2.NewTag()
				CAT1 = l2.NewTag()
				CAT2 = l2.NewTag()
				CAT3 = l2.NewTag()
			}).SetFlag(alog.Fdefault).SetPrefix("test ")

		l.SetFilter(BACK)

		b2.ReportAllocs()
		b2.ResetTimer()
		b2.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Print(alog.Lwarn, BACK, testLogMsg)
				l.Print(alog.Lwarn, FRNT, testLogMsg)
				l.Print(alog.Lwarn, CAT1, testLogMsg)
				l.Print(alog.Lwarn, CAT2, testLogMsg)
				l.Print(alog.Lwarn, CAT3, testLogMsg)
			}
		})
	})

	b.Run("5 eval 2 prints", func(b2 *testing.B) {
		// S1: 725 ns/op	       0 B/op	       0 allocs/op
		// S2: 482 ns/op	       0 B/op	       0 allocs/op
		// S2: 511 ns/op	       0 B/op	       0 allocs/op

		var BACK, FRNT, CAT1, CAT2, CAT3 alog.Tag

		l := alog.New(nil).SetTags(&BACK, &FRNT, &CAT1, &CAT2, &CAT3).SetFlag(alog.Fdefault).SetPrefix("test ")

		l.SetFilter(BACK | CAT1)

		b2.ReportAllocs()
		b2.ResetTimer()
		b2.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Printf(alog.Lwarn, BACK, testLogMsg)
				l.Printf(alog.Lwarn, FRNT, testLogMsg)
				l.Printf(alog.Lwarn, CAT1, testLogMsg)
				l.Printf(alog.Lwarn, CAT2, testLogMsg)
				l.Printf(alog.Lwarn, CAT3, testLogMsg)
			}
		})
	})
}

func BenchmarkLogger_Printf(b *testing.B) {
	b.Run("5 eval 1 prints", func(b2 *testing.B) {
		// S1: 698 ns/op	       0 B/op	       0 allocs/op
		// S2: 414 ns/op	       0 B/op	       0 allocs/op
		// S2: 380 ns/op	       0 B/op	       0 allocs/op

		var BACK, FRNT, CAT1, CAT2, CAT3 alog.Tag

		l := alog.New(nil).Do(
			func(l2 *alog.Logger) {
				BACK = l2.NewTag()
				FRNT = l2.NewTag()
				CAT1 = l2.NewTag()
				CAT2 = l2.NewTag()
				CAT3 = l2.NewTag()
			}).SetFlag(alog.Fdefault).SetPrefix("test ")

		l.SetFilter(BACK)

		b2.ReportAllocs()
		b2.ResetTimer()
		b2.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Printf(alog.Lwarn, BACK, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
				l.Printf(alog.Lwarn, FRNT, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
				l.Printf(alog.Lwarn, CAT1, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
				l.Printf(alog.Lwarn, CAT2, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
				l.Printf(alog.Lwarn, CAT3, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
			}
		})
	})

	b.Run("5 eval 2 prints", func(b2 *testing.B) {
		// S1: 1393 ns/op	       0 B/op	       0 allocs/op
		// S2:  770 ns/op	       0 B/op	       0 allocs/op

		var BACK, FRNT, CAT1, CAT2, CAT3 alog.Tag

		l := alog.New(nil).Do(
			func(l2 *alog.Logger) {
				BACK = l2.NewTag()
				FRNT = l2.NewTag()
				CAT1 = l2.NewTag()
				CAT2 = l2.NewTag()
				CAT3 = l2.NewTag()
			}).SetFlag(alog.Fdefault).SetPrefix("test ")

		l.SetFilter(BACK | CAT1)

		b2.ReportAllocs()
		b2.ResetTimer()
		b2.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Printf(alog.Lwarn, BACK, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
				l.Printf(alog.Lwarn, FRNT, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
				l.Printf(alog.Lwarn, CAT1, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
				l.Printf(alog.Lwarn, CAT2, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
				l.Printf(alog.Lwarn, CAT3, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
			}
		})
	})

	b.Run("print", func(b2 *testing.B) {
		// S1: 1393 ns/op	       0 B/op	       0 allocs/op
		// S2:  770 ns/op	       0 B/op	       0 allocs/op

		l := alog.New(nil)

		b2.ReportAllocs()
		b2.ResetTimer()
		b2.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Print(alog.Linfo, 0, "abc")
			}
		})
	})
	b.Run("printf with no arg", func(b2 *testing.B) {
		// S1: 1393 ns/op	       0 B/op	       0 allocs/op
		// S2:  770 ns/op	       0 B/op	       0 allocs/op

		l := alog.New(nil)

		b2.ReportAllocs()
		b2.ResetTimer()
		b2.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Printf(alog.Linfo, 0, "abc")
			}
		})
	})
}

// Added as of 0.1.1 update, 12/29/2020
func BenchmarkLogger_NewPrint(b *testing.B) {
	// S2:  293 ns/op	       0 B/op	       0 allocs/op

	var CAT1, CAT2 alog.Tag

	l := alog.New(nil).SetTags(&CAT1, &CAT2).SetFlag(alog.Fdefault).SetPrefix("test ")

	l.SetFilter(CAT1)
	WarnCAT1 := l.NewPrint(alog.Lwarn, CAT1, "CAT1w ")
	WarnCAT2 := l.NewPrint(alog.Lwarn, CAT2, "CAT2w ")
	TraceCAT1 := l.NewPrint(alog.Ltrace, CAT1, "CAT1t ")
	TraceCAT2 := l.NewPrint(alog.Ltrace, CAT2, "CAT2t ")

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			WarnCAT1("warn cat1 test")
			WarnCAT2("warn cat2 test")
			TraceCAT1("trace cat1 test")
			TraceCAT2("trace cat2 test")
		}
	})
}

// Added as of 0.1.2 update, 12/29/2020
func BenchmarkLogger_NewWriter(b *testing.B) {
	// S2: 401 ns/op	       0 B/op	       0 allocs/op
	// S2: 438 ns/op	       0 B/op	       0 allocs/op

	var TEST1, TEST2, TEST3 alog.Tag

	l := alog.New(nil).SetTags(&TEST1, &TEST2, &TEST3).SetFlag(alog.Fprefix | alog.Flevel).SetPrefix("nptest ")

	l.SetFilter(TEST2) // only show TEST2

	wT1D := l.NewWriter(alog.Ldebug, TEST1, "T1D ")
	wT1I := l.NewWriter(alog.Linfo, TEST1, "T1I ")
	wT2D := l.NewWriter(alog.Ldebug, TEST2, "T2D ")
	wT2I := l.NewWriter(alog.Linfo, TEST2, "T2I ")
	wT3D := l.NewWriter(alog.Ldebug, TEST3, "T3D ")
	wT3I := l.NewWriter(alog.Linfo, TEST3, "T3I ")

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = fmt.Fprintf(wT1D, "test: %s fprintf", "T1D") // Not printed
			_, _ = fmt.Fprintf(wT1I, "test: %s fprintf", "T1I") // Not printed
			_, _ = fmt.Fprintf(wT2D, "test: %s fprintf", "T2D") // Not printed
			_, _ = fmt.Fprintf(wT2I, "test: %s fprintf", "T2I") // Printed
			_, _ = fmt.Fprintf(wT3D, "test: %s fprintf", "T3D") // Not printed
			_, _ = fmt.Fprintf(wT3I, "test: %s fprintf", "T3I") // Not printed
		}
	})
}
