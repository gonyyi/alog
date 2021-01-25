package alog

// SubWriter is a writer with predefined Level and Tag.
type SubWriter struct {
	l   *Logger
	lvl Level
	tag Tag
}

// Write is to be used as io.Writer interface
func (w *SubWriter) Write(b []byte) (n int, err error) {
	return w.l.logb(w.lvl, w.tag, b) // todo: this byte to string conversion need to be optimized someday
}

func (w *SubWriter) WriteString(s string) (n int, err error) {
	return w.l.Log(w.lvl, w.tag, s) // todo: this byte to string conversion need to be optimized someday
}
