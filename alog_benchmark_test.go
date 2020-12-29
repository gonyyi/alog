// (c) 2020 Gon Y Yi. <https://gonyyi.com/copyright.txt>
// Version 2, 12/21/2020

package alog_test

import (
	"github.com/gonyyi/alog"
	"io/ioutil"
	"log"
	"testing"
)

/*
System1 (S1): MacBook Pro (13-inch, 2016, Two Thunderbolt 3 ports)
- MacOS 11.1 (20C69) BigSur
- 2 GHz Dual-Core Intel Core i5
- 8 GB 1867 MHz LPDDR3
- Intel Iris Graphics 540 1536 MB

System2 (S2): MacBook Pro (15-inch, 2018)
- MacOS 10.15.7 Catalina
- 2.9 GHz 6-Core Intel Core i9
- 32 GB 2400 MHz DDR4
- Radeon Pro 560X 4 GB / Intel UHD Graphics 630 1536 MB


Baseline, Go's Standard Error (S2, Print, Printf)
-----------------------------------
BenchmarkBuiltinLoggerBasic-12    	 2883834	       410 ns/op	      80 B/op	       2 allocs/op
BenchmarkBuiltinLoggerFmt-12      	 2396258	       508 ns/op	      88 B/op	       3 allocs/op


ALog (S2)
-----------------------------------
BenchmarkAlogPrintf-12            	 3314240	       361 ns/op	       0 B/op	       0 allocs/op
BenchmarkAlogPrint-12             	 4725747	       250 ns/op	       0 B/op	       0 allocs/op
BenchmarkAlogInfof-12             	 3195420	       365 ns/op	       0 B/op	       0 allocs/op
BenchmarkAlogInfo-12              	 4710241	       258 ns/op	       0 B/op	       0 allocs/op
Benchmark_ALog_Level_3-12         	 1482876	       808 ns/op	       0 B/op	       0 allocs/op
Benchmark_ALog_Levelf_3-12        	  974910	      1191 ns/op	       0 B/op	       0 allocs/op
Benchmark_ALog_Levelf5_0-12       	211717482	      5.68 ns/op	       0 B/op	       0 allocs/op
Benchmark_ALog_Levelf1_0-12       	1000000000	      1.19 ns/op	       0 B/op	       0 allocs/op
Benchmark_ALogPrint_Cat5_1-12     	 4502719	       263 ns/op	       0 B/op	       0 allocs/op
Benchmark_ALogPrint_Cat5_2-12     	 2311810	       518 ns/op	       0 B/op	       0 allocs/op
Benchmark_ALogPrintf_Cat5_1-12    	 2869698	       403 ns/op	       0 B/op	       0 allocs/op
Benchmark_ALogPrintf_Cat5_2-12    	 1564646	       774 ns/op	       0 B/op	       0 allocs/op

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
func BenchmarkAlogPrintf(b *testing.B) {
	// S1: 667 ns/op	       0 B/op	       0 allocs/op
	// S2: 385 ns/op	       0 B/op	       0 allocs/op
	l := alog.New(ioutil.Discard, "", alog.Fdefault)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Printf(alog.Linfo, 0, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
		}
	})
}
func BenchmarkAlogPrint(b *testing.B) {
	// S1: 359 ns/op	       0 B/op	       0 allocs/op
	// S2: 249 ns/op	       0 B/op	       0 allocs/op
	l := alog.New(ioutil.Discard, "", alog.Fdefault)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Print(alog.Linfo, 0, testLogMsg)
		}
	})
}
func BenchmarkAlogInfof(b *testing.B) {
	// S1: 651 ns/op	       0 B/op	       0 allocs/op
	// S2: 373 ns/op	       0 B/op	       0 allocs/op
	l := alog.New(ioutil.Discard, "", alog.Fdefault)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Infof(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
		}
	})
}
func BenchmarkAlogInfo(b *testing.B) {
	// S1: 360 ns/op	       0 B/op	       0 allocs/op
	// S2: 246 ns/op	       0 B/op	       0 allocs/op
	l := alog.New(ioutil.Discard, "", alog.Fdefault)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Info(testLogMsg)
		}
	})
}

func Benchmark_ALog_Level_3(b *testing.B) {
	// S1: 1121 ns/op	       0 B/op	       0 allocs/op
	// S2:  762 ns/op	       0 B/op	       0 allocs/op
	l := alog.New(ioutil.Discard, "test ", alog.Fdefault)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Trace(testLogMsg)
			l.Debug(testLogMsg)
			l.Info(testLogMsg)  // print
			l.Warn(testLogMsg)  // print
			l.Error(testLogMsg) // print
		}
	})
}

func Benchmark_ALog_Levelf_3(b *testing.B) {
	// S1: 2013 ns/op	       0 B/op	       0 allocs/op
	// S2: 1119 ns/op	       0 B/op	       0 allocs/op
	l := alog.New(ioutil.Discard, "test ", alog.Fdefault)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Tracef(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
			l.Debugf(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
			l.Infof(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)  // print
			l.Warnf(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)  // print
			l.Errorf(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3) // print
		}
	})
}
func Benchmark_ALog_Levelf5_0(b *testing.B) {
	// S1: 19.4 ns/op	       0 B/op	       0 allocs/op
	// S2: 5.30 ns/op	       0 B/op	       0 allocs/op
	l := alog.New(ioutil.Discard, "test ", alog.Fdefault)
	l.SetLevel(alog.Lfatal)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// prints none
			l.Tracef(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
			l.Debugf(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
			l.Infof(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
			l.Warnf(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
			l.Errorf(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
		}
	})

}
func Benchmark_ALog_Levelf1_0(b *testing.B) {
	// S1: 4.09 ns/op	       0 B/op	       0 allocs/op
	// S2: 1.12 ns/op	       0 B/op	       0 allocs/op
	l := alog.New(ioutil.Discard, "test ", alog.Fdefault)
	l.SetLevel(alog.Lfatal)
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// prints none
			l.Errorf(testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
		}
	})
}
func Benchmark_ALogPrint_Cat5_1(b *testing.B) {
	// S1: 382 ns/op	       0 B/op	       0 allocs/op
	// S2: 269 ns/op	       0 B/op	       0 allocs/op
	l := alog.New(ioutil.Discard, "test ", alog.Fdefault)

	cat := alog.NewCategory()
	BACK := cat.Add()
	FRNT := cat.Add()
	CAT1 := cat.Add()
	CAT2 := cat.Add()
	CAT3 := cat.Add()
	l.SetCategory(BACK)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Print(alog.Lwarn, BACK, testLogMsg)
			l.Print(alog.Lwarn, FRNT, testLogMsg)
			l.Print(alog.Lwarn, CAT1, testLogMsg)
			l.Print(alog.Lwarn, CAT2, testLogMsg)
			l.Print(alog.Lwarn, CAT3, testLogMsg)
		}
	})
}
func Benchmark_ALogPrint_Cat5_2(b *testing.B) {
	// S1: 725 ns/op	       0 B/op	       0 allocs/op
	// S2: 482 ns/op	       0 B/op	       0 allocs/op
	l := alog.New(ioutil.Discard, "test ", alog.Fdefault)

	cat := alog.NewCategory()
	BACK := cat.Add()
	FRNT := cat.Add()
	CAT1 := cat.Add()
	CAT2 := cat.Add()
	CAT3 := cat.Add()

	l.SetCategory(BACK | CAT2)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Print(alog.Lwarn, BACK, testLogMsg)
			l.Print(alog.Lwarn, FRNT, testLogMsg)
			l.Print(alog.Lwarn, CAT1, testLogMsg)
			l.Print(alog.Lwarn, CAT2, testLogMsg)
			l.Print(alog.Lwarn, CAT3, testLogMsg)
		}
	})
}
func Benchmark_ALogPrintf_Cat5_1(b *testing.B) {
	// S1: 698 ns/op	       0 B/op	       0 allocs/op
	// S2: 414 ns/op	       0 B/op	       0 allocs/op
	l := alog.New(ioutil.Discard, "test ", alog.Fdefault)

	cat := alog.NewCategory()
	BACK := cat.Add()
	FRNT := cat.Add()
	CAT1 := cat.Add()
	CAT2 := cat.Add()
	CAT3 := cat.Add()
	l.SetCategory(BACK)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Printf(alog.Lwarn, BACK, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
			l.Printf(alog.Lwarn, FRNT, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
			l.Printf(alog.Lwarn, CAT1, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
			l.Printf(alog.Lwarn, CAT2, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
			l.Printf(alog.Lwarn, CAT3, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
		}
	})
}
func Benchmark_ALogPrintf_Cat5_2(b *testing.B) {
	// S1: 1393 ns/op	       0 B/op	       0 allocs/op
	// S2:  770 ns/op	       0 B/op	       0 allocs/op
	l := alog.New(ioutil.Discard, "test ", alog.Fdefault)

	cat := alog.NewCategory()
	BACK := cat.Add()
	FRNT := cat.Add()
	CAT1 := cat.Add()
	CAT2 := cat.Add()
	CAT3 := cat.Add()
	l.SetCategory(BACK | CAT1)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Printf(alog.Lwarn, BACK, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
			l.Printf(alog.Lwarn, FRNT, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
			l.Printf(alog.Lwarn, CAT1, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
			l.Printf(alog.Lwarn, CAT2, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
			l.Printf(alog.Lwarn, CAT3, testLogMsgFmt, testLogMsgFmt1, testLogMsgFmt2, testLogMsgFmt3)
		}
	})
}
// Added as of 0.1.1 update
func Benchmark_ALog_NewPrint(b *testing.B) {
	// S2:  293 ns/op	       0 B/op	       0 allocs/op
	l := alog.New(nil, "test ", alog.Fdefault)
	cat := alog.NewCategory()
	CAT1 := cat.Add()
	CAT2 := cat.Add()
	l.SetCategory(CAT1)
	WarnCAT1 := l.NewPrint(alog.Lwarn, CAT1)
	WarnCAT2 := l.NewPrint(alog.Lwarn, CAT2)
	TraceCAT1 := l.NewPrint(alog.Ltrace, CAT1)
	TraceCAT2 := l.NewPrint(alog.Ltrace, CAT2)

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

