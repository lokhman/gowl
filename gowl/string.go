package gowl

import (
	"regexp"
	"strings"
	"unicode"
)

// StringMap
type StringMap map[string]string

func (sm StringMap) Set(key string, value string) {
	sm[key] = value
}

func (sm StringMap) Get(key string) string {
	return sm[key]
}

func (sm StringMap) Has(key string) bool {
	_, ok := sm[key]
	return ok
}

func (sm StringMap) Delete(key string) {
	delete(sm, key)
}

func (sm StringMap) Copy() StringMap {
	smc := make(StringMap, len(sm))
	for key, value := range smc {
		smc[key] = value
	}
	return smc
}

// ...
func StringContainsUpperCase(s string) bool {
	for _, c := range s {
		if unicode.IsUpper(c) {
			return true
		}
	}
	return false
}

func StringInSlice(s string, slice []string) bool {
	return StringInSliceIndex(s, slice) != -1
}

func StringInSliceIndex(s string, slice []string) int {
	for i, str := range slice {
		if s == str {
			return i
		}
	}
	return -1
}

func ToUnderscore(s string) string {
	buf := new(strings.Builder)
	var r rune
	for i, c := range s {
		if unicode.IsUpper(c) {
			if i > 0 && r != '_' && r != '.' && !unicode.IsUpper(r) {
				buf.WriteByte('_')
			}
			buf.WriteRune(unicode.ToLower(c))
		} else if c != '_' || r != '_' {
			buf.WriteRune(c)
		}
		if i > 0 {
			r = c
		}
	}
	return buf.String()
}

func ReplaceAllStringSubmatchFunc(re *regexp.Regexp, s string, repl func(match []string, i int) string) string {
	lastIndex, result := 0, ""
	for _, v := range re.FindAllStringSubmatchIndex(s, -1) {
		n := len(v)
		match := make([]string, n/2)
		for i := 0; i < n; i += 2 {
			if v[i] != -1 && v[i+1] != -1 {
				match[i/2] = s[v[i]:v[i+1]]
			}
		}
		result += s[lastIndex:v[0]] + repl(match, v[0])
		lastIndex = v[1]
	}
	return result + s[lastIndex:]
}
