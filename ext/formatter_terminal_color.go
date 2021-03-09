package ext

import (
	"github.com/gonyyi/alog"
	"io"
	"strconv"
	"time"
)

// NewFormatterTerminalColor will take io.Writer and returns
// a Logger.Formatter object.
func NewFormatterTerminalColor() *fmtTxtColor {
	return &fmtTxtColor{}
}

const (
	fcCLEAR  = "\033[0m"
	fcDIM    = "\033[0;90m"
	fcBOLD   = "\033[0;1m"
	fcITALIC = "\033[0;1;3m"

	fcTRACE = "\033[100;37m"
	fcDEBUG = "\033[102;90m"
	fcINFO  = "\033[104;90m"
	fcWARN  = "\033[1;103;90m"
	fcERROR = "\033[1;101;90m"
	fcFATAL = "\033[1;105;90m"
)

type fmtTxtColor struct {
	out       io.Writer
	format    alog.Flag
	tagBucket *alog.TagBucket
}

func (f *fmtTxtColor) Init(w io.Writer, formatFlag alog.Flag, tagBucket *alog.TagBucket) {
	f.out = w
	if w == nil {
		f.out = io.Discard
	}

	f.format = formatFlag
	f.tagBucket = tagBucket
}

func (f *fmtTxtColor) Write(dst []byte) (int, error) {
	return f.out.Write(dst)
}

func (f *fmtTxtColor) Close() error {
	if c, ok := f.out.(io.Closer); ok && c != nil {
		return c.Close()
	}
	return nil
}

func (fmtTxtColor) Begin(dst []byte) []byte {
	return dst
}

func (f *fmtTxtColor) AddTime(dst []byte) []byte {
	if (alog.WithUnixTime|alog.WithDate|alog.WithTime)&f.format != 0 {
		return append(append(dst, time.Now().Format(fcDIM+"2006-0102 "+fcCLEAR+ "15:04:05")...), ' ')
	}
	return dst
}

func (fmtTxtColor) AddLevel(dst []byte, level alog.Level) []byte {

	switch level {
	case alog.TraceLevel:
		dst = append(dst, fcTRACE...)
	case alog.DebugLevel:
		dst = append(dst, fcDEBUG...)
	case alog.InfoLevel:
		dst = append(dst, fcINFO...)
	case alog.WarnLevel:
		dst = append(dst, fcWARN...)
	case alog.ErrorLevel:
		dst = append(dst, fcERROR...)
	case alog.FatalLevel:
		dst = append(dst, fcFATAL...)
	}
	return append(append(append(append(dst, ' '), level.NameShort()...), ' '), fcCLEAR+" "...)
}

func (f *fmtTxtColor) AddTag(dst []byte, tag alog.Tag) []byte {
	return append(f.tagBucket.AppendTag(append(dst, fcDIM+"["+fcCLEAR...), tag), fcDIM+"]"+fcCLEAR+" "...)
}

func (fmtTxtColor) AddMsg(dst []byte, s string) []byte {
	if s != "" {
		return append(append(dst, s...), ` // `...)
	}
	return dst
}

func (f *fmtTxtColor) AddKVs(dst []byte, kvs []alog.KeyValue) []byte {
	for i := 0; i < len(kvs); i++ {
		dst = append(append(append(dst, fcDIM...), kvs[i].Key...), "="+fcCLEAR...)
		switch kvs[i].Vtype {
		case alog.KvString:
			dst = f.addValStringUnsafe(dst, kvs[i].Vstr)
		case alog.KvBool:
			dst = f.addValBool(dst, kvs[i].Vbool)
		case alog.KvError:
			if kvs[i].Verr == nil {
				dst = append(dst, "null, "...)
			} else {
				dst = f.addValString(dst, kvs[i].Verr.Error())
			}
		case alog.KvInt:
			dst = f.addValInt(dst, kvs[i].Vint)
		case alog.KvFloat64:
			dst = f.addValFloat(dst, kvs[i].Vf64)
		default:
			dst = append(dst, `null, `...)
		}
	}
	return dst
}

func (fmtTxtColor) End(dst []byte) []byte {
	if len(dst) > 1 {
		dst[len(dst)-2] = '\n'
		return dst[:len(dst)-1]
	}
	return dst
}

func (fmtTxtColor) addKeyUnsafe(dst []byte, s string) []byte {
	return append(append(dst, s...), '=')
}

func (fmtTxtColor) addValString(dst []byte, s string) []byte {
	return append(strconv.AppendQuote(dst, s), ',', ' ')
}

func (fmtTxtColor) addValStringUnsafe(dst []byte, s string) []byte {
	return append(append(append(dst, '"'), s...), `", `...)
}

func (fmtTxtColor) addValBool(dst []byte, b bool) []byte {
	if b {
		return append(dst, `true, `...)
	}
	return append(dst, `false, `...)
}

func (fmtTxtColor) addValInt(dst []byte, i int64) []byte {
	return append(strconv.AppendInt(dst, i, 10), ',', ' ')
}

func (fmtTxtColor) addValFloat(dst []byte, f float64) []byte {
	return append(strconv.AppendFloat(dst, f, 'f', -1, 64), ',', ' ')
}
