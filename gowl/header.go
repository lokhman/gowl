package gowl

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const asciiSpaceSet = " \t\r\n"

// HeaderValue
type HeaderValue struct {
	Value  string
	Params StringMap
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
						Params: make(StringMap),
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
						params[k] = v
					} else {
						params[token] = ""
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

// AcceptHeaderValue
type AcceptHeaderValue struct {
	HeaderValue
	Weight float32
}

// AcceptHeaderValues
type AcceptHeaderValues []AcceptHeaderValue

var acceptHeaderValuesCompare = []func(p, q *AcceptHeaderValue) bool{
	func(p, q *AcceptHeaderValue) bool { // weight
		return p.Weight > q.Weight
	},
	func(p, q *AcceptHeaderValue) bool { // wildcard `*/*` (`*` for "Accept-*")
		return (p.Value != "*/*" && p.Value != "*") && (q.Value == "*/*" || q.Value == "*")
	},
	func(p, q *AcceptHeaderValue) bool { // suffix (ignored for "Accept-*")
		return !strings.HasSuffix(p.Value, "/*") && strings.HasSuffix(q.Value, "/*")
	},
}

func (v AcceptHeaderValues) Len() int {
	return len(v)
}

func (v AcceptHeaderValues) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v AcceptHeaderValues) Less(i, j int) bool {
	p, q := &v[i], &v[j]
	var k int
	for k = 0; k < len(acceptHeaderValuesCompare)-1; k++ {
		less := acceptHeaderValuesCompare[k]
		switch {
		case less(p, q):
			return true
		case less(q, p):
			return false
		}
	}
	return acceptHeaderValuesCompare[k](p, q)
}

func ParseAcceptHeader(header http.Header, key string) AcceptHeaderValues {
	headerValues := ParseHeaderMultiValues(header, key)
	values := make(AcceptHeaderValues, len(headerValues))
	for i, headerValue := range headerValues {
		values[i].Value = headerValue.Value
		values[i].Weight = 1.0

		params := make(StringMap)
		for k, v := range headerValue.Params {
			if k == "q" {
				v, err := strconv.ParseFloat(v, 32)
				if err == nil && v >= 0 && v <= 1 {
					values[i].Weight = float32(v)
				}
			} else {
				params[k] = v
			}
		}
		values[i].Params = params
	}
	sort.Stable(values)
	return values
}

func NegotiateAcceptHeader(header http.Header, key string, offers []string) string {
	values := ParseAcceptHeader(header, key)
	for _, value := range values {
		for _, offer := range offers {
			if value.Value == offer {
				return offer
			}
			if op := strings.IndexByte(offer, '/'); op != -1 {
				if value.Value == offer[:op+1]+"*" {
					return offer
				}
			}
			if value.Value == "*/*" || value.Value == "*" {
				return offer
			}
		}
	}
	return ""
}

func NegotiateRequestListener(offers []string, defaultToFirstOffer bool) func(event EventInterface) {
	return func(event EventInterface) {
		if len(offers) == 0 {
			return
		}

		ev := event.(*RequestEvent)
		request := ev.Request()

		offer := NegotiateAcceptHeader(request.Header, "Accept", offers)
		if offer == "" && !defaultToFirstOffer {
			link := make(HeaderValues, len(offers))
			for i, offer := range offers {
				link[i] = HeaderValue{
					Value:  "<" + request.URL.String() + ">",
					Params: StringMap{"type": offer},
				}
			}

			response := ErrorResponse(http.StatusNotAcceptable, "")
			response.Header().Add("Link", link.String())
			ev.SetResponse(response)
			return
		} else if defaultToFirstOffer {
			offer = offers[0]
		}

		request.params[":accept"] = offer
	}
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

func copyHeader(dst, src http.Header) {
	for k, v := range src {
		dst[k] = v
	}
}
