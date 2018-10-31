package gowl

import (
	"net"
	"net/http"
	"strings"
)

const (
	HEAD    = http.MethodHead
	GET     = http.MethodGet
	POST    = http.MethodPost
	PUT     = http.MethodPut
	PATCH   = http.MethodPatch
	DELETE  = http.MethodDelete
	OPTIONS = http.MethodOptions
	TRACE   = http.MethodTrace
	CONNECT = http.MethodConnect
)

type Request struct {
	*http.Request

	server *server
	params StringMap

	Data Data
}

func (r *Request) Param(name string) string {
	return r.params.Get(name)
}

func (r *Request) Params() StringMap {
	return r.params.Copy()
}

func (r *Request) ClientIP() (ip string) {
	if ip = r.Header.Get("X-Forwarded-For"); ip != "" {
		if p := strings.IndexByte(ip, ','); p != -1 {
			ip = ip[:p]
		}
		if ip = strings.TrimSpace(ip); ip != "" {
			return
		}
	}
	if ip = r.Header.Get("X-Real-IP"); ip != "" {
		return
	}
	ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	return
}

func (r *Request) GetURL(name string, params StringMap, absolute bool) string {
	url := r.server.router.url(name, params)
	if absolute {
		url.Scheme = r.URL.Scheme
		url.Host = r.Host
	}
	return url.String()
}

func (r *Request) TemplateName() string {
	if path := r.params.Get(":route"); path != "" {
		return strings.Replace(path, ".", "/", -1) + r.server.config.TemplateFileExt
	}
	return ""
}
