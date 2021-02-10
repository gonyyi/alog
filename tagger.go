package alog

// Tagger
type Tagger struct {
	Filter       filters
	numTagIssued int        // numTagIssued stores number of Tag issued.
	tagNames     [64]string // tagNames stores Tag names.
}

func (t *Tagger) AppendTagNames(dst []byte, delimiter byte, quote bool, tag Tag) []byte {
	for i := 0; i < t.numTagIssued; i++ {
		if tag&(1<<i) != 0 {
			if quote { // redundant; as speed matter rather than the binary size
				dst = append(dst, '"')
				dst = append(dst, t.tagNames[i]...)
				dst = append(dst, '"', delimiter)
			} else {
				dst = append(dst, t.tagNames[i]...)
				dst = append(dst, delimiter)
			}
		}
	}
	return dst[:len(dst)-1]
}

// SetNewTags will initialize the wTag.
// Initialized wTag(s) can be retreated by GetTag(string) or MustGetTag(string)
func (t *Tagger) NewTags(names ...string) {
	for _, name := range names {
		t.MustGetTag(name)
	}
}

// GetTag returns a tag if found
func (t *Tagger) GetTag(name string) (tag Tag, ok bool) {
	for i := 0; i < t.numTagIssued; i++ {
		if t.tagNames[i] == name {
			return 1 << i, true
		}
	}
	return 0, false
}

// MustGetTag returns a tag if found. If not, create a new tag.
func (t *Tagger) MustGetTag(name string) Tag {
	if tag, ok := t.GetTag(name); ok {
		return tag
	}

	// create a new tag if tag not found
	t.tagNames[t.numTagIssued] = name
	tag := t.numTagIssued // this is the value to be printed.

	t.numTagIssued += 1
	return 1 << tag
}

type filters struct {
	fn  FilterFn
	lvl Level
	tag Tag
}

// SetFilter will define what level or tags to show.
// Integer 0 can be used, and when it's used, it will not Filter anything.
func (f *filters) Set(lv Level, tags Tag) {
	f.lvl = lv
	f.tag = tags
}

func (f *filters) SetFn(fn FilterFn) {
	// didn't check for nil, because if it's nil, it will simple remove current one.
	f.fn = fn
}

// check will check if Level and Tag given is good to be printed.
func (f *filters) check(lvl Level, tag Tag) bool {
	switch {
	case f.fn != nil: // FilterFn has the highest order if Set.
		return f.fn(lvl, tag)
	case f.lvl > lvl: // if wLevel is below wLevel limit, the do not print
		return false
	case f.tag != 0 && f.tag&tag == 0: // if filterTag is Set but Tag is not matching, then do not print
		return false
	default:
		return true
	}
}
