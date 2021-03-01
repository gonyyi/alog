package alog

// Tag is a bit-formatFlag used to show only necessary part of process to show
// in the log. For instance, if there's an web service, there can be different
// Tag such as UI, HTTP request, HTTP response, etc. By alConf a Tag
// for each log using `Print` or `Printf`, a user can only print certain
// Tag of log messages for better debugging.
type Tag uint64

// TagBucket can issue a tag and also holds the total number
// of tags issued AND also names given to each tag.
// Not that TagBucket is not using any mutex as it is designed
// to be set at the very beginning of the process.
// Also, the maximum number of tag can be issue is limited to 63.
type TagBucket struct {
	count int        // count stores number of Tag issued.
	names [64]string // names stores Tag names.
}

// GetTag returns a tag if found
func (t TagBucket) GetTag(name string) (tag Tag, ok bool) {
	for i := 0; i < t.count; i++ {
		if t.names[i] == name {
			return 1 << i, true
		}
	}
	return 0, false
}

// MustGetTag returns a tag if found. If not, create a new tag.
func (t *TagBucket) MustGetTag(name string) Tag {
	// If a tag is found, return it.
	if tag, ok := t.GetTag(name); ok {
		return tag
	}
	// If the tag is not found, issue a tag using most recently created.
	// When the maximum capacity of tag has met, return 0.
	if t.count >= 63 {
		return 0
	}
	// Create a new tag and return the tag.
	t.names[t.count] = name
	tag := t.count // this is the value to be printed.
	t.count += 1
	return 1 << tag
}

func (t *TagBucket) AppendTag(dst []byte, tag Tag) []byte {
	for i := 0; i < t.count; i++ {
		if tag&(1<<i) != 0 {
			dst = append(append(dst, t.names[i]...), ',')
		}
	}
	return dst
}
func (t *TagBucket) AppendTagForJSON(dst []byte, tag Tag) []byte {
	origLen := len(dst)
	for i := 0; i < t.count; i++ {
		if tag&(1<<i) != 0 {
			dst = append(append(append(dst, '"'), t.names[i]...), '"', ',')
		}
	}
	if len(dst) > origLen {
		return dst[0 : len(dst)-1]
	}
	return dst
}
