package alog

import "io"

// devNull is a type for discard
type devNull int

// discard is defined here to get rid of needs to import of ioutil package.
var discard io.Writer = devNull(0)

// Write discards everything
func (devNull) Write([]byte) (int, error) {
	return 0, nil
}

// HookFn is a type for a function designed to run when certain condition meets
type HookFn func(lvl Level, tag Tag, msg []byte)

// FilterFn is a function type to be used with SetFilter.
type FilterFn func(Level, Tag) bool
