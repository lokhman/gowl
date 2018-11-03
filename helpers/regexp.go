package helpers

import (
	"regexp"
)

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
