package alog

var Default = defaults{}

type defaults struct{}

func (d defaults) Format() Format {
	return Fdefault
}

func (d defaults) NewFmtJSON() fmtJSON {
	return fmtJSON{}
}

func (d defaults) NewFmtText() fmtText {
	return fmtText{}
}
