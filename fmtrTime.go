package alog

import "time"

type FormatterDateTime interface {
	// LogTime adds millisecond level time (hhmmss000); JSON key: "t"
	LogTime(dst []byte, t time.Time) []byte
	// LogTimeDate adds CCYYMMDD to buffer; JSON key: "d"
	LogTimeDate(dst []byte, t time.Time) []byte
	// LogTimeDay adds weekday to buffer;
	// JSON: "wd":3 (wednesday)
	// TEXT: "Wed"
	LogTimeDay(dst []byte, t time.Time) []byte
}
