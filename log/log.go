package log

import (
	"github.com/gonyyi/alog"
	"github.com/gonyyi/alog/ext"
	"io"
	"os"
)

var log *alog.Logger

func init() {
	log = alog.New(os.Stdout)
	log.SetFormatter(ext.FormatterTextColor())
}

// SetOutput will set the logger's output,
// default will be standard out.
func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

// Do will run (series of) function(s) and is used for
// quick macro like settings for the logger.
func Do(fns ...func(*alog.Logger)) {
	log.Do(fns...)
}

// MustGetTag will return a tag. If a required tag is not exists,
// it will create one.
func MustGetTag(name string) (tag alog.Tag) {
	return log.MustGetTag(name)
}

// SetControl will set logging level and tag.
// Note that this is an OR condition: if level has met the minimum logging level OR
// tag is met, the logger will log. For any precise control, use SetControlFn.
func SetControl(lv alog.Level, tag alog.Tag) {
	log.SetControl(lv, tag)
}

// Iferr method will log an error when argument err is not nil.
// This also returns true/false if error is or not nil.
func Iferr(err error, tag alog.Tag, msg string) bool {
	return log.Iferr(err, tag, msg)
}

// Trace records a msg with a trace level with optional additional variables
func Trace(tag alog.Tag, msg string, a ...interface{}) {
	log.Trace(tag, msg, a...)
}

// Debug records a msg with a debug level with optional additional variables
func Debug(tag alog.Tag, msg string, a ...interface{}) {
	log.Debug(tag, msg, a...)
}

// Info records a msg with an info level with optional additional variables
// And info level is default log level of Alog.
func Info(tag alog.Tag, msg string, a ...interface{}) {
	log.Info(tag, msg, a...)
}

// Warn records a msg with a warning level with optional additional variables
func Warn(tag alog.Tag, msg string, a ...interface{}) {
	log.Warn(tag, msg, a...)
}

// Error records a msg with an error level with optional additional variables
func Error(tag alog.Tag, msg string, a ...interface{}) {
	log.Error(tag, msg, a...)
}

// Fatal records a msg with a fatal level with optional additional variables.
// Unlike other logger, Alog will NOT terminal the program with a Fatal method.
// A user need to handle what to do.
func Fatal(tag alog.Tag, msg string, a ...interface{}) {
	log.Fatal(tag, msg, a...)
}
