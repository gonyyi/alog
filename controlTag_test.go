package alog_test

import (
	"github.com/gonyyi/alog"
	"strconv"
	"testing"
)

func TestTagBucket_GetTag(t *testing.T) {
	b := alog.TagBucket{}
	if tag, ok := b.GetTag("GON"); ok || tag != 0 {
		// This should return not okay with 0 value for tag
		t.Errorf("TagBucket.GetTag() 1")
	}
	tagGon := b.MustGetTag("GON")

	if tag, ok := b.GetTag("GON"); !ok {
		// This should return not okay with 0 value for tag
		t.Errorf("TagBucket.GetTag() 2")
	} else {
		if tag == 0 {
			t.Errorf("TagBucket.GetTag() 3")
		}
		if tag != tagGon {
			t.Errorf("TagBucket.GetTag() 4")
		}
	}
}

func TestTagBucket_MustGetTag(t *testing.T) {
	{
		b := alog.TagBucket{}
		tag := b.MustGetTag("ABC")
		if tag == 0 {
			t.Errorf("TagBucket.MustGetTag() 1")
		}
		tag2 := b.MustGetTag("ABC2")
		if tag2 == 0 {
			t.Errorf("TagBucket.MustGetTag() 2")
		}
		if tag == tag2 {
			t.Errorf("TagBucket.MustGetTag() 3")
		}

		tag3 := b.MustGetTag("ABC")
		tag4 := b.MustGetTag("ABC2")
		if tag3 == 0 || tag4 == 0 {
			t.Errorf("TagBucket.MustGetTag() 4")
		}

		if tag != tag3 || tag2 != tag4 || tag3 == tag4 {
			t.Errorf("TagBucket.MustGetTag() 5")
		}
	}

	{
		b := alog.TagBucket{}
		var tags []alog.Tag
		for i := 0; i < 70; i++ {
			tags = append(tags, b.MustGetTag("tag"+strconv.Itoa(i)))
		}

		for i := 0; i < len(tags); i++ {
			if i < 64 {
				if tags[i] == 0 {
					t.Errorf("TagBucket.MustGetTag() 6 // %d: %d", i, tags[i])
				}
			} else {
				if tags[i] != 0 {
					t.Errorf("TagBucket.MustGetTag() 7 // %d: %d", i, tags[i])
				}
			}
		}
	}
}

func TestTagBucket_AppendTag(t *testing.T) {
	tb := alog.TagBucket{}
	var out []byte
	out = tb.AppendTag(out, 0)
	if string(out) != "" {
		t.Errorf("TagBucket.AppendTag() 1 // out=%s", string(out))
	}

	out = out[:0]
	out = tb.AppendTagForJSON(out, 0)
	if string(out) != "" {
		t.Errorf("TagBucket.AppendTag() 2 // out=%s", string(out))
	}

	gon := tb.MustGetTag("GON")
	jon := tb.MustGetTag("JON")
	_, _ = gon, jon
	out = out[:0]
	out = tb.AppendTag(out, gon)
	if string(out) != "GON" {
		t.Errorf("TagBucket.AppendTag() 3 // out=%s", string(out))
	}

	out = out[:0]
	out = tb.AppendTag(out, gon|jon)
	if string(out) != "GON,JON" {
		t.Errorf("TagBucket.AppendTag() 4 // out=%s", string(out))
	}

	out = out[:0]
	out = tb.AppendTagForJSON(out, gon|jon)
	if string(out) != `"GON","JON"` {
		t.Errorf("TagBucket.AppendTag() 5 // out=%s", string(out))
	}

	out = out[:0]
	out = tb.AppendTagForJSON(out, jon)
	if string(out) != `"JON"` {
		t.Errorf("TagBucket.AppendTag() 6 // out=%s", string(out))
	}
}
