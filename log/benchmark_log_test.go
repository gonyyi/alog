package log_test

import (
	"github.com/gonyyi/alog/log"
	"testing"
)

func BenchmarkInfo(b *testing.B) {
	b.ReportAllocs()
	log.SetOutput(nil)
	for i := 0; i < b.N; i++ {
		log.Info(0, "abc", "def", "name", "age", 123)
	}
}
