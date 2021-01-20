// (c) 2020 Gon Y Yi. <https://gonyyi.com>
// Version 0.1.3, 12/29/2020

package alog_test

//
// import (
// 	"errors"
// 	"fmt"
// 	"github.com/gonyyi/alog"
// 	"os"
// )
//
// func ExampleNew() {
// 	// Default alog will record date (YYYYMMDD) and time.
// 	// So disable date and time, and only show level for output comparison
// 	l := alog.New(os.Stdout).SetFormat(alog.Flevel)
//
// 	// Trace/Debug will NOT be printed because default level is Info.
// 	l.Trace("hello trace")
// 	l.Debug("hello debug")
//
// 	// Info/Warn/Error/Fatal will be printed
// 	l.Info("hello info")
// 	l.Warn("hello warn")
// 	l.Error("hello error")
//
// 	// Output:
// 	// [INF] hello info
// 	// [WRN] hello warn
// 	// [ERR] hello error
// }
//
// func ExampleLogger_Do() {
// 	// Do func act as a script configuration.
//
// 	// myconf1 will set the prefix, and set the logging level to DEBUG
// 	// only shows prefix and level.
// 	myconf1 := func(l *alog.Logger) {
// 		l.SetPrefix("log ")
// 		l.SetLogLevel(alog.Ldebug).SetFormat(alog.Fprefix | alog.Flevel)
// 	}
//
// 	// myconf2 will set logging level prefixes. Instead of default 3 char long,
// 	// this will use full character
// 	myconf2 := func(l *alog.Logger) {
// 		l.setLevelPrefix(
// 			"[TRACE] ",
// 			"[DEBUG] ",
// 			"[INFO]  ",
// 			"[WARN]  ",
// 			"[ERROR] ",
// 			"[FATAL] ")
// 	}
//
// 	// Do function can take zero to many do functions - func(*alog.Logger).
// 	l := alog.New(os.Stdout).Do(myconf1, myconf2)
//
// 	l.Print(alog.Ltrace, 0, "testTrace")
// 	l.Print(alog.Ldebug, 0, "testDebug")
// 	l.Print(alog.Linfo, 0, "testInfo")
// 	l.Print(alog.Lwarn, 0, "testWarn")
// 	l.Print(alog.Lerror, 0, "testError")
//
// 	// Output:
// 	// log [DEBUG] testDebug
// 	// log [INFO]  testInfo
// 	// log [WARN]  testWarn
// 	// log [ERROR] testError
// }
//
// func ExampleLogger_SetOutput() {
// 	// When nil is used for output, it will use ioutil.Discard.
// 	// Therefore, below will not be printed because it will output to ioutil.Discard
//
// 	// v0.1.6 Code:
// 	// l := alog.New(nil, "test ", alog.Flevel|alog.Fprefix)
// 	// l.SetLogLevel(alog.Ltrace)
//
// 	// v0.2.0 Code:
// 	l := alog.New(nil).SetPrefix("test ").SetFormat(alog.Flevel | alog.Fprefix).SetLogLevel(alog.Ltrace)
//
// 	l.Info("Hello 1")
//
// 	// this will be printed as standard error
// 	l.SetOutput(os.Stdout)
// 	l.Info("Hello 2")
//
// 	// This will not be printed as new output is now ioutil.Discard again.
// 	l.SetOutput(nil)
// 	l.Info("Hello 3")
//
// 	// Output:
// 	// test [INF] Hello 2
// }
//
// func ExampleLogger_SetPrefix() {
// 	// When nil is used for output, it will use ioutil.Discard.
// 	// Therefore, below will not be printed because it will output to ioutil.Discard
//
// 	// v0.1.6 Code:
// 	// l := alog.New(os.Stdout, "test1 ", alog.Flevel|alog.Fprefix)
//
// 	// v0.2.0 Code:
// 	l := alog.New(os.Stdout).SetPrefix("test1 ").SetFormat(alog.Flevel | alog.Fprefix)
//
// 	l.Info("Hello 1")
// 	l.SetPrefix("test2 ") // change prefix to "test2 " (with trailing space)
// 	l.Info("Hello 2")
//
// 	// Output:
// 	// test1 [INF] Hello 1
// 	// test2 [INF] Hello 2
// }
//
// func ExampleLogger_SetFlag() {
// 	// v0.1.6 Code:
// 	// l := alog.New(os.Stdout, "test1 ", 0) // 0 is equal to no flag
//
// 	// v0.2.0 Code:
// 	// A flag `alog.Fnone` will reset flags.
// 	l := alog.New(os.Stdout).SetPrefix("test1 ").SetFormat(alog.Fnone)
//
// 	// As no flag is given, below will print the strings given as parameters.
// 	// info test 1
// 	// info test 2
// 	l.Info("info test 1")
// 	l.Info("info test 2")
//
// 	// Setting flag to show Level and prefix
// 	// Next log will have the prefix and log Level populated
// 	l.SetFormat(alog.Flevel | alog.Fprefix)
//
// 	l.Info("info test 3")
// 	l.Info("info test 4")
//
// 	// Output:
// 	// info test 1
// 	// info test 2
// 	// test1 [INF] info test 3
// 	// test1 [INF] info test 4
// }
//
// func ExampleLogger_SetLevelPrefix() {
// 	// By default, Level prefixes are:
// 	// 	  Trace: "[TRC] "
// 	// 	  Debug: "[DBG] "
// 	// 	  Info:  "[INF] "
// 	// 	  Warn:  "[WRN] "
// 	// 	  Error: "[ERR] "
// 	// 	  Fatal: "[FTL] "
// 	// but, this can be overwritten using setLevelPrefix()
//
// 	// v0.1.6 Code:
// 	// l := alog.New(os.Stdout, "test ", alog.Flevel|alog.Fprefix)
//
// 	// v0.2.0 Code:
// 	l := alog.New(os.Stdout).SetPrefix("test ").SetFormat(alog.Flevel | alog.Fprefix)
//
// 	l.Trace("1 trace")
// 	l.Debug("1 debug")
// 	l.Info("1 info")
// 	l.Warn("1 warn")
// 	l.Error("1 error")
// 	// if logger.Fatal() is called, it will exit with 1,
// 	// so for testing, using Print() instead.
// 	l.Print(alog.Lfatal, 0, "1 fatal")
//
// 	// Not switch Level prefix:
// 	l.setLevelPrefix("[TRACE] ", "[DEBUG] ", "[INFO] ", "[WARN] ", "[ERROR] ", "[FATAL] ")
// 	l.Trace("2 trace")
// 	l.Debug("2 debug")
// 	l.Info("2 info")
// 	l.Warn("2 warn")
// 	l.Error("2 error")
// 	l.Print(alog.Lfatal, 0, "2 fatal")
//
// 	// Output:
// 	// test [INF] 1 info
// 	// test [WRN] 1 warn
// 	// test [ERR] 1 error
// 	// test [FTL] 1 fatal
// 	// test [INFO] 2 info
// 	// test [WARN] 2 warn
// 	// test [ERROR] 2 error
// 	// test [FATAL] 2 fatal
// }
//
// func ExampleLogger_SetLogLevel() {
// 	// Create a new logger with default Level WARN
//
// 	l := alog.New(os.Stdout).SetPrefix("test ").SetFormat(alog.Flevel | alog.Fprefix).SetLogLevel(alog.Lwarn)
//
// 	// This will NOT be printed because of log Level is WARN
// 	l.Info("test info 1")
//
// 	// Override the log Level config to INFO (Linfo)
// 	l.SetLogLevel(alog.Linfo)
//
// 	// This WILL BE printed.
// 	l.Info("test info 2")
//
// 	// Output:
// 	// test [INF] test info 2
// }
// func ExampleLogger_SetLogTag() {
// 	// Crate a logger
// 	l := alog.New(os.Stdout).SetFormat(alog.Flevel)
//
// 	// Create 4 tags: SYSTEM, DISK, REQUEST, RESPONSE
// 	SYSTEM := l.NewTag("system")
// 	DISK := l.NewTag("disk")
// 	REQUEST := l.NewTag("request")
// 	RESPONSE := l.NewTag("response")
//
// 	// SetLogTag will only logging tags that meets the tag flag.
// 	// In this example, only REQUEST and RESPONSE tag will be shown.
// 	l.SetLogTag(REQUEST | RESPONSE)
//
// 	// Default log level is INFO. Anything equal to above INFO
// 	// and tag matches REQUEST or RESPONSE will be printed.
// 	l.Print(alog.Linfo, SYSTEM, "1 info + system")
// 	l.Print(alog.Linfo, DISK, "1 info + disk")
// 	l.Print(alog.Linfo, REQUEST, "1 info + request")
// 	l.Print(alog.Linfo, RESPONSE, "1 info + response")
//
// 	// Output:
// 	// [INF] 1 info + request
// 	// [INF] 1 info + response
// }
//
// func ExampleLogger_SetLogFn() {
// 	l := alog.New(os.Stdout).SetFormat(alog.Flevel)
//
// 	// Create tags and enable them by UseTag method.
// 	SYSTEM, DISK, REQUEST, RESPONSE := l.NewTag("system"),l.NewTag("disk"),l.NewTag("request"),l.NewTag("response")
//
// 	// Instead of SetLogTag and/or SetLogLevel method,
// 	// a function will be used for SetLogFn method as below.
// 	l.SetLogFn(func(level alog.Level, tag alog.Tag) bool {
// 		// Log anything with log level INFO or above,
// 		// OR, log REQUEST (regardless of level)
// 		if level >= alog.Linfo || tag&REQUEST != 0 {
// 			return true
// 		}
// 		return false
// 	})
//
// 	// Debug logs won't be printed
// 	l.Print(alog.Ldebug, SYSTEM, "DEBUG + SYSTEM")
// 	l.Print(alog.Ldebug, DISK, "DEBUG + DISK")
// 	l.Print(alog.Ldebug, REQUEST, "DEBUG + REQ")
// 	l.Print(alog.Ldebug, RESPONSE, "DEBUG + RESP")
// 	l.Print(alog.Linfo, SYSTEM, "INFO + SYSTEM")
// 	l.Print(alog.Linfo, DISK, "INFO + DISK")
// 	l.Print(alog.Linfo, REQUEST, "INFO + REQ")
// 	l.Print(alog.Linfo, RESPONSE, "INFO + RESP")
//
// 	// Output:
// 	// [DBG] DEBUG + REQ
// 	// [INF] INFO + SYSTEM
// 	// [INF] INFO + DISK
// 	// [INF] INFO + REQ
// 	// [INF] INFO + RESP
// }
//
// func ExampleLogger_Output() {
// 	l := alog.New(os.Stdout).SetPrefix("test ").SetFormat(alog.Flevel | alog.Fprefix)
//
// 	// this will print: "test [INF] hello"
// 	l.Output(alog.Linfo, 0, []byte("hello1"))
//
// 	// Output:
// 	// test [INF] hello1
// }
//
// func ExampleLogger_Print() {
// 	l := alog.New(os.Stdout).SetPrefix("test ").SetFormat(alog.Flevel | alog.Fprefix)
//
// 	// this will print: "test [INF] hello"
// 	l.Print(alog.Linfo, 0, "hello1")
//
// 	// Output:
// 	// test [INF] hello1
// }
//
// func ExampleLogger_Printf() {
// 	l := alog.New(os.Stdout).SetPrefix("test ").SetFormat(alog.Flevel | alog.Fprefix)
//
// 	// There can be slightly different parsing when it's float32 vs float64 due to rounding.
// 	// Also, float will be always at 2nd decimal point.
// 	l.Printf(alog.Linfo, 0, "hello %s, your ID is %d, you are %f", "JON", 124, 5.8000)
// 	l.Printf(alog.Linfo, 0, "hello %s, your ID is %d, you are %f", "GON", 123, float32(5.8000))
//
// 	// Output:
// 	// test [INF] hello JON, your ID is 124, you are 5.79
// 	// test [INF] hello GON, your ID is 123, you are 5.80
// }
// func ExampleLogger_Trace() {
// 	// Debug will only print IF log level is Debug or Trace.
// 	// However default is set to INFO.
// 	l := alog.New(os.Stdout).SetFormat(alog.Flevel)
//
// 	// This won't be printed.
// 	l.Trace("this will not be printed")
//
// 	// By changing the log level, next debug message will be printed
// 	l.SetLogLevel(alog.Ltrace)
// 	l.Trace("this will be printed")
//
// 	// Output:
// 	// [TRC] this will be printed
// }
// func ExampleLogger_Tracef() {
// 	l := alog.New(os.Stdout).SetFormat(alog.Flevel | alog.Fprefix)
// 	l.SetLogLevel(alog.Ltrace)
// 	l.Tracef("%s=%s", "Level", "trace")
//
// 	// Output:
// 	// [TRC] Level=trace
// }
// func ExampleLogger_Debug() {
// 	// Debug will only print IF log level is Debug or Trace.
// 	// However default is set to INFO.
// 	l := alog.New(os.Stdout).SetFormat(alog.Flevel)
//
// 	// This won't be printed.
// 	l.Debug("this will not be printed")
//
// 	// By changing the log level, next debug message will be printed
// 	l.SetLogLevel(alog.Ldebug)
// 	l.Debug("this will be printed")
//
// 	// Output:
// 	// [DBG] this will be printed
// }
// func ExampleLogger_Debugf() {
// 	l := alog.New(os.Stdout).SetFormat(alog.Flevel | alog.Fprefix)
// 	l.SetLogLevel(alog.Ldebug)
// 	l.Debugf("%s=%s", "Level", "debug")
//
// 	// Output:
// 	// [DBG] Level=debug
// }
// func ExampleLogger_Info() {
// 	// Debug will only print IF log level is Debug or Trace.
// 	// However default is set to INFO.
// 	l := alog.New(os.Stdout).SetFormat(alog.Flevel)
//
// 	// This will be printed.
// 	l.Info("this will be printed")
//
// 	// Output:
// 	// [INF] this will be printed
// }
// func ExampleLogger_Infof() {
// 	l := alog.New(os.Stdout).SetFormat(alog.Flevel)
// 	l.Infof("%s=%s", "Level", "info")
//
// 	// Output:
// 	// [INF] Level=info
// }
//
// func ExampleLogger_Warn() {
// 	// Debug will only print IF log level is Debug or Trace.
// 	// However default is set to INFO.
// 	l := alog.New(os.Stdout).SetFormat(alog.Flevel)
//
// 	// This won't be printed.
// 	l.Warn("this will be printed")
//
// 	// Output:
// 	// [WRN] this will be printed
// }
// func ExampleLogger_Warnf() {
// 	l := alog.New(os.Stdout).SetFormat(alog.Flevel | alog.Fprefix)
// 	l.Warnf("%s=%s", "Level", "warn")
//
// 	// Output:
// 	// [WRN] Level=warn
// }
//
// func ExampleLogger_Error() {
// 	// Debug will only print IF log level is Debug or Trace.
// 	// However default is set to INFO.
// 	l := alog.New(os.Stdout).SetFormat(alog.Flevel)
// 	l.Error("this will be printed")
//
// 	// Output:
// 	// [ERR] this will be printed
// }
// func ExampleLogger_Errorf() {
// 	l := alog.New(os.Stdout).SetFormat(alog.Flevel | alog.Fprefix)
// 	l.Errorf("%s=%s", "Level", "error")
//
// 	// Output:
// 	// [ERR] Level=error
// }
//
// func ExampleLogger_IfError() {
// 	l := alog.New(os.Stdout).SetFormat(alog.Flevel | alog.Fprefix)
//
// 	var myErr error
//
// 	// When myErr was created, it should be nil at this point.
// 	// Therefore, nothing will be logged.
// 	l.LogIferr(myErr)
//
// 	// Now setting myErr with an error,
// 	// This will be logged.
// 	myErr = errors.New("my test error")
// 	l.LogIferr(myErr)
//
// 	// Output:
// 	// [ERR] my test error
// }
//
// func ExampleLogger_IfFatal() {
// 	l := alog.New(os.Stdout).SetFormat(alog.Flevel | alog.Fprefix)
//
// 	var myErr error
//
// 	// When myErr was created, it should be nil at this point.
// 	// Therefore, nothing will be logged. But, if actual error is there,
// 	// this will print the error and exit the code.
// 	l.IfFatal(myErr)
//
// 	// Output:
// 	//
// }
//
// func ExampleLogger_Writer() {
// 	l := alog.New(os.Stdout).SetFormat(alog.Flevel)
// 	out := l.GetWriter() // this will be os.Stdout
//
// 	l.Info("info log")
// 	out.Write([]byte("writing test using io.GetWriter"))
//
// 	// Output:
// 	// [INF] info log
// 	// writing test using io.GetWriter
// }
// func ExampleLogger_Close() {
// 	l := alog.New(os.Stdout).SetFormat(alog.Flevel)
// 	l.Info("info log")
//
// 	// if io.GetWriter used to create the log was a file,
// 	// this will close the file
// 	l.Close()
//
// 	// Output:
// 	// [INF] info log
// }
// func ExampleLogger_NewPrint() {
// 	var CAT1, CAT2 alog.Tag
//
// 	// v0.1.6 Code:
// 	// l := alog.New(os.Stdout, "nptest ", alog.Fprefix|alog.Flevel) // Default Level is INFO and higher
//
// 	// v0.2.0 Code:
// 	l := alog.New(os.Stdout).Do(
// 		func(l2 *alog.Logger) {
// 			CAT1 = l2.NewTag("cat1")
// 			CAT2 = l2.NewTag("cat2")
// 		}).SetFormat(alog.Fprefix | alog.Flevel).SetPrefix("nptest ").SetLogTag(CAT1)
//
// 	l.SetLogTag(CAT1) // Print only CAT1
// 	WarnCAT1 := l.NewPrint(alog.Lwarn, CAT1, "CAT1W ")
// 	WarnCAT2 := l.NewPrint(alog.Lwarn, CAT2, "CAT2W ")
// 	TraceCAT1 := l.NewPrint(alog.Ltrace, CAT1, "CAT1T ")
// 	TraceCAT2 := l.NewPrint(alog.Ltrace, CAT2, "CAT2T ")
//
// 	// Since category is set to CAT1, and default Level is INFO,
// 	// only item(s) with CAT1 and INFO and above will be printed.
// 	WarnCAT1("warn cat1 test")
// 	WarnCAT2("warn cat2 test")
// 	TraceCAT1("trace cat1 test")
// 	TraceCAT2("trace cat2 test")
//
// 	// Output:
// 	// nptest [WRN] CAT1W warn cat1 test
// }
//
// func ExampleLogger_NewWriter() {
// 	var TEST1, TEST2, TEST3 alog.Tag
// 	l := alog.New(os.Stdout).Do(
// 		func(l2 *alog.Logger) {
// 			TEST1 = l2.NewTag("test1")
// 			TEST2 = l2.NewTag("test2")
// 			TEST3 = l2.NewTag("test3")
// 		}).SetFormat(alog.Fprefix | alog.Flevel).SetPrefix("nptest ").SetLogLevel(alog.Ldebug)
//
// 	// only show TEST2 here.
// 	// Therefore only DEBUG/INFO with TEST2 will be printed
// 	l.SetLogTag(TEST2)
//
// 	wT1D := l.NewWriter(alog.Ldebug, TEST1, "T1D ")
// 	wT1I := l.NewWriter(alog.Linfo, TEST1, "T1I ")
// 	wT2D := l.NewWriter(alog.Ldebug, TEST2, "T2D ")
// 	wT2I := l.NewWriter(alog.Linfo, TEST2, "T2I ")
// 	wT3D := l.NewWriter(alog.Ldebug, TEST3, "T3D ")
// 	wT3I := l.NewWriter(alog.Linfo, TEST3, "T3I ")
//
// 	_, _ = fmt.Fprintf(wT1D, "test: %s fprintf", "T1D")
// 	_, _ = fmt.Fprintf(wT1I, "test: %s fprintf", "T1I")
// 	_, _ = fmt.Fprintf(wT2D, "test: %s fprintf", "T2D")
// 	_, _ = fmt.Fprintf(wT2I, "test: %s fprintf", "T2I")
// 	_, _ = fmt.Fprintf(wT3D, "test: %s fprintf", "T3D")
// 	_, _ = fmt.Fprintf(wT3I, "test: %s fprintf", "T3I")
//
// 	// Output:
// 	// nptest [DBG] T2D test: T2D fprintf
// 	// nptest [INF] T2I test: T2I fprintf
// }
//
// func ExampleLogger_NewTag() {
// 	l := alog.New(os.Stdout).SetPrefix("test ").SetFormat(alog.Flevel | alog.Fprefix)
//
// 	// Creating new tags.
// 	BACK := l.NewTag("back")
// 	FRONT := l.NewTag("front")
// 	USER := l.NewTag("user")
//
// 	l.SetLogTag(BACK | FRONT) // only show BACK and FRONT
//
// 	f := func(c alog.Tag, s string) {
// 		l.Printf(alog.Ltrace, c, "%s.trace", s)
// 		l.Printf(alog.Ldebug, c, "%s.debug", s)
// 		l.Printf(alog.Linfo, c, "%s.info", s)
// 		l.Printf(alog.Lwarn, c, "%s.warn", s)
// 		l.Printf(alog.Lerror, c, "%s.error", s)
// 	}
//
// 	f(BACK, "BACK")
// 	f(FRONT, "FRONT")
// 	f(USER, "USER")
//
// 	// Output:
// 	// test [INF] BACK.info
// 	// test [WRN] BACK.warn
// 	// test [ERR] BACK.error
// 	// test [INF] FRONT.info
// 	// test [WRN] FRONT.warn
// 	// test [ERR] FRONT.error
// }
