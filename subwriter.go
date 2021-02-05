package alog

// SubWriter is a writer with predefined Level and Tag.
type SubWriter struct {
	l            *Logger
	defaultLevel Level
	defaultTag   Tag
}

// Write is to be used as io.Writer interface
func (w *SubWriter) Write(b []byte) (n int, err error) {
	return w.l.logb(w.defaultLevel, w.defaultTag, b)
}

// WriteString enables to write with string param
func (w *SubWriter) WriteString(s string) (n int, err error) {
	return w.l.Log(w.defaultLevel, w.defaultTag, s)
}
