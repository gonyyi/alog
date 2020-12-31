// (c) 2020 Gon Y Yi. <https://gonyyi.com>
// Version 0.1.3, 12/29/2020

package alog_test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gonyyi/alog"
	"testing"
)

// added as v0.1.6c3, 12/30/2020
func TestLogger_IfError(t *testing.T) {
	out := &bytes.Buffer{}

	// v0.1.6 Code:
	// l := alog.New(out, "log ", alog.Fprefix|alog.Flevel)

	// v0.1.7 Code:
	l := alog.New(out).SetPrefix("log ").SetFlag(alog.Fprefix | alog.Flevel)

	// create a new error object, at this point this should be nil, and calling err.Error() will cause a panic
	var err error
	l.IfError(err)
	err = errors.New("test error") // now error exists
	l.IfError(err)
	err = nil // empty err again with nil
	l.IfError(err)
	err = errors.New("again another one") // now error exists again
	l.IfError(err)

	expected := "log [ERR] test error\nlog [ERR] again another one\n"

	if expected != out.String() {
		t.Errorf("expected=<%s>, actual=<%s>", expected, out.String())
	}
}

func TestBasic(t *testing.T) {
	// NoLevel will use INFO as its level.
	t.Run("Print,NoLevel", func(t2 *testing.T) {
		out := &bytes.Buffer{}

		// v0.1.6 Code:
		// l := alog.New(out, "log ", alog.Fprefix|alog.Flevel)

		// v0.1.7 Code:
		l := alog.New(out).SetPrefix("log ").SetFlag(alog.Fprefix | alog.Flevel)

		l.Print(alog.Ltrace, 0, "testTrace")
		l.Print(alog.Ldebug, 0, "testDebug")
		l.Print(alog.Linfo, 0, "testInfo")
		l.Print(alog.Lwarn, 0, "testWarn")
		l.Print(alog.Lerror, 0, "testError")
		expect := "log [INF] testInfo\nlog [WRN] testWarn\nlog [ERR] testError\n"
		actual := out.String()
		if expect != actual {
			t2.Errorf("expected=<%s>, actual=<%s>", expect, actual)
		}
	})
	t.Run("Printf,NoLevel", func(t2 *testing.T) {
		out := &bytes.Buffer{}

		// v0.1.6 Code:
		// l := alog.New(out, "log ", alog.Fprefix|alog.Flevel)

		// v0.1.7 Code:
		l := alog.New(out).SetPrefix("log ").SetFlag(alog.Fprefix | alog.Flevel)

		l.Printf(alog.Ltrace, 0, "test%s", "Trace")
		l.Printf(alog.Ldebug, 0, "test%s", "Debug")
		l.Printf(alog.Linfo, 0, "test%s", "Info")
		l.Printf(alog.Lwarn, 0, "test%s", "Warn")
		l.Printf(alog.Lerror, 0, "test%s", "Error")
		expect := "log [INF] testInfo\nlog [WRN] testWarn\nlog [ERR] testError\n"
		actual := out.String()
		if expect != actual {
			t2.Errorf("expected=<%s>, actual=<%s>", expect, actual)
		}
	})
	t.Run("Print,Trace", func(t2 *testing.T) {
		out := &bytes.Buffer{}

		// v0.1.6 Code:
		// l := alog.New(out, "log ", alog.Fprefix|alog.Flevel)

		// v0.1.7 Code:
		l := alog.New(out).SetPrefix("log ").SetFlag(alog.Fprefix | alog.Flevel)

		l.SetLevel(alog.Ltrace)
		l.Print(alog.Ltrace, 0, "testTrace")
		l.Print(alog.Ldebug, 0, "testDebug")
		l.Print(alog.Linfo, 0, "testInfo")
		l.Print(alog.Lwarn, 0, "testWarn")
		l.Print(alog.Lerror, 0, "testError")
		expect := "log [TRC] testTrace\nlog [DBG] testDebug\nlog [INF] testInfo\nlog [WRN] testWarn\nlog [ERR] testError\n"
		actual := out.String()
		if expect != actual {
			t2.Errorf("expected=<%s>, actual=<%s>", expect, actual)
		}
	})
	t.Run("Printf,Trace", func(t2 *testing.T) {
		out := &bytes.Buffer{}

		// v0.1.6 Code:
		// l := alog.New(out, "log ", alog.Fprefix|alog.Flevel)

		// v0.1.7 Code:
		l := alog.New(out).SetPrefix("log ").SetFlag(alog.Fprefix | alog.Flevel)

		l.SetLevel(alog.Ltrace)
		l.Printf(alog.Ltrace, 0, "test%s", "Trace")
		l.Printf(alog.Ldebug, 0, "test%s", "Debug")
		l.Printf(alog.Linfo, 0, "test%s", "Info")
		l.Printf(alog.Lwarn, 0, "test%s", "Warn")
		l.Printf(alog.Lerror, 0, "test%s", "Error")
		expect := "log [TRC] testTrace\nlog [DBG] testDebug\nlog [INF] testInfo\nlog [WRN] testWarn\nlog [ERR] testError\n"
		actual := out.String()
		if expect != actual {
			t2.Errorf("expected=<%s>, actual=<%s>", expect, actual)
		}
	})
	t.Run("Print,Fatal", func(t2 *testing.T) {
		out := &bytes.Buffer{}

		// v0.1.6 Code:
		// l := alog.New(out, "log ", alog.Fprefix|alog.Flevel)

		// v0.1.7 Code:
		l := alog.New(out).SetPrefix("log ").SetFlag(alog.Fprefix | alog.Flevel)

		l.SetLevel(alog.Lfatal)
		l.Print(alog.Ltrace, 0, "testTrace")
		l.Print(alog.Ldebug, 0, "testDebug")
		l.Print(alog.Linfo, 0, "testInfo")
		l.Print(alog.Lwarn, 0, "testWarn")
		l.Print(alog.Lerror, 0, "testError")
		expect := ""
		actual := out.String()
		if expect != actual {
			t2.Errorf("expected=<%s>, actual=<%s>", expect, actual)
		}
	})
	t.Run("Printf,Fatal", func(t2 *testing.T) {
		out := &bytes.Buffer{}

		// v0.1.6 Code:
		// l := alog.New(out, "log ", alog.Fprefix|alog.Flevel)
		// l.SetLevel(alog.Lfatal)

		// v0.1.7 Code:
		l := alog.New(out).SetPrefix("log ").SetFlag(alog.Fprefix | alog.Flevel).SetLevel(alog.Lfatal)

		l.Printf(alog.Ltrace, 0, "test%s", "Trace")
		l.Printf(alog.Ldebug, 0, "test%s", "Debug")
		l.Printf(alog.Linfo, 0, "test%s", "Info")
		l.Printf(alog.Lwarn, 0, "test%s", "Warn")
		l.Printf(alog.Lerror, 0, "test%s", "Error")
		expect := ""
		actual := out.String()
		if expect != actual {
			t2.Errorf("expected=<%s>, actual=<%s>", expect, actual)
		}
	})
	t.Run("Predefined,NoLevel", func(t2 *testing.T) {
		out := &bytes.Buffer{}

		// v0.1.6 Code:
		// l := alog.New(out, "log ", alog.Fprefix|alog.Flevel)
		// l.SetLevel(alog.Linfo)

		// v0.1.7 Code:
		l := alog.New(out).SetPrefix("log ").SetFlag(alog.Fprefix | alog.Flevel).SetLevel(alog.Linfo)

		l.Trace("testForTrace")
		l.Debug("testForDebug")
		l.Info("testForInfo")
		l.Warn("testForWarn")
		l.Error("testForError")
		expect := "log [INF] testForInfo\nlog [WRN] testForWarn\nlog [ERR] testForError\n"
		actual := out.String()
		if expect != actual {
			t2.Errorf("expected=<%s>, actual=<%s>", expect, actual)
		}
	})
	t.Run("Predefined,Formatted,level", func(t2 *testing.T) {
		out := &bytes.Buffer{}

		// v0.1.6 Code:
		// l := alog.New(out, "log ", alog.Fprefix|alog.Flevel)
		// l.SetLevel(alog.Linfo)

		// v0.1.7 Code:
		// level info is default level and can be ignored
		l := alog.New(out).SetPrefix("log ").SetFlag(alog.Fprefix | alog.Flevel)

		l.Tracef("testFor%s", "Trace")
		l.Debugf("testFor%s", "Debug")
		l.Infof("testFor%s", "Info")
		l.Warnf("testFor%s", "Warn")
		l.Errorf("testFor%s", "Error")
		expect := "log [INF] testForInfo\nlog [WRN] testForWarn\nlog [ERR] testForError\n"
		actual := out.String()
		if expect != actual {
			t2.Errorf("expected=<%s>, actual=<%s>", expect, actual)
		}
	})
}

// v0.1.1, 12/29/2020
func TestNewPrint(t *testing.T) {
	t.Run("NewPrint", func(t2 *testing.T) {
		out := &bytes.Buffer{}

		// v0.1.6 Code:
		// l := alog.New(out, "nptest ", alog.Fprefix|alog.Flevel)

		// v0.1.7 Code:
		l := alog.New(out).SetPrefix("nptest ").SetFlag(alog.Fprefix | alog.Flevel)

		cat := alog.NewCategory()
		CAT1 := cat.Add()
		CAT2 := cat.Add()

		l.SetCategory(CAT1) // Print only CAT1
		WarnCAT1 := l.NewPrint(alog.Lwarn, CAT1, "CAT1W ")
		WarnCAT2 := l.NewPrint(alog.Lwarn, CAT2, "CAT2W ")
		TraceCAT1 := l.NewPrint(alog.Ltrace, CAT1, "CAT1T ")
		TraceCAT2 := l.NewPrint(alog.Ltrace, CAT2, "CAT2T ")
		// WarnCAT1 := l.NewPrint(alog.Lwarn, CAT1, []byte("CAT1W "))
		// WarnCAT2 := l.NewPrint(alog.Lwarn, CAT2, []byte("CAT2W "))
		// TraceCAT1 := l.NewPrint(alog.Ltrace, CAT1, []byte("CAT1T "))
		// TraceCAT2 := l.NewPrint(alog.Ltrace, CAT2, []byte("CAT2T "))

		WarnCAT1("warn cat1 test")
		WarnCAT2("warn cat2 test")
		TraceCAT1("trace cat1 test")
		TraceCAT2("trace cat2 test")

		actual := out.String()
		expect := "nptest [WRN] CAT1W warn cat1 test\n"
		if expect != actual {
			t2.Errorf("expected=<%s>, actual=<%s>", expect, actual)
		}
	})
}

// v0.1.2, 12/29/2020
func TestLogger_NewWriter(t *testing.T) {
	t.Run("NewWriter", func(t2 *testing.T) {
		out := &bytes.Buffer{}

		// v0.1.6 Code:
		// l := alog.New(out, "nwtest ", alog.Fprefix|alog.Flevel) // Default level is INFO and higher

		// v0.1.7 Code:
		l := alog.New(out).SetPrefix("nwtest ").SetFlag(alog.Fprefix | alog.Flevel)

		cat := alog.NewCategory()
		TEST1 := cat.Add()
		TEST2 := cat.Add()
		TEST3 := cat.Add()

		l.SetCategory(TEST2) // only show TEST2

		wT1D := l.NewWriter(alog.Ldebug, TEST1, "T1D ")
		wT1I := l.NewWriter(alog.Linfo, TEST1, "T1I ")
		wT2D := l.NewWriter(alog.Ldebug, TEST2, "T2D ")
		wT2I := l.NewWriter(alog.Linfo, TEST2, "T2I ")
		wT3D := l.NewWriter(alog.Ldebug, TEST3, "T3D ")
		wT3I := l.NewWriter(alog.Linfo, TEST3, "T3I ")

		_, _ = fmt.Fprintf(wT1D, "test: %s fprintf", "T1D") // Not printed
		_, _ = fmt.Fprintf(wT1I, "test: %s fprintf", "T1I") // Not printed
		_, _ = fmt.Fprintf(wT2D, "test: %s fprintf", "T2D") // Not printed
		_, _ = fmt.Fprintf(wT2I, "test: %s fprintf", "T2I") // Printed
		_, _ = fmt.Fprintf(wT3D, "test: %s fprintf", "T3D") // Not printed
		_, _ = fmt.Fprintf(wT3I, "test: %s fprintf", "T3I") // Not printed

		expect := "nwtest [INF] T2I test: T2I fprintf\n"
		actual := out.String()

		if expect != actual {
			t2.Errorf("expected=<%s>, actual=<%s>", expect, actual)
		}
	})
}

func TestEnum(t *testing.T) {
	if alog.Fall != 4294967295 {
		t.Fatalf("Unexpected alog.Fall value, exp=%d vs act=%d", 4294967295, alog.Fall)
	}
	if alog.Fnone != 0 {
		t.Fatalf("Unexpected alog.Fnone value, exp=%d vs act=%d", 0, alog.Fall)
	}
}
