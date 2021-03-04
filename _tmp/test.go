package main

import (
	"github.com/gonyyi/alog"
	"os"
)

func main() {

	al := alog.New(nil)
	al = al.SetOutput(os.Stderr)

	al.Info(0).Str("gonSaid", "hello").Write("")

}
