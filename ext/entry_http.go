package ext

import (
	"github.com/gonyyi/alog"
	"net/http"
)

var EntryHttp entryHttp

type entryHttp struct{}

// ReqRx is for Request Received
func (entryHttp) ReqRx(r *http.Request) alog.EntryFn {
	ipAddr := r.Header.Get("X-Real-Ip")
	if ipAddr == "" {
		ipAddr = r.Header.Get("X-Forwarded-For")
	}
	if ipAddr == "" {
		ipAddr = r.RemoteAddr
	}

	return func(e *alog.Entry) *alog.Entry {
		return e.Str("method", r.Method).
			Str("uri", r.RequestURI).
			Str("ip", ipAddr)
	}
}
