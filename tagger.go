package alog

// TODO: separate TAG part from logger.

// tagger
type tagger struct {
	filter       filters
	numTagIssued int        // numTagIssued stores number of Tag issued.
	tagNames     [64]string // tagNames stores Tag names.
}

// SetNewTags will initialize the wTag.
// Initialized wTag(s) can be retreated by GetTag(string) or MustGetTag(string)
func (t *tagger) newTags(names ...string) {
	for _, name := range names {
		t.mustGetTag(name)
	}
}

// getTag returns a tag if found
func (t *tagger) getTag(name string) (tag Tag, ok bool) {
	for i := 0; i < t.numTagIssued; i++ {
		if t.tagNames[i] == name {
			return 1 << i, true
		}
	}
	return 0, false
}

// mustGetTag returns a tag if found. If not, create a new tag.
func (t *tagger) mustGetTag(name string) Tag {
	if tag, ok := t.getTag(name); ok {
		return tag
	}

	// create a new tag if tag not found
	t.tagNames[t.numTagIssued] = name
	tag := t.numTagIssued // this is the value to be printed.

	t.numTagIssued += 1
	return 1 << tag
}

type filters struct {
	fn  func(Level, Tag) bool
	lvl Level
	tag Tag
}

// SetFilter will define what level or tags to show.
// Integer 0 can be used, and when it's used, it will not filter anything.
func (f *filters) set(lv Level, tags Tag) {
	f.lvl = lv
	f.tag = tags
}

func (f *filters) setFn(fn FilterFn) {
	// didn't check for nil, because if it's nil, it will simple remove current one.
	f.fn = fn
}

// check will check if Level and Tag given is good to be printed.
func (f *filters) check(lvl Level, tag Tag) bool {
	switch {
	case f.fn != nil: // filterFn has the highest order if set.
		return f.fn(lvl, tag)
	case f.lvl > lvl: // if wLevel is below wLevel limit, the do not print
		return false
	case f.tag != 0 && f.tag&tag == 0: // if filterTag is set but Tag is not matching, then do not print
		return false
	default:
		return true
	}
}
