package httputil

import (
	"net/http"
	"strings"

	"github.com/lokhman/gowl/types"
)

const asciiSpaceSet = " \t\r\n"

func CopyHeader(dst, src http.Header) {
	for k, v := range src {
		dst[k] = v
	}
}

// HeaderValue
type HeaderValue struct {
	Value  string
	Params types.StringMap
}

func (v HeaderValue) String() string {
	buf := new(strings.Builder)
	buf.WriteString(v.Value)
	if len(v.Params) > 0 {
		buf.WriteByte(';')
		for key, value := range v.Params {
			buf.WriteString(key)
			buf.Write([]byte{'=', '"'})
			buf.WriteString(value)
			buf.WriteByte('"')
		}
	}
	return buf.String()
}

// HeaderValues
type HeaderValues []HeaderValue

func (v HeaderValues) String() string {
	buf := new(strings.Builder)
	for i, value := range v {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(value.String())
	}
	return buf.String()
}

func ParseHeaderMultiValues(header http.Header, key string) HeaderValues {
	values := make(HeaderValues, 0)
	count := 0

	key = http.CanonicalHeaderKey(key)
	for _, value := range header[key] {
		isValue := true
		lastIndex := 0
		tokenStart := -1

		value += ","
		for i := 0; skipSpace(value, &i); i++ {
			if value[i] == ',' || value[i] == ';' {
				// skip empty parts
				if tokenStart == -1 {
					continue
				}

				// extract token
				token := value[tokenStart : lastIndex+1]
				tokenStart = -1

				if isValue {
					// append new header value
					values = append(values, HeaderValue{
						Value:  token,
						Params: make(types.StringMap),
					})
					count++
				} else {
					// parse value parameters
					params := values[count-1].Params
					if pos := strings.IndexByte(token, '='); pos != -1 {
						k := strings.TrimRight(token[:pos], asciiSpaceSet)
						v := strings.TrimLeft(token[pos+1:], asciiSpaceSet)

						// unquote parameter value
						if n := len(v); n > 1 && v[0] == '"' && v[n-1] == '"' {
							v = v[1 : n-1]
						}
						params.Set(k, v)
					} else {
						params.Set(token, "")
					}
				}
				isValue = value[i] == ','
			} else if tokenStart == -1 {
				tokenStart = i
			}
			lastIndex = i
		}
	}
	return values
}

// ...
func skipSpace(s string, i *int) bool {
	n := len(s)
	for *i < n {
		if strings.IndexByte(asciiSpaceSet, s[*i]) == -1 {
			break
		}
		*i++
	}
	return *i < n
}
