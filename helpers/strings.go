package helpers

import (
	"strings"
	"unicode"
)

func IndexString(s string, slice []string) int {
	for i, str := range slice {
		if s == str {
			return i
		}
	}
	return -1
}

func IndexUpper(s string) int {
	for i, c := range s {
		if unicode.IsUpper(c) {
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
