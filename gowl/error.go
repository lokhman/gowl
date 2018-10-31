package gowl

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func getMainStackTrace(stack errors.StackTrace) errors.StackTrace {
	start := 0
	for i, frame := range stack {
		fn := runtime.FuncForPC(uintptr(frame) - 1)
		if fn == nil {
			continue
		}
		if strings.HasPrefix(fn.Name(), "main.") {
			start = i
			break
		}
	}
	return stack[start:]
}

func NewErrorResponse(statusCode int, serverName, debug string) ResponseInterface {
	response := NewResponse(statusCode, func(w io.Writer) error {
		return errorTemplate.Execute(w, StringMap{
			"name":   fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
			"server": serverName,
			"debug":  debug,
		})
	})
	response.header.Set("Content-Type", "text/html; charset=utf-8")
	response.header.Set("X-Content-Type-Options", "nosniff")
	return response
}

var errorTemplate = template.Must(template.New("error").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>{{.name}}</title>
</head>
<body style="font-family: sans-serif">
    <h1 style="font-size: x-large">{{.name}}</h1>
	{{if .debug -}}
		<pre>{{.debug}}</pre>
	{{- end}}
	{{if .server -}}
		<hr style="border-style: outset">
		<small>{{.server}}</small>
	{{- end}}
</body>
</html>`))
