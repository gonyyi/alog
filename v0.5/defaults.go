package alog

import "io"

var Defaults defaults

type defaults struct{}

func (d defaults) BufferSingle() *bufBasic {
	return &bufBasic{}
}

func (d defaults) BufferSyncPool() *bufSyncPool {
	return &bufSyncPool{}
}

func (d defaults) ConverterSimple() *formatterConvBasic {
	return &formatterConvBasic{}
}

func (d defaults) FormatterJSON() *formatterJSON {
	fj := formatterJSON{
		conv: &formatterConvBasic{},
		esc:  &formatterEscBasic{},
	}
	fj.esc.Init()
	return &fj
}

func (d defaults) EscaperKey() *formatterEscBasic {
	return &formatterEscBasic{}
}

func (d defaults) ToAlWriter(w io.Writer) AlWriter {
	return iowToAlw(w)
}
