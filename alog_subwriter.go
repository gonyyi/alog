package alog

// SubWriter is a writer with predefined Level and Tag.
type SubWriter struct {
	l   *Logger
	lvl Level
	tag Tag
}

// Write is to be used as io.Writer interface
func (w *SubWriter) Write(b []byte) (n int, err error) {
	if w.l.check(w.lvl, w.tag) {
		w.l.mu.Lock()
		w.l.header(w.lvl, w.tag)

		w.l.buf = append(w.l.buf, b...) // todo: check if this works with JSON
		n, err := w.l.finalize()
		w.l.mu.Unlock()
		return n, err
	}
	return 0, nil
}
