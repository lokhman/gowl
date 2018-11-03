package gowl

import (
	"html/template"
	"net"
	"net/http"
	"strings"

	"github.com/lokhman/gowl/types"
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
	params types.StringMap

	Data types.Data
}

func (r *Request) Param(name string) string {
	return r.params.Get(name)
}

func (r *Request) Params() types.StringMap {
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

func (r *Request) GetURL(name string, params types.StringMap, absolute bool) string {
	url := r.server.router.url(name, params)
	if absolute {
		url.Scheme = r.URL.Scheme
		url.Host = r.Host
	}
	return url.String()
}

func (r *Request) Template() *template.Template {
	if path := r.params.Get(":route"); path != "" {
		name := strings.Replace(path, ".", "/", -1)
		return r.server.templates[name+r.server.config.TemplateFileExt]
	}
	return nil
}
