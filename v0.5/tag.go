package alog

// Tag is a bit-formatFlag used to show only necessary part of process to show
// in the log. For instance, if there's an web service, there can be different
// Tag such as UI, HTTP request, HTTP response, etc. By alConf a Tag
// for each log using `Print` or `Printf`, a user can only print certain
// Tag of log messages for better debugging.
type Tag uint64

func (tag Tag) Has(t Tag) bool {
	if tag&t != 0 {
		return true
	}
	return false
}

type TagBucket struct {
	count int        // count stores number of Tag issued.
	names [64]string // names stores Tag names.
}

// GetTag returns a tag if found
func (t *TagBucket) GetTag(name string) (tag Tag, ok bool) {
	for i := 0; i < t.count; i++ {
		if t.names[i] == name {
			return 1 << i, true
		}
	}
	return 0, false
}

// MustGetTag returns a tag if found. If not, create a new tag.
func (t *TagBucket) MustGetTag(name string) Tag {
	if tag, ok := t.GetTag(name); ok {
		return tag
	}

	// create a new tag if tag not found
	t.names[t.count] = name
	tag := t.count // this is the value to be printed.

	t.count += 1
	return 1 << tag
}

func (t *TagBucket) AppendSelectedTags(dst []byte, delimiter byte, quote bool, tag Tag) []byte {
	if tag == 0 {
		return dst
	}
	cntDst := len(dst)
	for i := 0; i < t.count; i++ {
		if tag&(1<<i) != 0 {
			if quote { // redundant; as speed matter rather than the binary size
				dst = append(dst, '"')
				dst = append(dst, t.names[i]...)
				dst = append(dst, '"', delimiter)
			} else {
				dst = append(dst, t.names[i]...)
				dst = append(dst, delimiter)
			}
		}
	}
	if cntDst < len(dst) {
		return dst[:len(dst)-1] // last delimiter to be omitted
	}
	return dst
}
