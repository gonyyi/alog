package ext

import "github.com/gonyyi/alog"

var DoFmt doFormatter

type doFormatter struct{}

func (doFormatter) TXT() alog.DoFn {
	return func(l alog.Logger) alog.Logger {
		l = l.SetFormatter(NewFormatterTerminal())
		return l
	}
}

func (doFormatter) TXTColor() alog.DoFn {
	return func(l alog.Logger) alog.Logger {
		l = l.SetFormatter(NewFormatterTerminalColor())
		return l
	}
}
