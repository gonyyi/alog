package alog

// control determines if it should be printed as a log or not.
// Target speed: single 170, multi 79, noLog 5
type control struct {
	TagBucket *TagBucket // this is 1032 bytes, better to be used as a pointer
	Fn        ControlFn
	Level     Level
	Tag       Tag
}

// Check will check if level and tag given is good to be printed.
func (c control) Check(lvl Level, tag Tag) bool {
	if c.Level <= lvl || c.Tag&tag != 0 {
		return true
	}
	return false
}

// CheckFn will check if level and tag given is good to be printed.
func (c control) CheckFn(lvl Level, tag Tag) bool {
	if c.Fn != nil {
		return c.Fn(lvl, tag)
	}
	return false
}
