// (c) 2020 Gon Y Yi. <https://gonyyi.com/copyright.txt>
// Version 2, 12/21/2020

package alog_test

import (
	"bytes"
	"github.com/gonyyi/alog"
	"testing"
)

func TestBasic(t *testing.T) {
	// NoLevel will use INFO as its level.
	t.Run("Print,NoLevel", func(t2 *testing.T) {
		out := &bytes.Buffer{}
		l := alog.New(out, "log ", alog.Fprefix|alog.Flevel)
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
		l := alog.New(out, "log ", alog.Fprefix|alog.Flevel)
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
		l := alog.New(out, "log ", alog.Fprefix|alog.Flevel)
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
		l := alog.New(out, "log ", alog.Fprefix|alog.Flevel)
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
		l := alog.New(out, "log ", alog.Fprefix|alog.Flevel)
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
		l := alog.New(out, "log ", alog.Fprefix|alog.Flevel)
		l.SetLevel(alog.Lfatal)
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
		l := alog.New(out, "log ", alog.Fprefix|alog.Flevel)
		l.SetLevel(alog.Linfo)
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
		l := alog.New(out, "log ", alog.Fprefix|alog.Flevel)
		l.SetLevel(alog.Linfo)
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

func TestNewPrint(t *testing.T) {
	t.Run("NewPrint", func(t2 *testing.T) {
		out := &bytes.Buffer{}

		l := alog.New(out, "nptest ", alog.Fprefix|alog.Flevel) // Default level is INFO and higher

		cat := alog.NewCategory()
		CAT1 := cat.Add()
		CAT2 := cat.Add()

		l.SetCategory(CAT1) // Print only CAT1
		WarnCAT1 := l.NewPrint(alog.Lwarn, CAT1)
		WarnCAT2 := l.NewPrint(alog.Lwarn, CAT2)
		TraceCAT1 := l.NewPrint(alog.Ltrace, CAT1)
		TraceCAT2 := l.NewPrint(alog.Ltrace, CAT2)

		WarnCAT1("warn cat1 test")
		WarnCAT2("warn cat2 test")
		TraceCAT1("trace cat1 test")
		TraceCAT2("trace cat2 test")

		actual := out.String()
		expect := "nptest [WRN] warn cat1 test\n"
		if expect != actual {
			t2.Errorf("expected=<%s>, actual=<%s>", expect, actual)
		}
	})
}
