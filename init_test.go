package alog_test

import (
	"github.com/gonyyi/alog"
	"testing"
)

func TestLevel(t *testing.T) {
	check := func(a alog.Level, expName string, expNameShort string) {
		if a.Name() != expName || a.NameShort() != expNameShort {
			println("err")
			t.Errorf("unexpected level // a=<%d>, a.Name=<%s>, a.NameShort=<%s>",
				a, a.Name(), a.NameShort())
		}
	}

	check(0, "", "")
	check(1, "trace", "TRC")
	check(2, "debug", "DBG")
	check(3, "info", "INF")
	check(4, "warn", "WRN")
	check(5, "error", "ERR")
	check(6, "fatal", "FTL")
	check(7, "", "")
}
