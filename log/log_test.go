package log_test

import (
	"github.com/gonyyi/alog/log"
	"testing"
)

func TestLog(t *testing.T) {
	log.Info(0).Write("Hello")
}
