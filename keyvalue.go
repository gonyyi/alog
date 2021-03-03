package alog

type kvType uint8

const (
	kvInt kvType = iota + 1
	kvFloat64
	kvString
	kvBool
	kvError
)

// KeyVal holds key and value info
type KeyVal struct {
	key   string
	vType kvType
	vInt  int64
	vF64  float64
	vStr  string
	vBool bool
	vErr  error
}
