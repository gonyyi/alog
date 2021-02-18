package alog

// will hold default objects
var Defaults *defaults = newDefaults()

func newDefaults() *defaults {
	d := defaults{}
	return &d
}

type defaults struct {
	converter Converter
	formatter Formatter
}

func (d *defaults) Converter() Converter {
	if d.converter == nil {
		d.converter = &convert{}
		d.converter.Init()
	}
	return d.converter
}

func (d *defaults) Formatter() Formatter {
	if d.formatter == nil {
		d.formatter = &formatJSON{}
		d.formatter.Init()
	}
	return d.formatter
}
