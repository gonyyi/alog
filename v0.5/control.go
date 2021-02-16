package alog

type control struct {
	Tags   TagBucket
	hook   HookFn
	Filter hFilters
}

func (c *control) SetHook(h HookFn) {
	c.hook = h
}

// HookFn is a type for a function designed to run when certain condition meets
type HookFn func(lvl Level, tag Tag, p []byte)

// FilterFn is a function type to be used with SetFilter.
type FilterFn func(Level, Tag) bool

type hFilters struct {
	fn    FilterFn
	level Level
	tag   Tag
}

// SetFilter will define what level or tags to show.
// Integer 0 can be used, and when it's used, it will not Filter anything.
func (f *hFilters) Set(lv Level, tags Tag) {
	f.level = lv
	f.tag = tags
}

func (f *hFilters) SetFn(fn FilterFn) {
	// didn't check for nil, because if it's nil, it will simple remove current one.
	f.fn = fn
}

// check will check if Level and Tag given is good to be printed.
func (f *hFilters) Check(lvl Level, tag Tag) bool {
	switch {
	case f.fn != nil: // FilterFn has the highest order if Set.
		return f.fn(lvl, tag)
	case f.level > lvl: // if wLevel is below wLevel limit, the do not print
		return false
	case f.tag != 0 && f.tag&tag == 0: // if filterTag is Set but Tag is not matching, then do not print
		return false
	default:
		return true
	}
}
