package benchmark

import (
	"github.com/gonyyi/alog"
	"github.com/gonyyi/alog/ext"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"
)

// Benchmark Summary:
// - Ran a benchmark against two most popular loggers out there. Golang's built-in logger and Zerolog.
// - Removed the Zerolog's output as it doesn't really have plain text format. (except fancy colored one)
//
// Hardware:   2020 Mac Mini (M1 processor) with 8 GB ram.
// Go version: go1.16 darwin/arm64
//
// Versions
// 	  Alog    0.5.0
// 	  Zerolog 1.20.0
// 	  Builtin 1.16.0
//
// Alog (Extension: Text)
// Benchmark_Compare3Loggers_Text/al-0-8         	 3545655	       320.2 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Compare3Loggers_Text/al-1-8         	 3736027	       317.8 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Compare3Loggers_Text/al-2-8         	 3740014	       317.4 ns/op	       0 B/op	       0 allocs/op
//
// Builtin (Text)
// Benchmark_Compare3Loggers_Text/gl-0-8         	 4939922	       238.3 ns/op	      16 B/op	       1 allocs/op
// Benchmark_Compare3Loggers_Text/gl-1-8         	 4910499	       235.5 ns/op	      16 B/op	       1 allocs/op
// Benchmark_Compare3Loggers_Text/gl-2-8         	 4926660	       241.7 ns/op	      16 B/op	       1 allocs/op
//
// 2021/02/18 20:03:03 INF alog test // name="alog"
// {"level":"info","name":"zlog","time":"2021-02-18T20:03:03-06:00","message":"zerolog"}
// 2021/02/18 20:03:03 defaultnamedef
//

func Benchmark_Compare3Loggers_Text(b *testing.B) {
	al := alog.New(nil).SetFormat(alog.Fdate | alog.Ftime | alog.Flevel).SetFormatter(ext.FormatterText())
	zl := zlog.Output(zerolog.ConsoleWriter{Out: ioutil.Discard})
	gl := log.New(ioutil.Discard, "", log.Ltime|log.Ldate)

	for rnd := 0; rnd < 3; rnd++ {
		b.Run("al-"+strconv.Itoa(rnd), func(b2 *testing.B) {
			b2.ReportAllocs()
			for i := 0; i < b2.N; i++ {
				al.Info(0, "alog test", "name", "alog")
			}
		})

		b.Run("zl-"+strconv.Itoa(rnd), func(b2 *testing.B) {
			b2.ReportAllocs()
			for i := 0; i < b2.N; i++ {
				zl.Info().Str("name", "zlog").Msg("zerolog")
			}
		})

		b.Run("gl-"+strconv.Itoa(rnd), func(b2 *testing.B) {
			b2.ReportAllocs()
			for i := 0; i < b2.N; i++ {
				gl.Print("default", "name", "def")
			}
		})
	}

	al.SetOutput(os.Stderr)
	zl = zl.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	gl.SetOutput(os.Stderr)
	al.Info(0, "alog test", "name", "alog")
	zl.Info().Str("name", "zlog").Msg("zerolog")
	gl.Print("default", "name", "def")
}

// Benchmark Summary:
// - Ran a benchmark against two most popular loggers out there. Golang's built-in logger and Zerolog.
// - Removed the Zerolog's output as it doesn't really have plain text format. (except fancy colored one)
//
// Hardware:   2020 Mac Mini (M1 processor) with 8 GB ram.
// Go version: go1.16 darwin/arm64
//
// Versions
// 	  Alog    0.5.0
// 	  Zerolog 1.20.0
// 	  Builtin 1.16.0
//
// Alog (ANSI color text)
// Benchmark_Compare3Loggers_TextColor/al-0-8         	 3539659	       329.8 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Compare3Loggers_TextColor/al-1-8         	 3646382	       325.6 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Compare3Loggers_TextColor/al-2-8         	 3646386	       326.6 ns/op	       0 B/op	       0 allocs/op
//
// Zerolog (ANSI color text)
// Benchmark_Compare3Loggers_TextColor/zl-0-8         	  350596	      3391 ns/op	    1953 B/op	      47 allocs/op
// Benchmark_Compare3Loggers_TextColor/zl-1-8         	  352564	      3410 ns/op	    1953 B/op	      47 allocs/op
// Benchmark_Compare3Loggers_TextColor/zl-2-8         	  351139	      3417 ns/op	    1953 B/op	      47 allocs/op
//
// Builtin (Text)
// Benchmark_Compare3Loggers_TextColor/gl-0-8         	 4995476	       235.4 ns/op	      16 B/op	       1 allocs/op
// Benchmark_Compare3Loggers_TextColor/gl-1-8         	 5006557	       238.0 ns/op	      16 B/op	       1 allocs/op
// Benchmark_Compare3Loggers_TextColor/gl-2-8         	 4994456	       230.7 ns/op	      16 B/op	       1 allocs/op
//
func Benchmark_Compare3Loggers_TextColor(b *testing.B) {
	al := alog.New(nil).SetFormat(alog.Fdate | alog.Ftime | alog.Flevel).SetFormatter(ext.FormatterTextColor())
	zl := zlog.Output(zerolog.ConsoleWriter{Out: ioutil.Discard})
	gl := log.New(ioutil.Discard, "", log.Ltime|log.Ldate)

	for rnd := 0; rnd < 3; rnd++ {
		b.Run("al-"+strconv.Itoa(rnd), func(b2 *testing.B) {
			b2.ReportAllocs()
			for i := 0; i < b2.N; i++ {
				al.Info(0, "alog test", "name", "alog")
			}
		})

		b.Run("zl-"+strconv.Itoa(rnd), func(b2 *testing.B) {
			b2.ReportAllocs()
			for i := 0; i < b2.N; i++ {
				zl.Info().Str("name", "zlog").Msg("zerolog")
			}
		})

		b.Run("gl-"+strconv.Itoa(rnd), func(b2 *testing.B) {
			b2.ReportAllocs()
			for i := 0; i < b2.N; i++ {
				gl.Print("default", "name", "def")
			}
		})
	}

	al.SetOutput(os.Stderr)
	zl = zl.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	gl.SetOutput(os.Stderr)
	al.Info(0, "alog test", "name", "alog")
	zl.Info().Str("name", "zlog").Msg("zerolog")
	gl.Print("default", "name", "def")
}
