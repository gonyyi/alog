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
