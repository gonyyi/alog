package main_test

import (
	"github.com/gonyyi/alog"
	"github.com/gonyyi/alog/ext"
	"os"
	"testing"
)

func TestExtLogMode(t *testing.T) {
  if skip_print {
    t.SkipNow()
  }
	al := alog.New(nil).SetOutput(os.Stderr).Ext(ext.LogMode.Test("test.log"))
	tOS := al.NewTag("OS")
	tSYS := al.NewTag("SYS")
	al.Info(tOS).Str("status", "starting").Bool("isError", false).Write("starting")
	al.Info(tOS|tSYS).Str("status", "reading sys").Write("init")
	al.Close()
}
