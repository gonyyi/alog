package ext

import (
	"github.com/gonyyi/alog"
	"os"
)

var LogMode logMode

type logMode struct{}

func (logMode) Prod(filename string) alog.LoggerFn {
	return func(l alog.Logger) alog.Logger {
		l.Control.Level = alog.InfoLevel
		l.Flag = alog.WithDefault | alog.WithUnixTimeMs
		bw, err := NewBufWriter(filename)
		if err != nil {
			l.Error(0).Err( err).Write("failed to open")
		} else {
			l = l.SetOutput(bw)
		}
		return l
	}
}

func (logMode) Dev(filename string) alog.LoggerFn {
	return func(l alog.Logger) alog.Logger {
		l.Control.Level = alog.TraceLevel
		l.Flag = alog.WithTimeMs | alog.WithDefault
		if fo, err := os.Create(filename); err != nil {
			l.Error(0).Err( err).Write("cannot create file")
		} else {
			l = l.SetOutput(fo)
		}
		return l
	}
}

func (logMode) Test(filename string) alog.LoggerFn {
	return func(l alog.Logger) alog.Logger {
		l.Control.Level = alog.TraceLevel
		l.Flag = alog.WithTimeMs | alog.WithTag | alog.WithLevel
		l = l.SetFormatter(NewFormatterTerminalColor())
		return l
	}
}
