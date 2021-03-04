package main_test

import (
	"errors"
	"github.com/gonyyi/alog"
	"github.com/rs/zerolog"
	"os"
	"testing"
)

var dataComp = struct {
	StrSlice []string
	Str1     string
	Float    float64
	Int      int
	Error    error
	Msg      string
}{
	StrSlice: []string{"a", "b", "c", "d", "e"},
	Error:    errors.New("err test"),
	Msg:      "test message",
}

func fal(i int) {
	al.Info(0).
		Str("name", "gonal").
		Int("count", i).
		Str("block", dataComp.StrSlice[i%5]).
		Write(dataComp.Msg)
}
func fzl(i int) {
	zl.Info().
		Str("name", "gonzl").
		Int("count", i).
		Str("block", dataComp.StrSlice[i%5]).
		Msg(dataComp.Msg)
}

var al = alog.New(nil)
var zl = zerolog.New(nil)

func init() {
	al.Flag = alog.Flevel
	al.Control.Level = alog.Linfo
	zl = zl.Level(zerolog.InfoLevel)
}

func TestCompSimple(t *testing.T) {
	al.SetOutput(os.Stderr)
	zl = zl.Output(os.Stderr)
	fal(123)
	fzl(123)
	al.SetOutput(nil)
	zl = zl.Output(nil)
}

func BenchmarkCompSingleThread(b *testing.B) {
	for rep := 0; rep < 5; rep++ {
		b.Run("zl", func(c *testing.B) {
			c.ReportAllocs()
			for i := 0; i < c.N; i++ {
				fzl(i)
			}
		})
	}
	for rep := 0; rep < 5; rep++ {
		b.Run("al", func(c *testing.B) {
			c.ReportAllocs()
			for i := 0; i < c.N; i++ {
				fal(i)
			}
		})
	}
}

func BenchmarkCompParallel(b *testing.B) {
	for rep := 0; rep < 5; rep++ {
		b.Run("al", func(c *testing.B) {
			c.ReportAllocs()
			c.RunParallel(func(p *testing.PB) {
				for p.Next() {
					fal(rep)
				}
			})
		})
	}
	for rep := 0; rep < 5; rep++ {
		b.Run("zl", func(c *testing.B) {
			c.ReportAllocs()
			c.RunParallel(func(p *testing.PB) {
				for p.Next() {
					fzl(rep)
				}
			})
		})
	}
}

func BenchmarkCompCheck(b *testing.B) {
	zl = zl.Level(zerolog.FatalLevel)
	al.Control.Level = alog.Lfatal

	for rep := 0; rep < 5; rep++ {
		b.Run("al", func(c *testing.B) {
			c.ReportAllocs()
			c.RunParallel(func(p *testing.PB) {
				for p.Next() {
					fal(rep)
				}
			})
		})
	}
	for rep := 0; rep < 5; rep++ {
		b.Run("zl", func(c *testing.B) {
			c.ReportAllocs()
			c.RunParallel(func(p *testing.PB) {
				for p.Next() {
					fzl(rep)
				}
			})
		})
	}
	zl = zl.Level(zerolog.InfoLevel)
	al.Control.Level = alog.Lfatal
}
