package alog_test

import (
	"github.com/gonyyi/alog"
	"testing"
)

func TestControl_Check(t *testing.T) {
	ctl := log.Control // since control is unexported struct, use it from logger.

	if bucket := ctl.Bucket(); bucket == nil {
		t.Errorf("control.bucket is nil")
		t.SkipNow()
	}
	if ctl.Check(alog.InfoLevel, 0) != true {
		t.Errorf("unexpected control.bucket.Check() result for Info")
	}
	if ctl.Check(alog.TraceLevel, 0) != false {
		t.Errorf("unexpected control.bucket.Check() result for Trace")
	}

	{
		if ok, ok2 := ctl.CheckFn(alog.InfoLevel, 0); ok {
			t.Errorf("unexpected control.CheckFn() 1")
			_ = ok2
		}

		ctl.Fn = func(level alog.Level, tag alog.Tag) bool {
			return true
		}

		if ok, ok2 := ctl.CheckFn(alog.InfoLevel, 0); !ok {
			t.Errorf("unexpected control.CheckFn() 2")
		} else if ok2 != true {
			t.Errorf("unexpected control.CheckFn() 3 output")
		}

		ctl.Fn = func(level alog.Level, tag alog.Tag) bool {
			return false
		}

		if ok, ok2 := ctl.CheckFn(alog.InfoLevel, 0); !ok {
			t.Errorf("unexpected control.CheckFn() 4")
		} else if ok2 != false {
			t.Errorf("unexpected control.CheckFn() 5 output")
		}

		ctl.Fn = nil
		if ok, ok2 := ctl.CheckFn(alog.InfoLevel, 0); ok || ok2 {
			t.Errorf("unexpected control.CheckFn() 6")
		}
	}

	{
		log.Control.Tags = 123
		tmp := log.Control
		if tmp.Tags != 123 {
			t.Errorf("log.Cotrol.Tags does not retain updated value")
		}
	}
}
