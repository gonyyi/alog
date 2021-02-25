package alog

type kvType uint8

const (
	kvInt kvType = iota + 1
	kvInt64
	kvUint
	kvUint64
	kvFloat32
	kvFloat64
	kvString
	kvBool
	kvError
)

// KeyVal holds key and value info
type KeyVal struct {
	k   string
	t   kvType
	i   int
	i64 int64
	u   uint
	u64 uint64
	f32 float32
	f64 float64
	s   string
	b   bool
	e   error
}

type keyvals struct{}

func (keyvals) Vint(key string, val int) KeyVal     { return KeyVal{k: key, t: kvInt, i: val} }
func (keyvals) Vi64(key string, val int64) KeyVal   { return KeyVal{k: key, t: kvInt64, i64: val} }
func (keyvals) Vuint(key string, val uint) KeyVal   { return KeyVal{k: key, t: kvUint, u: val} }
func (keyvals) Vu64(key string, val uint64) KeyVal  { return KeyVal{k: key, t: kvUint64, u64: val} }
func (keyvals) Vf32(key string, val float32) KeyVal { return KeyVal{k: key, t: kvFloat32, f32: val} }
func (keyvals) Vf64(key string, val float64) KeyVal { return KeyVal{k: key, t: kvFloat64, f64: val} }
func (keyvals) Vstr(key string, val string) KeyVal  { return KeyVal{k: key, t: kvString, s: val} }
func (keyvals) Vbool(key string, val bool) KeyVal   { return KeyVal{k: key, t: kvBool, b: val} }
func (keyvals) Verr(key string, val error) KeyVal   { return KeyVal{k: key, t: kvError, e: val} }
