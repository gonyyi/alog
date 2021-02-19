package benchmark

import (
	"github.com/gonyyi/alog"
	zlog "github.com/rs/zerolog/log"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"
)

// Benchmark Summary:
// Ran a benchmark against two most popular loggers out there. Golang's built-in logger and Zerolog.
//
// Hardware:   2020 Mac Mini (M1 processor) with 8 GB ram.
// Go version: go1.16 darwin/arm64
//
// Versions
// 	  Alog    0.5.0
// 	  Zerolog 1.20.0
// 	  Builtin 1.16.0
//
// Alog (JSON)
// Benchmark_Compare3Loggers/al-0-8         	 4247889	       271.2 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Compare3Loggers/al-1-8         	 4410723	       268.3 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Compare3Loggers/al-2-8         	 4404769	       268.2 ns/op	       0 B/op	       0 allocs/op
//
// Zerolog (JSON)
// Benchmark_Compare3Loggers/zl-0-8         	 3938845	       301.5 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Compare3Loggers/zl-1-8         	 3934148	       301.6 ns/op	       0 B/op	       0 allocs/op
// Benchmark_Compare3Loggers/zl-2-8         	 3922243	       301.5 ns/op	       0 B/op	       0 allocs/op
//
// Builtin (Text)
// Benchmark_Compare3Loggers/gl-0-8         	 4991648	       238.9 ns/op	      16 B/op	       1 allocs/op
// Benchmark_Compare3Loggers/gl-1-8         	 5002394	       238.8 ns/op	      16 B/op	       1 allocs/op
// Benchmark_Compare3Loggers/gl-2-8         	 4994046	       239.7 ns/op	      16 B/op	       1 allocs/op
//
// {"d":20210218,"t":195446,"lv":"info","msg":"alog test","name":"alog"}
// {"level":"info","name":"zlog","time":"2021-02-18T19:54:46-06:00","message":"zerolog"}
// 2021/02/18 19:54:46 defaultnamedef

func Benchmark_Compare3Loggers(b *testing.B) {
	al := alog.New(nil).SetFormat(alog.Fdate | alog.Ftime | alog.Flevel)
	zl := zlog.Output(nil)

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
	zl = zl.Output(os.Stderr)
	gl.SetOutput(os.Stderr)
	al.Info(0, "alog test", "name", "alog")
	zl.Info().Str("name", "zlog").Msg("zerolog")
	gl.Print("default", "name", "def")
}
