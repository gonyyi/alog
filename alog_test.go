package alog_test

import (
	"github.com/gonyyi/alog"
	"os"
	"testing"
)

func BenchmarkNew(b *testing.B) {
	al := alog.New(os.Stderr)
	al.Log(alog.Lerror, 0, "msg for err")

	al.SetOutput(nil)

	b.Run("al-kv", func(c *testing.B) {
		c.ReportAllocs()
		for i := 0; i < c.N; i++ {
			al.Log(alog.Lerror, 0, "msg for err")
		}
	})
}
