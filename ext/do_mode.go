package ext

import (
	"github.com/gonyyi/alog"
	"os"
)

var DoMode doMode

type doMode struct{}

func (doMode) PROD(filename string) alog.DoFn {
	return func(l alog.Logger) alog.Logger {
		l.Control.Level = alog.InfoLevel
		l.Flag = alog.UseDefault | alog.UseUnixTimeMs
		bw, err := NewBufWriter(filename)
		if err != nil {
			l.Error(0).Err("err", err).Write("failed to open")
		} else {
			l = l.SetOutput(bw)
		}
		return l
	}
}

func (doMode) DEV(filename string) alog.DoFn {
	return func(l alog.Logger) alog.Logger {
		l.Control.Level = alog.TraceLevel
		l.Flag = alog.UseTimeMs | alog.UseDefault
		if fo, err := os.Create(filename); err != nil {
			l.Error(0).Err("error", err).Write("cannot create file")
		} else {
			l = l.SetOutput(fo)
		}
		return l
	}
}

func (doMode) TEST(filename string) alog.DoFn {
	return func(l alog.Logger) alog.Logger {
		l.Control.Level = alog.TraceLevel
		l.Flag = alog.UseTimeMs | alog.UseTag | alog.UseLevel
		l = l.SetFormatter(NewFormatterTerminalColor())
		return l
	}
}
