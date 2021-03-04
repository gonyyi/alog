package main

import (
	"github.com/gonyyi/alog"
	"os"
	"runtime/pprof"
)

func main() {
}
func profile() {
	out, err := os.Create("cpu.pprof")
	if err != nil {
		println(err.Error())
	}
	pprof.StartCPUProfile(out)
	defer pprof.StopCPUProfile()

	{
		al := alog.New(nil)
		for i := 0; i < 1_000_000; i++ {
			al.Error(0).Writes("hello")
		}
	}

}
