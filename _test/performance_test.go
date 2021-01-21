package main

import (
	"github.com/gonyyi/alog"
	"github.com/rs/zerolog"
	"testing"
)

func BenchmarkAll(b *testing.B) {
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
