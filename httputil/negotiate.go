package httputil

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/lokhman/gowl/types"
)

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

		params := make(types.StringMap)
		for k, v := range headerValue.Params {
			if k == "q" {
				v, err := strconv.ParseFloat(v, 32)
				if err == nil && v >= 0 && v <= 1 {
					values[i].Weight = float32(v)
				}
			} else {
				params.Set(k, v)
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
