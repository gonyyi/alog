package alog

type Formatter interface {
	AppendPrefix(dst []byte, prefix []byte) []byte
	AppendTime(dst []byte, format Format) []byte
	AppendTag(dst []byte, tb *TagBucket, tag Tag) []byte
	AppendMsg(dst []byte, s string) []byte
	AppendMsgBytes(dst []byte, p []byte) []byte
	AppendAdd(dst []byte, a ...interface{}) []byte
	AppendSuffix(dst []byte, suffix []byte) []byte
	TrimLast(dst []byte, b byte) []byte
}
