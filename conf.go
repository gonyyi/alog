package alog

// Conf is a global variable holding basic configuration.
// By setting Conf before creating a logger will let the
// Conf's value as its default.
var Conf *conf = func() *conf {
	d := conf{
		BufferHead:   1024,
		BufferBody:   2048,
		FormatFlag:   Fdate | Ftime | Flevel | Ftag,
		ControlLevel: Linfo,
	}
	return &d
}()

type conf struct {
	converter Converter
	formatter Formatter
	buffer    Buffer

	BufferHead   int
	BufferBody   int
	ControlLevel Level
	FormatFlag   Format
}

func (d *conf) SetBuffer(buf Buffer) {
	d.buffer = buf
}
func (d *conf) Buffer() Buffer {
	if d.buffer == nil {
		d.buffer = &bufSyncPool{}
	}
	return d.buffer
}
func (d *conf) SetConverter(conv Converter) {
	d.converter = conv
}
func (d *conf) Converter() Converter {
	if d.converter == nil {
		d.converter = &convert{}
		d.converter.Init()
	}
	return d.converter
}
func (d *conf) SetFormatter(fmtr Formatter) {
	d.formatter = fmtr
}
func (d *conf) Formatter() Formatter {
	if d.formatter == nil {
		d.formatter = &formatJSON{}
		d.formatter.Init()
	}
	return d.formatter
}
