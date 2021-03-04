package alog

import "io"

// Formatter is an interface for a combination of formatter and writer.
type Formatter interface {
	// Init will initialize or update the formatter setting.
	// This will be run by Alog when used.
	Init(writer io.Writer, format Flag, tagBucket TagBucket)

	// Begin will be used for formats requiring prefix such as `{` in JSON.
	Begin([]byte) []byte

	// AddTime will add time to buffer.
	AddTime([]byte) []byte

	// AddLevel will add level to buffer.
	AddLevel([]byte, Level) []byte

	// AddTag will add a tag to buffer.
	AddTag([]byte, Tag) []byte

	// AddMsg will add default messages to buffer.
	AddMsg([]byte, string) []byte

	// AddKVs will add key value items to buffer
	AddKVs([]byte, []KeyValue) []byte

	// End will be used as suffix such as `}` and/or newline in JSON.
	End([]byte) []byte

	// Write will let the formatter write the buffer.
	// Depend on the setting of formatter, this can be
	// totally different than the writer than the logger's.
	// Using Write, AddLevel (and AddTag),
	// this Formatter interface will be used to create
	// multi output logger such as for a syslog.
	Write([]byte) (int, error)

	// Close will, if applicable, close the writer.
	Close() error
}
