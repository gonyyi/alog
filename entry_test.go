package alog_test

import (
	"encoding/json"
	"errors"
	"github.com/gonyyi/alog"
	"testing"
)

func TestEntry(t *testing.T) {
	tag1 = log.NewTag("TAG1")
	tag2 = log.NewTag("TAG2")
	log.Info(tag1|tag2).Str("hello", "gon").
		Bool("isSingle", true).
		Float("height", 5.8).
		Int("age", 50).
		Int64("age2", int64(50)).
		Err("myErr", errors.New("testMyErr")).
		Write("done")
	check(t, `{"level":"info","tag":["TAG1","TAG2"],"message":"done","hello":"gon","isSingle":true,"height":5.8,"age":50,"age2":50,"myErr":"testMyErr"}`)

	// UseTime
	{
		a := struct {
			Time    int    `json:"time"`
			TS      int    `json:"ts"`
			Date    int    `json:"date"`
			Day     int    `json:"day"`
			Message string `json:"message"`
		}{}
		resetA := func() {
			out.Reset()
			a = struct {
				Time    int    `json:"time"`
				TS      int    `json:"ts"`
				Date    int    `json:"date"`
				Day     int    `json:"day"`
				Message string `json:"message"`
			}{}
		}

		{
			log.Flag = alog.UseTime
			log.Info(0).Write("done")
			json.Unmarshal(out.Bytes(), &a)
			if a.Time < 1 {
				t.Errorf("TestEntry: flag UseTime 1")
			}
			if a.Message != "done" {
				t.Errorf("TestEntry: flag UseTime 2")
			}
			out.Reset()
		}

		{
			resetA()
			log.Flag = alog.UseUnixTimeMs
			log.Info(0).Write("done2")
			json.Unmarshal(out.Bytes(), &a)
			if a.TS < 1000 || a.Message != "done2" {
				t.Errorf("TestEntry: flag UseTime 3")
			}
		}

		{
			resetA()
			log.Flag = alog.UseUnixTime
			log.Info(0).Write("done3")
			json.Unmarshal(out.Bytes(), &a)
			if a.TS < 1000 || a.Message != "done3" {
				t.Errorf("TestEntry: flag UseTime 4 // TS: <%d>, MSG: <%s>", a.TS, a.Message)
			}
		}

		{
			resetA()
			log.Flag = alog.UseDate | alog.UseDay | alog.UseUTC
			log.Info(0).Write("done4")
			json.Unmarshal(out.Bytes(), &a)
			if a.TS != 0 || a.Message != "done4" || a.Date < 1 || a.Day > 6 {
				t.Errorf("TestEntry: flag UseTime 5 // TS: <%d>, MSG: <%s>, Date: <%d>, Day: <%d>", a.TS, a.Message, a.Date, a.Day)
			}
		}

		{
			resetA()
			log.Flag = alog.UseTimeMs
			log.Info(0).Str("k", "v\t").
				Err("er1", nil).
				Err("er", errors.New("a\tb")).
				Bool("b\t1", true).Bool("b2", false).
				Write("done5\ta")
			json.Unmarshal(out.Bytes(), &a)
			if a.TS != 0 || a.Message != "done5\ta" || a.Date != 0 || a.Day != 0 ||
				a.Time < 1 {
				t.Errorf("TestEntry: flag UseTime 5 // TS: <%d>, MSG: <%s>, Date: <%d>, Day: <%d> Time: <%d>", a.TS, a.Message, a.Date, a.Day, a.Time)
			}
		}

		out.Reset()
		check(t, ``)
	}
}

type fakeData struct {
	Name      string
	City      string
	State     string
	Postal    string
	Lat       float64
	Lon       float64
	Age       int
	IsCurrent bool
}

func fakeEntryFn(d fakeData) alog.EntryFn {
	return func(entry *alog.Entry) *alog.Entry {
		if entry == nil {
			return entry
		}
		return entry.Str("name", d.Name).
			Str("city", d.City).
			Str("state", d.State).
			Str("postal", d.Postal).
			Float("lat", d.Lat).
			Float("lon", d.Lon).
			Int("age", d.Age).
			Bool("isCurrent", d.IsCurrent)
	}
}

func TestEntry_Fn(t *testing.T) {
	data := fakeData{
		Name:      "Jon",
		City:      "Goncity",
		State:     "Gonstate",
		Postal:    "12345-1234",
		Lat:       5.10000001,
		Lon:       -5.20000002,
		Age:       50,
		IsCurrent: false,
	}
	reset()
	log.Info(0).Ext(fakeEntryFn(data)).Write("added fake data")
	check(t, `{"level":"info","tag":[],"message":"added fake data","name":"Jon","city":"Goncity","state":"Gonstate","postal":"12345-1234","lat":5.10000001,"lon":-5.20000002,"age":50,"isCurrent":false}`)
}
