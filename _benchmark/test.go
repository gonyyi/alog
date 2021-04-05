package main

import (
	"github.com/gonyyi/alog"
	"github.com/gonyyi/alog/ext"
	"os"
)
func main() {
	al := alog.New(os.Stderr)
	al.Info().Write("test")
}
func t1() {
	al := alog.New(os.Stderr)
	al.Flag = al.Flag | alog.WithTimeMs
	al = al.Ext(ext.LogFmt.TextColor())

	tagDisk := al.NewTag("Disk")
	tagDB := al.NewTag("DB")

	al.Control.Level = alog.TraceLevel

	al.Trace(tagDisk).
		Int("testId", 1).Write("ok")
	al.Debug(tagDisk).
		Int("testId", 2).Write("ok")
	al.Info(tagDisk).
		Int("testId", 3).Write("ok")
	al.Warn().
		Int("testId", 4).Write("ok")
	al.Error(tagDisk, tagDB).
		Int("testId", 5).Write("ok")
	al.Fatal(tagDisk|tagDB).
		Int("testId", 6).Write("ok")
}
