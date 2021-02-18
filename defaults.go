package alog

// will hold default objects
var Defaults *defaults = newDefaults()

func newDefaults() *defaults {
	d := defaults{}
	d.Converter.Init()
	d.FormatterJSON = &formatJSON{}
	d.FormatterText = &formatText{}
	return &d
}

type defaults struct {
	Converter     convert
	FormatterJSON Formatter
	FormatterText Formatter
}
