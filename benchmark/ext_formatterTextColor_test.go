package benchmark

import (
	"github.com/gonyyi/alog"
	"github.com/gonyyi/alog/ext"
	"os"
	"testing"
)

func Test_ExtFormatterTextColor(t *testing.T) {
	al := alog.New(os.Stdout)
	al.SetFormatter(ext.FormatterTextColor())
	al.SetControl(0, 0)
	al.Trace(0, "test alog with color text formatter ext", "extUsed", "ext.FormatterTextColor()", "alogVer", "0.5.0")
	al.Debug(0, "test alog with color text formatter ext", "extUsed", "ext.FormatterTextColor()", "alogVer", "0.5.0")
	al.Info(0, "test alog with color text formatter ext", "extUsed", "ext.FormatterTextColor()", "alogVer", "0.5.0")
	al.Warn(0, "test alog with color text formatter ext", "extUsed", "ext.FormatterTextColor()", "alogVer", "0.5.0")
	al.Error(0, "test alog with color text formatter ext", "extUsed", "ext.FormatterTextColor()", "alogVer", "0.5.0")
	al.Fatal(0, "test alog with color text formatter ext", "extUsed", "ext.FormatterTextColor()", "alogVer", "0.5.0")
}

func Test_ExtFormatterTextColor2(t *testing.T) {
	al := alog.New(os.Stdout)
	IO := al.MustGetTag("IO")
	SYS := al.MustGetTag("SYS")
	NET := al.MustGetTag("NET")
	USER := al.MustGetTag("USER")
	BACKEND := al.MustGetTag("BACKEND")
	DISK := al.MustGetTag("DISK")

	// To make the compiler not complaining about variables not being used.
	_, _, _, _, _, _ = IO, SYS, NET, USER, BACKEND, DISK

	al.SetFormatter(ext.FormatterTextColor())
	al.SetControl(0, 0)
	al.Trace(IO, "test alog for IO", "extUsed", "ext.FormatterTextColor()", "alogVer", "0.5.0")
	al.Debug(SYS, "test alog for SYS", "extUsed", "ext.FormatterTextColor()", "alogVer", "0.5.0")
	al.Info(NET, "test alog for NET", "extUsed", "ext.FormatterTextColor()", "alogVer", "0.5.0")
	al.Warn(USER, "test alog for USER", "extUsed", "ext.FormatterTextColor()", "alogVer", "0.5.0")
	al.Error(IO|NET|DISK, "test alog for IO|NET|DISK", "extUsed", "ext.FormatterTextColor()", "alogVer", "0.5.0")
	al.Error(BACKEND, "test alog for BACKEND", "extUsed", "ext.FormatterTextColor()", "alogVer", "0.5.0")
	al.Fatal(DISK, "test alog for DISK", "extUsed", "ext.FormatterTextColor()", "alogVer", "0.5.0")
}
