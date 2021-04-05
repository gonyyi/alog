package alog

import "testing"

func TestDiscard_Write(t *testing.T) {
	w := Discard{}
	w.Write([]byte("abc"))
	w.WriteLt([]byte("abc"), 0, 0)
}
