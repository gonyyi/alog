// (c) 2020 Gon Y Yi. <https://gonyyi.com>
// Version 0.2.0, 12/31/2020

package alog_test

import (
	"github.com/gonyyi/alog"
	"strconv"
	"testing"
)

var (
	tStr   = "hello this is gon"
	tInt   = 123
	tInt64 = 123
	tFloat = 1.234
	tBool  = true
)

func BenchmarkTest(b *testing.B) {
	var buf []byte
	s := "test\n\"escapes\""

	b.Run("byte by byte", func(b2 *testing.B) {
		for i := 0; i < b2.N; i++ {
			buf = buf[:0]
			for j := 0; j < len(s); j++ {
				switch s[j] {
				case '\\':
					buf = append(buf, '\\')
				case '"':
					buf = append(buf, '\\', '"')
				case '\n':
					buf = append(buf, '\\', 'n')
				case '\t':
					buf = append(buf, '\\', 't')
				case '\r':
					buf = append(buf, '\\', 'r')
				case '\b':
					buf = append(buf, '\\', 'b')
				case '\f':
					buf = append(buf, '\\', 'f')
				default:
					buf = append(buf, s[j])
				}
			}
		}
	})
	println("byte2byte", string(buf))

	b.Run("strconv", func(b2 *testing.B) {
		for i := 0; i < b2.N; i++ {
			buf = buf[:0]
			buf = strconv.AppendQuote(buf, s)
		}
	})
	println("strconv", string(buf))
}

func BenchmarkLogger_Info(b *testing.B) {
	al := alog.New(nil)
	{
		b.Run("fmt=json,arg=", func(b2 *testing.B) {
			b2.ReportAllocs()
			al.SetFormat(alog.Flevel | alog.Ftag | alog.Fjson)
			for i := 0; i < b2.N; i++ {
				al.Info(0, "test")
			}
		})

		b.Run("fmt=json,arg=int", func(b2 *testing.B) {
			b2.ReportAllocs()
			al.SetFormat(alog.Flevel | alog.Ftag | alog.Fjson)
			for i := 0; i < b2.N; i++ {
				al.Info(0, "test", "a1", 123)
			}
		})

		b.Run("fmt=json,arg=str", func(b2 *testing.B) {
			b2.ReportAllocs()
			al.SetFormat(alog.Flevel | alog.Ftag | alog.Fjson)
			for i := 0; i < b2.N; i++ {
				al.Info(0, "test", "a1", "str1")
			}
		})
		b.Run("fmt=json,arg=str+int", func(b2 *testing.B) {
			b2.ReportAllocs()
			al.SetFormat(alog.Flevel | alog.Ftag | alog.Fjson)
			for i := 0; i < b2.N; i++ {
				al.Info(0, "test", "a1", "str1", "b1", 123)
			}
		})
	}

	{
		b.Run("fmt=text,arg=", func(b2 *testing.B) {
			b2.ReportAllocs()
			al.SetFormat(alog.Flevel | alog.Ftag)
			for i := 0; i < b2.N; i++ {
				al.Info(0, "test")
			}
		})

		b.Run("fmt=text,arg=int", func(b2 *testing.B) {
			b2.ReportAllocs()
			al.SetFormat(alog.Flevel | alog.Ftag)
			for i := 0; i < b2.N; i++ {
				al.Info(0, "test", "a2", 123)
			}
		})
		b.Run("fmt=text,arg=str", func(b2 *testing.B) {
			b2.ReportAllocs()
			al.SetFormat(alog.Flevel | alog.Ftag)
			for i := 0; i < b2.N; i++ {
				al.Info(0, "test", "a1", "str1")
			}
		})
		b.Run("fmt=text,arg=str+int", func(b2 *testing.B) {
			b2.ReportAllocs()
			al.SetFormat(alog.Flevel | alog.Ftag)
			for i := 0; i < b2.N; i++ {
				al.Info(0, "test", "a1", "str1", "a2", 123)
			}
		})
	}
}
