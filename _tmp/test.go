package main

import (
	"github.com/gonyyi/alog"
	"github.com/gonyyi/alog/ext"
	"os"
)

func main() {
	al := alog.New(os.Stderr).Ext(ext.LogFmt.Text())
	tagDisk := al.NewTag("Disk")
	tagDB := al.NewTag("DB")
	al.Info(tagDisk).Str("action", "reading disk").Write()
	al.Info(tagDB).Str("id", "myID").Str("pwd", "myPasswd").Write("Login")
	al.Info(tagDisk, tagDB).Int("status", 200).Write("Login")
	al.Info(tagDisk|tagDB).Int("status", 200).Write("Logout")
}
