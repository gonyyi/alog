package alog

type kvType uint8

const (
	kvInt kvType = iota + 1
	kvFloat64
	kvString
	kvBool
)

// KeyVal holds key and value info
type KeyVal struct {
	k   string
	t   kvType
	i   int64
	f64 float64
	s   string
	b   bool
}

func Kint(key string, val int64) KeyVal       { return KeyVal{k: key, t: kvInt, i: val} }
func Kfloat64(key string, val float64) KeyVal { return KeyVal{k: key, t: kvFloat64, f64: val} }
func Kstr(key string, val string) KeyVal      { return KeyVal{k: key, t: kvString, s: val} }
func Kbool(key string, val bool) KeyVal       { return KeyVal{k: key, t: kvBool, b: val} }
func Kerror(key string, err error) KeyVal {
	if err != nil {
		return KeyVal{k: key, t: kvString, s: err.Error()}
	}
	return KeyVal{k: key, t: kvString, s: ""}
}
