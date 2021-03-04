package alog

type kvType uint8

const (
	KvInt kvType = iota + 1
	KvFloat64
	KvString
	KvBool
	KvError
)

// KeyValue holds Key and value info
type KeyValue struct {
	Key   string
	Vtype kvType
	Vint  int64
	Vf64  float64
	Vstr  string
	Vbool bool
	Verr  error
}
