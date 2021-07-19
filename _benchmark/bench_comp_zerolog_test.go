package main_test

import (
	"errors"
	"github.com/gonyyi/alog"
	"github.com/gonyyi/alog/ext"
	"github.com/rs/zerolog"
	"os"
	"testing"
)

var skip_write_to_file = true
var skip_print = true 

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
	al.Info().
		Str("name", "gonal").
		Int("count", i).
		Str("block", dataComp.StrSlice[i%5]).
		Write(dataComp.StrSlice[i%5])
}

func fzl(i int) {
	zl.Info().
		Str("name", "gonzl").
		Int("count", i).
		Str("block", dataComp.StrSlice[i%5]).
		Msg(dataComp.StrSlice[i%5])
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
	if skip_write_to_file {
		t.SkipNow()
	}
	alo, _ := os.Create("testOutAl.log")
	zlo, _ := os.Create("testOutZl.log")
	al := alog.New(alo)
	al.Flag = alog.WithLevel
	zl := zl.Output(zlo)

	for i := 0; i < 1_000_000; i++ {
		al.Info().
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

func BenchmarkZlogWrite(b *testing.B) {
	// 287367	      3774 ns/op	       0 B/op	       0 allocs/op
	// 305935	      4031 ns/op	       0 B/op	       0 allocs/op
	// 303000	      3944 ns/op	       0 B/op	       0 allocs/op
	// M1: 918042	      1266 ns/op	       0 B/op	       0 allocs/op
	// M1: 901912	      1283 ns/op	       0 B/op	       0 allocs/op
	// M1: 920892	      1275 ns/op	       0 B/op	       0 allocs/op

	if skip_write_to_file {
		b.SkipNow()
	}
	out, _ := os.Create("./test-zl.log")
	zl := zerolog.New(out)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		switch i % 5 {
		case 0:
			zl.Trace().Int("count", i).Send()
		case 1:
			zl.Debug().Int("count", i).Send()
		case 2:
			zl.Info().Int("count", i).Send()
		case 3:
			zl.Warn().Int("count", i).Send()
		case 4:
			zl.Error().Int("count", i).Send()
		}
	}
}

func BenchmarkAlWriter(b *testing.B) {
	// 310676	      3689 ns/op	       0 B/op	       0 allocs/op
	// 309205	      3783 ns/op	       0 B/op	       0 allocs/op
	// 327218	      3853 ns/op	       0 B/op	       0 allocs/op
	// M1: 898946	      1297 ns/op	       0 B/op	       0 allocs/op
	// M1: 948003	      1236 ns/op	       0 B/op	       0 allocs/op
	// M1: 930661	      1258 ns/op	       0 B/op	       0 allocs/op
	if skip_write_to_file {
		b.SkipNow()
	}
	out, _ := os.Create("./test-al.log")
	al := alog.New(ext.NewFilterWriter(out, alog.TraceLevel, 0))
	al.Control.Level = alog.TraceLevel
	al.Flag = alog.WithLevel

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		switch i % 5 {
		case 0:
			al.Trace().Int("count", i).Write("")
		case 1:
			al.Debug().Int("count", i).Write("")
		case 2:
			al.Info().Int("count", i).Write("")
		case 3:
			al.Warn().Int("count", i).Write("")
		case 4:
			al.Error().Int("count", i).Write("")
			//case 5:
			//	al.Fatal(fail).Int("count", i).write("")
		}
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
