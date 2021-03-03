package alog_test

import (
	"errors"
	"github.com/gonyyi/alog"
	"os"
	"testing"
)

func BenchmarkNew(b *testing.B) {
	al := alog.New(os.Stderr)
	e1 := errors.New("error msg my")
	al.Info(0).Err("err1", nil).Err("err2", e1).Str("ok", "yes okay").Write("log starts")
	al.SetOutput(nil)

	//b.Run("al-kv", func(c *testing.B) {
	//	c.ReportAllocs()
	//	for i := 0; i < c.N; i++ {
	//		al.Err(0).Err("err", nil).Write("err msg here")
	//		//al.Log(alog.Lerror, 0, "msg for err")
	//	}
	//})
}
