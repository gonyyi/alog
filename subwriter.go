package alog

// SubWriter is a writer with predefined Level and Tag.
type SubWriter struct {
	l      *Logger
	wLevel Level // writer's level
	wTag   Tag   // writer's tag
}

// Write is to be used as io.Writer interface
func (w *SubWriter) Write(b []byte) (n int, err error) {
	return w.l.logb(w.wLevel, w.wTag, b)
}

// WriteString enables to write with string param
func (w *SubWriter) WriteString(s string) (n int, err error) {
	return w.l.Log(w.wLevel, w.wTag, s)
}
