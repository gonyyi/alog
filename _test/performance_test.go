package main

import (
	"github.com/gonyyi/alog"
	"github.com/rs/zerolog"
	"testing"
)

func BenchmarkSubWriter(b *testing.B) {
	b.ReportAllocs()

	l := alog.New(nil).SetNewTags("mytag").SetFormat(alog.Fjson | alog.Ftag)
	sw := l.NewWriter(alog.Lerror, l.MustGetTag("mytag"))

	msg := []byte(`my test`)
	for i := 0; i < b.N; i++ {
		sw.Write(msg)
	}
	// l.SetOutput(os.Stdout)
	// sw.Write(msg)
}

func BenchmarkCompare(b *testing.B) {
	b.Run("t1: alog, msg", func(c *testing.B) {
		l := alog.New(nil).SetNewTags("t1", "t2", "t3")
		l.SetFormat(alog.Fjson)
		for i := 0; i < c.N; i++ {
			l.Error(0, "hello error 1")
			// l.Error(0, "hello error 1", "a", "b")
		}
	})

	b.Run("t1: zlog, msg", func(c *testing.B) {
		l := zerolog.New(nil)
		for i := 0; i < c.N; i++ {
			l.Error().Msg("hello error 1")
			// l.Error().Str("a", "b").Msg("hello error 1")
		}
	})

	b.Run("t2: alog, msg + s", func(c *testing.B) {
		l := alog.New(nil).SetNewTags("t1", "t2", "t3")
		l.SetFormat(alog.Fjson)
		for i := 0; i < c.N; i++ {
			// l.Error(0, "hello error 1")
			l.Error(0, "hello error 1", "a", "b")
		}
	})

	b.Run("t2: zlog, msg + s", func(c *testing.B) {
		l := zerolog.New(nil)
		for i := 0; i < c.N; i++ {
			// l.Error().Msg("hello error 1")
			l.Error().Str("a", "b").Msg("hello error 1")
		}
	})

}
