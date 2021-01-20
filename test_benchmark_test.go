// (c) 2020 Gon Y Yi. <https://gonyyi.com>
// Version 0.2.0, 12/31/2020

package alog_test

var (
	tStr   = "hello this is gon"
	tInt   = 123
	tInt64 = 123
	tFloat = 1.234
	tBool  = true
)

// func BenchmarkLogger_Info(b *testing.B) {
// 	b.Run("fmt=json, error=str", func(b2 *testing.B) {
// 		l := alog.New(nil).SetFormat(alog.Fjson)
// 		b2.ReportAllocs()
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Error(testLogMsg)
// 			}
// 		})
// 	})
// 	b.Run("fmt=json+utc, error=str", func(b2 *testing.B) {
// 		l := alog.New(nil).SetFormat(alog.Fjson | alog.FtimeUTC)
// 		b2.ReportAllocs()
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Error(testLogMsg)
// 			}
// 		})
// 	})
// 	b.Run("fmt=json+time, error=str", func(b2 *testing.B) {
// 		l := alog.New(nil).SetFormat(alog.Fjson | alog.Ftime)
// 		b2.ReportAllocs()
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Error(testLogMsg)
// 			}
// 		})
// 	})
// 	b.Run("fmt=json+time, errorf=str", func(b2 *testing.B) {
// 		l := alog.New(nil).SetFormat(alog.Fjson | alog.Ftime)
// 		b2.ReportAllocs()
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Errorf("test: %s", "abc")
// 			}
// 		})
// 	})
// 	b.Run("fmt=json+time, errorf=int", func(b2 *testing.B) {
// 		l := alog.New(nil).SetFormat(alog.Fjson | alog.Ftime)
// 		b2.ReportAllocs()
// 		b2.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				// prints none
// 				l.Errorf("test: %d", 1234)
// 			}
// 		})
// 	})
// }
