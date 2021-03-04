package alog

type control struct {
	TagBucket *TagBucket
	Hook      func(Level, Tag, []byte)
	Fn        func(Level, Tag) bool
	Level     Level
	Tag       Tag
}

// check will check if level and tag given is good to be printed.
func (c *control) Check(lvl Level, tag Tag) bool {
	if c.Level <= lvl || c.Tag&tag != 0 {
		return true
	}
	return false
}

func (c *control) CheckFn(lvl Level, tag Tag) bool {
	if c.Fn != nil {
		return c.Fn(lvl, tag)
	}
	return false
}
