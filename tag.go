package alog

// Tag is a bit-flag used to show only necessary part of process to show
// in the log. For instance, if there's an web service, there can be different
// Tag such as UI, HTTP request, HTTP response, etc. By setting a Tag
// for each log using `Print` or `Printf`, a user can only print certain
// Tag of log messages for better debugging.
type Tag uint64

func (tag Tag) Has(t Tag) bool {
	if tag&t != 0 {
		return true
	}
	return false
}
