package ext

import "github.com/gonyyi/alog"

var LogFmt logFormatter

type logFormatter struct{}

func (logFormatter) None() alog.LoggerFn {
	return func(l alog.Logger) alog.Logger {
		l = l.SetFormatter(nil)
		return l
	}
}

func (logFormatter) Text() alog.LoggerFn {
	return func(l alog.Logger) alog.Logger {
		l = l.SetFormatter(NewFormatterTerminal())
		return l
	}
}

func (logFormatter) TextColor() alog.LoggerFn {
	return func(l alog.Logger) alog.Logger {
		l = l.SetFormatter(NewFormatterTerminalColor())
		return l
	}
}
