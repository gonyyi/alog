package main_test

import (
	"errors"
	"github.com/gonyyi/alog"
	"github.com/rs/zerolog"
	"os"
	"testing"
	"time"
)

func BenchmarkNew(b *testing.B) {
	al := alog.New(nil)
	al.Flag = alog.Flevel
	al.Control.Level = alog.Linfo
	// al.SetFormat(alog.Flevel|alog.Ftime|alog.Fdate)
	//al.SetFormat(alog.Flevel)
	zl := zerolog.New(nil)
	zl = zl.Level(zerolog.InfoLevel)

	movingStr := []string{"a", "b", "c", "d", "e"}
	myerr := errors.New("err test")
	_, _ = movingStr, myerr
	msg := `all good characters shown`

	fal := func(i int) {
		al.Info(0).
			Str("name", "gonal").
			Int("count", i).
			Str("block", movingStr[i%5]).
			Writes(msg)
	}
	fzl := func(i int) {
		zl.Info().
			Str("name", "gonzl").
			Int("count", i).
			Str("block", movingStr[i%5]).
			Msg(msg)
	}

	if 1 == 1 { // 120 vs 92 = al is 23% faster

		for rep := 0; rep < 5; rep++ {
			b.Run("zl", func(c *testing.B) {
				c.ReportAllocs()
				for i := 0; i < c.N; i++ {
					fzl(i)
				}
			})
			b.Run("al", func(c *testing.B) {
				c.ReportAllocs()
				for i := 0; i < c.N; i++ {
					fal(i)
				}
			})

		}
	}

	if 1 == 0 { // 25->18; al is 28% faster
		time.Sleep(time.Second)
		for rep := 0; rep < 5; rep++ {
			b.Run("al", func(c *testing.B) {
				c.ReportAllocs()
				c.RunParallel(func(p *testing.PB) {
					for p.Next() {
						fal(rep)
					}
				})
			})
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

	if 1 == 0 { // 25->18; al is 28% faster
		zl = zl.Level(zerolog.FatalLevel)
		al.Control.Level = alog.Lfatal
		time.Sleep(time.Second)
		for rep := 0; rep < 5; rep++ {

			b.Run("al", func(c *testing.B) {
				c.ReportAllocs()
				c.RunParallel(func(p *testing.PB) {
					for p.Next() {
						fal(rep)
					}
				})
			})
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
	// Print
	{
		al.SetOutput(os.Stderr)
		zl = zl.Output(os.Stderr)
		fal(123)
		fzl(123)
	}
}
