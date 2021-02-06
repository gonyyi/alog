package alog

import "time"

// Fmtr is an interface for Alog formatter.
type Fmtr interface {
	// Begin is used for beginning of log entry. This method takes prefix as well.
	// If it is set, it will be the first message to print regardless of JSON or text.
	// For JSON, it'd be '{' and none for text.
	Begin(dst []byte, prefix []byte) []byte
	// End is used for end of log entry. For JSON, add `}\n`; for text, add '\n'
	End(dst []byte) []byte
	// Space adds delimiter; eg. for JSON: "," and for TEXT: " "
	Space(dst []byte) []byte
	// LogLevel adds level to buffeer
	LogLevel(dst []byte, lv Level) []byte
	// LogTime adds millisecond level time (hhmmss000); JSON key: "t"
	LogTime(dst []byte, t time.Time) []byte
	// LogTimeDate adds CCYYMMDD to buffer; JSON key: "d"
	LogTimeDate(dst []byte, t time.Time) []byte
	// LogTimeDay adds weekday to buffer;
	// JSON: "wd":3 (wednesday)
	// TEXT: "Wed"
	LogTimeDay(dst []byte, t time.Time) []byte
	// LogTimeUnix adds unix time in second; JSON key: "ts"
	LogTimeUnix(dst []byte, t time.Time) []byte
	// LogTimeUnixMs adds unix time in millisecond precision; JSON key: "ts"; key is identical to LogTimeUnix
	LogTimeUnixMs(dst []byte, t time.Time) []byte
	// LogTag adds wTag. When empty it will show "[]"
	LogTag(dst []byte, tag Tag, alogTagStr *[64]string, alogTagIssued int) []byte
	// LogMsg is a main message; JSON key: "msg"
	LogMsg(dst []byte, s string, suffix byte) []byte
	// LogMsgb is a main message in byte slice
	LogMsgb(dst []byte, b []byte, suffix byte) []byte
	// LogCustomHeader should be used as a user-definable function.
	LogCustomHeader(dst []byte) []byte
	// Nil adds Nil/Null
	Nil(dst []byte, k string) []byte
	// Error adds error
	Error(dst []byte, k string, v error) []byte
	// Errors adds slice of errors
	Errors(dst []byte, k string, v *[]error) []byte

	// Bool adds bool to buffer
	Bool(dst []byte, k string, v bool) []byte
	// String adds string to buffer
	String(dst []byte, k string, v string) []byte
	// Int adds int value to buffer
	Int(dst []byte, k string, v int) []byte
	// Int8 adds int8 value to buffer
	Int8(dst []byte, k string, v int8) []byte
	// Int16 adds int16 value to buffer
	Int16(dst []byte, k string, v int16) []byte
	// Int32 adds int32 value to buffer
	Int32(dst []byte, k string, v int32) []byte
	// Int64 adds int64 value to buffer
	Int64(dst []byte, k string, v int64) []byte
	// Uint adds uint value to buffer
	Uint(dst []byte, k string, v uint) []byte
	// Uint8 adds uint8 value to buffer
	Uint8(dst []byte, k string, v uint8) []byte
	// Uint16 adds uint16 value to buffer
	Uint16(dst []byte, k string, v uint16) []byte
	// Uint32 adds uint32 value to buffer
	Uint32(dst []byte, k string, v uint32) []byte
	// Uint64 adds uint64 value to buffer
	Uint64(dst []byte, k string, v uint64) []byte
	// Float32 adds float32 value to buffer
	Float32(dst []byte, k string, v float32) []byte
	// Float64 adds float64 value to buffer
	Float64(dst []byte, k string, v float64) []byte

	// Bools add slice of bool values to buffer
	Bools(dst []byte, k string, v *[]bool) []byte
	// Strings add slice of strings values to buffer
	Strings(dst []byte, k string, v *[]string) []byte
	// Ints add slice of ints values to buffer
	Ints(dst []byte, k string, v *[]int) []byte
	// Int32s add slice of int32s values to buffer
	Int32s(dst []byte, k string, v *[]int32) []byte
	// Int64s add slice of int64s values to buffer
	Int64s(dst []byte, k string, v *[]int64) []byte
	// Uints add slice of uints values to buffer
	Uints(dst []byte, k string, v *[]uint) []byte
	// Uint8s add slice of uint8s values to buffer
	Uint8s(dst []byte, k string, v *[]uint8) []byte
	// Uint32s add slice of uint32s values to buffer
	Uint32s(dst []byte, k string, v *[]uint32) []byte
	// Uint64s add slice of uint64s values to buffer
	Uint64s(dst []byte, k string, v *[]uint64) []byte
	// Float32s add slice of float32s values to buffer
	Float32s(dst []byte, k string, v *[]float32) []byte
	// Float64s add slice of float64s values to buffer
	Float64s(dst []byte, k string, v *[]float64) []byte
}

// TimeFmtr is used to format the time.
// As default time format for alog is very simple,
// a user can create his/her own time formatter
// using this interface.
type TimeFmtr interface {
	// LogTime adds millisecond level time (hhmmss000); JSON key: "t"
	LogTime(dst []byte, t time.Time) []byte
	// LogTimeDate adds CCYYMMDD to buffer; JSON key: "d"
	LogTimeDate(dst []byte, t time.Time) []byte
	// LogTimeDay adds weekday to buffer; JSON key: "wd"
	LogTimeDay(dst []byte, t time.Time) []byte
}
