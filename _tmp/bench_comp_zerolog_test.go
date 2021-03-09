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

const repeat = 3

func init() {
	al.Flag = alog.WithLevel
	al.Control.Level = alog.InfoLevel
	zl = zl.Level(zerolog.InfoLevel)
}

func TestWriteFile(t *testing.T) {
	alo, _ := os.Create("testOutAl.log")
	zlo, _ := os.Create("testOutZl.log")
	al := alog.New(alo)
	al.Flag = alog.WithLevel
	zl := zl.Output(zlo)

	for i := 0; i < 1_000_000; i++ {
		al.Info(0).
			Int64("count1", int64(i)).
			Int("count2", i).
			Str("randomStr", dataComp.StrSlice[i%5]).
			Float("float64", dataComp.Float+float64(i)).
			Bool("b", false).Write("")
		zl.Info().
			Int64("count1", int64(i)).
			Int("count2", i).
			Str("randomStr", dataComp.StrSlice[i%5]).
			Float64("float64", dataComp.Float+float64(i)).
			Bool("b", false).Send()
	}
}

func TestCompSimple(t *testing.T) {
	al = al.SetOutput(os.Stderr)
	zl = zl.Output(os.Stderr)
	fal(123)
	fzl(123)
	al = al.SetOutput(nil)
	zl = zl.Output(nil)
}

func BenchmarkCompSingleThread(b *testing.B) {
	for rep := 0; rep < repeat; rep++ {
		b.Run("zl", func(c *testing.B) {
			c.ReportAllocs()
			for i := 0; i < c.N; i++ {
				fzl(i)
			}
		})
	}
	for rep := 0; rep < repeat; rep++ {
		b.Run("al", func(c *testing.B) {
			c.ReportAllocs()
			for i := 0; i < c.N; i++ {
				fal(i)
			}
		})
	}
}

func BenchmarkCompParallel(b *testing.B) {
	for rep := 0; rep < repeat; rep++ {
		b.Run("al", func(c *testing.B) {
			c.ReportAllocs()
			c.RunParallel(func(p *testing.PB) {
				for p.Next() {
					fal(rep)
				}
			})
		})
	}
	for rep := 0; rep < repeat; rep++ {
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
	al.Control.Level = alog.FatalLevel

	for rep := 0; rep < repeat; rep++ {
		b.Run("al", func(c *testing.B) {
			c.ReportAllocs()
			c.RunParallel(func(p *testing.PB) {
				for p.Next() {
					fal(rep)
				}
			})
		})
	}
	for rep := 0; rep < repeat; rep++ {
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
	al.Control.Level = alog.FatalLevel
}
