package alog_test

import (
	alog "github.com/gonyyi/alog/v0.5"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	al := alog.New(os.Stdout)
	al.Log(alog.Lerror, 0, "test")
	al.Loga(alog.Lerror, 0, "name", "gon", "age", 39)
	al.Loga(alog.Lerror, 0, "name")
}

func BenchmarkNew(b *testing.B) {
	al := alog.New(nil)

	b.Run("Check", func(b2 *testing.B) {
		b2.ReportAllocs()
		for i := 0; i < b2.N; i++ {
			al.Log(alog.Lerror, 0, "test")
		}
	})
	b.Run("Loga", func(b2 *testing.B) {
		b2.ReportAllocs()
		for i := 0; i < b2.N; i++ {
			al.Loga(alog.Linfo, 0, "name", "gon")
		}
	})
}
