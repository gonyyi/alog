package alog

// control determines if it should be printed as a log or not.
// Target speed: single 170, multi 79, noLog 5
type control struct {
	bucket *TagBucket // this is 1032 bytes, better to be used as a pointer
	Fn     ControlFn
	Level  Level
	Tags   Tag
}

func newControl() control {
	return control{
		bucket: &TagBucket{},
		Fn:     nil,
		Level:  InfoLevel,
		Tags:   0,
	}
}

// Bucket will return the pointer of TagBucket in control.
func (c control) Bucket() *TagBucket {
	return c.bucket
}

// Check will check if level and tag given is good to be printed.
func (c control) Check(lvl Level, tag Tag) bool {
	if c.Level <= lvl || c.Tags&tag != 0 {
		return true
	}
	return false
}

// CheckFn will check if level and tag given is good to be printed.
func (c control) CheckFn(lvl Level, tag Tag) (bool, bool) {
	if c.Fn != nil {
		return true, c.Fn(lvl, tag)
	}
	return false, false
}
