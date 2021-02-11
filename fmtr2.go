package alog

import "time"

// (*Writer) WriteHeader(bufHeader []byte, Flag, Level, &Tag, string) (bufHeader []byte)
// (*Writer) WriteBody(bufBody []byte, msg string, a ...interface{}) (bufBody []byte)
type Fmtr2 interface {
	SetHook(HookFn)
	SetTagger(*Tagger)
	SetFormat(Format)
	SetPrefix(string)
	Format(bufHead, bufBody []byte, lv Level, tag Tag, msg string, a ...interface{}) ([]byte, []byte)
	Formatb(bufHead, bufBody []byte, lv Level, tag Tag, msg []byte) ([]byte, []byte)
}

func NewFmtr2JSON() *Fmtr2JSON {
	return &Fmtr2JSON{
		conv: conv,
	}
}

type Fmtr2JSON struct {
	hook   HookFn
	tagger *Tagger
	format Format
	prefix string
	time   time.Time
	conv   Converter
}

// TODO: Converter as an interface??

// TODO: create Escaper interface
// 		KeyValue(dst []byte, key, value string, join, suffix byte, quote bool)
//		String(dst []byte, s string, suffix byte, quote bool, ignoreSpcChar bool)
//		Byte(dst []byte, b []byte, suffix byte, quote bool, ignoreSpcChar bool)

func (f *Fmtr2JSON) Format(bufHead, bufBody []byte, lv Level, tag Tag, msg string, a ...interface{}) ([]byte, []byte) {
	if f.format&Fprefix != 0 {
		bufHead = append(bufHead, f.prefix...)
	}
	bufHead = append(bufHead, '{')

	if f.format&(FtimeUnix|FtimeUnixMs) != 0 {
		f.time = time.Now()
		if f.format&FtimeUnixMs != 0 { // MS
			// s.bufHeader = l.fmtr.LogTimeUnixMs(s.bufHeader, f.time)
		} else { // Just Unix
			// s.bufHeader = l.fmtr.LogTimeUnix(s.bufHeader, f.time)
		}
	} else if f.format&(Fdate|FdateDay|Ftime|FtimeUTC) != 0 {
		// at least one item will be printed here, so just check once.
		f.time = time.Now()
		if f.format&FtimeUTC != 0 {
			f.time = f.time.UTC() // todo: is this running all the time? or need just once?
		}

		if f.format&Fdate != 0 {
			// firstItem = false
			// s.bufHeader = l.fmtr.LogTimeDate(s.bufHeader, l.time)
		}
		if f.format&FdateDay != 0 {
			// if !firstItem {
			//	s.bufHeader = l.fmtr.Space(s.bufHeader)
			// }
			// firstItem = false
			// s.bufHeader = l.fmtr.LogTimeDay(s.bufHeader, l.time)
		}
		if f.format&Ftime != 0 {
			// if !firstItem {
			//	s.bufHeader = l.fmtr.Space(s.bufHeader)
			// }
			// s.bufHeader = l.fmtr.LogTime(s.bufHeader, l.time)
		}
		// firstItem = false
	}

	if f.format&Flevel != 0 {
		bufHead = append(bufHead, `"level":"`...)
		bufHead = append(bufHead, lv.String()...)
		bufHead = append(bufHead, '"', ',')
	}

	if f.format&Ftag != 0 {
		bufHead = append(bufHead, `"tag":[`...)
		bufHead = f.tagger.AppendTagNames(bufHead, ',', true, tag)
		bufHead = append(bufHead, ']', ',')
	}

	if len(msg) > 0 {
		bufBody = append(bufBody, `"msg":`...)
		// bufBody = append(bufBody, msg...) // 22.49 ns
		// bufBody = keyable.String(bufBody, msg) // 51.43 ns
		bufBody = escapeString(bufBody, msg, true, ',') // 64.25 ns, updated one at 48.57 ns
	}

	if f.hook != nil {
		f.hook(lv, tag, nil)
	}

	if len(bufBody) > 0 && bufBody[len(bufBody)-1] == ',' {
		bufBody = bufBody[0 : len(bufBody)-1]
	}
	// Ending

	return bufHead, append(bufBody, '}')
}

func (f *Fmtr2JSON) Formatb(bufHead, bufBody []byte, lv Level, tag Tag, msg []byte) ([]byte, []byte) {
	return nil, nil
}

func (f *Fmtr2JSON) SetHook(fn HookFn) {
	f.hook = fn
}

func (f *Fmtr2JSON) SetTagger(t *Tagger) {
	f.tagger = t
}

func (f *Fmtr2JSON) SetFormat(fmt Format) {
	f.format = fmt
}

func (f *Fmtr2JSON) SetPrefix(prefix string) {
	f.prefix = prefix
}
