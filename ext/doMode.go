package ext

import (
	"github.com/gonyyi/alog"
)

func DoModeProd(filename string) alog.DoFn {
	return func(l alog.Logger) alog.Logger {
		l.Control.Level = alog.Linfo
		l.Flag = alog.Fdefault|alog.FtimeUnixMs
		bw, err := NewBufWriter(filename)
		if err != nil {
			l.Error(0).Err("err", err).Write("failed to open")
		} else {
			l = l.SetOutput(bw)
		}
		return l
	}
}

func DoModeDev(filenamePlaceHolder string) alog.DoFn {
	return func(l alog.Logger) alog.Logger {
		l.Control.Level = alog.Ltrace
		l = l.SetFormatter(NewFormatterTerminalColor())
		return l
	}
}
