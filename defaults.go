package alog

// will hold default objects
var Defaults *defaults = newDefaults()

func newDefaults() *defaults {
	d := defaults{}
	return &d
}

type defaults struct {
	converter     Converter
	formatterText Formatter
}

func (d *defaults) Converter() Converter {
	if d.converter == nil {
		d.converter = &convert{}
		d.converter.Init()
	}
	return d.converter
}

func (d *defaults) FormatterText() Formatter {
	if d.formatterText == nil {
		d.formatterText = &formatText{}
		d.formatterText.Init()
	}
	return d.formatterText
}
