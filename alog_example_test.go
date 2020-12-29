// (c) 2020 Gon Y Yi. <https://gonyyi.com/copyright.txt>
// Version 0.1.3, 12/29/2020

package alog_test

import (
	"fmt"
	"github.com/gonyyi/alog"
	"os"
)

func ExampleNew() {
	l := alog.New(os.Stdout, "test ", alog.Fprefix|alog.Flevel)

	// Trace/Debug will NOT be printed
	l.Trace("hello trace")
	l.Debug("hello debug")

	// Info/Warn/Error/Fatal will be printed
	l.Info("hello info")
	l.Warn("hello warn")
	l.Error("hello error")

	// Output:
	// test [INF] hello info
	// test [WRN] hello warn
	// test [ERR] hello error
}

func ExampleLogger_NewWriter() {
	l := alog.New(os.Stdout, "nptest ", alog.Fprefix|alog.Flevel) // Default level is INFO and higher

	l.SetLevel(alog.Ldebug) // set logging level to DEBUG

	cat := alog.NewCategory()
	TEST1 := cat.Add()
	TEST2 := cat.Add()
	TEST3 := cat.Add()

	// only show TEST2 here.
	// Therefore only DEBUG/INFO with TEST2 will be printed
	l.SetCategory(TEST2)

	wT1D := l.NewWriter(alog.Ldebug, TEST1, "T1D ")
	wT1I := l.NewWriter(alog.Linfo, TEST1, "T1I ")
	wT2D := l.NewWriter(alog.Ldebug, TEST2, "T2D ")
	wT2I := l.NewWriter(alog.Linfo, TEST2, "T2I ")
	wT3D := l.NewWriter(alog.Ldebug, TEST3, "T3D ")
	wT3I := l.NewWriter(alog.Linfo, TEST3, "T3I ")

	fmt.Fprintf(wT1D, "test: %s fprintf", "T1D")
	fmt.Fprintf(wT1I, "test: %s fprintf", "T1I")
	fmt.Fprintf(wT2D, "test: %s fprintf", "T2D")
	fmt.Fprintf(wT2I, "test: %s fprintf", "T2I")
	fmt.Fprintf(wT3D, "test: %s fprintf", "T3D")
	fmt.Fprintf(wT3I, "test: %s fprintf", "T3I")

	// Output:
	// nptest [DBG] T2D test: T2D fprintf
	// nptest [INF] T2I test: T2I fprintf
}

func ExampleLogger_NewPrint() {
	l := alog.New(os.Stdout, "nptest ", alog.Fprefix|alog.Flevel) // Default level is INFO and higher

	cat := alog.NewCategory()
	CAT1 := cat.Add()
	CAT2 := cat.Add()

	l.SetCategory(CAT1) // Print only CAT1
	WarnCAT1 := l.NewPrint(alog.Lwarn, CAT1, "CAT1W ")
	WarnCAT2 := l.NewPrint(alog.Lwarn, CAT2, "CAT2W ")
	TraceCAT1 := l.NewPrint(alog.Ltrace, CAT1, "CAT1T ")
	TraceCAT2 := l.NewPrint(alog.Ltrace, CAT2, "CAT2T ")

	// Since category is set to CAT1, and default level is INFO,
	// only item(s) with CAT1 and INFO and above will be printed.
	WarnCAT1("warn cat1 test")
	WarnCAT2("warn cat2 test")
	TraceCAT1("trace cat1 test")
	TraceCAT2("trace cat2 test")

	// Output:
	// nptest [WRN] CAT1W warn cat1 test
}

func ExampleLogger_SetOutput() {
	// When nil is used for output, it will use ioutil.Discard.
	// Therefore, below will not be printed because it will output to ioutil.Discard
	l := alog.New(nil, "test ", alog.Flevel|alog.Fprefix)
	l.SetLevel(alog.Ltrace)

	l.Info("Hello 1")

	// this will be printed as standard error
	l.SetOutput(os.Stdout)
	l.Info("Hello 2")

	// This will not be printed as new output is now ioutil.Discard again.
	l.SetOutput(nil)
	l.Info("Hello 3")

	// Output:
	// test [INF] Hello 2
}

func ExampleLogger_SetPrefix() {
	// When nil is used for output, it will use ioutil.Discard.
	// Therefore, below will not be printed because it will output to ioutil.Discard
	l := alog.New(os.Stdout, "test1 ", alog.Flevel|alog.Fprefix)
	l.Info("Hello 1")
	l.SetPrefix("test2 ") // change prefix to "test2 " (with trailing space)
	l.Info("Hello 2")

	// Output:
	// test1 [INF] Hello 1
	// test2 [INF] Hello 2
}

func ExampleLogger_SetFlag() {
	l := alog.New(os.Stdout, "test1 ", 0) // 0 is equal to no flag

	// As no flag is given, below will print the strings given as parameters.
	// info test 1
	// info test 2
	l.Info("info test 1")
	l.Info("info test 2")

	// Setting flag to show level and prefix
	// Next log will have the prefix and log level populated
	l.SetFlag(alog.Flevel | alog.Fprefix)

	l.Info("info test 3")
	l.Info("info test 4")

	// Output:
	// info test 1
	// info test 2
	// test1 [INF] info test 3
	// test1 [INF] info test 4
}

func ExampleLogger_SetLevelPrefix() {
	// By default, level prefixes are:
	// 	  Trace: "[TRC] "
	// 	  Debug: "[DBG] "
	// 	  Info:  "[INF] "
	// 	  Warn:  "[WRN] "
	// 	  Error: "[ERR] "
	// 	  Fatal: "[FTL] "
	// but, this can be overwritten using SetLevelPrefix()
	l := alog.New(os.Stdout, "test ", alog.Flevel|alog.Fprefix)
	l.Trace("1 trace")
	l.Debug("1 debug")
	l.Info("1 info")
	l.Warn("1 warn")
	l.Error("1 error")
	// if logger.Fatal() is called, it will exit with 1,
	// so for testing, using Print() instead.
	l.Print(alog.Lfatal, 0, "1 fatal")

	// Not switch level prefix:
	l.SetLevelPrefix("[TRACE] ", "[DEBUG] ", "[INFO] ", "[WARN] ", "[ERROR] ", "[FATAL] ")
	l.Trace("2 trace")
	l.Debug("2 debug")
	l.Info("2 info")
	l.Warn("2 warn")
	l.Error("2 error")
	l.Print(alog.Lfatal, 0, "2 fatal")

	// Output:
	// test [INF] 1 info
	// test [WRN] 1 warn
	// test [ERR] 1 error
	// test [FTL] 1 fatal
	// test [INFO] 2 info
	// test [WARN] 2 warn
	// test [ERROR] 2 error
	// test [FATAL] 2 fatal
}

func ExampleLogger_SetCategory() {
	l := alog.New(os.Stdout, "test ", alog.Flevel|alog.Fprefix)

	// Assume there are 4 categories: SYSTEM, DISK, REQUEST, RESPONSE
	mycat := alog.NewCategory()
	SYSTEM := mycat.Add()
	DISK := mycat.Add()
	REQUEST := mycat.Add()
	RESPONSE := mycat.Add()

	// And the user only wants to see REQUEST and RESPONSE logs with INFO or above level.
	l.SetLevel(alog.Linfo)
	l.SetCategory(REQUEST | RESPONSE)

	// Debug logs won't be printed
	l.Print(alog.Ldebug, SYSTEM, "1 debug + system")
	l.Print(alog.Ldebug, DISK, "1 debug + disk")
	l.Print(alog.Ldebug, REQUEST, "1 debug + request")
	l.Print(alog.Ldebug, RESPONSE, "1 debug + response")

	// INFO will be printed, however, SYSTEM and DISK category won't.
	l.Print(alog.Linfo, SYSTEM, "1 info + system")
	l.Print(alog.Linfo, DISK, "1 info + disk")

	// Below two will be only log items that will be printed
	l.Print(alog.Linfo, REQUEST, "1 info + request")
	l.Print(alog.Linfo, RESPONSE, "1 info + response")

	// Output:
	// test [INF] 1 info + request
	// test [INF] 1 info + response
}

func ExampleLogger_SetLevel() {
	// Create a new logger with default level WARN (alog.CLevelWarn)
	l := alog.New(os.Stdout, "test ", alog.Flevel|alog.Fprefix)
	l.SetLevel(alog.Lwarn)

	// This will NOT be printed because of log level is WARN
	l.Info("test info 1")

	// Override the log level config to INFO (Linfo)
	l.SetLevel(alog.Linfo)

	// This WILL BE printed.
	l.Info("test info 2")

	// Output: test [INF] test info 2
}

func ExampleLogger_Print() {
	l := alog.New(os.Stdout, "test ", alog.Flevel|alog.Fprefix)

	// this will print: "test [INF] hello"
	l.Print(alog.Linfo, 0, "hello1")

	// Output:
	// test [INF] hello1
}

func ExampleLogger_Printf() {
	l := alog.New(os.Stdout, "test ", alog.Flevel|alog.Fprefix)
	// There can be slightly different parsing when it's float32 vs float64 due to rounding.
	// Also, float will be always at 2nd decimal point.
	l.Printf(alog.Linfo, 0, "hello %s, your ID is %d, you are %f", "JON", 124, 5.8000)
	l.Printf(alog.Linfo, 0, "hello %s, your ID is %d, you are %f", "GON", 123, float32(5.8000))
	// Output:
	// test [INF] hello JON, your ID is 124, you are 5.79
	// test [INF] hello GON, your ID is 123, you are 5.80
}

func ExampleLevel() {
	l := alog.New(os.Stdout, "test ", alog.Flevel|alog.Fprefix)
	l.Debugf("%s=%s", "level", "debug") // default logging level is INFO. Hence, this won't bee printed.
	l.Infof("%s=%s", "level", "info")
	l.Errorf("%s=%s", "level", "error")

	// Output:
	// test [INF] level=info
	// test [ERR] level=error
}

func ExampleNewCategory() {
	cat := alog.NewCategory()
	BACK := cat.Add()
	FRONT := cat.Add()
	USER := cat.Add()
	l := alog.New(os.Stdout, "test ", alog.Flevel|alog.Fprefix)
	l.SetLevel(alog.Ltrace)

	// Assume I want to see BACK and FRONT with a level INFO or above.
	l.SetCategory(BACK | FRONT) // only show BACK and FRONT
	l.SetLevel(alog.Linfo)      // this will override config "alog.CLevelTrace"

	f := func(c alog.Category, s string) {
		l.Printf(alog.Ltrace, c, "%s.trace", s)
		l.Printf(alog.Ldebug, c, "%s.debug", s)
		l.Printf(alog.Linfo, c, "%s.info", s)
		l.Printf(alog.Lwarn, c, "%s.warn", s)
		l.Printf(alog.Lerror, c, "%s.error", s)
	}

	f(BACK, "BACK")
	f(FRONT, "FRONT")
	f(USER, "USER")

	// Output:
	// test [INF] BACK.info
	// test [WRN] BACK.warn
	// test [ERR] BACK.error
	// test [INF] FRONT.info
	// test [WRN] FRONT.warn
	// test [ERR] FRONT.error
}
