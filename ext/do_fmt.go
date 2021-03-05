package ext

import "github.com/gonyyi/alog"

var LogFmt logFormatter

type logFormatter struct{}

func (logFormatter) NONE() alog.LoggerFn {
	return func(l alog.Logger) alog.Logger {
		l = l.SetFormatter(nil)
		return l
	}
}

func (logFormatter) TXT() alog.LoggerFn {
	return func(l alog.Logger) alog.Logger {
		l = l.SetFormatter(NewFormatterTerminal())
		return l
	}
}

func (logFormatter) TXTColor() alog.LoggerFn {
	return func(l alog.Logger) alog.Logger {
		l = l.SetFormatter(NewFormatterTerminalColor())
		return l
	}
}
